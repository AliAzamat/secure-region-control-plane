package reconcile

import (
	"errors"
	"testing"

	"github.com/khwajalabs/regionctl/internal/spec"
)

// crashingProvider fails its Nth Create to simulate a control-plane crash
// partway through applying a plan.
type crashingProvider struct {
	*fakeProvider
	failOn int
	seen   int
}

func (c *crashingProvider) Create(r Resource) error {
	c.seen++
	if c.seen == c.failOn {
		return errors.New("control plane crashed mid-apply")
	}
	return c.fakeProvider.Create(r)
}

func TestChaosConvergesAfterMidApplyCrash(t *testing.T) {
	s := spec.RegionSpec{
		Name: "au-sensitive-1",
		Networks: []spec.Network{
			{CIDR: "10.0.0.0/16"}, {CIDR: "10.1.0.0/16"},
		},
		Services: []spec.Service{
			{Name: "kms", Encrypted: true}, {Name: "logs", Encrypted: true},
		},
	}
	// Same underlying store survives the "crash" — like a restarted process
	// reconnecting to the same region.
	store := newFake()
	crashy := &crashingProvider{fakeProvider: store, failOn: 3}

	// First bootstrap attempt: crashes on the 3rd create.
	plan, _ := Plan(s, crashy)
	if len(plan) != 4 {
		t.Fatalf("want 4 planned, got %d", len(plan))
	}
	if err := Apply(plan, crashy); err == nil {
		t.Fatal("expected a crash on the 3rd create")
	}
	// Two resources made it in before the crash.
	if store.creates != 2 {
		t.Fatalf("want 2 created before crash, got %d", store.creates)
	}

	// RESTART: a healthy control plane re-plans against the same store and
	// converges the rest. No crash injection this time.
	replan, _ := Plan(s, store)
	if len(replan) != 2 {
		t.Fatalf("restart should plan the 2 remaining, got %d", len(replan))
	}
	if err := Apply(replan, store); err != nil {
		t.Fatal(err)
	}

	// The region is fully, compliantly converged and no resource was duplicated.
	final, _ := Plan(s, store)
	if len(final) != 0 {
		t.Fatalf("region must be fully converged, still %d to do", len(final))
	}
	if store.creates != 4 {
		t.Fatalf("exactly 4 resources total; got %d (duplicate?)", store.creates)
	}
}
