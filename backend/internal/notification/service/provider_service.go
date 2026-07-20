package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notify"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/secrets"
)

var (
	ErrProviderNotFound = errors.New("communication provider not found")
	ErrInvalidProvider  = errors.New("invalid provider configuration")
)

type ProviderService struct {
	repo   *repository.ProviderRepository
	box    *secrets.Box
	logger *zap.Logger
}

func NewProviderService(repo *repository.ProviderRepository, box *secrets.Box, logger *zap.Logger) *ProviderService {
	return &ProviderService{repo: repo, box: box, logger: logger}
}

func (s *ProviderService) List(ctx context.Context, orgID, channel string) ([]dto.ProviderResponse, error) {
	items, err := s.repo.List(ctx, orgID, channel)
	if err != nil {
		return nil, err
	}
	out := make([]dto.ProviderResponse, 0, len(items))
	for i := range items {
		out = append(out, toProviderResponse(&items[i]))
	}
	return out, nil
}

func (s *ProviderService) Get(ctx context.Context, orgID, id string) (*dto.ProviderResponse, error) {
	p, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrProviderNotFound
	}
	resp := toProviderResponse(p)
	return &resp, nil
}

func (s *ProviderService) Create(ctx context.Context, orgID, userID string, req dto.ProviderUpsertRequest) (*dto.ProviderResponse, error) {
	if err := validateProviderRequest(req); err != nil {
		return nil, err
	}
	cfgBytes, _ := json.Marshal(req.Config)
	if cfgBytes == nil {
		cfgBytes = []byte("{}")
	}
	var secretsEnc []byte
	if len(req.Secrets) > 0 {
		enc, err := s.box.EncryptJSON(req.Secrets)
		if err != nil {
			return nil, err
		}
		secretsEnc = enc
	}
	isDefault := req.IsDefault != nil && *req.IsDefault
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	if isDefault {
		_ = s.repo.ClearDefault(ctx, orgID, req.Channel)
	}
	p := &entity.Provider{
		OrganizationID:   orgID,
		Channel:          req.Channel,
		ProviderType:     req.ProviderType,
		Name:             req.Name,
		Config:           cfgBytes,
		SecretsEncrypted: secretsEnc,
		IsDefault:        isDefault,
		IsActive:         isActive,
		CreatedBy:        &userID,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	resp := toProviderResponse(p)
	return &resp, nil
}

func (s *ProviderService) Update(ctx context.Context, orgID, id string, req dto.ProviderUpsertRequest) (*dto.ProviderResponse, error) {
	existing, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrProviderNotFound
	}
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.ProviderType != "" {
		existing.ProviderType = req.ProviderType
	}
	if req.Config != nil {
		cfgBytes, _ := json.Marshal(req.Config)
		existing.Config = cfgBytes
	}
	if len(req.Secrets) > 0 {
		enc, err := s.box.EncryptJSON(req.Secrets)
		if err != nil {
			return nil, err
		}
		existing.SecretsEncrypted = enc
	}
	if req.IsDefault != nil {
		if *req.IsDefault {
			_ = s.repo.ClearDefault(ctx, orgID, existing.Channel)
		}
		existing.IsDefault = *req.IsDefault
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	resp := toProviderResponse(existing)
	return &resp, nil
}

func (s *ProviderService) Delete(ctx context.Context, orgID, id string) error {
	return s.repo.Delete(ctx, orgID, id)
}

func (s *ProviderService) Test(ctx context.Context, orgID, id string, req dto.ProviderTestRequest) error {
	p, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		return err
	}
	if p == nil {
		return ErrProviderNotFound
	}
	provider, err := s.BuildNotifyProvider(p)
	if err != nil {
		msg := err.Error()
		_ = s.repo.MarkHealth(ctx, id, &msg)
		return err
	}
	subject := "CRM Lite provider test"
	body := "This is a test message from CRM Lite communication providers."
	_, sendErr := provider.Send(ctx, notify.Message{
		Channel:  notify.Channel(p.Channel),
		To:       req.To,
		Subject:  subject,
		Body:     body,
		HTMLBody: "<p>" + body + "</p>",
	})
	if sendErr != nil {
		msg := sendErr.Error()
		_ = s.repo.MarkHealth(ctx, id, &msg)
		return sendErr
	}
	_ = s.repo.MarkHealth(ctx, id, nil)
	return nil
}

