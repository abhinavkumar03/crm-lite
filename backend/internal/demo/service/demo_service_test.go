package service

import (
	"encoding/json"
	"testing"

	"github.com/abhinavkumar03/crm-lite/backend/internal/demo/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/demo/repository"
)

func TestNextStepKey(t *testing.T) {
	steps := []repository.Step{
		{Key: "a"}, {Key: "b"}, {Key: "c"},
	}
	if got := nextStepKey(steps, "a"); got != "b" {
		t.Fatalf("expected b, got %s", got)
	}
	if got := nextStepKey(steps, "c"); got != "" {
		t.Fatalf("expected empty at end, got %s", got)
	}
}

func TestRunValidatorUnknownFails(t *testing.T) {
	s := &Service{}
	step := &repository.Step{
		ValidatorKey:    "not_a_real_validator",
		ValidatorParams: json.RawMessage(`{}`),
	}
	ok, msg := s.runValidator(nil, "", step, "", nil)
	if ok {
		t.Fatal("unknown validator must fail")
	}
	if msg == "" {
		t.Fatal("expected failure message")
	}
}

func TestRunValidatorAcknowledge(t *testing.T) {
	s := &Service{}
	step := &repository.Step{
		ValidatorKey:    "acknowledge",
		ValidatorParams: json.RawMessage(`{}`),
	}
	ok, _ := s.runValidator(nil, "", step, "", nil)
	if !ok {
		t.Fatal("acknowledge should pass")
	}
}

func TestRunValidatorRouteVisited(t *testing.T) {
	s := &Service{}
	step := &repository.Step{
		ValidatorKey:    "route_visited",
		ValidatorParams: json.RawMessage(`{"route":"/settings/modules"}`),
	}
	ok, _ := s.runValidator(nil, "", step, "/settings/modules", nil)
	if !ok {
		t.Fatal("matching route should pass")
	}
	ok, _ = s.runValidator(nil, "", step, "/dashboard", nil)
	if ok {
		t.Fatal("wrong route should fail")
	}
}

func TestRunValidatorUIClick(t *testing.T) {
	s := &Service{}
	step := &repository.Step{
		ValidatorKey:    "ui_click",
		ValidatorParams: json.RawMessage(`{"selector":"[data-tutorial-action=\"create-module\"]"}`),
	}
	ok, _ := s.runValidator(nil, "", step, "", &dto.ClientEvent{
		Type: "click", Selector: `[data-tutorial-action="create-module"]`,
	})
	if !ok {
		t.Fatal("matching click should pass")
	}
	ok, _ = s.runValidator(nil, "", step, "", nil)
	if ok {
		t.Fatal("missing click should fail")
	}
}
