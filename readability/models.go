package readability

// An Article represents a Readability article object.
type Article struct {
	Author        interface{} `json:"author"`
	Content       string      `json:"content"`
	DatePublished interface{} `json:"date_published"`
	Dek           interface{} `json:"dek"`
	Direction     string      `json:"direction"`
	Domain        string      `json:"domain"`
	Excerpt       string      `json:"excerpt"`
	LeadImageURL  string      `json:"lead_image_url"`
	NextPageID    interface{} `json:"next_page_id"`
	RenderedPages int         `json:"rendered_pages"`
	ShortURL      string      `json:"short_url"`
	Title         string      `json:"title"`
	TotalPages    int         `json:"total_pages"`
	URL           string      `json:"url"`
	WordCount     int         `json:"word_count"`
}
