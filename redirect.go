package openid

import (
	"net/url"
	"strings"
)

func RedirectURL(id, callbackURL, realm string) (string, error) {
	return redirectURL(id, callbackURL, realm, urlGetter)
}

func redirectURL(id, callbackURL, realm string, getter httpGetter) (string, error) {
	opEndpoint, opLocalID, claimedID, err := discover(id, getter)
	if err != nil {
		return "", err
	}
	return buildRedirectURL(opEndpoint, opLocalID, claimedID, callbackURL, realm)
}

func buildRedirectURL(opEndpoint, opLocalID, claimedID, returnTo, realm string) (string, error) {
	values := make(url.Values)
	values.Add("openid.ns", "http://specs.openid.net/auth/2.0")
	values.Add("openid.mode", "checkid_setup")
	values.Add("openid.return_to", returnTo)

	if len(claimedID) > 0 {
		values.Add("openid.claimed_id", claimedID)
		if len(opLocalID) > 0 {
			values.Add("openid.identity", opLocalID)
		} else {
			values.Add("openid.identity",
				"http://specs.openid.net/auth/2.0/identifier_select")
		}
	} else {
		values.Add("openid.identity",
			"http://specs.openid.net/auth/2.0/identifier_select")
	}

	if len(realm) > 0 {
		values.Add("openid.realm", realm)
	}

	if strings.Contains(opEndpoint, "?") {
		return opEndpoint + "&" + values.Encode(), nil
	}
	return opEndpoint + "?" + values.Encode(), nil
}
