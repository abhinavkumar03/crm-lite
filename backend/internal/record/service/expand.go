package service

import (
	"context"

	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record/dto"
)

// expand resolves lookup and user references on the given records into
// human-readable labels, attaching them to each response's Relations map. It
// batches lookups per target module so a page of records costs one query per
// referenced module plus one for users.
func (s *Service) expand(
	ctx context.Context,
	orgID string,
	fields []fieldentity.Field,
	records []*dto.RecordResponse,
) error {
	lookupFields, userFields := relationFields(fields)
	if len(lookupFields) == 0 && len(userFields) == 0 {
		return nil
	}

	// Resolve lookup references, grouped by target module.
	for _, f := range lookupFields {
		targetModule := *f.LookupModuleID
		ids := collectRefIDs(records, f.APIName)
		if len(ids) == 0 {
			continue
		}

		displayField, err := s.displayFieldFor(ctx, orgID, targetModule)
		if err != nil {
			return err
		}

		labels, err := s.repo.DisplayValues(ctx, orgID, targetModule, ids, displayField)
		if err != nil {
			return err
		}
		assignRelations(records, f.APIName, labels)
	}

	// Resolve user references across all user-type fields at once.
	if len(userFields) > 0 {
		var allIDs []string
		for _, f := range userFields {
			allIDs = append(allIDs, collectRefIDs(records, f.APIName)...)
		}
		labels, err := s.repo.UserDisplays(ctx, dedupe(allIDs))
		if err != nil {
			return err
		}
		for _, f := range userFields {
			assignRelations(records, f.APIName, labels)
		}
	}

	return nil
}

func relationFields(fields []fieldentity.Field) (lookups, users []fieldentity.Field) {
	for _, f := range fields {
		switch f.FieldType {
		case fieldentity.TypeLookup:
			if f.LookupModuleID != nil && *f.LookupModuleID != "" {
				lookups = append(lookups, f)
			}
		case fieldentity.TypeUser:
			users = append(users, f)
		}
	}
	return lookups, users
}

// displayFieldFor picks the best label field of a target module: a field named
// "name", else the first searchable field, else the first text field.
func (s *Service) displayFieldFor(ctx context.Context, orgID, moduleID string) (string, error) {
	fields, err := s.fields.List(ctx, orgID, moduleID)
	if err != nil {
		return "", err
	}

	var firstSearchable, firstText string
	for _, f := range fields {
		if f.APIName == "name" {
			return "name", nil
		}
		if firstSearchable == "" && f.IsSearchable {
			firstSearchable = f.APIName
		}
		if firstText == "" && f.FieldType == fieldentity.TypeText {
			firstText = f.APIName
		}
	}
	if firstSearchable != "" {
		return firstSearchable, nil
	}
	return firstText, nil
}

func collectRefIDs(records []*dto.RecordResponse, apiName string) []string {
	seen := map[string]struct{}{}
	var ids []string
	for _, r := range records {
		if id, ok := r.Data[apiName].(string); ok && id != "" {
			if _, dup := seen[id]; !dup {
				seen[id] = struct{}{}
				ids = append(ids, id)
			}
		}
	}
	return ids
}

func assignRelations(records []*dto.RecordResponse, apiName string, labels map[string]string) {
	for _, r := range records {
		id, ok := r.Data[apiName].(string)
		if !ok || id == "" {
			continue
		}
		label, found := labels[id]
		if !found {
			label = id
		}
		if r.Relations == nil {
			r.Relations = map[string]dto.RelationRef{}
		}
		r.Relations[apiName] = dto.RelationRef{ID: id, Label: label}
	}
}

func dedupe(in []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, s := range in {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}
