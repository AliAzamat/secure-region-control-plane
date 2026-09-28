// Package spec defines the declarative region spec and its content address.
package spec

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
)

// RegionSpec is the DESIRED state of a region. The control plane's only job is
// to make reality match this — nothing here describes HOW, only WHAT.
type RegionSpec struct {
	Name          string            `json:"name"`
	Classification string           `json:"classification"` // e.g. "sensitive"
	Isolated      bool              `json:"isolated"`        // no public egress
	Networks      []Network         `json:"networks"`
	Services      []Service         `json:"services"`
	Tags          map[string]string `json:"tags"`
}

type Network struct {
	CIDR         string `json:"cidr"`
	PublicEgress bool   `json:"public_egress"`
}

type Service struct {
	Name     string `json:"name"`
	Encrypted bool  `json:"encrypted"`
}

// ContentAddress is the SHA-256 of the spec's canonical JSON. Two specs with the
// same content have the same address; a version is named by WHAT it is, not when
// it was made. That is what makes the audit trail meaningful.
func (s RegionSpec) ContentAddress() string {
	sum := sha256.Sum256(s.canonicalJSON())
	return "sha256:" + hex.EncodeToString(sum[:])
}

func (s RegionSpec) canonicalJSON() []byte {
	// Sort tag keys so the JSON is deterministic — otherwise the same spec
	// hashes differently run to run and the content address is worthless.
	keys := make([]string, 0, len(s.Tags))
	for k := range s.Tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	ordered := make(map[string]string, len(keys))
	for _, k := range keys {
		ordered[k] = s.Tags[k]
	}
	clone := s
	clone.Tags = ordered
	b, _ := json.Marshal(clone)
	return b
}
