package openid

import (
	"bufio"
	"bytes"
	"errors"
	"net/http"
	"net/url"
)

type fakeGetter struct {
	urls      map[string]string
	redirects map[string]string
}

var testGetter = &fakeGetter{
	make(map[string]string), make(map[string]string)}

var testInstance = &OpenID{urlGetter: testGetter}

func (f *fakeGetter) Get(uri string, headers map[string]string) (resp *http.Response, err error) {
	key := uri
	for k, v := range headers {
		key += "#" + k + "#" + v
	}

	if doc, ok := f.urls[key]; ok {
		request, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			return nil, err
		}

		return http.ReadResponse(bufio.NewReader(
			bytes.NewBuffer([]byte(doc))), request)
	}
	if uri, ok := f.redirects[key]; ok {
		return f.Get(uri, headers)
	}

	return nil, errors.New("404 not found")
}

func (f *fakeGetter) Post(uri string, form url.Values) (resp *http.Response, err error) {
	return f.Get("POST@"+uri, nil)
}

func init() {
	// Prepare (http#header#header-val --> http response) pairs.

	// === For Yadis discovery ==================================
	// Directly reffers a valid XRDS document
	testGetter.urls["http://example.com/xrds#Accept#application/xrds+xml"] = `HTTP/1.0 200 OK
Content-Type: application/xrds+xml; charset=UTF-8

<?xml version="1.0" encoding="UTF-8"?>
<xrds:XRDS xmlns:xrds="xri://$xrds" xmlns="xri://$xrd*($v*2.0)"
xmlns:openid="http://openid.net/xmlns/1.0">
  <XRD>
    <Service xmlns="xri://$xrd*($v*2.0)">
      <Type>http://specs.openid.net/auth/2.0/signon</Type>
      <URI>foo</URI>
      <LocalID>bar</LocalID>
    </Service>
  </XRD>
</xrds:XRDS>`

	// Uses a X-XRDS-Location header to redirect to the valid XRDS document.
	testGetter.urls["http://example.com/xrds-loc#Accept#application/xrds+xml"] = `HTTP/1.0 200 OK
X-XRDS-Location: http://example.com/xrds

nothing interesting here`

	// Html document, with meta tag X-XRDS-Location. Points to the
	// previous valid XRDS document.
	testGetter.urls["http://example.com/xrds-meta#Accept#application/xrds+xml"] = `HTTP/1.0 200 OK
Content-Type: text/html

<html>
<head>
<meta http-equiv="X-XRDS-Location" content="http://example.com/xrds">`

	// === For HTML discovery ===================================
	testGetter.urls["http://example.com/html"] = `HTTP/1.0 200 OK

<html>
<head>
<link rel="openid2.provider" href="example.com/openid">
<link rel="openid2.local_id" href="bar-name">`

	testGetter.redirects["http://example.com/html-redirect"] = "http://example.com/html"
}
