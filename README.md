# API Documentation

All the endpoints available are listed below alongside the types of methods supported.

## Publication

A publication represents a specific edition of a work.
```
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
```

- [**GET** /api/publication]: Retrieves the entire list of publications from the database.
- [**POST** /api/publication]: Creates an entry in the publication table with the given attributes.
- [**DELETE** /api/publication]: Removes the entries in the publication table matching the given ids.

- [**GET** /api/publication/:id]: Retrieves the publication from the database matching the given id.
- [**PATCH** /api/publication/:id]: Updates the entry in the database matching pub.id with the given attributes.
- [**DELETE** /api/publication/:id]: Removes the entries in the publication table matching the given ids.

## Work

A work represents a literary work.
```
type Work struct {
	ID               int           `json:"id"`
	Description      string        `json:"description"`
	InitialPubDate   string        `json:"initialPubDate"`
	OriginalLanguage string        `json:"originalLanguage"`
	Title            string        `json:"title"`
	Author           author.Author `json:"author"`
}
```

- [**GET** /api/work]: Retrieves the entire list of works from the database.
- [**POST** /api/work]: Creates an entry in the work table with the given attributes.
- [**DELETE** /api/work]: Removes the entry in the work table matching the given id.

- [**GET** /api/work/:id]: Retrieves the work from the database matching the given id.
- [**PATCH** /api/work/:id]: Updates the entry in the database matching work.id with the given attributes.
- [**DELETE** /api/work/:id]: Removes the entries in the work table matching the given ids.

## Author

An author represents a writer of a work.
```
type Author struct {
	ID           int               `json:"id"`
	FirstName    string            `json:"firstName"`
	LastName     string            `json:"lastName"`
	Gender       string            `json:"gender"`
	DateOfBirth  string            `json:"dateOfBirth"`
	PlaceOfBirth location.Location `json:"placeOfBirth"`
}
```

- [**GET** /api/author]: Retrieves the entire list of authors from the database.
- [**POST** /api/author]: Creates an entry in the author table with the given attributes.
- [**DELETE** /api/author]: Removes the entries in the author table matching the given ids.

- [**GET** /api/author/:id]: Retrieves the author from the database matching the given id.
- [**PATCH** /api/author/:id]: Updates the entry in the database matching author.id with the given attributes.
- [**DELETE** /api/author/:id]: Removes the entry in the author table matching the given id.
