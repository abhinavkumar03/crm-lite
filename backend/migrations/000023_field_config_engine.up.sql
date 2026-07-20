-- Dynamic Field Configuration Engine: lock modes, ACL stubs, field types, list layouts.

-- Expand field_type CHECK to include Zoho-style engine types.
ALTER TABLE fields DROP CONSTRAINT IF EXISTS fields_field_type_check;
ALTER TABLE fields ADD CONSTRAINT fields_field_type_check CHECK (field_type IN (
    'text','textarea','email','phone','number','currency','percentage',
    'date','datetime','time','boolean','toggle','dropdown','multiselect',
    'radio','checkbox','url','file','image','user','lookup','formula','json',
    'richtext','gst','pan','address','auto_number','barcode','serial_number'
));

ALTER TABLE fields
    ADD COLUMN IF NOT EXISTS lock_mode VARCHAR(20) NOT NULL DEFAULT 'never'
        CHECK (lock_mode IN ('never', 'after_create', 'always')),
    ADD COLUMN IF NOT EXISTS editable_by VARCHAR(40) NOT NULL DEFAULT 'ALL',
    ADD COLUMN IF NOT EXISTS viewable_by VARCHAR(40) NOT NULL DEFAULT 'ALL';

-- Allow list layouts alongside form/detail.
ALTER TABLE layouts DROP CONSTRAINT IF EXISTS layouts_layout_type_check;
ALTER TABLE layouts ADD CONSTRAINT layouts_layout_type_check CHECK (layout_type IN ('form', 'detail', 'list'));
