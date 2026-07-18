-- Phase 17: composite indexes for hot list paths.
-- records: org + module + created_at covers default list/export paging.
-- notifications: org + created_at covers inbox / history listing.

CREATE INDEX IF NOT EXISTS idx_records_org_module_created
  ON records (organization_id, module_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_notifications_org_created
  ON notifications (organization_id, created_at DESC);
