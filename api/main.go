package main

import (
	"encoding/json"
	"net/http"
)

// Publication represents a specific edition of a work
type Publication struct {
	Image  string `json:"image"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var publications []Publication = []Publication{
	{
		Image:  "https://images-na.ssl-images-amazon.com/images/I/81X4R7QhFkL.jpg",
		Title:  "Normal People",
		Author: "Sally Rooney",
	},
	{
		Image:  "https://images-na.ssl-images-amazon.com/images/I/91twTG-CQ8L.jpg",
		Title:  "Little Fires Everywhere",
		Author: "Celeste Ng",
	},
	{
		Image:  "https://images-na.ssl-images-amazon.com/images/I/51j5p18mJNL.jpg",
		Title:  "Where the Crawdads Sing",
		Author: "Delia Owens",
	},
	{
		Image:  "https://images-na.ssl-images-amazon.com/images/I/81af+MCATTL.jpg",
		Title:  "The Great Gatsby",
		Author: "F. Scott Fitzgerald",
	},
	{
		Image:  "https://images-na.ssl-images-amazon.com/images/I/81iVsj91eQL.jpg",
		Title:  "American Dirt",
		Author: "Jeanine Cummins",
	},
	{
		Image:  "https://images-na.ssl-images-amazon.com/images/I/91Xq+S+F2jL.jpg",
		Title:  "Atomic Habits",
		Author: "James Clear",
	},
}

func main() {
	http.HandleFunc("/api/publications", PublicationsHandler)

	http.ListenAndServe(":8000", nil)
}

// PublicationsHandler handles requests made to /api/publications
func PublicationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		publications := getPublications()
		bytes, err := json.Marshal(publications)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

func getPublications() []Publication {
	return publications
}
