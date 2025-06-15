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

// GetRuleResponses returns all responses for a rule
func (h *RuleHandler) GetRuleResponses(c *gin.Context) {
	id := c.Param("id")
	ruleID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule ID"})
		return
	}

	responses, err := h.repo.GetRuleResponses(c.Request.Context(), ruleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    responses,
	})
}

// CreateRuleResponse creates a new response for a rule
func (h *RuleHandler) CreateRuleResponse(c *gin.Context) {
	id := c.Param("id")
	ruleID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule ID"})
		return
	}

	var response models.RuleResponse
	if err := c.ShouldBindJSON(&response); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response.RuleID = ruleID
	if err := h.repo.CreateRuleResponse(c.Request.Context(), &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    response,
	})
}

// UpdateRuleResponse updates a rule response
func (h *RuleHandler) UpdateRuleResponse(c *gin.Context) {
	responseID := c.Param("response_id")
	respID, err := strconv.Atoi(responseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid response ID"})
		return
	}

	var response models.RuleResponse
	if err := c.ShouldBindJSON(&response); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response.ID = respID
	if err := h.repo.UpdateRuleResponse(c.Request.Context(), &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// DeleteRuleResponse deletes a rule response
func (h *RuleHandler) DeleteRuleResponse(c *gin.Context) {
	responseID := c.Param("response_id")
	respID, err := strconv.Atoi(responseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid response ID"})
		return
	}

	if err := h.repo.DeleteRuleResponse(c.Request.Context(), respID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Response deleted successfully",
	})
}
