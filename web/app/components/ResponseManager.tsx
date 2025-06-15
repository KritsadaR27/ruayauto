'use client'

import React, { useState, useEffect } from 'react'
import { Response } from '../types/rule'

// Icons
const PlusIcon = () => (
  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
  </svg>
)

const TrashIcon = () => (
  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
  </svg>
)

const PhotoIcon = () => (
  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
  </svg>
)

const VideoIcon = () => (
  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
  </svg>
)

const TextIcon = () => (
  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16m-7 6h7" />
  </svg>
)

const CompositeIcon = () => (
  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
  </svg>
)

const WeightIcon = () => (
  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 6l3 1m0 0l-3 9a5.002 5.002 0 006.001 0M6 7l3 9M6 7l6-2m6 2l3-1m-3 1l-3 9a5.002 5.002 0 006.001 0M18 7l3 9m-3-9l-6-2m0-2v2m0 16V5m0 16H9m3 0h3" />
  </svg>
)

interface ResponseManagerProps {
  ruleId: number
  responses: Response[]
  onUpdateResponses: (responses: Response[]) => void
}

interface RuleResponse {
  id: number
  rule_id: number
  response_text: string
  response_type: 'text' | 'image' | 'video' | 'composite'
  media_url?: string
  has_media: boolean
  media_description?: string
  weight: number
  is_active: boolean
}

