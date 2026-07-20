"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import {
  DndContext,
  DragOverlay,
  PointerSensor,
  closestCorners,
  useSensor,
  useSensors,
  type DragEndEvent,
  type DragOverEvent,
  type DragStartEvent,
  useDroppable,
} from "@dnd-kit/core";
import {
  SortableContext,
  arrayMove,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import {
  GripVertical,
  Pencil,
  Plus,
  Trash2,
} from "lucide-react";

import type { ModuleField } from "@/features/metadata/types";
import type { LayoutSection } from "@/features/workspace/types";

const SYSTEM_KEYS = new Set([
  "owner_id",
  "assigned_to",
  "visibility",
  "created_at",
  "updated_at",
]);

function sectionId(key: string) {
  return `section:${key}`;
}

function fieldId(apiName: string) {
  return `field:${apiName}`;
}

function parseSectionId(id: string): string | null {
  return id.startsWith("section:") ? id.slice("section:".length) : null;
}

function parseFieldId(id: string): string | null {
  return id.startsWith("field:") ? id.slice("field:".length) : null;
}

/** Merge layout sections with known fields; orphans land in general (or first). */
export function buildEditableSections(
  layoutSections: LayoutSection[] | undefined,
  fields: ModuleField[]
): LayoutSection[] {
  const byApi = new Map(fields.map((f) => [f.api_name, f]));
  const known = new Set(fields.map((f) => f.api_name));
  const seen = new Set<string>();

  const base =
    layoutSections && layoutSections.length > 0
      ? layoutSections.map((s) => ({
          key: s.key,
          label: s.label,
          fields: (s.fields ?? []).filter((name) => {
            if (SYSTEM_KEYS.has(name)) {
              if (seen.has(name)) return false;
              seen.add(name);
              return true;
            }
            if (!known.has(name) || seen.has(name)) return false;
            seen.add(name);
            return true;
          }),
        }))
      : [{ key: "general", label: "General Information", fields: [] as string[] }];

  let orphanTarget = base.findIndex((s) => s.key !== "system");
  if (orphanTarget < 0) orphanTarget = 0;

  for (const f of fields) {
    if (!seen.has(f.api_name)) {
      base[orphanTarget].fields.push(f.api_name);
      seen.add(f.api_name);
    }
  }

  // Drop empty references that aren't fields (keep system placeholders).
  return base.map((s) => ({
    ...s,
    fields: s.fields.filter(
      (name) => SYSTEM_KEYS.has(name) || byApi.has(name)
    ),
  }));
}

type Props = {
  fields: ModuleField[];
  sections: LayoutSection[];
  saving?: boolean;
  onChange: (sections: LayoutSection[]) => void;
  onEditField: (field: ModuleField) => void;
  onDeleteField: (field: ModuleField) => void;
};

function SortableFieldRow({
  field,
  onEdit,
  onDelete,
}: {
  field: ModuleField;
  onEdit: () => void;
  onDelete: () => void;
}) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: fieldId(field.api_name) });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.4 : 1,
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      className="flex items-center gap-2 rounded-xl border border-slate-200 bg-white px-3 py-2"
    >
      <button
        type="button"
        className="cursor-grab touch-none text-slate-300 hover:text-slate-500 active:cursor-grabbing"
        aria-label="Drag field"
        {...attributes}
        {...listeners}
      >
        <GripVertical className="h-4 w-4" />
      </button>
      <div className="min-w-0 flex-1">
        <p className="truncate text-sm font-semibold text-slate-800">
          {field.label}
        </p>
        <p className="truncate text-xs text-slate-400">
          <code>{field.api_name}</code> · {field.field_type}
          {field.is_required ? " · required" : ""}
        </p>
      </div>
      <button
        type="button"
        onClick={onEdit}
        className="rounded-lg p-1.5 text-slate-400 hover:bg-slate-100 hover:text-slate-700"
        aria-label="Edit field"
      >
        <Pencil className="h-4 w-4" />
      </button>
      <button
        type="button"
        onClick={onDelete}
        disabled={field.is_system}
        className="rounded-lg p-1.5 text-red-400 hover:bg-red-50 hover:text-red-600 disabled:opacity-30"
        aria-label="Delete field"
      >
        <Trash2 className="h-4 w-4" />
      </button>
    </div>
  );
}

