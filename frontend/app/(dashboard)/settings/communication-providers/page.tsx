"use client";

import { useEffect, useState } from "react";
import { toast } from "sonner";
import { Plus, Trash2, Zap } from "lucide-react";

import PageHeader from "@/components/common/PageHeader";
import Modal from "@/components/common/Modal";
import FormInput from "@/components/common/form/FormInput";
import FormSelect from "@/components/common/form/FormSelect";

import {
  CommunicationProvider,
  ProviderChannel,
  createProvider,
  deleteProvider,
  listProviders,
  testProvider,
} from "@/features/notifications/providersApi";

const EMAIL_TYPES = ["smtp", "resend", "ses", "sendgrid", "mailgun", "simulation"];
const WA_TYPES = ["meta", "twilio", "gupshup", "interakt", "360dialog", "simulation"];

export default function CommunicationProvidersPage() {
  const [channel, setChannel] = useState<ProviderChannel>("email");
  const [items, setItems] = useState<CommunicationProvider[]>([]);
  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");
  const [providerType, setProviderType] = useState("smtp");
  const [smtpHost, setSmtpHost] = useState("");
  const [smtpPort, setSmtpPort] = useState("587");
  const [from, setFrom] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [apiKey, setApiKey] = useState("");
  const [phoneId, setPhoneId] = useState("");
  const [token, setToken] = useState("");
  const [accountSid, setAccountSid] = useState("");
  const [authToken, setAuthToken] = useState("");
  const [fromNumber, setFromNumber] = useState("");
  const [isDefault, setIsDefault] = useState(true);
  const [saving, setSaving] = useState(false);
  const [reloadKey, setReloadKey] = useState(0);
  const [testTo, setTestTo] = useState("");

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const rows = await listProviders(channel);
        if (active) setItems(rows);
      } catch {
        toast.error("Failed to load providers");
      }
    })();
    return () => {
      active = false;
    };
  }, [channel, reloadKey]);

  useEffect(() => {
    setProviderType(channel === "email" ? "smtp" : "meta");
  }, [channel]);

  async function handleCreate() {
    if (!name.trim()) {
      toast.error("Name is required");
      return;
    }
    try {
      setSaving(true);
      const config: Record<string, unknown> = {};
      const secrets: Record<string, unknown> = {};
      if (channel === "email") {
        config.from = from;
        config.default_sender = from;
        if (providerType === "smtp") {
          config.smtp_host = smtpHost;
          config.smtp_port = Number(smtpPort) || 587;
          config.encryption = "starttls";
          secrets.smtp_username = username;
          secrets.smtp_password = password;
        } else {
          secrets.api_key = apiKey;
        }
      } else {
        if (providerType === "meta") {
          config.phone_number_id = phoneId;
          config.api_url = "https://graph.facebook.com/v20.0";
          secrets.access_token = token;
        } else if (providerType === "twilio") {
          config.from_number = fromNumber;
          secrets.account_sid = accountSid;
          secrets.auth_token = authToken;
        }
      }
      await createProvider({
        channel,
        provider_type: providerType,
        name: name.trim(),
        config,
        secrets,
        is_default: isDefault,
        is_active: true,
      });
      toast.success("Provider saved");
      setOpen(false);
      setReloadKey((k) => k + 1);
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : "Failed to save provider");
    } finally {
      setSaving(false);
    }
  }

  async function handleTest(id: string) {
    if (!testTo.trim()) {
      toast.error("Enter a test recipient");
      return;
    }
    try {
      await testProvider(id, testTo.trim());
      toast.success("Provider test succeeded");
      setReloadKey((k) => k + 1);
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : "Provider test failed");
    }
  }

  async function handleDelete(id: string) {
    if (!confirm("Delete this provider?")) return;
    try {
      await deleteProvider(id);
      toast.success("Provider deleted");
      setReloadKey((k) => k + 1);
    } catch {
      toast.error("Failed to delete");
    }
  }

  const types = channel === "email" ? EMAIL_TYPES : WA_TYPES;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Communication Providers"
        description="Configure email and WhatsApp delivery. Secrets are encrypted and never returned by the API."
        action={
          <button
            type="button"
            onClick={() => setOpen(true)}
            className="inline-flex items-center gap-2 rounded-xl bg-emerald-600 px-4 py-2 text-sm font-semibold text-white hover:bg-emerald-700"
          >
            <Plus size={16} /> Add provider
          </button>
        }
      />

      <div className="flex flex-wrap items-end gap-3">
        <FormSelect
          label="Channel"
          value={channel}
          onChange={(e) => setChannel(e.target.value as ProviderChannel)}
        >
          <option value="email">Email</option>
          <option value="whatsapp">WhatsApp</option>
        </FormSelect>
        <FormInput
          label="Test recipient"
          value={testTo}
          onChange={(e) => setTestTo(e.target.value)}
          placeholder={channel === "email" ? "you@example.com" : "+15551234567"}
        />
      </div>

      <div className="overflow-hidden rounded-2xl border border-slate-200 bg-white">
        <table className="min-w-full text-sm">
          <thead className="bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500">
            <tr>
              <th className="px-4 py-3">Name</th>
              <th className="px-4 py-3">Type</th>
              <th className="px-4 py-3">Status</th>
              <th className="px-4 py-3">Health</th>
              <th className="px-4 py-3" />
            </tr>
          </thead>
          <tbody>
            {items.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-4 py-8 text-center text-slate-500">
                  No providers configured for this channel. Env bootstrap still works for local
                  development.
                </td>
              </tr>
            ) : (
              items.map((p) => (
                <tr key={p.id} className="border-t border-slate-100">
                  <td className="px-4 py-3 font-medium text-slate-800">
                    {p.name}
                    {p.is_default ? (
                      <span className="ml-2 rounded-full bg-emerald-50 px-2 py-0.5 text-xs text-emerald-700">
                        default
                      </span>
                    ) : null}
                  </td>
                  <td className="px-4 py-3 text-slate-600">{p.provider_type}</td>
                  <td className="px-4 py-3">
                    {p.is_active ? "Active" : "Inactive"}
                    {p.secrets_configured ? " · secrets set" : " · missing secrets"}
                  </td>
                  <td className="px-4 py-3 text-slate-500">
                    {p.last_error ? (
                      <span className="text-red-600">{p.last_error}</span>
                    ) : p.last_health_at ? (
                      `OK · ${new Date(p.last_health_at).toLocaleString()}`
                    ) : (
                      "—"
                    )}
                  </td>
                  <td className="px-4 py-3 text-right">
                    <button
                      type="button"
                      onClick={() => handleTest(p.id)}
                      className="mr-2 inline-flex items-center gap-1 rounded-lg border border-slate-200 px-2 py-1 text-xs font-medium text-slate-700 hover:bg-slate-50"
                    >
                      <Zap size={12} /> Test
                    </button>
                    <button
                      type="button"
                      onClick={() => handleDelete(p.id)}
                      className="inline-flex items-center gap-1 rounded-lg border border-red-100 px-2 py-1 text-xs font-medium text-red-600 hover:bg-red-50"
                    >
                      <Trash2 size={12} /> Delete
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      <Modal open={open} onClose={() => setOpen(false)} title="Add provider">
        <div className="space-y-3">
          <FormInput label="Name" value={name} onChange={(e) => setName(e.target.value)} />
          <FormSelect
            label="Provider type"
            value={providerType}
            onChange={(e) => setProviderType(e.target.value)}
          >
            {types.map((t) => (
              <option key={t} value={t}>
                {t}
              </option>
            ))}
          </FormSelect>
          {channel === "email" && providerType === "smtp" ? (
            <>
              <FormInput label="SMTP host" value={smtpHost} onChange={(e) => setSmtpHost(e.target.value)} />
              <FormInput label="SMTP port" value={smtpPort} onChange={(e) => setSmtpPort(e.target.value)} />
              <FormInput label="From" value={from} onChange={(e) => setFrom(e.target.value)} />
              <FormInput label="Username" value={username} onChange={(e) => setUsername(e.target.value)} />
              <FormInput
                label="Password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </>
          ) : null}
          {channel === "email" && providerType !== "smtp" && providerType !== "simulation" ? (
            <>
              <FormInput label="From" value={from} onChange={(e) => setFrom(e.target.value)} />
              <FormInput
                label="API key"
                type="password"
                value={apiKey}
                onChange={(e) => setApiKey(e.target.value)}
              />
            </>
          ) : null}
          {channel === "whatsapp" && providerType === "meta" ? (
            <>
              <FormInput label="Phone number ID" value={phoneId} onChange={(e) => setPhoneId(e.target.value)} />
              <FormInput
                label="Access token"
                type="password"
                value={token}
                onChange={(e) => setToken(e.target.value)}
              />
            </>
          ) : null}
          {channel === "whatsapp" && providerType === "twilio" ? (
            <>
              <FormInput label="From number" value={fromNumber} onChange={(e) => setFromNumber(e.target.value)} />
              <FormInput label="Account SID" value={accountSid} onChange={(e) => setAccountSid(e.target.value)} />
              <FormInput
                label="Auth token"
                type="password"
                value={authToken}
                onChange={(e) => setAuthToken(e.target.value)}
              />
            </>
          ) : null}
          <label className="flex items-center gap-2 text-sm text-slate-700">
            <input
              type="checkbox"
              checked={isDefault}
              onChange={(e) => setIsDefault(e.target.checked)}
            />
            Set as default for this channel
          </label>
          <button
            type="button"
            disabled={saving}
            onClick={handleCreate}
            className="w-full rounded-xl bg-emerald-600 px-4 py-2 text-sm font-semibold text-white hover:bg-emerald-700 disabled:opacity-60"
          >
            {saving ? "Saving…" : "Save provider"}
          </button>
        </div>
      </Modal>
    </div>
  );
}
