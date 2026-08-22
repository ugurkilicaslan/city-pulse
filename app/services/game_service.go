package services

import (
	"context"

	"city-pulse/app/clients"
	"city-pulse/app/models"
)

// GameService — Oyun iş mantığı
// CheapShark ücretsiz ve hızlı, cache'lemeye gerek yok
type GameService struct {
	client *clients.GameClient
}

func NewGameService(client *clients.GameClient) *GameService {
	return &GameService{client: client}
}

// GetDeals — İndirimli oyunları döner
func (s *GameService) GetDeals(ctx context.Context, maxPrice float64, pageSize int) ([]models.GameDeal, error) {
	raw, err := s.client.GetDeals(ctx, maxPrice, pageSize)
	if err != nil {
		return nil, err
	}

	deals := make([]models.GameDeal, len(raw))
	for i, d := range raw {
		deals[i] = models.GameDeal{
			DealID:          d.DealID,
			Title:           d.Title,
			StoreID:         d.StoreID,
			SalePrice:       d.SalePrice,
			NormalPrice:     d.NormalPrice,
			Savings:         d.Savings,
			MetacriticScore: d.MetacriticScore,
			SteamRatingText: d.SteamRatingText,
			Thumb:           d.Thumb,
		}
	}

	return deals, nil
}

// SearchGames — Oyun adına göre arama yapar
func (s *GameService) SearchGames(ctx context.Context, title string) ([]models.GameInfo, error) {
	raw, err := s.client.SearchGames(ctx, title)
	if err != nil {
		return nil, err
	}

	games := make([]models.GameInfo, len(raw))
	for i, g := range raw {
		games[i] = models.GameInfo{
			GameID:         g.GameID,
			Title:          g.External,
			CheapestPrice:  g.Cheapest,
			CheapestDealID: g.CheapestDealID,
			Thumb:          g.Thumb,
		}
	}

	return games, nil
}
