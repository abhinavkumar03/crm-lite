#!/usr/bin/env python3
"""Generate openapi.yaml for CRM Lite. Run from this directory: python3 generate_openapi.py"""

from __future__ import annotations

import textwrap
from pathlib import Path

try:
    import yaml
except ImportError:
    import json
    import sys

    print("PyYAML required: pip install pyyaml", file=sys.stderr)
    sys.exit(1)


def op(
    summary,
    tag,
    auth=True,
    req=None,
    responses=None,
    params=None,
    description=None,
    op_id=None,
):
    d = {"summary": summary, "tags": [tag]}
    if op_id:
        d["operationId"] = op_id
    if description:
        d["description"] = description
    if auth:
        d["security"] = [{"bearerAuth": []}]
    if params:
        d["parameters"] = params
    if req:
        d["requestBody"] = req
    d["responses"] = responses or {}
    return d


def json_body(schema_ref, example=None, required=True):
    content = {"application/json": {"schema": {"$ref": schema_ref}}}
    if example is not None:
        content["application/json"]["example"] = example
    return {"required": required, "content": content}


def multipart_body(parts):
    props, required = {}, []
    for name, schema, req in parts:
        props[name] = schema
        if req:
            required.append(name)
    return {
        "required": True,
        "content": {
            "multipart/form-data": {
                "schema": {
                    "type": "object",
                    "required": required,
                    "properties": props,
                }
            }
        },
    }


def ok(schema_ref=None, example=None, desc="Success", code="200"):
    schema = {"$ref": "#/components/schemas/APIResponse"}
    if schema_ref:
        schema = {
            "allOf": [
                {"$ref": "#/components/schemas/APIResponse"},
                {"type": "object", "properties": {"data": {"$ref": schema_ref}}},
            ]
        }
    content = {"application/json": {"schema": schema}}
    if example is not None:
        content["application/json"]["example"] = example
    return {code: {"description": desc, "content": content}}


# Validation field errors are returned as HTTP 400 (BadRequest) in handlers.
CODE = {
    "Unauthorized": "401",
    "Forbidden": "403",
    "NotFound": "404",
    "BadRequest": "400",
    "ValidationError": "400",
    "Conflict": "409",
    "InternalError": "500",
}


def err_refs(*names):
    return {CODE[n]: {"$ref": f"#/components/responses/{n}"} for n in names}


def q(name, typ="string", desc="", required=False, example=None):
    p = {
        "name": name,
        "in": "query",
        "required": required,
        "schema": {"type": typ},
        "description": desc,
    }
    if example is not None:
        p["schema"]["example"] = example
    return p


def path_param(name, desc=""):
    return {
        "name": name,
        "in": "path",
        "required": True,
        "schema": {"type": "string", "format": "uuid"},
        "description": desc,
    }


