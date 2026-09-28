package reconcile

import (
	"testing"

	"github.com/khwajalabs/regionctl/internal/spec"
)

func TestApplyIsIdempotent(t *testing.T) {
	s := spec.RegionSpec{
		Name:     "r1",
		Networks: []spec.Network{{CIDR: "10.0.0.0/16"}},
		Services: []spec.Service{{Name: "kms", Encrypted: true}},
	}
	p := newFake()

	// First converge.
	plan1, _ := Plan(s, p)
	if len(plan1) != 2 {
		t.Fatalf("want 2 planned actions, got %d", len(plan1))
	}
	if err := Apply(plan1, p); err != nil {
		t.Fatal(err)
	}
	if p.creates != 2 {
		t.Fatalf("want 2 real creates, got %d", p.creates)
	}

	// Re-plan after converging: nothing left to do.
	plan2, _ := Plan(s, p)
	if len(plan2) != 0 {
		t.Fatalf("converged region should plan 0 actions, got %d", len(plan2))
	}

	// Re-apply the STALE first plan: idempotent, does no new work.
	if err := Apply(plan1, p); err != nil {
		t.Fatal(err)
	}
	if p.creates != 2 {
		t.Fatalf("re-apply must not create again; creates=%d", p.creates)
	}
}
