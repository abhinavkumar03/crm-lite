package seeders

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WorkspaceSideDataSeeder attaches notes, attachments, and timeline activities
// (including call.* actions) to existing records across every demo workspace.
type WorkspaceSideDataSeeder struct{}

func (WorkspaceSideDataSeeder) Name() string { return "workspace_side_data" }

func (WorkspaceSideDataSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	orgIDs, err := listDemoOrgIDs(ctx, db)
	if err != nil {
		return err
	}
	owners, err := demoOwnerIDs(ctx, db)
	if err != nil {
		return err
	}

	modules := []string{"lead", "contact", "task", "company", "deal"}
	now := time.Now()

	for orgIdx, orgID := range orgIDs {
		r := rand.New(rand.NewSource(99 + int64(orgIdx)*500))

		// Idempotent: skip if this org already has RECORD notes.
		var noteCount int
		if err := db.QueryRow(ctx, `
			SELECT count(*) FROM notes
			WHERE organization_id = $1 AND entity_type = 'RECORD'
		`, orgID).Scan(&noteCount); err != nil {
			return err
		}
		if noteCount > 0 {
			continue
		}

		for _, apiName := range modules {
			moduleID, err := getModuleID(ctx, db, orgID, apiName)
			if err != nil {
				return err
			}
			recordIDs, err := listRecordIDs(ctx, db, orgID, moduleID, 40)
			if err != nil {
				return err
			}
			if len(recordIDs) == 0 {
				continue
			}

			limit := 12
			if len(recordIDs) < limit {
				limit = len(recordIDs)
			}
			for i := 0; i < limit; i++ {
				recordID := recordIDs[i]
				userID := owners[i%len(owners)]
				createdAt := spread(r, now)

				if err := insertNote(ctx, db, orgID, moduleID, recordID, userID, pick(r, noteBodies), createdAt); err != nil {
					return err
				}
				if err := insertAttachment(ctx, db, orgID, moduleID, recordID, userID, r, createdAt.Add(time.Hour)); err != nil {
					return err
				}

				acts := []struct {
					Action string
					Desc   string
				}{
					{"record.created", fmt.Sprintf("%s record created in demo seed", apiName)},
					{"record.updated", "Updated key fields after discovery"},
					{"note.added", "Follow-up note added"},
					{"attachment.uploaded", "Supporting document uploaded"},
					{pick(r, []string{"call.incoming", "call.outgoing", "call.missed", "call.busy", "call.completed"}), pick(r, callSummaries)},
				}
				if apiName == "task" {
					acts = append(acts, struct{ Action, Desc string }{"task.completed", "Marked task completed after customer confirmation"})
				}
				for j, a := range acts {
					meta, _ := json.Marshal(map[string]any{
						"seed": true, "module": apiName, "index": j,
					})
					if err := insertActivity(ctx, db, orgID, moduleID, recordID, userID, a.Action, a.Desc, meta, createdAt.Add(time.Duration(j)*time.Hour)); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func insertNote(ctx context.Context, db *pgxpool.Pool, orgID, moduleID, recordID, userID, body string, at time.Time) error {
	_, err := db.Exec(ctx, `
		INSERT INTO notes (
			id, entity_type, entity_id, note, title,
			organization_id, module_id, created_by, updated_by, created_at, updated_at
		) VALUES ($1,'RECORD',$2,$3,$4,$5,$6,$7,$7,$8,$8)
	`, uuid.NewString(), recordID, body, "Follow-up", orgID, moduleID, userID, at)
	return err
}

func insertAttachment(ctx context.Context, db *pgxpool.Pool, orgID, moduleID, recordID, userID string, r *rand.Rand, at time.Time) error {
	name := pick(r, attachmentNames)
	url := pick(r, attachmentURLs)
	resource := "raw"
	if stringsHasSuffix(name, ".jpg") || stringsHasSuffix(name, ".png") {
		resource = "image"
	} else if stringsHasSuffix(name, ".mp4") {
		resource = "video"
	}
	_, err := db.Exec(ctx, `
		INSERT INTO attachments (
			id, entity_type, entity_id, file_name, file_url, public_id,
			resource_type, file_size, uploaded_by, organization_id, module_id, created_at
		) VALUES ($1,'RECORD',$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, uuid.NewString(), recordID, name, url, "demo/"+name, resource, int64(50_000+r.Intn(500_000)), userID, orgID, moduleID, at)
	return err
}

func insertActivity(
	ctx context.Context, db *pgxpool.Pool,
	orgID, moduleID, recordID, userID, action, description string,
	metadata []byte, at time.Time,
) error {
	_, err := db.Exec(ctx, `
		INSERT INTO activities (
			id, entity_type, entity_id, action, description, performed_by,
			metadata, organization_id, module_id, created_at
		) VALUES ($1,'RECORD',$2,$3,$4,$5,$6,$7,$8,$9)
	`, uuid.NewString(), recordID, action, description, userID, metadata, orgID, moduleID, at)
	return err
}

func stringsHasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}
