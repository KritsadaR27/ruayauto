// Next.js 15 enhanced caching API functions
export async function getKeywords() {
  try {
    const response = await fetch('http://localhost:3006/api/keywords', {
      next: { 
        revalidate: 30,
        tags: ['keywords']
      }
    })
    
    if (!response.ok) {
      throw new Error('Failed to fetch keywords')
    }
    
    return await response.json()
  } catch (error) {
    console.error('Error fetching keywords:', error)
    return { pairs: [] }
  }
}

export async function getSystemHealth() {
  try {
    const response = await fetch('http://localhost:3006/health', {
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
