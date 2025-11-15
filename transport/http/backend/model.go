package backend

type PageResponse struct {
	Total *int `json:"total"` // Total number of objects available in the server
	Limit *int `json:"limit"` // Limit maximum number of objects that the server can return
}

type SearchResponse struct {
	PageResponse
	Items any `json:"items,omitempty"`
}
