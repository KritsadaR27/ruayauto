import { NextRequest, NextResponse } from 'next/server'

// Mock Facebook pages for development
const mockPages = [
  {
    id: '123456789',
    name: 'Test Page 1',
    access_token: 'mock_page_token_1',
    connected: true,
    status: 'active',
    last_connected: '2025-06-15T00:00:00Z'
  },
  {
    id: '987654321', 
    name: 'Test Page 2',
    access_token: 'mock_page_token_2',
    connected: false,
    status: 'active',
    last_connected: null
  },
  {
    id: '555666777',
    name: 'Demo Store Page',
    access_token: 'mock_page_token_3',
    connected: true,
    status: 'expired',
    last_connected: '2025-06-10T00:00:00Z'
  }
]

export async function GET() {
  try {
    // TODO: Get actual pages from database/Facebook API
    // const userSession = await getCurrentSession()
    // if (!userSession) {
    //   return NextResponse.json({ error: 'Not authenticated' }, { status: 401 })
    // }
    
    // const pages = await getUserPages(userSession.userId)

    return NextResponse.json({
      success: true,
      pages: mockPages
    })
  } catch (error) {
    console.error('Failed to load pages:', error)
    return NextResponse.json(
      { error: 'Failed to load pages' },
      { status: 500 }
    )
  }
}

export async function POST(request: NextRequest) {
  try {
    const { pageId } = await request.json()

    if (!pageId) {
      return NextResponse.json(
        { error: 'Page ID is required' },
        { status: 400 }
      )
    }

    // TODO: Connect page in database
    // await connectFacebookPage(pageId, userSession.userId)

    console.log(`Connecting page: ${pageId}`)

    return NextResponse.json({
      success: true,
      message: 'Page connected successfully'
    })
  } catch (error) {
    console.error('Failed to connect page:', error)
    return NextResponse.json(
      { error: 'Failed to connect page' },
      { status: 500 }
    )
  }
}
