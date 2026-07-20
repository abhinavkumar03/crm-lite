"use client";

import { useRef, useState } from "react";
import { ImagePlus, Loader2, Trash2 } from "lucide-react";
import { toast } from "sonner";

import { uploadFile } from "@/features/uploads/api";

const MAX_BYTES = 2 * 1024 * 1024;
const ACCEPT = "image/png,image/jpeg,image/webp,image/svg+xml,image/gif";

type Props = {
  value: string;
  orgName?: string;
  disabled?: boolean;
  onChange: (url: string) => void;
};

/** Upload / preview / clear organization logo (stores Cloudinary URL). */
export default function OrgLogoUploader({
  value,
  orgName,
  disabled,
  onChange,
}: Props) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [uploading, setUploading] = useState(false);

  async function handleFile(file: File | undefined) {
    if (!file) return;
    if (!file.type.startsWith("image/")) {
      toast.error("Please choose an image file");
      return;
    }
    if (file.size > MAX_BYTES) {
      toast.error("Logo must be 2MB or smaller");
      return;
    }
    try {
      setUploading(true);
      const uploaded = await uploadFile(file);
      onChange(uploaded.url);
      toast.success("Logo uploaded — save settings to apply");
    } catch {
      toast.error("Failed to upload logo");
    } finally {
      setUploading(false);
      if (inputRef.current) inputRef.current.value = "";
    }
  }

  const initial = (orgName?.trim()?.[0] ?? "C").toUpperCase();

  return (
    <div className="space-y-2">
      <p className="text-sm font-semibold text-slate-700">Organization logo</p>
      <div className="flex flex-wrap items-center gap-4">
        <div className="flex h-16 w-16 shrink-0 items-center justify-center overflow-hidden rounded-2xl border border-slate-200 bg-slate-50 shadow-sm">
          {value ? (
            // eslint-disable-next-line @next/next/no-img-element
            <img
              src={value}
              alt={orgName ? `${orgName} logo` : "Organization logo"}
              className="h-full w-full object-contain p-1"
            />
          ) : (
            <span className="text-xl font-bold text-emerald-600">{initial}</span>
          )}
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <input
            ref={inputRef}
            type="file"
            accept={ACCEPT}
            className="hidden"
            disabled={disabled || uploading}
            onChange={(e) => handleFile(e.target.files?.[0])}
          />
          <button
            type="button"
            disabled={disabled || uploading}
            onClick={() => inputRef.current?.click()}
            className="inline-flex items-center gap-2 rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm font-semibold text-slate-700 hover:bg-slate-50 disabled:opacity-50"
          >
            {uploading ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <ImagePlus className="h-4 w-4" />
            )}
            {uploading ? "Uploading…" : value ? "Replace logo" : "Upload logo"}
          </button>
          {value ? (
            <button
              type="button"
              disabled={disabled || uploading}
              onClick={() => onChange("")}
              className="inline-flex items-center gap-2 rounded-xl border border-red-200 bg-white px-3 py-2 text-sm font-semibold text-red-600 hover:bg-red-50 disabled:opacity-50"
            >
              <Trash2 className="h-4 w-4" />
              Remove
            </button>
          ) : null}
        </div>
      </div>
      <p className="text-xs text-slate-400">
        PNG, JPG, WebP, or SVG up to 2MB. Shown in the sidebar and workspace
        switcher after you save.
      </p>
    </div>
  );
}

/** Custom event so chrome (sidebar) refreshes after settings save. */
export const ORG_BRANDING_EVENT = "crm:org-branding-updated";

export type OrgBrandingDetail = {
  name: string;
  logo_url: string | null;
};

export function notifyOrgBrandingUpdated(detail: OrgBrandingDetail) {
  if (typeof window === "undefined") return;
  window.dispatchEvent(
    new CustomEvent(ORG_BRANDING_EVENT, { detail })
  );
}
