import { NextRequest, NextResponse } from 'next/server'

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
        redirect_uri: `${process.env.NEXTAUTH_URL}/auth/facebook/callback`,
        code
      })
    })

    if (!tokenResponse.ok) {
      throw new Error('Failed to exchange code for token')
    }

    const tokenData = await tokenResponse.json()

    // Get user information
    const userResponse = await fetch(`https://graph.facebook.com/v18.0/me?access_token=${tokenData.access_token}&fields=id,name,email`)
    const userData = await userResponse.json()

    // Get user's pages
    const pagesResponse = await fetch(`https://graph.facebook.com/v18.0/me/accounts?access_token=${tokenData.access_token}&fields=id,name,access_token`)
    const pagesData = await pagesResponse.json()

    // Store user session and pages
    // TODO: Implement actual session storage
    // await storeUserSession(userData, tokenData.access_token, pagesData.data)

    // For development, log the data
    console.log('Facebook OAuth Success:', {
      user: userData,
      pages: pagesData.data?.length || 0
    })

    // Redirect to main page with success
    return NextResponse.redirect(`${process.env.NEXTAUTH_URL}?facebook_connected=true`)

  } catch (error) {
    console.error('Facebook OAuth callback error:', error)
    return NextResponse.redirect(`${process.env.NEXTAUTH_URL}?error=oauth_failed`)
  }
}
