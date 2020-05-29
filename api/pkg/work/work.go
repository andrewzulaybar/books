package work

type Work struct {
	ID               int    `json:"id"`
	AuthorID         int    `json:"authorId"`
	Description      string `json:"description"`
	InitialPubDate   string `json:"initialPubDate"`
	OriginalLanguage string `json:"originalLanguage"`
	Title            string `json:"title"`
}

type Works []Work
