package services

import (
	"context"
	"fmt"
	"time"

	"city-pulse/app/clients"
	"city-pulse/internal/cache"
)

// WeatherData — Hava durumu özet verisi
type WeatherData struct {
	City        string  `json:"city"`
	Region      string  `json:"region,omitempty"`
	Country     string  `json:"country"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"windSpeed"`
	Condition   string  `json:"condition"`
	Emoji       string  `json:"emoji"`
	IsDay       bool    `json:"isDay"`
	FetchedAt   string  `json:"fetchedAt"`
}

// WeatherService — Hava durumu iş mantığı
type WeatherService struct {
	client *clients.WeatherClient
	cache  *cache.Cache
}

func NewWeatherService(client *clients.WeatherClient, c *cache.Cache) *WeatherService {
	return &WeatherService{client: client, cache: c}
}

// GetWeather — Şehir için hava durumu getirir (15dk cache)
func (s *WeatherService) GetWeather(ctx context.Context, city string) (*WeatherData, error) {
	cacheKey := "weather:" + city
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.(*WeatherData), nil
	}

	geo, err := s.client.Geocode(ctx, city)
	if err != nil {
		return nil, fmt.Errorf("geocode hatası: %w", err)
	}

	forecast, err := s.client.FetchForecast(ctx, geo.Latitude, geo.Longitude, geo.Timezone)
	if err != nil {
		return nil, fmt.Errorf("hava durumu hatası: %w", err)
	}

	cw := forecast.CurrentWeather
	data := &WeatherData{
		City:        geo.Name,
		Region:      geo.Admin1,
		Country:     geo.Country,
		Latitude:    geo.Latitude,
		Longitude:   geo.Longitude,
		Temperature: cw.Temperature,
		WindSpeed:   cw.WindSpeed,
		Condition:   weatherCodeToCondition(cw.WeatherCode),
		Emoji:       weatherCodeToEmoji(cw.WeatherCode),
		IsDay:       cw.IsDay == 1,
		FetchedAt:   time.Now().Format(time.RFC3339),
	}

	s.cache.Set(cacheKey, data, 15*time.Minute)
	return data, nil
}

func weatherCodeToCondition(code int) string {
	switch {
	case code == 0:
		return "Açık"
	case code <= 3:
		return "Parçalı Bulutlu"
	case code <= 9:
		return "Sis/Pus"
	case code <= 19:
		return "Hafif Yağmur"
	case code <= 29:
		return "Fırtına"
	case code <= 49:
		return "Sis"
	case code <= 59:
		return "Çisenti"
	case code <= 69:
		return "Yağmur"
	case code <= 79:
		return "Kar"
	case code <= 89:
		return "Sağanak Yağış"
	case code <= 99:
		return "Gök Gürültülü Fırtına"
	default:
		return "Belirsiz"
	}
}

func weatherCodeToEmoji(code int) string {
	switch {
	case code == 0:
		return "☀️"
	case code <= 3:
		return "⛅"
	case code <= 49:
		return "🌫️"
	case code <= 69:
		return "🌧️"
	case code <= 79:
		return "❄️"
	default:
		return "⛈️"
	}
}
