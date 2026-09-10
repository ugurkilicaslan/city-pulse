package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GithubTrending — GET /api/1.0/github/trending
func (h *Handler) GithubTrending(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	repos, err := h.githubSvc.GetTrendingRepos(h.ctx, limit)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "github trendleri alinamadi: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"repos": repos,
		"count": len(repos),
	})
}