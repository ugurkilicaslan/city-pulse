package dtos

// LatestRatesResponse — GET /exchange/latest cevabı
type LatestRatesResponse struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

// ConvertRequest — GET /exchange/convert query parametreleri
type ConvertRequest struct {
	From   string  `form:"from"   binding:"required"`
	To     string  `form:"to"     binding:"required"`
	Amount float64 `form:"amount" binding:"required,gt=0"`
}

// ConvertResponse — GET /exchange/convert cevabı
type ConvertResponse struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
	Result float64 `json:"result"`
	Date   string  `json:"date"`
}
