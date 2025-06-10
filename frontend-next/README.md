# Next.js 15 Frontend - ruayAutoMsg

## 🚀 Features

### Next.js 15 Features Implemented
- **Turbopack Development**: Enabled with `--turbo` flag for faster builds
- **React 19 RC Support**: Using React 19.0.0-rc for cutting-edge features
- **App Router**: Full App Router implementation with modern patterns
- **TypeScript Configuration**: Optimized for Next.js 15
- **Enhanced Caching**: Using Next.js 15 caching improvements

### Core Functionality
- **Dashboard Interface**: Real-time system status monitoring
- **Keyword Management**: Add/edit/view Facebook auto-reply keywords
- **API Proxy**: Seamless backend integration through Next.js rewrites
- **Health Monitoring**: Real-time database and system health checks
- **Responsive Design**: Tailwind CSS with mobile-first approach

## 🛠 Tech Stack

- **Next.js**: 15.0.1
- **React**: 19.0.0-rc-69d4b800-20241021
- **TypeScript**: ^5.3.3
- **Tailwind CSS**: ^3.3.6
- **PostCSS & Autoprefixer**: Latest versions

## 📁 Project Structure

```
frontend-next/
├── app/                    # Next.js 15 App Router
│   ├── layout.tsx         # Root layout with metadata
│   ├── page.tsx           # Homepage with dashboard
│   ├── globals.css        # Global styles with Tailwind
│   ├── api/               # API routes
│   │   ├── health/        # Health check proxy
│   │   └── keywords/      # Keywords management API
│   ├── components/        # React components
│   │   ├── KeywordManager.tsx
│   │   └── StatusCard.tsx
│   └── lib/               # Utility functions
│       └── api.ts         # API helpers with caching
├── next.config.ts         # Next.js configuration
├── tailwind.config.js     # Tailwind CSS config
├── tsconfig.json          # TypeScript config
├── Dockerfile             # Production Docker image
└── package.json           # Dependencies and scripts
```

## 🚀 Quick Start

### Development

```bash
# Install dependencies
npm install

# Start development server with Turbopack
npm run dev

# Open browser
open http://localhost:3000
```

### Production Build

```bash
# Build for production
npm run build

# Start production server
npm start
```

### Docker

```bash
# Build Docker image
docker build -t ruayautomsg-frontend-next .

# Run container
docker run -p 3000:3000 ruayautomsg-frontend-next
```

## 🔧 Configuration

### Environment Variables

```bash
# Production API URL (optional)
NEXT_PUBLIC_API_URL=http://your-backend-api:3006

# Disable telemetry
NEXT_TELEMETRY_DISABLED=1
```

### Backend Integration

The frontend automatically proxies API requests to the backend:

```
Frontend (Next.js) :3000
├── /health → Backend :3006/health
├── /api/keywords → Backend :3006/api/keywords
└── /webhook/* → Backend :3006/webhook/*
```

## 📊 Features Overview

### 1. Dashboard Interface
- Real-time system status cards
- Database health monitoring
- API connectivity status
- Webhook status tracking

### 2. Keyword Management
- Add new keyword-response pairs
- View all active keywords
- Support for Thai and English keywords
- Real-time form validation

### 3. API Integration
- Seamless backend communication
- Error handling and loading states
- Automatic retry mechanisms
- Response caching for performance

### 4. Modern UI/UX
- Clean, professional design
- Responsive mobile layout
- Loading indicators
- Error state handling
- Status indicators with emoji

## 🏗 Next.js 15 Specific Features

### Turbopack
```bash
# Development with Turbopack (faster builds)
npm run dev  # automatically uses --turbo flag
```

### Enhanced Caching
```typescript
// API routes with enhanced caching
const response = await fetch('/api/data', {
  next: { 
    revalidate: 30,
    tags: ['keywords']
  }
})
```

### TypeScript Integration
```typescript
// Strict TypeScript configuration
// Type-safe API calls
// Component prop validation
```

## 🐳 Docker Deployment

### Standalone Build
The project is configured for standalone deployment:

```typescript
// next.config.ts
export default {
  output: 'standalone',
  // ... other config
}
```

### Multi-stage Dockerfile
- Optimized for production
- Minimal image size
- Security best practices
- Non-root user execution

## 🔗 Integration with Backend

### API Proxy Configuration
```typescript
// next.config.ts
async rewrites() {
  return [
    {
      source: '/api/:path*',
      destination: 'http://localhost:3006/api/:path*',
    },
    // ... more routes
  ]
}
```

### Health Check Integration
```typescript
// Real-time health monitoring
useEffect(() => {
  const fetchHealth = async () => {
    const response = await fetch('/health')
    setHealthStatus(await response.json())
  }
  
  fetchHealth()
  const interval = setInterval(fetchHealth, 30000)
  return () => clearInterval(interval)
}, [])
```

## 🚀 Performance Features

1. **Turbopack**: Lightning-fast development builds
2. **Standalone Output**: Optimized production deployment
3. **Code Splitting**: Automatic chunk optimization
4. **Image Optimization**: Built-in Next.js image optimization
5. **Caching**: Enhanced API response caching

## 🎯 Production Ready

✅ **Docker Ready**: Multi-stage production Dockerfile  
✅ **TypeScript**: Full type safety  
✅ **Error Handling**: Comprehensive error boundaries  
✅ **Performance**: Optimized builds and caching  
✅ **Security**: Non-root container execution  
✅ **Scalability**: Standalone deployment ready  

## 📝 Next Steps

1. **Authentication**: Add user authentication system
2. **Real-time Updates**: WebSocket integration for live updates  
3. **Analytics Dashboard**: Extended analytics and reporting
4. **Mobile App**: React Native version using shared components
5. **PWA Features**: Service worker and offline support

The Next.js 15 frontend is now ready for production deployment with modern React 19 features and optimized performance! 🎉
