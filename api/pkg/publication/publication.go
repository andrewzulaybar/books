package publication

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/author"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/andrewzulaybar/books/api/pkg/work"
	"github.com/lib/pq"
)

// Columns is the comma-separated list of columns found in the publication table.
const Columns string = "edition_pub_date, format, image_url, isbn, isbn13, language, num_pages, publisher, work_id"

// Enum constants representing types of SQL statements.
const (
	Unknown postgres.Query = iota
	DeletePublication
	GetPublication
	GetPublications
	PatchPublication
	PostPublication
)

// Publication represents a specific edition of a work.
type Publication struct {
	ID             int       `json:"id"`
	EditionPubDate string    `json:"editionPubDate"`
	Format         string    `json:"format"`
	ImageURL       string    `json:"imageUrl"`
	ISBN           string    `json:"isbn"`
	ISBN13         string    `json:"isbn13"`
	Language       string    `json:"language"`
	NumPages       int       `json:"numPages"`
	Publisher      string    `json:"publisher"`
	Work           work.Work `json:"work"`
}

// Publications represents a list of publications.
type Publications []Publication

// Service wraps the database and other dependencies.
type Service struct {
	DB postgres.DB

	WorkService work.Service
}

// Query returns a SQL statement based on the postgres.Query value passed in.
func (s *Service) Query(query postgres.Query, args ...interface{}) string {
	switch query {
	case DeletePublication:
		return "DELETE FROM publication WHERE id = $1"
	case GetPublication:
		return fmt.Sprintf(
			`SELECT publication.id, %s, %s, %s
                        FROM publication
                        JOIN work ON publication.work_id=work.id
                        JOIN author ON work.author_id=author.id
                        WHERE publication.id = $1`,
			Columns,
			work.Columns,
			author.Columns,
		)
	case GetPublications:
		return fmt.Sprintf(
			`SELECT publication.id, %s, %s, %s
                        FROM publication
                        JOIN work ON publication.work_id=work.id
                        JOIN author ON work.author_id=author.id
                        ORDER BY publication.id`,
			Columns,
			work.Columns,
			author.Columns,
		)
	case PatchPublication:
		var hasUpdate bool
		query := "UPDATE publication SET "
		for column, value := range args[0].(map[string]interface{}) {
			switch value.(type) {
			case string:
				if value != "" {
					query += fmt.Sprintf("%s = '%s',", column, value)
					hasUpdate = true
				}
			case int:
				if value != 0 {
					query += fmt.Sprintf("%s = %d,", column, value)
					hasUpdate = true
				}
			}
		}
		if hasUpdate {
			return strings.TrimSuffix(query, ",") + fmt.Sprintf(" WHERE id = $1 RETURNING id, %s", Columns)
		}
		return ""
	case PostPublication:
		return fmt.Sprintf(
			`INSERT INTO publication (%s)
                        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
                        RETURNING id, %s`,
			Columns,
			Columns,
		)
	default:
		return ""
	}
}

// DeletePublication removes the entry in the publication table matching the given id.
func (s *Service) DeletePublication(id int) *status.Status {
	db := s.DB
	deletePublication := s.Query(DeletePublication)

	res, err := db.Exec(deletePublication, id)
	if err != nil {
		log.Printf("[DeletePublication] %s", err)
		return status.New(status.InternalServerError, err.Error())
	}

	numDeleted, err := res.RowsAffected()
	if err != nil {
		log.Printf("[DeletePublication] %s", err)
		return status.New(status.InternalServerError, err.Error())
	}
	if numDeleted == 0 {
		msg := fmt.Sprintf("Publication with id = %d does not exist", id)
		log.Printf("[DeletePublication] %s", msg)
		return status.New(status.OK, msg)
	}

	return status.New(status.NoContent, "")
}

// DeletePublications removes the entries in the publication table matching the given ids.
func (s *Service) DeletePublications(ids []int) (*status.Status, []int) {
	notFound := []int{}
	for _, id := range ids {
		if s := s.DeletePublication(id); s.Code() != status.NoContent {
			notFound = append(notFound, id)
		}
	}

	if len(notFound) > 0 {
		msg := fmt.Sprintf("The following publications could not be found: %v", notFound)
		log.Printf("[DeletePublications] %s", msg)
		return status.Newf(status.OK, msg), notFound
	}
	return status.New(status.NoContent, ""), nil
}

// GetPublication retrieves the publication from the database matching the given ID.
func (s *Service) GetPublication(id int) (*status.Status, *Publication) {
	db := s.DB
	getPublication := s.Query(GetPublication)

	row := db.QueryRow(getPublication, id)
	pub, err := s.getPublication(row)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := fmt.Sprintf("Publication with id = %d does not exist", id)
			log.Printf("[GetPublication] %s", msg)
			return status.Newf(status.NotFound, msg), nil
		}
		log.Printf("[GetPublication] %s", err)
		return status.New(status.InternalServerError, err.Error()), nil
	}
	return status.New(status.OK, ""), pub
}