export default function ResponseManager({ ruleId, responses, onUpdateResponses }: ResponseManagerProps) {
  const [ruleResponses, setRuleResponses] = useState<RuleResponse[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Load responses from backend
  useEffect(() => {
    loadRuleResponses()
  }, [ruleId])

  const loadRuleResponses = async () => {
    try {
      setLoading(true)
      const response = await fetch(`/api/rules/${ruleId}/responses`)
      const data = await response.json()
      
      if (data.success && data.data) {
        setRuleResponses(data.data)
      } else {
        setRuleResponses([])
      }
    } catch (err) {
      console.error('Failed to load rule responses:', err)
      setError('Failed to load responses')
    } finally {
      setLoading(false)
    }
  }

  const addResponse = async () => {
    const newResponse: Partial<RuleResponse> = {
      response_text: '',
      response_type: 'text',
      has_media: false,
      weight: 10,
      is_active: true
    }

    try {
      const response = await fetch(`/api/rules/${ruleId}/responses`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newResponse)
      })

      if (response.ok) {
        await loadRuleResponses()
      } else {
        setError('Failed to add response')
      }
    } catch (err) {
      console.error('Failed to add response:', err)
      setError('Failed to add response')
    }
  }

  const updateResponse = async (responseId: number, field: keyof RuleResponse, value: any) => {
    try {
      const response = await fetch(`/api/rules/${ruleId}/responses/${responseId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ [field]: value })
      })

      if (response.ok) {
        // Update local state immediately for better UX
        setRuleResponses(prev => 
          prev.map(resp => 
            resp.id === responseId 
              ? { ...resp, [field]: value }
              : resp
          )
        )
      } else {
        setError('Failed to update response')
      }
    } catch (err) {
      console.error('Failed to update response:', err)
      setError('Failed to update response')
    }
  }

  const deleteResponse = async (responseId: number) => {
    try {
      const response = await fetch(`/api/rules/${ruleId}/responses/${responseId}`, {
        method: 'DELETE'
      })

      if (response.ok) {
        setRuleResponses(prev => prev.filter(resp => resp.id !== responseId))
      } else {
        setError('Failed to delete response')
      }
    } catch (err) {
      console.error('Failed to delete response:', err)
      setError('Failed to delete response')
    }
  }

  const getResponseTypeIcon = (type: string) => {
    switch (type) {
      case 'image': return <PhotoIcon />
      case 'video': return <VideoIcon />
      case 'composite': return <CompositeIcon />
      default: return <TextIcon />
    }
  }

  const getWeightColor = (weight: number) => {
    if (weight >= 15) return 'text-green-600 bg-green-50'
    if (weight >= 10) return 'text-blue-600 bg-blue-50'
    if (weight >= 5) return 'text-yellow-600 bg-yellow-50'
    return 'text-gray-600 bg-gray-50'
  }

  if (loading) {
    return (
      <div className="bg-white rounded-lg border p-4">
        <div className="animate-pulse">
          <div className="h-4 bg-gray-200 rounded w-32 mb-3"></div>
          <div className="space-y-2">
            <div className="h-12 bg-gray-200 rounded"></div>
            <div className="h-12 bg-gray-200 rounded"></div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="bg-white rounded-lg border">
      <div className="p-4 border-b bg-gray-50">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <WeightIcon />
            <h3 className="font-medium text-gray-900">Response Management</h3>
            <span className="text-sm text-gray-500">({ruleResponses.length} responses)</span>
          </div>
          <button
            onClick={addResponse}
            className="flex items-center space-x-1 px-3 py-1.5 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors text-sm"
          >
            <PlusIcon />
            <span>Add Response</span>
          </button>
        </div>
      </div>

      <div className="p-4">
        {error && (
          <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-md">
            <p className="text-red-700 text-sm">{error}</p>
            <button 
              onClick={() => setError(null)}
              className="text-red-600 hover:text-red-800 text-sm underline"
            >
              Dismiss
            </button>
          </div>
        )}

        {ruleResponses.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            <WeightIcon />
            <p className="mt-2">No responses yet. Add your first response!</p>
          </div>
        ) : (
          <div className="space-y-4">
            {ruleResponses
              .sort((a, b) => b.weight - a.weight) // Sort by weight (highest first)
              .map((response, index) => (
              <div key={response.id} className="border rounded-lg p-4 bg-gray-50">
                <div className="flex items-start space-x-3">
                  {/* Response Type Indicator */}
                  <div className="flex-shrink-0">
                    <div className="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center">
                      {getResponseTypeIcon(response.response_type)}
                    </div>
                  </div>

                  <div className="flex-1 min-w-0">
                    {/* Response Type & Weight Controls */}
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center space-x-2">
                        <select
                          value={response.response_type}
                          onChange={(e) => updateResponse(response.id, 'response_type', e.target.value)}
                          className="text-sm border border-gray-300 rounded px-2 py-1"
                        >
                          <option value="text">Text Only</option>
                          <option value="composite">Text + Media</option>
                          <option value="image">Image Only</option>
                          <option value="video">Video Only</option>
                        </select>

                        <div className="flex items-center space-x-1">
                          <WeightIcon />
                          <input
                            type="number"
                            min="1"
                            max="100"
                            value={response.weight}
                            onChange={(e) => updateResponse(response.id, 'weight', parseInt(e.target.value))}
                            className={`w-16 text-sm border border-gray-300 rounded px-2 py-1 text-center ${getWeightColor(response.weight)}`}
                          />
                          <span className="text-xs text-gray-500">weight</span>
                        </div>
                      </div>

                      <div className="flex items-center space-x-2">
                        <label className="flex items-center space-x-1">
                          <input
                            type="checkbox"
                            checked={response.is_active}
                            onChange={(e) => updateResponse(response.id, 'is_active', e.target.checked)}
                            className="rounded"
                          />
                          <span className="text-sm text-gray-600">Active</span>
                        </label>
                        
                        <button
                          onClick={() => deleteResponse(response.id)}
                          className="p-1 text-red-600 hover:text-red-800 hover:bg-red-50 rounded"
                        >
                          <TrashIcon />
                        </button>
                      </div>
                    </div>

                    {/* Response Content */}
                    <div className="space-y-2">
                      <textarea
                        value={response.response_text}
                        onChange={(e) => updateResponse(response.id, 'response_text', e.target.value)}
                        placeholder="Enter response text..."
                        className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm resize-none"
                        rows={2}
                      />

                      {/* Media Fields for composite, image, or video types */}
                      {(response.response_type === 'composite' || response.response_type === 'image' || response.response_type === 'video') && (
                        <div className="space-y-2">
                          <input
                            type="url"
                            value={response.media_url || ''}
                            onChange={(e) => {
                              updateResponse(response.id, 'media_url', e.target.value)
                              updateResponse(response.id, 'has_media', e.target.value.length > 0)
                            }}
                            placeholder={`Enter ${response.response_type === 'composite' ? 'media' : response.response_type} URL...`}
                            className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                          />
                          <input
                            type="text"
                            value={response.media_description || ''}
                            onChange={(e) => updateResponse(response.id, 'media_description', e.target.value)}
                            placeholder="Media description (optional)"
                            className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                          />
                        </div>
                      )}
                    </div>

                    {/* Weight Indicator */}
                    <div className="mt-2 flex items-center justify-between text-xs text-gray-500">
                      <div className="flex items-center space-x-2">
                        <span>Probability: {((response.weight / (ruleResponses.reduce((sum, r) => sum + r.weight, 0))) * 100).toFixed(1)}%</span>
                        <span className={`px-2 py-0.5 rounded-full text-xs ${getWeightColor(response.weight)}`}>
                          Weight: {response.weight}
                        </span>
                      </div>
                      <span>#{index + 1}</span>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Total Weight Summary */}
        {ruleResponses.length > 0 && (
          <div className="mt-4 p-3 bg-blue-50 border border-blue-200 rounded-md">
            <div className="flex items-center justify-between text-sm">
              <span className="text-blue-800">Total Weight: {ruleResponses.reduce((sum, r) => sum + r.weight, 0)}</span>
              <span className="text-blue-600">Active Responses: {ruleResponses.filter(r => r.is_active).length}</span>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
