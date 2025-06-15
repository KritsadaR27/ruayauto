import { NextRequest, NextResponse } from 'next/server'

// Helper function to get current session
async function getCurrentSession(): Promise<{ userId: string } | null> {
    // TODO: Implement actual session retrieval
    return null
}

// Helper function to disconnect page via backend
async function disconnectPageFromBackend(userId: string, pageId: string) {
    try {
        const response = await fetch(`${process.env.CHATBOT_API_URL}/api/facebook/pages/${pageId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${userId}`,
                'Content-Type': 'application/json'
            }
        })

        if (!response.ok) {
            throw new Error('Failed to disconnect page from backend')
        }

        return await response.json()
    } catch (error) {
        console.error('Error disconnecting page from backend:', error)
        throw error
    }
}

export async function DELETE(
    request: NextRequest,
    { params }: { params: { pageId: string } }
) {
    try {
        const { pageId } = params

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

        // Disconnect page via backend API
        await disconnectPageFromBackend(userSession.userId, pageId)

        return NextResponse.json({
            success: true,
            message: 'Page disconnected successfully'
        })
    } catch (error) {
        console.error('Failed to disconnect page:', error)
        return NextResponse.json(
            { error: 'Failed to disconnect page' },
            { status: 500 }
        )
    }
}
