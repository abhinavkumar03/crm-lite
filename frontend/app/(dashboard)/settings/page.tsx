"use client";

import { useEffect, useState } from "react";
import { toast } from "sonner";
import { Save } from "lucide-react";

import FormInput from "@/components/common/form/FormInput";
import FormSelect from "@/components/common/form/FormSelect";

import { getSettings, updateSettings } from "@/features/settings/api";
import OrgLogoUploader, {
  notifyOrgBrandingUpdated,
} from "@/features/settings/components/OrgLogoUploader";
import { GeneralSettings, OrgSettings } from "@/features/settings/types";

const TIMEZONES = [
  "UTC",
  "America/New_York",
  "America/Chicago",
  "America/Los_Angeles",
  "Europe/London",
  "Europe/Berlin",
  "Asia/Kolkata",
  "Asia/Singapore",
  "Asia/Tokyo",
  "Australia/Sydney",
];

const DATE_FORMATS = ["YYYY-MM-DD", "DD/MM/YYYY", "MM/DD/YYYY", "DD MMM YYYY"];

const CURRENCIES = ["USD", "EUR", "GBP", "INR", "AUD", "CAD", "SGD", "JPY"];

const LOCALES = ["en-US", "en-GB", "en-IN", "de-DE", "fr-FR", "es-ES", "ja-JP"];

const COMPANY_SIZES = ["1-10", "11-50", "51-200", "201-500", "500+"];

