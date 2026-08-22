package services

import (
	"context"
	"strings"

	"city-pulse/app/clients"
	"city-pulse/app/daos"
	"city-pulse/app/models"
)

// ExchangeService — Döviz iş mantığı
// Önce MongoDB cache'e bakar, yoksa Frankfurter API'den çeker ve cache'e yazar
type ExchangeService struct {
	dao    *daos.ExchangeDAO
	client *clients.ExchangeClient
}

func NewExchangeService(dao *daos.ExchangeDAO, client *clients.ExchangeClient) *ExchangeService {
	return &ExchangeService{dao: dao, client: client}
}

// GetLatest — Base currency için güncel kurları döner
// Cache hit → direkt döner | Cache miss → API çek, cache'e yaz, döner
func (s *ExchangeService) GetLatest(ctx context.Context, base string, symbols []string) (*models.ExchangeRate, error) {
	base = strings.ToUpper(base)

	cached, err := s.dao.FindByBase(ctx, base)
	if err != nil {
		return nil, err
	}
	if cached != nil {
		return cached, nil
	}

	resp, err := s.client.GetLatest(ctx, base, symbols)
	if err != nil {
		return nil, err
	}

	rate := &models.ExchangeRate{
		Base:  resp.Base,
		Rates: resp.Rates,
		Date:  resp.Date,
	}

	if err := s.dao.Upsert(ctx, rate); err != nil {
		return nil, err
	}

	return rate, nil
}

// Convert — Para birimi dönüştürme, direkt API (cache'lenmez çünkü amount değişken)
func (s *ExchangeService) Convert(ctx context.Context, from, to string, amount float64) (float64, string, error) {
	resp, err := s.client.Convert(ctx, strings.ToUpper(from), strings.ToUpper(to), amount)
	if err != nil {
		return 0, "", err
	}

	result := resp.Rates[strings.ToUpper(to)]
	return result, resp.Date, nil
}
