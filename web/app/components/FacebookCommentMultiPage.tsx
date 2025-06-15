'use client'

import React, { useState, useEffect, useRef } from 'react'
import RuleCard from './RuleCard'

// Custom Icon Components
const PlusIcon = () => (
  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
  </svg>
)

const TrashIcon = () => (
  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
  </svg>
)

const ChevronDownIcon = () => (
  <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
    <path fillRule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clipRule="evenodd" />
  </svg>
)

const ChevronUpIcon = () => (
  <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
    <path fillRule="evenodd" d="M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4 4a1 1 0 011.414 0l4-4a1 1 0 010 1.414z" clipRule="evenodd" />
  </svg>
)

const PhotoIcon = () => (
  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
  </svg>
)

const InboxIcon = () => (
  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a1 1 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
  </svg>
)

const EyeOffIcon = () => (
  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
  </svg>
)

const ShuffleIcon = () => (
  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
  </svg>
)

const InfoIcon = () => (
  <svg
    style={{
      width: '1rem',
      height: '1rem',
      minWidth: '1rem',
      minHeight: '1rem',
      maxWidth: '1rem',
      maxHeight: '1rem',
      display: 'inline-block',
      verticalAlign: 'middle',
      pointerEvents: 'auto',
    }}
    fill="currentColor"
    viewBox="0 0 20 20"
  >
    <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
  </svg>
)

const CheckIcon = () => (
  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
  </svg>
)

// Types
interface ConnectedPage {
  id: string
  name: string
  pageId: string
  connected: boolean
  enabled: boolean
}

interface Response {
  text: string
  image?: string
}

interface Rule {
  id: number
  name: string
  keywords: string[]
  responses: Response[]
  enabled: boolean
  expanded: boolean
  selectedPages: string[]
  hideAfterReply: boolean
  sendToInbox: boolean
  inboxMessage: string
  inboxImage?: string
  hasManuallyEditedTitle?: boolean // Track if user manually edited title
}

interface FallbackRule {
  id: number
  name: string
  enabled: boolean
  expanded: boolean
  selectedPages: string[]
  responses: Response[]
  hideAfterReply: boolean
  sendToInbox: boolean
  inboxMessage: string
  inboxImage?: string
}