export default function GeneralSettingsPage() {
  const [settings, setSettings] = useState<OrgSettings | null>(null);
  const [name, setName] = useState("");
  const [logoUrl, setLogoUrl] = useState("");
  const [industry, setIndustry] = useState("");
  const [companySize, setCompanySize] = useState("");
  const [country, setCountry] = useState("");
  const [general, setGeneral] = useState<GeneralSettings | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const data = await getSettings();
        if (!active) return;
        setSettings(data);
        setName(data.name);
        setLogoUrl(data.logo_url ?? "");
        setIndustry(data.industry ?? "");
        setCompanySize(data.company_size ?? "");
        setCountry(data.country ?? "");
        setGeneral(data.general);
      } catch {
        toast.error("Failed to load settings");
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, []);

  function patchGeneral(patch: Partial<GeneralSettings>) {
    setGeneral((prev) => (prev ? { ...prev, ...patch } : prev));
  }

  async function handleSave() {
    if (!general) return;
    if (!name.trim()) {
      toast.error("Organization name is required");
      return;
    }
    try {
      setSaving(true);
      const updated = await updateSettings({
        name: name.trim(),
        logo_url: logoUrl.trim(),
        industry: industry.trim() || null,
        company_size: companySize.trim() || null,
        country: country.trim() || null,
        general,
      });
      setSettings(updated);
      setName(updated.name);
      setLogoUrl(updated.logo_url ?? "");
      setIndustry(updated.industry ?? "");
      setCompanySize(updated.company_size ?? "");
      setCountry(updated.country ?? "");
      setGeneral(updated.general);
      notifyOrgBrandingUpdated({
        name: updated.name,
        logo_url: updated.logo_url ?? null,
      });
      toast.success("Settings saved");
    } catch {
      toast.error("Failed to save settings");
    } finally {
      setSaving(false);
    }
  }

  if (loading || !general || !settings) {
    return (
      <div className="rounded-3xl border border-slate-200 bg-white p-8 text-slate-400 shadow-sm">
        Loading settings...
      </div>
    );
  }

  return (
    <div className="space-y-6" data-tutorial-surface="settings-home">
      <section className="space-y-5 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        <div>
          <h2 className="text-lg font-semibold text-slate-900">Organization</h2>
          <p className="text-sm text-slate-500">
            Workspace identity and profile. Timezone and currency live under
            Preferences (settings.general).
          </p>
        </div>

        <FormInput
          label="Organization name"
          value={name}
          requiredMark
          onChange={(e) => setName(e.target.value)}
        />

        <OrgLogoUploader
          value={logoUrl}
          orgName={name}
          disabled={saving}
          onChange={setLogoUrl}
        />

        <div className="grid gap-4 sm:grid-cols-2">
          <div className="space-y-1">
            <p className="text-sm font-semibold text-slate-700">Slug</p>
            <div className="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-500">
              {settings.slug}
            </div>
          </div>
          <div className="space-y-1">
            <p className="text-sm font-semibold text-slate-700">Plan</p>
            <div className="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm capitalize text-slate-500">
              {settings.subscription_plan || settings.plan}
            </div>
          </div>
        </div>

        <div className="grid gap-4 sm:grid-cols-2">
          <FormInput
            label="Industry"
            value={industry}
            placeholder="Technology"
            onChange={(e) => setIndustry(e.target.value)}
          />
          <FormSelect
            label="Company size"
            value={companySize}
            onChange={(e) => setCompanySize(e.target.value)}
          >
            <option value="">Select…</option>
            {COMPANY_SIZES.map((s) => (
              <option key={s} value={s}>
                {s}
              </option>
            ))}
          </FormSelect>
          <FormInput
            label="Country"
            value={country}
            placeholder="IN"
            onChange={(e) => setCountry(e.target.value)}
          />
          {settings.status ? (
            <div className="space-y-1">
              <p className="text-sm font-semibold text-slate-700">Status</p>
              <div className="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm capitalize text-slate-500">
                {settings.status}
              </div>
            </div>
          ) : null}
        </div>
      </section>

      <section className="space-y-5 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        <div>
          <h2 className="text-lg font-semibold text-slate-900">Preferences</h2>
          <p className="text-sm text-slate-500">
            Locale and formatting defaults applied across the workspace.
          </p>
        </div>

        <div className="grid gap-4 sm:grid-cols-2">
          <FormSelect
            label="Timezone"
            value={general.timezone}
            onChange={(e) => patchGeneral({ timezone: e.target.value })}
          >
            {TIMEZONES.map((tz) => (
              <option key={tz} value={tz}>
                {tz}
              </option>
            ))}
          </FormSelect>

          <FormSelect
            label="Date format"
            value={general.date_format}
            onChange={(e) => patchGeneral({ date_format: e.target.value })}
          >
            {DATE_FORMATS.map((f) => (
              <option key={f} value={f}>
                {f}
              </option>
            ))}
          </FormSelect>

          <FormSelect
            label="Time format"
            value={general.time_format}
            onChange={(e) =>
              patchGeneral({ time_format: e.target.value as "12h" | "24h" })
            }
          >
            <option value="24h">24-hour</option>
            <option value="12h">12-hour</option>
          </FormSelect>

          <FormSelect
            label="Currency"
            value={general.currency}
            onChange={(e) => patchGeneral({ currency: e.target.value })}
          >
            {CURRENCIES.map((c) => (
              <option key={c} value={c}>
                {c}
              </option>
            ))}
          </FormSelect>

          <FormSelect
            label="Locale"
            value={general.locale}
            onChange={(e) => patchGeneral({ locale: e.target.value })}
          >
            {LOCALES.map((l) => (
              <option key={l} value={l}>
                {l}
              </option>
            ))}
          </FormSelect>

          <FormSelect
            label="Week starts on"
            value={general.week_start}
            onChange={(e) =>
              patchGeneral({ week_start: e.target.value as "sunday" | "monday" })
            }
          >
            <option value="monday">Monday</option>
            <option value="sunday">Sunday</option>
          </FormSelect>
        </div>
      </section>

      <div className="flex justify-end">
        <button
          type="button"
          onClick={handleSave}
          disabled={saving}
          className="inline-flex items-center gap-2 rounded-full bg-emerald-500 px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:opacity-50"
        >
          <Save className="h-4 w-4" />
          {saving ? "Saving..." : "Save changes"}
        </button>
      </div>
    </div>
  );
}
