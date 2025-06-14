import { NextRequest } from 'next/server'

// Next.js 15 async request APIs demonstration
export async function GET(request: NextRequest) {
  try {
    // Using Next.js 15 async request APIs
    const { searchParams } = request.nextUrl
    const limit = searchParams.get('limit') || '10'
    
    // Enhanced caching with 'use cache' directive (Next.js 15 feature)
    const response = await fetch('http://chatbot:8090/api/keywords', {
      next: { 
        revalidate: 30,
        tags: ['keywords']
      }
    })
    
    if (!response.ok) {
      throw new Error('Failed to fetch from backend API')
    }
    
    const data = await response.json()
    
    return Response.json({
      success: true,
      data: data,
      timestamp: new Date().toISOString(),
      cached: false
    })
  } catch (error) {
    console.error('API Error:', error)
    return Response.json(
      { 
        error: 'Failed to fetch keywords',
        message: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    )
  }
}

export async function POST(request: NextRequest) {
  try {
    const body = await request.json()
    
    // Forward to backend API
    const response = await fetch('http://chatbot:8090/api/keywords', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body)
    })
    
    if (!response.ok) {
      throw new Error('Failed to save to backend API')
    }
    
    const data = await response.json()
    
    // Revalidate the keywords cache
    // In Next.js 15, this would use the enhanced caching system
    
    return Response.json({
      success: true,
      data: data,
      timestamp: new Date().toISOString()
    })
  } catch (error) {
    console.error('API Error:', error)
    return Response.json(
      { 
        error: 'Failed to save keywords',
        message: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    )
  }
}
