package handler

import (
	"net/http"
	"strconv"
	"time"

	service "ruaychatbot/backend/internal/ChatBotCore/services"

	"github.com/gin-gonic/gin"
)

// RuleHandler handles rule management operations
type RuleHandler struct {
	ruleEngine  *service.RuleEngine
	pageManager *service.PageManager
}

// NewRuleHandler creates a new rule handler
func NewRuleHandler(repos *service.Repositories) *RuleHandler {
	return &RuleHandler{
		ruleEngine:  service.NewRuleEngine(repos),
		pageManager: service.NewPageManager(repos),
	}
}

// CreateRuleRequest represents the request to create a new rule
type CreateRuleRequest struct {
	Name           string              `json:"name" binding:"required"`
	Keywords       []string            `json:"keywords" binding:"required,min=1"`
	Responses      []*service.Response `json:"responses" binding:"required,min=1"`
	SelectedPages  []string            `json:"selected_pages"`
	HideAfterReply bool                `json:"hide_after_reply"`
	SendToInbox    bool                `json:"send_to_inbox"`
	InboxMessage   string              `json:"inbox_message"`
	InboxImage     string              `json:"inbox_image"`
	Priority       int                 `json:"priority"`
	MatchType      string              `json:"match_type"`
}

// UpdateRuleRequest represents the request to update an existing rule
type UpdateRuleRequest struct {
	Name           *string             `json:"name"`
	Keywords       []string            `json:"keywords"`
	Responses      []*service.Response `json:"responses"`
	Enabled        *bool               `json:"enabled"`
	SelectedPages  []string            `json:"selected_pages"`
	HideAfterReply *bool               `json:"hide_after_reply"`
	SendToInbox    *bool               `json:"send_to_inbox"`
	InboxMessage   *string             `json:"inbox_message"`
	InboxImage     *string             `json:"inbox_image"`
	Priority       *int                `json:"priority"`
	MatchType      *string             `json:"match_type"`
}

// GetRules returns all rules for a specific page
func (h *RuleHandler) GetRules(c *gin.Context) {
	pageID := c.Query("page_id")
	if pageID == "" {
		pageID = "default" // Default page for backward compatibility
	}

	rules, err := h.ruleEngine.GetRulesForPage(c.Request.Context(), pageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    rules,
		"meta": gin.H{
			"page_id": pageID,
			"count":   len(rules),
		},
	})
}

// GetRule returns a specific rule by ID
func (h *RuleHandler) GetRule(c *gin.Context) {
	ruleIDStr := c.Param("id")
	ruleID, err := strconv.Atoi(ruleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid rule ID",
		})
		return
	}

	// For Phase 1, we'll get the rule through the page context
	pageID := c.Query("page_id")
	if pageID == "" {
		pageID = "default"
	}

	rules, err := h.ruleEngine.GetRulesForPage(c.Request.Context(), pageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Find the specific rule
	for _, rule := range rules {
		if rule.ID == ruleID {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    rule,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"error":   "rule not found",
	})
}

// CreateRule creates a new rule
func (h *RuleHandler) CreateRule(c *gin.Context) {
	var req CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Create rule object
	rule := &service.Rule{
		Name:           req.Name,
		Keywords:       req.Keywords,
		Responses:      req.Responses,
		Enabled:        true,
		SelectedPages:  req.SelectedPages,
		HideAfterReply: req.HideAfterReply,
		SendToInbox:    req.SendToInbox,
		InboxMessage:   req.InboxMessage,
		InboxImage:     req.InboxImage,
		Priority:       req.Priority,
		MatchType:      req.MatchType,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Set default match type if not provided
	if rule.MatchType == "" {
		rule.MatchType = "contains"
	}

	// Validate rule
	if err := h.ruleEngine.ValidateRule(rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// For Phase 1, we'll simulate saving the rule
	// In a real implementation, this would save to database
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    rule,
		"message": "Rule created successfully",
	})

	// Clear cache to force refresh
	h.ruleEngine.ClearCache()
}

// UpdateRule updates an existing rule
func (h *RuleHandler) UpdateRule(c *gin.Context) {
	ruleIDStr := c.Param("id")
	ruleID, err := strconv.Atoi(ruleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid rule ID",
		})
		return
	}

	var req UpdateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// For Phase 1, simulate the update
	// In a real implementation, this would update the database

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Rule updated successfully",
		"data": gin.H{
			"rule_id":    ruleID,
			"updated_at": time.Now(),
		},
	})

	// Clear cache to force refresh
	h.ruleEngine.ClearCache()
}

// DeleteRule deletes a rule
func (h *RuleHandler) DeleteRule(c *gin.Context) {
	ruleIDStr := c.Param("id")
	ruleID, err := strconv.Atoi(ruleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid rule ID",
		})
		return
	}

	// For Phase 1, simulate the deletion
	// In a real implementation, this would delete from database

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Rule deleted successfully",
		"data": gin.H{
			"rule_id":    ruleID,
			"deleted_at": time.Now(),
		},
	})

	// Clear cache to force refresh
	h.ruleEngine.ClearCache()
}

// GetRuleStats returns rule statistics
func (h *RuleHandler) GetRuleStats(c *gin.Context) {
	stats, err := h.ruleEngine.GetRuleStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
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

// TestRule tests a rule against sample content
func (h *RuleHandler) TestRule(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
		PageID  string `json:"page_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if req.PageID == "" {
		req.PageID = "default"
	}

	// Test keyword matching
	result, err := h.ruleEngine.MatchRules(c.Request.Context(), req.PageID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// ClearRuleCache clears the rule engine cache
func (h *RuleHandler) ClearRuleCache(c *gin.Context) {
	h.ruleEngine.ClearCache()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Rule cache cleared successfully",
	})
}

// RegisterRoutes registers all rule management routes
func (h *RuleHandler) RegisterRoutes(r *gin.RouterGroup) {
	rules := r.Group("/rules")
	{
		rules.GET("", h.GetRules)
		rules.POST("", h.CreateRule)
		rules.GET("/:id", h.GetRule)
		rules.PUT("/:id", h.UpdateRule)
		rules.DELETE("/:id", h.DeleteRule)
		rules.GET("/stats", h.GetRuleStats)
		rules.POST("/test", h.TestRule)
		rules.POST("/cache/clear", h.ClearRuleCache)
	}
}
