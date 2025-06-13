package line

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

// VerifyWebhook verifies LINE webhook signature
func VerifyWebhook(payload []byte, signature string, channelSecret string) error {
	if signature == "" {
		return errors.New("missing signature")
	}

	// Calculate expected signature
	mac := hmac.New(sha256.New, []byte(channelSecret))
	mac.Write(payload)
	expectedSig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	
	// Compare signatures
	if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
		return errors.New("signature verification failed")
	}
	
	return nil
}
