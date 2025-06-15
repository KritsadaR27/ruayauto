package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"ruaychatbot/backend/internal/ChatBotCore/channels"
	"ruaychatbot/backend/internal/ChatBotCore/repository"
)

// PageManager handles page configuration and management
type PageManager struct {
	conversationRepo repository.ConversationRepository
	messageRepo      repository.MessageRepository

	// In-memory page registry (future: move to database)
	pages map[string]*PageConfig
}

// PageConfig represents a configured page/channel
type PageConfig struct {
	ID          string        `json:"id"`
	Platform    string        `json:"platform"` // "facebook", "line", "tiktok"
	PageID      string        `json:"page_id"`  // Platform-specific page ID
	PageName    string        `json:"page_name"`
	AccessToken string        `json:"access_token,omitempty"`
	WebhookURL  string        `json:"webhook_url,omitempty"`
	IsEnabled   bool          `json:"is_enabled"`
	IsConnected bool          `json:"is_connected"`
	Settings    *PageSettings `json:"settings,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// PageSettings represents page-specific settings
type PageSettings struct {
	AutoReply              bool   `json:"auto_reply"`
	HideCommentsAfterReply bool   `json:"hide_comments_after_reply"`
	SendToInbox            bool   `json:"send_to_inbox"`
	WelcomeMessage         string `json:"welcome_message,omitempty"`
	FallbackMessage        string `json:"fallback_message,omitempty"`

	// Rate limiting
	MaxRepliesPerMinute int `json:"max_replies_per_minute"`

	// Advanced features
	IgnoreMentions bool          `json:"ignore_mentions"`
	IgnoreStickers bool          `json:"ignore_stickers"`
	WorkingHours   *WorkingHours `json:"working_hours,omitempty"`
}

// WorkingHours represents page working hours
type WorkingHours struct {
	Enabled   bool   `json:"enabled"`
	StartTime string `json:"start_time"` // "09:00"
	EndTime   string `json:"end_time"`   // "18:00"
	Timezone  string `json:"timezone"`   // "Asia/Bangkok"
	WeekDays  []int  `json:"week_days"`  // 1=Monday, 7=Sunday
}

// NewPageManager creates a new page manager
func NewPageManager(repos *repository.Repositories) *PageManager {
	pm := &PageManager{
		conversationRepo: repos.Conversation,
		messageRepo:      repos.Message,
		pages:            make(map[string]*PageConfig),
	}

	// Initialize with default Facebook page (backward compatibility)
	pm.initializeDefaultPages()

	return pm
}

// initializeDefaultPages sets up default page configurations
func (pm *PageManager) initializeDefaultPages() {
	// Default Facebook page for backward compatibility
	defaultPage := &PageConfig{
		ID:          "fb_default",
		Platform:    "facebook",
		PageID:      "default",
		PageName:    "Default Facebook Page",
		IsEnabled:   true,
		IsConnected: true,
		Settings: &PageSettings{
			AutoReply:              true,
			HideCommentsAfterReply: false,
			SendToInbox:            false,
			FallbackMessage:        "ขอบคุณสำหรับข้อความครับ เราจะตอบกลับโดยเร็วที่สุด",
			MaxRepliesPerMinute:    30,
			IgnoreMentions:         false,
			IgnoreStickers:         false,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	pm.pages[defaultPage.ID] = defaultPage
	log.Printf("📄 Initialized default page: %s", defaultPage.PageName)
}

// RegisterChannel registers a new channel configuration
func (pm *PageManager) RegisterChannel(config channels.ChannelConfig) error {
	pageConfig := &PageConfig{
		ID:          fmt.Sprintf("%s_%s", config.Platform, config.Name),
		Platform:    config.Platform,
		PageID:      config.Name, // For now, use name as page ID
		PageName:    config.Name,
		IsEnabled:   config.Enabled,
		IsConnected: true,
		Settings: &PageSettings{
			AutoReply:           true,
			MaxRepliesPerMinute: 30,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	pm.pages[pageConfig.ID] = pageConfig
	log.Printf("📄 Registered page: %s (%s)", pageConfig.PageName, pageConfig.Platform)

	return nil
}

// GetPage retrieves a page configuration by ID
func (pm *PageManager) GetPage(pageID string) (*PageConfig, error) {
	// First, try exact match
	if page, exists := pm.pages[pageID]; exists {
		return page, nil
	}

	// For backward compatibility, try to find by page ID
	for _, page := range pm.pages {
		if page.PageID == pageID {
			return page, nil
		}
	}

	return nil, fmt.Errorf("page not found: %s", pageID)
}

// GetPagesByPlatform retrieves all pages for a specific platform
func (pm *PageManager) GetPagesByPlatform(platform string) []*PageConfig {
	var pages []*PageConfig
	for _, page := range pm.pages {
		if page.Platform == platform {
			pages = append(pages, page)
		}
	}
	return pages
}

// GetAllPages retrieves all configured pages
func (pm *PageManager) GetAllPages() []*PageConfig {
	var pages []*PageConfig
	for _, page := range pm.pages {
		pages = append(pages, page)
	}
	return pages
}

// GetEnabledPages retrieves all enabled pages
func (pm *PageManager) GetEnabledPages() []*PageConfig {
	var pages []*PageConfig
	for _, page := range pm.pages {
		if page.IsEnabled {
			pages = append(pages, page)
		}
	}
	return pages
}

// GetConnectedPages retrieves all connected pages
func (pm *PageManager) GetConnectedPages() []*PageConfig {
	var pages []*PageConfig
	for _, page := range pm.pages {
		if page.IsConnected {
			pages = append(pages, page)
		}
	}
	return pages
}

// UpdatePage updates an existing page configuration
func (pm *PageManager) UpdatePage(pageID string, updates *PageConfig) error {
	page, err := pm.GetPage(pageID)
	if err != nil {
		return err
	}

	// Update fields
	if updates.PageName != "" {
		page.PageName = updates.PageName
	}

	if updates.AccessToken != "" {
		page.AccessToken = updates.AccessToken
	}

	page.IsEnabled = updates.IsEnabled
	page.IsConnected = updates.IsConnected

	if updates.Settings != nil {
		page.Settings = updates.Settings
	}

	page.UpdatedAt = time.Now()

	log.Printf("📄 Updated page: %s", page.PageName)
	return nil
}

// EnablePage enables a page
func (pm *PageManager) EnablePage(pageID string) error {
	page, err := pm.GetPage(pageID)
	if err != nil {
		return err
	}

	page.IsEnabled = true
	page.UpdatedAt = time.Now()

	log.Printf("✅ Enabled page: %s", page.PageName)
	return nil
}

// DisablePage disables a page
func (pm *PageManager) DisablePage(pageID string) error {
	page, err := pm.GetPage(pageID)
	if err != nil {
		return err
	}

	page.IsEnabled = false
	page.UpdatedAt = time.Now()

	log.Printf("❌ Disabled page: %s", page.PageName)
	return nil
}

// DeletePage removes a page configuration
func (pm *PageManager) DeletePage(pageID string) error {
	if _, err := pm.GetPage(pageID); err != nil {
		return err
	}

	delete(pm.pages, pageID)
	log.Printf("🗑️ Deleted page: %s", pageID)
	return nil
}

// IsPageEnabled checks if a page is enabled
func (pm *PageManager) IsPageEnabled(pageID string) bool {
	page, err := pm.GetPage(pageID)
	if err != nil {
		return false
	}
	return page.IsEnabled && page.IsConnected
}

// GetPageSettings retrieves settings for a specific page
func (pm *PageManager) GetPageSettings(pageID string) (*PageSettings, error) {
	page, err := pm.GetPage(pageID)
	if err != nil {
		return nil, err
	}

	if page.Settings == nil {
		// Return default settings
		return &PageSettings{
			AutoReply:           true,
			MaxRepliesPerMinute: 30,
			FallbackMessage:     "ขอบคุณสำหรับข้อความครับ เราจะตอบกลับโดยเร็วที่สุด",
		}, nil
	}

	return page.Settings, nil
}

// UpdatePageSettings updates settings for a specific page
func (pm *PageManager) UpdatePageSettings(pageID string, settings *PageSettings) error {
	page, err := pm.GetPage(pageID)
	if err != nil {
		return err
	}

	page.Settings = settings
	page.UpdatedAt = time.Now()

	log.Printf("⚙️ Updated settings for page: %s", page.PageName)
	return nil
}

// GetPageStats retrieves statistics for a specific page
func (pm *PageManager) GetPageStats(ctx context.Context, pageID string) (map[string]interface{}, error) {
	page, err := pm.GetPage(pageID)
	if err != nil {
		return nil, err
	}

	// Get conversation count for this page
	// Note: This is a simplified implementation
	// In a real scenario, you'd have proper page-conversation mapping

	return map[string]interface{}{
		"page_id":        page.PageID,
		"page_name":      page.PageName,
		"platform":       page.Platform,
		"is_enabled":     page.IsEnabled,
		"is_connected":   page.IsConnected,
		"created_at":     page.CreatedAt,
		"updated_at":     page.UpdatedAt,
		"conversations":  0, // TODO: Implement proper counting
		"messages_today": 0, // TODO: Implement proper counting
	}, nil
}

// GetManagerStats retrieves overall page manager statistics
func (pm *PageManager) GetManagerStats() map[string]interface{} {
	totalPages := len(pm.pages)
	enabledPages := len(pm.GetEnabledPages())
	connectedPages := len(pm.GetConnectedPages())

	platformCounts := make(map[string]int)
	for _, page := range pm.pages {
		platformCounts[page.Platform]++
	}

	return map[string]interface{}{
		"total_pages":     totalPages,
		"enabled_pages":   enabledPages,
		"connected_pages": connectedPages,
		"platforms":       platformCounts,
	}
}

// ValidatePageConfig validates a page configuration
func (pm *PageManager) ValidatePageConfig(config *PageConfig) error {
	if config.Platform == "" {
		return fmt.Errorf("platform is required")
	}

	if config.PageID == "" {
		return fmt.Errorf("page ID is required")
	}

	if config.PageName == "" {
		return fmt.Errorf("page name is required")
	}

	return nil
}

// AddPage adds a new page configuration
func (pm *PageManager) AddPage(config *PageConfig) error {
	if err := pm.ValidatePageConfig(config); err != nil {
		return fmt.Errorf("invalid page config: %w", err)
	}

	// Generate ID if not provided
	if config.ID == "" {
		config.ID = fmt.Sprintf("%s_%s", config.Platform, config.PageID)
	}

	// Check if page already exists
	if _, exists := pm.pages[config.ID]; exists {
		return fmt.Errorf("page already exists: %s", config.ID)
	}

	// Set timestamps
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	// Set default settings if not provided
	if config.Settings == nil {
		config.Settings = &PageSettings{
			AutoReply:           true,
			MaxRepliesPerMinute: 30,
			FallbackMessage:     "ขอบคุณสำหรับข้อความครับ เราจะตอบกลับโดยเร็วที่สุด",
		}
	}

	pm.pages[config.ID] = config
	log.Printf("➕ Added page: %s (%s)", config.PageName, config.Platform)

	return nil
}
