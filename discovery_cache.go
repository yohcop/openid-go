package openid

import (
	"sync"
	"time"
)

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

type SimpleDiscoveryCache struct {
	cache map[string]DiscoveredInfo
	mutex *sync.Mutex
}

func NewSimpleDiscoveryCache() *SimpleDiscoveryCache {
	return &SimpleDiscoveryCache{cache: map[string]DiscoveredInfo{}, mutex: &sync.Mutex{}}
}

func (s *SimpleDiscoveryCache) Put(id string, info DiscoveredInfo) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.cache[id] = info
}

func (s *SimpleDiscoveryCache) Get(id string) DiscoveredInfo {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if info, has := s.cache[id]; has {
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

type timedDiscoveredInfo struct {
	info DiscoveredInfo
	time time.Time
}

type TimedDiscoveryCache struct {
	cache map[string]timedDiscoveredInfo
	ttl time.Duration
	mutex *sync.Mutex
}

func NewTimedDiscoveryCache(ttl time.Duration) *TimedDiscoveryCache {
	return &TimedDiscoveryCache{cache: map[string]timedDiscoveredInfo{}, ttl: ttl, mutex: &sync.Mutex{}}
}

func (s *TimedDiscoveryCache) Put(id string, info DiscoveredInfo) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.cache[id] = timedDiscoveredInfo{info:info, time:time.Now()}
}

func (s *TimedDiscoveryCache) Get(id string) DiscoveredInfo {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Delete old cached while we are at it.
	newCache := map[string]timedDiscoveredInfo{}
	now := time.Now()
	for k, e := range s.cache {
		diff := now.Sub(e.time)
		if diff <= s.ttl {
			newCache[k] = e
		}
	}
	s.cache = newCache

	if info, has := s.cache[id]; has {
		return info.info
	}
	return nil
}
