import PageHeader from "@/components/common/PageHeader";
import SettingsNav from "@/components/settings/SettingsNav";

export default function SettingsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="space-y-8">
      <PageHeader
        badge="Settings Center"
        title="Settings"
        description="Configure your organization, shape your data model, and tune automation — all metadata-driven, no redeploys required."
      />

      <div className="grid gap-8 lg:grid-cols-[280px_1fr]">
        <aside className="lg:sticky lg:top-24 lg:self-start">
          <SettingsNav />
        </aside>

        <div className="min-w-0">{children}</div>
      </div>
    </div>
  );
}