def build_paths():
    paths = {}
    page_limit = [
        q("page", "integer", "1-based page", example=1),
        q("limit", "integer", "Page size (max 100)", example=20),
        q("search"),
        q("sort_by"),
        q("sort_order", example="desc"),
    ]

    paths["/health"] = {
        "get": op(
            "Health check",
            "Health",
            auth=False,
            op_id="healthCheck",
            description="Liveness probe. Not role-gated.",
            responses={
                **ok(
                    None,
                    {
                        "success": True,
                        "message": "Service is healthy",
                        "data": {"service": "crm-lite", "status": "UP"},
                    },
                ),
            },
        )
    }

    paths["/auth/register"] = {
        "post": op(
            "Register user",
            "Auth",
            auth=False,
            op_id="register",
            req=json_body(
                "#/components/schemas/RegisterRequest",
                {
                    "name": "Ada Lovelace",
                    "email": "ada@example.com",
                    "password": "Secret@12345",
                },
            ),
            responses={
                **ok(
                    "#/components/schemas/UserResponse",
                    {
                        "success": True,
                        "message": "Registered",
                        "data": {
                            "id": "11111111-1111-1111-1111-111111111111",
                            "name": "Ada Lovelace",
                            "email": "ada@example.com",
                        },
                    },
                ),
                **err_refs("BadRequest", "Conflict", "ValidationError", "InternalError"),
            },
        )
    }
    paths["/auth/login"] = {
        "post": op(
            "Login",
            "Auth",
            auth=False,
            op_id="login",
            req=json_body(
                "#/components/schemas/LoginRequest",
                {"email": "admin@crm.com", "password": "Admin@123"},
            ),
            responses={
                **ok(
                    "#/components/schemas/LoginResponse",
                    {
                        "success": True,
                        "message": "Logged in",
                        "data": {
                            "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
                            "user": {
                                "id": "…",
                                "name": "Admin",
                                "email": "admin@crm.com",
                            },
                        },
                    },
                ),
                **err_refs("Unauthorized", "BadRequest", "ValidationError"),
            },
        )
    }
    paths["/auth/profile"] = {
        "get": op(
            "Current user profile",
            "Auth",
            op_id="profile",
            responses={
                **ok("#/components/schemas/UserResponse"),
                **err_refs("Unauthorized"),
            },
        )
    }

    paths["/dashboard"] = {
        "get": op(
            "Dashboard metrics",
            "Dashboard",
            op_id="getDashboard",
            params=[q("refresh", "boolean", "Bypass Redis cache", example=False)],
            responses={
                **ok("#/components/schemas/DashboardResponse"),
                **err_refs("Unauthorized"),
            },
        )
    }
    paths["/search"] = {
        "get": op(
            "Global search",
            "Search",
            op_id="search",
            params=[q("q", required=True, example="acme")],
            responses={
                **ok("#/components/schemas/SearchResponse"),
                **err_refs("Unauthorized", "BadRequest"),
            },
        )
    }

    paths["/attachments/{attachmentId}"] = {
        "delete": op(
            "Delete attachment",
            "Attachments",
            op_id="deleteAttachment",
            params=[path_param("attachmentId")],
            responses={
                "204": {"description": "Deleted"},
                **err_refs("Unauthorized", "NotFound"),
            },
        ),
    }
    paths["/uploads"] = {
        "post": op(
            "Upload media file",
            "Media",
            op_id="uploadMedia",
            description="Multipart upload to the configured media provider.",
            req=multipart_body(
                [("file", {"type": "string", "format": "binary"}, True)]
            ),
            responses={
                **ok("#/components/schemas/UploadResponse", code="201"),
                **err_refs("Unauthorized", "BadRequest"),
            },
        ),
    }

    # Metadata engine
    paths["/modules"] = {
        "get": op(
            "List modules",
            "Modules",
            op_id="listModules",
            responses={
                **ok("#/components/schemas/ModuleList"),
                **err_refs("Unauthorized", "Forbidden"),
            },
        ),
        "post": op(
            "Create module",
            "Modules",
            op_id="createModule",
            req=json_body("#/components/schemas/CreateModuleRequest"),
            responses={
                **ok("#/components/schemas/ModuleResponse", code="201"),
                **err_refs("Unauthorized", "Forbidden", "BadRequest", "Conflict"),
            },
        ),
    }
    paths["/modules/reorder"] = {
        "post": op(
            "Reorder modules",
            "Modules",
            op_id="reorderModules",
            req=json_body("#/components/schemas/ReorderRequest"),
            responses={
                **ok(),
                **err_refs("Unauthorized", "Forbidden", "BadRequest"),
            },
        )
    }
    paths["/modules/{id}"] = {
        "get": op(
            "Get module",
            "Modules",
            op_id="getModule",
            params=[path_param("id")],
            responses={
                **ok("#/components/schemas/ModuleResponse"),
                **err_refs("Unauthorized", "Forbidden", "NotFound"),
            },
        ),
        "put": op(
            "Update module",
            "Modules",
            op_id="updateModule",
            params=[path_param("id")],
            req=json_body("#/components/schemas/UpdateModuleRequest"),
            responses={
                **ok("#/components/schemas/ModuleResponse"),
                **err_refs("Unauthorized", "Forbidden", "NotFound", "BadRequest"),
            },
        ),
        "delete": op(
            "Delete module",
            "Modules",
            op_id="deleteModule",
            params=[path_param("id")],
            responses={
                "204": {"description": "Deleted"},
                **err_refs("Unauthorized", "Forbidden", "NotFound", "Conflict"),
            },
        ),
    }
    paths["/modules/{id}/status"] = {
        "patch": op(
            "Enable/disable module",
            "Modules",
            op_id="setModuleStatus",
            params=[path_param("id")],
            req=json_body("#/components/schemas/SetStatusRequest", {"enabled": True}),
            responses={
                **ok("#/components/schemas/ModuleResponse"),
                **err_refs("Unauthorized", "Forbidden", "NotFound"),
            },
        )
    }
    paths["/navigation"] = {
        "get": op(
            "Sidebar navigation",
            "Modules",
            op_id="navigation",
            responses={
                **ok("#/components/schemas/NavigationList"),
                **err_refs("Unauthorized", "Forbidden"),
            },
        )
    }

    paths["/modules/{id}/fields"] = {
        "get": op(
            "List fields",
            "Fields",
            op_id="listFields",
            params=[path_param("id")],
            responses={
                **ok("#/components/schemas/FieldList"),
                **err_refs("Unauthorized", "Forbidden", "NotFound"),
            },
        ),
        "post": op(
            "Create field",
            "Fields",
            op_id="createField",
            params=[path_param("id")],
            req=json_body("#/components/schemas/CreateFieldRequest"),
            responses={
                **ok("#/components/schemas/FieldResponse", code="201"),
                **err_refs("Unauthorized", "Forbidden", "BadRequest", "Conflict"),
            },
        ),
    }
    paths["/modules/{id}/fields/reorder"] = {
        "post": op(
            "Reorder fields",
            "Fields",
            op_id="reorderFields",
            params=[path_param("id")],
            req=json_body("#/components/schemas/ReorderRequest"),
            responses={
                **ok(),
                **err_refs("Unauthorized", "Forbidden", "BadRequest"),
            },
        )
    }
    paths["/modules/{id}/fields/{fieldId}"] = {
        "get": op(
            "Get field",
            "Fields",
            op_id="getField",
            params=[path_param("id"), path_param("fieldId")],
            responses={
                **ok("#/components/schemas/FieldResponse"),
                **err_refs("Unauthorized", "Forbidden", "NotFound"),
            },
        ),
        "put": op(
            "Update field",
            "Fields",
            op_id="updateField",
            params=[path_param("id"), path_param("fieldId")],
            req=json_body("#/components/schemas/UpdateFieldRequest"),
            responses={
                **ok("#/components/schemas/FieldResponse"),
                **err_refs("Unauthorized", "Forbidden", "NotFound", "BadRequest"),
            },
        ),
        "delete": op(
            "Delete field",
            "Fields",
            op_id="deleteField",
            params=[path_param("id"), path_param("fieldId")],
            responses={
                "204": {"description": "Deleted"},
                **err_refs("Unauthorized", "Forbidden", "NotFound"),
            },
        ),
    }

    paths["/modules/{id}/validation-rules"] = {
        "get": op(
            "List validation rules",
            "Validation",
            op_id="listRules",
            params=[path_param("id")],
            responses={
                **ok("#/components/schemas/RuleList"),
                **err_refs("Unauthorized", "Forbidden"),
            },
        ),
        "post": op(
            "Create validation rule",
            "Validation",
            op_id="createRule",
            params=[path_param("id")],
            req=json_body(
                "#/components/schemas/CreateRuleRequest",
                {
                    "rule_type": "min_length",
                    "field_id": "…",
                    "params": {"value": 3},
                    "error_message": "Too short",
                    "is_active": True,
                },
            ),
            responses={
                **ok("#/components/schemas/RuleResponse", code="201"),
                **err_refs("Unauthorized", "Forbidden", "BadRequest"),
            },
        ),
    }
    paths["/modules/{id}/validation-rules/{ruleId}"] = {
        "get": op(
            "Get validation rule",
            "Validation",
            op_id="getRule",
            params=[path_param("id"), path_param("ruleId")],
            responses={
                **ok("#/components/schemas/RuleResponse"),
                **err_refs("Unauthorized", "Forbidden", "NotFound"),
            },
        ),
        "put": op(
            "Update validation rule",
            "Validation",
            op_id="updateRule",
            params=[path_param("id"), path_param("ruleId")],
            req=json_body("#/components/schemas/UpdateRuleRequest"),
            responses={
                **ok("#/components/schemas/RuleResponse"),
                **err_refs("Unauthorized", "Forbidden", "NotFound"),
            },
        ),
        "delete": op(
            "Delete validation rule",
            "Validation",
            op_id="deleteRule",
            params=[path_param("id"), path_param("ruleId")],
            responses={
                "204": {"description": "Deleted"},
                **err_refs("Unauthorized", "Forbidden", "NotFound"),
            },
        ),
    }
    paths["/modules/{id}/validation-schema"] = {
        "get": op(
            "Compiled validation schema",
            "Validation",
            op_id="validationSchema",
            params=[path_param("id")],
            responses={
                **ok("#/components/schemas/ValidationSchema"),
                **err_refs("Unauthorized", "Forbidden", "NotFound"),
            },
        )
    }
    paths["/modules/{id}/validate"] = {
        "post": op(
            "Dry-run validate payload",
            "Validation",
            op_id="validateRecord",
            params=[path_param("id")],
            req=json_body(
                "#/components/schemas/ValidateRequest",
                {"data": {"name": "Acme", "email": "bad"}},
            ),
            responses={
                **ok(
                    "#/components/schemas/ValidateResult",
                    {
                        "success": True,
                        "data": {
                            "valid": False,
                            "errors": [
                                {
                                    "field": "email",
                                    "message": "Must be a valid email",
                                }
                            ],
                        },
                    },
                ),
                **err_refs("Unauthorized", "Forbidden", "NotFound", "BadRequest"),
            },
        )
    }

    paths["/modules/{id}/views"] = {
        "get": op(
            "List saved views",
            "Views",
            op_id="listViews",
            params=[path_param("id")],
            responses={
                **ok("#/components/schemas/ViewList"),
                **err_refs("Unauthorized", "Forbidden"),
            },
        ),
        "post": op(
            "Create saved view",
            "Views",
            op_id="createView",
            params=[path_param("id")],
            req=json_body("#/components/schemas/CreateViewRequest"),
            responses={
                **ok("#/components/schemas/ViewResponse", code="201"),
                **err_refs("Unauthorized", "Forbidden", "BadRequest"),
            },
        ),
    }
    paths["/modules/{id}/views/{viewId}"] = {
        "get": op(
            "Get saved view",
            "Views",
            op_id="getView",
            params=[path_param("id"), path_param("viewId")],
            responses={
                **ok("#/components/schemas/ViewResponse"),
                **err_refs("Unauthorized", "NotFound"),
            },
        ),
        "put": op(
            "Update saved view",
            "Views",
            op_id="updateView",
            params=[path_param("id"), path_param("viewId")],
            req=json_body("#/components/schemas/UpdateViewRequest"),
            responses={
                **ok("#/components/schemas/ViewResponse"),
                **err_refs("Unauthorized", "NotFound", "BadRequest"),
            },
        ),
        "delete": op(
            "Delete saved view",
            "Views",
            op_id="deleteView",
            params=[path_param("id"), path_param("viewId")],
            responses={
                "204": {"description": "Deleted"},
                **err_refs("Unauthorized", "NotFound"),
            },
        ),
    }
    paths["/modules/{id}/views/{viewId}/default"] = {
        "post": op(
            "Set default view",
            "Views",
            op_id="setDefaultView",
            params=[path_param("id"), path_param("viewId")],
            responses={**ok(), **err_refs("Unauthorized", "NotFound")},
        )
    }

    record_params = [
        path_param("id"),
        q("page", "integer", example=1),
        q("page_size", "integer", example=20),
        q("search"),
        q("sort"),
        q("order", example="desc"),
        q("expand", "boolean"),
        q("filters", desc="JSON-encoded FilterClause[]"),
    ]
    paths["/modules/{id}/records"] = {
        "get": op(
            "List records",
            "Records",
            op_id="listRecords",
            params=record_params,
            responses={
                **ok("#/components/schemas/RecordListResult"),
                **err_refs("Unauthorized", "Forbidden", "NotFound"),
            },
        ),
        "post": op(
            "Create record",
            "Records",
            op_id="createRecord",
            params=[path_param("id")],
            req=json_body(
                "#/components/schemas/CreateRecordRequest",
                {"data": {"name": "Acme Inc", "amount": 12000}},
            ),
            responses={
                **ok("#/components/schemas/RecordResponse", code="201"),
                **err_refs(
                    "Unauthorized", "Forbidden", "BadRequest", "ValidationError"
                ),
            },
        ),
    }
    paths["/modules/{id}/records/{recordId}"] = {
        "get": op(
            "Get record",
            "Records",
            op_id="getRecord",
            params=[path_param("id"), path_param("recordId"), q("expand", "boolean")],
            responses={
                **ok("#/components/schemas/RecordResponse"),
                **err_refs("Unauthorized", "Forbidden", "NotFound"),
            },
        ),
        "put": op(
            "Update record",
            "Records",
            op_id="updateRecord",
            params=[path_param("id"), path_param("recordId")],
            req=json_body("#/components/schemas/UpdateRecordRequest"),
            responses={
                **ok("#/components/schemas/RecordResponse"),
                **err_refs(
                    "Unauthorized", "Forbidden", "NotFound", "ValidationError"
                ),
            },
        ),
        "delete": op(
            "Delete record",
            "Records",
            op_id="deleteRecord",
            params=[path_param("id"), path_param("recordId")],
            responses={
                "204": {"description": "Deleted"},
                **err_refs("Unauthorized", "Forbidden", "NotFound"),
            },
        ),
    }

    paths["/modules/{id}/imports/analyze"] = {
        "post": op(
            "Analyze import file",
            "Import",
            op_id="analyzeImport",
            params=[path_param("id")],
            req=multipart_body(
                [("file", {"type": "string", "format": "binary"}, True)]
            ),
            responses={
                **ok("#/components/schemas/AnalyzeResult"),
                **err_refs("Unauthorized", "Forbidden", "BadRequest"),
            },
        )
    }
    paths["/modules/{id}/imports"] = {
        "get": op(
            "List import jobs",
            "Import",
            op_id="listImports",
            params=[
                path_param("id"),
                q("page", "integer"),
                q("page_size", "integer"),
                q("status"),
            ],
            responses={
                **ok("#/components/schemas/ImportListResult"),
                **err_refs("Unauthorized", "Forbidden"),
            },
        ),
        "post": op(
            "Create import job",
            "Import",
            op_id="createImport",
            params=[path_param("id")],
            description="Multipart: file + mapping (JSON string) + options (JSON string).",
            req=multipart_body(
                [
                    ("file", {"type": "string", "format": "binary"}, True),
                    (
                        "mapping",
                        {
                            "type": "string",
                            "description": "JSON object header→api_name",
                            "example": '{"Name":"name"}',
                        },
                        True,
                    ),
                    (
                        "options",
                        {"type": "string", "description": "JSON options object"},
                        False,
                    ),
                ]
            ),
            responses={
                **ok("#/components/schemas/ImportResponse", code="201"),
                **err_refs("Unauthorized", "Forbidden", "BadRequest"),
            },
        ),
    }
    paths["/modules/{id}/imports/{importId}"] = {
        "get": op(
            "Get import job",
            "Import",
            op_id="getImport",
            params=[path_param("id"), path_param("importId")],
            responses={
                **ok("#/components/schemas/ImportResponse"),
                **err_refs("Unauthorized", "Forbidden", "NotFound"),
            },
        )
    }

    paths["/modules/{id}/export"] = {
        "get": op(
            "Synchronous export download",
            "Export",
            op_id="exportNow",
            description="Builds and streams a file immediately (small datasets).",
            params=[
                path_param("id"),
                q("format", example="csv"),
                q("search"),
                q("sort"),
                q("order"),
                q("expand", "boolean"),
                q("columns", desc="Comma-separated api_names"),
                q("filters", desc="JSON FilterClause[]"),
            ],
            responses={
                "200": {
                    "description": "File bytes",
                    "content": {
                        "text/csv": {
                            "schema": {"type": "string", "format": "binary"}
                        },
                        "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": {
                            "schema": {"type": "string", "format": "binary"}
                        },
                    },
                },
                **err_refs("Unauthorized", "Forbidden", "BadRequest", "NotFound"),
            },
        )
    }
    paths["/modules/{id}/exports"] = {
        "get": op(
            "List export jobs",
            "Export",
            op_id="listExports",
            params=[
                path_param("id"),
                q("page", "integer"),
                q("page_size", "integer"),
                q("status"),
            ],
            responses={
                **ok("#/components/schemas/ExportListResult"),
                **err_refs("Unauthorized", "Forbidden"),
            },
        ),
        "post": op(
            "Create async export job",
            "Export",
            op_id="createExport",
            params=[path_param("id")],
            req=json_body(
                "#/components/schemas/ExportSpec",
                {
                    "format": "xlsx",
                    "columns": ["name", "amount"],
                    "search": "",
                    "expand": False,
                },
            ),
            responses={
                **ok("#/components/schemas/ExportResponse", code="201"),
                **err_refs("Unauthorized", "Forbidden", "BadRequest"),
            },
        ),
    }
    paths["/modules/{id}/exports/{exportId}"] = {
        "get": op(
            "Get export job",
            "Export",
            op_id="getExport",
            params=[path_param("id"), path_param("exportId")],
            responses={
                **ok("#/components/schemas/ExportResponse"),
                **err_refs("Unauthorized", "Forbidden", "NotFound"),
            },
        )
    }
    paths["/modules/{id}/exports/{exportId}/download"] = {
        "get": op(
            "Download export file",
            "Export",
            op_id="downloadExport",
            params=[path_param("id"), path_param("exportId")],
            responses={
                "200": {
                    "description": "Generated file",
                    "content": {
                        "application/octet-stream": {
                            "schema": {"type": "string", "format": "binary"}
                        }
                    },
                },
                **err_refs("Unauthorized", "Forbidden", "NotFound", "Conflict"),
            },
        )
    }
    paths["/modules/{id}/export-templates"] = {
        "get": op(
            "List export templates",
            "Export",
            op_id="listExportTemplates",
            params=[path_param("id")],
            responses={
                **ok("#/components/schemas/TemplateList"),
                **err_refs("Unauthorized", "Forbidden"),
            },
        ),
        "post": op(
            "Create export template",
            "Export",
            op_id="createExportTemplate",
            params=[path_param("id")],
            req=json_body("#/components/schemas/CreateTemplateRequest"),
            responses={
                **ok("#/components/schemas/TemplateResponse", code="201"),
                **err_refs("Unauthorized", "Forbidden", "BadRequest"),
            },
        ),
    }
    paths["/modules/{id}/export-templates/{templateId}"] = {
        "put": op(
            "Update export template",
            "Export",
            op_id="updateExportTemplate",
            params=[path_param("id"), path_param("templateId")],
            req=json_body("#/components/schemas/UpdateTemplateRequest"),
            responses={
                **ok("#/components/schemas/TemplateResponse"),
                **err_refs("Unauthorized", "Forbidden", "NotFound"),
            },
        ),
        "delete": op(
            "Delete export template",
            "Export",
            op_id="deleteExportTemplate",
            params=[path_param("id"), path_param("templateId")],
            responses={
                "204": {"description": "Deleted"},
                **err_refs("Unauthorized", "Forbidden", "NotFound"),
            },
        ),
    }

    paths["/notifications"] = {
        "get": op(
            "List notifications",
            "Notifications",
            op_id="listNotifications",
            params=[
                q("page", "integer"),
                q("page_size", "integer"),
                q("status"),
                q("channel"),
            ],
            responses={
                **ok("#/components/schemas/NotificationListResult"),
                **err_refs("Unauthorized", "Forbidden"),
            },
        ),
        "post": op(
            "Send notification",
            "Notifications",
            op_id="sendNotification",
            req=json_body(
                "#/components/schemas/SendNotificationRequest",
                {
                    "channel": "email",
                    "to": "ada@example.com",
                    "subject": "Welcome",
                    "body": "Thanks for joining",
                    "template": "lead_welcome",
                },
            ),
            responses={
                **ok("#/components/schemas/NotificationResponse", code="201"),
                **err_refs("Unauthorized", "Forbidden", "BadRequest"),
            },
        ),
    }
    paths["/notifications/{id}"] = {
        "get": op(
            "Get notification",
            "Notifications",
            op_id="getNotification",
            params=[path_param("id")],
            responses={
                **ok("#/components/schemas/NotificationResponse"),
                **err_refs("Unauthorized", "Forbidden", "NotFound"),
            },
        )
    }

    paths["/tour"] = {
        "get": op(
            "Get tour progress",
            "Tour",
            op_id="getTour",
            params=[q("tour_key", example="crm_onboarding")],
            responses={
                **ok("#/components/schemas/ProgressResponse"),
                **err_refs("Unauthorized"),
            },
        ),
        "put": op(
            "Update tour progress",
            "Tour",
            op_id="updateTour",
            req=json_body(
                "#/components/schemas/UpdateProgressRequest",
                {
                    "tour_key": "crm_onboarding",
                    "status": "active",
                    "current_step": 2,
                    "completed_steps": ["welcome", "forms"],
                },
            ),
            responses={
                **ok("#/components/schemas/ProgressResponse"),
                **err_refs("Unauthorized", "BadRequest"),
            },
        ),
    }
    paths["/tour/restart"] = {
        "post": op(
            "Restart tour",
            "Tour",
            op_id="restartTour",
            req=json_body(
                "#/components/schemas/RestartRequest",
                {"tour_key": "crm_onboarding"},
            ),
            responses={
                **ok("#/components/schemas/ProgressResponse"),
                **err_refs("Unauthorized", "BadRequest"),
            },
        )
    }

    paths["/me/organizations"] = {
        "get": op(
            "List my workspaces",
            "Organizations",
            op_id="listMyOrganizations",
            description="Returns every non-deleted organization the caller belongs to.",
            responses={
                **ok("#/components/schemas/OrgSummaryList"),
                **err_refs("Unauthorized"),
            },
        ),
    }
    paths["/me/organizations/switch"] = {
        "post": op(
            "Switch active workspace",
            "Organizations",
            op_id="switchOrganization",
            description="Updates users.active_organization_id. JWT stays user-only.",
            req=json_body("#/components/schemas/SwitchOrgRequest"),
            responses={
                **ok(),
                **err_refs("Unauthorized", "Forbidden", "BadRequest"),
            },
        )
    }
    paths["/organizations"] = {
        "post": op(
            "Create workspace",
            "Organizations",
            op_id="createOrganization",
            description="Bootstraps roles, full module catalog, and owner membership.",
            req=json_body("#/components/schemas/CreateOrgRequest"),
            responses={
                **ok(None, {"id": "11111111-1111-1111-1111-111111111111"}, code="201"),
                **err_refs("Unauthorized", "BadRequest"),
            },
        )
    }
    paths["/organizations/current"] = {
        "get": op(
            "Get current workspace",
            "Organizations",
            op_id="getCurrentOrganization",
            responses={
                **ok("#/components/schemas/OrgDetail"),
                **err_refs("Unauthorized", "NotFound"),
            },
        ),
        "patch": op(
            "Update current workspace",
            "Organizations",
            op_id="updateCurrentOrganization",
            description="Requires `organization.manage`.",
            req=json_body("#/components/schemas/UpdateOrgRequest"),
            responses={
                **ok("#/components/schemas/OrgDetail"),
                **err_refs("Unauthorized", "Forbidden", "BadRequest", "NotFound"),
            },
        ),
        "delete": op(
            "Soft-delete current workspace",
            "Organizations",
            op_id="deleteCurrentOrganization",
            description="Requires `organization.manage`. Blocked when it is the caller's last workspace.",
            responses={
                **ok(),
                **err_refs("Unauthorized", "Forbidden", "BadRequest", "NotFound"),
            },
        ),
    }

    paths["/settings"] = {
        "get": op(
            "Get organization settings",
            "Settings",
            op_id="getSettings",
            responses={
                **ok("#/components/schemas/SettingsResponse"),
                **err_refs("Unauthorized"),
            },
        ),
        "put": op(
            "Update organization settings",
            "Settings",
            op_id="updateSettings",
            description="Requires `settings.manage`.",
            req=json_body("#/components/schemas/UpdateSettingsRequest"),
            responses={
                **ok("#/components/schemas/SettingsResponse"),
                **err_refs("Unauthorized", "Forbidden", "BadRequest"),
            },
        ),
    }

    paths["/me/access"] = {
        "get": op(
            "Caller's effective RBAC access",
            "Roles",
            op_id="meAccess",
            responses={
                **ok("#/components/schemas/MeResponse"),
                **err_refs("Unauthorized"),
            },
        )
    }
    paths["/permissions"] = {
        "get": op(
            "Permission catalog",
            "Roles",
            op_id="listPermissions",
            responses={
                **ok("#/components/schemas/PermissionList"),
                **err_refs("Unauthorized", "Forbidden"),
            },
        )
    }
    paths["/roles"] = {
        "get": op(
            "List roles",
            "Roles",
            op_id="listRoles",
            responses={
                **ok("#/components/schemas/RoleSummaryList"),
                **err_refs("Unauthorized", "Forbidden"),
            },
        ),
        "post": op(
            "Create role",
            "Roles",
            op_id="createRole",
            req=json_body(
                "#/components/schemas/CreateRoleRequest",
                {"name": "Custom", "slug": "custom_role"},
            ),
            responses={
                **ok("#/components/schemas/RoleDetail", code="201"),
                **err_refs(
                    "Unauthorized", "Forbidden", "BadRequest", "Conflict"
                ),
            },
        ),
    }
    paths["/roles/{id}"] = {
        "get": op(
            "Get role",
            "Roles",
            op_id="getRole",
            params=[path_param("id")],
            responses={
                **ok("#/components/schemas/RoleDetail"),
                **err_refs("Unauthorized", "Forbidden", "NotFound"),
            },
        ),
        "put": op(
            "Update role",
            "Roles",
            op_id="updateRole",
            params=[path_param("id")],
            req=json_body("#/components/schemas/UpdateRoleRequest"),
            responses={
                **ok("#/components/schemas/RoleDetail"),
                **err_refs("Unauthorized", "Forbidden", "NotFound"),
            },
        ),
        "delete": op(
            "Delete role",
            "Roles",
            op_id="deleteRole",
            params=[path_param("id")],
            responses={
                "204": {"description": "Deleted"},
                **err_refs("Unauthorized", "Forbidden", "NotFound", "Conflict"),
            },
        ),
    }
    paths["/roles/{id}/permissions"] = {
        "put": op(
            "Replace role permissions",
            "Roles",
            op_id="setRolePermissions",
            params=[path_param("id")],
            req=json_body(
                "#/components/schemas/SetPermissionsRequest",
                {"permissions": ["module.view", "record.view", "export.run"]},
            ),
            responses={
                **ok("#/components/schemas/RoleDetail"),
                **err_refs("Unauthorized", "Forbidden", "NotFound", "BadRequest"),
            },
        )
    }
    paths["/roles/{id}/module-access"] = {
        "put": op(
            "Replace module ACL",
            "Roles",
            op_id="setModuleAccess",
            params=[path_param("id")],
            req=json_body("#/components/schemas/SetModuleAccessRequest"),
            responses={
                **ok("#/components/schemas/RoleDetail"),
                **err_refs("Unauthorized", "Forbidden", "NotFound", "BadRequest"),
            },
        )
    }
    paths["/roles/{id}/field-access"] = {
        "put": op(
            "Replace field ACL",
            "Roles",
            op_id="setFieldAccess",
            params=[path_param("id")],
            req=json_body("#/components/schemas/SetFieldAccessRequest"),
            responses={
                **ok("#/components/schemas/RoleDetail"),
                **err_refs("Unauthorized", "Forbidden", "NotFound", "BadRequest"),
            },
        )
    }

    paths["/openapi.yaml"] = {
        "get": op(
            "OpenAPI 3 specification",
            "Docs",
            auth=False,
            op_id="openapiSpec",
            responses={
                "200": {
                    "description": "YAML document",
                    "content": {
                        "application/yaml": {"schema": {"type": "string"}}
                    },
                }
            },
        )
    }
    paths["/docs"] = {
        "get": op(
            "Swagger UI",
            "Docs",
            auth=False,
            op_id="swaggerUI",
            responses={
                "200": {
                    "description": "HTML Swagger UI",
                    "content": {"text/html": {"schema": {"type": "string"}}},
                }
            },
        )
    }
    return paths


