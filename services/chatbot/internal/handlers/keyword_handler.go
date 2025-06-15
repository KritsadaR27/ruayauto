package handlers

import (
	"log"
	"net/http"
	"strconv"

	"chatbot/internal/models"
	"chatbot/internal/repository"

	"github.com/gin-gonic/gin"
)

type KeywordHandler struct {
	repo *repository.KeywordRepository
}

func NewKeywordHandler(repo *repository.KeywordRepository) *KeywordHandler {
	return &KeywordHandler{repo: repo}
}

// GetKeywords returns all active keywords
func (h *KeywordHandler) GetKeywords(c *gin.Context) {
	keywords, err := h.repo.GetActive(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"pairs": keywords,
		},
		"timestamp": "now",
	})
}

// CreateKeyword creates a new keyword
func (h *KeywordHandler) CreateKeyword(c *gin.Context) {
	var req struct {
		Keyword  string `json:"keyword" binding:"required"`
		Response string `json:"response" binding:"required"`
		Priority int    `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	keyword := &models.Keyword{
		Keyword:   req.Keyword,
		Response:  req.Response,
		IsActive:  true,
		Priority:  req.Priority,
		MatchType: "contains", // default match type
	}

	if err := h.repo.Create(c.Request.Context(), keyword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    keyword,
	})
}

// UpdateKeyword updates an existing keyword
func (h *KeywordHandler) UpdateKeyword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid keyword ID"})
		return
	}

	var req struct {
		Keyword  string `json:"keyword"`
		Response string `json:"response"`
		IsActive *bool  `json:"is_active"`
		Priority int    `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ Update keyword bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("📝 Updating keyword %d with data: %+v", id, req)

	keyword := &models.Keyword{
		ID:       id,
		Keyword:  req.Keyword,
		Response: req.Response,
		Priority: req.Priority,
	}

	if req.IsActive != nil {
		keyword.IsActive = *req.IsActive
	}

	if err := h.repo.Update(c.Request.Context(), keyword); err != nil {
		log.Printf("❌ Update keyword DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    keyword,
	})
}

// DeleteKeyword deletes a keyword
func (h *KeywordHandler) DeleteKeyword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid keyword ID"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "keyword deleted successfully",
	})
}

// BulkCreateKeywords creates multiple keywords at once (for web interface)
func (h *KeywordHandler) BulkCreateKeywords(c *gin.Context) {
	var req struct {
		Pairs []struct {
			Keyword  string `json:"keyword"`
			Response string `json:"response"`
		} `json:"pairs"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// First, clear existing keywords to avoid duplicates
	if err := h.repo.DeleteAll(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear existing keywords"})
		return
	}

	// Create new keywords
	var keywords []*models.Keyword
	for i, pair := range req.Pairs {
		keyword := &models.Keyword{
			Keyword:   pair.Keyword,
			Response:  pair.Response,
			IsActive:  true,
			Priority:  i + 1,
			MatchType: "contains",
		}

		if err := h.repo.Create(c.Request.Context(), keyword); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		keywords = append(keywords, keyword)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"pairs": keywords,
			"count": len(keywords),
		},
	})
}
