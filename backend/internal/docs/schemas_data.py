"""OpenAPI 3 component schemas for CRM Lite. Consumed by generate_openapi.py."""

# UUID / datetime helpers reused via copy for clarity in property dicts.
_UUID = {"type": "string", "format": "uuid", "example": "11111111-1111-1111-1111-111111111111"}
_DT = {"type": "string", "format": "date-time", "example": "2026-07-18T10:00:00Z"}
_EMAIL = {"type": "string", "format": "email", "example": "ada@example.com"}

SCHEMAS = {
    # ── Shared envelope / errors ──────────────────────────────────────────
    "APIResponse": {
        "type": "object",
        "required": ["success"],
        "properties": {
            "success": {"type": "boolean", "example": True},
            "message": {"type": "string", "example": "OK"},
            "data": {"description": "Endpoint-specific payload", "nullable": True},
            "errors": {
                "nullable": True,
                "description": "Field errors or app error items when success is false",
                "oneOf": [
                    {"type": "array", "items": {"$ref": "#/components/schemas/FieldValidationError"}},
                    {"type": "array", "items": {"$ref": "#/components/schemas/AppErrorItem"}},
                ],
            },
        },
    },
    "ErrorResponse": {
        "type": "object",
        "required": ["success", "message"],
        "properties": {
            "success": {"type": "boolean", "example": False},
            "message": {"type": "string", "example": "Validation failed"},
            "data": {"nullable": True},
            "errors": {
                "nullable": True,
                "oneOf": [
                    {"type": "array", "items": {"$ref": "#/components/schemas/FieldValidationError"}},
                    {"type": "array", "items": {"$ref": "#/components/schemas/AppErrorItem"}},
                ],
            },
        },
    },
    "ErrorCode": {
        "type": "string",
        "enum": [
            "BAD_REQUEST",
            "UNAUTHORIZED",
            "FORBIDDEN",
            "NOT_FOUND",
            "CONFLICT",
            "VALIDATION_ERROR",
            "INTERNAL_SERVER_ERROR",
        ],
        "example": "BAD_REQUEST",
    },
    "FieldValidationError": {
        "type": "object",
        "required": ["field", "message"],
        "properties": {
            "field": {"type": "string", "example": "email"},
            "message": {"type": "string", "example": "Must be a valid email"},
        },
    },
    "AppErrorItem": {
        "type": "object",
        "required": ["code"],
        "properties": {
            "code": {"$ref": "#/components/schemas/ErrorCode"},
        },
    },
    # ── Auth ──────────────────────────────────────────────────────────────
    "RegisterRequest": {
        "type": "object",
        "required": ["name", "email", "password"],
        "properties": {
            "name": {"type": "string", "example": "Ada Lovelace"},
            "email": _EMAIL,
            "password": {"type": "string", "format": "password", "minLength": 8, "example": "Secret@12345"},
        },
    },
    "LoginRequest": {
        "type": "object",
        "required": ["email", "password"],
        "properties": {
            "email": _EMAIL,
            "password": {"type": "string", "format": "password", "example": "Admin@12345"},
        },
    },
    "UserResponse": {
        "type": "object",
        "required": ["id", "name", "email"],
        "properties": {
            "id": _UUID,
            "name": {"type": "string", "example": "Ada Lovelace"},
            "email": _EMAIL,
        },
    },
    "LoginResponse": {
        "type": "object",
        "required": ["access_token", "user"],
        "properties": {
            "access_token": {
                "type": "string",
                "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            },
            "user": {"$ref": "#/components/schemas/UserResponse"},
        },
    },
    # ── Leads ─────────────────────────────────────────────────────────────
    "CreateLeadRequest": {
        "type": "object",
        "required": ["name"],
        "properties": {
            "name": {"type": "string", "example": "Acme Corp"},
            "email": {"type": "string", "format": "email", "example": "hello@acme.com"},
            "phone": {"type": "string", "example": "+91-9876543210"},
            "company": {"type": "string", "example": "Acme"},
            "notes": {"type": "string", "example": "Inbound demo request"},
        },
    },
    "UpdateLeadRequest": {
        "type": "object",
        "properties": {
            "name": {"type": "string"},
            "email": {"type": "string", "format": "email"},
            "phone": {"type": "string"},
            "company": {"type": "string"},
            "notes": {"type": "string"},
            "status": {
                "type": "string",
                "enum": ["NEW", "CONTACTED", "QUALIFIED", "WON", "LOST"],
                "example": "CONTACTED",
            },
        },
    },
    "LeadResponse": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "name": {"type": "string"},
            "email": {"type": "string"},
            "phone": {"type": "string"},
            "company": {"type": "string"},
            "status": {
                "type": "string",
                "enum": ["NEW", "CONTACTED", "QUALIFIED", "WON", "LOST"],
            },
            "notes": {"type": "string"},
            "owner_id": {"type": "string", "format": "uuid"},
            "created_at": _DT,
            "updated_at": _DT,
        },
    },
    "ListLeadResponse": {
        "type": "object",
        "properties": {
            "data": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/LeadResponse"},
            },
            "page": {"type": "integer", "example": 1},
            "limit": {"type": "integer", "example": 20},
            "total": {"type": "integer", "format": "int64", "example": 100},
            "total_pages": {"type": "integer", "example": 5},
        },
    },
    # ── Contacts ──────────────────────────────────────────────────────────
    "CreateContactRequest": {
        "type": "object",
        "required": ["first_name"],
        "properties": {
            "first_name": {"type": "string", "example": "Ada"},
            "last_name": {"type": "string", "example": "Lovelace"},
            "email": _EMAIL,
            "phone": {"type": "string", "example": "+91-9876543210"},
            "company": {"type": "string", "example": "Analytical Engines"},
            "job_title": {"type": "string", "example": "Mathematician"},
            "notes": {"type": "string"},
        },
    },
    "UpdateContactRequest": {"$ref": "#/components/schemas/CreateContactRequest"},
    "ContactResponse": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "first_name": {"type": "string"},
            "last_name": {"type": "string"},
            "email": {"type": "string"},
            "phone": {"type": "string"},
            "company": {"type": "string"},
            "job_title": {"type": "string"},
            "notes": {"type": "string"},
            "owner_id": {"type": "string", "format": "uuid"},
            "created_at": _DT,
            "updated_at": _DT,
        },
    },
    "ListContactResponse": {
        "type": "object",
        "properties": {
            "data": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/ContactResponse"},
            },
            "page": {"type": "integer"},
            "limit": {"type": "integer"},
            "total": {"type": "integer", "format": "int64"},
            "total_pages": {"type": "integer"},
        },
    },
    # ── Tasks ─────────────────────────────────────────────────────────────
    "CreateTaskRequest": {
        "type": "object",
        "required": ["title"],
        "properties": {
            "title": {"type": "string", "example": "Call prospect"},
            "description": {"type": "string"},
            "lead_id": {"type": "string", "format": "uuid", "nullable": True},
            "contact_id": {"type": "string", "format": "uuid", "nullable": True},
            "due_date": {"type": "string", "format": "date-time", "nullable": True},
        },
    },
    "UpdateTaskRequest": {
        "type": "object",
        "properties": {
            "title": {"type": "string"},
            "description": {"type": "string"},
            "lead_id": {"type": "string", "format": "uuid", "nullable": True},
            "contact_id": {"type": "string", "format": "uuid", "nullable": True},
            "due_date": {"type": "string", "format": "date-time", "nullable": True},
            "status": {
                "type": "string",
                "enum": ["PENDING", "IN_PROGRESS", "COMPLETED"],
                "example": "IN_PROGRESS",
            },
        },
    },
    "TaskResponse": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "title": {"type": "string"},
            "description": {"type": "string"},
            "status": {
                "type": "string",
                "enum": ["PENDING", "IN_PROGRESS", "COMPLETED"],
            },
            "lead_id": {"type": "string", "format": "uuid", "nullable": True},
            "contact_id": {"type": "string", "format": "uuid", "nullable": True},
            "due_date": {"type": "string", "format": "date-time", "nullable": True},
            "owner_id": {"type": "string", "format": "uuid"},
            "created_at": _DT,
            "updated_at": _DT,
        },
    },
    "ListTaskResponse": {
        "type": "object",
        "properties": {
            "data": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/TaskResponse"},
            },
            "page": {"type": "integer"},
            "limit": {"type": "integer"},
            "total": {"type": "integer", "format": "int64"},
            "total_pages": {"type": "integer"},
        },
    },
    # ── Dashboard / Search (dynamic modules) ──────────────────────────────
    "ModuleCount": {
        "type": "object",
        "properties": {
            "module_id": _UUID,
            "api_name": {"type": "string"},
            "plural_label": {"type": "string"},
            "icon": {"type": "string"},
            "color": {"type": "string"},
            "record_count": {"type": "integer", "format": "int64"},
        },
    },
    "RecentRecord": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "module_id": _UUID,
            "module_label": {"type": "string"},
            "api_name": {"type": "string"},
            "title": {"type": "string"},
            "created_at": _DT,
        },
    },
    "DashboardResponse": {
        "type": "object",
        "properties": {
            "total_modules": {"type": "integer", "format": "int64"},
            "total_records": {"type": "integer", "format": "int64"},
            "module_counts": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/ModuleCount"},
            },
            "recent_records": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/RecentRecord"},
            },
        },
    },
    "SearchHit": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "module_id": _UUID,
            "module_label": {"type": "string"},
            "api_name": {"type": "string"},
            "title": {"type": "string"},
            "subtitle": {"type": "string"},
        },
    },
    "SearchResponse": {
        "type": "object",
        "properties": {
            "results": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/SearchHit"},
            },
        },
    },
    # ── Notes ─────────────────────────────────────────────────────────────
    "CreateNoteBody": {
        "type": "object",
        "required": ["note"],
        "properties": {
            "note": {"type": "string", "example": "Followed up by phone"},
        },
    },
    "UpdateNoteRequest": {"$ref": "#/components/schemas/CreateNoteBody"},
    "NoteResponse": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "entity_type": {
                "type": "string",
                "enum": ["LEAD", "CONTACT", "TASK"],
            },
            "entity_id": {"type": "string", "format": "uuid"},
            "note": {"type": "string"},
            "created_by": {"type": "string", "format": "uuid"},
            "updated_by": {"type": "string", "format": "uuid", "nullable": True},
            "created_at": _DT,
            "updated_at": _DT,
            "user": {
                "type": "object",
                "properties": {
                    "id": {"type": "string", "format": "uuid"},
                    "name": {"type": "string"},
                },
            },
        },
    },
    "NoteList": {
        "type": "array",
        "items": {"$ref": "#/components/schemas/NoteResponse"},
    },
    # ── Call logs ─────────────────────────────────────────────────────────
    "CreateCallLogRequest": {
        "type": "object",
        "required": ["direction", "status"],
        "properties": {
            "direction": {
                "type": "string",
                "enum": ["INCOMING", "OUTGOING"],
                "example": "OUTGOING",
            },
            "status": {
                "type": "string",
                "enum": [
                    "COMPLETED",
                    "MISSED",
                    "NO_ANSWER",
                    "BUSY",
                    "VOICEMAIL",
                    "CANCELLED",
                ],
                "example": "COMPLETED",
            },
            "summary": {"type": "string", "example": "Discussed pricing"},
            "duration_seconds": {"type": "integer", "example": 180},
            "follow_up_at": {"type": "string", "format": "date-time", "nullable": True},
        },
    },
    "UpdateCallLogRequest": {"$ref": "#/components/schemas/CreateCallLogRequest"},
    "CallLogResponse": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "entity_type": {
                "type": "string",
                "enum": ["LEAD", "CONTACT", "TASK"],
            },
            "entity_id": {"type": "string", "format": "uuid"},
            "direction": {"type": "string", "enum": ["INCOMING", "OUTGOING"]},
            "status": {
                "type": "string",
                "enum": [
                    "COMPLETED",
                    "MISSED",
                    "NO_ANSWER",
                    "BUSY",
                    "VOICEMAIL",
                    "CANCELLED",
                ],
            },
            "summary": {"type": "string"},
            "duration_seconds": {"type": "integer"},
            "follow_up_at": {"type": "string", "format": "date-time", "nullable": True},
            "created_by": {"type": "string", "format": "uuid"},
            "updated_by": {"type": "string", "format": "uuid", "nullable": True},
            "created_at": _DT,
            "updated_at": _DT,
        },
    },
    "CallLogList": {
        "type": "array",
        "items": {"$ref": "#/components/schemas/CallLogResponse"},
    },
    # ── Attachments / Activities / Media ───────────────────────────────────
    "CreateAttachmentRequest": {
        "type": "object",
        "required": ["file_name", "file_url"],
        "properties": {
            "file_name": {"type": "string", "example": "proposal.pdf"},
            "file_url": {
                "type": "string",
                "format": "uri",
                "example": "https://cdn.example.com/proposal.pdf",
            },
            "public_id": {"type": "string", "example": "crm/uploads/proposal"},
            "resource_type": {"type": "string", "example": "raw"},
            "file_size": {"type": "integer", "format": "int64", "example": 102400},
        },
    },
    "AttachmentResponse": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "entity_type": {
                "type": "string",
                "enum": ["LEAD", "CONTACT", "TASK"],
            },
            "entity_id": {"type": "string", "format": "uuid"},
            "file_name": {"type": "string"},
            "file_url": {"type": "string"},
            "public_id": {"type": "string"},
            "resource_type": {"type": "string"},
            "file_size": {"type": "integer", "format": "int64"},
            "uploaded_by": {"type": "string", "format": "uuid"},
            "created_at": _DT,
        },
    },
    "AttachmentList": {
        "type": "array",
        "items": {"$ref": "#/components/schemas/AttachmentResponse"},
    },
    "ActivityResponse": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "action": {"type": "string", "example": "LEAD_CREATED"},
            "description": {"type": "string"},
            "performed_by": {"type": "string", "format": "uuid"},
            "metadata": {"type": "object", "additionalProperties": True, "nullable": True},
            "created_at": _DT,
        },
    },
    "ActivityList": {
        "type": "array",
        "items": {"$ref": "#/components/schemas/ActivityResponse"},
    },
    "UploadResponse": {
        "type": "object",
        "properties": {
            "url": {"type": "string", "format": "uri"},
            "public_id": {"type": "string"},
            "resource_type": {"type": "string", "example": "image"},
            "format": {"type": "string", "example": "png"},
            "bytes": {"type": "integer", "example": 20480},
        },
    },
    # ── Modules ───────────────────────────────────────────────────────────
    "ModuleResponse": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "api_name": {"type": "string", "example": "deals"},
            "singular_label": {"type": "string", "example": "Deal"},
            "plural_label": {"type": "string", "example": "Deals"},
            "description": {"type": "string", "nullable": True},
            "icon": {"type": "string", "nullable": True, "example": "handshake"},
            "color": {"type": "string", "nullable": True, "example": "#0F766E"},
            "storage_strategy": {"type": "string", "enum": ["native", "dynamic"]},
            "native_table": {"type": "string", "nullable": True},
            "is_system": {"type": "boolean"},
            "is_enabled": {"type": "boolean"},
            "is_visible_sidebar": {"type": "boolean"},
            "sort_order": {"type": "integer"},
            "default_sort_field": {"type": "string"},
            "default_sort_order": {"type": "string", "enum": ["asc", "desc"]},
            "created_at": _DT,
            "updated_at": _DT,
        },
    },
    "ModuleList": {
        "type": "array",
        "items": {"$ref": "#/components/schemas/ModuleResponse"},
    },
    "NavigationItem": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "api_name": {"type": "string"},
            "singular_label": {"type": "string"},
            "plural_label": {"type": "string"},
            "icon": {"type": "string", "nullable": True},
            "color": {"type": "string", "nullable": True},
            "sort_order": {"type": "integer"},
        },
    },
    "NavigationList": {
        "type": "array",
        "items": {"$ref": "#/components/schemas/NavigationItem"},
    },
    "CreateModuleRequest": {
        "type": "object",
        "required": ["api_name", "singular_label", "plural_label"],
        "properties": {
            "api_name": {"type": "string", "example": "deals"},
            "singular_label": {"type": "string", "example": "Deal"},
            "plural_label": {"type": "string", "example": "Deals"},
            "description": {"type": "string"},
            "icon": {"type": "string"},
            "color": {"type": "string"},
            "is_visible_sidebar": {"type": "boolean", "default": True},
            "default_sort_field": {"type": "string"},
            "default_sort_order": {"type": "string", "enum": ["asc", "desc"]},
        },
    },
    "UpdateModuleRequest": {
        "type": "object",
        "properties": {
            "singular_label": {"type": "string"},
            "plural_label": {"type": "string"},
            "description": {"type": "string", "nullable": True},
            "icon": {"type": "string", "nullable": True},
            "color": {"type": "string", "nullable": True},
            "is_visible_sidebar": {"type": "boolean"},
            "default_sort_field": {"type": "string"},
            "default_sort_order": {"type": "string", "enum": ["asc", "desc"]},
        },
    },
    "SetStatusRequest": {
        "type": "object",
        "required": ["enabled"],
        "properties": {
            "enabled": {"type": "boolean", "example": True},
        },
    },
    "ReorderRequest": {
        "type": "object",
        "required": ["items"],
        "properties": {
            "items": {
                "type": "array",
                "items": {
                    "type": "object",
                    "required": ["id", "sort_order"],
                    "properties": {
                        "id": {"type": "string", "format": "uuid"},
                        "sort_order": {"type": "integer"},
                    },
                },
            },
        },
    },
    # ── Fields ────────────────────────────────────────────────────────────
    "FieldOption": {
        "type": "object",
        "required": ["label", "value"],
        "properties": {
            "label": {"type": "string", "example": "Hot"},
            "value": {"type": "string", "example": "hot"},
        },
    },
    "FieldResponse": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "module_id": {"type": "string", "format": "uuid"},
            "api_name": {"type": "string", "example": "amount"},
            "label": {"type": "string", "example": "Amount"},
            "field_type": {
                "type": "string",
                "enum": [
                    "text",
                    "textarea",
                    "email",
                    "phone",
                    "number",
                    "currency",
                    "date",
                    "datetime",
                    "boolean",
                    "dropdown",
                    "multiselect",
                    "radio",
                    "checkbox",
                    "url",
                    "file",
                    "image",
                    "user",
                    "lookup",
                    "formula",
                    "json",
                    "richtext",
                ],
            },
            "is_required": {"type": "boolean"},
            "is_unique": {"type": "boolean"},
            "is_read_only": {"type": "boolean"},
            "is_visible": {"type": "boolean"},
            "is_searchable": {"type": "boolean"},
            "is_filterable": {"type": "boolean"},
            "is_nullable": {"type": "boolean"},
            "is_indexed": {"type": "boolean"},
            "is_system": {"type": "boolean"},
            "default_value": {"type": "string", "nullable": True},
            "placeholder": {"type": "string", "nullable": True},
            "description": {"type": "string", "nullable": True},
            "help_text": {"type": "string", "nullable": True},
            "regex": {"type": "string", "nullable": True},
            "validation_message": {"type": "string", "nullable": True},
            "lookup_module_id": {"type": "string", "format": "uuid", "nullable": True},
            "min_length": {"type": "integer", "nullable": True},
            "max_length": {"type": "integer", "nullable": True},
            "sort_order": {"type": "integer"},
            "options": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/FieldOption"},
            },
            "storage": {
                "type": "object",
                "properties": {
                    "kind": {"type": "string", "enum": ["column", "jsonb"]},
                    "path": {"type": "string"},
                },
            },
            "created_at": _DT,
            "updated_at": _DT,
        },
    },
    "FieldList": {
        "type": "array",
        "items": {"$ref": "#/components/schemas/FieldResponse"},
    },
    "CreateFieldRequest": {
        "type": "object",
        "required": ["api_name", "label", "field_type"],
        "properties": {
            "api_name": {"type": "string", "example": "amount"},
            "label": {"type": "string", "example": "Amount"},
            "field_type": {
                "type": "string",
                "enum": [
                    "text",
                    "textarea",
                    "email",
                    "phone",
                    "number",
                    "currency",
                    "date",
                    "datetime",
                    "boolean",
                    "dropdown",
                    "multiselect",
                    "radio",
                    "checkbox",
                    "url",
                    "file",
                    "image",
                    "user",
                    "lookup",
                    "formula",
                    "json",
                    "richtext",
                ],
            },
            "is_required": {"type": "boolean"},
            "is_unique": {"type": "boolean"},
            "is_read_only": {"type": "boolean"},
            "is_visible": {"type": "boolean", "nullable": True},
            "is_searchable": {"type": "boolean"},
            "is_filterable": {"type": "boolean"},
            "default_value": {"type": "string", "nullable": True},
            "placeholder": {"type": "string", "nullable": True},
            "description": {"type": "string", "nullable": True},
            "help_text": {"type": "string", "nullable": True},
            "regex": {"type": "string", "nullable": True},
            "validation_message": {"type": "string", "nullable": True},
            "lookup_module_id": {"type": "string", "format": "uuid", "nullable": True},
            "min_length": {"type": "integer", "nullable": True},
            "max_length": {"type": "integer", "nullable": True},
            "sort_order": {"type": "integer"},
            "options": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/FieldOption"},
            },
        },
    },
    "UpdateFieldRequest": {
        "type": "object",
        "properties": {
            "label": {"type": "string"},
            "is_required": {"type": "boolean"},
            "is_unique": {"type": "boolean"},
            "is_read_only": {"type": "boolean"},
            "is_visible": {"type": "boolean"},
            "is_searchable": {"type": "boolean"},
            "is_filterable": {"type": "boolean"},
            "default_value": {"type": "string", "nullable": True},
            "placeholder": {"type": "string", "nullable": True},
            "description": {"type": "string", "nullable": True},
            "help_text": {"type": "string", "nullable": True},
            "regex": {"type": "string", "nullable": True},
            "validation_message": {"type": "string", "nullable": True},
            "min_length": {"type": "integer", "nullable": True},
            "max_length": {"type": "integer", "nullable": True},
            "sort_order": {"type": "integer"},
            "options": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/FieldOption"},
            },
        },
    },
    # ── Validation ────────────────────────────────────────────────────────
    "CreateRuleRequest": {
        "type": "object",
        "required": ["rule_type"],
        "properties": {
            "field_id": {"type": "string", "format": "uuid", "nullable": True},
            "rule_type": {
                "type": "string",
                "enum": [
                    "required",
                    "min_length",
                    "max_length",
                    "min",
                    "max",
                    "pattern",
                    "email",
                    "url",
                    "in",
                    "not_in",
                    "required_if",
                ],
                "example": "min_length",
            },
            "params": {
                "type": "object",
                "additionalProperties": True,
                "example": {"value": 3},
            },
            "error_message": {"type": "string", "example": "Too short"},
            "is_active": {"type": "boolean", "default": True},
            "sort_order": {"type": "integer"},
        },
    },
    "UpdateRuleRequest": {
        "type": "object",
        "properties": {
            "params": {"type": "object", "additionalProperties": True},
            "error_message": {"type": "string", "nullable": True},
            "is_active": {"type": "boolean"},
            "sort_order": {"type": "integer"},
        },
    },
    "RuleResponse": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "module_id": {"type": "string", "format": "uuid"},
            "field_id": {"type": "string", "format": "uuid", "nullable": True},
            "rule_type": {"type": "string"},
            "params": {"type": "object", "additionalProperties": True},
            "error_message": {"type": "string", "nullable": True},
            "is_active": {"type": "boolean"},
            "sort_order": {"type": "integer"},
            "created_at": _DT,
            "updated_at": _DT,
        },
    },
    "RuleList": {
        "type": "array",
        "items": {"$ref": "#/components/schemas/RuleResponse"},
    },
    "ValidateRequest": {
        "type": "object",
        "required": ["data"],
        "properties": {
            "data": {
                "type": "object",
                "additionalProperties": True,
                "example": {"name": "Acme", "email": "bad"},
            },
        },
    },
    "ValidateResult": {
        "type": "object",
        "properties": {
            "valid": {"type": "boolean", "example": False},
            "errors": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/FieldValidationError"},
            },
        },
    },
    "FieldSchema": {
        "type": "object",
        "properties": {
            "api_name": {"type": "string"},
            "label": {"type": "string"},
            "type": {"type": "string"},
            "required": {"type": "boolean"},
            "min_length": {"type": "integer"},
            "max_length": {"type": "integer"},
            "min": {"type": "number"},
            "max": {"type": "number"},
            "pattern": {"type": "string"},
            "format": {"type": "string", "enum": ["email", "url"]},
            "options": {"type": "array", "items": {"type": "string"}},
            "multiple": {"type": "boolean"},
            "messages": {
                "type": "object",
                "additionalProperties": {"type": "string"},
            },
        },
    },
    "ValidationSchema": {
        "type": "object",
        "properties": {
            "module_id": {"type": "string", "format": "uuid"},
            "fields": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/FieldSchema"},
            },
        },
    },
    "FilterClause": {
        "type": "object",
        "required": ["field", "operator"],
        "properties": {
            "field": {"type": "string", "example": "amount"},
            "operator": {
                "type": "string",
                "enum": ["eq", "ne", "contains", "gt", "lt", "gte", "lte", "in"],
                "example": "gte",
            },
            "value": {"description": "Filter value; type depends on field"},
        },
    },
    # ── Views ─────────────────────────────────────────────────────────────
    "ViewResponse": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "module_id": {"type": "string", "format": "uuid"},
            "name": {"type": "string", "example": "Open deals"},
            "columns": {
                "type": "array",
                "items": {"type": "string"},
                "example": ["name", "amount", "status"],
            },
            "filters": {
                "type": "array",
                "items": {
                    "type": "object",
                    "properties": {
                        "field": {"type": "string"},
                        "operator": {"type": "string"},
                        "value": {},
                    },
                },
            },
            "sort": {
                "type": "object",
                "properties": {
                    "field": {"type": "string"},
                    "order": {"type": "string", "enum": ["asc", "desc"]},
                },
            },
            "is_default": {"type": "boolean"},
            "is_public": {"type": "boolean"},
            "owner_id": {"type": "string", "format": "uuid", "nullable": True},
            "is_owner": {"type": "boolean"},
            "created_at": _DT,
            "updated_at": _DT,
        },
    },
    "ViewList": {
        "type": "array",
        "items": {"$ref": "#/components/schemas/ViewResponse"},
    },
    "CreateViewRequest": {
        "type": "object",
        "required": ["name", "columns"],
        "properties": {
            "name": {"type": "string", "example": "Open deals"},
            "columns": {
                "type": "array",
                "items": {"type": "string"},
                "example": ["name", "amount"],
            },
            "filters": {
                "type": "array",
                "items": {
                    "type": "object",
                    "properties": {
                        "field": {"type": "string"},
                        "operator": {"type": "string"},
                        "value": {},
                    },
                },
            },
            "sort": {
                "type": "object",
                "properties": {
                    "field": {"type": "string"},
                    "order": {"type": "string", "enum": ["asc", "desc"]},
                },
            },
            "is_public": {"type": "boolean"},
        },
    },
    "UpdateViewRequest": {"$ref": "#/components/schemas/CreateViewRequest"},
    # ── Records ───────────────────────────────────────────────────────────
    "RecordResponse": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "module_id": {"type": "string", "format": "uuid"},
            "data": {
                "type": "object",
                "additionalProperties": True,
                "example": {"name": "Acme Inc", "amount": 12000},
            },
            "owner_id": {"type": "string", "format": "uuid", "nullable": True},
            "created_by": {"type": "string", "format": "uuid", "nullable": True},
            "updated_by": {"type": "string", "format": "uuid", "nullable": True},
            "relations": {
                "type": "object",
                "additionalProperties": {
                    "type": "object",
                    "properties": {
                        "id": {"type": "string"},
                        "label": {"type": "string"},
                    },
                },
                "nullable": True,
            },
            "created_at": _DT,
            "updated_at": _DT,
        },
    },
    "RecordListResult": {
        "type": "object",
        "properties": {
            "records": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/RecordResponse"},
            },
            "page": {"type": "integer"},
            "page_size": {"type": "integer"},
            "total": {"type": "integer", "format": "int64"},
            "total_pages": {"type": "integer"},
        },
    },
    "CreateRecordRequest": {
        "type": "object",
        "required": ["data"],
        "properties": {
            "data": {
                "type": "object",
                "additionalProperties": True,
                "example": {"name": "Acme Inc", "amount": 12000},
            },
            "owner_id": {"type": "string", "format": "uuid"},
        },
    },
    "UpdateRecordRequest": {"$ref": "#/components/schemas/CreateRecordRequest"},
    # ── Import ────────────────────────────────────────────────────────────
    "AnalyzeResult": {
        "type": "object",
        "properties": {
            "headers": {
                "type": "array",
                "items": {"type": "string"},
                "example": ["Name", "Email"],
            },
            "sample_rows": {
                "type": "array",
                "items": {"type": "object", "additionalProperties": True},
            },
            "suggested_mapping": {
                "type": "object",
                "additionalProperties": {"type": "string"},
                "example": {"Name": "name", "Email": "email"},
            },
            "row_count": {"type": "integer", "example": 250},
        },
    },
    "ImportResponse": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "module_id": {"type": "string", "format": "uuid"},
            "filename": {"type": "string"},
            "status": {
                "type": "string",
                "enum": ["pending", "processing", "completed", "failed"],
            },
            "mapping": {
                "type": "object",
                "additionalProperties": {"type": "string"},
            },
            "total_rows": {"type": "integer"},
            "processed_rows": {"type": "integer"},
            "success_rows": {"type": "integer"},
            "error_rows": {"type": "integer"},
            "errors": {
                "type": "array",
                "items": {
                    "type": "object",
                    "properties": {
                        "row": {"type": "integer"},
                        "field": {"type": "string"},
                        "message": {"type": "string"},
                    },
                },
            },
            "created_by": {"type": "string", "format": "uuid", "nullable": True},
            "started_at": {"type": "string", "format": "date-time", "nullable": True},
            "finished_at": {"type": "string", "format": "date-time", "nullable": True},
            "created_at": _DT,
            "updated_at": _DT,
        },
    },
    "ImportListResult": {
        "type": "object",
        "properties": {
            "imports": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/ImportResponse"},
            },
            "page": {"type": "integer"},
            "page_size": {"type": "integer"},
            "total": {"type": "integer", "format": "int64"},
            "total_pages": {"type": "integer"},
        },
    },
    # ── Export ────────────────────────────────────────────────────────────
    "ExportSpec": {
        "type": "object",
        "properties": {
            "format": {"type": "string", "enum": ["csv", "xlsx"], "example": "xlsx"},
            "columns": {
                "type": "array",
                "items": {"type": "string"},
                "example": ["name", "amount"],
            },
            "filters": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/FilterClause"},
            },
            "search": {"type": "string"},
            "sort": {"type": "string"},
            "order": {"type": "string", "enum": ["asc", "desc"]},
            "expand": {"type": "boolean", "example": False},
        },
    },
    "ExportResponse": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "module_id": {"type": "string", "format": "uuid"},
            "filename": {"type": "string"},
            "format": {"type": "string", "enum": ["csv", "xlsx"]},
            "status": {
                "type": "string",
                "enum": ["pending", "processing", "completed", "failed"],
            },
            "columns": {"type": "array", "items": {"type": "string"}},
            "row_count": {"type": "integer"},
            "byte_size": {"type": "integer"},
            "error": {"type": "string", "nullable": True},
            "created_by": {"type": "string", "format": "uuid", "nullable": True},
            "started_at": {"type": "string", "format": "date-time", "nullable": True},
            "finished_at": {"type": "string", "format": "date-time", "nullable": True},
            "created_at": _DT,
            "updated_at": _DT,
        },
    },
    "ExportListResult": {
        "type": "object",
        "properties": {
            "exports": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/ExportResponse"},
            },
            "page": {"type": "integer"},
            "page_size": {"type": "integer"},
            "total": {"type": "integer", "format": "int64"},
            "total_pages": {"type": "integer"},
        },
    },
    "CreateTemplateRequest": {
        "type": "object",
        "required": ["name", "columns"],
        "properties": {
            "name": {"type": "string", "example": "Finance export"},
            "format": {"type": "string", "enum": ["csv", "xlsx"]},
            "columns": {
                "type": "array",
                "items": {"type": "string"},
                "example": ["name", "amount"],
            },
            "filters": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/FilterClause"},
            },
            "sort": {
                "type": "object",
                "properties": {
                    "field": {"type": "string"},
                    "order": {"type": "string", "enum": ["asc", "desc"]},
                },
            },
        },
    },
    "UpdateTemplateRequest": {"$ref": "#/components/schemas/CreateTemplateRequest"},
    "TemplateResponse": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "module_id": {"type": "string", "format": "uuid"},
            "name": {"type": "string"},
            "format": {"type": "string", "enum": ["csv", "xlsx"]},
            "columns": {"type": "array", "items": {"type": "string"}},
            "filters": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/FilterClause"},
            },
            "sort": {
                "type": "object",
                "properties": {
                    "field": {"type": "string"},
                    "order": {"type": "string", "enum": ["asc", "desc"]},
                },
            },
            "created_by": {"type": "string", "format": "uuid", "nullable": True},
            "created_at": _DT,
            "updated_at": _DT,
        },
    },
    "TemplateList": {
        "type": "array",
        "items": {"$ref": "#/components/schemas/TemplateResponse"},
    },
    # ── Notifications ─────────────────────────────────────────────────────
    "SendNotificationRequest": {
        "type": "object",
        "required": ["channel", "to", "body"],
        "properties": {
            "channel": {
                "type": "string",
                "enum": ["email", "whatsapp"],
                "example": "email",
            },
            "to": {"type": "string", "example": "ada@example.com"},
            "subject": {"type": "string", "example": "Welcome"},
            "body": {"type": "string", "example": "Thanks for joining"},
            "template": {"type": "string", "example": "lead_welcome"},
            "data": {"type": "object", "additionalProperties": True},
            "entity_type": {"type": "string"},
            "entity_id": {"type": "string", "format": "uuid"},
        },
    },
    "NotificationResponse": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "channel": {"type": "string", "enum": ["email", "whatsapp"]},
            "recipient": {"type": "string"},
            "subject": {"type": "string", "nullable": True},
            "body": {"type": "string"},
            "template": {"type": "string", "nullable": True},
            "data": {"type": "object", "additionalProperties": True},
            "status": {
                "type": "string",
                "enum": ["queued", "sent", "failed"],
            },
            "provider": {"type": "string", "nullable": True},
            "error": {"type": "string", "nullable": True},
            "entity_type": {"type": "string", "nullable": True},
            "entity_id": {"type": "string", "format": "uuid", "nullable": True},
            "created_by": {"type": "string", "format": "uuid", "nullable": True},
            "sent_at": {"type": "string", "format": "date-time", "nullable": True},
            "created_at": _DT,
            "updated_at": _DT,
        },
    },
    "NotificationListResult": {
        "type": "object",
        "properties": {
            "notifications": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/NotificationResponse"},
            },
            "page": {"type": "integer"},
            "page_size": {"type": "integer"},
            "total": {"type": "integer", "format": "int64"},
            "total_pages": {"type": "integer"},
        },
    },
    # ── Tour ──────────────────────────────────────────────────────────────
    "UpdateProgressRequest": {
        "type": "object",
        "properties": {
            "tour_key": {"type": "string", "example": "crm_onboarding"},
            "status": {
                "type": "string",
                "enum": ["active", "completed", "skipped"],
                "example": "active",
            },
            "current_step": {"type": "integer", "example": 2},
            "completed_steps": {
                "type": "array",
                "items": {"type": "string"},
                "example": ["welcome", "leads"],
            },
        },
    },
    "RestartRequest": {
        "type": "object",
        "properties": {
            "tour_key": {"type": "string", "example": "crm_onboarding"},
        },
    },
    "ProgressResponse": {
        "type": "object",
        "properties": {
            "tour_key": {"type": "string"},
            "status": {
                "type": "string",
                "enum": ["active", "completed", "skipped"],
            },
            "current_step": {"type": "integer"},
            "completed_steps": {"type": "array", "items": {"type": "string"}},
            "started_at": _DT,
            "completed_at": {"type": "string", "format": "date-time", "nullable": True},
            "updated_at": _DT,
        },
    },
    # ── Settings ──────────────────────────────────────────────────────────
    "GeneralSettings": {
        "type": "object",
        "properties": {
            "timezone": {"type": "string", "example": "Asia/Kolkata"},
            "date_format": {"type": "string", "example": "DD/MM/YYYY"},
            "time_format": {"type": "string", "enum": ["12h", "24h"], "example": "24h"},
            "currency": {"type": "string", "example": "INR"},
            "locale": {"type": "string", "example": "en-IN"},
            "week_start": {
                "type": "string",
                "enum": ["sunday", "monday"],
                "example": "monday",
            },
        },
    },
    "AutomationSettings": {
        "type": "object",
        "properties": {
            "notifications_enabled": {"type": "boolean"},
            "default_channel": {
                "type": "string",
                "enum": ["whatsapp", "email"],
            },
            "daily_digest": {"type": "boolean"},
        },
    },
    "SettingsResponse": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "name": {"type": "string", "example": "Acme Org"},
            "slug": {"type": "string", "example": "acme"},
            "plan": {"type": "string", "example": "pro"},
            "general": {"$ref": "#/components/schemas/GeneralSettings"},
            "automation": {"$ref": "#/components/schemas/AutomationSettings"},
            "updated_at": _DT,
        },
    },
    "UpdateSettingsRequest": {
        "type": "object",
        "properties": {
            "name": {"type": "string"},
            "general": {"$ref": "#/components/schemas/GeneralSettings"},
            "automation": {"$ref": "#/components/schemas/AutomationSettings"},
        },
    },
    # ── Roles / RBAC ──────────────────────────────────────────────────────
    "PermissionResponse": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "key": {"type": "string", "example": "record.view"},
            "category": {"type": "string", "example": "records"},
            "description": {"type": "string", "nullable": True},
        },
    },
    "PermissionList": {
        "type": "array",
        "items": {"$ref": "#/components/schemas/PermissionResponse"},
    },
    "ModuleAccess": {
        "type": "object",
        "required": ["module_id", "can_view", "can_create", "can_update", "can_delete"],
        "properties": {
            "module_id": {"type": "string", "format": "uuid"},
            "can_view": {"type": "boolean"},
            "can_create": {"type": "boolean"},
            "can_update": {"type": "boolean"},
            "can_delete": {"type": "boolean"},
        },
    },
    "FieldAccess": {
        "type": "object",
        "required": ["field_id", "access"],
        "properties": {
            "field_id": {"type": "string", "format": "uuid"},
            "module_id": {"type": "string", "format": "uuid"},
            "access": {
                "type": "string",
                "enum": ["hidden", "read", "write"],
                "example": "read",
            },
        },
    },
    "RoleSummary": {
        "type": "object",
        "properties": {
            "id": _UUID,
            "name": {"type": "string", "example": "Sales Rep"},
            "slug": {"type": "string", "example": "sales_rep"},
            "description": {"type": "string", "nullable": True},
            "is_system": {"type": "boolean"},
            "member_count": {"type": "integer", "example": 4},
            "created_at": _DT,
            "updated_at": _DT,
        },
    },
    "RoleSummaryList": {
        "type": "array",
        "items": {"$ref": "#/components/schemas/RoleSummary"},
    },
    "RoleDetail": {
        "allOf": [
            {"$ref": "#/components/schemas/RoleSummary"},
            {
                "type": "object",
                "properties": {
                    "permissions": {
                        "type": "array",
                        "items": {"type": "string"},
                        "example": ["module.view", "record.view", "export.run"],
                    },
                    "module_access": {
                        "type": "array",
                        "items": {"$ref": "#/components/schemas/ModuleAccess"},
                    },
                    "field_access": {
                        "type": "array",
                        "items": {"$ref": "#/components/schemas/FieldAccess"},
                    },
                },
            },
        ],
    },
    "CreateRoleRequest": {
        "type": "object",
        "required": ["name", "slug"],
        "properties": {
            "name": {"type": "string", "example": "Custom"},
            "slug": {"type": "string", "example": "custom_role"},
            "description": {"type": "string"},
        },
    },
    "UpdateRoleRequest": {
        "type": "object",
        "properties": {
            "name": {"type": "string"},
            "description": {"type": "string", "nullable": True},
        },
    },
    "SetPermissionsRequest": {
        "type": "object",
        "required": ["permissions"],
        "properties": {
            "permissions": {
                "type": "array",
                "items": {"type": "string"},
                "example": ["module.view", "record.view", "export.run"],
            },
        },
    },
    "SetModuleAccessRequest": {
        "type": "object",
        "required": ["access"],
        "properties": {
            "access": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/ModuleAccess"},
            },
        },
    },
    "SetFieldAccessRequest": {
        "type": "object",
        "required": ["access"],
        "properties": {
            "access": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/FieldAccess"},
            },
        },
    },
    "MeResponse": {
        "type": "object",
        "properties": {
            "role_id": {"type": "string", "format": "uuid"},
            "role_slug": {"type": "string", "example": "admin"},
            "permissions": {
                "type": "array",
                "items": {"type": "string"},
            },
            "module_access": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/ModuleAccess"},
            },
            "field_access": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/FieldAccess"},
            },
        },
    },
}
