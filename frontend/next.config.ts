import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: '/api/search',
        destination: 'http://localhost:8083/search',
      },
      {
        source: '/api/alerts',
        destination: 'http://localhost:8085/alerts',
      },
    ];
  },
};

export default nextConfig;
