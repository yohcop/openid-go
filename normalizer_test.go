package openid

import (
	"testing"
)

func TestNormalize(t *testing.T) {
	// OpenID 2.0 spec Appendix A.1. Normalization
	doNormalize(t, "example.com", "http://example.com/", true)
	doNormalize(t, "http://example.com", "http://example.com/", true)
	doNormalize(t, "https://example.com/", "https://example.com/", true)
	doNormalize(t, "http://example.com/user", "http://example.com/user", true)
	doNormalize(t, "http://example.com/user/", "http://example.com/user/", true)
	doNormalize(t, "http://example.com/", "http://example.com/", true)
	doNormalize(t, "=example", "=example", false)       // XRI not supported
	doNormalize(t, "(=example)", "(=example)", false)   // XRI not supported
	doNormalize(t, "xri://=example", "=example", false) // XRI not supported

	// Empty
	doNormalize(t, "", "", false)
	doNormalize(t, " ", "", false)
	doNormalize(t, "	", "", false)
	doNormalize(t, "xri://", "", false)
	doNormalize(t, "http://", "", false)
	doNormalize(t, "https://", "", false)

	// Padded with spacing
	doNormalize(t, " example.com  ", "http://example.com/", true)
	doNormalize(t, " 	http://example.com		 ", "http://example.com/", true)

	// XRI not supported
	doNormalize(t, "xri://asdf", "asdf", false)
	doNormalize(t, "=asdf", "=asdf", false)
	doNormalize(t, "@asdf", "@asdf", false)

	// HTTP
	doNormalize(t, "foo.com", "http://foo.com/", true)
	doNormalize(t, "http://foo.com", "http://foo.com/", true)
	doNormalize(t, "https://foo.com", "https://foo.com/", true)

	// Fragment need to be removed
	doNormalize(t, "http://foo.com#bar", "http://foo.com/", true)
	doNormalize(t, "http://foo.com/page#bar", "http://foo.com/page", true)
}

func doNormalize(t *testing.T, idIn, idOut string, succeed bool) {
	if id, err := Normalize(idIn); err != nil && succeed {
		t.Errorf("unexpected normalize error: gave %v, expected %v, got %v - %v", idIn, idOut, id, err)
	} else if err == nil && !succeed {
		t.Errorf("unexpected normalize success: gave %v, expected %v, got %v", idIn, idOut, id)
	} else if id != idOut {
		t.Errorf("unexpected normalize result: gave %v, expected %v, got %v", idIn, idOut, id)
	}
}