def main():
    from schemas_data import SCHEMAS

    doc = {
        "openapi": "3.0.3",
        "info": {
            "title": "CRM Lite API",
            "version": "2.0.0",
            "description": textwrap.dedent(
                """\
                Portfolio CRM REST API (Go/Gin + PostgreSQL + Redis).

                All JSON endpoints return the shared envelope:

                ```json
                { "success": true|false, "message": "...", "data": {}, "errors": [] }
                ```

                Authenticate with `Authorization: Bearer <access_token>` from `POST /auth/login`.
                Organization scope is resolved server-side from the token (tenant middleware).
                """
            ),
        },
        # Placeholder is rewritten per request in internal/docs/docs.go so
        # Swagger "Try it out" always targets the host the client opened.
        "servers": [
            {
                "url": "__API_SERVER_URL__",
                "description": "Current host (injected at request time)",
            }
        ],
        "tags": [
            {"name": n}
            for n in [
                "Health",
                "Auth",
                "Docs",
                "Dashboard",
                "Search",
                "Attachments",
                "Activities",
                "Media",
                "Modules",
                "Fields",
                "Validation",
                "Views",
                "Records",
                "Import",
                "Export",
                "Notifications",
                "Tour",
                "Organizations",
                "Settings",
                "Roles",
            ]
        ],
        "paths": build_paths(),
        "components": {
            "securitySchemes": {
                "bearerAuth": {
                    "type": "http",
                    "scheme": "bearer",
                    "bearerFormat": "JWT",
                    "description": "JWT from POST /auth/login",
                }
            },
            "responses": error_responses(),
            "schemas": SCHEMAS,
        },
    }

    out = Path(__file__).with_name("openapi.yaml")
    with out.open("w", encoding="utf-8") as f:
        f.write(
            "# Generated by generate_openapi.py — edit schemas_data.py / this script, then regenerate.\n"
        )
        yaml.dump(
            doc,
            f,
            sort_keys=False,
            allow_unicode=True,
            default_flow_style=False,
            width=100,
        )
    print(f"Wrote {out} ({out.stat().st_size} bytes)")


