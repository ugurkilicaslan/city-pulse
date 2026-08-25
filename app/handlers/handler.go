package handlers

import (
	"context"

	"city-pulse/app/services"
)

type Handler struct {
	exchangeSvc  *services.ExchangeService
	newsSvc      *services.NewsService
	gameSvc      *services.GameService
	citySvc      *services.CityService
	userSvc      *services.UserService
	weatherSvc   *services.WeatherService
	cryptoSvc    *services.CryptoService
	analyticsSvc *services.AnalyticsService
	prefSvc      *services.PreferenceService
	bookmarkSvc  *services.BookmarkService
	alertSvc     *services.AlertService
	ctx          context.Context
}

func New(
	exchangeSvc *services.ExchangeService,
	newsSvc *services.NewsService,
	gameSvc *services.GameService,
	citySvc *services.CityService,
	userSvc *services.UserService,
	weatherSvc *services.WeatherService,
	cryptoSvc *services.CryptoService,
	analyticsSvc *services.AnalyticsService,
	prefSvc *services.PreferenceService,
	bookmarkSvc *services.BookmarkService,
	alertSvc *services.AlertService,
	ctx context.Context,
) *Handler {
	return &Handler{
		exchangeSvc:  exchangeSvc,
		newsSvc:      newsSvc,
		gameSvc:      gameSvc,
		citySvc:      citySvc,
		userSvc:      userSvc,
		weatherSvc:   weatherSvc,
		cryptoSvc:    cryptoSvc,
		analyticsSvc: analyticsSvc,
		prefSvc:      prefSvc,
		bookmarkSvc:  bookmarkSvc,
		alertSvc:     alertSvc,
		ctx:          ctx,
	}
}
