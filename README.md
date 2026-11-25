# sandboxes-core

<img src="docs/assets/gyrel-shark.png" alt="Gyrel shark logo" width="180">


## Import path migration

The canonical Go module path is `github.com/gyrel/sandboxes-core`. Update imports and `go.mod` requirements when moving to a release that includes this change. Previous module paths are not kept as compatibility aliases.

Status: pre-alpha.

Domain model, state machines, operation journal contracts, workspace rules, and
ports implemented by infrastructure modules.

## Boundaries

- no SQLite, containerd, HTTP, or os/exec;
- no dependency on runtime, gateway, daemon, or SDK;
- domain invariants live here and are enforced by the daemon.
