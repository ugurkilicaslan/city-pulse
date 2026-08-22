package dtos

// CitySnapshotResponse — GET /city/snapshot birleşik cevabı
// Döviz + Haberler + Oyun + NASA APOD + GitHub Trending
type CitySnapshotResponse struct {
	City        string               `json:"city"`
	Exchange    *LatestRatesResponse `json:"exchange"`
	News        []ArticleDTO         `json:"news"`
	GameDeals   []GameDealDTO        `json:"gameDeals"`
	NasaApod    *ApodDTO             `json:"nasaApod"`
	GithubTrend []GitHubRepoDTO      `json:"githubTrend"`
}
