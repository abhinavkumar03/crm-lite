package secrets

import "testing"

func TestBoxRoundTrip(t *testing.T) {
	box, err := NewBox("unit-test-communication-secrets-key!!")
	if err != nil {
		t.Fatal(err)
	}
	in := map[string]any{"api_key": "secret-123", "password": "p@ss"}
	ct, err := box.EncryptJSON(in)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]any{}
	if err := box.DecryptJSON(ct, &out); err != nil {
		t.Fatal(err)
	}
	if out["api_key"] != "secret-123" || out["password"] != "p@ss" {
		t.Fatalf("round trip mismatch: %#v", out)
	}
}
