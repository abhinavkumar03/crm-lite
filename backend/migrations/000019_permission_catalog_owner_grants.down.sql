-- Leave the permission catalog in place (other roles may reference it).
-- Only remove the backfilled grants for system full-access roles that were
-- added by this migration — cannot distinguish cleanly, so down is a no-op
-- for grants. Catalog rows are kept.
SELECT 1;
