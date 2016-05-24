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

func TestVerifyDiscovered(t *testing.T) {
	dc := NewSimpleDiscoveryCache()
	vals := url.Values{"openid.ns": []string{"http://specs.openid.net/auth/2.0"},
		"openid.mode":        []string{"id_res"},
		"openid.op_endpoint": []string{"http://example.com/openid/login"},
		"openid.claimed_id":  []string{"http://example.com/openid/id/foo"},
		"openid.identity":    []string{"http://example.com/openid/id/foo"}}

	// Make sure we fail with no discovery handler
	if err := testInstance.verifyDiscovered(nil, vals, dc); err == nil {
		t.Errorf("verifyDiscovered succeeded unexpectedly with no discovery")
	}

	// Add the discovery handler
	testGetter.urls["http://example.com/openid/id/foo#Accept#application/xrds+xml"] = `HTTP/1.0 200 OK
Content-Type: application/xrds+xml; charset=UTF-8

<?xml version="1.0" encoding="UTF-8"?>
<xrds:XRDS xmlns:xrds="xri://$xrds" xmlns="xri://$xrd*($v*2.0)">
	<XRD>
		<Service priority="0">
			<Type>http://specs.openid.net/auth/2.0/signon</Type>		
			<URI>http://example.com/openid/login</URI>
		</Service>
	</XRD>
</xrds:XRDS>`

	// Make sure we succeed now
	if err := testInstance.verifyDiscovered(nil, vals, dc); err != nil {
		t.Errorf("verifyDiscovered failed unexpectedly: %v", err)
	}

	// Remove the discovery handler
	delete(testGetter.urls, "http://example.com/openid/id/foo#Accept#application/xrds+xml")

	// Make sure we still succeed thanks to the discovery cache
	if err := testInstance.verifyDiscovered(nil, vals, dc); err != nil {
		t.Errorf("verifyDiscovered failed unexpectedly: %v", err)
	}
}
