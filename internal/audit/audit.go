// Package audit is an append-only, hash-chained log of control-plane actions.
// Each entry commits to the previous entry's hash, so altering any past entry
// breaks every hash after it — tampering is detectable, not preventable.
package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

type Entry struct {
	Seq       int       `json:"seq"`
	Time      time.Time `json:"time"`
	Actor     string    `json:"actor"`
	Action    string    `json:"action"`     // e.g. "apply", "policy-deny", "break-glass"
	SpecAddr  string    `json:"spec_addr"`  // content address of the spec involved
	Detail    string    `json:"detail"`
	PrevHash  string    `json:"prev_hash"`  // hash of entry Seq-1
	Hash      string    `json:"hash"`       // hash of THIS entry (excluding Hash)
}

type Log struct {
	entries []Entry
}

func (l *Log) Append(e Entry) Entry {
	e.Seq = len(l.entries)
	e.Time = time.Now().UTC()
	if len(l.entries) > 0 {
		e.PrevHash = l.entries[len(l.entries)-1].Hash
	} else {
		e.PrevHash = "genesis"
	}
	e.Hash = hashEntry(e)
	l.entries = append(l.entries, e)
	return e
}

func hashEntry(e Entry) string {
	e.Hash = "" // hash excludes the Hash field itself
	b, _ := json.Marshal(e)
	sum := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(sum[:])
}

// Verify walks the chain and returns the first sequence number whose hash or
// prev-link does not check out. -1 means the whole chain is intact.
func (l *Log) Verify() int {
	prev := "genesis"
	for i, e := range l.entries {
		if e.PrevHash != prev {
			return i // chain link broken
		}
		if e.Hash != hashEntry(e) {
			return i // entry content was altered
		}
		prev = e.Hash
	}
	return -1
}

func (l *Log) Entries() []Entry { return l.entries }
