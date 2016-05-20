package openid

import (
	"net/url"
	"testing"
)

func TestBuildRedirectUrl(t *testing.T) {
	expectURL(t, "https://endpoint/a", "opLocalId", "claimedId", "returnTo", "realm",
		"https://endpoint/a?"+
			"openid.ns=http://specs.openid.net/auth/2.0"+
			"&openid.mode=checkid_setup"+
			"&openid.return_to=returnTo"+
			"&openid.claimed_id=claimedId"+
			"&openid.identity=opLocalId"+
			"&openid.realm=realm")
	// No realm.
	expectURL(t, "https://endpoint/a", "opLocalId", "claimedId", "returnTo", "",
		"https://endpoint/a?"+
			"openid.ns=http://specs.openid.net/auth/2.0"+
			"&openid.mode=checkid_setup"+
			"&openid.return_to=returnTo"+
			"&openid.claimed_id=claimedId"+
			"&openid.identity=opLocalId")
	// No realm, no localId
	expectURL(t, "https://endpoint/a", "", "claimedId", "returnTo", "",
		"https://endpoint/a?"+
			"openid.ns=http://specs.openid.net/auth/2.0"+
			"&openid.mode=checkid_setup"+
			"&openid.return_to=returnTo"+
			"&openid.claimed_id=claimedId"+
			"&openid.identity=claimedId")
	// No realm, no claimedId
	expectURL(t, "https://endpoint/a", "opLocalId", "", "returnTo", "",
		"https://endpoint/a?"+
			"openid.ns=http://specs.openid.net/auth/2.0"+
			"&openid.mode=checkid_setup"+
			"&openid.return_to=returnTo"+
			"&openid.claimed_id="+
			"http://specs.openid.net/auth/2.0/identifier_select"+
			"&openid.identity="+
			"http://specs.openid.net/auth/2.0/identifier_select")
}

func expectURL(t *testing.T, opEndpoint, opLocalID, claimedID, returnTo, realm, expected string) {
	url, err := BuildRedirectURL(opEndpoint, opLocalID, claimedID, returnTo, realm)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	compareUrls(t, url, expected)
}

func TestRedirectWithDiscovery(t *testing.T) {
	expected := "foo?" +
		"openid.ns=http://specs.openid.net/auth/2.0" +
		"&openid.mode=checkid_setup" +
		"&openid.return_to=mysite/cb" +
		"&openid.claimed_id=" +
		"http://specs.openid.net/auth/2.0/identifier_select" +
		"&openid.identity=" +
		"http://specs.openid.net/auth/2.0/identifier_select"

	// They all redirect to the same XRDS document
	expectRedirect(t, "http://example.com/xrds",
		"mysite/cb", "", expected, false)
	expectRedirect(t, "http://example.com/xrds-loc",
		"mysite/cb", "", expected, false)
	expectRedirect(t, "http://example.com/xrds-meta",
		"mysite/cb", "", expected, false)
}

func expectRedirect(t *testing.T, uri, callback, realm, exRedirect string, exErr bool) {
	redirect, err := testInstance.RedirectURL(uri, callback, realm)
	if (err != nil) != exErr {
		t.Errorf("Unexpected error: '%s'", err)
		return
	}
	compareUrls(t, redirect, exRedirect)
}

func compareUrls(t *testing.T, url1, expected string) {
	p1, err1 := url.Parse(url1)
	p2, err2 := url.Parse(expected)
	if err1 != nil {
		t.Errorf("Url1 non parsable: %s", err1)
		return
	}
	if err2 != nil {
		t.Errorf("ExpectedUrl non parsable: %s", err2)
		return
	}
	if p1.Scheme != p2.Scheme ||
		p1.Host != p2.Host ||
		p1.Path != p2.Path {
		t.Errorf("URLs don't match: %s vs %s", url1, expected)
	}
	q1, _ := url.ParseQuery(p1.RawQuery)
	q2, _ := url.ParseQuery(p2.RawQuery)
	if err := compareQueryParams(q1, q2); err != nil {
		t.Errorf("URLs query params don't match: %s: %s vs %s", err, url1, expected)
	}
}
