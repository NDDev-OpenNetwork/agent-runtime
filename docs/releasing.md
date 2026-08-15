# Releasing

A release is a tag. Everything else follows from it.

## What a tag does on its own

Pushing `vMAJOR.MINOR.PATCH` publishes the library: the Go module proxy fetches
it within seconds and serves it to `go get` forever after. The proxy records an
immutable hash, so a tag that is later moved or deleted does not change what a
consumer already resolved — and does not free the version number either.

**A tag is therefore irreversible.** Cut it only from a commit whose CI is
already green on `main`.

## Steps

1. Move `## [Unreleased]` in `CHANGELOG.md` to a dated `## [X.Y.Z]` heading and
   merge that through a pull request like any other change.
2. Confirm CI is green on the resulting `main` commit.
3. Tag it and push:

   ```
   git tag -a vX.Y.Z -m "vX.Y.Z"
   git push origin vX.Y.Z
   ```

The `Release` workflow then builds the CLI for linux and darwin on both
architectures, writes `SHA256SUMS`, and creates the GitHub Release with
generated notes.

## Verifying a release

For the library, `go get` is the verification: the proxy hash and `go.sum` catch
any substitution without further ceremony.

For the CLI binaries, download the archive and `SHA256SUMS` from the release and
check them:

```
sha256sum --check --ignore-missing SHA256SUMS
```
