package webhook

type Webhook struct {
	Id     string   `json:"id"`     // Id webhook id
	ApiKey string   `json:"apiKey"` // ApiKey passed by webhook
	Events []string `json:"events"` // Events on which webhook is fired
	Url    string   `json:"url"`    // Url webhook url
}

type SearchResponse struct {
	Items []*Webhook `json:"items"`
}
