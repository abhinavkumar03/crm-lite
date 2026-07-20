package service

import (
	"testing"

	"github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
)

func TestCanonicalizeType(t *testing.T) {
	tests := map[string]string{
		"select":       entity.TypeDropdown,
		"multi_select": entity.TypeMultiselect,
		"user_lookup":  entity.TypeUser,
		"percent":      entity.TypePercentage,
		"GST":          entity.TypeGST, // lowercased
		"text":         entity.TypeText,
	}
	for in, want := range tests {
		got := canonicalizeType(in)
		if got != want {
			t.Fatalf("canonicalizeType(%q)=%q want %q", in, got, want)
		}
	}
}

func TestNormalizeLockMode(t *testing.T) {
	mode, err := normalizeLockMode("")
	if err != nil || mode != entity.LockNever {
		t.Fatalf("empty -> never, got %q %v", mode, err)
	}
	mode, err = normalizeLockMode("after_create")
	if err != nil || mode != entity.LockAfterCreate {
		t.Fatalf("got %q %v", mode, err)
	}
	_, err = normalizeLockMode("sometime")
	if err != ErrInvalidLockMode {
		t.Fatalf("expected ErrInvalidLockMode, got %v", err)
	}
}

func TestNormalizeACL(t *testing.T) {
	ed, vw, err := normalizeACL("", "")
	if err != nil || ed != "ALL" || vw != "ALL" {
		t.Fatalf("got %s %s %v", ed, vw, err)
	}
	_, _, err = normalizeACL("OWNER", "ALL")
	if err != ErrInvalidACL {
		t.Fatalf("expected ErrInvalidACL, got %v", err)
	}
}
