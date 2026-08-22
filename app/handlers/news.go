package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"city-pulse/app/dtos"
)

// NewsHeadlines — GET /api/1.0/news/headlines?lang=tr
func (h *Handler) NewsHeadlines(c *gin.Context) {
	lang := c.DefaultQuery("lang", "en")

	articles, err := h.newsSvc.GetHeadlines(c.Request.Context(), lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	articleDTOs := make([]dtos.ArticleDTO, len(articles))
	for i, a := range articles {
		articleDTOs[i] = dtos.ArticleDTO{
			Title:       a.Title,
			Description: a.Description,
			URL:         a.URL,
			Image:       a.Image,
			PublishedAt: a.PublishedAt,
			Source:      dtos.SourceDTO{Name: a.Source.Name, URL: a.Source.URL},
		}
	}

	c.JSON(http.StatusOK, dtos.HeadlinesResponse{
		TotalArticles: len(articleDTOs),
		Articles:      articleDTOs,
	})
}

// NewsSearch — GET /api/1.0/news/search?q=yapay+zeka&lang=tr
func (h *Handler) NewsSearch(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'q' parametresi zorunlu"})
		return
	}

	lang := c.DefaultQuery("lang", "en")

	articles, err := h.newsSvc.Search(c.Request.Context(), query, lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	articleDTOs := make([]dtos.ArticleDTO, len(articles))
	for i, a := range articles {
		articleDTOs[i] = dtos.ArticleDTO{
			Title:       a.Title,
			Description: a.Description,
			URL:         a.URL,
			Image:       a.Image,
			PublishedAt: a.PublishedAt,
			Source:      dtos.SourceDTO{Name: a.Source.Name, URL: a.Source.URL},
		}
	}

	c.JSON(http.StatusOK, dtos.HeadlinesResponse{
		TotalArticles: len(articleDTOs),
		Articles:      articleDTOs,
	})
}
