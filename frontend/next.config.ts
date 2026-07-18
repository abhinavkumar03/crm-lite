import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async redirects() {
    return [
      { source: "/leads", destination: "/tables", permanent: false },
      { source: "/leads/:path*", destination: "/tables", permanent: false },
      { source: "/contacts", destination: "/tables", permanent: false },
      { source: "/contacts/:path*", destination: "/tables", permanent: false },
      { source: "/tasks", destination: "/tables", permanent: false },
      { source: "/tasks/:path*", destination: "/tables", permanent: false },
    ];
  },
};

export default nextConfig;
