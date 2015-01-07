package openid

import (
	"errors"
	"flag"
	"fmt"
	"sync"
	"time"
)

var max_nonce_age = flag.Duration("openid-max-nonce-age",
	60*time.Second,
	"Maximum accepted age for openid nonces. The bigger, the more"+
		"memory is needed to store used nonces.")

var nonceStore = &SimpleNonceStore{Store: make(map[string][]*Nonce)}

type NonceStore interface {
	// Returns nil if accepted, an error otherwise.
	Accept(endpoint, nonce string) error
}

type Nonce struct {
	T time.Time
	S string
}

type SimpleNonceStore struct {
	Store map[string][]*Nonce
	mutex sync.Mutex
}

func (d *SimpleNonceStore) Accept(endpoint, nonce string) error {
	// Value: A string 255 characters or less in length, that MUST be
	// unique to this particular successful authentication response.
	if len(nonce) < 20 || len(nonce) > 256 {
		return errors.New("Invalid nonce")
	}

	// The nonce MUST start with the current time on the server, and MAY
	// contain additional ASCII characters in the range 33-126 inclusive
	// (printable non-whitespace characters), as necessary to make each
	// response unique. The date and time MUST be formatted as specified in
	// section 5.6 of [RFC3339], with the following restrictions:

	// All times must be in the UTC timezone, indicated with a "Z".  No
	// fractional seconds are allowed For example:
	// 2005-05-15T17:11:51ZUNIQUE
	ts, err := time.Parse(time.RFC3339, nonce[0:20])
	if err != nil {
		return err
	}
	now := time.Now()
	diff := now.Sub(ts)
	if diff > *max_nonce_age {
		return fmt.Errorf("Nonce too old: %ds", diff.Seconds())
	}

	s := nonce[20:]

	// Meh.. now we have to use a mutex, to protect that map from
	// concurrent access. Could put a go routine in charge of it
	// though.
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if nonces, hasOp := d.Store[endpoint]; hasOp {
		// Delete old nonces while we are at it.
		newNonces := []*Nonce{}
		for _, n := range nonces {
			if n.T == ts && n.S == s {
				// If return early, just ignore the filtered list
				// we have been building so far...
				return errors.New("Nonce already used")
			}
			if now.Sub(n.T) < *max_nonce_age {
				newNonces = append(newNonces, n)
			}
		}
		d.Store[endpoint] = newNonces
	} else {
		d.Store[endpoint] = []*Nonce{{ts, s}}
	}
	return nil
}
