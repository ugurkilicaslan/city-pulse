package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const githubBaseURL = "https://api.github.com"

// GitHubClient — GitHub Search API istemcisi (key gerektirmez)
type GitHubClient struct {
	httpClient *http.Client
}

func NewGitHubClient() *GitHubClient {
	return &GitHubClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// GitHubOwner — Repo sahibi
type GitHubOwner struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

// GitHubRepoRaw — GitHub API ham repo nesnesi
type GitHubRepoRaw struct {
	FullName    string      `json:"full_name"`
	Description string      `json:"description"`
	Stars       int         `json:"stargazers_count"`
	Forks       int         `json:"forks_count"`
	Language    string      `json:"language"`
	HTMLURL     string      `json:"html_url"`
	Owner       GitHubOwner `json:"owner"`
}

// GitHubSearchResponse — GitHub Search API ham cevabı
type GitHubSearchResponse struct {
	TotalCount int             `json:"total_count"`
	Items      []GitHubRepoRaw `json:"items"`
}

// GetTrending — Son 7 günde oluşturulan, en yıldızlı repoları getirir
// Bu GitHub'ın resmi trending API'si olmasa da aynı sonucu verir
func (c *GitHubClient) GetTrending(ctx context.Context, limit int) ([]GitHubRepoRaw, error) {
	if limit <= 0 {
		limit = 10
	}

	since := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	url := fmt.Sprintf(
		"%s/search/repositories?q=created:>%s&sort=stars&order=desc&per_page=%d",
		githubBaseURL, since, limit,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "city-pulse/1.0 (learning-project)")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github trending isteği başarısız: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github API hata kodu: %d", resp.StatusCode)
	}

	var result GitHubSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("github JSON parse hatası: %w", err)
	}

	return result.Items, nil
}
