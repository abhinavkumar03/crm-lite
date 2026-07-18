# Product roadmap (tenancy-first)

Canonical ordering for reviewers and future SaaS work. Production CRMs establish
**tenant isolation → org structure → authorization → record visibility** before
metadata-driven CRM features.

Historical implementation phases in [`COMMANDS.md`](../COMMANDS.md) keep their
original numbers (Phase 3 schema, Phase 16 RBAC, …). Use the mapping table below.

## Status legend

| Status | Meaning |
| --- | --- |
| Done | Shipped in repo |
| Partial | Core exists; gaps listed |
| Next | Active engineering focus |

## Roadmap

### Phase 1 — Multi-tenant foundation

**Status: Done (harden)**

- Organizations + `organization_id` on business metadata/runtime tables
- Tenant middleware (`internal/tenant`) injects org context after JWT
- Org-scoped repositories and Settings JSON (`settings.general` for locale/currency/timezone)

**Harden:** richer org profile columns (`logo_url`, `industry`, `company_size`, `country`, `status`, `created_by`).

### Phase 2 — Organization and membership

**Status: Partial → Next**

- Membership model (multi-org capable schema)
- Active organization + switch APIs
- Organization create + **bootstrap service**
- Invitations (email simulation → accept → join)
- Members list / assign role

### Phase 3 — Roles, permissions, and hierarchy

**Status: Partial → Next**

- Org-scoped roles + global permission catalog + module/field ACL (**done**)
- Role `hierarchy_level` + reporting tree on members
- Departments, teams, branches
- Expand permission catalog (`organization.manage`, analytics, …)

### Phase 4 — Record visibility engine

**Status: Next**

- Ownership columns on `records` (`assigned_to`, `team_id`, `department_id`, `visibility`)
- Shared access service used by record list/get/update/delete
- Hierarchy / manager / department / organization visibility modes

### Phase 5 — Migration and seeder framework

**Status: Partial → Next**

- Existing golang-migrate + seed runner (**done**)
- Extend seed: org profile, reporting tree, departments, visibility demo records

### Phase 6+ — Metadata CRM (already shipped)

Dynamic modules, fields, validation, views, records, notifications, import/export,
guided tour, settings UI, performance, OpenAPI, documentation.

## Mapping: new roadmap ↔ historical build order

| Roadmap | Historical / packages |
| --- | --- |
| Phase 1 | `000003` orgs/members, `internal/tenant`, settings |
| Phase 2 | New: org APIs, invites, bootstrap, `active_organization_id` |
| Phase 3 | Phase 16 RBAC + new hierarchy/depts |
| Phase 4 | New: `internal/access` + record filters |
| Phase 5 | `cmd/seed` + seeders |
| Phase 6+ | Old Phases 6–15, 17–19 (metadata → docs) |

## Isolation rule

Every business query must filter by `organization_id` from tenant context.
Cross-tenant joins are forbidden. Row visibility is an additional filter **within**
the active organization.
