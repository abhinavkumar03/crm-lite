package service

import (
	"context"
	"errors"
	"testing"

	"github.com/hibiken/asynq"

	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/importer/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/importer/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
)

type fakeRepo struct {
	created  *entity.ImportJob
	finished bool
}

func (f *fakeRepo) Create(_ context.Context, j *entity.ImportJob) error {
	j.ID = "import-1"
	j.Status = entity.StatusPending
	f.created = j
	return nil
}
func (f *fakeRepo) GetByID(context.Context, string, string) (*entity.ImportJob, error) {
	return f.created, nil
}
func (f *fakeRepo) List(context.Context, string, string, dto.ListQuery) ([]entity.ImportJob, int, error) {
	return nil, 0, nil
}
func (f *fakeRepo) Finish(context.Context, string, string, int, int, int, []byte) error {
	f.finished = true
	return nil
}

type fakeFields struct {
	strategy string
	fields   []fieldentity.Field
}

func (f *fakeFields) ModuleStorage(context.Context, string, string) (string, bool, error) {
	if f.strategy == "" {
		return "", false, nil
	}
	return f.strategy, true, nil
}
func (f *fakeFields) List(context.Context, string, string) ([]fieldentity.Field, error) {
	return f.fields, nil
}

type fakeEnqueuer struct {
	published int
	err       error
}

func (f *fakeEnqueuer) Publish(context.Context, jobs.Job, ...asynq.Option) error {
	f.published++
	return f.err
}

func dynamicFields() *fakeFields {
	return &fakeFields{
		strategy: "dynamic",
		fields: []fieldentity.Field{
			{APIName: "first_name", Label: "First Name", FieldType: fieldentity.TypeText},
		},
	}
}

func TestCreateHappyPath(t *testing.T) {
	repo := &fakeRepo{}
	enq := &fakeEnqueuer{}
	svc := New(repo, dynamicFields(), enq)

	csv := []byte("first_name\nAda\nGrace\n")
	mapping := map[string]string{"first_name": "first_name"}

	resp, err := svc.Create(context.Background(), "org", "mod", "user", "leads.csv", csv, mapping, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != entity.StatusPending {
		t.Fatalf("status = %q, want pending", resp.Status)
	}
	if resp.TotalRows != 2 {
		t.Fatalf("total_rows = %d, want 2", resp.TotalRows)
	}
	if enq.published != 1 {
		t.Fatalf("published = %d, want 1", enq.published)
	}
	if repo.created == nil {
		t.Fatal("expected job to be persisted")
	}
}

func TestCreateRejectsNativeModule(t *testing.T) {
	svc := New(&fakeRepo{}, &fakeFields{strategy: "native"}, &fakeEnqueuer{})
	_, err := svc.Create(context.Background(), "org", "mod", "user", "x.csv", []byte("a\n1\n"), map[string]string{"a": "b"}, nil)
	if !errors.Is(err, ErrNotDynamic) {
		t.Fatalf("err = %v, want ErrNotDynamic", err)
	}
}

func TestCreateRejectsEmptyMapping(t *testing.T) {
	svc := New(&fakeRepo{}, dynamicFields(), &fakeEnqueuer{})
	// Mapping targets a field that does not exist -> sanitized to empty.
	_, err := svc.Create(context.Background(), "org", "mod", "user", "x.csv", []byte("first_name\nAda\n"), map[string]string{"first_name": "ghost"}, nil)
	if !errors.Is(err, ErrNoMapping) {
		t.Fatalf("err = %v, want ErrNoMapping", err)
	}
}

func TestCreateMarksFailedWhenEnqueueFails(t *testing.T) {
	repo := &fakeRepo{}
	enq := &fakeEnqueuer{err: errors.New("redis down")}
	svc := New(repo, dynamicFields(), enq)

	_, err := svc.Create(context.Background(), "org", "mod", "user", "x.csv", []byte("first_name\nAda\n"), map[string]string{"first_name": "first_name"}, nil)
	if err == nil {
		t.Fatal("expected error when enqueue fails")
	}
	if !repo.finished {
		t.Fatal("expected job to be marked failed after enqueue failure")
	}
}

func TestAnalyzeSuggestsMapping(t *testing.T) {
	svc := New(&fakeRepo{}, dynamicFields(), &fakeEnqueuer{})
	res, err := svc.Analyze(context.Background(), "org", "mod", "leads.csv", []byte("First Name\nAda\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.SuggestedMapping["First Name"] != "first_name" {
		t.Fatalf("suggested mapping = %v", res.SuggestedMapping)
	}
	if res.RowCount != 1 {
		t.Fatalf("row_count = %d, want 1", res.RowCount)
	}
}
