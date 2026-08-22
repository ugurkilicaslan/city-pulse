package dtos

// GameDealDTO — Oyun indirimi response'u
type GameDealDTO struct {
	DealID          string `json:"dealID"`
	Title           string `json:"title"`
	SalePrice       string `json:"salePrice"`
	NormalPrice     string `json:"normalPrice"`
	SavingsPercent  string `json:"savingsPercent"`
	MetacriticScore string `json:"metacriticScore"`
	SteamRating     string `json:"steamRating"`
	Thumb           string `json:"thumb"`
}

// GameSearchDTO — Oyun arama sonucu response'u
type GameSearchDTO struct {
	GameID        string `json:"gameID"`
	Title         string `json:"title"`
	CheapestPrice string `json:"cheapestPrice"`
	Thumb         string `json:"thumb"`
}
