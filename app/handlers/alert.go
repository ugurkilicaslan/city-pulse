package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"city-pulse/app/models"
	"city-pulse/app/services"
)

// AlertList — GET /api/1.0/me/alerts
func (h *Handler) AlertList(c *gin.Context) {
	uid, err := userID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	alerts, err := h.alertSvc.List(h.ctx, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"alerts": alerts, "count": len(alerts)})
}

// AlertCreate — POST /api/1.0/me/alerts
func (h *Handler) AlertCreate(c *gin.Context) {
	uid, err := userID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var a models.PriceAlert
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if a.Pair == "" || a.TargetRate <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parite ve hedef kur zorunludur"})
		return
	}
	if a.Condition != models.AlertAbove && a.Condition != models.AlertBelow {
		a.Condition = models.AlertAbove
	}
	a.UserID = uid

	if err := h.alertSvc.Create(h.ctx, &a); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"alert": a})
}

// AlertDelete — DELETE /api/1.0/me/alerts/:id
func (h *Handler) AlertDelete(c *gin.Context) {
	uid, err := userID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gecersiz id"})
		return
	}

	if err := h.alertSvc.Delete(h.ctx, id, uid); err != nil {
		if errors.Is(err, services.ErrAlertNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "uyari bulunamadi"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "uyari silindi"})
}