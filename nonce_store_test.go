package openid

import (
	"testing"
	"time"
)

func TestDefaultNonceStore(t *testing.T) {
	*maxNonceAge = 60 * time.Second
	now := time.Now().UTC()
	// 30 seconds ago
	now30s := now.Add(-30 * time.Second)
	// 2 minutes ago
	now2m := now.Add(-2 * time.Minute)

	now30sStr := now30s.Format(time.RFC3339)
	now2mStr := now2m.Format(time.RFC3339)

	ns := NewSimpleNonceStore()
	reject(t, ns, "1", "foo")                        // invalid nonce
	reject(t, ns, "1", "fooBarBazLongerThan20Chars") // invalid nonce

	accept(t, ns, "1", now30sStr+"asd")
	reject(t, ns, "1", now30sStr+"asd") // same nonce
	accept(t, ns, "1", now30sStr+"xxx") // different nonce
	reject(t, ns, "1", now30sStr+"xxx") // different nonce again to verify storage of multiple nonces per endpoint
	accept(t, ns, "2", now30sStr+"asd") // different endpoint

	reject(t, ns, "1", now2mStr+"old") // too old
	reject(t, ns, "3", now2mStr+"old") // too old
}

func accept(t *testing.T, ns NonceStore, op, nonce string) {
	e := ns.Accept(op, nonce)
	if e != nil {
		t.Errorf("Should accept %s nonce %s", op, nonce)
	}
}

func reject(t *testing.T, ns NonceStore, op, nonce string) {
	e := ns.Accept(op, nonce)
	if e == nil {
		t.Errorf("Should reject %s nonce %s", op, nonce)
	}
}
