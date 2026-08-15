# Agent instructions

The organization-wide contract in `NDDev-OpenNetwork/.github` applies. This file
adds only what is specific to this repository.

## What this module is

Two work-unit contracts — a bounded **Task** with a versioned manifest, and a
long-running **Goal** with a durable journal — plus an optional `observability`
package that derives lifecycle events from them. It is deliberately not a model
provider, tool protocol, scheduler, or sandbox. A change that adds one of those
belongs in a different repository.

## Checks

```
go test -race ./...
go vet ./...
go run honnef.co/go/tools/cmd/staticcheck@v0.7.0 ./...
go build ./cmd/agent-runtime
go run ./cmd/check-fuzz
go run ./cmd/check-cold-compile
```

`check-fuzz` starts every target named in `fuzz/v1alpha1.json` for a bounded
number of iterations, and fails if that list and the targets in the tree
disagree. Adding a `Fuzz*` function therefore means adding its entry.

`check-cold-compile` builds every platform this module claims to support, from a
cold cache. It is the only thing that makes the claim true, since CI runs on two
of them.

## Two things worth knowing before changing the schemas

`schemas/` is a published contract: anyone can read this repository's output
without running its Go. `schema_validation_test.go` runs real documents through
a real validator and then mutates them, so a schema and its producer cannot
drift apart silently. Change both halves in the same commit.

Error strings open with the domain's proper nouns — `Task`, `Goal` — because
those are the exported type names the schemas and the documentation use. That is
why `ST1005` is disabled in `staticcheck.conf`; do not "fix" the casing.
