package handler

import (
	"net/http"
	"strconv"

	"ruaychatbot/backend/internal/ChatBotCore/models"
	"ruaychatbot/backend/internal/ChatBotCore/repository"

	"github.com/gin-gonic/gin"
)

// KeywordHandler handles keyword operations (backward compatibility)
type KeywordHandler struct {
	keywordRepo repository.KeywordRepository
}

// NewKeywordHandler creates a new keyword handler
func NewKeywordHandler(keywordRepo repository.KeywordRepository) *KeywordHandler {
	return &KeywordHandler{
		keywordRepo: keywordRepo,
	}
}

// GetKeywords returns all keywords
func (h *KeywordHandler) GetKeywords(c *gin.Context) {
	keywords, err := h.keywordRepo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    keywords,
		"meta": gin.H{
			"count": len(keywords),
		},
	})
}

// CreateKeywordRequest represents the request to create a keyword
type CreateKeywordRequest struct {
	Keyword   string `json:"keyword" binding:"required"`
	Response  string `json:"response" binding:"required"`
	MatchType string `json:"match_type"`
	Priority  int    `json:"priority"`
}

// CreateKeyword creates a new keyword
func (h *KeywordHandler) CreateKeyword(c *gin.Context) {
	var req CreateKeywordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Set default match type
	if req.MatchType == "" {
		req.MatchType = "contains"
	}

	keyword := models.SimpleKeyword{
		Keyword:   req.Keyword,
		Response:  req.Response,
		IsActive:  true,
		MatchType: req.MatchType,
		Priority:  req.Priority,
	}

	if err := h.keywordRepo.Create(c.Request.Context(), keyword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    keyword,
		"message": "Keyword created successfully",
	})
}

// UpdateKeywordRequest represents the request to update a keyword
type UpdateKeywordRequest struct {
	Keyword   *string `json:"keyword"`
	Response  *string `json:"response"`
	IsActive  *bool   `json:"is_active"`
	MatchType *string `json:"match_type"`
	Priority  *int    `json:"priority"`
}

// UpdateKeyword updates an existing keyword
func (h *KeywordHandler) UpdateKeyword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid keyword ID",
		})
		return
	}

	var req UpdateKeywordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Get existing keyword
	existing, err := h.keywordRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "keyword not found",
		})
		return
	}

	// Apply updates
	updated := *existing
	if req.Keyword != nil {
		updated.Keyword = *req.Keyword
	}
	if req.Response != nil {
		updated.Response = *req.Response
	}
	if req.IsActive != nil {
		updated.IsActive = *req.IsActive
	}
	if req.MatchType != nil {
		updated.MatchType = *req.MatchType
	}
	if req.Priority != nil {
		updated.Priority = *req.Priority
	}

	if err := h.keywordRepo.Update(c.Request.Context(), id, updated); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
		"message": "Keyword updated successfully",
	})
}

// DeleteKeyword deletes a keyword
func (h *KeywordHandler) DeleteKeyword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid keyword ID",
		})
		return
	}

	if err := h.keywordRepo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Keyword deleted successfully",
	})
}
