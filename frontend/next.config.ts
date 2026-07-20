import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async redirects() {
    return [
      { source: "/leads", destination: "/m/leads", permanent: false },
      { source: "/leads/:path*", destination: "/m/leads/:path*", permanent: false },
      { source: "/contacts", destination: "/m/contacts", permanent: false },
      { source: "/contacts/:path*", destination: "/m/contacts/:path*", permanent: false },
      { source: "/tasks", destination: "/m/tasks", permanent: false },
      { source: "/tasks/:path*", destination: "/m/tasks/:path*", permanent: false },
      { source: "/forms", destination: "/settings/forms", permanent: false },
      { source: "/tables", destination: "/dashboard", permanent: false },
      {
        source: "/tables/:path*",
        destination: "/settings/tables/:path*",
        permanent: false,
      },
      { source: "/imports", destination: "/settings/imports", permanent: false },
      { source: "/exports", destination: "/settings/exports", permanent: false },
      {
        source: "/settings/data",
        destination: "/settings/imports",
        permanent: false,
      },
    ];
  },
};

export default nextConfig;
