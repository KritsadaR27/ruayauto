import { NextRequest, NextResponse } from 'next/server'

// Helper function to store user session and pages via backend
async function storeUserSessionInBackend(userData: any, accessToken: string, pagesData: any[]) {
  try {
    const response = await fetch(`${process.env.CHATBOT_API_URL}/api/facebook/auth/callback`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        user: userData,
        access_token: accessToken,
        pages: pagesData
      })
    })

    if (!response.ok) {
      throw new Error('Failed to store session in backend')
    }

    return await response.json()
  } catch (error) {
    console.error('Error storing session in backend:', error)
    throw error
  }
}

// Helper function to create user session
async function createUserSession(userData: any): Promise<string> {
  // TODO: Implement actual session creation (JWT, cookies, etc.)
  // For now, return a mock session ID
  return `session_${userData.id}_${Date.now()}`
}

export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)
    const code = searchParams.get('code')
    const state = searchParams.get('state')
    const error = searchParams.get('error')

    if (error) {
      // User denied permission or other error
      return NextResponse.redirect(`${process.env.NEXTAUTH_URL}?error=facebook_denied`)
    }

    if (!code) {
      return NextResponse.redirect(`${process.env.NEXTAUTH_URL}?error=no_code`)
    }

    // Exchange code for access token
    const tokenResponse = await fetch('https://graph.facebook.com/v18.0/oauth/access_token', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: new URLSearchParams({
        client_id: process.env.FACEBOOK_APP_ID!,
        client_secret: process.env.FACEBOOK_APP_SECRET!,
        redirect_uri: `${process.env.NEXTAUTH_URL}/api/auth/facebook/callback`,
        code
      })
    })

    if (!tokenResponse.ok) {
      throw new Error('Failed to exchange code for token')
    }

    const tokenData = await tokenResponse.json()

    // Get user information
    const userResponse = await fetch(`https://graph.facebook.com/v18.0/me?access_token=${tokenData.access_token}&fields=id,name,email`)
    if (!userResponse.ok) {
      throw new Error('Failed to get user information')
    }
    const userData = await userResponse.json()

    // Get user's pages
    const pagesResponse = await fetch(`https://graph.facebook.com/v18.0/me/accounts?access_token=${tokenData.access_token}&fields=id,name,access_token`)
    if (!pagesResponse.ok) {
      throw new Error('Failed to get user pages')
    }
    const pagesData = await pagesResponse.json()

    // Store user session and pages in backend
    await storeUserSessionInBackend(userData, tokenData.access_token, pagesData.data || [])

    // Create frontend session
    const sessionId = await createUserSession(userData)

    // TODO: Set session cookie
    // cookies().set('session_id', sessionId, { httpOnly: true, secure: true })

    console.log('Facebook OAuth Success:', {
      user: userData.name,
      userId: userData.id,
      pages: pagesData.data?.length || 0
    })

    // Redirect to main page with success
    return NextResponse.redirect(`${process.env.NEXTAUTH_URL}?facebook_connected=true`)

  } catch (error) {
    console.error('Facebook OAuth callback error:', error)
    return NextResponse.redirect(`${process.env.NEXTAUTH_URL}?error=oauth_failed`)
  }
}
