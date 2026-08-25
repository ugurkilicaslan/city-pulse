package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	geocodeURL  = "https://geocoding-api.open-meteo.com/v1/search"
	forecastURL = "https://api.open-meteo.com/v1/forecast"
)

// WeatherClient — Open-Meteo API istemcisi (ücretsiz, API key gerektirmez)
type WeatherClient struct {
	http *http.Client
}

func NewWeatherClient() *WeatherClient {
	return &WeatherClient{
		http: &http.Client{Timeout: 8 * time.Second},
	}
}

// GeoLocation — Şehir koordinat bilgisi
type GeoLocation struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
	Timezone  string  `json:"timezone"`
	Admin1    string  `json:"admin1"`
}

// CurrentWeather — Anlık hava durumu
type CurrentWeather struct {
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"windspeed"`
	WeatherCode int     `json:"weathercode"`
	IsDay       int     `json:"is_day"`
	Time        string  `json:"time"`
}

// OpenMeteoResponse — Open-Meteo API yanıtı
type OpenMeteoResponse struct {
	Latitude       float64        `json:"latitude"`
	Longitude      float64        `json:"longitude"`
	Timezone       string         `json:"timezone"`
	CurrentWeather CurrentWeather `json:"current_weather"`
}

// Geocode — Şehir adını koordinata çevirir
func (c *WeatherClient) Geocode(ctx context.Context, city string) (*GeoLocation, error) {
	u := fmt.Sprintf("%s?name=%s&count=1&language=tr&format=json", geocodeURL, url.QueryEscape(city))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("geocode isteği başarısız: %w", err)
	}
	defer resp.Body.Close()

	var payload struct {
		Results []GeoLocation `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if len(payload.Results) == 0 {
		return nil, fmt.Errorf("şehir bulunamadı: %s", city)
	}
	return &payload.Results[0], nil
}

// FetchForecast — Koordinat için hava durumu getirir
func (c *WeatherClient) FetchForecast(ctx context.Context, lat, lon float64, tz string) (*OpenMeteoResponse, error) {
	if tz == "" {
		tz = "auto"
	}
	u := fmt.Sprintf("%s?latitude=%.4f&longitude=%.4f&current_weather=true&timezone=%s",
		forecastURL, lat, lon, url.QueryEscape(tz))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("hava durumu isteği başarısız: %w", err)
	}
	defer resp.Body.Close()

	var result OpenMeteoResponse
	return &result, json.NewDecoder(resp.Body).Decode(&result)
}
