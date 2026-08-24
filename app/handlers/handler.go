package handlers

import (
	"context"

	"city-pulse/app/services"
)

type Handler struct {
	exchangeSvc *services.ExchangeService
	newsSvc     *services.NewsService
	gameSvc     *services.GameService
	citySvc     *services.CityService
	userSvc     *services.UserService
	ctx         context.Context
}

func New(
	exchangeSvc *services.ExchangeService,
	newsSvc *services.NewsService,
	gameSvc *services.GameService,
	citySvc *services.CityService,
	userSvc *services.UserService,
	ctx context.Context,
) *Handler {
	return &Handler{
		exchangeSvc: exchangeSvc,
		newsSvc:     newsSvc,
		gameSvc:     gameSvc,
		citySvc:     citySvc,
		userSvc:     userSvc,
		ctx:         ctx,
	}
}

