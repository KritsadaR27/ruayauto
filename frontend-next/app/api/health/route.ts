import { NextRequest } from 'next/server'

export async function GET(request: NextRequest) {
  try {
    // Proxy to backend health endpoint
    const response = await fetch('http://localhost:3006/health', {
      next: { revalidate: 10 } // Cache for 10 seconds
    })
    
    if (!response.ok) {
      // Return healthy status even if backend health endpoint fails
      return Response.json({
        status: 'healthy',
        backend: 'checking',
        frontend: 'healthy',
        nextjs_version: '15.0.1',
        react_version: '19.0.0-rc',
        timestamp: new Date().toISOString(),
        message: 'Frontend is running, backend status unknown'
      })
    }
    
    const healthData = await response.json()
    
    return Response.json({
      ...healthData,
      frontend: 'healthy',
      nextjs_version: '15.0.1',
      react_version: '19.0.0-rc',
      timestamp: new Date().toISOString()
    })
  } catch (error) {
    return Response.json(
      { 
        status: 'healthy',
        frontend: 'healthy',
        backend: 'unavailable',
        error: error instanceof Error ? error.message : 'Unknown error',
        timestamp: new Date().toISOString()
      },
      { status: 200 }
    )
  }
}
