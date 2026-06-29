import Navbar from "@/components/home/Navbar";
import Hero from "@/components/home/Hero";

export default function Home() {
  return (
    <main className="min-h-screen overflow-hidden bg-slate-50">
      <Navbar />
      <Hero />
    </main>
  );
}