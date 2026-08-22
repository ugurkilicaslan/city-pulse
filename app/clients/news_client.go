package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const gnewsBaseURL = "https://gnews.io/api/v4"

// NewsClient — GNews API HTTP istemcisi
type NewsClient struct {
	httpClient *http.Client
	apiKey     string
}

func NewNewsClient(apiKey string) *NewsClient {
	return &NewsClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiKey:     apiKey,
	}
}

// GNewsSource — GNews kaynak nesnesi
type GNewsSource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// GNewsArticle — GNews tek makale
type GNewsArticle struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Content     string      `json:"content"`
	URL         string      `json:"url"`
	Image       string      `json:"image"`
	PublishedAt string      `json:"publishedAt"`
	Source      GNewsSource `json:"source"`
}

// GNewsResponse — GNews API ham cevabı
type GNewsResponse struct {
	TotalArticles int            `json:"totalArticles"`
	Articles      []GNewsArticle `json:"articles"`
}

// GetTopHeadlines — Güncel manşetleri çeker
// lang: "tr", "en" vb. | maxItems: en fazla kaç makale
func (c *NewsClient) GetTopHeadlines(ctx context.Context, lang string, maxItems int) (*GNewsResponse, error) {
	if lang == "" {
		lang = "en"
	}
	if maxItems <= 0 {
		maxItems = 10
	}

	endpoint := fmt.Sprintf("%s/top-headlines?lang=%s&max=%d&token=%s",
		gnewsBaseURL, lang, maxItems, c.apiKey)

	return c.doRequest(ctx, endpoint)
}

// Search — Anahtar kelimeye göre haber arar
// query: arama terimi | lang: dil kodu | maxItems: makale sayısı
func (c *NewsClient) Search(ctx context.Context, query, lang string, maxItems int) (*GNewsResponse, error) {
	if lang == "" {
		lang = "en"
	}
	if maxItems <= 0 {
		maxItems = 10
	}

	endpoint := fmt.Sprintf("%s/search?q=%s&lang=%s&max=%d&token=%s",
		gnewsBaseURL, url.QueryEscape(query), lang, maxItems, c.apiKey)

	return c.doRequest(ctx, endpoint)
}

func (c *NewsClient) doRequest(ctx context.Context, endpoint string) (*GNewsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gnews isteği başarısız: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gnews API hata kodu: %d", resp.StatusCode)
	}

	var result GNewsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("gnews JSON parse hatası: %w", err)
	}

	return &result, nil
}
