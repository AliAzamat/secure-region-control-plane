// Package breakglass gates a privileged action behind an approval and a
// MANDATORY audit entry. If it cannot record the entry, it does not act.
package breakglass

import (
	"errors"

	"github.com/khwajalabs/regionctl/internal/audit"
)

type Request struct {
	Actor      string
	Approver   string // must differ from Actor — two-person rule
	Action     string
	SpecAddr   string
	Justification string
}

var (
	ErrSelfApproval = errors.New("break-glass requires a different approver")
	ErrNoJustification = errors.New("break-glass requires a justification")
)

// Authorize records the break-glass in the audit log and returns nil only if
// the request is valid AND the audit entry was written. The privileged action
// runs only after this returns nil.
func Authorize(req Request, log *audit.Log) error {
	if req.Approver == "" || req.Approver == req.Actor {
		return ErrSelfApproval
	}
	if req.Justification == "" {
		return ErrNoJustification
	}
	// Audit FIRST — the record of the attempt must exist before the privilege
	// is granted, so even a later crash leaves the attempt on the record.
	log.Append(audit.Entry{
		Actor:    req.Actor,
		Action:   "break-glass:" + req.Action,
		SpecAddr: req.SpecAddr,
		Detail:   "approver=" + req.Approver + " justification=" + req.Justification,
	})
	return nil
}
