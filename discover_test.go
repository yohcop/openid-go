package openid

import (
	"testing"
)

func TestDiscoverWithYadis(t *testing.T) {
	// They all redirect to the same XRDS document
	expectOpIdErr(t, "http://example.com/xrds",
		"foo", identifier_select, identifier_select, false)
	expectOpIdErr(t, "http://example.com/xrds-loc",
		"foo", identifier_select, identifier_select, false)
	expectOpIdErr(t, "http://example.com/xrds-meta",
		"foo", identifier_select, identifier_select, false)
}

func TestDiscoverWithHtml(t *testing.T) {
	// Yadis discovery will fail, and fall back to html.
	expectOpIdErr(t, "http://example.com/html",
		"example.com/openid", "bar-name", "http://example.com/html",
		false)
	// The first url redirects to a different URL. The redirected-to
	// url should be used as claimedID.
	expectOpIdErr(t, "http://example.com/html-redirect",
		"example.com/openid", "bar-name", "http://example.com/html",
		false)
}

func TestDiscoverBadUrl(t *testing.T) {
	expectOpIdErr(t, "http://example.com/404", "", "", "", true)
}

func expectOpIdErr(t *testing.T, uri, exOpEndpoint, exOpLocalId, exClaimedId string, exErr bool) {
	opEndpoint, opLocalId, claimedId, err := discover(uri, testGetter)
	if (err != nil) != exErr {
		t.Errorf("Unexpected error: '%s'", err)
	} else {
		if opEndpoint != exOpEndpoint {
			t.Errorf("Extracted Endpoint does not match: Exepect %s, Got %s",
				exOpEndpoint, opEndpoint)
		}
		if opLocalId != exOpLocalId {
			t.Errorf("Extracted LocalId does not match: Exepect %s, Got %s",
				exOpLocalId, opLocalId)
		}
		if claimedId != exClaimedId {
			t.Errorf("Extracted ClaimedID does not match: Exepect %s, Got %s",
				exClaimedId, claimedId)
		}
	}
}
