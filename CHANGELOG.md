# Changelog

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
