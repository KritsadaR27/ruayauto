package handlers

import (
	"net/http"
	"strconv"

	"chatbot/internal/models"
	"chatbot/internal/repository"

	"github.com/gin-gonic/gin"
)

type RuleHandler struct {
	repo *repository.RuleRepository
}

func NewRuleHandler(repo *repository.RuleRepository) *RuleHandler {
	return &RuleHandler{repo: repo}
}

// GetRules returns all active rules
func (h *RuleHandler) GetRules(c *gin.Context) {
	rules, err := h.repo.GetActive(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"pairs": rules,
		},
		"timestamp": "now",
	})
}

// GetRuleByID returns a rule by its ID
func (h *RuleHandler) GetRuleByID(c *gin.Context) {
	id := c.Param("id")
	ruleID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule ID"})
		return
	}

	rule, err := h.repo.GetByID(c.Request.Context(), ruleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rule == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"pair": rule,
		},
		"timestamp": "now",
	})
}

// CreateRule creates a new rule
func (h *RuleHandler) CreateRule(c *gin.Context) {
	var input models.Rule
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.Create(c.Request.Context(), &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": input})
}

// UpdateRule updates an existing rule
func (h *RuleHandler) UpdateRule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule ID"})
		return
	}
	var input models.Rule
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.ID = id
	if err := h.repo.Update(c.Request.Context(), &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": input})
}

// DeleteRule deletes a rule
func (h *RuleHandler) DeleteRule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule ID"})
		return
	}
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "rule deleted successfully"})
}
