# Entity-relationship diagrams

PostgreSQL schema from migrations `000001`–`000011`. Diagrams use Mermaid ER
syntax (render in GitHub / most Markdown previewers).

## Identity & tenancy

```mermaid
erDiagram
  users ||--o{ organization_members : has
  organizations ||--o{ organization_members : has
  roles ||--o{ organization_members : "optional role"
  organizations ||--o{ organization_invitations : invites
  organizations ||--o{ departments : owns
  organizations ||--o{ teams : owns
  organizations ||--o{ branches : owns
  departments ||--o{ teams : groups

  users {
    uuid id PK
    text name
    text email UK
    text password_hash
    uuid active_organization_id FK
  }
  organizations {
    uuid id PK
    text name
    text slug UK
    text plan
    text logo_url
    text industry
    text company_size
    text country
    text status
    uuid created_by FK
    jsonb settings
  }
  organization_members {
    uuid id PK
    uuid organization_id FK
    uuid user_id FK
    uuid role_id FK
    uuid manager_user_id FK
    uuid department_id FK
    uuid team_id FK
    uuid branch_id FK
    text designation
    int hierarchy_level
    text status
  }
  organization_invitations {
    uuid id PK
    uuid organization_id FK
    text email
    uuid role_id FK
    text token UK
    text status
    timestamptz expires_at
  }
  departments {
    uuid id PK
    uuid organization_id FK
    text name
  }
  teams {
    uuid id PK
    uuid organization_id FK
    uuid department_id FK
    text name
  }
  branches {
    uuid id PK
    uuid organization_id FK
    text name
    text location
  }
  roles {
    uuid id PK
    uuid organization_id FK
    text name
    text slug
    boolean is_system
    int hierarchy_level
    uuid parent_role_id FK
  }
```

## RBAC

```mermaid
erDiagram
  roles ||--o{ role_permissions : grants
  permissions ||--o{ role_permissions : granted_by
  roles ||--o{ role_module_access : restricts
  modules ||--o{ role_module_access : for
  roles ||--o{ role_field_access : restricts
  fields ||--o{ role_field_access : for

  permissions {
    uuid id PK
    text key UK
    text category
  }
  role_permissions {
    uuid role_id PK
    uuid permission_id PK
  }
  role_module_access {
    uuid role_id PK
    uuid module_id PK
    boolean can_view
    boolean can_create
    boolean can_update
    boolean can_delete
  }
  role_field_access {
    uuid role_id PK
    uuid field_id PK
    text access
  }
```

`access` ∈ `hidden` | `read` | `write`. No ACL row ⇒ unrestricted at that layer.

## Native CRM

```mermaid
erDiagram
  users ||--o{ leads : owns
  users ||--o{ contacts : owns
  users ||--o{ tasks : owns
  leads ||--o{ tasks : "optional"
  contacts ||--o{ tasks : "optional"
  leads ||--o{ activity_logs : has

  leads {
    uuid id PK
    uuid owner_id FK
    text name
    text email
    text status
  }
  contacts {
    uuid id PK
    uuid owner_id FK
    text first_name
    text last_name
    text email
  }
  tasks {
    uuid id PK
    uuid owner_id FK
    uuid lead_id FK
    uuid contact_id FK
    text title
    text status
    timestamptz due_date
  }
  activity_logs {
    uuid id PK
    uuid lead_id FK
    text activity_type
    text message
  }
```

Polymorphic side-tables (no FK to leads/contacts/tasks — keyed by
`entity_type` + `entity_id`):

| Table | Purpose |
| --- | --- |
| `notes` | Free-text notes |
| `call_logs` | Call direction / status / duration |
| `attachments` | Cloudinary metadata |
| `activities` | Audit trail (`action`, `metadata` JSONB) |

## Metadata engine

```mermaid
erDiagram
  organizations ||--o{ modules : owns
  modules ||--o{ fields : defines
  modules ||--o{ layouts : has
  modules ||--o{ views : has
  modules ||--o{ validation_rules : has
  modules ||--o{ records : stores
  fields ||--o{ validation_rules : "optional field"

  modules {
    uuid id PK
    uuid organization_id FK
    text api_name
    text storage_strategy
    text native_table
    boolean is_system
    boolean is_enabled
  }
  fields {
    uuid id PK
    uuid module_id FK
    text api_name
    text field_type
    boolean is_required
    jsonb options
    uuid lookup_module_id FK
  }
  validation_rules {
    uuid id PK
    uuid module_id FK
    uuid field_id FK
    text rule_type
    jsonb params
    boolean is_active
  }
  views {
    uuid id PK
    uuid module_id FK
    text name
    jsonb columns
    jsonb filters
    jsonb sort
  }
  records {
    uuid id PK
    uuid module_id FK
    jsonb data
    uuid owner_id FK
    uuid assigned_to FK
    uuid team_id FK
    uuid department_id FK
    text visibility
  }
```

`visibility` ∈ `private` | `owner` | `manager` | `hierarchy` | `department` | `organization` | `team` | `public`.

`storage_strategy`:

- `native` — data in a first-class table (`native_table`)
- `dynamic` — data in `records.data` (GIN-indexed)

Also present: `layouts`, `automation_rules`, `import_templates`, `export_templates`.

## Jobs & notifications

```mermaid
erDiagram
  organizations ||--o{ notifications : owns
  organizations ||--o{ import_jobs : owns
  organizations ||--o{ export_jobs : owns
  modules ||--o{ import_jobs : targets
  modules ||--o{ export_jobs : targets

  notifications {
    uuid id PK
    text channel
    text recipient
    text status
    text provider
    jsonb data
  }
  import_jobs {
    uuid id PK
    text filename
    text status
    jsonb mapping
    jsonb source_rows
    jsonb errors
    int total_rows
  }
  export_jobs {
    uuid id PK
    text filename
    text format
    text status
    bytea content
    int row_count
  }
```

Statuses: import/export `pending` → `processing` → `completed` | `failed`;
notifications `queued` → `sent` | `failed`.

## Tour

```mermaid
erDiagram
  organizations ||--o{ tour_steps : "optional org"
  organizations ||--o{ tour_progress : tracks
  users ||--o{ tour_progress : tracks

  tour_steps {
    uuid id PK
    uuid organization_id FK
    text step_key
    text title
    text target_selector
    text route
    int sort_order
  }
  tour_progress {
    uuid id PK
    uuid organization_id FK
    uuid user_id FK
    text tour_key
    text status
    int current_step
    jsonb completed_steps
  }
```

## Indexes of note (Phase 17)

- `records (organization_id, module_id, created_at DESC)`
- `notifications (organization_id, created_at DESC)`
- GIN on `records.data` (`jsonb_path_ops`)

See `backend/migrations/` for the authoritative SQL.
