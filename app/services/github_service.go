package services

import (
	"context"

	"city-pulse/app/clients"
	"city-pulse/app/dtos"
)

type GithubService struct {
	client *clients.GitHubClient
}

func NewGithubService(client *clients.GitHubClient) *GithubService {
	return &GithubService{
		client: client,
	}
}

func (s *GithubService) GetTrendingRepos(ctx context.Context, limit int) ([]dtos.GitHubRepoDTO, error) {
	rawRepos, err := s.client.GetTrending(ctx, limit)
	if err != nil {
		return nil, err
	}

	var dtoList []dtos.GitHubRepoDTO
	for _, r := range rawRepos {
		dtoList = append(dtoList, dtos.GitHubRepoDTO{
			FullName:    r.FullName,
			Description: r.Description,
			Stars:       r.Stars,
			Language:    r.Language,
			URL:         r.HTMLURL,
			AvatarURL:   r.Owner.AvatarURL,
		})
	}
	return dtoList, nil
}
