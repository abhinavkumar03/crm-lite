package parser

import (
	"strings"
	"testing"
)

func TestParseCSV(t *testing.T) {
	in := "Name, Email ,Amount\nAda,ada@example.com,100\nGrace, grace@example.com ,\n"
	p, err := ParseCSV(strings.NewReader(in))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantHeaders := []string{"Name", "Email", "Amount"}
	if len(p.Headers) != len(wantHeaders) {
		t.Fatalf("headers = %v, want %v", p.Headers, wantHeaders)
	}
	for i, h := range wantHeaders {
		if p.Headers[i] != h {
			t.Fatalf("header[%d] = %q, want %q", i, p.Headers[i], h)
		}
	}

	if len(p.Rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(p.Rows))
	}
	if p.Rows[0]["Email"] != "ada@example.com" {
		t.Fatalf("row0 email = %q", p.Rows[0]["Email"])
	}
	// Values are trimmed and short rows are padded with empty strings.
	if p.Rows[1]["Email"] != "grace@example.com" {
		t.Fatalf("row1 email = %q, want trimmed", p.Rows[1]["Email"])
	}
	if p.Rows[1]["Amount"] != "" {
		t.Fatalf("row1 amount = %q, want empty", p.Rows[1]["Amount"])
	}
}

func TestParseCSVDuplicateHeaders(t *testing.T) {
	p, err := ParseCSV(strings.NewReader("name,name\na,b\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Headers[0] != "name" || p.Headers[1] != "name_2" {
		t.Fatalf("headers = %v, want [name name_2]", p.Headers)
	}
}

func TestParseCSVBlankRowsSkipped(t *testing.T) {
	p, err := ParseCSV(strings.NewReader("a,b\n1,2\n , \n3,4\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Rows) != 2 {
		t.Fatalf("rows = %d, want 2 (blank row skipped)", len(p.Rows))
	}
}

func TestParseCSVEmpty(t *testing.T) {
	if _, err := ParseCSV(strings.NewReader("only,headers\n")); err != ErrEmptyFile {
		t.Fatalf("err = %v, want ErrEmptyFile", err)
	}
}

func TestParseFileUnsupported(t *testing.T) {
	if _, err := ParseFile("data.txt", []byte("x")); err != ErrUnsupported {
		t.Fatalf("err = %v, want ErrUnsupported", err)
	}
}