const FacebookCommentMultiPage = () => {
  // Facebook Pages with connection status
  const [connectedPages, setConnectedPages] = useState<ConnectedPage[]>([
    { id: 'fb1', name: 'ร้านขายของออนไลน์', pageId: '123456', connected: true, enabled: true },
    { id: 'fb2', name: 'แฟนเพจสินค้าแฮนด์เมด', pageId: '789012', connected: true, enabled: true },
    { id: 'fb3', name: 'ร้านอาหารสุขภาพ', pageId: '345678', connected: true, enabled: true },
    { id: 'fb4', name: 'ร้านเครื่องสำอาง', pageId: '901234', connected: false, enabled: false }
  ])

  const [showPageManager, setShowPageManager] = useState(false)

  // Initialize with empty array and load from database
  const [rules, setRules] = useState<Rule[]>([])
  const [loading, setLoading] = useState(true)

  // Load data from API when component mounts
  useEffect(() => {
    const loadRulesFromAPI = async () => {
      try {
        console.log('🚀 Loading rules from API...')
        setLoading(true)
        const response = await fetch(`/api/rules?_t=${Date.now()}`, {
          cache: 'no-store',
          headers: {
            'Cache-Control': 'no-cache'
          }
        })
        if (!response.ok) {
          throw new Error('Failed to load rules')
        }

        const result = await response.json()
        console.log('📦 API Response:', result)
        if (result.success && result.data?.data?.pairs) {
          console.log('✅ Found pairs data:', result.data.data.pairs)
          
          // Transform the new API format (rules with keywords array) to frontend format
          const transformedRules: Rule[] = result.data.data.pairs.map((item: any, index: number) => {
            // Split comma-separated responses into separate Response objects
            const responsesArray = item.response.split(', ').map((responseText: string) => ({
              text: responseText.trim()
            }))
            
            // Use the keywords array directly
            const keywords = item.keywords || []
            
            return {
              id: item.id,
              name: keywords.length > 1 ? `${keywords[0]} และอีก ${keywords.length - 1} คำ` : (keywords[0] || `รายการที่ ${index + 1}`),
              keywords: keywords, // Multiple keywords from the keywords array
              responses: responsesArray, // Multiple responses from comma-separated string
              enabled: item.is_active,
              expanded: false,
              selectedPages: ['fb1', 'fb2'], // Default pages
              hideAfterReply: item.hide_after_reply || false,
              sendToInbox: item.send_to_inbox || false,
              inboxMessage: item.inbox_message || '',
              inboxImage: item.inbox_image || undefined,
              hasManuallyEditedTitle: false
            }
          })
          console.log('🔄 Transformed rules:', transformedRules)
          setRules(transformedRules)
        } else {
          console.log('⚠️ No data found or invalid format')
          // If no data, keep empty array
          setRules([])
        }
      } catch (error) {
        console.error('Error loading rules:', error)
        // Keep empty array if loading fails
        setRules([])
      } finally {
        setLoading(false)
      }
    }

    loadRulesFromAPI()
  }, [])

  const [fallbackRules, setFallbackRules] = useState<FallbackRule[]>([
    {
      id: 1,
      name: 'คำตอบเมื่อไม่ตรงคีย์เวิร์ด 1',
      enabled: true,
      expanded: true,
      selectedPages: ['fb1', 'fb2'],
      responses: [
        { text: 'ขอบคุณสำหรับข้อความครับ จะรีบตอบกลับโดยเร็วที่สุด' }
      ],
      hideAfterReply: false,
      sendToInbox: false,
      inboxMessage: '',
      inboxImage: undefined
    }
  ])

  const [showAdvanced, setShowAdvanced] = useState(false)

  // Helper function to add response
  const addResponse = (ruleId: number) => {
    setRules(prevRules => prevRules.map(rule => {
      if (rule.id === ruleId) {
        return { ...rule, responses: [...rule.responses, { text: '' }] }
      }
      return rule
    }))
  }

  // Debounce function to prevent too many API calls
  const debounceTimeouts = useRef<Map<number, NodeJS.Timeout>>(new Map())

  const debouncedSave = (rule: Rule, delay: number = 1000) => {
    // Clear existing timeout for this rule
    const existingTimeout = debounceTimeouts.current.get(rule.id)
    if (existingTimeout) {
      clearTimeout(existingTimeout)
    }

    // Set new timeout
    const newTimeout = setTimeout(async () => {
      try {
        if (rule.keywords.length > 0 && rule.responses.length > 0 && rule.responses[0].text.trim()) {
          await saveRuleToDatabase(rule)
        }
      } catch (error) {
        console.error('Failed to save rule:', error)
      }
      debounceTimeouts.current.delete(rule.id)
    }, delay)

    debounceTimeouts.current.set(rule.id, newTimeout)
  }

  // Function to save rule to database
  const saveRuleToDatabase = async (rule: Rule) => {
    try {
      // Check if this is a new rule (timestamp-based ID) or existing rule
      const isNewRule = rule.id > 1000000000000 // Timestamp-based IDs are much larger

      const keywords = rule.keywords.filter(k => k.trim() !== '')
      const responses = rule.responses.filter(r => r.text.trim() !== '')
      const responseText = responses.map(r => r.text).join(', ')
      
      const payload = {
        rule_name: rule.name || '',
        keywords: keywords, // Send keywords as an array
        response: responseText, // All responses as a comma-separated string
        is_active: rule.enabled,
        priority: 1,
        match_type: 'contains',
        hide_after_reply: rule.hideAfterReply || false,
        send_to_inbox: rule.sendToInbox || false,
        inbox_message: rule.inboxMessage || '',
        inbox_image: rule.inboxImage || ''
      }

      if (isNewRule) {
        // Create new rule
        const response = await fetch('/api/rules', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(payload)
        })
        if (!response.ok) {
          throw new Error('Failed to save rule to database')
        }
        const result = await response.json()
        if (!result.success) {
          throw new Error(result.error || 'Failed to save rule')
        }
      } else {
        // Update existing rule
        const response = await fetch(`/api/rules/${rule.id}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(payload)
        })
        if (!response.ok) {
          throw new Error('Failed to update rule in database')
        }
        const result = await response.json()
        if (!result.success) {
          throw new Error(result.error || 'Failed to update rule')
        }
      }
    } catch (error) {
      console.error('Error saving rule to database:', error)
      throw error
    }
  }

  const addRule = async () => {
    let newRuleIndex = 0
    
    setRules(prevRules => {
      newRuleIndex = prevRules.length
      const newRule: Rule = {
        id: Date.now(),
        name: `รายการที่ ${prevRules.length + 1}`, // Start with default title
        keywords: [],
        responses: [{ text: '' }],
        enabled: true,
        expanded: true, // ✅ เปลี่ยนจาก false เป็น true
        selectedPages: [],
        hideAfterReply: false,
        sendToInbox: false,
        inboxMessage: '',
        inboxImage: undefined,
        hasManuallyEditedTitle: false // Track manual editing
      }

      // Add to local state first
      return [...prevRules, newRule]
    })

    // Save to database (will be saved when user adds content)
    // Note: We'll save when user actually enters keywords/responses

    // ✅ ปรับปรุง auto-focus logic
    setTimeout(() => {
      const newCard = document.querySelector(`[data-rule-index="${newRuleIndex}"]`)
      if (newCard) {
        const keywordInput = newCard.querySelector('input[placeholder*="คีย์เวิร์ด"]') as HTMLInputElement
        if (keywordInput) {
          keywordInput.focus()
        }
      }
    }, 200) // เพิ่ม delay เป็น 200ms
  }

  // Manual save function for individual rules
  const saveRule = async (id: number) => {
    const rule = rules.find(r => r.id === id)
    if (!rule) return

    try {
      console.log('💾 Manually saving rule:', id)
      await saveRuleToDatabase(rule)
      console.log('✅ Rule saved successfully')
      
      // Optional: Show success feedback
      // You could add a toast notification here
    } catch (error) {
      console.error('❌ Failed to save rule:', error)
      // Optional: Show error feedback
      // You could add error toast notification here
    }
  }

  const updateRule = async (id: number, field: keyof Rule, value: any) => {
    let updatedRule: Rule | undefined

    // Update local state immediately for responsive UI using functional update
    setRules(prevRules => {
      const newRules = prevRules.map(rule => {
        if (rule.id === id) {
          const updated = { ...rule, [field]: value }
          updatedRule = updated
          return updated
        }
        return rule
      })
      return newRules
    })

    // Use debounced save instead of immediate save
    if (updatedRule) {
      debouncedSave(updatedRule)
    }
  }

  const updateResponse = async (ruleId: number, responseIdx: number, field: keyof Response, value: any) => {
    let updatedRule: Rule | undefined

    // Update local state immediately using functional update
    setRules(prevRules => {
      const newRules = prevRules.map(rule => {
        if (rule.id === ruleId) {
          const newResponses = [...rule.responses]
          newResponses[responseIdx] = { ...newResponses[responseIdx], [field]: value }
          const updated = { ...rule, responses: newResponses }
          updatedRule = updated
          return updated
        }
        return rule
      })
      return newRules
    })

    // Use debounced save instead of immediate save
    if (updatedRule) {
      debouncedSave(updatedRule)
    }
  }

  const deleteRule = async (id: number) => {
    // Remove from local state first using functional update
    setRules(prevRules => prevRules.filter(rule => rule.id !== id))

    // Delete from database (only if it's not a timestamp-based new rule)
    if (id < 1000000000000) {
      try {
        const response = await fetch(`/api/rules/${id}`, {
          method: 'DELETE',
        })

        if (!response.ok) {
          throw new Error('Failed to delete rule from database')
        }

        const result = await response.json()
        if (!result.success) {
          throw new Error(result.error || 'Failed to delete rule')
        }
      } catch (error) {
        console.error('Failed to delete rule:', error)
        // Could show error message and restore the rule
      }
    }
  }

  const toggleRule = async (id: number) => {
    let updatedRule: Rule | undefined

    // Update local state first using functional update
    setRules(prevRules => {
      const newRules = prevRules.map(rule => {
        if (rule.id === id) {
          const updated = { ...rule, enabled: !rule.enabled }
          updatedRule = updated
          return updated
        }
        return rule
      })
      return newRules
    })

    // Save to database (only if it's not a timestamp-based new rule)
    if (id < 1000000000000 && updatedRule) {
      try {
        await saveRuleToDatabase(updatedRule)
      } catch (error) {
        console.error('Failed to toggle rule:', error)
      }
    }
  }

  const toggleExpand = (id: number) => {
    setRules(prevRules => prevRules.map(rule =>
      rule.id === id ? { ...rule, expanded: !rule.expanded } : rule
    ))
  }

  const handleImageUpload = (ruleId: number, responseIdx: number, event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (e) => {
        updateResponse(ruleId, responseIdx, 'image', e.target?.result as string)
      }
      reader.readAsDataURL(file)
    }
  }

  const handleInboxImageUpload = (ruleId: number, event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (e) => {
        updateRule(ruleId, 'inboxImage', e.target?.result as string)
      }
      reader.readAsDataURL(file)
    }
  }

  const togglePageSelection = (ruleId: number, pageId: string) => {
    setRules(prevRules => prevRules.map(rule => {
      if (rule.id === ruleId) {
        const selectedPages = rule.selectedPages.includes(pageId)
          ? rule.selectedPages.filter(id => id !== pageId)
          : [...rule.selectedPages, pageId]
        return { ...rule, selectedPages }
      }
      return rule
    }))
  }

  const addFallbackRule = () => {
    const newFallback: FallbackRule = {
      id: Date.now(),
      name: `คำตอบเมื่อไม่ตรงคีย์เวิร์ด ${fallbackRules.length + 1}`,
      enabled: true,
      expanded: true,
      selectedPages: [],
      responses: [{ text: '' }],
      hideAfterReply: false,
      sendToInbox: false,
      inboxMessage: '',
      inboxImage: undefined
    }
    setFallbackRules([...fallbackRules, newFallback])
  }

  const updateFallbackRule = (id: number, field: keyof FallbackRule, value: any) => {
    setFallbackRules(fallbackRules.map(rule =>
      rule.id === id ? { ...rule, [field]: value } : rule
    ))
  }

  const updateFallbackResponse = (ruleId: number, responseIdx: number, field: keyof Response, value: any) => {
    setFallbackRules(fallbackRules.map(rule => {
      if (rule.id === ruleId) {
        const newResponses = [...rule.responses]
        newResponses[responseIdx] = { ...newResponses[responseIdx], [field]: value }
        return { ...rule, responses: newResponses }
      }
      return rule
    }))
  }

  const deleteFallbackRule = (id: number) => {
    if (fallbackRules.length > 1) {
      setFallbackRules(fallbackRules.filter(rule => rule.id !== id))
    }
  }

  const toggleFallbackPageSelection = (ruleId: number, pageId: string) => {
    setFallbackRules(fallbackRules.map(rule => {
      if (rule.id === ruleId) {
        const selectedPages = rule.selectedPages.includes(pageId)
          ? rule.selectedPages.filter(id => id !== pageId)
          : [...rule.selectedPages, pageId]
        return { ...rule, selectedPages }
      }
      return rule
    }))
  }

  const addFallbackResponse = (ruleId: number) => {
    setFallbackRules(fallbackRules.map(rule => {
      if (rule.id === ruleId) {
        return { ...rule, responses: [...rule.responses, { text: '' }] }
      }
      return rule
    }))
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 p-4">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <div className="bg-white rounded-2xl shadow-xl mb-8 overflow-hidden">
          <div className="bg-gradient-to-r from-blue-600 to-indigo-600 p-8">
            <div className="flex items-center justify-between">
              <div>
                <h1 className="text-3xl font-bold text-white mb-2">
                  ระบบตอบกลับอัตโนมัติ Facebook Comments ✨
                </h1>
                <p className="text-blue-100">
                  ตั้งค่าคำตอบอัตโนมัติเมื่อมีคนคอมเมนต์ในเพจของคุณ • พร้อมฟีเจอร์ Auto-Title ใหม่!
                </p>
              </div>
              <button
                onClick={() => setShowPageManager(true)}
                className="bg-white/20 hover:bg-white/30 text-white px-4 py-2 rounded-lg font-medium transition-colors"
              >
                ⚙️ จัดการเพจ
              </button>
            </div>
          </div>

          {/* Connected Pages Summary */}
          <div className="p-6 bg-gray-50">
            <h3 className="text-sm font-medium text-gray-700 mb-3">เพจที่เชื่อมต่อ:</h3>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {connectedPages.map(page => {
                const enabledRulesCount = rules.filter(rule =>
                  rule.enabled && rule.selectedPages.includes(page.id)
                ).length

                return (
                  <div key={page.id} className="bg-white rounded-lg p-4 border border-gray-200">
                    <div className="flex items-center justify-between mb-2">
                      <span className="font-medium text-sm truncate">{page.name}</span>
                      {page.connected && (
                        <span className={`text-xs px-2 py-1 rounded-full ${page.enabled ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-600'
                          }`}>
                          {page.enabled ? 'เปิด' : 'ปิด'}
                        </span>
                      )}
                    </div>
                    <div className="text-xs text-gray-500">
                      {page.connected ? (
                        <>
                          <span className="text-green-600">● เชื่อมต่อแล้ว</span>
                          {page.enabled && (
                            <div className="mt-1">{enabledRulesCount} กฎที่ใช้งาน</div>
                          )}
                        </>
                      ) : (
                        <span className="text-gray-400">● ยังไม่เชื่อมต่อ</span>
                      )}
                    </div>
                  </div>
                )
              })}
            </div>
          </div>
        </div>

        {/* Page Manager Modal */}
        {showPageManager && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded-2xl shadow-2xl max-w-3xl w-full max-h-[90vh] overflow-hidden">
              <div className="bg-gradient-to-r from-blue-600 to-indigo-600 p-6 text-white">
                <div className="flex items-center justify-between">
                  <h2 className="text-2xl font-bold">จัดการเพจ Facebook</h2>
                  <button
                    onClick={() => setShowPageManager(false)}
                    className="text-white hover:text-gray-200"
                  >
                    ✕
                  </button>
                </div>
              </div>

              <div className="p-6 overflow-y-auto max-h-[calc(90vh-120px)]">
                <div className="space-y-3">
                  {connectedPages.map(page => (
                    <div key={page.id} className={`border rounded-lg p-4 ${page.connected ? 'border-gray-200' : 'border-gray-300 bg-gray-50'
                      }`}>
                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-3">
                          <button
                            onClick={() => {
                              const updatedPages = connectedPages.map(p =>
                                p.id === page.id ? { ...p, enabled: !p.enabled } : p
                              )
                              setConnectedPages(updatedPages)
                            }}
                            disabled={!page.connected}
                            className={`w-12 h-6 rounded-full transition-colors duration-200 ${page.enabled && page.connected ? 'bg-green-500' : 'bg-gray-300'
                              } ${!page.connected ? 'opacity-50 cursor-not-allowed' : ''}`}
                          >
                            <div className={`w-5 h-5 bg-white rounded-full shadow-md transform transition-transform duration-200 ${page.enabled && page.connected ? 'translate-x-6' : 'translate-x-0.5'
                              }`} />
                          </button>

                          <div>
                            <div className="font-medium">{page.name}</div>
                            <div className="text-sm text-gray-500">Page ID: {page.pageId}</div>
                          </div>
                        </div>

                        <div className="flex items-center space-x-2">
                          {page.connected ? (
                            <span className="text-sm text-green-600 flex items-center">
                              <CheckIcon />
                              <span className="ml-1">เชื่อมต่อแล้ว</span>
                            </span>
                          ) : (
                            <button className="text-sm px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600">
                              เชื่อมต่อ
                            </button>
                          )}
                          <button className="text-red-500 hover:text-red-700">
                            <TrashIcon />
                          </button>
                        </div>
                      </div>
                    </div>
                  ))}

                  <button className="w-full py-2 border-2 border-dashed border-gray-300 rounded-lg text-gray-500 hover:border-blue-500 hover:text-blue-500">
                    + เพิ่มเพจใหม่
                  </button>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Fallback Responses */}
        <div className="mb-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-bold text-gray-800 flex items-center">
              💬 คำตอบเมื่อไม่ตรงกับคีย์เวิร์ดใดๆ
            </h2>
            <button
              onClick={addFallbackRule}
              className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center text-sm"
            >
              <PlusIcon />
              <span className="ml-1">เพิ่มชุดคำตอบ</span>
            </button>
          </div>

          <div className="space-y-4">
            {fallbackRules.map((fallbackRule) => (
              <div key={fallbackRule.id} className={`bg-white rounded-xl shadow-lg overflow-hidden transition-all duration-300 ${!fallbackRule.enabled ? 'opacity-60' : ''
                }`}>
                {/* Fallback Rule Header */}
                <div className="p-4 bg-amber-50 flex items-center justify-between">
                  <div className="flex items-center space-x-3">
                    <button
                      onClick={() => updateFallbackRule(fallbackRule.id, 'enabled', !fallbackRule.enabled)}
                      className={`w-12 h-6 rounded-full transition-colors duration-200 ${fallbackRule.enabled ? 'bg-green-500' : 'bg-gray-300'
                        }`}
                    >
                      <div className={`w-5 h-5 bg-white rounded-full shadow-md transform transition-transform duration-200 ${fallbackRule.enabled ? 'translate-x-6' : 'translate-x-0.5'
                        }`} />
                    </button>

                    <input
                      type="text"
                      value={fallbackRule.name}
                      onChange={(e) => updateFallbackRule(fallbackRule.id, 'name', e.target.value)}
                      className="text-lg font-semibold bg-transparent border-b-2 border-transparent hover:border-gray-300 focus:border-amber-500 focus:outline-none px-1"
                      placeholder="ชื่อชุดคำตอบ"
                    />

                    <span className="text-sm text-gray-500">
                      {fallbackRule.selectedPages.length} เพจที่เลือก
                    </span>
                  </div>

                  <div className="flex items-center space-x-2">
                    <button
                      onClick={() => updateFallbackRule(fallbackRule.id, 'expanded', !fallbackRule.expanded)}
                      className="p-2 hover:bg-amber-100 rounded-lg transition-colors"
                    >
                      {fallbackRule.expanded ? <ChevronUpIcon /> : <ChevronDownIcon />}
                    </button>

                    {fallbackRules.length > 1 && (
                      <button
                        onClick={() => deleteFallbackRule(fallbackRule.id)}
                        className="p-2 text-red-500 hover:bg-red-50 rounded-lg transition-colors"
                      >
                        <TrashIcon />
                      </button>
                    )}
                  </div>
                </div>

                {/* Fallback Rule Content - This will be continued in the next part */}
                {fallbackRule.expanded && (
                  <div className="p-6 space-y-6">
                    {/* Page Selection */}
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-3">
                        ใช้กับเพจ:
                      </label>
                      <div className="space-y-2">
                        {connectedPages.filter(page => page.connected && page.enabled).map(page => (
                          <label key={page.id} className="flex items-center space-x-3 p-3 border rounded-lg hover:bg-gray-50 cursor-pointer">
                            <input
                              type="checkbox"
                              checked={fallbackRule.selectedPages.includes(page.id)}
                              onChange={() => toggleFallbackPageSelection(fallbackRule.id, page.id)}
                              className="w-4 h-4 text-amber-600 rounded"
                            />
                            <span className="font-medium">{page.name}</span>
                            <span className="text-sm text-gray-500">({page.pageId})</span>
                          </label>
                        ))}
                      </div>
                    </div>

                    {/* Responses Section */}
                    <div>
                      <div className="flex items-center justify-between mb-3">
                        <label className="block text-sm font-medium text-gray-700">
                          คำตอบเมื่อไม่ตรงคีย์เวิร์ด
                        </label>
                        <div className="flex items-center text-sm text-gray-500">
                          <ShuffleIcon />
                          <span className="ml-1">ระบบจะสุ่มเลือก 1 คำตอบ</span>
                          <InfoIcon />
                        </div>
                      </div>
                      <div className="space-y-3 mb-3">
                        {fallbackRule.responses.map((response, idx) => (
                          <div key={idx} className="border border-amber-200 rounded-lg p-4 space-y-3 bg-amber-50">
                            <div className="flex items-start space-x-2">
                              <div className="flex-1">
                                <textarea
                                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500 resize-none"
                                  value={response.text}
                                  onChange={(e) => updateFallbackResponse(fallbackRule.id, idx, 'text', e.target.value)}
                                  placeholder="ข้อความตอบกลับเมื่อไม่ตรงคีย์เวิร์ด"
                                  rows={3}
                                />
                              </div>
                              {fallbackRule.responses.length > 1 && (
                                <button
                                  onClick={() => {
                                    const newResponses = fallbackRule.responses.filter((_, i) => i !== idx)
                                    updateFallbackRule(fallbackRule.id, 'responses', newResponses)
                                  }}
                                  className="p-1 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors"
                                >
                                  <TrashIcon />
                                </button>
                              )}
                            </div>

                            <div className="flex items-center justify-between">
                              <div className="flex items-center space-x-2">
                                <input
                                  type="file"
                                  accept="image/*"
                                  onChange={(e) => {
                                    const file = e.target.files?.[0]
                                    if (file) {
                                      const reader = new FileReader()
                                      reader.onload = (event) => {
                                        updateFallbackResponse(fallbackRule.id, idx, 'image', event.target?.result as string)
                                      }
                                      reader.readAsDataURL(file)
                                    }
                                  }}
                                  className="hidden"
                                  id={`fallback-image-${fallbackRule.id}-${idx}`}
                                />
                                <label
                                  htmlFor={`fallback-image-${fallbackRule.id}-${idx}`}
                                  className="inline-flex items-center px-3 py-1 text-sm bg-gray-100 hover:bg-gray-200 border border-gray-300 rounded-md cursor-pointer transition-colors"
                                >
                                  <PhotoIcon />
                                  <span className="ml-1">{response.image ? 'เปลี่ยนรูป' : 'เลือกรูป'}</span>
                                </label>
                                {response.image && (
                                  <div className="flex items-center space-x-2">
                                    <img
                                      src={response.image}
                                      alt="Fallback response preview"
                                      className="h-8 w-8 object-cover rounded border"
                                    />
                                    <button
                                      onClick={() => updateFallbackResponse(fallbackRule.id, idx, 'image', undefined)}
                                      className="text-red-500 hover:text-red-700 text-xs"
                                    >
                                      ลบรูป
                                    </button>
                                  </div>
                                )}
                              </div>
                            </div>
                          </div>
                        ))}
                      </div>
                      <div className="mt-3 pt-3 border-t border-gray-100">
                        <button
                          onClick={() => addFallbackResponse(fallbackRule.id)}
                          className="px-4 py-2 bg-amber-500 text-white text-sm rounded hover:bg-amber-600 flex items-center"
                        >
                          <PlusIcon />
                          <span className="ml-1">เพิ่มคำตอบ</span>
                        </button>
                      </div>
                    </div>

                    {/* Feature Options */}
                    <div className="bg-amber-50 border border-amber-200 rounded-lg p-4 space-y-4">
                      <h4 className="text-sm font-semibold text-amber-900 flex items-center">
                        ⚡ ฟีเจอร์เสริม
                      </h4>

                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        {/* Hide Comments */}
                        <label className="flex items-start space-x-3 cursor-pointer bg-white rounded-lg p-3 border border-amber-100 hover:border-amber-300 transition-colors">
                          <input
                            type="checkbox"
                            className="mt-1 rounded border-gray-300 text-amber-600 focus:ring-amber-500"
                            checked={fallbackRule.hideAfterReply || false}
                            onChange={(e) => updateFallbackRule(fallbackRule.id, 'hideAfterReply', e.target.checked)}
                          />
                          <div>
                            <span className="text-sm font-medium text-gray-900">
                              <EyeOffIcon />
                              ซ่อนคอมเมนต์หลังตอบ
                            </span>
                            <p className="text-xs text-gray-500">
                              ซ่อนคอมเมนต์อัตโนมัติหลังจากตอบกลับ
                            </p>
                          </div>
                        </label>

                        {/* Inbox Integration */}
                        <label className="flex items-start space-x-3 cursor-pointer bg-white rounded-lg p-3 border border-amber-100 hover:border-amber-300 transition-colors">
                          <input
                            type="checkbox"
                            className="mt-1 rounded border-gray-300 text-amber-600 focus:ring-amber-500"
                            checked={fallbackRule.sendToInbox || false}
                            onChange={(e) => updateFallbackRule(fallbackRule.id, 'sendToInbox', e.target.checked)}
                          />
                          <div className="flex-1">
                            <span className="text-sm font-medium text-gray-900">
                              <InboxIcon />
                              ส่งไปแชท Inbox
                            </span>
                            <p className="text-xs text-gray-500">
                              ส่งข้อความไปยัง Inbox หลังจากตอบในคอมเมนต์
                            </p>
                          </div>
                        </label>
                      </div>

                      {/* Inbox Settings */}
                      {fallbackRule.sendToInbox && (
                        <div className="bg-white rounded-lg p-4 border border-amber-200 space-y-3">
                          <h5 className="text-sm font-medium text-gray-700">📧 ตั้งค่า Inbox</h5>

                          <div>
                            <label className="block text-xs font-medium text-gray-600 mb-2">
                              ข้อความใน Inbox:
                            </label>
                            <textarea
                              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500 text-sm"
                              value={fallbackRule.inboxMessage || ''}
                              onChange={(e) => updateFallbackRule(fallbackRule.id, 'inboxMessage', e.target.value)}
                              placeholder="ข้อความที่จะส่งใน Inbox หลังจากตอบในคอมเมนต์"
                              rows={2}
                            />
                          </div>

                          <div>
                            <label className="block text-xs font-medium text-gray-600 mb-2">
                              <PhotoIcon />
                              รูปภาพแนบใน Inbox:
                            </label>
                            <div className="flex items-center space-x-2">
                              <input
                                type="file"
                                accept="image/*"
                                onChange={(e) => {
                                  const file = e.target.files?.[0]
                                  if (file) {
                                    const reader = new FileReader()
                                    reader.onload = (event) => {
                                      updateFallbackRule(fallbackRule.id, 'inboxImage', event.target?.result as string)
                                    }
                                    reader.readAsDataURL(file)
                                  }
                                }}
                                className="hidden"
                                id={`fallback-inbox-image-${fallbackRule.id}`}
                              />
                              <label
                                htmlFor={`fallback-inbox-image-${fallbackRule.id}`}
                                className="inline-flex items-center px-3 py-1 text-xs bg-gray-100 hover:bg-gray-200 border border-gray-300 rounded-md cursor-pointer transition-colors"
                              >
                                <PhotoIcon />
                                <span className="ml-1">{fallbackRule.inboxImage ? 'เปลี่ยนรูปใน Inbox' : 'เลือกรูปใน Inbox'}</span>
                              </label>
                              {fallbackRule.inboxImage && (
                                <div className="flex items-center space-x-2">
                                  <img
                                    src={fallbackRule.inboxImage}
                                    alt="Fallback inbox preview"
                                    className="h-8 w-8 object-cover rounded border"
                                  />
                                  <button
                                    onClick={() => updateFallbackRule(fallbackRule.id, 'inboxImage', undefined)}
                                    className="text-red-500 hover:text-red-700 text-xs"
                                  >
                                    ลบ
                                  </button>
                                </div>
                              )}
                            </div>
                          </div>
                        </div>
                      )}
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>

        {/* Rules List */}
        <div className="space-y-4 mb-6">
          {loading ? (
            <div className="flex items-center justify-center py-12">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
              <span className="ml-3 text-gray-600">กำลังโหลดข้อมูล...</span>
            </div>
          ) : rules.length === 0 ? (
            <div className="text-center py-12 text-gray-500">
              <p className="text-lg mb-2">ยังไม่มีกฎการตอบกลับ</p>
              <p className="text-sm">เริ่มต้นสร้างกฎใหม่โดยคลิกปุ่ม "เพิ่มรายการใหม่"</p>
            </div>
          ) : (
            rules.map((rule, index) => (
              <RuleCard
                key={rule.id}
                rule={rule}
                index={index}
                connectedPages={connectedPages}
                onUpdateRule={updateRule}
                onUpdateResponse={updateResponse}
                onDeleteRule={deleteRule}
                onToggleRule={toggleRule}
                onToggleExpand={toggleExpand}
                onAddResponse={addResponse}
                onTogglePageSelection={togglePageSelection}
                onImageUpload={handleImageUpload}
                onInboxImageUpload={handleInboxImageUpload}
                onSaveRule={saveRule}
              />
            ))
          )}
        </div>

        {/* Add New Rule Button */}
        <button
          onClick={addRule}
          className="w-full mb-6 py-4 bg-gradient-to-r from-blue-500 to-indigo-500 text-white rounded-xl font-medium hover:from-blue-600 hover:to-indigo-600 transition-all duration-200 shadow-lg hover:shadow-xl flex items-center justify-center space-x-2"
        >
          <PlusIcon />
          <span>เพิ่มกฎใหม่</span>
        </button>

        {/* Advanced Settings */}
        <div className="bg-white rounded-xl shadow-lg overflow-hidden mb-6">
          <button
            onClick={() => setShowAdvanced(!showAdvanced)}
            className="w-full p-4 text-left flex items-center justify-between hover:bg-gray-50 transition-colors"
          >
            <span className="font-medium text-gray-900">ตั้งค่าขั้นสูง (ใช้กับทุกเพจ)</span>
            {showAdvanced ? <ChevronUpIcon /> : <ChevronDownIcon />}
          </button>

          {showAdvanced && (
            <div className="p-6 border-t border-gray-200 space-y-4">
              <label className="flex items-center space-x-3">
                <input type="checkbox" className="w-4 h-4 text-blue-600 rounded" />
                <span className="text-gray-700">ไม่ตอบคอมเมนต์ที่มีการแท็ก (@)</span>
              </label>
              <label className="flex items-center space-x-3">
                <input type="checkbox" className="w-4 h-4 text-blue-600 rounded" />
                <span className="text-gray-700">ไม่ตอบคอมเมนต์ที่มีสติกเกอร์</span>
              </label>
            </div>
          )}
        </div>

        {/* Save Button */}
        <div className="mt-8 text-center">
          <button className="px-8 py-3 bg-gradient-to-r from-green-500 to-emerald-500 text-white rounded-xl font-medium hover:from-green-600 hover:to-emerald-600 transition-all duration-200 shadow-lg hover:shadow-xl">
            บันทึกการตั้งค่า
          </button>
        </div>

        {/* Note */}
        <div className="mt-6 p-4 bg-blue-50 rounded-lg text-center">
          <p className="text-sm text-blue-800">
            📌 หน้านี้สำหรับตั้งค่าตอบคอมเมนต์ Facebook •
            สำหรับตอบข้อความแชท (Facebook, LINE, TikTok) อยู่ในหน้าตั้งค่าแยกต่างหาก
          </p>
        </div>
      </div>
    </div>
  )
}

export default FacebookCommentMultiPage
