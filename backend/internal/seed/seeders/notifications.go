package seeders

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NotificationDemoSeeder seeds templates and realistic delivery history for demo orgs.
type NotificationDemoSeeder struct{}

func (NotificationDemoSeeder) Name() string { return "notification_demo" }

func (NotificationDemoSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	orgIDs, err := listDemoOrgIDs(ctx, db)
	if err != nil {
		return err
	}
	for _, orgID := range orgIDs {
		if err := seedOrgNotifications(ctx, db, orgID); err != nil {
			return err
		}
	}
	return nil
}

func seedOrgNotifications(ctx context.Context, db *pgxpool.Pool, orgID string) error {
	var userID *string
	_ = db.QueryRow(ctx, `
		SELECT u.id FROM users u
		JOIN organization_members om ON om.user_id = u.id
		WHERE om.organization_id = $1 AND om.status = 'active'
		ORDER BY u.created_at
		LIMIT 1
	`, orgID).Scan(&userID)

	templates := []struct {
		channel, name, category, subject, body string
	}{
		{"email", "Welcome Email", "welcome", "Welcome to {{workspace.name}}", "Hello {{lead.name}},\n\nThank you for your interest in {{workspace.name}}."},
		{"email", "Follow-up Email", "follow_up", "Following up — {{lead.company}}", "Hi {{lead.name}},\n\nJust checking in on our last conversation."},
		{"email", "Proposal Email", "proposal", "Proposal for {{lead.company}}", "Hi {{lead.name}},\n\nPlease find our proposal attached."},
		{"whatsapp", "Welcome WhatsApp", "welcome", "", "Hello {{lead.name}}, thanks for connecting with {{workspace.name}}!"},
		{"whatsapp", "Reminder WhatsApp", "reminder", "", "Hi {{lead.name}}, this is a friendly reminder from {{owner.name}}."},
		{"email", "Invoice Email", "invoice", "Invoice from {{workspace.name}}", "Hello {{contact.name}},\n\nYour invoice is ready."},
	}

	templateIDs := map[string]string{}
	for _, t := range templates {
		var id string
		err := db.QueryRow(ctx, `
			INSERT INTO notification_templates (
				organization_id, channel, name, category, subject, body, variables, is_active, created_by
			) VALUES ($1,$2,$3,$4,$5,$6,'["lead.name","workspace.name","owner.name","contact.name"]'::jsonb, TRUE, $7)
			ON CONFLICT (organization_id, channel, name) DO UPDATE
			SET body = EXCLUDED.body, subject = EXCLUDED.subject, updated_at = NOW()
			RETURNING id
		`, orgID, t.channel, t.name, t.category, nullIfEmpty(t.subject), t.body, userID).Scan(&id)
		if err != nil {
			return err
		}
		templateIDs[t.channel+":"+t.name] = id
	}

	var moduleID, recordID *string
	_ = db.QueryRow(ctx, `
		SELECT m.id, r.id FROM modules m
		JOIN records r ON r.module_id = m.id AND r.organization_id = m.organization_id
		WHERE m.organization_id = $1 AND m.api_name = 'lead'
		ORDER BY r.created_at LIMIT 1
	`, orgID).Scan(&moduleID, &recordID)

	now := time.Now().UTC()
	samples := []struct {
		channel, to, subject, body, status string
		hoursAgo                           int
		fail                               bool
		scheduled                          bool
		draft                              bool
	}{
		{"email", "dana@example.com", "Welcome to CRM Lite", "Hello Dana, welcome aboard.", "delivered", 2, false, false, false},
		{"email", "sam@example.com", "Follow-up", "Hi Sam, following up.", "opened", 26, false, false, false},
		{"whatsapp", "+15551234567", "", "Hi John, can we schedule a demo?", "delivered", 1, false, false, false},
		{"whatsapp", "+15557654321", "", "Reminder: call tomorrow 11 AM", "sent", 5, false, false, false},
		{"email", "fail@example.com", "Proposal Shared", "Please review the proposal.", "failed", 3, true, false, false},
		{"email", "later@example.com", "Scheduled check-in", "Checking in next week.", "scheduled", -48, false, true, false},
		{"whatsapp", "+15550001111", "", "Draft message for later", "draft", 0, false, false, true},
		{"email", "retry@example.com", "Retry sample", "This one is retrying.", "retrying", 4, true, false, false},
	}

	for _, s := range samples {
		created := now.Add(-time.Duration(s.hoursAgo) * time.Hour)
		var scheduledAt *time.Time
		var queuedAt, sentAt, deliveredAt, openedAt *time.Time
		status := s.status
		var errMsg *string
		if s.fail {
			msg := "simulated provider timeout"
			errMsg = &msg
		}
		if s.scheduled {
			t := now.Add(48 * time.Hour)
			scheduledAt = &t
		}
		switch status {
		case "queued", "processing", "sent", "delivered", "opened", "read", "failed", "retrying":
			queuedAt = &created
		}
		switch status {
		case "sent", "delivered", "opened", "read":
			sentAt = &created
		}
		switch status {
		case "delivered", "opened", "read":
			deliveredAt = &created
		}
		switch status {
		case "opened", "read":
			openedAt = &created
		}

		var provider *string
		if status == "sent" || status == "delivered" || status == "opened" || status == "read" {
			p := "simulation"
			provider = &p
		}
		retryCount := 0
		if status == "retrying" {
			retryCount = 1
		}
		var entityType *string
		if moduleID != nil && recordID != nil {
			et := "RECORD"
			entityType = &et
		}

		vars, _ := json.Marshal(map[string]any{"name": "Dana", "workspace.name": "Demo Org"})
		var id string
		err := db.QueryRow(ctx, `
			INSERT INTO notifications (
				organization_id, channel, recipient, subject, body, template, data, variables_used,
				status, provider, error, last_error, entity_type, entity_id, module_id, created_by,
				scheduled_at, queued_at, sent_at, delivered_at, opened_at, retry_count, created_at, updated_at
			) VALUES (
				$1,$2,$3,$4,$5,$6,'{}'::jsonb,$7,
				$8,$9,$10,$10,$11,$12,$13,$14,
				$15,$16,$17,$18,$19,$20,$21,$21
			)
			RETURNING id
		`,
			orgID, s.channel, s.to, nullIfEmpty(s.subject), s.body, nullIfEmpty("demo"),
			vars, status, provider, errMsg, entityType, recordID, moduleID, userID,
			scheduledAt, queuedAt, sentAt, deliveredAt, openedAt, retryCount, created,
		).Scan(&id)
		if err != nil {
			return err
		}

		events := []string{"queued"}
		switch status {
		case "delivered", "opened", "read":
			events = []string{"queued", "processing", "sent", "delivered"}
			if status == "opened" || status == "read" {
				events = append(events, "opened")
			}
		case "sent":
			events = []string{"queued", "processing", "sent"}
		case "failed", "retrying":
			events = []string{"queued", "processing", "failed"}
		case "scheduled":
			events = []string{"scheduled"}
		case "draft":
			events = []string{"draft"}
		}
		for _, ev := range events {
			var evProvider *string
			if ev == "sent" || ev == "delivered" || ev == "opened" {
				p := "simulation"
				evProvider = &p
			}
			_, _ = db.Exec(ctx, `
				INSERT INTO notification_delivery_events (organization_id, notification_id, event, provider, payload)
				VALUES ($1,$2,$3,$4,'{}'::jsonb)
			`, orgID, id, ev, evProvider)
		}

		if moduleID != nil && recordID != nil && userID != nil && (status == "delivered" || status == "sent" || status == "failed") {
			action := "EMAIL_SENT"
			if s.channel == "whatsapp" {
				action = "WHATSAPP_SENT"
			}
			if status == "failed" {
				action = "NOTIFICATION_FAILED"
			}
			_, _ = db.Exec(ctx, `
				INSERT INTO activities (
					id, entity_type, entity_id, action, description, performed_by,
					metadata, organization_id, module_id, created_at
				) VALUES (
					gen_random_uuid(), 'RECORD', $1, $2, $3, $4,
					jsonb_build_object('notification_id', $5::text, 'channel', $6),
					$7, $8, $9
				)
			`, *recordID, action, "Demo "+action+" to "+s.to, *userID, id, s.channel, orgID, *moduleID, created)
		}
	}

	_ = templateIDs
	return nil
}

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
