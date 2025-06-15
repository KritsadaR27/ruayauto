package facebook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// VerifyWebhook verifies Facebook webhook signature
func VerifyWebhook(payload []byte, signature string, appSecret string) error {
	if signature == "" {
		return errors.New("missing signature")
	}

	// Remove 'sha256=' prefix
	if !strings.HasPrefix(signature, "sha256=") {
		return errors.New("invalid signature format")
	}

	expectedSig := strings.TrimPrefix(signature, "sha256=")

	// Calculate expected signature
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write(payload)
	calculatedSig := hex.EncodeToString(mac.Sum(nil))

	// Compare signatures
	if !hmac.Equal([]byte(expectedSig), []byte(calculatedSig)) {
		return fmt.Errorf("signature verification failed")
	}

	return nil
}

// VerifyToken verifies the webhook verification challenge
func VerifyToken(challengeToken, verifyToken string) error {
	if challengeToken != verifyToken {
		return errors.New("verify token mismatch")
	}
	return nil
}
