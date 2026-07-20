"use client";

import {
  DndContext,
  PointerSensor,
  closestCenter,
  useSensor,
  useSensors,
  type DragEndEvent,
} from "@dnd-kit/core";
import {
  SortableContext,
  arrayMove,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { GripVertical, Eye, EyeOff, Lock } from "lucide-react";

import type { ListColumn } from "@/features/workspace/types";

type Props = {
  columns: ListColumn[];
  saving?: boolean;
  onReorder: (next: ListColumn[]) => void;
  onToggle: (col: ListColumn) => void;
};

function columnLabel(col: ListColumn): string {
  if (col.label) return col.label;
  if (col.field_key === "_actions") return "Actions";
  return col.field_key;
}

function SortableColumnRow({
  col,
  saving,
  onToggle,
}: {
  col: ListColumn;
  saving?: boolean;
  onToggle: () => void;
}) {
  const isActions = col.field_key === "_actions" || col.system;
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: col.field_key,
    disabled: isActions || saving,
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.45 : 1,
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      className="flex items-center gap-3 border-b border-slate-100 px-4 py-3 last:border-0"
    >
      <button
        type="button"
        className="cursor-grab touch-none text-slate-300 hover:text-slate-500 active:cursor-grabbing disabled:cursor-not-allowed disabled:opacity-30"
        aria-label={`Drag ${columnLabel(col)}`}
        disabled={isActions || saving}
        {...attributes}
        {...listeners}
      >
        <GripVertical className="h-4 w-4" />
      </button>

      <div className="min-w-0 flex-1">
        <span className="font-medium text-slate-800">{columnLabel(col)}</span>
        {col.locked && (
          <span className="ml-2 inline-flex items-center gap-1 rounded bg-amber-50 px-1.5 py-0.5 text-[10px] font-semibold uppercase text-amber-700">
            <Lock className="h-2.5 w-2.5" />
            Locked
          </span>
        )}
        {col.system && !col.locked && (
          <span className="ml-2 rounded bg-slate-100 px-1.5 py-0.5 text-[10px] font-semibold uppercase text-slate-500">
            System
          </span>
        )}
        {col.field_key !== columnLabel(col) && (
          <span className="ml-2 text-xs text-slate-400">{col.field_key}</span>
        )}
      </div>

      <button
        type="button"
        disabled={col.locked || col.system || saving}
        onClick={onToggle}
        className="inline-flex items-center gap-1.5 rounded-lg border border-slate-200 px-2.5 py-1.5 text-xs font-semibold text-slate-600 hover:bg-slate-50 disabled:opacity-40"
      >
        {col.visible ? (
          <Eye className="h-3.5 w-3.5" />
        ) : (
          <EyeOff className="h-3.5 w-3.5" />
        )}
        {col.visible ? "Visible" : "Hidden"}
      </button>
    </div>
  );
}

/** Drag-and-drop editor for org-default listing columns. */
export default function ListingColumnsEditor({
  columns,
  saving,
  onReorder,
  onToggle,
}: Props) {
  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 6 } })
  );

  const sortableIds = columns
    .filter((c) => !c.system && c.field_key !== "_actions")
    .map((c) => c.field_key);

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event;
    if (!over || active.id === over.id) return;

    const from = columns.findIndex((c) => c.field_key === active.id);
    const to = columns.findIndex((c) => c.field_key === over.id);
    if (from < 0 || to < 0) return;

    const target = columns[to];
    if (target.system || target.field_key === "_actions") return;

    onReorder(arrayMove(columns, from, to));
  }

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      onDragEnd={handleDragEnd}
    >
      <SortableContext items={sortableIds} strategy={verticalListSortingStrategy}>
        <div className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
          <div className="flex items-center gap-3 border-b border-slate-100 bg-slate-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-slate-500">
            <span className="w-4" />
            <span className="flex-1">Column</span>
            <span>Visibility</span>
          </div>
          {columns.map((col) => (
            <SortableColumnRow
              key={col.field_key}
              col={col}
              saving={saving}
              onToggle={() => onToggle(col)}
            />
          ))}
        </div>
      </SortableContext>
    </DndContext>
  );
}
