import { NextRequest, NextResponse } from 'next/server'

// Helper function to get current session
async function getCurrentSession(): Promise<{ userId: string, facebookUserId?: string } | null> {
    // TODO: Implement actual session retrieval from cookies/JWT
    // This should check for valid session tokens
    return null
}

// Helper function to check Facebook auth status via backend
async function checkFacebookAuthStatus(userId: string) {
    try {
        const response = await fetch(`${process.env.CHATBOT_API_URL}/api/facebook/auth/status`, {
            headers: {
                'Authorization': `Bearer ${userId}`,
                'Content-Type': 'application/json'
            }
        })

        if (!response.ok) {
            throw new Error('Failed to check auth status from backend')
        }

        return await response.json()
    } catch (error) {
        console.error('Error checking Facebook auth status:', error)
        return { authenticated: false, user: null }
    }
}

export async function GET() {
    try {
        // Check if user has valid session
        const session = await getCurrentSession()

        if (!session) {
            return NextResponse.json({
                authenticated: false,
                user: null
            })
        }

        // Check Facebook authentication status via backend
        const authStatus = await checkFacebookAuthStatus(session.userId)

        return NextResponse.json(authStatus)
    } catch (error) {
        console.error('Auth status check failed:', error)
        return NextResponse.json(
            { error: 'Failed to check authentication status' },
            { status: 500 }
        )
    }
}
