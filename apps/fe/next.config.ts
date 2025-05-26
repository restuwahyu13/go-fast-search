import type { NextConfig } from 'next'
import os from 'node:os'

const nextConfig: NextConfig = {
  compress: true,
  reactStrictMode: true,
  productionBrowserSourceMaps: false,
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL
  },
  experimental: {
    caseSensitiveRoutes: true,
    optimizeCss: true,
    cpus: os.cpus().length / 2,
  },
}

export default nextConfig