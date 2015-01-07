package openid

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
)

func Verify(uri string, cache DiscoveryCache, nonceStore NonceStore) (id string, err error) {
	return verify(uri, cache, urlGetter, nonceStore)
}

func verify(uri string, cache DiscoveryCache, getter httpGetter, nonceStore NonceStore) (id string, err error) {
	parsedUrl, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	values, err := url.ParseQuery(parsedUrl.RawQuery)
	if err != nil {
		return "", err
	}

	// 11.  Verifying Assertions
	// When the Relying Party receives a positive assertion, it MUST
	// verify the following before accepting the assertion:

	// - The value of "openid.return_to" matches the URL of the current
	//   request (Section 11.1)
	if err = verifyReturnTo(parsedUrl, values); err != nil {
		return "", err
	}

	// - Discovered information matches the information in the assertion
	//   (Section 11.2)
	if err = verifyDiscovered(parsedUrl, values, cache, getter); err != nil {
		return "", err
	}

	// - An assertion has not yet been accepted from this OP with the
	//   same value for "openid.response_nonce" (Section 11.3)
	if err = verifyNonce(values, nonceStore); err != nil {
		return "", err
	}

	// - The signature on the assertion is valid and all fields that are
	//   required to be signed are signed (Section 11.4)
	if err = verifySignature(uri, values, getter); err != nil {
		return "", err
	}

	// If all four of these conditions are met, assertion is now
	// verified. If the assertion contained a Claimed Identifier, the
	// user is now authenticated with that identifier.
	return values.Get("openid.claimed_id"), nil
}

// 11.1.  Verifying the Return URL
// To verify that the "openid.return_to" URL matches the URL that is processing this assertion:
// - The URL scheme, authority, and path MUST be the same between the two
//   URLs.
// - Any query parameters that are present in the "openid.return_to" URL
//   MUST also be present with the same values in the URL of the HTTP
//   request the RP received.
func verifyReturnTo(uri *url.URL, vals url.Values) error {
	returnTo := vals.Get("openid.return_to")
	rp, err := url.Parse(returnTo)
	if err != nil {
		return err
	}
	if uri.Scheme != rp.Scheme ||
		uri.Host != rp.Host ||
		uri.Path != rp.Path {
		return errors.New(
			"Scheme, host or path don't match in return_to URL")
	}
	qp, err := url.ParseQuery(rp.RawQuery)
	if err != nil {
		return err
	}
	return compareQueryParams(qp, vals)
}

// Any parameter in q1 must also be present in q2, and values must match.
func compareQueryParams(q1, q2 url.Values) error {
	for k := range q1 {
		v1 := q1.Get(k)
		v2 := q2.Get(k)
		if v1 != v2 {
			return fmt.Errorf(
				"URLs query params don't match: Param %s different: %s vs %s",
				k, v1, v2)
		}
	}
	return nil
}

func verifyDiscovered(uri *url.URL, vals url.Values, cache DiscoveryCache, getter httpGetter) error {
	version := vals.Get("openid.ns")
	if version != "http://specs.openid.net/auth/2.0" {
		return errors.New("Bad protocol version")
	}

	endpoint := vals.Get("openid.op_endpoint")
	if len(endpoint) == 0 {
		return errors.New("missing openid.op_endpoint url param")
	}
	localId := vals.Get("openid.identity")
	if len(localId) == 0 {
		return errors.New("no localId to verify")
	}
	claimedId := vals.Get("openid.claimed_id")
	if len(claimedId) == 0 {
		// If no Claimed Identifier is present in the response, the
		// assertion is not about an identifier and the RP MUST NOT use the
		// User-supplied Identifier associated with the current OpenID
		// authentication transaction to identify the user. Extension
		// information in the assertion MAY still be used.
		// --- This library does not support this case. So claimed
		//     identifier must be present.
		return errors.New("no claimed_id to verify")
	}

	// 11.2.  Verifying Discovered Information

	// If the Claimed Identifier in the assertion is a URL and contains a
	// fragment, the fragment part and the fragment delimiter character "#"
	// MUST NOT be used for the purposes of verifying the discovered
	// information.
	claimedIdVerify := claimedId
	if fragmentIndex := strings.Index(claimedId, "#"); fragmentIndex != -1 {
		claimedIdVerify = claimedId[0:fragmentIndex]
	}

	discovered := cache.Get(endpoint)
	discoveredClaimedId := ""
	discoveredLocalId := ""
	if discovered != nil {
		// If the Claimed Identifier is included in the assertion, it
		// MUST have been discovered by the Relying Party and the
		// information in the assertion MUST be present in the
		// discovered information. The Claimed Identifier MUST NOT be an
		// OP Identifier.
		discoveredClaimedId = discovered.ClaimedId()
		discoveredLocalId = discovered.OpLocalId()
	}

	// If the Claimed Identifier was not previously discovered by the
	// Relying Party (the "openid.identity" in the request was
	// "http://specs.openid.net/auth/2.0/identifier_select" or a different
	// Identifier, or if the OP is sending an unsolicited positive
	// assertion), the Relying Party MUST perform discovery on the Claimed
	// Identifier in the response to make sure that the OP is authorized to
	// make assertions about the Claimed Identifier.
	if discovered == nil ||
		discoveredClaimedId ==
			"http://specs.openid.net/auth/2.0/identifier_select" ||
		claimedIdVerify != discoveredClaimedId {
		if ep, lid, dci, err := discover(claimedId, getter); err == nil {
			discoveredClaimedId = dci
			discoveredLocalId = lid

			if ep == endpoint {
				// This claimed ID points to the same endpoint, therefore this
				// endpoint is authorized to make assertions about that claimed ID.
				// TODO: There may be multiple undpoints found during discovery.
				// They should all be checked.
				return nil
			}
		}
		return errors.New("Could not verify that the claimed ID")
	}

	// we know here that claimedIdVerify == discoveredClaimedId.
	if localId != discoveredLocalId {
		return errors.New("Discovered local ID does not match identity")
	}
	return nil
}

func verifyNonce(vals url.Values, store NonceStore) error {
	nonce := vals.Get("openid.response_nonce")
	endpoint := vals.Get("openid.endpoint")
	return store.Accept(endpoint, nonce)
}

func verifySignature(uri string, vals url.Values, getter httpGetter) error {
	// To have the signature verification performed by the OP, the
	// Relying Party sends a direct request to the OP. To verify the
	// signature, the OP uses a private association that was generated
	// when it issued the positive assertion.

	// 11.4.2.1.  Request Parameters
	params := make(url.Values)
	// openid.mode: Value: "check_authentication"
	params.Add("openid.mode", "check_authentication")
	// Exact copies of all fields from the authentication response,
	// except for "openid.mode".
	for k, vs := range vals {
		if k == "openid.mode" {
			continue
		}
		for _, v := range vs {
			params.Add(k, v)
		}
	}
	resp, err := getter.Post(vals.Get("openid.op_endpoint"), params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	response := string(content)
	lines := strings.Split(response, "\n")

	isValid := false
	nsValid := false
	for _, l := range lines {
		if l == "is_valid:true" {
			isValid = true
		} else if l == "ns:http://specs.openid.net/auth/2.0" {
			nsValid = true
		}
	}
	if isValid && nsValid {
		// Yay !
		return nil
	}

	return errors.New("Could not verify assertion with provider")
}
