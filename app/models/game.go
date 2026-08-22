package models

// GameDeal — CheapShark'tan gelen oyun indirimi
type GameDeal struct {
	DealID          string `json:"dealID"`
	Title           string `json:"title"`
	StoreID         string `json:"storeID"`
	SalePrice       string `json:"salePrice"`
	NormalPrice     string `json:"normalPrice"`
	Savings         string `json:"savings"`
	MetacriticScore string `json:"metacriticScore"`
	SteamRatingText string `json:"steamRatingText"`
	Thumb           string `json:"thumb"`
}

// GameInfo — CheapShark oyun arama sonucu
type GameInfo struct {
	GameID         string `json:"gameID"`
	Title          string `json:"title"`
	CheapestPrice  string `json:"cheapestPrice"`
	CheapestDealID string `json:"cheapestDealID"`
	Thumb          string `json:"thumb"`
}
