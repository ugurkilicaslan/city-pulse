package services

import (
	"context"
	"time"

	"city-pulse/app/clients"
	"city-pulse/app/daos"
	"city-pulse/app/dtos"
	"city-pulse/app/models"
)

type NasaService struct {
	client *clients.NASAClient
	dao    *daos.NasaDAO
}

func NewNasaService(client *clients.NASAClient, dao *daos.NasaDAO) *NasaService {
	return &NasaService{
		client: client,
		dao:    dao,
	}
}

func (s *NasaService) GetApod(ctx context.Context) (*dtos.ApodDTO, error) {

	today := time.Now().Format("2006-01-02")

	cached, err := s.dao.GetByDate(ctx, today)
	if err == nil && cached != nil {
		return &dtos.ApodDTO{
			Title:       cached.Title,
			Explanation: cached.Explanation,
			URL:         cached.URL,
			HdURL:       cached.HdURL,
			MediaType:   cached.MediaType,
			Date:        cached.Date,
		}, nil
	}

	resp, err := s.client.GetAPOD(ctx)
	if err != nil {
		return nil, err
	}

	entry := &models.ApodEntry{
		Title:       resp.Title,
		Explanation: resp.Explanation,
		URL:         resp.URL,
		HdURL:       resp.HdURL,
		MediaType:   resp.MediaType,
		Date:        resp.Date,
		CachedAt:    time.Now(),
	}
	_ = s.dao.Save(ctx, entry)

	return &dtos.ApodDTO{
		Title:       resp.Title,
		Explanation: resp.Explanation,
		URL:         resp.URL,
		HdURL:       resp.HdURL,
		MediaType:   resp.MediaType,
		Date:        resp.Date,
	}, nil
}
