package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/abhinavkumar03/crm-lite/backend/internal/workspace/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workspace/repository"
)

var (
	ErrNotFound = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ensureRecord(ctx context.Context, orgID, moduleID, recordID string) error {
	ok, err := s.repo.RecordExists(ctx, orgID, moduleID, recordID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

func (s *Service) GetDetailLayout(ctx context.Context, orgID, moduleID string) (*dto.LayoutResponse, error) {
	l, err := s.repo.GetDefaultDetailLayout(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	if l == nil {
		cfg, err := s.buildDefaultConfig(ctx, orgID, moduleID)
		if err != nil {
			return nil, err
		}
		l, err = s.repo.UpsertDefaultDetailLayout(ctx, orgID, moduleID, cfg)
		if err != nil {
			return nil, err
		}
	}
	return &dto.LayoutResponse{
		ID: l.ID, Name: l.Name, Type: l.Type, IsDefault: l.IsDefault, Config: l.Config,
	}, nil
}

func (s *Service) buildDefaultConfig(ctx context.Context, orgID, moduleID string) (json.RawMessage, error) {
	fields, err := s.repo.ListVisibleFieldAPINames(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	cfg := map[string]any{
		"sections": []map[string]any{
			{"key": "general", "label": "General Information", "fields": fields},
			{"key": "system", "label": "System Fields", "fields": []string{"owner_id", "assigned_to", "visibility", "created_at", "updated_at"}},
		},
		"tabs": []string{"overview", "notes", "attachments", "timeline", "related"},
	}
	return json.Marshal(cfg)
}

func (s *Service) ListNotes(ctx context.Context, orgID, moduleID, recordID string) ([]dto.NoteResponse, error) {
	if err := s.ensureRecord(ctx, orgID, moduleID, recordID); err != nil {
		return nil, err
	}
	notes, err := s.repo.ListNotes(ctx, orgID, moduleID, recordID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.NoteResponse, 0, len(notes))
	for _, n := range notes {
		out = append(out, dto.NoteResponse{
			ID: n.ID, Title: n.Title, Body: n.Body, CreatedBy: n.CreatedBy,
			AuthorName: n.AuthorName, CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt,
		})
	}
	return out, nil
}

func (s *Service) CreateNote(ctx context.Context, orgID, moduleID, recordID, userID string, req dto.CreateNoteRequest) (*dto.NoteResponse, error) {
	if err := s.ensureRecord(ctx, orgID, moduleID, recordID); err != nil {
		return nil, err
	}
	body := strings.TrimSpace(req.Body)
	if body == "" {
		return nil, errors.New("body required")
	}
	n, err := s.repo.CreateNote(ctx, orgID, moduleID, recordID, userID, body, req.Title)
	if err != nil {
		return nil, err
	}
	_ = s.repo.CreateActivity(ctx, orgID, moduleID, recordID, userID, "NOTE_ADDED", "Note added", nil)
	return &dto.NoteResponse{
		ID: n.ID, Title: n.Title, Body: n.Body, CreatedBy: n.CreatedBy,
		AuthorName: n.AuthorName, CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt,
	}, nil
}

func (s *Service) DeleteNote(ctx context.Context, orgID, moduleID, recordID, noteID, userID string) error {
	if err := s.ensureRecord(ctx, orgID, moduleID, recordID); err != nil {
		return err
	}
	ok, err := s.repo.DeleteNote(ctx, orgID, noteID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	_ = s.repo.CreateActivity(ctx, orgID, moduleID, recordID, userID, "NOTE_DELETED", "Note deleted", nil)
	return nil
}

func (s *Service) ListAttachments(ctx context.Context, orgID, moduleID, recordID string) ([]dto.AttachmentResponse, error) {
	if err := s.ensureRecord(ctx, orgID, moduleID, recordID); err != nil {
		return nil, err
	}
	items, err := s.repo.ListAttachments(ctx, orgID, moduleID, recordID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.AttachmentResponse, 0, len(items))
	for _, a := range items {
		out = append(out, dto.AttachmentResponse{
			ID: a.ID, FileName: a.FileName, FileURL: a.FileURL, PublicID: a.PublicID,
			ResourceType: a.ResourceType, FileSize: a.FileSize, UploadedBy: a.UploadedBy,
			UploaderName: a.UploaderName, CreatedAt: a.CreatedAt,
		})
	}
	return out, nil
}

func (s *Service) CreateAttachment(ctx context.Context, orgID, moduleID, recordID, userID string, req dto.CreateAttachmentRequest) (*dto.AttachmentResponse, error) {
	if err := s.ensureRecord(ctx, orgID, moduleID, recordID); err != nil {
		return nil, err
	}
	a, err := s.repo.CreateAttachment(
		ctx, orgID, moduleID, recordID, userID,
		req.FileName, req.FileURL, req.PublicID, req.ResourceType, req.FileSize,
	)
	if err != nil {
		return nil, err
	}
	meta, _ := json.Marshal(map[string]any{"file_name": req.FileName})
	_ = s.repo.CreateActivity(ctx, orgID, moduleID, recordID, userID, "ATTACHMENT_ADDED", "Attachment uploaded: "+req.FileName, meta)
	return &dto.AttachmentResponse{
		ID: a.ID, FileName: a.FileName, FileURL: a.FileURL, PublicID: a.PublicID,
		ResourceType: a.ResourceType, FileSize: a.FileSize, UploadedBy: a.UploadedBy,
		UploaderName: a.UploaderName, CreatedAt: a.CreatedAt,
	}, nil
}

func (s *Service) DeleteAttachment(ctx context.Context, orgID, moduleID, recordID, id, userID string) error {
	if err := s.ensureRecord(ctx, orgID, moduleID, recordID); err != nil {
		return err
	}
	ok, err := s.repo.DeleteAttachment(ctx, orgID, id)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	_ = s.repo.CreateActivity(ctx, orgID, moduleID, recordID, userID, "ATTACHMENT_DELETED", "Attachment deleted", nil)
	return nil
}

func (s *Service) ListActivities(ctx context.Context, orgID, moduleID, recordID string) ([]dto.ActivityResponse, error) {
	if err := s.ensureRecord(ctx, orgID, moduleID, recordID); err != nil {
		return nil, err
	}
	items, err := s.repo.ListActivities(ctx, orgID, moduleID, recordID, 50)
	if err != nil {
		return nil, err
	}
	out := make([]dto.ActivityResponse, 0, len(items))
	for _, a := range items {
		out = append(out, dto.ActivityResponse{
			ID: a.ID, Action: a.Action, Description: a.Description,
			PerformedBy: a.PerformedBy, ActorName: a.ActorName,
			Metadata: a.Metadata, CreatedAt: a.CreatedAt,
		})
	}
	return out, nil
}

// LogRecordActivity is used by the record service on CUD.
func (s *Service) LogRecordActivity(ctx context.Context, orgID, moduleID, recordID, userID, action, description string, metadata map[string]any) error {
	var raw json.RawMessage
	if metadata != nil {
		b, err := json.Marshal(metadata)
		if err != nil {
			return err
		}
		raw = b
	}
	return s.repo.CreateActivity(ctx, orgID, moduleID, recordID, userID, action, description, raw)
}

func (s *Service) ListRelated(ctx context.Context, orgID, moduleID string) ([]dto.RelatedDescriptorResponse, error) {
	items, err := s.repo.ListRelatedDescriptors(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.RelatedDescriptorResponse, 0, len(items))
	for _, d := range items {
		out = append(out, dto.RelatedDescriptorResponse{
			ChildModuleID: d.ChildModuleID, ChildModuleName: d.ChildModuleName,
			ChildAPIName: d.ChildAPIName, LookupFieldAPI: d.LookupFieldAPI,
			LookupFieldLabel: d.LookupFieldLabel,
		})
	}
	return out, nil
}

func (s *Service) RelatedLookupField(ctx context.Context, orgID, childModuleID, parentModuleID string) (string, error) {
	api, err := s.repo.LookupFieldAPI(ctx, orgID, childModuleID, parentModuleID)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrNotFound, err)
	}
	return api, nil
}
