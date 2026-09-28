// Package policy evaluates a RegionSpec against compliance rules. It runs
// BEFORE any provisioning — a non-compliant spec never reaches the reconciler.
package policy

import (
	"fmt"

	"github.com/khwajalabs/regionctl/internal/spec"
)

// Violation is a single failed rule. A non-empty slice means the spec is denied.
type Violation struct {
	Rule    string
	Message string
}

type Rule func(spec.RegionSpec) []Violation

// Evaluate runs every rule and collects all violations. It does NOT stop at the
// first — an operator should see every problem at once, not fix-and-retry.
func Evaluate(s spec.RegionSpec, rules []Rule) []Violation {
	var out []Violation
	for _, r := range rules {
		out = append(out, r(s)...)
	}
	return out
}

// DefaultRules for a sensitive region.
func DefaultRules() []Rule {
	return []Rule{
		noPublicEgressWhenIsolated,
		allServicesEncrypted,
		requireOwnerTag,
	}
}

func noPublicEgressWhenIsolated(s spec.RegionSpec) []Violation {
	if !s.Isolated {
		return nil
	}
	var v []Violation
	for _, n := range s.Networks {
		if n.PublicEgress {
			v = append(v, Violation{
				Rule:    "no-public-egress-when-isolated",
				Message: fmt.Sprintf("network %s has public egress in an isolated region", n.CIDR),
			})
		}
	}
	return v
}

func allServicesEncrypted(s spec.RegionSpec) []Violation {
	var v []Violation
	for _, svc := range s.Services {
		if !svc.Encrypted {
			v = append(v, Violation{
				Rule:    "all-services-encrypted",
				Message: fmt.Sprintf("service %s is not encrypted", svc.Name),
			})
		}
	}
	return v
}

func requireOwnerTag(s spec.RegionSpec) []Violation {
	if _, ok := s.Tags["owner"]; !ok {
		return []Violation{{Rule: "require-owner-tag", Message: "spec is missing an 'owner' tag"}}
	}
	return nil
}
