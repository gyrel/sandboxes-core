// Package sandbox defines the sandbox aggregate and persistence port.
package sandbox

import (
	"context"
	"time"
)

// NetworkRule is a CIDR or domain/protocol/port egress rule. Domain and CIDR
// are mutually exclusive.
type NetworkRule struct {
	Protocol string
	CIDR     string
	Port     uint16
	Domain   string
	RuleID   string
}

// Meta stores non-lifecycle metadata and the clone configuration/status record.
type Meta struct {
	Agent         string
	Detached      bool
	PolicyProfile string
	LastUsedAt    time.Time
	Clone         *CloneState
}

// NetworkPolicy is the create-time allow/deny policy snapshot.
type NetworkPolicy struct {
	Allow []NetworkRule
	Deny  []NetworkRule
}

// Port maps a host loopback address/port to a guest port.
type Port struct {
	Host     uint16
	Guest    uint16
	Protocol string
	HostIP   string
}

// Workspace maps a host path into the guest. An empty GuestPath preserves the
// legacy behavior and mounts Path at the same absolute path inside the guest.
type Workspace struct {
	Path      string
	GuestPath string
	ReadOnly  bool
}

// CloneOptions describes a private guest checkout copied from one workspace.
// A nil GuestWorkdir lets the daemon use the guest image's default directory.
type CloneOptions struct {
	SourceWorkspacePath string
	GuestWorkdir        *string
}

// GitEndpoint is the host-loopback endpoint that forwards to a guest Git daemon.
type GitEndpoint struct {
	HostIP     string
	Port       uint16
	Repository string
}

// CloneInfo is the persisted clone directory and current guest Git endpoint.
// GitEndpoint is nil while the sandbox is stopped.
type CloneInfo struct {
	GuestWorkdir string
	RemoteName   string
	GitEndpoint  *GitEndpoint
}

// CloneState persists the clone request and live clone endpoint together.
type CloneState struct {
	Options CloneOptions
	Info    *CloneInfo
}

// Sandbox is the durable sandbox record.
type Sandbox struct {
	ID         string
	Name       string
	State      State
	Generation int64
	Image      string
	Workspace  string
	Workspaces []Workspace
	Ports      []Port
	Policy     *NetworkPolicy
	Meta       *Meta
	CPUs       uint32
	MemoryMiB  uint32
	StoppedAt  time.Time
	Error      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ListOptions controls pagination.
type ListOptions struct {
	PageSize  int
	PageToken string
}

// Store is implemented by the SQLite repository.
type Store interface {
	Create(ctx context.Context, sb Sandbox) error
	Get(ctx context.Context, id string) (Sandbox, error)
	GetByName(ctx context.Context, name string) (Sandbox, error)
	List(ctx context.Context, opts ListOptions) ([]Sandbox, string, error)
	Update(ctx context.Context, sb Sandbox) error
	Delete(ctx context.Context, id string) error
}
