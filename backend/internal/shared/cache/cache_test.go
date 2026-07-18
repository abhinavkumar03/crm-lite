package cache

import (
	"context"
	"testing"
)

func TestKey(t *testing.T) {
	got := Key("rbac", "perms", "role-1")
	if got != "rbac:perms:role-1" {
		t.Fatalf("unexpected key %q", got)
	}
}

func TestNilCacheIsSafe(t *testing.T) {
	var c *Cache
	ctx := context.Background()

	var dest string
	if c.GetJSON(ctx, "x", &dest) {
		t.Fatal("nil cache should miss")
	}
	c.SetJSON(ctx, "x", "y", TTLShort)
	c.Delete(ctx, "x")
	c.InvalidateDashboard(ctx, "u1")
	c.InvalidateRole(ctx, "r1")
	c.InvalidateMembership(ctx, "u1")
}

func TestNewNilClient(t *testing.T) {
	if New(nil) != nil {
		t.Fatal("New(nil) should return nil Cache")
	}
}
