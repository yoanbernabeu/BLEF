package blef

import "time"

// BLEFDocument represents the root structure of a BLEF file
type BLEFDocument struct {
	Format      string       `json:"format"`
	Version     string       `json:"version"`
	ExportedAt  time.Time    `json:"exported_at"`
	User        *User        `json:"user,omitempty"`
	Books       []Book       `json:"books"`
	Collections []Collection `json:"collections"`
	Entries     []Entry      `json:"entries"`
}

// User represents optional user information
type User struct {
	ID       string                 `json:"id,omitempty"`
	Name     string                 `json:"name,omitempty"`
	Email    string                 `json:"email,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Book represents a unique bibliographic work
type Book struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Subtitle    string                 `json:"subtitle,omitempty"`
	Authors     []Author               `json:"authors"`
	Identifiers Identifiers            `json:"identifiers"`
	Language    string                 `json:"language,omitempty"`
	Description string                 `json:"description,omitempty"`
	CoverURL    string                 `json:"cover_url,omitempty"`
	Edition     *Edition               `json:"edition,omitempty"`
	Series      *Series                `json:"series,omitempty"`
	Subjects    []string               `json:"subjects,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Author represents a book author
type Author struct {
	Name        string            `json:"name"`
	Role        string            `json:"role,omitempty"`
	Identifiers map[string]string `json:"identifiers,omitempty"`
}

// Identifiers holds various book identifiers
type Identifiers struct {
	ISBN13      string                 `json:"isbn13,omitempty"`
	ISBN10      string                 `json:"isbn10,omitempty"`
	ASIN        string                 `json:"asin,omitempty"`
	OpenLibrary string                 `json:"openlibrary,omitempty"`
	Wikidata    string                 `json:"wikidata,omitempty"`
	Goodreads   string                 `json:"goodreads,omitempty"`
	Other       map[string]interface{} `json:"other,omitempty"`
}

// Edition represents edition information
type Edition struct {
	Publisher     string `json:"publisher,omitempty"`
	PublishedDate string `json:"published_date,omitempty"`
	Format        string `json:"format,omitempty"`
	Pages         int    `json:"pages,omitempty"`
	EditionNumber string `json:"edition_number,omitempty"`
}

// Series represents book series information
type Series struct {
	Name   string      `json:"name"`
	Volume interface{} `json:"volume,omitempty"` // can be number or string
}

// Collection represents a user's shelf/list
type Collection struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Description string                 `json:"description,omitempty"`
	IsPublic    bool                   `json:"is_public"`
	CreatedAt   *time.Time             `json:"created_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Entry links a book to user-specific data
type Entry struct {
	BookID        string                 `json:"book_id"`
	CollectionIDs []string               `json:"collection_ids"`
	UserData      UserData               `json:"user_data"`
	Ownership     *Ownership             `json:"ownership,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// UserData contains user-specific book information
type UserData struct {
	Status       string     `json:"status"`
	Rating       float64    `json:"rating,omitempty"`
	Review       string     `json:"review,omitempty"`
	PrivateNotes string     `json:"private_notes,omitempty"`
	Tags         []string   `json:"tags,omitempty"`
	Favorite     bool       `json:"favorite,omitempty"`
	ReadDates    []ReadDate `json:"read_dates,omitempty"`
	AddedAt      *time.Time `json:"added_at,omitempty"`
}

// ReadDate represents reading history
type ReadDate struct {
	Started  string `json:"started,omitempty"`
	Finished string `json:"finished,omitempty"`
	Progress int    `json:"progress,omitempty"`
}

// Ownership represents book ownership and lending
type Ownership struct {
	Owned  bool    `json:"owned,omitempty"`
	Loaned *Loaned `json:"loaned,omitempty"`
}

// Loaned represents lending information
type Loaned struct {
	Status bool   `json:"status"`
	To     string `json:"to,omitempty"`
	Date   string `json:"date,omitempty"`
	Notes  string `json:"notes,omitempty"`
}
