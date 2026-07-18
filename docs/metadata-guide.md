# Metadata guide

CRM Lite’s customizable objects are described in metadata tables and executed by
engines (module, field, validation, view, record). Settings UI under
`/settings/*` is a thin admin surface over those engines.

## Concepts

| Concept | Table | Meaning |
| --- | --- | --- |
| Module | `modules` | An object type (`deal`, `company`, …) |
| Field | `fields` | Typed attribute on a module (`api_name`, `field_type`) |
| Record | `records` | One row of a **dynamic** module (`data` JSONB) |
| Validation rule | `validation_rules` | DB-driven constraint (field or module level) |
| View | `views` | Saved list columns / filters / sort |
| Layout | `layouts` | Form/detail layout config (JSONB) |

Everything above is **organization-scoped** (`organization_id`).

## Storage strategies

| `storage_strategy` | Persistence | APIs |
| --- | --- | --- |
| `native` | Dedicated SQL table (`native_table`) | Lead/contact/task style modules |
| `dynamic` | `records` + `data` JSONB | Record runtime, import, export |

Creating a module via Settings defaults to **dynamic**. System native modules
are seeded and usually not deleted.

```text
modules.api_name = "deal"
modules.storage_strategy = "dynamic"
fields: name (text), amount (currency), stage (dropdown)
records.data = { "name": "Acme", "amount": 12000, "stage": "proposal" }
```

## Field types

Common `field_type` values (see migration CHECK / field package):

`text`, `textarea`, `number`, `currency`, `email`, `phone`, `url`, `date`,
`datetime`, `boolean`, `dropdown`, `radio`, `multiselect`, `lookup`, `user`.

Notable flags: `is_required`, `is_unique`, `is_searchable`, `is_filterable`,
`is_visible`, `is_system`. Dropdowns store `options` as JSON
`[{label,value}, …]`. Lookups reference `lookup_module_id`.

## Validation engine

Rules live in `validation_rules`:

- **Field-level** — `field_id` set (`required`, `min_length`, `pattern`, …)
- **Module-level** — `field_id` null (e.g. `required_if`)

Params are JSONB (shape depends on `rule_type`). The engine merges field
metadata + active rules for:

- `POST /modules/:id/validate` — dry-run
- `GET /modules/:id/validation-schema` — frontend schema
- Record create/update
- Import row validation (`LoadSpec` once per job)

Custom `error_message` wins over defaults.

## Views

A view stores:

```json
{
  "name": "Open deals",
  "columns": ["name", "amount", "stage"],
  "filters": [{ "field": "stage", "operator": "ne", "value": "closed" }],
  "sort": { "field": "amount", "order": "desc" },
  "is_default": true
}
```

Filter operators: `eq`, `ne`, `contains`, `gt`, `lt`, `gte`, `lte`, `in`.

## Record runtime

Base: `/api/v1/modules/:moduleId/records`

| Method | Action |
| --- | --- |
| `GET` | List — `page`, `page_size`, `search`, `sort`, `order`, `filters`, `expand` |
| `POST` | Create — `{ data, owner_id? }` |
| `GET /:id` | Get (optional `expand` for lookup labels) |
| `PUT /:id` | Replace `data` |
| `DELETE /:id` | Delete |

`expand=true` resolves lookup/user fields into `relations.{api_name}.{id,label}`.

RBAC: `record.*` permissions **and** optional module/field ACL. Hidden fields
are stripped; read-only fields reject writes.

## Settings surfaces

| UI route | Engine |
| --- | --- |
| `/settings/modules` | Module CRUD / enable / reorder |
| `/settings/fields` | Field CRUD / reorder |
| `/settings/validation` | Rule CRUD |
| `/settings/roles` | Permissions + ACL |
| `/tables`, `/forms` | Consume metadata + records |

Frontend caches modules/fields/schemas for 60s (`features/metadata/cache.ts`);
Settings mutations invalidate the cache.

## Feature flags

| Env | Effect |
| --- | --- |
| `FEATURE_DYNAMIC_MODULES` | Gate dynamic module UX |
| `FEATURE_RBAC` | Gate roles UI / enforcement packaging |

## Related

- [ERD](./erd.md) — metadata tables
- [Import / export guide](./import-export-guide.md)
- [Sequences](./sequences.md) — record create
- OpenAPI tag **Modules / Fields / Validation / Views / Records**
