import type { NextConfig } from 'next'

const nextConfig: NextConfig = {
  // Enable Turbopack for development
  experimental: {
    // Enable Turbopack for faster development builds
    turbo: {
      rules: {
        '*.svg': {
          loaders: ['@svgr/webpack'],
          as: '*.js',
        },
      },
    },
  },

  // TypeScript configuration
  typescript: {
    // Dangerously allow production builds to successfully complete even if your project has type errors
    ignoreBuildErrors: false,
  },

  // Rewrites for API proxy to backend
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://chatbot:8090/api/:path*',
      },
      {
        source: '/webhook/:path*',
        destination: 'http://webhook:8091/webhook/:path*',
      },
      {
        source: '/health',
        destination: 'http://chatbot:8090/health',
      },
    ]
  },

  // Image optimization
  images: {
    domains: ['localhost'],
    unoptimized: process.env.NODE_ENV === 'development',
  },

  // Output configuration for Docker deployment
  output: 'standalone',

  // Bundle analyzer (optional)
  env: {
    NEXT_TELEMETRY_DISABLED: '1',
  },

  // Webpack configuration for additional optimizations
  webpack: (config, { dev, isServer }) => {
    // Add any custom webpack configuration here
    if (!dev && !isServer) {
      // Production client-side optimizations
      config.optimization.splitChunks.chunks = 'all'
    }
    return config
  },
}

export default nextConfig
