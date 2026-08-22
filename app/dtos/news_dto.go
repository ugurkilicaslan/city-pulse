package dtos

// SourceDTO — Haber kaynağı özet bilgisi
type SourceDTO struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// ArticleDTO — Tek bir haber makalesi response'u
type ArticleDTO struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Image       string    `json:"image"`
	PublishedAt string    `json:"publishedAt"`
	Source      SourceDTO `json:"source"`
}

// HeadlinesResponse — GET /news/headlines ve /news/search cevabı
type HeadlinesResponse struct {
	TotalArticles int          `json:"totalArticles"`
	Articles      []ArticleDTO `json:"articles"`
}
