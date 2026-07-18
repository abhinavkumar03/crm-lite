package notify

import "testing"

func TestRender(t *testing.T) {
	data := map[string]any{"name": "Dana", "amount": 4200}

	cases := map[string]struct {
		in, want string
	}{
		"simple":        {"Hi {{name}}", "Hi Dana"},
		"spaced token":  {"Hi {{ name }}", "Hi Dana"},
		"numeric value": {"Total {{amount}}", "Total 4200"},
		"unknown key":   {"Hi {{missing}}!", "Hi !"},
		"no tokens":     {"plain text", "plain text"},
		"multiple":      {"{{name}} owes {{amount}}", "Dana owes 4200"},
	}

	for name, tc := range cases {
		if got := Render(tc.in, data); got != tc.want {
			t.Errorf("%s: Render(%q) = %q, want %q", name, tc.in, got, tc.want)
		}
	}
}

func TestRender_EmptyInputs(t *testing.T) {
	if got := Render("", map[string]any{"a": 1}); got != "" {
		t.Errorf("empty text should stay empty, got %q", got)
	}
	if got := Render("hi {{a}}", nil); got != "hi {{a}}" {
		t.Errorf("nil data should leave text untouched, got %q", got)
	}
}

func TestBuildWhatsAppProvider_FallsBackToSimulation(t *testing.T) {
	// "meta" without credentials must not produce a live provider.
	p := BuildWhatsAppProvider(WhatsAppConfig{Provider: "meta"}, nil)
	if p.Name() != "simulation" {
		t.Fatalf("expected simulation fallback, got %q", p.Name())
	}

	full := BuildWhatsAppProvider(WhatsAppConfig{
		Provider: "meta", APIURL: "https://x", Token: "t", PhoneID: "p",
	}, nil)
	if full.Name() != "meta-cloud" {
		t.Fatalf("expected meta-cloud provider, got %q", full.Name())
	}
	if full.Channel() != ChannelWhatsApp {
		t.Fatalf("expected whatsapp channel")
	}
}
