package scraper

type Job struct {
	Title       string   `json:"title"`
	Company     string   `json:"company"`
	URL         string   `json:"url"`
	Link        string   `json:"link"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
}
