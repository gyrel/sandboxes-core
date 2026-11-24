package ids

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestMonotonicAndRollback(t *testing.T) {
	now := time.UnixMilli(1700000000000)
	g := generator{now: func() time.Time { return now }}
	prev := g.newID("req_")
	for i := 0; i < 10000; i++ {
		if i == 5000 {
			now = now.Add(-time.Hour)
		}
		next := g.newID("req_")
		if next <= prev {
			t.Fatalf("IDs not increasing: %s >= %s", prev, next)
		}
		if err := Validate(next); err != nil {
			t.Fatal(err)
		}
		prev = next
	}
}

func TestEntropyCarry(t *testing.T) {
	g := generator{now: func() time.Time { return time.UnixMilli(100) }, ms: 100, initialized: true}
	for i := range g.entropy {
		g.entropy[i] = 255
	}
	g.newID("op_")
	if g.ms != 101 || g.entropy != [10]byte{} {
		t.Fatal("entropy overflow did not carry to timestamp")
	}
}

func TestValidate(t *testing.T) {
	valid := strings.Repeat("0", 26)
	for _, tc := range []struct {
		id   string
		want error
	}{
		{"op_" + valid, nil}, {"req_" + valid, nil}, {"other_" + valid, ErrUnknownKind},
		{"_" + valid, ErrMalformed}, {"", ErrMalformed}, {"op_" + valid[:25], ErrMalformed},
		{"op_" + valid + "0", ErrMalformed}, {"op_" + strings.ToUpper(strings.Repeat("a", 26)), ErrMalformed},
		{"op_" + valid[:25] + "u", ErrMalformed}, {"op_8" + valid[:25], ErrMalformed},
		{"op_" + valid[:25] + "i", ErrMalformed}, {"OP_" + valid, ErrUnknownKind},
	} {
		if err := Validate(tc.id); !errors.Is(err, tc.want) {
			t.Errorf("Validate(%q) = %v, want %v", tc.id, err, tc.want)
		}
	}
	if err := ValidateParentReqID("op_" + valid); !errors.Is(err, ErrUnknownKind) {
		t.Fatal(err)
	}
	if err := ValidateParentReqID("req_" + valid); err != nil {
		t.Fatal(err)
	}
}

func FuzzValidate(f *testing.F) {
	f.Add("req_" + strings.Repeat("0", 26))
	f.Add("")
	f.Fuzz(func(t *testing.T, id string) {
		err := Validate(id)
		if err == nil {
			if !strings.HasPrefix(id, "op_") && !strings.HasPrefix(id, "req_") {
				t.Fatal("invalid prefix")
			}
			_, s, _ := strings.Cut(id, "_")
			if len(s) != 26 || s[0] > '7' {
				t.Fatal("invalid length or timestamp")
			}
			for _, c := range s {
				if !strings.ContainsRune(alphabet, c) {
					t.Fatal("invalid character")
				}
			}
		}
	})
}

func TestContext(t *testing.T) {
	op, req := NewOpID(), NewReqID()
	ctx := WithReqID(WithOpID(context.Background(), op), req)
	if gotOp, gotReq := FromContext(ctx); gotOp != op || gotReq != req {
		t.Fatal(gotOp, gotReq)
	}
	var nilContext context.Context
	if a, b := FromContext(nilContext); a != "" || b != "" {
		t.Fatal(a, b)
	}
	if a, _ := FromContext(WithOpID(nilContext, op)); a != op {
		t.Fatal(a)
	}
	if _, b := FromContext(WithReqID(nilContext, req)); b != req {
		t.Fatal(b)
	}
}

func BenchmarkNewOpID(b *testing.B) {
	for b.Loop() {
		NewOpID()
	}
}
func BenchmarkNewReqID(b *testing.B) {
	for b.Loop() {
		NewReqID()
	}
}

func TestKnownULIDEncoding(t *testing.T) {
	g := generator{now: func() time.Time { return time.UnixMilli(1469922850259) }, ms: 1469922850259, initialized: true,
		entropy: [10]byte{0x7a, 0xc9, 0x5d, 0x64, 0x8a, 0xf1, 0x12, 0xe6, 0x8f, 0xe5}}
	if got := g.newID("req_"); got != "req_01arz3ndekfb4nts4ay49ed3z6" {
		t.Fatalf("unexpected encoding: %s", got)
	}
}
