package openid

import (
	"bytes"
	"testing"
)

func TestFindMetaXrdsLocation(t *testing.T) {
	searchMeta(t, `
      <html>
        <head>
          <meta http-equiv="X-XRDS-Location" content="foo.com">
      `, "foo.com", false)
	searchMeta(t, `
      <html>
        <head>
          <meta http-equiv="other" content="blah.com">
          <meta http-equiv="X-XRDS-Location" content="foo.com">
      `, "foo.com", false)
}

func TestMetaXrdsLocationOutsideHead(t *testing.T) {
	searchMeta(t, `
      <html>
        <meta http-equiv="X-XRDS-Location" content="foo.com">
      `, "", true)
	searchMeta(t, `
      <html>
        <head></head>
        <meta http-equiv="X-XRDS-Location" content="foo.com">
      `, "", true)
}

func TestNoMetaXrdsLocation(t *testing.T) {
	searchMeta(t, `
      <html><head>
        <meta http-equiv="bad-tag" content="foo.com">
      `, "", true)
}

func searchMeta(t *testing.T, doc, loc string, err bool) {
	r := bytes.NewReader([]byte(doc))
	res, e := findMetaXrdsLocation(r)
	if (e != nil) != err {
		t.Errorf("Unexpected error: '%s'", e)
	} else if e == nil {
		if res != loc {
			t.Errorf("Found bad location: Expected %s, Got %s", loc, res)
		}
	}
}
