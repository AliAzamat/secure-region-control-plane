package audit

import "testing"

func TestVerifyDetectsTampering(t *testing.T) {
	l := &Log{}
	l.Append(Entry{Actor: "op", Action: "apply", SpecAddr: "sha256:aaa"})
	l.Append(Entry{Actor: "op", Action: "apply", SpecAddr: "sha256:bbb"})
	l.Append(Entry{Actor: "op", Action: "break-glass:apply", SpecAddr: "sha256:ccc"})

	if bad := l.Verify(); bad != -1 {
		t.Fatalf("fresh chain should verify, broke at %d", bad)
	}

	// Forge history: rewrite entry 1's detail to hide it, keep its stored hash.
	l.entries[1].Detail = "nothing to see here"

	if bad := l.Verify(); bad != 1 {
		t.Fatalf("tampering with entry 1 should be caught at 1, got %d", bad)
	}
}
