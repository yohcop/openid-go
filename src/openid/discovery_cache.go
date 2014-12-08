package openid

type DiscoveredInfo interface {
	OpEndpoint() string
	OpLocalId() string
	ClaimedId() string
	// ProtocolVersion: it's always openId 2.
}

type DiscoveryCache interface {
	Put(id string, info DiscoveredInfo)
	// Return a discovered info, or nil.
	Get(id string) DiscoveredInfo
}

type SimpleDiscoveredInfo struct {
	opEndpoint string
	opLocalId  string
	claimedId  string
}

func (s *SimpleDiscoveredInfo) OpEndpoint() string {
	return s.opEndpoint
}

func (s *SimpleDiscoveredInfo) OpLocalId() string {
	return s.opLocalId
}

func (s *SimpleDiscoveredInfo) ClaimedId() string {
	return s.claimedId
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

func compareDiscoveredInfo(a DiscoveredInfo, opEndpoint, opLocalId, claimedId string) bool {
	return a != nil &&
		a.OpEndpoint() == opEndpoint &&
		a.OpLocalId() == opLocalId &&
		a.ClaimedId() == claimedId
}
