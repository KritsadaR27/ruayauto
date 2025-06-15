'use server'

import type { KeywordData, Pair, ApiResponse } from '../types/keyword'

// Server action สำหรับ save rule data
export async function saveRuleData(data: KeywordData): Promise<ApiResponse> {
  try {
    // Transform data for backend API
    const payload = {
      pairs: data.pairs.map(pair => ({
        keyword: pair.keywords.join(','),
        response: pair.responses.join('|')
      }))
    }

    const response = await fetch('http://chatbot:8090/api/rules', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
      // Next.js 15 caching
      next: { revalidate: 0 }
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const result = await response.json()
    
    return {
      success: true,
      data: result,
      timestamp: new Date().toISOString()
    }
  } catch (error) {
    console.error('Save keyword data error:', error)
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Failed to save data',
      timestamp: new Date().toISOString()
    }
  }
}

// Server action สำหรับ load rule data
export async function loadRuleData(): Promise<ApiResponse<KeywordData>> {
  try {
    const response = await fetch('http://chatbot:8090/api/rules', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      // Next.js 15 caching
      next: { revalidate: 30, tags: ['keywords'] }
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const result = await response.json()
    
    // Transform backend data to frontend format
    const transformedData: KeywordData = {
      pairs: result.pairs?.map((item: any) => ({
        keywords: item.keyword ? item.keyword.split(',').filter(Boolean) : [],
        responses: item.response ? item.response.split('|').filter(Boolean) : []
      })) || [],
      defaultResponses: ['ขอบคุณสำหรับข้อความทดสอบครับ', 'ได้รับข้อความทดสอบแล้วครับ'],
      enableDefault: true,
      noTag: false,
      noSticker: false,
      hideCommentsAfterReply: false,
      enableInboxIntegration: false,
      inboxResponse: '',
      filterSettings: { skipMentions: false, skipStickers: false },
      fallbackSettings: { enabled: false, responses: [{ text: '' }] }
    }
    
    return {
      success: true,
      data: transformedData,
      timestamp: new Date().toISOString()
    }
  } catch (error) {
    console.error('Load keyword data error:', error)
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Failed to load data',
      timestamp: new Date().toISOString()
    }
  }
}

// Server action สำหรับ health check
export async function checkSystemHealth(): Promise<ApiResponse> {
  try {
    const response = await fetch('http://chatbot:8090/health', {
      method: 'GET',
      // Next.js 15 caching
      next: { revalidate: 10, tags: ['health'] }
    })

    if (!response.ok) {
      throw new Error(`Health check failed: ${response.status}`)
    }

    const healthData = await response.json()
    
    return {
      success: true,
      data: {
        ...healthData,
        frontend: 'healthy',
        nextjs_version: '15.0.1'
      },
      timestamp: new Date().toISOString()
    }
  } catch (error) {
    console.error('Health check error:', error)
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Health check failed',
      data: {
        status: 'error',
        frontend: 'healthy',
        backend: 'unhealthy'
      },
      timestamp: new Date().toISOString()
    }
  }
}
