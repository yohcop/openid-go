package integration

// These tests fetch real data from google.com and other OpenID
// providers. If they change the files returned, or endpoints, or
// whatever, they will fail. It's ok though, they are full tests.

import (
  "testing"
  . "openid"
)

func TestGoogleCom(t *testing.T) {
  endpoint, localId, claimedId, err := Discover("https://www.google.com/accounts/o8/id")
  if err != nil {
    t.Errorf("Google.com discovery failed: %s", err)
  }
  if endpoint != "https://www.google.com/accounts/o8/ud" {
    t.Errorf("Unexpected Google.com endpoint: %s", endpoint)
  }
  if localId != "http://specs.openid.net/auth/2.0/identifier_select" {
    t.Errorf("Unexpected Google.com localId: %s", localId)
  }
  if claimedId != "http://specs.openid.net/auth/2.0/identifier_select" {
    t.Errorf("Unexpected Google.com localId: %s", claimedId)
  }
}

func TestYohcop(t *testing.T) {
  endpoint, localId, claimedId, err := Discover("http://yohcop.net")
  if err != nil {
    t.Errorf("Yohcop.net discovery failed")
  }
  if endpoint != "https://www.google.com/accounts/o8/ud?source=profiles" {
    t.Errorf("Unexpected yohcop.net endpoint: %s", endpoint)
  }
  if localId != "http://www.google.com/profiles/yohcop" {
    t.Errorf("Unexpected yohcop.net localId: %s", localId)
  }
  if claimedId != "http://yohcop.net" {
    t.Errorf("Unexpected yohcop.net claimedId: %s", claimedId)
  }
}
