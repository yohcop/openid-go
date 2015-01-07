package openid

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"golang.org/x/net/html"
)

var yadisHeaders = map[string]string{
	"Accept": "application/xrds+xml"}

func yadisDiscovery(id string, getter httpGetter) (opEndpoint string, opLocalId string, err error) {
	// Section 6.2.4 of Yadis 1.0 specifications.
	// The Yadis Protocol is initiated by the Relying Party Agent
	// with an initial HTTP request using the Yadis URL.

	// This request MUST be either a GET or a HEAD request.

	// A GET or HEAD request MAY include an HTTP Accept
	// request-header (HTTP 14.1) specifying MIME media type,
	// application/xrds+xml.
	resp, err := getter.Get(id, yadisHeaders)
	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close()

	// Section 6.2.5 from Yadis 1.0 spec: Response

	contentType := resp.Header.Get("Content-Type")

	// The response MUST be one of:
	// (see 6.2.6 for precedence)
	if l := resp.Header.Get("X-XRDS-Location"); l != "" {
		// 2. HTTP response-headers that include an X-XRDS-Location
		// response-header, together with a document
		return getYadisResourceDescriptor(l, getter)
	} else if strings.Contains(contentType, "text/html") {
		// 1. An HTML document with a <head> element that includes a
		// <meta> element with http-equiv attribute, X-XRDS-Location,

		if metaContent, err := findMetaXrdsLocation(resp.Body); err == nil {
			return getYadisResourceDescriptor(metaContent, getter)
		} else {
			return "", "", err
		}
	} else if strings.Contains(contentType, "application/xrds+xml") {
		// 4. A document of MIME media type, application/xrds+xml.
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			return parseXrds(body)
		} else {
			return "", "", err
		}
	}
	// 3. HTTP response-headers only, which MAY include an
	// X-XRDS-Location response-header, a content-type
	// response-header specifying MIME media type,
	// application/xrds+xml, or both.
	//   (this is handled by one of the 2 previous if statements)
	return "", "", errors.New("No expected header, or content type")
}

// Similar as above, but we expect an absolute Yadis document URL.
func getYadisResourceDescriptor(id string, getter httpGetter) (opEndpoint string, opLocalId string, err error) {
	resp, err := getter.Get(id, yadisHeaders)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	// 4. A document of MIME media type, application/xrds+xml.
	if body, err := ioutil.ReadAll(resp.Body); err == nil {
		return parseXrds(body)
	} else {
		return "", "", err
	}
	return
}

// Search for
// <head>
//    <meta http-equiv="X-XRDS-Location" content="....">
func findMetaXrdsLocation(input io.Reader) (location string, err error) {
	tokenizer := html.NewTokenizer(input)
	inHead := false
	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			return "", tokenizer.Err()
		case html.StartTagToken, html.EndTagToken:
			tk := tokenizer.Token()
			if tk.Data == "head" {
				if tt == html.StartTagToken {
					inHead = true
				} else {
					return "", errors.New("Meta X-XRDS-Location not found")
				}
			} else if inHead && tk.Data == "meta" {
				ok := false
				content := ""
				for _, attr := range tk.Attr {
					if attr.Key == "http-equiv" &&
						attr.Val == "X-XRDS-Location" {
						ok = true
					} else if attr.Key == "content" {
						content = attr.Val
					}
				}
				if ok && len(content) > 0 {
					return content, nil
				}
			}
		}
	}
	return "", errors.New("Meta X-XRDS-Location not found")
}
