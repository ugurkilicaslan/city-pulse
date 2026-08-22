package handlers

import (
	"context"

	"city-pulse/app/services"
)

// Handler — Tüm HTTP handler'larını tutan ana struct
// Dependency injection: servisler constructor'dan geçirilir
type Handler struct {
	exchangeSvc *services.ExchangeService
	newsSvc     *services.NewsService
	gameSvc     *services.GameService
	citySvc     *services.CityService
	ctx         context.Context
}

// New — Handler constructor, route-init.go'dan çağrılır
func New(
	exchangeSvc *services.ExchangeService,
	newsSvc *services.NewsService,
	gameSvc *services.GameService,
	citySvc *services.CityService,
	ctx context.Context,
) *Handler {
	return &Handler{
		exchangeSvc: exchangeSvc,
		newsSvc:     newsSvc,
		gameSvc:     gameSvc,
		citySvc:     citySvc,
		ctx:         ctx,
	}
}
