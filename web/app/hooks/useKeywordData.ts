'use client'

import { useState, useEffect } from 'react'
import type { RuleData, RulePair, Settings, FilterSettings, FallbackSettings } from '../types/rule'

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

// Remove initialData and use undefined as initial state
export function useRuleData() {
  const [data, setData] = useState<RuleData | undefined>(undefined)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Load data from API
  const loadData = async () => {
    try {
      setLoading(true)
      setError(null)

      const response = await fetch('/api/rules')
      if (!response.ok) {
        throw new Error('Failed to load rule data')
      }

      const result = await response.json()
      if (result.success && result.data) {
        // Transform backend data to frontend format
        const transformedData: RuleData = {
          rules: result.data.pairs || [],
          ...defaultSettings
        }
        setData(transformedData)
      }
    } catch (err) {
      console.error('Error loading rule data:', err)
      setError(err instanceof Error ? err.message : 'Failed to load data')
    } finally {
      setLoading(false)
    }
  }

  // Save data to API
  const saveData = async (newData: RuleData) => {
    try {
      setLoading(true)
      setError(null)

      const response = await fetch('/api/rules', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          pairs: newData.rules.map(rule => ({
            rule: rule.rules.join(','),
            response: rule.responses.map(r => r.text).join('|')
          }))
        })
      })

      if (!response.ok) {
        throw new Error('Failed to save rule data')
      }

      const result = await response.json()
      if (result.success) {
        setData(newData)
      } else {
        throw new Error(result.error || 'Save failed')
      }
    } catch (err) {
      console.error('Error saving rule data:', err)
      setError(err instanceof Error ? err.message : 'Failed to save data')
      throw err
    } finally {
      setLoading(false)
    }
  }

  // Update a specific rule
  const updateRule = (index: number, rule: RulePair) => {
    if (!data) return
    const newData: RuleData = {
      ...data,
      rules: data.rules.map((r, i) => i === index ? rule : r)
    }
    setData(newData)
  }

  // Add a new rule
  const addRule = (rule: RulePair) => {
    if (!data) return
    const newData: RuleData = {
      ...data,
      rules: [...data.rules, rule]
    }
    setData(newData)
  }

  // Delete a rule
  const deleteRule = (index: number) => {
    if (!data) return
    const newData: RuleData = {
      ...data,
      rules: data.rules.filter((_, i) => i !== index)
    }
    setData(newData)
  }

  // Update settings
  const updateSettings = (settings: Settings) => {
    if (!data) return
    const newData: RuleData = {
      ...data,
      ...settings,
      rules: data.rules
    }
    setData(newData)
  }

  // Update filter settings
  const updateFilterSettings = (filterSettings: FilterSettings) => {
    if (!data) return
    const newData: RuleData = {
      ...data,
      filterSettings,
      rules: data.rules
    }
    setData(newData)
  }

  // Update fallback settings
  const updateFallbackSettings = (fallbackSettings: FallbackSettings) => {
    if (!data) return
    const newData: RuleData = {
      ...data,
      fallbackSettings,
      rules: data.rules
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
    updateRule,
    addRule,
    deleteRule,
    updateSettings,
    updateFilterSettings,
    updateFallbackSettings
  }
}
