package clients

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

const coingeckoURL = "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum,binancecoin,solana,cardano,ripple,dogecoin&vs_currencies=usd,eur,try&include_24hr_change=true&include_market_cap=true"

// CryptoClient — CoinGecko API istemcisi (ücretsiz tier)
type CryptoClient struct {
	http *http.Client
}

func NewCryptoClient() *CryptoClient {
	return &CryptoClient{
		http: &http.Client{Timeout: 10 * time.Second},
	}
}

// CoinData — Tek para birimi verisi
type CoinData struct {
	USD          float64 `json:"usd"`
	EUR          float64 `json:"eur"`
	TRY          float64 `json:"try"`
	USD24hChg    float64 `json:"usd_24h_change"`
	EUR24hChg    float64 `json:"eur_24h_change"`
	USDMarketCap float64 `json:"usd_market_cap"`
}

// RawCoinPrices — coin_id -> CoinData
type RawCoinPrices map[string]CoinData

// FetchPrices — Kripto para fiyatlarını çeker
func (c *CryptoClient) FetchPrices(ctx context.Context) (RawCoinPrices, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, coingeckoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "city-pulse/2.0 (educational)")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var prices RawCoinPrices
	return prices, json.NewDecoder(resp.Body).Decode(&prices)
}
