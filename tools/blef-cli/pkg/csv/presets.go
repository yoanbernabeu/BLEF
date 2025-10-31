package csv

// ColumnMapping defines how CSV columns map to BLEF fields
type ColumnMapping struct {
	BookID        string
	ISBN13        string
	ISBN10        string
	Title         string
	Author        string
	Language      string
	Publisher     string
	PublishedDate string
	Pages         string
	Rating        string
	Review        string
	Status        string
	DateRead      string
	DateAdded     string
	Tags          string
	Shelf         string
}
