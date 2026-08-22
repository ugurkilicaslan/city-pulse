package services

import (
	"context"

	"city-pulse/app/clients"
	"city-pulse/app/daos"
	"city-pulse/app/models"
)

// NewsService — Haber iş mantığı
// Headlines: cache'lenir (15 dk TTL) | Search: her seferinde API'ye gider
type NewsService struct {
	dao    *daos.NewsDAO
	client *clients.NewsClient
}

func NewNewsService(dao *daos.NewsDAO, client *clients.NewsClient) *NewsService {
	return &NewsService{dao: dao, client: client}
}

// GetHeadlines — Güncel manşetleri döner (cache öncelikli)
func (s *NewsService) GetHeadlines(ctx context.Context, lang string) ([]models.NewsArticle, error) {
	if lang == "" {
		lang = "en"
	}

	cached, err := s.dao.FindByLang(ctx, lang)
	if err != nil {
		return nil, err
	}
	if len(cached) > 0 {
		return cached, nil
	}

	resp, err := s.client.GetTopHeadlines(ctx, lang, 20)
	if err != nil {
		return nil, err
	}

	articles := make([]models.NewsArticle, len(resp.Articles))
	for i, a := range resp.Articles {
		articles[i] = models.NewsArticle{
			Title:       a.Title,
			Description: a.Description,
			Content:     a.Content,
			URL:         a.URL,
			Image:       a.Image,
			PublishedAt: a.PublishedAt,
			Source: models.NewsSource{
				Name: a.Source.Name,
				URL:  a.Source.URL,
			},
		}
	}

	if err := s.dao.BulkUpsert(ctx, articles, lang); err != nil {
		return nil, err
	}

	return articles, nil
}

// Search — Anahtar kelimeye göre haber arar (cache'lenmez)
func (s *NewsService) Search(ctx context.Context, query, lang string) ([]models.NewsArticle, error) {
	if lang == "" {
		lang = "en"
	}

	resp, err := s.client.Search(ctx, query, lang, 10)
	if err != nil {
		return nil, err
	}

	articles := make([]models.NewsArticle, len(resp.Articles))
	for i, a := range resp.Articles {
		articles[i] = models.NewsArticle{
			Title:       a.Title,
			Description: a.Description,
			Content:     a.Content,
			URL:         a.URL,
			Image:       a.Image,
			PublishedAt: a.PublishedAt,
			Source: models.NewsSource{
				Name: a.Source.Name,
				URL:  a.Source.URL,
			},
		}
	}

	return articles, nil
}
