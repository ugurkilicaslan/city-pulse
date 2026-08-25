package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"city-pulse/app/models"
	"city-pulse/app/services"
)

// userID — JWT context'inden kullanıcı ObjectID'sini çeker
func userID(c *gin.Context) (bson.ObjectID, error) {
	raw, exists := c.Get("userID")
	if !exists {
		return bson.NilObjectID, errors.New("userID bulunamadı")
	}
	str, ok := raw.(string)
	if !ok {
		return bson.NilObjectID, errors.New("userID formatı geçersiz")
	}
	id, err := bson.ObjectIDFromHex(str)
	if err != nil {
		return bson.NilObjectID, errors.New("userID parse hatası")
	}
	return id, nil
}

// MeProfile — GET /api/1.0/me
func (h *Handler) MeProfile(c *gin.Context) {
	uid, err := userID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	pref, _ := h.prefSvc.Get(h.ctx, uid)
	bCount, _ := h.bookmarkSvc.Count(h.ctx, uid)

	c.JSON(http.StatusOK, gin.H{
		"id":             c.GetString("userID"),
		"email":          c.GetString("email"),
		"username":       c.GetString("username"),
		"preferences":    pref,
		"bookmark_count": bCount,
	})
}

// MePreferences — GET /api/1.0/me/preferences
func (h *Handler) MePreferences(c *gin.Context) {
	uid, err := userID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	pref, err := h.prefSvc.Get(h.ctx, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pref)
}

// MeUpdatePreferences — PUT /api/1.0/me/preferences
func (h *Handler) MeUpdatePreferences(c *gin.Context) {
	uid, err := userID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var pref models.UserPreference
	if err := c.ShouldBindJSON(&pref); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pref.UserID = uid

	if err := h.prefSvc.Update(h.ctx, &pref); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "tercihler güncellendi", "preferences": pref})
}

// MeResetPreferences — DELETE /api/1.0/me/preferences
func (h *Handler) MeResetPreferences(c *gin.Context) {
	uid, err := userID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	defaults, err := h.prefSvc.Reset(h.ctx, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "tercihler sıfırlandı", "preferences": defaults})
}

// BookmarkList — GET /api/1.0/me/bookmarks
func (h *Handler) BookmarkList(c *gin.Context) {
	uid, err := userID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	btype := c.Query("type")
	var list []models.Bookmark
	if btype != "" {
		list, err = h.bookmarkSvc.ListByType(h.ctx, uid, models.BookmarkType(btype))
	} else {
		list, err = h.bookmarkSvc.List(h.ctx, uid)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"bookmarks": list, "count": len(list)})
}

// BookmarkCreate — POST /api/1.0/me/bookmarks
func (h *Handler) BookmarkCreate(c *gin.Context) {
	uid, err := userID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var b models.Bookmark
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	b.UserID = uid

	if err := h.bookmarkSvc.Add(h.ctx, &b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"bookmark": b})
}

// BookmarkDelete — DELETE /api/1.0/me/bookmarks/:id
func (h *Handler) BookmarkDelete(c *gin.Context) {
	uid, err := userID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "geçersiz id"})
		return
	}

	if err := h.bookmarkSvc.Remove(h.ctx, id, uid); err != nil {
		if errors.Is(err, services.ErrBookmarkNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "yer imi bulunamadı"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "yer imi silindi"})
}
