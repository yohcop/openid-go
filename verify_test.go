package openid

import (
	"net/url"
	"testing"
	"time"
)

func TestVerifyNonce(t *testing.T) {
	timeStr := time.Now().UTC().Format(time.RFC3339)
	ns := NewSimpleNonceStore()
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

func TestVerifySignedFields(t *testing.T) {
	// No claimed_id/identity, properly signed
	doVerifySignedFields(t,
		url.Values{"openid.signed": []string{"signed,op_endpoint,return_to,response_nonce,assoc_handle"}},
		true)

	// Everything properly signed, even empty claimed_id/identity
	doVerifySignedFields(t,
		url.Values{"openid.signed": []string{"signed,op_endpoint,claimed_id,identity,return_to,response_nonce,assoc_handle"}},
		true)

	// With claimed_id/identity, properly signed
	doVerifySignedFields(t,
		url.Values{"openid.signed": []string{"signed,op_endpoint,claimed_id,identity,return_to,response_nonce,assoc_handle"},
			"openid.claimed_id": []string{"foo"},
			"openid.identity":   []string{"foo"}},
		true)

	// With claimed_id/identity, but those two not signed
	doVerifySignedFields(t,
		url.Values{"openid.signed": []string{"signed,op_endpoint,return_to,response_nonce,assoc_handle"},
			"openid.claimed_id": []string{"foo"},
			"openid.identity":   []string{"foo"}},
		false)

	// Missing signature for op_endpoint
	doVerifySignedFields(t,
		url.Values{"openid.signed": []string{"signed,claimed_id,identity,return_to,response_nonce,assoc_handle"},
			"openid.claimed_id": []string{"foo"},
			"openid.identity":   []string{"foo"}},
		false)

	// Missing signature for return_to
	doVerifySignedFields(t,
		url.Values{"openid.signed": []string{"signed,op_endpoint,claimed_id,identity,response_nonce,assoc_handle"},
			"openid.claimed_id": []string{"foo"},
			"openid.identity":   []string{"foo"}},
		false)

	// Missing signature for response_nonce
	doVerifySignedFields(t,
		url.Values{"openid.signed": []string{"signed,op_endpoint,claimed_id,identity,return_to,assoc_handle"},
			"openid.claimed_id": []string{"foo"},
			"openid.identity":   []string{"foo"}},
		false)

	// Missing signature for assoc_handle
	doVerifySignedFields(t,
		url.Values{"openid.signed": []string{"signed,op_endpoint,claimed_id,identity,return_to,response_nonce"},
			"openid.claimed_id": []string{"foo"},
			"openid.identity":   []string{"foo"}},
		false)
}

func doVerifySignedFields(t *testing.T, v url.Values, succeed bool) {
	if err := verifySignedFields(v); err == nil && !succeed {
		t.Errorf("verifySignedFields succeeded unexpectedly: %v - %v", v, err)
	} else if err != nil && succeed {
		t.Errorf("verifySignedFields failed unexpectedly: %v - %v", v, err)
	}
}
