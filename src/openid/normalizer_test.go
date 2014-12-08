package openid

import (
	"testing"
)

func TestNormalizeXRI(t *testing.T) {
	if _, err := Normalize("xri://asdf"); err == nil {
		t.Errorf("XRI not supported")
	}
	if _, err := Normalize("=asdf"); err == nil {
		t.Errorf("XRI not supported")
	}
	if _, err := Normalize("@asdf"); err == nil {
		t.Errorf("XRI not supported")
	}
}

func TestNormalizeHttp(t *testing.T) {
	if n, err := Normalize("foo.com"); err != nil || n != "http://foo.com" {
		t.Errorf("http:// should be added")
	}
	if n, err := Normalize("http://foo.com"); err != nil || n != "http://foo.com" {
		t.Errorf("valid URL should not be modified")
	}
	if n, err := Normalize("https://foo.com"); err != nil || n != "https://foo.com" {
		t.Errorf("https:// URLs are valid")
	}
}

func TestNormalizeFragment(t *testing.T) {
	if n, err := Normalize("http://foo.com#bar"); err != nil || n != "http://foo.com" {
		t.Errorf("URL fragments must be removed")
	}
}
