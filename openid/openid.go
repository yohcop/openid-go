package openid

import (
	"encoding/xml"
	"exp/html"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Normalize(id string) string {
	var normalized string
	if strings.HasPrefix(id, "xri://") {
		normalized = id[6:]
	} else if strings.HasPrefix(id, "xri://$ip") {
		normalized = id[9:]
	} else if strings.HasPrefix(id, "xri://$dns*") {
		normalized = id[10:]
	} else {
		normalized = id
	}
	if normalized[0] == '=' || normalized[0] == '@' || normalized[0] == '$' || normalized[0] == '!' {
		return normalized
	}
	if strings.HasPrefix(id, "http://") || strings.HasPrefix(id, "https://") {
		return id
	}
	return "http://" + id
}

type DiscoveryError struct {
	str string
}

func (e *DiscoveryError) Error() string {
	return e.str
}

func DiscoverXml(id string) (*string, error) {
	resp, err := http.Get(id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	parser := xml.NewParser(resp.Body)
	inURI := false
	for {
		t, err := parser.Token()
		if err != nil {
			return nil, err
		}
		switch tt := t.(type) {
		case xml.StartElement:
			if tt.Name.Local == "URI" {
				inURI = true
			}
		case xml.CharData:
			if inURI {
				s := string([]byte(tt))
				return &s, nil
			}
		}
	}
	return nil, &DiscoveryError{str: "URI not found"}
}

func DiscoverHtml(id string) (*string, error) {
	resp, err := http.Get(id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			log.Println("Error: ", tokenizer.Err())
			return nil, tokenizer.Err()
		case html.StartTagToken, html.EndTagToken:
			tk := tokenizer.Token()
			if tk.Data == "link" {
				ok := false
				for _, attr := range tk.Attr {
					if attr.Key == "rel" && attr.Val == "openid2.provider" {
						log.Println(tk.String())
						ok = true
					} else if attr.Key == "href" && ok {
						return &attr.Val, nil
					}
				}
			}
		}
	}
	return nil, &DiscoveryError{str: "provider not found"}
}

func PrepareRedirect(url *string, returnTo string) (*string, error) {
	redirect := *url + "?openid.ns=http://specs.openid.net/auth/2.0"
	redirect += "&openid.claimed_id=http://specs.openid.net/auth/2.0/identifier_select"
	redirect += "&openid.identity=http://specs.openid.net/auth/2.0/identifier_select"
	redirect += "&openid.return_to=" + returnTo
	redirect += "&openid.realm=" + returnTo
	redirect += "&openid.mode=checkid_setup"
	// To ask for email, add:
	//redirect += "&openid.ns.ax=http://openid.net/srv/ax/1.0"
	//redirect += "&openid.ax.mode=fetch_request"
	//redirect += "&openid.ax.type.email=http://axschema.org/contact/email"
	//redirect += "&openid.ax.required=email"
	return &redirect, nil
}

func ValidateLogin(params map[string]string, expectedURL string, thisURL string) (*string, error) {
	if v, ok := params["openid.mode"]; !ok || v != "id_res" {
		return nil, &DiscoveryError{str: "Open ID connection failed"}
	}

	if v, ok := params["openid.return_to"]; !ok || v != expectedURL || v != thisURL {
		return nil, &DiscoveryError{str: "return_to URL doesn't match: " + v + " vs " + expectedURL + " vs " + thisURL}
	}

	// TODO: verify that openid.op_endpoint matches the discoverd info, etc.
	// See http://openid.net/specs/openid-authentication-2_0.html#verification
	// section 11.2.
	if !verifyAssertion(params) {
		return nil, &DiscoveryError{str: "Verification failed."}
	}

	id, ok := params["openid.identity"]
	if !ok {
		return nil, &DiscoveryError{str: "Could not find openId identity in openId response"}
	}
	return &id, nil
}

func verifyAssertion(params map[string]string) bool {
	fields := map[string][]string{"openid.mode": []string{"check_authentication"}}
	for k, v := range params {
		if k != "openid.mode" {
			fields[k] = []string{v}
		}
	}

	resp, err := http.PostForm(params["openid.op_endpoint"], fields)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	response := string(content)
	log.Println(response)
	lines := strings.Split(response, "\n")
	for _, l := range lines {
		if l == "is_valid:true" {
			return true
		}
	}
	return false
}
