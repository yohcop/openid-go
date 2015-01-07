package integration

// These tests fetch real data from google.com and other OpenID
// providers. If they change the files returned, or endpoints, or
// whatever, they will fail. It's ok though, they are full tests.

import (
	. "openid"
	"testing"
)

func TestGoogleCom(t *testing.T) {
	expectDiscovery(t, "https://www.google.com/accounts/o8/id",
		"https://www.google.com/accounts/o8/ud",
		"http://specs.openid.net/auth/2.0/identifier_select",
		"http://specs.openid.net/auth/2.0/identifier_select")
}

func TestYahoo(t *testing.T) {
	expectDiscovery(t, "https://me.yahoo.com",
		"https://open.login.yahooapis.com/openid/op/auth",
		"http://specs.openid.net/auth/2.0/identifier_select",
		"http://specs.openid.net/auth/2.0/identifier_select")
}

func TestYohcop(t *testing.T) {
	expectDiscovery(t, "http://yohcop.net",
		"https://www.google.com/accounts/o8/ud?source=profiles",
		"http://www.google.com/profiles/yohcop",
		"http://yohcop.net")
}

func expectDiscovery(t *testing.T, uri, expectOp, expectLocalId, expectClaimedId string) {
	endpoint, localId, claimedId, err := Discover(uri)
	if err != nil {
		t.Errorf("Discovery failed")
	}
	if endpoint != expectOp {
		t.Errorf("Unexpected endpoint: %s", endpoint)
	}
	if localId != expectLocalId {
		t.Errorf("Unexpected localId: %s", localId)
	}
	if claimedId != expectClaimedId {
		t.Errorf("Unexpected claimedId: %s", claimedId)
	}
}
