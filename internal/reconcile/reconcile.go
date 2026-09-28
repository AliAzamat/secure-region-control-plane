// Package reconcile converges a region toward a RegionSpec. Every action is
// idempotent, so a re-run after a crash is safe and does no harm.
package reconcile

import (
	"fmt"

	"github.com/khwajalabs/regionctl/internal/spec"
)

// Resource is a piece of a region the control plane can create.
type Resource struct {
	Kind string
	ID   string
}

// Provider is the effectful boundary — the real one talks to a cloud API; a
// fake one is used in tests. Both must be idempotent: creating an existing
// resource is a no-op success, not an error.
type Provider interface {
	Exists(r Resource) (bool, error)
	Create(r Resource) error
}

// Action is one planned change. A plan is a list of these — the operator sees
// the plan before any of it runs.
type Action struct {
	Resource Resource
	Reason   string
}

// Plan computes the gap between desired (spec) and actual (provider). It makes
// NO changes — pure read. This is the "terraform plan" half of the split.
func Plan(s spec.RegionSpec, p Provider) ([]Action, error) {
	var actions []Action
	for _, n := range s.Networks {
		r := Resource{Kind: "network", ID: n.CIDR}
		exists, err := p.Exists(r)
		if err != nil {
			return nil, fmt.Errorf("checking %v: %w", r, err)
		}
		if !exists {
			actions = append(actions, Action{Resource: r, Reason: "network absent"})
		}
	}
	for _, svc := range s.Services {
		r := Resource{Kind: "service", ID: svc.Name}
		exists, err := p.Exists(r)
		if err != nil {
			return nil, fmt.Errorf("checking %v: %w", r, err)
		}
		if !exists {
			actions = append(actions, Action{Resource: r, Reason: "service absent"})
		}
	}
	return actions, nil
}

// Apply executes a plan. Because Create is idempotent and we re-check Exists,
// running Apply twice converges to the same state — the second run is a no-op.
func Apply(actions []Action, p Provider) error {
	for _, a := range actions {
		exists, err := p.Exists(a.Resource)
		if err != nil {
			return err
		}
		if exists {
			continue // already there — idempotent skip
		}
		if err := p.Create(a.Resource); err != nil {
			return fmt.Errorf("creating %v: %w", a.Resource, err)
		}
	}
	return nil
}
