package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"city-pulse/app/dtos"
)

// GameDeals — GET /api/1.0/games/deals?maxPrice=15&pageSize=20
func (h *Handler) GameDeals(c *gin.Context) {
	maxPrice, _ := strconv.ParseFloat(c.DefaultQuery("maxPrice", "0"), 64)
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	deals, err := h.gameSvc.GetDeals(c.Request.Context(), maxPrice, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dealDTOs := make([]dtos.GameDealDTO, len(deals))
	for i, d := range deals {
		dealDTOs[i] = dtos.GameDealDTO{
			DealID:          d.DealID,
			Title:           d.Title,
			SalePrice:       d.SalePrice,
			NormalPrice:     d.NormalPrice,
			SavingsPercent:  d.Savings,
			MetacriticScore: d.MetacriticScore,
			SteamRating:     d.SteamRatingText,
			Thumb:           d.Thumb,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(dealDTOs),
		"deals": dealDTOs,
	})
}

// GameSearch — GET /api/1.0/games/search?title=witcher
func (h *Handler) GameSearch(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'title' parametresi zorunlu"})
		return
	}

	games, err := h.gameSvc.SearchGames(c.Request.Context(), title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gameDTOs := make([]dtos.GameSearchDTO, len(games))
	for i, g := range games {
		gameDTOs[i] = dtos.GameSearchDTO{
			GameID:        g.GameID,
			Title:         g.Title,
			CheapestPrice: g.CheapestPrice,
			Thumb:         g.Thumb,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(gameDTOs),
		"games": gameDTOs,
	})
}
