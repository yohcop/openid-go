package openid

import (
	"net/url"
	"testing"
	"time"
)

func TestVerifyNonce(t *testing.T) {
	timeStr := time.Now().UTC().Format(time.RFC3339)
	ns := &SimpleNonceStore{Store: make(map[string][]*Nonce)}
	v := url.Values{}

	// Initial values
	v.Set("openid.op_endpoint", "1")
	v.Set("openid.response_nonce", timeStr+"foo")
	if err := verifyNonce(v, ns); err != nil {
		t.Errorf("verifyNonce failed unexpectedly: %v", err)
	}

	// Different nonce
	v.Set("openid.response_nonce", timeStr+"bar")
	if err := verifyNonce(v, ns); err != nil {
		t.Errorf("verifyNonce failed unexpectedly: %v", err)
	}

	// Different endpoint
	v.Set("openid.op_endpoint", "2")
	if err := verifyNonce(v, ns); err != nil {
		t.Errorf("verifyNonce failed unexpectedly: %v", err)
	}
}
