package service

import (
	"context"
	"testing"

	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/entity"
)

// --- fakes -----------------------------------------------------------------

type fakeFields struct{ fields []fieldentity.Field }

func (f fakeFields) ModuleStorage(_ context.Context, _, _ string) (string, bool, error) {
	return "dynamic", true, nil
}
func (f fakeFields) List(_ context.Context, _, _ string) ([]fieldentity.Field, error) {
	return f.fields, nil
}

type fakeRules struct{ rules []entity.ValidationRule }

func (f fakeRules) ModuleExists(context.Context, string, string) (bool, error) { return true, nil }
func (f fakeRules) FieldExists(context.Context, string, string, string) (bool, error) {
	return true, nil
}
func (f fakeRules) Create(context.Context, *entity.ValidationRule) error { return nil }
func (f fakeRules) List(context.Context, string, string) ([]entity.ValidationRule, error) {
	return f.rules, nil
}
func (f fakeRules) ActiveByModule(context.Context, string, string) ([]entity.ValidationRule, error) {
	return f.rules, nil
}
func (f fakeRules) GetByID(context.Context, string, string, string) (*entity.ValidationRule, error) {
	return nil, nil
}
func (f fakeRules) Update(context.Context, *entity.ValidationRule) error         { return nil }
func (f fakeRules) Delete(context.Context, string, string, string) (bool, error) { return true, nil }

// --- helpers ---------------------------------------------------------------

func ptrInt(i int) *int       { return &i }
func ptrStr(s string) *string { return &s }

func errorFields(t *testing.T, svc *Service, data map[string]any) map[string]string {
	t.Helper()
	res, err := svc.Validate(context.Background(), "org", "mod", data)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	out := map[string]string{}
	for _, e := range res.Errors {
		out[e.Field] = e.Message
	}
	return out
}

// --- tests -----------------------------------------------------------------

func TestValidate_MetadataAndRules(t *testing.T) {
	fields := []fieldentity.Field{
		{ID: "f1", APIName: "name", Label: "Name", FieldType: fieldentity.TypeText, IsRequired: true, MinLength: ptrInt(3)},
		{ID: "f2", APIName: "email", Label: "Email", FieldType: fieldentity.TypeEmail},
		{ID: "f3", APIName: "status", Label: "Status", FieldType: fieldentity.TypeDropdown, Options: []byte(`["NEW","WON"]`)},
		{ID: "f4", APIName: "age", Label: "Age", FieldType: fieldentity.TypeNumber},
	}
	rules := []entity.ValidationRule{
		{FieldID: ptrStr("f4"), RuleType: entity.RuleMin, Params: []byte(`{"value":18}`), ErrorMessage: ptrStr("Must be 18+"), IsActive: true},
	}
	svc := New(fakeRules{rules: rules}, fakeFields{fields: fields})

	t.Run("valid payload passes", func(t *testing.T) {
		res, err := svc.Validate(context.Background(), "org", "mod", map[string]any{
			"name":   "Ada",
			"email":  "ada@example.com",
			"status": "NEW",
			"age":    float64(25),
		})
		if err != nil {
			t.Fatalf("validate: %v", err)
		}
		if !res.Valid || len(res.Errors) != 0 {
			t.Fatalf("expected valid, got %+v", res)
		}
	})

	t.Run("required + minlength + email + option + custom min message", func(t *testing.T) {
		errs := errorFields(t, svc, map[string]any{
			"name":   "Al",
			"email":  "not-an-email",
			"status": "LOST",
			"age":    float64(10),
		})
		if errs["name"] == "" {
			t.Errorf("expected name length error")
		}
		if errs["email"] == "" {
			t.Errorf("expected email format error")
		}
		if errs["status"] == "" {
			t.Errorf("expected invalid option error")
		}
		if errs["age"] != "Must be 18+" {
			t.Errorf("expected custom min message, got %q", errs["age"])
		}
	})

	t.Run("missing required field", func(t *testing.T) {
		errs := errorFields(t, svc, map[string]any{"email": "ada@example.com"})
		if errs["name"] == "" {
			t.Errorf("expected required error for name")
		}
	})
}

func TestValidate_RequiredIfModuleRule(t *testing.T) {
	fields := []fieldentity.Field{
		{ID: "f1", APIName: "type", Label: "Type", FieldType: fieldentity.TypeDropdown, Options: []byte(`["person","company"]`)},
		{ID: "f2", APIName: "company_name", Label: "Company Name", FieldType: fieldentity.TypeText},
	}
	rules := []entity.ValidationRule{
		{RuleType: entity.RuleRequiredIf, Params: []byte(`{"field":"type","equals":"company","target":"company_name"}`), IsActive: true},
	}
	svc := New(fakeRules{rules: rules}, fakeFields{fields: fields})

	t.Run("triggers when condition met", func(t *testing.T) {
		errs := errorFields(t, svc, map[string]any{"type": "company"})
		if errs["company_name"] == "" {
			t.Errorf("expected company_name required")
		}
	})

	t.Run("skips when condition not met", func(t *testing.T) {
		errs := errorFields(t, svc, map[string]any{"type": "person"})
		if len(errs) != 0 {
			t.Errorf("expected no errors, got %+v", errs)
		}
	})
}

func TestSchema_Compiles(t *testing.T) {
	fields := []fieldentity.Field{
		{ID: "f1", APIName: "name", Label: "Name", FieldType: fieldentity.TypeText, IsRequired: true, MaxLength: ptrInt(50)},
		{ID: "f2", APIName: "status", Label: "Status", FieldType: fieldentity.TypeDropdown, Options: []byte(`["NEW","WON"]`)},
	}
	rules := []entity.ValidationRule{
		{FieldID: ptrStr("f1"), RuleType: entity.RuleMinLength, Params: []byte(`{"value":2}`), ErrorMessage: ptrStr("Too short"), IsActive: true},
	}
	svc := New(fakeRules{rules: rules}, fakeFields{fields: fields})

	schema, err := svc.Schema(context.Background(), "org", "mod")
	if err != nil {
		t.Fatalf("schema: %v", err)
	}
	if len(schema.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(schema.Fields))
	}

	name := schema.Fields[0]
	if !name.Required || name.MinLength == nil || *name.MinLength != 2 || name.MaxLength == nil || *name.MaxLength != 50 {
		t.Errorf("unexpected name schema: %+v", name)
	}
	if name.Messages["min_length"] != "Too short" {
		t.Errorf("expected custom min_length message, got %+v", name.Messages)
	}

	status := schema.Fields[1]
	if len(status.Options) != 2 {
		t.Errorf("expected status options, got %+v", status.Options)
	}
}
