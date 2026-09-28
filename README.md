# Amazon | Region Bootstrap — A Secure, Auditable Control Plane for Sensitive Workloads

An advanced infrastructure capstone modeled on AWS Region Services' work standing up highly secure, isolated environments for sensitive government workloads. You build the Go control plane that turns a declarative region spec into a running region — safely. Store every region config as an immutable, content-addressed version, evaluate a policy engine that BLOCKS a non-compliant spec before a single resource is touched, drive a reconcile loop that converges desired state idempotently (and is safe to re-run after a crash), and record every decision in a hash-chained audit log an auditor can verify end to end. Finish with an operator CLI, a break-glass approval gate for privileged changes, and a chaos drill that kills the control plane mid-bootstrap and proves the region still converges to the intended, compliant state.

Built step-by-step with [KhwajaLabs Build](https://khwajalabs.com).

## Stack
- Go
- PostgreSQL
- Policy-as-Code
- SHA-256
- TypeScript
