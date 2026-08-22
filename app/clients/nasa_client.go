package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const nasaBaseURL = "https://api.nasa.gov"

// NASAClient — NASA API HTTP istemcisi
type NASAClient struct {
	httpClient *http.Client
	apiKey     string
}

func NewNASAClient(apiKey string) *NASAClient {
	if apiKey == "" {
		apiKey = "DEMO_KEY"
	}
	return &NASAClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiKey:     apiKey,
	}
}

// ApodResponse — NASA API'sinden gelen ham APOD verisi
type ApodResponse struct {
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	URL         string `json:"url"`
	HdURL       string `json:"hdurl"`
	MediaType   string `json:"media_type"` // "image" veya "video"
	Date        string `json:"date"`
}

// GetAPOD — Günün astronomi fotoğrafını getirir
func (c *NASAClient) GetAPOD(ctx context.Context) (*ApodResponse, error) {
	url := fmt.Sprintf("%s/planetary/apod?api_key=%s", nasaBaseURL, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("nasa apod isteği başarısız: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nasa apod hata kodu: %d", resp.StatusCode)
	}

	var result ApodResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("nasa apod JSON parse hatası: %w", err)
	}

	return &result, nil
}
