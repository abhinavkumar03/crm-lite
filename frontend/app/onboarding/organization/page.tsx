"use client";

import { useMemo, useState, type ReactNode } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import {
  ArrowLeft,
  ArrowRight,
  Building2,
  Check,
  Loader2,
  Mail,
} from "lucide-react";

import { createOrganization } from "@/features/organization/api";
import { useAuth } from "@/context/AuthContext";

const inputClass =
  "w-full rounded-2xl border border-slate-200 px-4 py-3 text-sm outline-none focus:border-emerald-500 focus:ring-4 focus:ring-emerald-100";

const STEPS = [
  "Welcome",
  "Name",
  "Slug",
  "Industry",
  "Size",
  "Timezone",
  "Currency",
  "Country",
  "Finish",
] as const;

const INDUSTRIES = [
  "Technology",
  "Healthcare",
  "Finance",
  "Education",
  "Retail",
  "Manufacturing",
  "Consulting",
  "Other",
];

const COMPANY_SIZES = ["1-10", "11-50", "51-200", "201-500", "500+"];

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

const CURRENCIES = ["USD", "EUR", "GBP", "INR", "AUD", "CAD", "SGD", "JPY"];

function slugify(value: string): string {
  return value
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "")
    .slice(0, 80);
}

export default function OrganizationOnboardingPage() {
  const router = useRouter();
  const auth = useAuth();

  const [step, setStep] = useState(0);
  const [mode, setMode] = useState<"create" | "invite" | null>(null);
  const [name, setName] = useState("");
  const [slug, setSlug] = useState("");
  const [slugTouched, setSlugTouched] = useState(false);
  const [industry, setIndustry] = useState("Technology");
  const [companySize, setCompanySize] = useState("11-50");
  const [timezone, setTimezone] = useState("Asia/Kolkata");
  const [currency, setCurrency] = useState("INR");
  const [country, setCountry] = useState("IN");
  const [saving, setSaving] = useState(false);

  const progress = useMemo(
    () => Math.round(((step + 1) / STEPS.length) * 100),
    [step]
  );

  function goNext() {
    if (step === 0 && mode !== "create") {
      if (mode === "invite") return;
      toast.error("Choose how you want to continue");
      return;
    }
    if (step === 1 && !name.trim()) {
      toast.error("Organization name is required");
      return;
    }
    if (step === 2 && !slugify(slug || name)) {
      toast.error("Slug is required");
      return;
    }
    setStep((s) => Math.min(s + 1, STEPS.length - 1));
  }

  function goBack() {
    setStep((s) => Math.max(s - 1, 0));
  }

  async function handleFinish() {
    if (!name.trim()) {
      toast.error("Organization name is required");
      setStep(1);
      return;
    }
    try {
      setSaving(true);
      await createOrganization({
        name: name.trim(),
        slug: slugify(slug || name),
        industry,
        company_size: companySize,
        country: country.trim(),
        general: {
          timezone,
          currency,
          locale: timezone.startsWith("Asia/Kolkata") ? "en-IN" : "en-US",
        },
      });
      toast.success("Workspace ready");
      router.replace("/dashboard");
    } catch (err: unknown) {
      const message =
        (err as { response?: { data?: { message?: string } } })?.response?.data
          ?.message ?? "Could not create workspace";
      toast.error(message);
    } finally {
      setSaving(false);
    }
  }

  return (
    <div className="mx-auto flex min-h-screen w-full max-w-2xl flex-col px-4 py-10 sm:px-6">
      <div className="mb-8">
        <p className="text-sm font-semibold text-emerald-600">CRM Lite</p>
        <h1 className="mt-2 text-3xl font-semibold tracking-tight text-slate-900">
          Welcome{auth.user?.name ? `, ${auth.user.name.split(" ")[0]}` : ""}!
        </h1>
        <p className="mt-2 text-slate-500">
          Let&apos;s set up your workspace before you open the CRM.
        </p>
      </div>

      <div className="mb-6 h-2 overflow-hidden rounded-full bg-slate-200">
        <div
          className="h-full rounded-full bg-emerald-500 transition-all"
          style={{ width: `${progress}%` }}
        />
      </div>

      <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm sm:p-8">
        {step === 0 && (
          <div className="space-y-4">
            <h2 className="text-xl font-semibold text-slate-900">
              How do you want to get started?
            </h2>
            <button
              type="button"
              onClick={() => setMode("create")}
              className={`flex w-full items-start gap-4 rounded-2xl border px-4 py-4 text-left transition ${
                mode === "create"
                  ? "border-emerald-300 bg-emerald-50"
                  : "border-slate-200 hover:bg-slate-50"
              }`}
            >
              <Building2 className="mt-0.5 h-5 w-5 text-emerald-600" />
              <span>
                <span className="block font-semibold text-slate-900">
                  Create organization
                </span>
                <span className="text-sm text-slate-500">
                  Start a new workspace. You become the Owner.
                </span>
              </span>
            </button>
            <button
              type="button"
              onClick={() => setMode("invite")}
              className={`flex w-full items-start gap-4 rounded-2xl border px-4 py-4 text-left transition ${
                mode === "invite"
                  ? "border-emerald-300 bg-emerald-50"
                  : "border-slate-200 hover:bg-slate-50"
              }`}
            >
              <Mail className="mt-0.5 h-5 w-5 text-emerald-600" />
              <span>
                <span className="block font-semibold text-slate-900">
                  Join via invite
                </span>
                <span className="text-sm text-slate-500">
                  Coming soon — ask your admin for an invitation link.
                </span>
              </span>
            </button>
            {mode === "invite" && (
              <div className="rounded-2xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-800">
                Invite acceptance UI is a placeholder for this MVP. Ask an Owner
                to invite you from Settings → Members, then use the invite token
                endpoint when available.
              </div>
            )}
          </div>
        )}

        {step === 1 && (
          <Field
            label="Organization name"
            hint="Example: Acme Technologies"
          >
            <input
              value={name}
              onChange={(e) => {
                setName(e.target.value);
                if (!slugTouched) setSlug(slugify(e.target.value));
              }}
              placeholder="Acme Technologies"
              className={inputClass}
              autoFocus
            />
          </Field>
        )}

        {step === 2 && (
          <Field label="Workspace slug" hint="Used in URLs and identifiers">
            <input
              value={slug}
              onChange={(e) => {
                setSlugTouched(true);
                setSlug(slugify(e.target.value));
              }}
              placeholder="acme-technologies"
              className={inputClass}
              autoFocus
            />
          </Field>
        )}

        {step === 3 && (
          <Field label="Industry">
            <select
              value={industry}
              onChange={(e) => setIndustry(e.target.value)}
              className={inputClass}
            >
              {INDUSTRIES.map((i) => (
                <option key={i} value={i}>
                  {i}
                </option>
              ))}
            </select>
          </Field>
        )}

        {step === 4 && (
          <Field label="Company size">
            <select
              value={companySize}
              onChange={(e) => setCompanySize(e.target.value)}
              className={inputClass}
            >
              {COMPANY_SIZES.map((s) => (
                <option key={s} value={s}>
                  {s}
                </option>
              ))}
            </select>
          </Field>
        )}

        {step === 5 && (
          <Field label="Timezone">
            <select
              value={timezone}
              onChange={(e) => setTimezone(e.target.value)}
              className={inputClass}
            >
              {TIMEZONES.map((tz) => (
                <option key={tz} value={tz}>
                  {tz}
                </option>
              ))}
            </select>
          </Field>
        )}

        {step === 6 && (
          <Field label="Currency">
            <select
              value={currency}
              onChange={(e) => setCurrency(e.target.value)}
              className={inputClass}
            >
              {CURRENCIES.map((c) => (
                <option key={c} value={c}>
                  {c}
                </option>
              ))}
            </select>
          </Field>
        )}

        {step === 7 && (
          <Field label="Country" hint="ISO country code or name">
            <input
              value={country}
              onChange={(e) => setCountry(e.target.value)}
              placeholder="IN"
              className={inputClass}
              autoFocus
            />
          </Field>
        )}

        {step === 8 && (
          <div className="space-y-4">
            <h2 className="text-xl font-semibold text-slate-900">
              Review & create
            </h2>
            <ul className="space-y-2 text-sm text-slate-600">
              <li>
                <strong className="text-slate-800">Name:</strong> {name}
              </li>
              <li>
                <strong className="text-slate-800">Slug:</strong>{" "}
                {slugify(slug || name)}
              </li>
              <li>
                <strong className="text-slate-800">Industry:</strong> {industry}
              </li>
              <li>
                <strong className="text-slate-800">Size:</strong> {companySize}
              </li>
              <li>
                <strong className="text-slate-800">Timezone:</strong> {timezone}
              </li>
              <li>
                <strong className="text-slate-800">Currency:</strong> {currency}
              </li>
              <li>
                <strong className="text-slate-800">Country:</strong> {country}
              </li>
            </ul>
            <p className="text-sm text-slate-500">
              We&apos;ll create your organization, assign you as Owner, seed
              default roles, and set up Companies & Deals modules.
            </p>
          </div>
        )}

        <div className="mt-8 flex items-center justify-between gap-3">
          <button
            type="button"
            onClick={goBack}
            disabled={step === 0 || saving}
            className="inline-flex items-center gap-2 rounded-full border border-slate-200 px-4 py-2 text-sm font-semibold text-slate-600 disabled:opacity-40"
          >
            <ArrowLeft className="h-4 w-4" />
            Back
          </button>

          {step < STEPS.length - 1 ? (
            <button
              type="button"
              onClick={goNext}
              disabled={mode === "invite" && step === 0}
              className="inline-flex items-center gap-2 rounded-full bg-emerald-500 px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:opacity-40"
            >
              Continue
              <ArrowRight className="h-4 w-4" />
            </button>
          ) : (
            <button
              type="button"
              onClick={handleFinish}
              disabled={saving}
              className="inline-flex items-center gap-2 rounded-full bg-emerald-500 px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:opacity-50"
            >
              {saving ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <Check className="h-4 w-4" />
              )}
              {saving ? "Creating…" : "Create workspace"}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}

function Field({
  label,
  hint,
  children,
}: {
  label: string;
  hint?: string;
  children: ReactNode;
}) {
  return (
    <div className="space-y-2">
      <h2 className="text-xl font-semibold text-slate-900">{label}</h2>
      {hint ? <p className="text-sm text-slate-500">{hint}</p> : null}
      {children}
    </div>
  );
}
