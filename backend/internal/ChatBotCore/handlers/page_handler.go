package handler

import (
	"net/http"

	service "ruaychatbot/backend/internal/ChatBotCore/services"

	"github.com/gin-gonic/gin"
)

// PageHandler handles page management operations
type PageHandler struct {
	pageManager *service.PageManager
}

// NewPageHandler creates a new page handler
func NewPageHandler(repos *service.Repositories) *PageHandler {
	return &PageHandler{
		pageManager: service.NewPageManager(repos),
	}
}

// CreatePageRequest represents the request to create a new page
type CreatePageRequest struct {
	Platform    string                `json:"platform" binding:"required"`
	PageID      string                `json:"page_id" binding:"required"`
	PageName    string                `json:"page_name" binding:"required"`
	AccessToken string                `json:"access_token"`
	Settings    *service.PageSettings `json:"settings"`
}

// UpdatePageRequest represents the request to update a page
type UpdatePageRequest struct {
	PageName    *string               `json:"page_name"`
	AccessToken *string               `json:"access_token"`
	IsEnabled   *bool                 `json:"is_enabled"`
	IsConnected *bool                 `json:"is_connected"`
	Settings    *service.PageSettings `json:"settings"`
}

// GetPages returns all configured pages
func (h *PageHandler) GetPages(c *gin.Context) {
	platform := c.Query("platform")

	var pages []*service.PageConfig

	if platform != "" {
		pages = h.pageManager.GetPagesByPlatform(platform)
	} else {
		pages = h.pageManager.GetAllPages()
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    pages,
		"meta": gin.H{
			"count":    len(pages),
			"platform": platform,
		},
	})
}

// GetPage returns a specific page by ID
func (h *PageHandler) GetPage(c *gin.Context) {
	pageID := c.Param("id")

	page, err := h.pageManager.GetPage(pageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    page,
	})
}

// GetEnabledPages returns all enabled pages
func (h *PageHandler) GetEnabledPages(c *gin.Context) {
	pages := h.pageManager.GetEnabledPages()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    pages,
		"meta": gin.H{
			"count": len(pages),
		},
	})
}

// GetConnectedPages returns all connected pages
func (h *PageHandler) GetConnectedPages(c *gin.Context) {
	pages := h.pageManager.GetConnectedPages()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    pages,
		"meta": gin.H{
			"count": len(pages),
		},
	})
}

// CreatePage creates a new page configuration
func (h *PageHandler) CreatePage(c *gin.Context) {
	var req CreatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Create page config
	pageConfig := &service.PageConfig{
		Platform:    req.Platform,
		PageID:      req.PageID,
		PageName:    req.PageName,
		AccessToken: req.AccessToken,
		IsEnabled:   true,
		IsConnected: req.AccessToken != "", // Connected if access token is provided
		Settings:    req.Settings,
	}

	if err := h.pageManager.AddPage(pageConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    pageConfig,
		"message": "Page created successfully",
	})
}

// UpdatePage updates an existing page
func (h *PageHandler) UpdatePage(c *gin.Context) {
	pageID := c.Param("id")

	var req UpdatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Get existing page
	page, err := h.pageManager.GetPage(pageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Apply updates
	updates := &service.PageConfig{
		IsEnabled:   page.IsEnabled,
		IsConnected: page.IsConnected,
	}

	if req.PageName != nil {
		updates.PageName = *req.PageName
	}

	if req.AccessToken != nil {
		updates.AccessToken = *req.AccessToken
		updates.IsConnected = *req.AccessToken != ""
	}

	if req.IsEnabled != nil {
		updates.IsEnabled = *req.IsEnabled
	}

	if req.IsConnected != nil {
		updates.IsConnected = *req.IsConnected
	}

	if req.Settings != nil {
		updates.Settings = req.Settings
	}

	if err := h.pageManager.UpdatePage(pageID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Get updated page
	updatedPage, _ := h.pageManager.GetPage(pageID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updatedPage,
		"message": "Page updated successfully",
	})
}

// EnablePage enables a page
func (h *PageHandler) EnablePage(c *gin.Context) {
	pageID := c.Param("id")

	if err := h.pageManager.EnablePage(pageID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Page enabled successfully",
	})
}

// DisablePage disables a page
func (h *PageHandler) DisablePage(c *gin.Context) {
	pageID := c.Param("id")

	if err := h.pageManager.DisablePage(pageID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Page disabled successfully",
	})
}

// DeletePage deletes a page
func (h *PageHandler) DeletePage(c *gin.Context) {
	pageID := c.Param("id")

	if err := h.pageManager.DeletePage(pageID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Page deleted successfully",
	})
}

// GetPageSettings returns settings for a specific page
func (h *PageHandler) GetPageSettings(c *gin.Context) {
	pageID := c.Param("id")

	settings, err := h.pageManager.GetPageSettings(pageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    settings,
	})
}

// UpdatePageSettings updates settings for a specific page
func (h *PageHandler) UpdatePageSettings(c *gin.Context) {
	pageID := c.Param("id")

	var settings service.PageSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := h.pageManager.UpdatePageSettings(pageID, &settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    settings,
		"message": "Page settings updated successfully",
	})
}

// GetPageStats returns statistics for a specific page
func (h *PageHandler) GetPageStats(c *gin.Context) {
	pageID := c.Param("id")

	stats, err := h.pageManager.GetPageStats(c.Request.Context(), pageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// GetManagerStats returns overall page manager statistics
func (h *PageHandler) GetManagerStats(c *gin.Context) {
	stats := h.pageManager.GetManagerStats()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// RegisterRoutes registers all page management routes
func (h *PageHandler) RegisterRoutes(r *gin.RouterGroup) {
	pages := r.Group("/pages")
	{
		pages.GET("", h.GetPages)
		pages.POST("", h.CreatePage)
		pages.GET("/enabled", h.GetEnabledPages)
		pages.GET("/connected", h.GetConnectedPages)
		pages.GET("/stats", h.GetManagerStats)

		// Single page operations
		pages.GET("/:id", h.GetPage)
		pages.PUT("/:id", h.UpdatePage)
		pages.DELETE("/:id", h.DeletePage)
		pages.POST("/:id/enable", h.EnablePage)
		pages.POST("/:id/disable", h.DisablePage)
		pages.GET("/:id/stats", h.GetPageStats)

		// Page settings
		pages.GET("/:id/settings", h.GetPageSettings)
		pages.PUT("/:id/settings", h.UpdatePageSettings)
	}
}