function SystemFieldChip({ name }: { name: string }) {
  return (
    <div className="rounded-xl border border-dashed border-slate-200 bg-slate-50 px-3 py-2 text-xs font-medium text-slate-500">
      {name.replace(/_/g, " ")}
    </div>
  );
}

function SortableSectionCard({
  section,
  fieldsByApi,
  onRename,
  onRemove,
  onEditField,
  onDeleteField,
}: {
  section: LayoutSection;
  fieldsByApi: Map<string, ModuleField>;
  onRename: (label: string, persist?: boolean) => void;
  onRemove: () => void;
  onEditField: (field: ModuleField) => void;
  onDeleteField: (field: ModuleField) => void;
}) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: sectionId(section.key) });
  const { setNodeRef: setDropRef } = useDroppable({
    id: `drop:${section.key}`,
    disabled: section.key === "system",
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  const fieldIds = section.fields
    .filter((n) => !SYSTEM_KEYS.has(n))
    .map(fieldId);

  const canRemove =
    section.key !== "general" &&
    section.key !== "system" &&
    section.fields.every((n) => SYSTEM_KEYS.has(n) || !fieldsByApi.has(n));

  return (
    <div
      ref={setNodeRef}
      style={style}
      className="rounded-3xl border border-slate-200 bg-white p-4 shadow-sm"
    >
      <div className="mb-3 flex items-center gap-2">
        <button
          type="button"
          className="cursor-grab touch-none text-slate-300 hover:text-slate-500 active:cursor-grabbing"
          aria-label="Drag section"
          {...attributes}
          {...listeners}
        >
          <GripVertical className="h-5 w-5" />
        </button>
        <input
          value={section.label}
          onChange={(e) => onRename(e.target.value)}
          onBlur={(e) => onRename(e.target.value, true)}
          className="min-w-0 flex-1 rounded-xl border border-transparent bg-transparent px-2 py-1 text-sm font-semibold text-slate-900 hover:border-slate-200 focus:border-emerald-400 focus:outline-none"
        />
        <code className="hidden text-xs text-slate-400 sm:inline">
          {section.key}
        </code>
        {canRemove && (
          <button
            type="button"
            onClick={onRemove}
            className="rounded-lg p-1.5 text-slate-400 hover:bg-red-50 hover:text-red-500"
            aria-label="Remove section"
          >
            <Trash2 className="h-4 w-4" />
          </button>
        )}
      </div>

      <div ref={setDropRef} className="space-y-2 min-h-[44px]">
        <SortableContext items={fieldIds} strategy={verticalListSortingStrategy}>
          {section.fields.map((name) => {
            if (SYSTEM_KEYS.has(name)) {
              return <SystemFieldChip key={name} name={name} />;
            }
            const field = fieldsByApi.get(name);
            if (!field) return null;
            return (
              <SortableFieldRow
                key={name}
                field={field}
                onEdit={() => onEditField(field)}
                onDelete={() => onDeleteField(field)}
              />
            );
          })}
        </SortableContext>
        {section.fields.length === 0 && (
          <p className="rounded-xl border border-dashed border-slate-200 px-3 py-4 text-center text-xs text-slate-400">
            Drop fields here
          </p>
        )}
      </div>
    </div>
  );
}

export default function FieldSectionsEditor({
  fields,
  sections: initialSections,
  saving,
  onChange,
  onEditField,
  onDeleteField,
}: Props) {
  const [sections, setSections] = useState(initialSections);
  const [activeId, setActiveId] = useState<string | null>(null);
  const sectionsRef = useRef(sections);
  /** Last layout successfully sent to (or loaded from) the parent/API. */
  const persistedRef = useRef(initialSections);

  useEffect(() => {
    setSections(initialSections);
    sectionsRef.current = initialSections;
    persistedRef.current = initialSections;
  }, [initialSections]);

  useEffect(() => {
    sectionsRef.current = sections;
  }, [sections]);

  const fieldsByApi = useMemo(
    () => new Map(fields.map((f) => [f.api_name, f])),
    [fields]
  );

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 6 } })
  );

  const sectionIds = sections.map((s) => sectionId(s.key));
  const dragStartRef = useRef<LayoutSection[]>([]);

  function sectionsEqual(a: LayoutSection[], b: LayoutSection[]): boolean {
    if (a.length !== b.length) return false;
    return a.every(
      (s, i) =>
        s.key === b[i].key &&
        s.label === b[i].label &&
        s.fields.length === b[i].fields.length &&
        s.fields.every((f, j) => f === b[i].fields[j])
    );
  }

  function commit(next: LayoutSection[], persist = true) {
    sectionsRef.current = next;
    setSections(next);
    // Compare against last persisted snapshot — not local state — so renames
    // typed via updateLocal still save on blur.
    if (persist && !sectionsEqual(persistedRef.current, next)) {
      persistedRef.current = next;
      onChange(next);
    }
  }

  function updateLocal(next: LayoutSection[]) {
    sectionsRef.current = next;
    setSections(next);
  }

  function findContainerForField(
    apiName: string,
    list: LayoutSection[] = sectionsRef.current
  ): string | null {
    for (const s of list) {
      if (s.fields.includes(apiName)) return s.key;
    }
    return null;
  }

  function resolveSectionTarget(
    overId: string,
    list: LayoutSection[] = sectionsRef.current
  ): string | null {
    const asSection = parseSectionId(overId);
    if (asSection) return asSection;
    if (overId.startsWith("drop:")) return overId.slice("drop:".length);
    const asField = parseFieldId(overId);
    if (asField) return findContainerForField(asField, list);
    return null;
  }

  function handleDragStart(event: DragStartEvent) {
    dragStartRef.current = sectionsRef.current.map((s) => ({
      ...s,
      fields: [...s.fields],
    }));
    setActiveId(String(event.active.id));
  }

  function handleDragOver(event: DragOverEvent) {
    const { active, over } = event;
    if (!over) return;

    const activeField = parseFieldId(String(active.id));
    if (!activeField) return;

    const overField = parseFieldId(String(over.id));
    const overSection = resolveSectionTarget(String(over.id));
    if (!overSection || overSection === "system") return;

    const from = findContainerForField(activeField);
    if (!from || from === "system") return;

    // Same-section reorder: move while dragging so order updates before drop.
    if (from === overSection) {
      if (!overField || overField === activeField) return;
      const prev = sectionsRef.current;
      const sec = prev.find((s) => s.key === from);
      if (!sec) return;
      const oldIndex = sec.fields.indexOf(activeField);
      const newIndex = sec.fields.indexOf(overField);
      if (oldIndex < 0 || newIndex < 0 || oldIndex === newIndex) return;
      updateLocal(
        prev.map((s) =>
          s.key !== from
            ? s
            : { ...s, fields: arrayMove(s.fields, oldIndex, newIndex) }
        )
      );
      return;
    }

    // Cross-section move
    const prev = sectionsRef.current;
    const next = prev.map((s) => ({ ...s, fields: [...s.fields] }));
    const fromSec = next.find((s) => s.key === from);
    const toSec = next.find((s) => s.key === overSection);
    if (!fromSec || !toSec) return;
    fromSec.fields = fromSec.fields.filter((n) => n !== activeField);
    const overIndex = overField
      ? toSec.fields.indexOf(overField)
      : toSec.fields.length;
    const insertAt = overIndex >= 0 ? overIndex : toSec.fields.length;
    toSec.fields.splice(insertAt, 0, activeField);
    updateLocal(next);
  }

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event;
    setActiveId(null);
    const current = sectionsRef.current;

    if (!over) {
      // Cancel: restore snapshot from drag start.
      if (!sectionsEqual(dragStartRef.current, current)) {
        updateLocal(dragStartRef.current);
      }
      return;
    }

    const activeSec = parseSectionId(String(active.id));
    if (activeSec) {
      const overSec = resolveSectionTarget(String(over.id), current);
      if (overSec && activeSec !== overSec) {
        const oldIndex = current.findIndex((s) => s.key === activeSec);
        const newIndex = current.findIndex((s) => s.key === overSec);
        if (oldIndex >= 0 && newIndex >= 0) {
          commit(arrayMove(current, oldIndex, newIndex));
          return;
        }
      }
      // No section move — don't persist unless something else changed.
      if (!sectionsEqual(dragStartRef.current, current)) {
        commit(current);
      }
      return;
    }

    const activeField = parseFieldId(String(active.id));
    if (!activeField) return;

    // Persist whatever local order dragOver already applied (same or cross section).
    if (!sectionsEqual(dragStartRef.current, current)) {
      commit(current);
      return;
    }

    // Fallback: drop landed on another field but dragOver never ran (rare).
    const overField = parseFieldId(String(over.id));
    if (overField && activeField !== overField) {
      const container = findContainerForField(activeField, current);
      const overContainer = findContainerForField(overField, current);
      if (container && container === overContainer && container !== "system") {
        const next = current.map((s) => {
          if (s.key !== container) return s;
          const oldIndex = s.fields.indexOf(activeField);
          const newIndex = s.fields.indexOf(overField);
          if (oldIndex < 0 || newIndex < 0) return s;
          return {
            ...s,
            fields: arrayMove(s.fields, oldIndex, newIndex),
          };
        });
        commit(next);
      }
    }
  }

  function addSection() {
    const list = sectionsRef.current;
    const base = `section_${list.length + 1}`;
    let key = base;
    let n = 1;
    while (list.some((s) => s.key === key)) {
      n += 1;
      key = `${base}_${n}`;
    }
    commit([
      ...list,
      { key, label: `Section ${list.length + 1}`, fields: [] },
    ]);
  }

  const activeFieldName = activeId ? parseFieldId(activeId) : null;
  const activeField = activeFieldName
    ? fieldsByApi.get(activeFieldName)
    : undefined;

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <div>
          <h3 className="text-sm font-semibold text-slate-800">Field sections</h3>
          <p className="text-xs text-slate-500">
            Drag fields within a section (or between sections) to set form and
            overview order.
            {saving ? " Saving…" : ""}
          </p>
        </div>
        <button
          type="button"
          onClick={addSection}
          className="inline-flex items-center gap-1.5 rounded-full border border-slate-200 px-3 py-1.5 text-xs font-semibold text-slate-700 hover:bg-slate-50"
        >
          <Plus className="h-3.5 w-3.5" />
          Add section
        </button>
      </div>

      <DndContext
        sensors={sensors}
        collisionDetection={closestCorners}
        onDragStart={handleDragStart}
        onDragOver={handleDragOver}
        onDragEnd={handleDragEnd}
      >
        <SortableContext items={sectionIds} strategy={verticalListSortingStrategy}>
          <div className="space-y-4">
            {sections.map((section) => (
              <SortableSectionCard
                key={section.key}
                section={section}
                fieldsByApi={fieldsByApi}
                onRename={(label, persist) => {
                  const trimmed = label.trim();
                  if (persist) {
                    // Reject empty names; restore last persisted label for this section.
                    const fallback =
                      persistedRef.current.find((s) => s.key === section.key)
                        ?.label ?? section.label;
                    const nextLabel = trimmed || fallback;
                    const next = sectionsRef.current.map((s) =>
                      s.key === section.key ? { ...s, label: nextLabel } : s
                    );
                    commit(next);
                    return;
                  }
                  const next = sectionsRef.current.map((s) =>
                    s.key === section.key ? { ...s, label } : s
                  );
                  updateLocal(next);
                }}
                onRemove={() =>
                  commit(
                    sectionsRef.current.filter((s) => s.key !== section.key)
                  )
                }
                onEditField={onEditField}
                onDeleteField={onDeleteField}
              />
            ))}
          </div>
        </SortableContext>

        <DragOverlay>
          {activeField ? (
            <div className="rounded-xl border border-emerald-200 bg-white px-3 py-2 text-sm font-semibold shadow-lg">
              {activeField.label}
            </div>
          ) : null}
        </DragOverlay>
      </DndContext>
    </div>
  );
}
