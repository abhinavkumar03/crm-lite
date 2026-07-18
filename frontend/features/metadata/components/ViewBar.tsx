import { useState } from "react";
import { Plus, Star, Trash2 } from "lucide-react";

import { SavedView } from "../types";

interface Props {
  views: SavedView[];
  activeViewId: string | null;
  onApply: (view: SavedView) => void;
  onSaveCurrent: (name: string, isPublic: boolean) => void | Promise<void>;
  onSetDefault: (view: SavedView) => void | Promise<void>;
  onDelete: (view: SavedView) => void | Promise<void>;
}

// ViewBar lets users switch between saved table configurations, save the current
// layout as a new view, mark a default and remove views they own.
export default function ViewBar({
  views,
  activeViewId,
  onApply,
  onSaveCurrent,
  onSetDefault,
  onDelete,
}: Props) {
  const [saving, setSaving] = useState(false);
  const [name, setName] = useState("");
  const [isPublic, setIsPublic] = useState(true);

  async function handleSave() {
    const trimmed = name.trim();
    if (!trimmed) return;
    await onSaveCurrent(trimmed, isPublic);
    setName("");
    setSaving(false);
  }

  return (
    <div className="flex flex-wrap items-center gap-2">
      {views.map((view) => {
        const active = view.id === activeViewId;
        return (
          <div
            key={view.id}
            className={`group flex items-center gap-1 rounded-full border px-3 py-1.5 text-sm transition ${
              active
                ? "border-emerald-400 bg-emerald-50 text-emerald-700"
                : "border-slate-200 bg-white text-slate-600 hover:border-slate-300"
            }`}
          >
            <button
              type="button"
              onClick={() => onApply(view)}
              className="flex items-center gap-1.5 font-medium"
            >
              <Star
                className={`h-3.5 w-3.5 ${
                  view.is_default
                    ? "fill-amber-400 text-amber-400"
                    : "text-slate-300"
                }`}
                onClick={(e) => {
                  e.stopPropagation();
                  if (!view.is_default) void onSetDefault(view);
                }}
              />
              {view.name}
              {!view.is_public && (
                <span className="text-[10px] uppercase text-slate-400">
                  private
                </span>
              )}
            </button>
            {view.is_owner && (
              <button
                type="button"
                onClick={() => onDelete(view)}
                className="text-slate-300 opacity-0 transition hover:text-red-500 group-hover:opacity-100"
                aria-label={`Delete ${view.name}`}
              >
                <Trash2 className="h-3.5 w-3.5" />
              </button>
            )}
          </div>
        );
      })}

      {saving ? (
        <div className="flex items-center gap-2 rounded-full border border-slate-200 bg-white px-3 py-1">
          <input
            autoFocus
            value={name}
            onChange={(e) => setName(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") void handleSave();
              if (e.key === "Escape") setSaving(false);
            }}
            placeholder="View name"
            className="w-32 text-sm text-slate-700 focus:outline-none"
          />
          <label className="flex items-center gap-1 text-[11px] text-slate-500">
            <input
              type="checkbox"
              checked={isPublic}
              onChange={(e) => setIsPublic(e.target.checked)}
              className="h-3.5 w-3.5 accent-emerald-500"
            />
            shared
          </label>
          <button
            type="button"
            onClick={handleSave}
            className="rounded-full bg-emerald-500 px-2.5 py-0.5 text-xs font-semibold text-white hover:bg-emerald-600"
          >
            Save
          </button>
        </div>
      ) : (
        <button
          type="button"
          onClick={() => setSaving(true)}
          className="flex items-center gap-1 rounded-full border border-dashed border-slate-300 px-3 py-1.5 text-sm text-slate-500 hover:border-emerald-400 hover:text-emerald-600"
        >
          <Plus className="h-3.5 w-3.5" />
          Save view
        </button>
      )}
    </div>
  );
}
