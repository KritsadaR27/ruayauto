package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"chatbot/internal/models"
	"chatbot/internal/repository"
)

type FacebookService interface {
	// OAuth flow
	GenerateAuthURL(redirectURI, state string) string
	HandleCallback(code, state, redirectURI string) (*models.FacebookSession, error)
	GetAuthStatus(facebookUserID string) (*models.FacebookAuthStatus, error)
	Logout(facebookUserID string) error

	// Page management
	GetUserPages(facebookUserID string) ([]models.FacebookPageInfo, error)
	ConnectPage(facebookUserID, pageID string) error
	DisconnectPage(facebookUserID, pageID string) error
	GetConnectedPages(facebookUserID string) ([]models.FacebookPage, error)

	// Webhook management
	SetupWebhook(pageID, webhookURL, verifyToken string) error
	VerifyWebhook(pageID, verifyToken string) bool
}

type facebookService struct {
	facebookRepo repository.FacebookRepository
	appID        string
	appSecret    string
	apiVersion   string
}

func NewFacebookService(facebookRepo repository.FacebookRepository, appID, appSecret, apiVersion string) FacebookService {
	return &facebookService{
		facebookRepo: facebookRepo,
		appID:        appID,
		appSecret:    appSecret,
		apiVersion:   apiVersion,
	}
}

