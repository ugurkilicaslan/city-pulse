package dtos

// GitHubRepoDTO — Tek bir GitHub trending repo response'u
type GitHubRepoDTO struct {
	FullName    string `json:"fullName"`
	Description string `json:"description"`
	Stars       int    `json:"stars"`
	Language    string `json:"language"`
	URL         string `json:"url"`
	AvatarURL   string `json:"avatarURL"`
}
