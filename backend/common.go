package backend

import "time"

type Image struct {
	Url       string `json:"url,omitempty"`       // Url of the image
	Thumbnail string `json:"thumbnail,omitempty"` // Thumbnail of the image
	Category  string `json:"category,omitempty"`  // Category
	Type      string `json:"type,omitempty"`      // Type file type (jpeg, giff etc)
	Width     int    `json:"width,omitempty"`     // Width size
	Height    int    `json:"height,omitempty"`    // Height size
}

// BusinessDetails is data received by remote platform
type BusinessDetails struct {
	Name    string `json:"name,omitempty"`    // Name company name
	Website string `json:"website,omitempty"` // Website company website
	Logo    *Image `json:"logo,omitempty"`    // Logo company logo
	Inn     string `json:"inn,omitempty"`     // Inn (tax number) of the company (Yandex extension of the protocol. Isn't supported by OCPI)
}

type PageResponse struct {
	Total *int `json:"total,omitempty"` // Total number of items available by request
	Limit *int `json:"limit,omitempty"` // Limit number of retrieved items
}

type DisplayText struct {
	Language string `json:"language,omitempty"` // Language
	Text     string `json:"text,omitempty"`     // Text description
}

type PullRequest struct {
	From *time.Time `json:"from"`
	To   *time.Time `json:"to"`
}
