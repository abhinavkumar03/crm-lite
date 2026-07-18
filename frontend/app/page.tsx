import Navbar from "@/components/home/Navbar";
import Hero from "@/components/home/Hero";
import WorkspaceBridge from "@/components/home/WorkspaceBridge";

export default function Home() {
  return (
    <main className="min-h-screen overflow-hidden bg-slate-50">
      <Navbar />
      <Hero />
      <WorkspaceBridge />
    </main>
  );
}
