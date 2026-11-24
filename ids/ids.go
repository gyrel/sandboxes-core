// Package ids provides correlation identifiers for operations and requests.
// IDs are not identity or authorization material and must never be used as such.
package ids

import (
	"context"
	"crypto/rand"
	"errors"
	"io"
	"strings"
	"sync"
	"time"
)

var (
	// ErrMalformed indicates an ID is not a canonical ULID-based identifier.
	ErrMalformed = errors.New("malformed correlation ID")
	// ErrUnknownKind indicates an otherwise well-formed ID has an unknown prefix.
	ErrUnknownKind = errors.New("unknown correlation ID kind")
)

const alphabet = "0123456789abcdefghjkmnpqrstvwxyz"

// generator serializes the timestamp and 80-bit entropy across both kinds.
type generator struct {
	mu          sync.Mutex
	now         func() time.Time
	ms          uint64
	entropy     [10]byte
	initialized bool
}

var defaultGenerator = generator{now: time.Now}

// NewOpID returns a new, process-monotonic operation correlation ID.
func NewOpID() string { return defaultGenerator.newID("op_") }

// NewReqID returns a new, process-monotonic request correlation ID.
func NewReqID() string { return defaultGenerator.newID("req_") }

func (g *generator) newID(prefix string) string {
	g.mu.Lock()
	defer g.mu.Unlock()
	ms := uint64(g.now().UnixMilli())
	if !g.initialized || ms > g.ms {
		g.reseed(ms)
	} else {
		g.increment()
	}
	var raw [16]byte
	for i := 5; i >= 0; i-- {
		raw[i] = byte((g.ms >> (8 * (5 - i))) & 0xff)
	}
	copy(raw[6:], g.entropy[:])
	var result [30]byte
	n := copy(result[:], prefix)
	// The first symbol contains three bits; the rest consume five at a time.
	result[n] = alphabet[raw[0]>>5]
	bits := uint32(raw[0] & 31)
	count := 5
	out := n + 1
	for _, v := range raw[1:] {
		bits = bits<<8 | uint32(v)
		count += 8
		for count >= 5 {
			count -= 5
			result[out] = alphabet[(bits>>count)&31]
			out++
		}
	}
	return string(result[:n+26])
}

// reseed installs fresh entropy for a new timestamp. crypto/rand does not fail
// on supported platforms; if a read ever did, the generator stays monotonic
// with a deterministic fallback instead of panicking inside a library.
// Correlation IDs are never identity or authorization material.
func (g *generator) reseed(ms uint64) {
	if _, err := io.ReadFull(rand.Reader, g.entropy[:]); err != nil {
		if g.initialized {
			g.increment()
			return
		}
		g.entropy = [10]byte{}
	}
	g.ms = ms
	g.initialized = true
}

// increment carries through the big-endian entropy and into time on overflow.
func (g *generator) increment() {
	for i := len(g.entropy) - 1; i >= 0; i-- {
		g.entropy[i]++
		if g.entropy[i] != 0 {
			return
		}
	}
	g.ms++
}

// Validate checks a canonical operation or request correlation ID. An unknown
// nonempty prefix with a valid ULID returns ErrUnknownKind; all other invalid
// forms return ErrMalformed.
func Validate(id string) error {
	kind, value, ok := strings.Cut(id, "_")
	if !ok || kind == "" || len(value) != 26 || value[0] > '7' {
		return ErrMalformed
	}
	for i := range len(value) {
		if !strings.ContainsRune(alphabet, rune(value[i])) {
			return ErrMalformed
		}
	}
	if kind != "op" && kind != "req" {
		return ErrUnknownKind
	}
	return nil
}

// ValidateParentReqID validates a parent_req_id as a request correlation ID.
func ValidateParentReqID(id string) error {
	if err := Validate(id); err != nil {
		return err
	}
	if !strings.HasPrefix(id, "req_") {
		return ErrUnknownKind
	}
	return nil
}

type contextKey uint8

const (
	opKey contextKey = iota
	reqKey
)

// WithOpID returns a context carrying id. A nil context starts from Background.
// Callers should validate IDs received from outside the process.
func WithOpID(ctx context.Context, id string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, opKey, id)
}

// WithReqID returns a context carrying id. A nil context starts from Background.
// Callers should validate IDs received from outside the process.
func WithReqID(ctx context.Context, id string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, reqKey, id)
}

// FromContext returns correlation IDs, or empty strings when absent or nil.
func FromContext(ctx context.Context) (op, req string) {
	if ctx == nil {
		return "", ""
	}
	op, _ = ctx.Value(opKey).(string)
	req, _ = ctx.Value(reqKey).(string)
	return op, req
}