// BuildNotifyProvider constructs a live notify.Provider from a stored row.
func (s *ProviderService) BuildNotifyProvider(p *entity.Provider) (notify.Provider, error) {
	cfgMap := map[string]any{}
	_ = json.Unmarshal(p.Config, &cfgMap)
	secMap := map[string]any{}
	if len(p.SecretsEncrypted) > 0 {
		if err := s.box.DecryptJSON(p.SecretsEncrypted, &secMap); err != nil {
			return nil, fmt.Errorf("decrypt secrets: %w", err)
		}
	}

	switch p.Channel {
	case entity.ChannelEmail:
		emailCfg := notify.EmailConfig{
			Provider:     p.ProviderType,
			SMTPHost:     str(cfgMap, "smtp_host", "host"),
			SMTPPort:     intVal(cfgMap, "smtp_port", "port"),
			SMTPUsername: str(secMap, "smtp_username", "username"),
			SMTPPassword: str(secMap, "smtp_password", "password"),
			SMTPFrom:     str(cfgMap, "from", "default_sender"),
			Encryption:   str(cfgMap, "encryption"),
			APIKey:       str(secMap, "api_key", "apiKey"),
			APIURL:       str(cfgMap, "api_url"),
			From:         str(cfgMap, "from", "default_sender"),
			ReplyTo:      str(cfgMap, "reply_to"),
		}
		if emailCfg.SMTPPort == 0 {
			emailCfg.SMTPPort = 587
		}
		return notify.BuildEmailProvider(emailCfg, s.logger), nil
	case entity.ChannelWhatsApp:
		waCfg := notify.WhatsAppConfig{
			Provider:   p.ProviderType,
			APIURL:     str(cfgMap, "api_url", "apiURL"),
			Token:      str(secMap, "access_token", "token", "api_key"),
			PhoneID:    str(cfgMap, "phone_number_id", "phone_id"),
			AccountSID: str(secMap, "account_sid", "accountSID"),
			AuthToken:  str(secMap, "auth_token", "authToken"),
			FromNumber: str(cfgMap, "from_number", "from"),
		}
		if waCfg.APIURL == "" {
			waCfg.APIURL = "https://graph.facebook.com/v20.0"
		}
		return notify.BuildWhatsAppProvider(waCfg, s.logger), nil
	default:
		return nil, fmt.Errorf("%w: unknown channel %s", ErrInvalidProvider, p.Channel)
	}
}

// Resolve implements notify.OrgProviderResolver.
func (s *ProviderService) Resolve(ctx context.Context, orgID string, channel notify.Channel) (notify.Provider, error) {
	p, err := s.repo.GetDefault(ctx, orgID, string(channel))
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	return s.BuildNotifyProvider(p)
}

func (s *ProviderService) ListSenders(ctx context.Context, orgID, channel string) ([]dto.SenderResponse, error) {
	items, err := s.repo.ListSenders(ctx, orgID, channel)
	if err != nil {
		return nil, err
	}
	out := make([]dto.SenderResponse, 0, len(items))
	for i := range items {
		out = append(out, dto.SenderResponse{
			ID: items[i].ID, ProviderID: items[i].ProviderID, Channel: items[i].Channel,
			DisplayName: items[i].DisplayName, FromAddress: items[i].FromAddress,
			ReplyTo: items[i].ReplyTo, IsDefault: items[i].IsDefault,
			CreatedAt: items[i].CreatedAt, UpdatedAt: items[i].UpdatedAt,
		})
	}
	return out, nil
}

func (s *ProviderService) CreateSender(ctx context.Context, orgID string, req dto.SenderUpsertRequest) (*dto.SenderResponse, error) {
	isDefault := req.IsDefault != nil && *req.IsDefault
	row := &entity.SenderIdentity{
		OrganizationID: orgID,
		ProviderID:     req.ProviderID,
		Channel:        req.Channel,
		DisplayName:    req.DisplayName,
		FromAddress:    req.FromAddress,
		ReplyTo:        req.ReplyTo,
		IsDefault:      isDefault,
	}
	if err := s.repo.CreateSender(ctx, row); err != nil {
		return nil, err
	}
	resp := dto.SenderResponse{
		ID: row.ID, ProviderID: row.ProviderID, Channel: row.Channel,
		DisplayName: row.DisplayName, FromAddress: row.FromAddress, ReplyTo: row.ReplyTo,
		IsDefault: row.IsDefault, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}
	return &resp, nil
}

func (s *ProviderService) DeleteSender(ctx context.Context, orgID, id string) error {
	return s.repo.DeleteSender(ctx, orgID, id)
}

func validateProviderRequest(req dto.ProviderUpsertRequest) error {
	ch := strings.ToLower(req.Channel)
	if ch != "email" && ch != "whatsapp" {
		return fmt.Errorf("%w: channel", ErrInvalidProvider)
	}
	if strings.TrimSpace(req.ProviderType) == "" || strings.TrimSpace(req.Name) == "" {
		return ErrInvalidProvider
	}
	return nil
}

func toProviderResponse(p *entity.Provider) dto.ProviderResponse {
	cfg := map[string]any{}
	_ = json.Unmarshal(p.Config, &cfg)
	return dto.ProviderResponse{
		ID: p.ID, Channel: p.Channel, ProviderType: p.ProviderType, Name: p.Name,
		Config: cfg, SecretsConfigured: p.SecretsConfigured || len(p.SecretsEncrypted) > 0,
		IsDefault: p.IsDefault, IsActive: p.IsActive, LastHealthAt: p.LastHealthAt,
		LastError: p.LastError, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func str(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			switch t := v.(type) {
			case string:
				if t != "" {
					return t
				}
			}
		}
	}
	return ""
}

func intVal(m map[string]any, keys ...string) int {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			switch t := v.(type) {
			case float64:
				return int(t)
			case int:
				return t
			case json.Number:
				i, _ := t.Int64()
				return int(i)
			}
		}
	}
	return 0
}
