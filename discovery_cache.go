package openid

type DiscoveredInfo interface {
	OpEndpoint() string
	OpLocalID() string
	ClaimedID() string
	// ProtocolVersion: it's always openId 2.
}

type DiscoveryCache interface {
	Put(id string, info DiscoveredInfo)
	// Return a discovered info, or nil.
	Get(id string) DiscoveredInfo
}

type SimpleDiscoveredInfo struct {
	opEndpoint string
	opLocalID  string
	claimedID  string
}

func (s *SimpleDiscoveredInfo) OpEndpoint() string {
	return s.opEndpoint
}

func (s *SimpleDiscoveredInfo) OpLocalID() string {
	return s.opLocalID
}

func (s *SimpleDiscoveredInfo) ClaimedID() string {
	return s.claimedID
}

type SimpleDiscoveryCache map[string]DiscoveredInfo

func (s SimpleDiscoveryCache) Put(id string, info DiscoveredInfo) {
	s[id] = info
}

func (s SimpleDiscoveryCache) Get(id string) DiscoveredInfo {
	if info, has := s[id]; has {
		return info
	}
	return nil
}

func compareDiscoveredInfo(a DiscoveredInfo, opEndpoint, opLocalID, claimedID string) bool {
	return a != nil &&
		a.OpEndpoint() == opEndpoint &&
		a.OpLocalID() == opLocalID &&
		a.ClaimedID() == claimedID
}
