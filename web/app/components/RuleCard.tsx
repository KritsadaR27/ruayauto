'use client'

import React, { useState, useRef, useEffect } from 'react'
import SimpleResponseManager from './SimpleResponseManager'
import { Response } from '../types/rule'

// Icons
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

const ShuffleIcon = () => (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
    </svg>
)

const InfoIcon = () => (
    <svg
        className="w-4 h-4 inline ml-1"
        fill="currentColor"
        viewBox="0 0 20 20"
    >
        <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
    </svg>
)

const InboxIcon = () => (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
    </svg>
)

const EyeOffIcon = () => (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
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

interface RuleCardProps {
    rule: Rule
    index: number
    connectedPages: ConnectedPage[]
    onUpdateRule: (id: number, field: keyof Rule, value: any) => void
    onUpdateResponse: (ruleId: number, responseIdx: number, field: keyof Response, value: any) => void
    onDeleteRule: (id: number) => void
    onToggleRule: (id: number) => void
    onToggleExpand: (id: number) => void
    onAddResponse: (ruleId: number) => void
    onTogglePageSelection: (ruleId: number, pageId: string) => void
    onImageUpload: (ruleId: number, responseIdx: number, event: React.ChangeEvent<HTMLInputElement>) => void
    onInboxImageUpload: (ruleId: number, event: React.ChangeEvent<HTMLInputElement>) => void
    onSaveRule: (id: number) => void
}

const RuleCard: React.FC<RuleCardProps> = ({
    rule,
    index,
    connectedPages,
    onUpdateRule,
    onUpdateResponse,
    onDeleteRule,
    onToggleRule,
    onToggleExpand,
    onAddResponse,
    onTogglePageSelection,
    onImageUpload,
    onInboxImageUpload,
    onSaveRule
}) => {
    // Auto-Title feature states
    const [isEditingTitle, setIsEditingTitle] = useState(false)
    const [localTitle, setLocalTitle] = useState(rule.name)
    const titleInputRef = useRef<HTMLInputElement>(null)

    // Keyword editing states
    const [editingKeywordIndex, setEditingKeywordIndex] = useState<number | null>(null)
    const editingKeywordRefs = useRef<{ [key: number]: HTMLInputElement | null }>({})

    // Auto-update title based on first keyword (Auto-Title feature)
    useEffect(() => {
        const firstKeyword = rule.keywords.find(k => k.trim() !== '')
        const defaultTitle = `รายการที่ ${index + 1}`

        if (firstKeyword && !rule.hasManuallyEditedTitle && !isEditingTitle) {
            const newTitle = `${firstKeyword} (รายการที่ ${index + 1})`
            if (localTitle !== newTitle) {
                setLocalTitle(newTitle)
                // Temporarily disabled to prevent infinite loop
                // onUpdateRule(rule.id, 'name', newTitle)
            }
        } else if (!firstKeyword && !rule.hasManuallyEditedTitle && !isEditingTitle) {
            if (localTitle !== defaultTitle) {
                setLocalTitle(defaultTitle)
                // Temporarily disabled to prevent infinite loop
                // onUpdateRule(rule.id, 'name', defaultTitle)
            }
        }
    }, [rule.keywords, index, rule.hasManuallyEditedTitle, isEditingTitle, rule.id, localTitle])

    // Auto-focus when editing title
    useEffect(() => {
        if (isEditingTitle && titleInputRef.current) {
            titleInputRef.current.focus()
            titleInputRef.current.select()
        }
    }, [isEditingTitle])

    const handleTitleClick = () => {
        setIsEditingTitle(true)
    }

    const handleTitleChange = (value: string) => {
        setLocalTitle(value)
    }

    const handleTitleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter') {
            saveTitleChanges()
        } else if (e.key === 'Escape') {
            cancelTitleEdit()
        }
    }

    const handleTitleBlur = () => {
        saveTitleChanges()
    }

    const saveTitleChanges = () => {
        setIsEditingTitle(false)
        onUpdateRule(rule.id, 'name', localTitle)
        onUpdateRule(rule.id, 'hasManuallyEditedTitle', true)
    }

    const cancelTitleEdit = () => {
        setIsEditingTitle(false)
        setLocalTitle(rule.name)
    }

    // Keyword editing functions
    const handleKeywordChange = (kidx: number, value: string) => {
        const newKeywords = rule.keywords.map((k, i) => i === kidx ? value : k)
        onUpdateRule(rule.id, 'keywords', newKeywords)
    }

    const startEditingKeyword = (kidx: number) => {
        setEditingKeywordIndex(kidx)
        setTimeout(() => {
            const inputRef = editingKeywordRefs.current[kidx]
            if (inputRef) {
                inputRef.focus()
                inputRef.select()
            }
        }, 10)
    }

    const finishEditingKeyword = () => {
        setEditingKeywordIndex(null)
    }

    const handleEditingKeywordKeyDown = (kidx: number, e: React.KeyboardEvent) => {
        if (e.key === 'Enter') {
            e.preventDefault()
            setEditingKeywordIndex(null)
        } else if (e.key === 'Escape') {
            e.preventDefault()
            // Reset to original keyword value
            const originalKeyword = rule.keywords[kidx]
            const newKeywords = rule.keywords.map((k, i) => i === kidx ? originalKeyword : k)
            onUpdateRule(rule.id, 'keywords', newKeywords)
            setEditingKeywordIndex(null)
        }
    }

    return (
        <div
            className={`bg-white rounded-xl shadow-lg overflow-hidden transition-all duration-300 ${!rule.enabled ? 'opacity-60' : ''
                }`}
            data-rule-index={index}
        >
            {/* Rule Header */}
            <div className="p-4 bg-blue-50 flex items-center justify-between">
                <div className="flex items-center space-x-3">
                    {/* Toggle Switch */}
                    <button
                        onClick={() => onToggleRule(rule.id)}
                        className={`w-12 h-6 rounded-full transition-colors duration-200 ${rule.enabled ? 'bg-green-500' : 'bg-gray-300'
                            }`}
                    >
                        <div className={`w-5 h-5 bg-white rounded-full shadow-md transform transition-transform duration-200 ${rule.enabled ? 'translate-x-6' : 'translate-x-0.5'
                            }`} />
                    </button>

                    {/* Editable Title */}
                    {isEditingTitle ? (
                        <input
                            ref={titleInputRef}
                            type="text"
                            value={localTitle}
                            onChange={(e) => handleTitleChange(e.target.value)}
                            onKeyDown={handleTitleKeyDown}
                            onBlur={handleTitleBlur}
                            className="text-lg font-semibold bg-white border-2 border-blue-500 rounded px-2 py-1 focus:outline-none"
                            placeholder="ชื่อกฎ"
                        />
                    ) : (
                        <div
                            onClick={handleTitleClick}
                            className="text-lg font-semibold cursor-pointer hover:text-blue-600 hover:underline transition-colors px-1"
                            title="คลิกเพื่อแก้ไขชื่อ"
                        >
                            {localTitle}
                        </div>
                    )}

                    {/* Stats */}
                    <span className="text-sm text-gray-500">
                        {rule.keywords.length} คีย์เวิร์ด • {rule.responses.length} คำตอบ • {rule.selectedPages.length} เพจ
                    </span>
                </div>

                <div className="flex items-center space-x-2">
                    <button
                        onClick={() => onSaveRule(rule.id)}
                        className="px-3 py-1 text-sm bg-green-500 hover:bg-green-600 text-white rounded-md transition-colors"
                        title="บันทึกการเปลี่ยนแปลง"
                    >
                        บันทึก
                    </button>

                    <button
                        onClick={() => onToggleExpand(rule.id)}
                        className="p-2 hover:bg-blue-100 rounded-lg transition-colors"
                    >
                        {rule.expanded ? <ChevronUpIcon /> : <ChevronDownIcon />}
                    </button>

                    <button
                        onClick={() => onDeleteRule(rule.id)}
                        className="p-2 text-red-500 hover:bg-red-50 rounded-lg transition-colors"
                    >
                        <TrashIcon />
                    </button>
                </div>
            </div>

            {/* Rule Content */}
            {rule.expanded && (
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
                                        checked={rule.selectedPages.includes(page.id)}
                                        onChange={() => onTogglePageSelection(rule.id, page.id)}
                                        className="w-4 h-4 text-blue-600 rounded"
                                    />
                                    <span className="font-medium">{page.name}</span>
                                    <span className="text-sm text-gray-500">({page.pageId})</span>
                                </label>
                            ))}
                        </div>
                    </div>

                    {/* Keywords Section */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-3">
                            คีย์เวิร์ด (คำที่ต้องการให้ตอบ)
                            <span className="text-gray-500 font-normal text-xs ml-2">คลิกคีย์เวิร์ดเพื่อแก้ไข</span>
                        </label>
                        <div className="flex flex-wrap gap-2 mb-3">
                            {rule.keywords.map((keyword, idx) => (
                                <div
                                    key={idx}
                                    className="inline-flex items-center bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm hover:bg-blue-200 transition-colors"
                                >
                                    {editingKeywordIndex === idx ? (
                                        <input
                                            ref={(ref) => { editingKeywordRefs.current[idx] = ref }}
                                            type="text"
                                            value={keyword}
                                            onChange={(e) => handleKeywordChange(idx, e.target.value)}
                                            onKeyDown={(e) => handleEditingKeywordKeyDown(idx, e)}
                                            onBlur={finishEditingKeyword}
                                            className="bg-transparent text-blue-800 text-sm focus:outline-none min-w-[60px] max-w-[200px]"
                                            autoFocus
                                        />
                                    ) : (
                                        <span
                                            onClick={(e) => {
                                                e.preventDefault()
                                                e.stopPropagation()
                                                startEditingKeyword(idx)
                                            }}
                                            className="cursor-pointer hover:underline select-none"
                                            title="คลิกเพื่อแก้ไข"
                                        >
                                            {keyword}
                                        </span>
                                    )}
                                    <button
                                        onClick={() => {
                                            const newKeywords = rule.keywords.filter((_, i) => i !== idx)
                                            onUpdateRule(rule.id, 'keywords', newKeywords)
                                        }}
                                        className="ml-2 text-blue-600 hover:text-red-600 font-bold text-lg leading-none"
                                        title="ลบคีย์เวิร์ด"
                                    >
                                        ×
                                    </button>
                                </div>
                            ))}
                        </div>
                        <input
                            type="text"
                            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                            onKeyDown={(e) => {
                                if ((e.key === 'Enter' || e.key === 'Tab' || e.key === ',') && e.currentTarget.value.trim()) {
                                    e.preventDefault()
                                    const newKeyword = e.currentTarget.value.trim()
                                    if (!rule.keywords.includes(newKeyword)) {
                                        onUpdateRule(rule.id, 'keywords', [...rule.keywords, newKeyword])
                                    }
                                    e.currentTarget.value = ''
                                }
                            }}
                        />
                    </div>

                    {/* Modern Response Management */}
                    <div className="mt-6">
                        <SimpleResponseManager
                            ruleId={rule.id}
                            responses={rule.responses}
                            onUpdateResponses={(responses: Response[]) => onUpdateRule(rule.id, 'responses', responses)}
                        />
                    </div>

                    {/* Feature Options */}
                    <div className="space-y-3">
                        {/* Hide Comments */}
                        <label className="flex items-center space-x-3 cursor-pointer">
                            <input
                                type="checkbox"
                                className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                checked={rule.hideAfterReply || false}
                                onChange={(e) => onUpdateRule(rule.id, 'hideAfterReply', e.target.checked)}
                            />
                            <span className="text-sm text-gray-700 flex items-center">
                                <EyeOffIcon />
                                <span className="ml-2">ซ่อนคอมเมนต์หลังตอบ</span>
                            </span>
                        </label>

                        {/* Inbox Integration */}
                        <label className="flex items-center space-x-3 cursor-pointer">
                            <input
                                type="checkbox"
                                className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                checked={rule.sendToInbox || false}
                                onChange={(e) => onUpdateRule(rule.id, 'sendToInbox', e.target.checked)}
                            />
                            <span className="text-sm text-gray-700 flex items-center">
                                <InboxIcon />
                                <span className="ml-2">ส่งข้อความเข้า Inbox</span>
                            </span>
                        </label>

                        {/* Inbox Message Form - Shows when sendToInbox is checked */}
                        {rule.sendToInbox && (
                            <div className="ml-7 space-y-3 p-4 bg-blue-50 rounded-lg border border-blue-200">
                                <div className="text-sm font-medium text-blue-800 mb-2">
                                    📩 ข้อความที่จะส่งใน Inbox (ไม่ใส่ = ใช้ข้อความจากคำตอบ)
                                </div>

                                <textarea
                                    value={rule.inboxMessage || ''}
                                    onChange={(e) => onUpdateRule(rule.id, 'inboxMessage', e.target.value)}
                                    className="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-none bg-white"
                                    rows={2}
                                    placeholder="ข้อความที่จะส่งใน Inbox (ไม่ใส่ = ใช้ข้อความจากคำตอบ)"
                                />

                                {/* Inbox Image Upload */}
                                <div className="flex items-center space-x-2">
                                    <input
                                        type="file"
                                        accept="image/*"
                                        onChange={(e) => onInboxImageUpload(rule.id, e)}
                                        className="hidden"
                                        id={`inbox-image-${rule.id}`}
                                    />
                                    <label
                                        htmlFor={`inbox-image-${rule.id}`}
                                        className="inline-flex items-center px-3 py-1 text-sm bg-white hover:bg-gray-50 border border-gray-300 rounded-md cursor-pointer transition-colors"
                                    >
                                        <PhotoIcon />
                                        <span className="ml-1">{rule.inboxImage ? 'เปลี่ยนรูปใน Inbox' : 'เพิ่มรูปใน Inbox'}</span>
                                    </label>
                                    {rule.inboxImage && (
                                        <div className="flex items-center space-x-2">
                                            <img
                                                src={rule.inboxImage}
                                                alt="Inbox preview"
                                                className="h-8 w-8 object-cover rounded border"
                                            />
                                            <button
                                                onClick={() => onUpdateRule(rule.id, 'inboxImage', undefined)}
                                                className="text-red-500 hover:text-red-700 text-xs"
                                            >
                                                ลบรูป
                                            </button>
                                        </div>
                                    )}
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            )}
        </div>
    )
}

export default RuleCard
