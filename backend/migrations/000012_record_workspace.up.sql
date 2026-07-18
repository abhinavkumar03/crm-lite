-- Record Workspace: polymorphic side-entities for dynamic records + layout seeds.

-- =====================================================================
-- Notes / attachments / activities — org + module scoping for RECORD
-- =====================================================================

ALTER TABLE notes
    ADD COLUMN IF NOT EXISTS organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS module_id UUID REFERENCES modules(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS title VARCHAR(200);

ALTER TABLE attachments
    ADD COLUMN IF NOT EXISTS organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS module_id UUID REFERENCES modules(id) ON DELETE CASCADE;

ALTER TABLE activities
    ADD COLUMN IF NOT EXISTS organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS module_id UUID REFERENCES modules(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_notes_org_module_entity
    ON notes(organization_id, module_id, entity_id)
    WHERE organization_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_attachments_org_module_entity
    ON attachments(organization_id, module_id, entity_id)
    WHERE organization_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_activities_org_module_entity
    ON activities(organization_id, module_id, entity_id)
    WHERE organization_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_notes_entity
    ON notes(entity_type, entity_id);

CREATE INDEX IF NOT EXISTS idx_attachments_entity
    ON attachments(entity_type, entity_id);

CREATE INDEX IF NOT EXISTS idx_activities_entity
    ON activities(entity_type, entity_id);

-- =====================================================================
-- Default detail layouts for every existing dynamic module
-- =====================================================================

INSERT INTO layouts (organization_id, module_id, name, layout_type, is_default, config)
SELECT
    m.organization_id,
    m.id,
    'Default Detail',
    'detail',
    TRUE,
    jsonb_build_object(
        'sections', jsonb_build_array(
            jsonb_build_object(
                'key', 'general',
                'label', 'General Information',
                'fields', COALESCE((
                    SELECT jsonb_agg(f.api_name ORDER BY f.sort_order, f.api_name)
                    FROM fields f
                    WHERE f.module_id = m.id
                      AND f.is_visible = TRUE
                      AND COALESCE(f.is_system, FALSE) = FALSE
                ), '[]'::jsonb)
            ),
            jsonb_build_object(
                'key', 'system',
                'label', 'System Fields',
                'fields', jsonb_build_array('owner_id', 'assigned_to', 'visibility', 'created_at', 'updated_at')
            )
        ),
        'tabs', jsonb_build_array('overview', 'notes', 'attachments', 'timeline', 'related')
    )
FROM modules m
WHERE m.storage_strategy = 'dynamic'
  AND NOT EXISTS (
      SELECT 1 FROM layouts l
      WHERE l.module_id = m.id AND l.layout_type = 'detail' AND l.is_default = TRUE
  );
