ALTER TABLE layouts DROP CONSTRAINT IF EXISTS layouts_layout_type_check;
ALTER TABLE layouts ADD CONSTRAINT layouts_layout_type_check CHECK (layout_type IN ('form', 'detail'));

DELETE FROM layouts WHERE layout_type = 'list';

ALTER TABLE fields DROP COLUMN IF EXISTS viewable_by;
ALTER TABLE fields DROP COLUMN IF EXISTS editable_by;
ALTER TABLE fields DROP COLUMN IF EXISTS lock_mode;

-- Revert any rows using new types before restoring the narrower CHECK.
UPDATE fields SET field_type = 'number' WHERE field_type = 'percentage';
UPDATE fields SET field_type = 'text' WHERE field_type IN ('time', 'gst', 'pan', 'address', 'auto_number', 'barcode', 'serial_number');
UPDATE fields SET field_type = 'boolean' WHERE field_type = 'toggle';

ALTER TABLE fields DROP CONSTRAINT IF EXISTS fields_field_type_check;
ALTER TABLE fields ADD CONSTRAINT fields_field_type_check CHECK (field_type IN (
    'text','textarea','email','phone','number','currency',
    'date','datetime','boolean','dropdown','multiselect',
    'radio','checkbox','url','file','image','user','lookup',
    'formula','json','richtext'
));