// GetPublications retrieves the entire list of publications from the database.
func (s *Service) GetPublications() (*status.Status, Publications) {
	db := s.DB
	getPublications := s.Query(GetPublications)

	rows, err := db.Query(getPublications)
	if err != nil {
		log.Printf("[GetPublications] %s", err)
		return status.New(status.InternalServerError, err.Error()), nil
	}

	publications := Publications{}
	for rows.Next() {
		pub, err := s.getPublication(rows)
		if err != nil {
			log.Printf("[GetPublications] %s", err)
			return status.New(status.InternalServerError, err.Error()), nil
		}
		publications = append(publications, *pub)
	}
	return status.New(status.OK, ""), publications
}

// PatchPublication updates the entry in the database matching pub.id with the given attributes.
func (s *Service) PatchPublication(pub *Publication) (*status.Status, *Publication) {
	if pub.Work != (work.Work{}) {
		if s := s.handleWork(pub); s.Err() != nil {
			log.Printf("[PatchPublication] %s", s.Err())
			return status.New(s.Code(), s.Message()), nil
		}
	}

	db := s.DB
	p := map[string]interface{}{
		"edition_pub_date": pub.EditionPubDate,
		"format":           pub.Format,
		"image_url":        pub.ImageURL,
		"isbn":             pub.ISBN,
		"isbn13":           pub.ISBN13,
		"language":         pub.Language,
		"num_pages":        pub.NumPages,
		"publisher":        pub.Publisher,
		"work_id":          pub.Work.ID,
	}

	patchPublication := s.Query(PatchPublication, p)
	if patchPublication != "" {
		var pb Publication
		row := db.QueryRow(patchPublication, pub.ID)
		if err := row.Scan(
			&pb.ID, &pb.EditionPubDate, &pb.Format, &pb.ImageURL, &pb.ISBN, &pb.ISBN13,
			&pb.Language, &pb.NumPages, &pb.Publisher, &pb.Work.ID,
		); err != nil {
			if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
				log.Printf("[PatchPublication] %s", err)
				return status.New(status.Conflict, err.Error()), nil
			}
			log.Printf("[PatchPublication] %s", err)
			return status.New(status.InternalServerError, err.Error()), nil
		}
		return status.New(status.OK, ""), &pb
	}
	return status.New(status.OK, "No fields in publication to update"), nil
}

// PostPublication creates an entry in the publication table with the given attributes.
func (s *Service) PostPublication(pub *Publication) (*status.Status, *Publication) {
	if pub.ID != 0 {
		if s, l := s.GetPublication(pub.ID); s.Code() == status.OK {
			msg := fmt.Sprintf("Publication with id = %d already exists", pub.ID)
			log.Printf("[PostWork] %s", msg)
			return status.Newf(status.Conflict, msg), l
		}
	}

	if pub.Work != (work.Work{}) {
		if s := s.handleWork(pub); s.Err() != nil {
			log.Printf("[PostWork] %s", s.Err())
			return status.New(s.Code(), s.Message()), nil
		}
	}

	db := s.DB
	postPublication := s.Query(PostPublication)

	var pb Publication
	row := db.QueryRow(
		postPublication,
		pub.EditionPubDate,
		pub.Format,
		pub.ImageURL,
		pub.ISBN,
		pub.ISBN13,
		pub.Language,
		pub.NumPages,
		pub.Publisher,
		pub.Work.ID,
	)
	if err := row.Scan(
		&pb.ID,
		&pb.EditionPubDate,
		&pb.Format,
		&pb.ImageURL,
		&pb.ISBN,
		&pb.ISBN13,
		&pb.Language,
		&pb.NumPages,
		&pb.Publisher,
		&pb.Work.ID,
	); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			log.Printf("[PostPublication] %s", err)
			return status.New(status.Conflict, err.Error()), nil
		}
		log.Printf("[PostPublication] %s", err)
		return status.New(status.UnprocessableEntity, err.Error()), nil
	}
	return status.New(status.Created, ""), &pb
}

func (s *Service) getPublication(row interface {
	Scan(dest ...interface{}) error
}) (*Publication, error) {
	var p Publication
	var w *work.Work = &p.Work
	var a *author.Author = &w.Author
	err := row.Scan(
		&p.ID,
		&p.EditionPubDate,
		&p.Format,
		&p.ImageURL,
		&p.ISBN,
		&p.ISBN13,
		&p.Language,
		&p.NumPages,
		&p.Publisher,
		&w.ID,
		&w.Description,
		&w.InitialPubDate,
		&w.OriginalLanguage,
		&w.Title,
		&a.ID,
		&a.FirstName,
		&a.LastName,
		&a.Gender,
		&a.DateOfBirth,
		&a.PlaceOfBirth.ID,
	)
	return &p, err
}

func (s *Service) handleWork(pub *Publication) *status.Status {
	wk := &pub.Work

	if wk.ID != 0 {
		return status.New(status.OK, "")
	}

	stat, work := s.WorkService.PostWork(wk)
	if stat.Err() != nil {
		if stat.Code() != status.Conflict {
			return status.New(stat.Code(), stat.Message())
		}

		if stat, work = s.WorkService.GetWork(wk.ID); stat.Err() != nil {
			stat, work = s.WorkService.FindWork(wk.Title, wk.Author.ID)
			if stat.Err() != nil {
				return status.New(stat.Code(), stat.Message())
			}
		}
	}

	wk.ID = work.ID
	return status.New(status.OK, "")
}
