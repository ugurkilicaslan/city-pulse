package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const frankfurterBaseURL = "https://api.frankfurter.app"

// ExchangeClient — Frankfurter API HTTP istemcisi (API key gerektirmez, tamamen ücretsiz)
type ExchangeClient struct {
	httpClient *http.Client
}

func NewExchangeClient() *ExchangeClient {
	return &ExchangeClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// FrankfurterResponse — Frankfurter API'nin ham cevabı
type FrankfurterResponse struct {
	Amount float64            `json:"amount"`
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  map[string]float64 `json:"rates"`
}

// GetLatest — Belirli bir base currency için güncel kurları çeker
// Örnek: GetLatest(ctx, "USD", []string{"TRY", "EUR", "GBP"})
func (c *ExchangeClient) GetLatest(ctx context.Context, base string, symbols []string) (*FrankfurterResponse, error) {
	url := fmt.Sprintf("%s/latest?base=%s", frankfurterBaseURL, base)
	if len(symbols) > 0 {
		url += "&symbols=" + strings.Join(symbols, ",")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("frankfurter isteği başarısız: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("frankfurter API hata kodu: %d", resp.StatusCode)
	}

	var result FrankfurterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("frankfurter JSON parse hatası: %w", err)
	}

	return &result, nil
}

// Convert — Para birimi dönüştürme
// Örnek: Convert(ctx, "USD", "TRY", 100.0)
func (c *ExchangeClient) Convert(ctx context.Context, from, to string, amount float64) (*FrankfurterResponse, error) {
	url := fmt.Sprintf("%s/latest?amount=%.2f&from=%s&to=%s", frankfurterBaseURL, amount, from, to)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("frankfurter convert isteği başarısız: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("frankfurter convert hata kodu: %d", resp.StatusCode)
	}

	var result FrankfurterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("frankfurter convert JSON parse hatası: %w", err)
	}

	return &result, nil
}
