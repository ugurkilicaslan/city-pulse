package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const cheapSharkBaseURL = "https://www.cheapshark.com/api/1.0"

// GameClient — CheapShark API HTTP istemcisi (tamamen ücretsiz, key gerektirmez)
type GameClient struct {
	httpClient *http.Client
}

func NewGameClient() *GameClient {
	return &GameClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// CheapSharkDeal — CheapShark'tan gelen ham indirim nesnesi
type CheapSharkDeal struct {
	DealID          string `json:"dealID"`
	Title           string `json:"title"`
	StoreID         string `json:"storeID"`
	SalePrice       string `json:"salePrice"`
	NormalPrice     string `json:"normalPrice"`
	Savings         string `json:"savings"`
	MetacriticScore string `json:"metacriticScore"`
	SteamRatingText string `json:"steamRatingText"`
	Thumb           string `json:"thumb"`
}

// CheapSharkGame — CheapShark arama sonucu ham nesnesi
type CheapSharkGame struct {
	GameID         string `json:"gameID"`
	External       string `json:"external"` // Oyun adı
	Cheapest       string `json:"cheapest"`
	CheapestDealID string `json:"cheapestDealID"`
	Thumb          string `json:"thumb"`
}

// GetDeals — İndirimli oyunları çeker
// maxPrice: 0 ise filtre yok | pageSize: kaç sonuç dönsün
func (c *GameClient) GetDeals(ctx context.Context, maxPrice float64, pageSize int) ([]CheapSharkDeal, error) {
	if pageSize <= 0 {
		pageSize = 20
	}

	endpoint := fmt.Sprintf("%s/deals?pageSize=%d", cheapSharkBaseURL, pageSize)
	if maxPrice > 0 {
		endpoint += fmt.Sprintf("&upperPrice=%.2f", maxPrice)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "city-pulse/1.0 (learning-project)")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cheapshark deals isteği başarısız: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cheapshark API hata kodu: %d", resp.StatusCode)
	}

	var deals []CheapSharkDeal
	if err := json.NewDecoder(resp.Body).Decode(&deals); err != nil {
		return nil, fmt.Errorf("cheapshark JSON parse hatası: %w", err)
	}

	return deals, nil
}

// SearchGames — Oyun adına göre arama yapar
// title: oyun adı (kısmi eşleşme desteklenir)
func (c *GameClient) SearchGames(ctx context.Context, title string) ([]CheapSharkGame, error) {
	endpoint := fmt.Sprintf("%s/games?title=%s&limit=10", cheapSharkBaseURL, url.QueryEscape(title))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "city-pulse/1.0 (learning-project)")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cheapshark search isteği başarısız: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cheapshark search hata kodu: %d", resp.StatusCode)
	}

	var games []CheapSharkGame
	if err := json.NewDecoder(resp.Body).Decode(&games); err != nil {
		return nil, fmt.Errorf("cheapshark search JSON parse hatası: %w", err)
	}

	return games, nil
}