// OAuth flow
func (s *facebookService) GenerateAuthURL(redirectURI, state string) string {
	baseURL := "https://www.facebook.com/v19.0/dialog/oauth"
	params := url.Values{}
	params.Add("client_id", s.appID)
	params.Add("redirect_uri", redirectURI)
	params.Add("scope", "pages_manage_posts,pages_read_engagement,pages_manage_metadata")
	params.Add("response_type", "code")
	if state != "" {
		params.Add("state", state)
	}
	
	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

func (s *facebookService) HandleCallback(code, state, redirectURI string) (*models.FacebookSession, error) {
	// Exchange code for access token
	tokenURL := fmt.Sprintf("https://graph.facebook.com/v19.0/oauth/access_token")
	params := url.Values{}
	params.Add("client_id", s.appID)
	params.Add("client_secret", s.appSecret)
	params.Add("redirect_uri", redirectURI)
	params.Add("code", code)

	resp, err := http.PostForm(tokenURL, params)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token response: %w", err)
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	// Get user info
	userResp, err := http.Get(fmt.Sprintf("https://graph.facebook.com/v19.0/me?access_token=%s", tokenResponse.AccessToken))
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer userResp.Body.Close()

	userBody, err := io.ReadAll(userResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user response: %w", err)
	}

	var userInfo struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	if err := json.Unmarshal(userBody, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	// Create or update session
	var expiresAt *time.Time
	if tokenResponse.ExpiresIn > 0 {
		expiry := time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)
		expiresAt = &expiry
	}

	session := &models.FacebookSession{
		FacebookUserID: userInfo.ID,
		AccessToken:    tokenResponse.AccessToken,
		ExpiresAt:      expiresAt,
	}

	// Check if session exists
	existingSession, err := s.facebookRepo.GetSession(userInfo.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing session: %w", err)
	}

	if existingSession != nil {
		// Update existing session
		session.ID = existingSession.ID
		err = s.facebookRepo.UpdateSession(session)
	} else {
		// Create new session
		err = s.facebookRepo.CreateSession(session)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	return session, nil
}

func (s *facebookService) GetAuthStatus(facebookUserID string) (*models.FacebookAuthStatus, error) {
	session, err := s.facebookRepo.GetSession(facebookUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		return &models.FacebookAuthStatus{
			IsAuthenticated: false,
		}, nil
	}

	// Check if token is expired
	if session.ExpiresAt != nil && session.ExpiresAt.Before(time.Now()) {
		return &models.FacebookAuthStatus{
			IsAuthenticated: false,
		}, nil
	}

	// Get connected pages
	pages, err := s.facebookRepo.GetPagesByUser(facebookUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pages: %w", err)
	}

	return &models.FacebookAuthStatus{
		IsAuthenticated: true,
		UserID:          facebookUserID,
		Pages:           pages,
	}, nil
}

func (s *facebookService) Logout(facebookUserID string) error {
	return s.facebookRepo.DeleteSession(facebookUserID)
}

// Page management
func (s *facebookService) GetUserPages(facebookUserID string) ([]models.FacebookPageInfo, error) {
	session, err := s.facebookRepo.GetSession(facebookUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		return nil, fmt.Errorf("user not authenticated")
	}

	// Get pages from Facebook API
	pagesURL := fmt.Sprintf("https://graph.facebook.com/v19.0/me/accounts?access_token=%s", session.AccessToken)
	resp, err := http.Get(pagesURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get pages from Facebook: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read pages response: %w", err)
	}

	var pagesResponse models.FacebookPagesResponse
	if err := json.Unmarshal(body, &pagesResponse); err != nil {
		return nil, fmt.Errorf("failed to parse pages response: %w", err)
	}

	return pagesResponse.Data, nil
}

func (s *facebookService) ConnectPage(facebookUserID, pageID string) error {
	// Get user pages from Facebook
	pages, err := s.GetUserPages(facebookUserID)
	if err != nil {
		return fmt.Errorf("failed to get user pages: %w", err)
	}

	// Find the page
	var pageInfo *models.FacebookPageInfo
	for _, page := range pages {
		if page.ID == pageID {
			pageInfo = &page
			break
		}
	}

	if pageInfo == nil {
		return fmt.Errorf("page not found or user doesn't have access")
	}

	// Create or update page in database
	page := &models.FacebookPage{
		PageID:           pageInfo.ID,
		FacebookPageID:   &pageInfo.ID,
		PageName:         pageInfo.Name,
		Name:             &pageInfo.Name,
		AccessToken:      &pageInfo.AccessToken,
		FacebookUserID:   &facebookUserID,
		IsActive:         true,
		Connected:        &[]bool{true}[0],
		Enabled:          &[]bool{true}[0],
		WebhookVerified:  &[]bool{false}[0],
	}

	return s.facebookRepo.CreatePage(page)
}

func (s *facebookService) DisconnectPage(facebookUserID, pageID string) error {
	return s.facebookRepo.SetPageConnectionStatus(pageID, false)
}

func (s *facebookService) GetConnectedPages(facebookUserID string) ([]models.FacebookPage, error) {
	return s.facebookRepo.GetPagesByUser(facebookUserID)
}

// Webhook management
func (s *facebookService) SetupWebhook(pageID, webhookURL, verifyToken string) error {
	page, err := s.facebookRepo.GetPage(pageID)
	if err != nil {
		return fmt.Errorf("failed to get page: %w", err)
	}

	if page == nil || page.AccessToken == nil {
		return fmt.Errorf("page not found or no access token")
	}

	// Subscribe to webhook
	subscribeURL := fmt.Sprintf("https://graph.facebook.com/v19.0/%s/subscribed_apps", pageID)
	params := url.Values{}
	params.Add("access_token", *page.AccessToken)
	params.Add("subscribed_fields", "feed,mention,message_deliveries,messaging_optouts,messaging_postbacks,messaging_pre_checkouts")

	resp, err := http.PostForm(subscribeURL, params)
	if err != nil {
		return fmt.Errorf("failed to subscribe to webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook subscription failed: %s", string(body))
	}

	// Save webhook subscription
	subscription := &models.FacebookWebhookSubscription{
		FacebookPageID: pageID,
		WebhookURL:     webhookURL,
		VerifyToken:    verifyToken,
		Fields:         []string{"comments", "mentions", "messages"},
		Active:         true,
	}

	return s.facebookRepo.CreateWebhookSubscription(subscription)
}

func (s *facebookService) VerifyWebhook(pageID, verifyToken string) bool {
	subscription, err := s.facebookRepo.GetWebhookSubscription(pageID)
	if err != nil {
		return false
	}

	if subscription == nil {
		return false
	}

	return subscription.VerifyToken == verifyToken && subscription.Active
}
