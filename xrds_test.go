package openid

import (
	"testing"
)

func TestXrds(t *testing.T) {
	testExpectOpID(t, []byte(`
<?xml version="1.0" encoding="UTF-8"?>
<xrds:XRDS xmlns:xrds="xri://$xrds" xmlns="xri://$xrd*($v*2.0)"
xmlns:openid="http://openid.net/xmlns/1.0">
  <XRD>
    <Service priority="10">
      <Type>http://openid.net/signon/1.0</Type>
      <URI>http://www.myopenid.com/server</URI>
      <openid:Delegate>http://smoker.myopenid.com/</openid:Delegate>
    </Service>
    <Service priority="50">
      <Type>http://openid.net/signon/1.0</Type>
      <URI>http://www.livejournal.com/openid/server.bml</URI>
      <openid:Delegate>
        http://www.livejournal.com/users/frank/
      </openid:Delegate>
    </Service>
    <Service priority="20">
      <Type>http://lid.netmesh.org/sso/2.0</Type>
    </Service>
    <Service>
      <Type>http://specs.openid.net/auth/2.0/server</Type>
      <URI>foo</URI>
    </Service>
  </XRD>
</xrds:XRDS>
    `), "foo", "")

	testExpectOpID(t, []byte(`
<?xml version="1.0" encoding="UTF-8"?>
<xrds:XRDS xmlns:xrds="xri://$xrds" xmlns="xri://$xrd*($v*2.0)"
xmlns:openid="http://openid.net/xmlns/1.0">
  <XRD>
    <Service xmlns="xri://$xrd*($v*2.0)">
      <Type>http://specs.openid.net/auth/2.0/signon</Type>
      <URI>https://www.exampleprovider.com/endpoint/</URI>
      <LocalID>https://exampleuser.exampleprovider.com/</LocalID>
    </Service>
  </XRD>
</xrds:XRDS>
    `),
		"https://www.exampleprovider.com/endpoint/",
		"https://exampleuser.exampleprovider.com/")

	// OP Identifier Element has priority over Claimed Identifier Element
	testExpectOpID(t, []byte(`
<?xml version="1.0" encoding="UTF-8"?>
<xrds:XRDS xmlns:xrds="xri://$xrds" xmlns="xri://$xrd*($v*2.0)"
xmlns:openid="http://openid.net/xmlns/1.0">
  <XRD>
    <Service xmlns="xri://$xrd*($v*2.0)">
      <Type>http://specs.openid.net/auth/2.0/signon</Type>
      <URI>https://www.exampleprovider.com/endpoint-signon/</URI>
    </Service>
    <Service xmlns="xri://$xrd*($v*2.0)">
      <Type>http://specs.openid.net/auth/2.0/server</Type>
      <URI>https://www.exampleprovider.com/endpoint-server/</URI>
    </Service>
  </XRD>
</xrds:XRDS>
    `),
		"https://www.exampleprovider.com/endpoint-server/",
		"")
}

func testExpectOpID(t *testing.T, xrds []byte, op, id string) {
	receivedOp, receivedID, err := parseXrds(xrds)
	if err != nil {
		t.Errorf("Got an error parsing XRDS (%s): %s", string(xrds), err)
	} else {
		if receivedOp != op {
			t.Errorf("Extracted OP does not match: Exepect %s, Got %s",
				op, receivedOp)
		}
		if receivedID != id {
			t.Errorf("Extracted ID does not match: Exepect %s, Got %s",
				id, receivedID)
		}
	}
}
