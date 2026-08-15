# Changelog

All notable changes are documented here. The project follows Semantic
Versioning and uses an unstable `v1alpha1` manifest until its first stable
contract.

## [0.1.1] - 2026-08-16

First release of `agent-runtime` as an open-source module under
`github.com/NDDev-OpenNetwork/agent-runtime`. The version line starts here: this
is a new module path, so the numbering of the predecessor it grew out of does
not carry over and would only mislead anyone resolving it.

### Added

- **Task**: a bounded, atomic unit of work with explicit acceptance, described
  by a versioned Task manifest. `runner.go` executes one, `manifest.go` types it,
  and `schemas/task-manifest-v1alpha1.schema.json` is the published contract for
  anyone reading the output without running the Go.
- **Goal**: a complex or long-running unit with a durable, revision-based
  journal, a living acceptance checklist, ordered phases and machine-verifiable
  receipts. `goal/store.go` persists it; `goal/types.go` types it.
- **Workspace resolution and process execution** with typed termination, so a
  caller can tell a cancelled run from a failed one and a timeout from either.
- **`observability`**: provider-neutral lifecycle events derived from those work
  units, with canonical JSON, correlation and causation identity, deterministic
  redaction, memory and JSONL sinks, deduplication and replay. It never replaces
  a Task result or a Goal journal.
- **`agent-runtime` CLI** over both contracts, useful in scripts, contract tests
  and restart-safe agent workflows.
- **`check-fuzz`**: starts every target named in `fuzz/v1alpha1.json` for a
  bounded number of iterations and fails when that list and the targets in the
  tree disagree, so a target cannot be added or lost silently.
- **`check-cold-compile`**: builds every platform this module claims to support,
  from a cold cache. CI runs on two of them, so this is what makes the claim true.

### Notes

The module is deliberately not a model provider, tool protocol, scheduler or
operating-system sandbox. It describes work and its lifecycle; running agents is
someone else's contract.
