package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"city-pulse/app/dtos"
)

// ExchangeLatest — GET /api/1.0/exchange/latest?base=USD&symbols=TRY,EUR,GBP
func (h *Handler) ExchangeLatest(c *gin.Context) {
	base := c.DefaultQuery("base", "USD")

	symbolsStr := c.Query("symbols")
	var symbols []string
	if symbolsStr != "" {
		symbols = strings.Split(symbolsStr, ",")
	}

	rates, err := h.exchangeSvc.GetLatest(c.Request.Context(), base, symbols)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dtos.LatestRatesResponse{
		Base:  rates.Base,
		Date:  rates.Date,
		Rates: rates.Rates,
	})
}

// ExchangeConvert — GET /api/1.0/exchange/convert?from=USD&to=TRY&amount=100
func (h *Handler) ExchangeConvert(c *gin.Context) {
	var req dtos.ConvertRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, date, err := h.exchangeSvc.Convert(c.Request.Context(), req.From, req.To, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dtos.ConvertResponse{
		From:   strings.ToUpper(req.From),
		To:     strings.ToUpper(req.To),
		Amount: req.Amount,
		Result: result,
		Date:   date,
	})
}
