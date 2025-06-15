'use client'

import React, { useState, useEffect } from 'react'

// Icons
const FacebookIcon = () => (
  <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
    <path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"/>
  </svg>
)

const CheckIcon = () => (
  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
  </svg>
)

const XIcon = () => (
  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
  </svg>
)

const RefreshIcon = () => (
  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
  </svg>
)

interface FacebookPage {
  id: string
  name: string
  access_token: string
  connected: boolean
  expires_at?: string
  last_connected?: string
  status: 'active' | 'expired' | 'error'
}

interface FacebookAuthProps {
  onPagesUpdate: (pages: FacebookPage[]) => void
}

export default function FacebookAuth({ onPagesUpdate }: FacebookAuthProps) {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [user, setUser] = useState<any>(null)
  const [pages, setPages] = useState<FacebookPage[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Check authentication status on load
  useEffect(() => {
    checkAuthStatus()
  }, [])

  const checkAuthStatus = async () => {
    try {
      setError(null)
      const response = await fetch('/api/auth/facebook/status')
      if (response.ok) {
        const data = await response.json()
        setIsAuthenticated(data.authenticated || false)
        setUser(data.user || null)
        if (data.authenticated) {
          loadPages()
        } else {
          setPages([])
          onPagesUpdate([])
        }
      } else {
        console.error('Failed to check auth status')
        setIsAuthenticated(false)
        setUser(null)
        setPages([])
        onPagesUpdate([])
      }
    } catch (error) {
      console.error('Failed to check auth status:', error)
      setIsAuthenticated(false)
      setUser(null)
      setPages([])
      onPagesUpdate([])
    }
  }

  const loadPages = async () => {
    try {
      setLoading(true)
      setError(null)
      const response = await fetch('/api/facebook/pages')
      
      if (response.ok) {
        const data = await response.json()
        if (data.success) {
          setPages(data.pages || [])
          onPagesUpdate(data.pages || [])
        } else {
          setError(data.error || 'Failed to load pages')
          setPages([])
          onPagesUpdate([])
        }
      } else if (response.status === 401) {
        setError('Please authenticate with Facebook first')
        setPages([])
        onPagesUpdate([])
      } else {
        const errorData = await response.json().catch(() => ({}))
        setError(errorData.error || 'Failed to load pages')
        setPages([])
        onPagesUpdate([])
      }
    } catch (error) {
      console.error('Failed to load pages:', error)
      setError('Unable to connect to server')
      setPages([])
      onPagesUpdate([])
    } finally {
      setLoading(false)
    }
  }

  const handleLogin = () => {
    const params = new URLSearchParams({
      client_id: process.env.NEXT_PUBLIC_FACEBOOK_APP_ID || '',
      redirect_uri: `${window.location.origin}/api/auth/facebook/callback`,
      scope: 'pages_show_list,pages_read_engagement,pages_manage_metadata',
      response_type: 'code',
      state: Math.random().toString(36).substring(7)
    })

    window.location.href = `https://www.facebook.com/v18.0/dialog/oauth?${params.toString()}`
  }

  const handleLogout = async () => {
    try {
      setLoading(true)
      const response = await fetch('/api/auth/facebook/logout', {
        method: 'POST'
      })
      
      if (response.ok) {
        setIsAuthenticated(false)
        setUser(null)
        setPages([])
        onPagesUpdate([])
      }
    } catch (error) {
      console.error('Logout failed:', error)
      setError('Logout failed')
    } finally {
      setLoading(false)
    }
  }

  const connectPage = async (pageId: string) => {
    try {
      setLoading(true)
      setError(null)
      const response = await fetch('/api/facebook/pages/connect', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ pageId })
      })

      if (response.ok) {
        const data = await response.json()
        if (data.success) {
          loadPages() // Reload pages
        } else {
          setError(data.error || 'Failed to connect page')
        }
      } else if (response.status === 401) {
        setError('Please authenticate with Facebook first')
      } else {
        const errorData = await response.json().catch(() => ({}))
        setError(errorData.error || 'Failed to connect page')
      }
    } catch (error) {
      console.error('Failed to connect page:', error)
      setError('Unable to connect to server')
    } finally {
      setLoading(false)
    }
  }

  const disconnectPage = async (pageId: string) => {
    try {
      setLoading(true)
      setError(null)
      const response = await fetch(`/api/facebook/pages/${pageId}`, {
        method: 'DELETE'
      })

      if (response.ok) {
        const data = await response.json()
        if (data.success) {
          loadPages() // Reload pages
        } else {
          setError(data.error || 'Failed to disconnect page')
        }
      } else if (response.status === 401) {
        setError('Please authenticate with Facebook first')
      } else {
        const errorData = await response.json().catch(() => ({}))
        setError(errorData.error || 'Failed to disconnect page')
      }
    } catch (error) {
      console.error('Failed to disconnect page:', error)
      setError('Unable to connect to server')
    } finally {
      setLoading(false)
    }
  }

  const refreshPageToken = async (pageId: string) => {
    try {
      setLoading(true)
      const response = await fetch(`/api/facebook/pages/${pageId}/refresh`, {
        method: 'POST'
      })

      if (response.ok) {
        loadPages() // Reload pages
      } else {
        setError('Failed to refresh token')
      }
    } catch (error) {
      console.error('Failed to refresh token:', error)
      setError('Failed to refresh token')
    } finally {
      setLoading(false)
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'text-green-600 bg-green-100'
      case 'expired': return 'text-yellow-600 bg-yellow-100'
      case 'error': return 'text-red-600 bg-red-100'
      default: return 'text-gray-600 bg-gray-100'
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'active': return <CheckIcon />
      case 'expired': return <RefreshIcon />
      case 'error': return <XIcon />
      default: return null
    }
  }

  return (
    <div className="bg-white rounded-xl shadow-lg p-6">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-xl font-bold text-gray-800 flex items-center">
          <FacebookIcon />
          <span className="ml-2">Facebook Integration</span>
        </h2>
        
        {isAuthenticated && (
          <button
            onClick={handleLogout}
            disabled={loading}
            className="px-4 py-2 text-sm bg-red-500 hover:bg-red-600 text-white rounded-lg disabled:opacity-50 transition-colors"
          >
            Logout
          </button>
        )}
      </div>

      {error && (
        <div className="mb-4 p-3 bg-red-100 border border-red-300 rounded-lg text-red-700 text-sm">
          {error}
          <button 
            onClick={() => setError(null)}
            className="ml-2 text-red-500 hover:text-red-700"
          >
            ✕
          </button>
        </div>
      )}

      {!isAuthenticated ? (
        <div className="text-center py-8">
          <div className="mb-4">
            <FacebookIcon />
            <h3 className="text-lg font-semibold text-gray-700 mb-2">Connect Facebook</h3>
            <p className="text-gray-600 text-sm mb-6">
              เชื่อมต่อกับ Facebook เพื่อจัดการคอมเมนต์อัตโนมัติ
            </p>
          </div>
          
          <button
            onClick={handleLogin}
            disabled={loading}
            className="inline-flex items-center px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg disabled:opacity-50 transition-colors"
          >
            <FacebookIcon />
            <span className="ml-2">Connect with Facebook</span>
          </button>
        </div>
      ) : (
        <div className="space-y-6">
          {/* User Info */}
          <div className="flex items-center p-4 bg-blue-50 rounded-lg">
            <div className="w-10 h-10 bg-blue-500 rounded-full flex items-center justify-center text-white font-bold">
              {user?.name?.charAt(0) || 'U'}
            </div>
            <div className="ml-3">
              <p className="font-medium text-gray-800">{user?.name || 'Facebook User'}</p>
              <p className="text-sm text-gray-600">Connected to Facebook</p>
            </div>
          </div>

          {/* Pages Management */}
          <div>
            <h3 className="text-lg font-semibold text-gray-800 mb-4">Manage Pages</h3>
            
            {loading ? (
              <div className="text-center py-8 text-gray-500">
                <RefreshIcon />
                <span className="ml-2">Loading pages...</span>
              </div>
            ) : pages.length === 0 ? (
              <div className="text-center py-8">
                <div className="text-gray-500 mb-4">
                  <p className="text-lg mb-2">📄 No Facebook pages found</p>
                  <p className="text-sm">This could happen because:</p>
                  <ul className="text-xs mt-2 space-y-1 text-left max-w-xs mx-auto">
                    <li>• You haven't connected to Facebook yet</li>
                    <li>• You don't have admin access to any pages</li>
                    <li>• Network connection issues</li>
                    <li>• Facebook API is temporarily unavailable</li>
                  </ul>
                </div>
                <div className="space-x-3">
                  <button
                    onClick={loadPages}
                    disabled={loading}
                    className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg transition-colors disabled:opacity-50"
                  >
                    🔄 Reload Pages
                  </button>
                  <button
                    onClick={checkAuthStatus}
                    disabled={loading}
                    className="px-4 py-2 bg-gray-500 hover:bg-gray-600 text-white rounded-lg transition-colors disabled:opacity-50"
                  >
                    🔐 Check Auth
                  </button>
                </div>
              </div>
            ) : (
              <div className="space-y-3">
                {pages.map((page) => (
                  <div key={page.id} className="flex items-center justify-between p-4 border rounded-lg">
                    <div className="flex items-center space-x-3">
                      <div className="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white text-sm font-bold">
                        {page.name.charAt(0)}
                      </div>
                      <div>
                        <p className="font-medium text-gray-800">{page.name}</p>
                        <p className="text-sm text-gray-600">ID: {page.id}</p>
                        {page.last_connected && (
                          <p className="text-xs text-gray-500">
                            Connected: {new Date(page.last_connected).toLocaleDateString()}
                          </p>
                        )}
                      </div>
                    </div>

                    <div className="flex items-center space-x-3">
                      {/* Status Badge */}
                      <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(page.status)}`}>
                        {getStatusIcon(page.status)}
                        <span className="ml-1 capitalize">{page.status}</span>
                      </span>

                      {/* Action Buttons */}
                      {page.connected ? (
                        <div className="flex space-x-2">
                          {page.status === 'expired' && (
                            <button
                              onClick={() => refreshPageToken(page.id)}
                              disabled={loading}
                              className="px-3 py-1 text-sm bg-yellow-500 hover:bg-yellow-600 text-white rounded transition-colors"
                              title="Refresh Token"
                            >
                              <RefreshIcon />
                            </button>
                          )}
                          <button
                            onClick={() => disconnectPage(page.id)}
                            disabled={loading}
                            className="px-3 py-1 text-sm bg-red-500 hover:bg-red-600 text-white rounded transition-colors"
                          >
                            Disconnect
                          </button>
                        </div>
                      ) : (
                        <button
                          onClick={() => connectPage(page.id)}
                          disabled={loading}
                          className="px-3 py-1 text-sm bg-green-500 hover:bg-green-600 text-white rounded transition-colors"
                        >
                          Connect
                        </button>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}
