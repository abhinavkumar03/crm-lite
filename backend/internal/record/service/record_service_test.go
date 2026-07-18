package service

import (
	"context"
	"encoding/json"
	"testing"

	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record/repository"
	vdto "github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/dto"
)

// --- fakes -----------------------------------------------------------------

type fakeFields struct {
	strategy map[string]string              // moduleID -> storage strategy
	byModule map[string][]fieldentity.Field // moduleID -> fields
}

func (f *fakeFields) ModuleStorage(_ context.Context, _, moduleID string) (string, bool, error) {
	s, ok := f.strategy[moduleID]
	return s, ok, nil
}

func (f *fakeFields) List(_ context.Context, _, moduleID string) ([]fieldentity.Field, error) {
	return f.byModule[moduleID], nil
}

type fakeValidator struct{ result vdto.ValidateResult }

func (v *fakeValidator) Validate(_ context.Context, _, _ string, _ map[string]any) (vdto.ValidateResult, error) {
	return v.result, nil
}

type fakeRepo struct {
	created   *entity.Record
	list      []entity.Record
	display   map[string]string
	users     map[string]string
	displayFn string
}

func (r *fakeRepo) Create(_ context.Context, rec *entity.Record) error {
	rec.ID = "new-id"
	r.created = rec
	return nil
}
func (r *fakeRepo) GetByID(_ context.Context, _, _, _ string) (*entity.Record, error) {
	return nil, nil
}
func (r *fakeRepo) Update(_ context.Context, _ *entity.Record) error { return nil }
func (r *fakeRepo) Delete(_ context.Context, _, _, _ string) (bool, error) {
	return true, nil
}
func (r *fakeRepo) List(_ context.Context, _, _ string, _ dto.ListQuery, _ map[string]repository.FieldMeta, _ repository.ExtraWhere) ([]entity.Record, int, error) {
	return r.list, len(r.list), nil
}
func (r *fakeRepo) DisplayValues(_ context.Context, _, _ string, _ []string, displayField string) (map[string]string, error) {
	r.displayFn = displayField
	return r.display, nil
}
func (r *fakeRepo) UserDisplays(_ context.Context, _ []string) (map[string]string, error) {
	return r.users, nil
}

// --- tests -----------------------------------------------------------------

func TestCreate_RejectsInvalidPayload(t *testing.T) {
	svc := New(
		&fakeRepo{},
		&fakeFields{strategy: map[string]string{"m": "dynamic"}},
		&fakeValidator{result: vdto.ValidateResult{
			Valid:  false,
			Errors: []vdto.FieldError{{Field: "name", Message: "required"}},
		}},
		nil,
		nil,
	)

	_, err := svc.Create(context.Background(), "o", "m", "u1", dto.CreateRecordRequest{Data: map[string]any{}})

	var verr *ValidationError
	if err == nil || !asValidation(err, &verr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
	if len(verr.Errors) != 1 || verr.Errors[0].Field != "name" {
		t.Fatalf("unexpected errors: %#v", verr.Errors)
	}
}

func TestCreate_DefaultsOwnerToActingUser(t *testing.T) {
	repo := &fakeRepo{}
	svc := New(repo,
		&fakeFields{strategy: map[string]string{"m": "dynamic"}},
		&fakeValidator{result: vdto.ValidateResult{Valid: true}},
		nil,
		nil,
	)

	resp, err := svc.Create(context.Background(), "o", "m", "u1", dto.CreateRecordRequest{
		Data: map[string]any{"name": "Acme"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "new-id" {
		t.Fatalf("expected persisted id, got %q", resp.ID)
	}
	if repo.created.OwnerID == nil || *repo.created.OwnerID != "u1" {
		t.Fatalf("owner should default to acting user, got %v", repo.created.OwnerID)
	}
	if repo.created.CreatedBy == nil || *repo.created.CreatedBy != "u1" {
		t.Fatalf("created_by should be acting user")
	}
}

func TestCreate_RejectsNativeModule(t *testing.T) {
	svc := New(&fakeRepo{},
		&fakeFields{strategy: map[string]string{"m": "native"}},
		&fakeValidator{result: vdto.ValidateResult{Valid: true}},
		nil,
		nil,
	)

	_, err := svc.Create(context.Background(), "o", "m", "u1", dto.CreateRecordRequest{Data: map[string]any{}})
	if err != ErrNotDynamic {
		t.Fatalf("expected ErrNotDynamic, got %v", err)
	}
}

func TestList_ExpandsLookupAndUserRelations(t *testing.T) {
	lookupModuleID := "companies"
	fields := []fieldentity.Field{
		{APIName: "company", FieldType: fieldentity.TypeLookup, LookupModuleID: &lookupModuleID},
		{APIName: "assignee", FieldType: fieldentity.TypeUser},
	}
	targetFields := []fieldentity.Field{
		{APIName: "name", FieldType: fieldentity.TypeText},
	}

	repo := &fakeRepo{
		list: []entity.Record{
			{ID: "r1", Data: mustJSON(map[string]any{"company": "c1", "assignee": "u9"})},
		},
		display: map[string]string{"c1": "Acme Inc"},
		users:   map[string]string{"u9": "Dana Dev"},
	}
	svc := New(repo,
		&fakeFields{
			strategy: map[string]string{"m": "dynamic"},
			byModule: map[string][]fieldentity.Field{"m": fields, lookupModuleID: targetFields},
		},
		&fakeValidator{result: vdto.ValidateResult{Valid: true}},
		nil,
		nil,
	)

	result, err := svc.List(context.Background(), "o", "m", "u1", dto.ListQuery{Expand: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(result.Records))
	}

	rels := result.Records[0].Relations
	if rels["company"].Label != "Acme Inc" || rels["company"].ID != "c1" {
		t.Fatalf("lookup relation not expanded: %#v", rels["company"])
	}
	if rels["assignee"].Label != "Dana Dev" {
		t.Fatalf("user relation not expanded: %#v", rels["assignee"])
	}
	if repo.displayFn != "name" {
		t.Fatalf("expected display field 'name', got %q", repo.displayFn)
	}
}

func mustJSON(m map[string]any) []byte {
	b, _ := json.Marshal(m)
	return b
}

func asValidation(err error, target **ValidationError) bool {
	v, ok := err.(*ValidationError)
	if ok {
		*target = v
	}
	return ok
}