def error_responses():
    return {
        "BadRequest": {
            "description": "Malformed request or business validation failure",
            "content": {
                "application/json": {
                    "schema": {"$ref": "#/components/schemas/ErrorResponse"},
                    "examples": {
                        "simple": {
                            "value": {
                                "success": False,
                                "message": "Invalid request",
                                "errors": None,
                            }
                        },
                        "fieldErrors": {
                            "value": {
                                "success": False,
                                "message": "Validation failed",
                                "errors": [
                                    {
                                        "field": "email",
                                        "message": "Must be a valid email",
                                    }
                                ],
                            }
                        },
                        "appError": {
                            "value": {
                                "success": False,
                                "message": "Invalid request",
                                "errors": [{"code": "BAD_REQUEST"}],
                            }
                        },
                    },
                }
            },
        },
        "Unauthorized": {
            "description": "Missing or invalid Bearer token",
            "content": {
                "application/json": {
                    "schema": {"$ref": "#/components/schemas/ErrorResponse"},
                    "example": {
                        "success": False,
                        "message": "Unauthorized",
                        "errors": None,
                    },
                }
            },
        },
        "Forbidden": {
            "description": "Authenticated but lacking permission / module ACL",
            "content": {
                "application/json": {
                    "schema": {"$ref": "#/components/schemas/ErrorResponse"},
                    "example": {
                        "success": False,
                        "message": "Forbidden",
                        "errors": None,
                    },
                }
            },
        },
        "NotFound": {
            "description": "Resource not found in the caller's organization",
            "content": {
                "application/json": {
                    "schema": {"$ref": "#/components/schemas/ErrorResponse"},
                    "example": {
                        "success": False,
                        "message": "Not found",
                        "errors": [{"code": "NOT_FOUND"}],
                    },
                }
            },
        },
        "Conflict": {
            "description": "Conflict (duplicate slug, system role, job not ready, etc.)",
            "content": {
                "application/json": {
                    "schema": {"$ref": "#/components/schemas/ErrorResponse"},
                    "example": {
                        "success": False,
                        "message": "Conflict",
                        "errors": [{"code": "CONFLICT"}],
                    },
                }
            },
        },
        "ValidationError": {
            "description": "Field-level validation errors (HTTP 400; record engine / import)",
            "content": {
                "application/json": {
                    "schema": {"$ref": "#/components/schemas/ErrorResponse"},
                    "example": {
                        "success": False,
                        "message": "Validation failed",
                        "errors": [
                            {"field": "name", "message": "This field is required"},
                            {
                                "field": "email",
                                "message": "Must be a valid email",
                            },
                        ],
                    },
                }
            },
        },
        "InternalError": {
            "description": "Unexpected server error",
            "content": {
                "application/json": {
                    "schema": {"$ref": "#/components/schemas/ErrorResponse"},
                    "example": {
                        "success": False,
                        "message": "Internal server error",
                        "errors": None,
                    },
                }
            },
        },
    }


if __name__ == "__main__":
    main()
