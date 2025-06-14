// Next.js 15 enhanced caching API functions
export async function getRules() {
  try {
    const response = await fetch('http://chatbot:8090/api/rules', {
      next: { 
        revalidate: 30,
        tags: ['rules']
      }
    })
    
    if (!response.ok) {
      throw new Error('Failed to fetch rules')
    }
    
    return await response.json()
  } catch (error) {
    console.error('Error fetching rules:', error)
    return { pairs: [] }
  }
}

export async function getSystemHealth() {
  try {
    const response = await fetch('http://chatbot:8090/health', {
      next: { 
        revalidate: 10,
        tags: ['health']
      }
    })
    
    if (!response.ok) {
      throw new Error('Health check failed')
    }
    
    return await response.json()
  } catch (error) {
    console.error('Error fetching health status:', error)
    return { status: 'error', database: { is_healthy: false } }
  }
}

// Utility function for API calls with error handling
export async function apiCall<T>(
  endpoint: string, 
  options?: RequestInit
): Promise<T> {
  const baseURL = process.env.NODE_ENV === 'production' 
    ? process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3006'
    : 'http://localhost:3006'
    
  const url = `${baseURL}${endpoint}`
  
  try {
    const response = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
      ...options,
    })
    
    if (!response.ok) {
      throw new Error(`API call failed: ${response.status} ${response.statusText}`)
    }
    
    return await response.json()
  } catch (error) {
    console.error(`API call to ${endpoint} failed:`, error)
    throw error
  }
}
