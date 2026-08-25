package services

import (
	"context"
	"sort"
	"time"

	"city-pulse/app/clients"
	"city-pulse/internal/cache"
)

// CoinInfo — Kripto para birimi bilgisi
type CoinInfo struct {
	ID           string  `json:"id"`
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	LogoURL      string  `json:"logoUrl"`
	PriceUSD     float64 `json:"priceUsd"`
	PriceEUR     float64 `json:"priceEur"`
	PriceTRY     float64 `json:"priceTry"`
	Change24hUSD float64 `json:"change24hUsd"`
	MarketCapUSD float64 `json:"marketCapUsd"`
}

// CryptoService — Kripto para fiyat iş mantığı
type CryptoService struct {
	client *clients.CryptoClient
	cache  *cache.Cache
}

func NewCryptoService(client *clients.CryptoClient, c *cache.Cache) *CryptoService {
	return &CryptoService{client: client, cache: c}
}

// coinMeta — Sembol ve isim bilgileri
var coinMeta = map[string]struct {
	Symbol  string
	Name    string
	LogoURL string
}{
	"bitcoin":     {Symbol: "BTC", Name: "Bitcoin", LogoURL: "https://assets.coingecko.com/coins/images/1/small/bitcoin.png"},
	"ethereum":    {Symbol: "ETH", Name: "Ethereum", LogoURL: "https://assets.coingecko.com/coins/images/279/small/ethereum.png"},
	"binancecoin": {Symbol: "BNB", Name: "BNB", LogoURL: "https://assets.coingecko.com/coins/images/825/small/bnb-icon2_2x.png"},
	"solana":      {Symbol: "SOL", Name: "Solana", LogoURL: "https://assets.coingecko.com/coins/images/4128/small/solana.png"},
	"cardano":     {Symbol: "ADA", Name: "Cardano", LogoURL: "https://assets.coingecko.com/coins/images/975/small/cardano.png"},
	"ripple":      {Symbol: "XRP", Name: "XRP", LogoURL: "https://assets.coingecko.com/coins/images/44/small/xrp-symbol-white-128.png"},
	"dogecoin":    {Symbol: "DOGE", Name: "Dogecoin", LogoURL: "https://assets.coingecko.com/coins/images/5/small/dogecoin.png"},
}

// GetPrices — Tüm kripto fiyatlarını getirir (3dk cache)
func (s *CryptoService) GetPrices(ctx context.Context) ([]CoinInfo, error) {
	const key = "crypto:prices"
	if cached, ok := s.cache.Get(key); ok {
		return cached.([]CoinInfo), nil
	}

	raw, err := s.client.FetchPrices(ctx)
	if err != nil {
		return nil, err
	}

	coins := make([]CoinInfo, 0, len(raw))
	for id, data := range raw {
		meta, ok := coinMeta[id]
		if !ok {
			continue
		}
		coins = append(coins, CoinInfo{
			ID:           id,
			Symbol:       meta.Symbol,
			Name:         meta.Name,
			LogoURL:      meta.LogoURL,
			PriceUSD:     data.USD,
			PriceEUR:     data.EUR,
			PriceTRY:     data.TRY,
			Change24hUSD: data.USD24hChg,
			MarketCapUSD: data.USDMarketCap,
		})
	}

	// Market cap'e göre sırala
	sort.Slice(coins, func(i, j int) bool {
		return coins[i].MarketCapUSD > coins[j].MarketCapUSD
	})

	s.cache.Set(key, coins, 3*time.Minute)
	return coins, nil
}
