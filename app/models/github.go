package models

// GitHubRepo — GitHub trending repo modeli (DB'ye kaydedilmiyor)
type GitHubRepo struct {
	FullName    string `json:"fullName"`
	Description string `json:"description"`
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	Language    string `json:"language"`
	URL         string `json:"url"`
	OwnerLogin  string `json:"ownerLogin"`
	AvatarURL   string `json:"avatarURL"`
}
