package services

import (
	"context"
	"sync"

	"city-pulse/app/dtos"
	"city-pulse/app/models"
)

// CityService — Aggregator Pattern: Birden fazla servisi paralel çalıştırıp tek cevap döner
type CityService struct {
	exchangeSvc *ExchangeService
	newsSvc     *NewsService
	gameSvc     *GameService
	nasaSvc     *NasaService
	githubSvc   *GithubService
}

func NewCityService(exc *ExchangeService, news *NewsService, game *GameService, nasa *NasaService, gh *GithubService) *CityService {
	return &CityService{
		exchangeSvc: exc,
		newsSvc:     news,
		gameSvc:     game,
		nasaSvc:     nasa,
		githubSvc:   gh,
	}
}

// GetSnapshot — 3 servisi goroutine ile paralel çalıştırır, hepsini bekler, birleştirir
// Bir servis hata verse bile snapshot tamamen durmuyor, o kısım boş geliyor
func (s *CityService) GetSnapshot(ctx context.Context, city, lang string) (*dtos.CitySnapshotResponse, error) {
	var (
		rates    *models.ExchangeRate
		articles []models.NewsArticle
		deals    []models.GameDeal
		apod     *dtos.ApodDTO
		repos    []dtos.GitHubRepoDTO
		wg       sync.WaitGroup
	)

	wg.Add(5)

	go func() {
		defer wg.Done()
		rates, _ = s.exchangeSvc.GetLatest(ctx, "USD", []string{"TRY", "EUR", "CHF"})
	}()

	go func() {
		defer wg.Done()
		arts, err := s.newsSvc.GetHeadlines(ctx, lang)
		if err == nil && len(arts) > 5 {
			arts = arts[:5]
		}
		articles = arts
	}()

	go func() {
		defer wg.Done()
		deals, _ = s.gameSvc.GetDeals(ctx, 0, 5)
	}()

	go func() {
		defer wg.Done()
		apod, _ = s.nasaSvc.GetApod(ctx)
	}()

	go func() {
		defer wg.Done()
		repos, _ = s.githubSvc.GetTrendingRepos(ctx, 5)
	}()

	wg.Wait()

	resp := &dtos.CitySnapshotResponse{City: city}

	if rates != nil {
		resp.Exchange = &dtos.LatestRatesResponse{
			Base:  rates.Base,
			Date:  rates.Date,
			Rates: rates.Rates,
		}
	}

	if len(articles) > 0 {
		articleDTOs := make([]dtos.ArticleDTO, len(articles))
		for i, a := range articles {
			articleDTOs[i] = dtos.ArticleDTO{
				Title:       a.Title,
				Description: a.Description,
				URL:         a.URL,
				Image:       a.Image,
				PublishedAt: a.PublishedAt,
				Source:      dtos.SourceDTO{Name: a.Source.Name, URL: a.Source.URL},
			}
		}
		resp.News = articleDTOs
	}

	if len(deals) > 0 {
		dealDTOs := make([]dtos.GameDealDTO, len(deals))
		for i, d := range deals {
			dealDTOs[i] = dtos.GameDealDTO{
				DealID:          d.DealID,
				Title:           d.Title,
				SalePrice:       d.SalePrice,
				NormalPrice:     d.NormalPrice,
				SavingsPercent:  d.Savings,
				MetacriticScore: d.MetacriticScore,
				SteamRating:     d.SteamRatingText,
				Thumb:           d.Thumb,
			}
		}
		resp.GameDeals = dealDTOs
	}

	if apod != nil {
		resp.NasaApod = apod
	}

	if len(repos) > 0 {
		resp.GithubTrend = repos
	}

	return resp, nil
}
