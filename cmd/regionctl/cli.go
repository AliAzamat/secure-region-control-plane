package main

import (
	"fmt"

	"github.com/khwajalabs/regionctl/internal/audit"
	"github.com/khwajalabs/regionctl/internal/policy"
	"github.com/khwajalabs/regionctl/internal/reconcile"
	"github.com/khwajalabs/regionctl/internal/spec"
)

// Bootstrap is the full safe pipeline: policy gate -> plan -> (only if approved)
// apply, auditing each stage. Returns an exit code: 0 ok, 2 policy-denied,
// 1 error.
func Bootstrap(s spec.RegionSpec, p reconcile.Provider, log *audit.Log, apply bool) int {
	// 1. Policy gate — deny before anything else.
	if v := policy.Evaluate(s, policy.DefaultRules()); len(v) > 0 {
		log.Append(audit.Entry{Action: "policy-deny", SpecAddr: s.ContentAddress(),
			Detail: fmt.Sprintf("%d violation(s)", len(v))})
		for _, viol := range v {
			fmt.Printf("DENY [%s] %s\n", viol.Rule, viol.Message)
		}
		return 2
	}

	// 2. Plan — pure read, always shown.
	actions, err := reconcile.Plan(s, p)
	if err != nil {
		return 1
	}
	fmt.Printf("plan: %d change(s) for region %s\n", len(actions), s.Name)
	for _, a := range actions {
		fmt.Printf("  + %s %s (%s)\n", a.Resource.Kind, a.Resource.ID, a.Reason)
	}

	// 3. Apply — ONLY when explicitly asked. Default is plan-only.
	if !apply {
		fmt.Println("(plan-only; re-run with --apply to converge)")
		return 0
	}
	if err := reconcile.Apply(actions, p); err != nil {
		return 1
	}
	log.Append(audit.Entry{Actor: "operator", Action: "apply", SpecAddr: s.ContentAddress(),
		Detail: fmt.Sprintf("applied %d action(s)", len(actions))})
	fmt.Printf("applied %d change(s)\n", len(actions))
	return 0
}
