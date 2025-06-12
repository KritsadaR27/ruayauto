package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"ruaymanagement/backend/internal/ChatBotCore/models"
	"ruaymanagement/backend/internal/ChatBotCore/repository"
)

// KeywordHandler handles HTTP requests for keyword operations
type KeywordHandler struct {
	keywordRepo repository.KeywordRepository
}

// NewKeywordHandler creates a new keyword handler
func NewKeywordHandler(keywordRepo repository.KeywordRepository) *KeywordHandler {
	return &KeywordHandler{
		keywordRepo: keywordRepo,
	}
}

// CreateKeyword creates a new keyword
// POST /api/v1/keywords
func (h *KeywordHandler) CreateKeyword(c *gin.Context) {
	var keyword models.Keyword
	if err := c.ShouldBindJSON(&keyword); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request format: " + err.Error(),
			},
		})
		return
	}

	// Validate required fields
	if keyword.Keyword == "" || keyword.ResponseContent == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "MISSING_REQUIRED_FIELDS",
				Message: "keyword and response_content are required",
			},
		})
		return
	}

	// Set defaults
	if keyword.MatchType == "" {
		keyword.MatchType = "contains"
	}
	if keyword.ResponseType == "" {
		keyword.ResponseType = "text"
	}
	if keyword.Priority == 0 {
		keyword.Priority = 1
	}

	err := h.keywordRepo.Create(c.Request.Context(), &keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to create keyword: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    keyword,
	})
}

// GetKeywords retrieves keywords with pagination
// GET /api/v1/keywords
func (h *KeywordHandler) GetKeywords(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	activeOnly := c.Query("active_only") == "true"

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	var keywords []models.Keyword
	var err error

	if activeOnly {
		keywords, err = h.keywordRepo.GetActive(c.Request.Context())
		// For active keywords, we don't apply pagination since it's usually for matching
		if len(keywords) > offset {
			end := offset + perPage
			if end > len(keywords) {
				end = len(keywords)
			}
			keywords = keywords[offset:end]
		} else {
			keywords = []models.Keyword{}
		}
	} else {
		keywords, err = h.keywordRepo.List(c.Request.Context(), offset, perPage)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to retrieve keywords: " + err.Error(),
			},
		})
		return
	}

	// Get total count
	total, err := h.keywordRepo.Count(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to count keywords: " + err.Error(),
			},
		})
		return
	}

	// Calculate total pages
	totalPages := (total + perPage - 1) / perPage

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    keywords,
		Meta: &models.Meta{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// GetKeyword retrieves a specific keyword
// GET /api/v1/keywords/:id
func (h *KeywordHandler) GetKeyword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_ID",
				Message: "Invalid keyword ID",
			},
		})
		return
	}

	keyword, err := h.keywordRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "KEYWORD_NOT_FOUND",
				Message: "Keyword not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    keyword,
	})
}

// UpdateKeyword updates an existing keyword
// PUT /api/v1/keywords/:id
func (h *KeywordHandler) UpdateKeyword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_ID",
				Message: "Invalid keyword ID",
			},
		})
		return
	}

	var keyword models.Keyword
	if err := c.ShouldBindJSON(&keyword); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request format: " + err.Error(),
			},
		})
		return
	}

	keyword.ID = id
	err = h.keywordRepo.Update(c.Request.Context(), &keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to update keyword: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    keyword,
	})
}

// DeleteKeyword deletes a keyword
// DELETE /api/v1/keywords/:id
func (h *KeywordHandler) DeleteKeyword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_ID",
				Message: "Invalid keyword ID",
			},
		})
		return
	}

	err = h.keywordRepo.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to delete keyword: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    gin.H{"message": "Keyword deleted successfully"},
	})
}

// ToggleKeywordStatus toggles keyword active status
// PATCH /api/v1/keywords/:id/toggle
func (h *KeywordHandler) ToggleKeywordStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_ID",
				Message: "Invalid keyword ID",
			},
		})
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request format: " + err.Error(),
			},
		})
		return
	}

	err = h.keywordRepo.SetActive(c.Request.Context(), id, req.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to toggle keyword status: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    gin.H{"message": "Keyword status updated successfully"},
	})
}
