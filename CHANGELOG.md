# Changelog

## Unreleased

### Changed

- Rename the Go module path from `github.com/rheic/sandboxes-core` to
  `github.com/gyrel/sandboxes-core`. Breaking change; pre-1.0 alpha policy.

### Added

- `sandbox.NetworkRule` carries HTTP L7 conditions (`Methods`, `Path`,
  `Query []QueryCondition`) mirroring the gateway policy engine semantics.
  Zero values keep existing rules unconditioned; persistence encodes rules as
  whole structs, so stored rows from before this release decode to no HTTP
  conditions.

## v0.0.13 — 2026-11-16

### Changed

- Remove the mandatory external release-review gate while retaining candidate validation, CI, and immutable release-tag protections.
- Accept an `Unreleased` heading before the newest version during release validation.

## v0.0.9 — 2026-11-03

### Added

- Add domain-based network policy rules and expanded policy settings, secret, kit, and agent-preview models.

### Changed

- Normalize the Go module path to the lowercase `github.com/quoiner/sandboxes-core` namespace.

### Build

- Publish the module CycloneDX SBOM as a CI artifact.
