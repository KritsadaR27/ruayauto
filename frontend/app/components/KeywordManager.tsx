'use client'

import { useState, useEffect } from 'react'

interface Keyword {
  id: number
  keyword: string
  response: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export default function KeywordManager() {
  const [keywords, setKeywords] = useState<Keyword[]>([])
  const [newKeyword, setNewKeyword] = useState('')
  const [newResponse, setNewResponse] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Fetch keywords using Next.js 15 async patterns
  const fetchKeywords = async () => {
    try {
      setLoading(true)
      const response = await fetch('/api/keywords', {
        // Next.js 15 enhanced caching
        next: { revalidate: 30 }
      })
      
      if (!response.ok) {
        throw new Error('Failed to fetch keywords')
      }
      
      const data = await response.json()
      setKeywords(data.pairs || [])
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchKeywords()
  }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!newKeyword.trim() || !newResponse.trim()) {
      setError('Both keyword and response are required')
      return
    }

    try {
      setLoading(true)
      const response = await fetch('/api/keywords', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          pairs: [
            {
              keyword: newKeyword.trim(),
              response: newResponse.trim()
            }
          ]
        })
      })

      if (!response.ok) {
        throw new Error('Failed to save keyword')
      }

      // Reset form
      setNewKeyword('')
      setNewResponse('')
      setError(null)
      
      // Refresh keywords list
      await fetchKeywords()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save keyword')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-6">
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
          {error}
        </div>
      )}

      {/* Add new keyword form */}
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label htmlFor="keyword" className="block text-sm font-medium text-gray-700 mb-2">
              Keyword
            </label>
            <input
              type="text"
              id="keyword"
              value={newKeyword}
              onChange={(e) => setNewKeyword(e.target.value)}
              placeholder="Enter keyword (e.g., สวัสดี, hello)"
              className="input-field"
              disabled={loading}
            />
          </div>
          <div>
            <label htmlFor="response" className="block text-sm font-medium text-gray-700 mb-2">
              Auto Response
            </label>
            <input
              type="text"
              id="response"
              value={newResponse}
              onChange={(e) => setNewResponse(e.target.value)}
              placeholder="Enter auto reply message"
              className="input-field"
              disabled={loading}
            />
          </div>
        </div>
        <button
          type="submit"
          disabled={loading}
          className="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? 'Adding...' : 'Add Keyword'}
        </button>
      </form>

      {/* Keywords list */}
      <div>
        <h3 className="text-lg font-medium text-gray-900 mb-4">
          Current Keywords ({keywords.length})
        </h3>
        
        {loading && keywords.length === 0 ? (
          <div className="text-center py-8">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
            <p className="text-gray-500 mt-2">Loading keywords...</p>
          </div>
        ) : keywords.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            No keywords configured yet. Add your first keyword above.
          </div>
        ) : (
          <div className="space-y-3">
            {keywords.map((keyword) => (
              <div
                key={keyword.id}
                className="flex items-center justify-between p-4 bg-gray-50 rounded-lg border border-gray-200"
              >
                <div className="flex-1">
                  <div className="flex items-center space-x-3">
                    <span className="font-medium text-gray-900">
                      {keyword.keyword}
                    </span>
                    <span className={`px-2 py-1 text-xs rounded-full ${
                      keyword.is_active 
                        ? 'bg-green-100 text-green-800' 
                        : 'bg-gray-100 text-gray-800'
                    }`}>
                      {keyword.is_active ? 'Active' : 'Inactive'}
                    </span>
                  </div>
                  <p className="text-gray-600 mt-1">{keyword.response}</p>
                  <p className="text-xs text-gray-400 mt-2">
                    Created: {new Date(keyword.created_at).toLocaleDateString()}
                  </p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
