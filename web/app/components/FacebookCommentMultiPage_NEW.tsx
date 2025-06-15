'use client'

import React, { useState, useEffect } from 'react'
import RuleCard from './RuleCard'

// Custom Icon Components
const PlusIcon = () => (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
    </svg>
)

const SaveIcon = () => (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4" />
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
    hasManuallyEditedTitle?: boolean
}

const FacebookCommentMultiPage = () => {
    // Facebook Pages with connection status
    const [connectedPages, setConnectedPages] = useState<ConnectedPage[]>([
        { id: 'fb1', name: 'ร้านขายของออนไลน์', pageId: '123456', connected: true, enabled: true },
        { id: 'fb2', name: 'แฟนเพจสินค้าแฮนด์เมด', pageId: '789012', connected: true, enabled: true },
        { id: 'fb3', name: 'ร้านอาหารสุขภาพ', pageId: '345678', connected: true, enabled: true },
        { id: 'fb4', name: 'ร้านเครื่องสำอาง', pageId: '901234', connected: false, enabled: false }
    ])

    // Initialize with empty array and load from database
    const [rules, setRules] = useState<Rule[]>([])
    const [loading, setLoading] = useState(true)
    const [saving, setSaving] = useState(false)

    // Load data from API when component mounts
    useEffect(() => {
        const loadRulesFromAPI = async () => {
            try {
                console.log('🚀 Loading keywords from API...')
                setLoading(true)
                const response = await fetch(`/api/rules?_t=${Date.now()}`, {
                    cache: 'no-store',
                    headers: {
                        'Cache-Control': 'no-cache'
                    }
                })
                if (!response.ok) {
                    throw new Error('Failed to load keywords')
                }

                const result = await response.json()
                console.log('📦 API Response:', result)
                if (result.success && result.data?.data?.pairs) {
                    console.log('✅ Found pairs data:', result.data.data.pairs)
                    // Transform database format to frontend format
                    const transformedRules: Rule[] = result.data.data.pairs.map((item: any) => ({
                        id: item.id,
                        name: item.keyword, // Use keyword as name for now
                        keywords: [item.keyword], // Single keyword for now
                        responses: [{ text: item.response }], // Single response for now
                        enabled: item.is_active,
                        expanded: false,
                        selectedPages: ['fb1', 'fb2'], // Default pages
                        hideAfterReply: false,
                        sendToInbox: false,
                        inboxMessage: '',
                        inboxImage: undefined,
                        hasManuallyEditedTitle: false
                    }))
                    console.log('🔄 Transformed rules:', transformedRules)
                    setRules(transformedRules)
                } else {
                    console.log('⚠️ No data found or invalid format')
                    setRules([])
                }
            } catch (error) {
                console.error('Error loading rules:', error)
                setRules([])
            } finally {
                setLoading(false)
            }
        }

        loadRulesFromAPI()
    }, [])

    // Save all rules to database - called when user clicks save button
    const saveAllRules = async () => {
        setSaving(true)
        try {
            console.log('💾 Saving all rules to database...')

            for (const rule of rules) {
                // Skip empty rules
                if (!rule.keywords.length || !rule.responses.length || !rule.responses[0].text.trim()) {
                    continue
                }

                const isNewRule = rule.id > 1000000000000 // Timestamp-based IDs are much larger

                if (isNewRule) {
                    // For new rules, use POST to create
                    const response = await fetch('/api/rules', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({
                            pairs: [{
                                keyword: rule.keywords[0] || rule.name,
                                response: rule.responses[0]?.text || '',
                            }]
                        })
                    })

                    if (!response.ok) {
                        throw new Error(`Failed to save new rule: ${rule.name}`)
                    }

                    const result = await response.json()
                    if (!result.success) {
                        throw new Error(result.error || `Failed to save new rule: ${rule.name}`)
                    }
                } else {
                    // For existing rules, use PUT to update
                    const response = await fetch(`/api/rules/${rule.id}`, {
                        method: 'PUT',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({
                            keyword: rule.keywords[0] || rule.name,
                            response: rule.responses[0]?.text || '',
                            is_active: rule.enabled,
                            priority: 1
                        })
                    })

                    if (!response.ok) {
                        throw new Error(`Failed to update rule: ${rule.name}`)
                    }

                    const result = await response.json()
                    if (!result.success) {
                        throw new Error(result.error || `Failed to update rule: ${rule.name}`)
                    }
                }
            }

            alert('✅ บันทึกข้อมูลเรียบร้อยแล้ว!')
            console.log('✅ All rules saved successfully')
        } catch (error) {
            console.error('❌ Error saving rules:', error)
            alert(`❌ เกิดข้อผิดพลาด: ${error}`)
        } finally {
            setSaving(false)
        }
    }

    // Simple update functions - no auto-save
    const updateRule = (id: number, field: keyof Rule, value: any) => {
        setRules(prevRules =>
            prevRules.map(rule =>
                rule.id === id ? { ...rule, [field]: value } : rule
            )
        )
    }

    const updateResponse = (ruleId: number, responseIdx: number, field: keyof Response, value: any) => {
        setRules(prevRules =>
            prevRules.map(rule => {
                if (rule.id === ruleId) {
                    const newResponses = [...rule.responses]
                    newResponses[responseIdx] = { ...newResponses[responseIdx], [field]: value }
                    return { ...rule, responses: newResponses }
                }
                return rule
            })
        )
    }

    const addRule = () => {
        const newRule: Rule = {
            id: Date.now(),
            name: `รายการที่ ${rules.length + 1}`,
            keywords: [],
            responses: [{ text: '' }],
            enabled: true,
            expanded: true,
            selectedPages: [],
            hideAfterReply: false,
            sendToInbox: false,
            inboxMessage: '',
            inboxImage: undefined,
            hasManuallyEditedTitle: false
        }

        setRules(prevRules => [...prevRules, newRule])
    }

    const deleteRule = (id: number) => {
        setRules(prevRules => prevRules.filter(rule => rule.id !== id))
    }

    const toggleRule = (id: number) => {
        setRules(prevRules =>
            prevRules.map(rule =>
                rule.id === id ? { ...rule, enabled: !rule.enabled } : rule
            )
        )
    }

    const toggleExpand = (id: number) => {
        setRules(prevRules =>
            prevRules.map(rule =>
                rule.id === id ? { ...rule, expanded: !rule.expanded } : rule
            )
        )
    }

    const addResponse = (ruleId: number) => {
        setRules(prevRules =>
            prevRules.map(rule => {
                if (rule.id === ruleId) {
                    return { ...rule, responses: [...rule.responses, { text: '' }] }
                }
                return rule
            })
        )
    }

    const togglePageSelection = (ruleId: number, pageId: string) => {
        setRules(prevRules =>
            prevRules.map(rule => {
                if (rule.id === ruleId) {
                    const selectedPages = rule.selectedPages.includes(pageId)
                        ? rule.selectedPages.filter(id => id !== pageId)
                        : [...rule.selectedPages, pageId]
                    return { ...rule, selectedPages }
                }
                return rule
            })
        )
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

    return (
        <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 p-4">
            <div className="max-w-6xl mx-auto">
                {/* Header with Save Button */}
                <div className="bg-white rounded-2xl shadow-xl mb-8 overflow-hidden">
                    <div className="bg-gradient-to-r from-blue-600 to-indigo-600 p-8">
                        <div className="flex items-center justify-between">
                            <div>
                                <h1 className="text-3xl font-bold text-white mb-2">
                                    ระบบตอบกลับอัตโนมัติ Facebook Comments ✨
                                </h1>
                                <p className="text-blue-100">
                                    ตั้งค่าคำตอบอัตโนมัติเมื่อมีคนคอมเมนต์ในเพจของคุณ • ไม่มี Auto-Save แล้ว กดปุ่มบันทึกเพื่อเซฟ!
                                </p>
                            </div>
                            <button
                                onClick={saveAllRules}
                                disabled={saving}
                                className={`px-6 py-3 rounded-xl font-medium transition-all duration-200 shadow-lg hover:shadow-xl flex items-center space-x-2 ${saving
                                        ? 'bg-gray-400 cursor-not-allowed'
                                        : 'bg-green-500 hover:bg-green-600 text-white'
                                    }`}
                            >
                                <SaveIcon />
                                <span>{saving ? 'กำลังบันทึก...' : 'บันทึกทั้งหมด'}</span>
                            </button>
                        </div>
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
                            <p className="text-sm">เริ่มต้นสร้างกฎใหม่โดยคลิกปุ่ม "เพิ่มกฎใหม่"</p>
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

                {/* Save Button Again */}
                <div className="text-center mb-6">
                    <button
                        onClick={saveAllRules}
                        disabled={saving}
                        className={`px-8 py-3 rounded-xl font-medium transition-all duration-200 shadow-lg hover:shadow-xl flex items-center space-x-2 mx-auto ${saving
                                ? 'bg-gray-400 cursor-not-allowed'
                                : 'bg-gradient-to-r from-green-500 to-emerald-500 hover:from-green-600 hover:to-emerald-600 text-white'
                            }`}
                    >
                        <SaveIcon />
                        <span>{saving ? 'กำลังบันทึก...' : 'บันทึกการตั้งค่า'}</span>
                    </button>
                </div>

                {/* Note */}
                <div className="mt-6 p-4 bg-yellow-50 rounded-lg text-center border border-yellow-200">
                    <p className="text-sm text-yellow-800">
                        ⚠️ <strong>หมายเหตุสำคัญ:</strong> ระบบไม่ได้บันทึกอัตโนมัติแล้ว!
                        กรุณากดปุ่ม "บันทึกทั้งหมด" เพื่อเซฟข้อมูลลงฐานข้อมูล
                    </p>
                </div>
            </div>
        </div>
    )
}

export default FacebookCommentMultiPage
