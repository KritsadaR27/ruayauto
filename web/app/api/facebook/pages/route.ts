import { NextRequest, NextResponse } from 'next/server'

// Helper function to get user session (implement based on your auth system)
async function getCurrentSession(): Promise<{ userId: string } | null> {
    // TODO: Implement actual session retrieval
    // For now, return null to indicate no authenticated user
    return null
}

// Helper function to get Facebook pages from the backend API
async function getUserPages(userId: string) {
    try {
        const response = await fetch(`${process.env.CHATBOT_API_URL}/api/facebook/pages`, {
            headers: {
                'Authorization': `Bearer ${userId}`, // Use proper auth header
                'Content-Type': 'application/json'
            }
        })

        if (!response.ok) {
            throw new Error('Failed to fetch pages from backend')
        }

        return await response.json()
    } catch (error) {
        console.error('Error fetching pages from backend:', error)
        throw error
    }
}

export async function GET() {
    try {
        // Check if user is authenticated
        const userSession = await getCurrentSession()
        if (!userSession) {
            return NextResponse.json({
                error: 'Not authenticated',
                pages: []
            }, { status: 401 })
        }

        // Get pages from backend/database
        const pages = await getUserPages(userSession.userId)

        return NextResponse.json({
            success: true,
            pages: pages || []
        })
    } catch (error) {
        console.error('Failed to load pages:', error)
        return NextResponse.json(
            {
                error: 'Failed to load pages',
                pages: []
            },
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

        // Check if user is authenticated
        const userSession = await getCurrentSession()
        if (!userSession) {
            return NextResponse.json({
                error: 'Not authenticated'
            }, { status: 401 })
        }

        // Connect page via backend API
        const response = await fetch(`${process.env.CHATBOT_API_URL}/api/facebook/pages/connect`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${userSession.userId}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ pageId })
        })

        if (!response.ok) {
            throw new Error('Failed to connect page via backend')
        }

        const result = await response.json()

        return NextResponse.json({
            success: true,
            message: 'Page connected successfully',
            data: result
        })
    } catch (error) {
        console.error('Failed to connect page:', error)
        return NextResponse.json(
            { error: 'Failed to connect page' },
            { status: 500 }
        )
    }
}
