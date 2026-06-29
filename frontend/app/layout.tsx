import type { Metadata } from "next";
import { Geist } from "next/font/google";
import "./globals.css";

import { AuthProvider } from "@/context/AuthContext";

const geist = Geist({
  subsets: ["latin"],
  variable: "--font-geist",
});

export const metadata: Metadata = {
  title: "CRM Lite",
  description:
    "Production-ready CRM built with Go, PostgreSQL, Redis, Next.js and Docker.",
  icons: {
    icon: "/icon.png",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={`${geist.variable} antialiased`}>
        <AuthProvider>
          {children}
        </AuthProvider>
      </body>
    </html>
  );
}