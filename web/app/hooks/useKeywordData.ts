'use client'

import { useState, useEffect } from 'react'
import type { KeywordData, Pair, Settings, FilterSettings, FallbackSettings } from '../types/keyword'

const defaultSettings: Settings = {
  defaultResponses: ['ขอบคุณสำหรับข้อความทดสอบครับ', 'ได้รับข้อความทดสอบแล้วครับ'],
  enableDefault: true,
  noTag: false,
  noSticker: false,
  hideCommentsAfterReply: false,
  enableInboxIntegration: false,
  inboxResponse: '',
  inboxImage: undefined,
  filterSettings: {
    skipMentions: false,
    skipStickers: false
  },
  fallbackSettings: {
    enabled: false,
    responses: [{ text: '' }]
  }
}

const initialData: KeywordData = {
  pairs: [
    {
      keywords: ['xyztest123', 'qwertyuiop987', 'zxcvbnm456'],
      responses: [
        { text: 'This is a test response for greeting keywords' },
        { text: 'Test greeting response activated' }
      ]
    },
    {
      keywords: ['abcdefg789', 'mnbvcxz321', 'poiuytrewq654'],
      responses: [
        { text: 'This is a test response for pricing keywords' },
        { text: 'Test pricing response activated' }
      ]
    }
  ],
  ...defaultSettings
}

export function useKeywordData() {
  const [data, setData] = useState<KeywordData>(initialData)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Load data from API
  const loadData = async () => {
    try {
      setLoading(true)
      setError(null)

      const response = await fetch('/api/keywords')
      if (!response.ok) {
        throw new Error('Failed to load keyword data')
      }

      const result = await response.json()
      if (result.success && result.data) {
        // Transform backend data to frontend format
        const transformedData: KeywordData = {
          pairs: result.data.pairs || [],
          ...defaultSettings
        }
        setData(transformedData)
      }
    } catch (err) {
      console.error('Error loading keyword data:', err)
      setError(err instanceof Error ? err.message : 'Failed to load data')
    } finally {
      setLoading(false)
    }
  }

  // Save data to API
  const saveData = async (newData: KeywordData) => {
    try {
      setLoading(true)
      setError(null)

      const response = await fetch('/api/keywords', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          pairs: newData.pairs.map(pair => ({
            keyword: pair.keywords.join(','),
            response: pair.responses.join('|')
          }))
        })
      })

      if (!response.ok) {
        throw new Error('Failed to save keyword data')
      }

      const result = await response.json()
      if (result.success) {
        setData(newData)
      } else {
        throw new Error(result.error || 'Save failed')
      }
    } catch (err) {
      console.error('Error saving keyword data:', err)
      setError(err instanceof Error ? err.message : 'Failed to save data')
      throw err
    } finally {
      setLoading(false)
    }
  }

  // Update a specific pair
  const updatePair = (index: number, pair: Pair) => {
    const newData = {
      ...data,
      pairs: data.pairs.map((p, i) => i === index ? pair : p)
    }
    setData(newData)
  }

  // Add a new pair
  const addPair = (pair: Pair) => {
    const newData = {
      ...data,
      pairs: [...data.pairs, pair]
    }
    setData(newData)
  }

  // Delete a pair
  const deletePair = (index: number) => {
    const newData = {
      ...data,
      pairs: data.pairs.filter((_, i) => i !== index)
    }
    setData(newData)
  }

  // Update settings
  const updateSettings = (settings: Settings) => {
    const newData = {
      ...data,
      ...settings
    }
    setData(newData)
  }

  // Update filter settings
  const updateFilterSettings = (filterSettings: FilterSettings) => {
    const newData = {
      ...data,
      filterSettings
    }
    setData(newData)
  }

  // Update fallback settings
  const updateFallbackSettings = (fallbackSettings: FallbackSettings) => {
    const newData = {
      ...data,
      fallbackSettings
    }
    setData(newData)
  }

  // Load data on mount
  useEffect(() => {
    loadData()
  }, [])

  return {
    data,
    loading,
    error,
    loadData,
    saveData,
    updatePair,
    addPair,
    deletePair,
    updateSettings,
    updateFilterSettings,
    updateFallbackSettings
  }
}
