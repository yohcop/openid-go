package openid

import (
	"errors"
	"strings"
)

func Normalize(id string) (string, error) {
	// 7.2 from openID 2.0 spec.

	//If the user's input starts with the "xri://" prefix, it MUST be
	//stripped off, so that XRIs are used in the canonical form.
	if strings.HasPrefix(id, "xri://") {
		id = id[6:]
		return "", errors.New("XRI identifiers not supported")
	}

	// If the first character of the resulting string is an XRI
	// Global Context Symbol ("=", "@", "+", "$", "!") or "(", as
	// defined in Section 2.2.1 of [XRI_Syntax_2.0], then the input
	// SHOULD be treated as an XRI.
	if s := string(id[0]); s == "=" || s == "@" || s == "+" || s == "$" || s == "!" {
		return "", errors.New("XRI identifiers not supported")
	}

	// Otherwise, the input SHOULD be treated as an http URL; if it
	// does not include a "http" or "https" scheme, the Identifier
	// MUST be prefixed with the string "http://". If the URL
	// contains a fragment part, it MUST be stripped off together
	// with the fragment delimiter character "#". See Section 11.5.2 for
	// more information.
	if !strings.HasPrefix(id, "http://") && !strings.HasPrefix(id,
		"https://") {
		id = "http://" + id
	}
	if fragmentIndex := strings.Index(id, "#"); fragmentIndex != -1 {
		id = id[0:fragmentIndex]
	}

	// URL Identifiers MUST then be further normalized by both
	// following redirects when retrieving their content and finally
	// applying the rules in Section 6 of [RFC3986] to the final
	// destination URL. This final URL MUST be noted by the Relying
	// Party as the Claimed Identifier and be used when requesting
	// authentication.
	return id, nil
}
