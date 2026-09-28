package reconcile

// fakeProvider is an in-memory, idempotent Provider for tests.
type fakeProvider struct {
	created map[string]bool
	creates int // how many Create calls actually did work
}

func newFake() *fakeProvider { return &fakeProvider{created: map[string]bool{}} }

func (f *fakeProvider) Exists(r Resource) (bool, error) {
	return f.created[r.Kind+"/"+r.ID], nil
}

func (f *fakeProvider) Create(r Resource) error {
	key := r.Kind + "/" + r.ID
	if f.created[key] {
		return nil // idempotent: creating an existing resource is a no-op
	}
	f.created[key] = true
	f.creates++
	return nil
}
