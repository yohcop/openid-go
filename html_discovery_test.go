package openid

import (
	"bytes"
	"testing"
)

func TestFindEndpointFromLink(t *testing.T) {
	searchLink(t, `
      <html>
        <head>
          <link rel="openid2.provider" href="example.com/openid">
      `, "example.com/openid", "", false)
	searchLink(t, `
      <html>
        <head>
          <link rel="openid2.provider" href="foo.com">
          <link rel="openid2.local_id" href="bar-name">
        </head>
      </html>
      `, "foo.com", "bar-name", false)
	// Self-closing link
	searchLink(t, `
      <html>
        <head>
          <link rel="openid2.provider" href="selfclose.com" />
          <link rel="openid2.local_id" href="selfclose-name" />
        </head>
      </html>
      `, "selfclose.com", "selfclose-name", false)
}

func TestNoEndpointFromLink(t *testing.T) {
	searchLink(t, `
      <html>
        <head>
          <link rel="openid2.provider">
      `, "", "", true)
	// Outside of head.
	searchLink(t, `
      <html>
        <head></head>
        <link rel="openid2.provider" href="example.com/openid">
      `, "", "", true)
}

func searchLink(t *testing.T, doc, opEndpoint, claimedID string, err bool) {
	r := bytes.NewReader([]byte(doc))
	op, id, e := findProviderFromHeadLink(r)
	if (e != nil) != err {
		t.Errorf("Unexpected error: '%s'", e)
	} else if e == nil {
		if op != opEndpoint {
			t.Errorf("Found bad endpoint: Expected %s, Got %s",
				op, opEndpoint)
		}
		if id != claimedID {
			t.Errorf("Found bad id: Expected %s, Got %s",
				id, claimedID)
		}
	}
}
