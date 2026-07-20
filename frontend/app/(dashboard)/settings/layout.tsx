import PageHeader from "@/components/common/PageHeader";
import SettingsNav from "@/components/settings/SettingsNav";

export default function SettingsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="space-y-8">
      <div className="grid gap-8 lg:grid-cols-[280px_1fr]">
        <aside className="lg:sticky lg:self-start">
          <SettingsNav />
        </aside>

        <div className="min-w-0">{children}</div>
      </div>
    </div>
  );
}
