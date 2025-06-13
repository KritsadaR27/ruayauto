'use client'

import { useState, useRef, useEffect } from 'react'
import {
  PlusIcon,
  TrashIcon,
  PhotoIcon,
  EyeSlashIcon,
  InboxIcon,
  ChevronUpIcon,
  ChevronDownIcon,
  PencilIcon,
  CheckIcon,
  XMarkIcon
} from '@heroicons/react/24/outline'
import type { Pair, Response } from '../types/keyword'

// Custom Icon for Shuffle
const ShuffleIcon = ({ className = "w-5 h-5" }: { className?: string }) => (
  <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
  </svg>
)

// InfoIcon: force 1rem size, always local, no className
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

interface Props {
  pair: Pair
  index: number
  onUpdate: (index: number, pair: Pair) => void
  onDelete: (index: number) => void
}

export default function KeywordPairCard({ pair, index, onUpdate, onDelete }: Props) {
  const [editingKeywords, setEditingKeywords] = useState<string[]>(pair.keywords)
  const [editingResponses, setEditingResponses] = useState<Response[]>(() => {
    if (!pair.responses || pair.responses.length === 0) {
      return [{ text: '' }]
    }

    // Handle legacy string[] format
    if (typeof pair.responses[0] === 'string') {
      return (pair.responses as unknown as string[]).map(text => ({ text }))
    }

    // Handle new Response[] format
    return pair.responses as Response[]
  })
  const [newKeyword, setNewKeyword] = useState('')
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false)
  const [isEditingTitle, setIsEditingTitle] = useState(false)
  const [editingTitle, setEditingTitle] = useState(pair.title || `รายการที่ ${index + 1}`)
  const [hasManuallyEditedTitle, setHasManuallyEditedTitle] = useState(pair.hasManuallyEditedTitle || false)
  const [isExpanded, setIsExpanded] = useState(pair.expanded !== false) // Default to expanded
  const [isEnabled, setIsEnabled] = useState(pair.enabled !== false) // Default to enabled
  const [editingKeywordIndex, setEditingKeywordIndex] = useState<number | null>(null)

  const titleInputRef = useRef<HTMLInputElement>(null)
  const keywordInputRef = useRef<HTMLInputElement>(null)
  const responseImageRefs = useRef<{ [key: number]: HTMLInputElement | null }>({})
  const inboxImageRef = useRef<HTMLInputElement>(null)
  const editingKeywordRefs = useRef<{ [key: number]: HTMLInputElement | null }>({})

  // Auto-focus on title input when editing
  useEffect(() => {
    if (isEditingTitle && titleInputRef.current) {
      titleInputRef.current.focus()
      titleInputRef.current.select()
    }
  }, [isEditingTitle])

  // Auto-focus on keyword input when editing
  useEffect(() => {
    console.log('=== KEYWORD EDITING useEffect ===')
    console.log('editingKeywordIndex:', editingKeywordIndex)

    if (editingKeywordIndex !== null) {
      console.log('Trying to focus input at index:', editingKeywordIndex)
      const inputRef = editingKeywordRefs.current[editingKeywordIndex]
      console.log('Input ref found:', inputRef)

      if (inputRef) {
        console.log('Focusing and selecting input')
        setTimeout(() => {
          inputRef.focus()
          inputRef.select()
        }, 10)
      }
      // else: do nothing
    }
  }, [editingKeywordIndex])

  // Auto-update title based on first keyword if title hasn't been manually edited
  useEffect(() => {
    const firstKeyword = editingKeywords.find(k => k.trim() !== '')
    const defaultTitle = `รายการที่ ${index + 1}`

    if (firstKeyword && !hasManuallyEditedTitle && !isEditingTitle) {
      const newTitle = `${firstKeyword} (รายการที่ ${index + 1})`
      setEditingTitle(newTitle)
    } else if (!firstKeyword && !hasManuallyEditedTitle && !isEditingTitle) {
      setEditingTitle(defaultTitle)
    }
  }, [editingKeywords, index, isEditingTitle, hasManuallyEditedTitle])

  // Auto-save when data changes
  useEffect(() => {
    const updatedPair: Pair = {
      ...pair,
      title: editingTitle,
      hasManuallyEditedTitle,
      keywords: editingKeywords.filter(k => k.trim() !== ''),
      responses: editingResponses.filter(r => r.text.trim() !== ''),
      enabled: isEnabled,
      expanded: isExpanded
    }
    onUpdate(index, updatedPair)
  }, [editingKeywords, editingResponses, editingTitle, hasManuallyEditedTitle, isEnabled, isExpanded])

  const handleTitleChange = (value: string) => {
    setEditingTitle(value)
    setHasManuallyEditedTitle(true)
  }

  const handleTitleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      setIsEditingTitle(false)
    } else if (e.key === 'Escape') {
      setEditingTitle(pair.title || `รายการที่ ${index + 1}`)
      setIsEditingTitle(false)
    }
  }

  const addKeyword = (keyword: string) => {
    if (keyword.trim() && !editingKeywords.includes(keyword.trim())) {
      setEditingKeywords([...editingKeywords, keyword.trim()])
    }
  }

  const deleteKeyword = (kidx: number) => {
    setEditingKeywords(prev => prev.filter((_, i) => i !== kidx))
  }

  const handleKeywordChange = (kidx: number, value: string) => {
    setEditingKeywords(prev => prev.map((k, i) => i === kidx ? value : k))
  }

  const startEditingKeyword = (kidx: number) => {
    console.log('=== START EDITING KEYWORD ===')
    console.log('Index:', kidx)
    console.log('Current keyword:', editingKeywords[kidx])
    console.log('Current editingKeywordIndex:', editingKeywordIndex)
    setEditingKeywordIndex(kidx)
    console.log('Set editingKeywordIndex to:', kidx)
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
      // Reset to original value and finish editing
      setEditingKeywords(prev => prev.map((k, i) => i === kidx ? (pair.keywords[kidx] || k) : k))
      setEditingKeywordIndex(null)
    }
  }

  const handleResponseChange = (ridx: number, value: string) => {
    setEditingResponses(prev => prev.map((r, i) =>
      i === ridx ? { ...r, text: value } : r
    ))
  }

  const addResponse = () => {
    const newResponseObj: Response = { text: '' }
    setEditingResponses([...editingResponses, newResponseObj])
  }

  const deleteResponse = (ridx: number) => {
    setEditingResponses(prev => prev.filter((_, i) => i !== ridx))
  }

  const handlePairSettingChange = (key: keyof Pair, value: any) => {
    const updatedPair: Pair = {
      ...pair,
      title: editingTitle,
      keywords: editingKeywords.filter(k => k.trim() !== ''),
      responses: editingResponses.filter(r => r.text.trim() !== ''),
      enabled: isEnabled,
      expanded: isExpanded,
      [key]: value
    }
    onUpdate(index, updatedPair)
  }

  const handleResponseImageUpload = (ridx: number, event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (e) => {
        setEditingResponses(prev => prev.map((r, i) =>
          i === ridx ? { ...r, image: e.target?.result as string } : r
        ))
      }
      reader.readAsDataURL(file)
    }
  }

  const handleInboxImageUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (e) => {
        handlePairSettingChange('inboxImage', e.target?.result as string)
      }
      reader.readAsDataURL(file)
    }
  }

  const toggleEnabled = () => {
    setIsEnabled(!isEnabled)
  }

  const toggleExpanded = () => {
    setIsExpanded(!isExpanded)
  }

  return (
    <div
      className={`bg-white rounded-xl shadow-lg overflow-hidden transition-all duration-300 ${!isEnabled ? 'opacity-60' : ''
        }`}
      data-pair-index={index}
    >
      {/* Rule Header */}
      <div className="p-4 bg-gray-50 flex items-center justify-between">
        <div className="flex items-center space-x-3">
          {/* Toggle Switch */}
          <button
            onClick={toggleEnabled}
            className={`w-12 h-6 rounded-full transition-colors duration-200 ${isEnabled ? 'bg-green-500' : 'bg-gray-300'
              }`}
          >
            <div className={`w-5 h-5 bg-white rounded-full shadow-md transform transition-transform duration-200 ${isEnabled ? 'translate-x-6' : 'translate-x-0.5'
              }`} />
          </button>

          {/* Editable Title */}
          {isEditingTitle ? (
            <input
              ref={titleInputRef}
              type="text"
              value={editingTitle}
              onChange={(e) => handleTitleChange(e.target.value)}
              onKeyDown={handleTitleKeyDown}
              onBlur={() => setIsEditingTitle(false)}
              className="text-lg font-semibold bg-transparent border-b-2 border-blue-300 focus:border-blue-500 focus:outline-none px-1"
              placeholder="ชื่อกฎ"
            />
          ) : (
            <input
              type="text"
              value={editingTitle}
              onChange={(e) => handleTitleChange(e.target.value)}
              className="text-lg font-semibold bg-transparent border-b-2 border-transparent hover:border-gray-300 focus:border-blue-500 focus:outline-none px-1"
              placeholder="ชื่อกฎ"
            />
          )}

          {/* Stats */}
          <span className="text-sm text-gray-500">
            {editingKeywords.length} คีย์เวิร์ด • {editingResponses.length} คำตอบ
          </span>
        </div>

        <div className="flex items-center space-x-2">
          {/* Expand/Collapse Button */}
          <button
            onClick={toggleExpanded}
            className="p-2 hover:bg-gray-200 rounded-lg transition-colors"
          >
            {isExpanded ? <ChevronUpIcon className="w-5 h-5" /> : <ChevronDownIcon className="w-5 h-5" />}
          </button>

          {/* Delete Button */}
          {!showDeleteConfirm ? (
            <button
              onClick={() => setShowDeleteConfirm(true)}
              className="p-2 text-red-500 hover:bg-red-50 rounded-lg transition-colors"
            >
              <TrashIcon className="w-5 h-5" />
            </button>
          ) : (
            <div className="flex items-center space-x-2 bg-red-50 rounded-lg px-3 py-2">
              <span className="text-red-700 text-sm font-medium">ลบ?</span>
              <button
                onClick={() => onDelete(index)}
                className="p-1 text-red-600 hover:text-red-800"
              >
                <CheckIcon className="w-4 h-4" />
              </button>
              <button
                onClick={() => setShowDeleteConfirm(false)}
                className="p-1 text-gray-600 hover:text-gray-800"
              >
                <XMarkIcon className="w-4 h-4" />
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Rule Content */}
      {isExpanded && (
        <div className="p-6 space-y-6">
          {/* Keywords Section */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-3">
              คีย์เวิร์ด (คำที่ต้องการให้ตอบ)
              <span className="text-gray-500 font-normal text-xs ml-2">คลิกคีย์เวิร์ดเพื่อแก้ไข</span>
            </label>
            <div className="flex flex-wrap gap-2 mb-3">
              {editingKeywords.map((keyword, idx) => (
                <div
                  key={idx}
                  className="inline-flex items-center bg-blue-100 text-blue-800 rounded-full text-sm hover:bg-blue-200 transition-colors group"
                >
                  {editingKeywordIndex === idx ? (
                    <input
                      ref={(ref) => { editingKeywordRefs.current[idx] = ref }}
                      type="text"
                      value={keyword}
                      onChange={(e) => handleKeywordChange(idx, e.target.value)}
                      onKeyDown={(e) => handleEditingKeywordKeyDown(idx, e)}
                      onBlur={finishEditingKeyword}
                      className="bg-transparent text-blue-800 text-sm focus:outline-none min-w-[60px] max-w-[200px] w-auto px-2 py-1 border border-blue-300 rounded-full shadow-sm transition-all duration-150"
                      style={{ background: 'rgba(255,255,255,0.95)', position: 'relative', zIndex: 10 }}
                      tabIndex={0}
                      autoFocus
                    />
                  ) : (
                    <>
                      <div
                        onClick={(e) => {
                          e.preventDefault()
                          e.stopPropagation()
                          startEditingKeyword(idx)
                        }}
                        className="px-3 py-1 cursor-pointer hover:underline select-none flex-grow group-hover:bg-blue-200 rounded-full transition-all duration-150"
                        title="คลิกเพื่อแก้ไข"
                        style={{ minWidth: 40 }}
                      >
                        {keyword}
                      </div>
                      <button
                        onClick={(e) => {
                          e.preventDefault()
                          e.stopPropagation()
                          deleteKeyword(idx)
                        }}
                        className="px-2 py-1 text-blue-600 hover:text-red-600 font-bold text-lg leading-none flex-shrink-0"
                        title="ลบคีย์เวิร์ด"
                      >
                        ×
                      </button>
                    </>
                  )}
                </div>
              ))}
            </div>
            <input
              ref={keywordInputRef}
              type="text"
              value={newKeyword}
              onChange={(e) => setNewKeyword(e.target.value)}
              placeholder="พิมพ์คีย์เวิร์ดแล้วกด Enter, Tab หรือ , (เช่น สวัสดี, ราคา, สอบถาม)"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              onKeyDown={(e) => {
                if ((e.key === 'Enter' || e.key === 'Tab' || e.key === ',') && newKeyword.trim()) {
                  e.preventDefault();
                  addKeyword(newKeyword)
                  setNewKeyword('')
                }
              }}
            />
          </div>

          {/* Responses Section */}
          <div>
            <div className="flex items-center justify-between mb-3">
              <label className="block text-sm font-medium text-gray-700">
                คำตอบสำหรับคอมเมนต์
              </label>
              <div className="flex items-center text-sm text-gray-500">
                <ShuffleIcon className="w-4 h-4" />
                <span className="ml-1">ระบบจะสุ่มเลือก 1 คำตอบ</span>
                <InfoIcon />
              </div>
            </div>
            <div className="space-y-3 mb-3">
              {editingResponses.map((response, idx) => (
                <div key={idx} className="border border-gray-200 rounded-lg p-4 space-y-3 bg-gray-50">
                  <div className="flex items-start space-x-2">
                    <div className="flex-1">
                      <div className="text-xs font-medium text-gray-500 mb-1">คำตอบที่ {idx + 1}</div>
                      <textarea
                        value={response.text}
                        onChange={(e) => handleResponseChange(idx, e.target.value)}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-none bg-white"
                        rows={2}
                        placeholder="ข้อความตอบกลับ"
                      />
                    </div>
                    <button
                      onClick={() => deleteResponse(idx)}
                      className="p-2 text-red-500 hover:bg-red-50 rounded-lg"
                    >
                      <TrashIcon className="w-5 h-5" />
                    </button>
                  </div>

                  {/* Image Upload for Response */}
                  <div className="flex items-center space-x-2">
                    <input
                      ref={(ref) => { responseImageRefs.current[idx] = ref }}
                      type="file"
                      accept="image/*"
                      className="hidden"
                      onChange={(e) => handleResponseImageUpload(idx, e)}
                    />
                    <button
                      onClick={() => responseImageRefs.current[idx]?.click()}
                      className="flex items-center px-3 py-1 text-sm border border-gray-300 rounded-lg hover:bg-gray-50 bg-white"
                    >
                      <PhotoIcon className="w-5 h-5" />
                      <span className="ml-2">{response.image ? 'เปลี่ยนรูป' : 'เพิ่มรูป'}</span>
                    </button>

                    {response.image && (
                      <div className="flex items-center space-x-2">
                        <img
                          src={response.image}
                          alt="Response"
                          className="h-10 w-10 object-cover rounded"
                        />
                        <button
                          onClick={() => {
                            setEditingResponses(prev => prev.map((r, i) =>
                              i === idx ? { ...r, image: undefined } : r
                            ))
                          }}
                          className="text-red-500 text-sm hover:underline"
                        >
                          ลบรูป
                        </button>
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
            <button
              onClick={addResponse}
              className="w-full py-2 border-2 border-dashed border-gray-300 rounded-lg text-gray-500 hover:border-blue-500 hover:text-blue-500 transition-colors flex items-center justify-center"
            >
              <PlusIcon className="w-5 h-5 mr-2" />
              เพิ่มคำตอบ (เพื่อให้ระบบสุ่มเลือก)
            </button>
          </div>

          {/* Options for this rule */}
          <div className="border-t pt-4 space-y-3">
            <label className="flex items-center space-x-3">
              <input
                type="checkbox"
                checked={pair.hideCommentsAfterReply || false}
                onChange={(e) => handlePairSettingChange('hideCommentsAfterReply', e.target.checked)}
                className="w-4 h-4 text-blue-600 rounded"
              />
              <span className="text-gray-700 flex items-center">
                <EyeSlashIcon className="w-5 h-5" />
                <span className="ml-2">ซ่อนคอมเมนต์หลังตอบ</span>
              </span>
            </label>

            <label className="flex items-center space-x-3">
              <input
                type="checkbox"
                checked={pair.enableInboxIntegration || false}
                onChange={(e) => handlePairSettingChange('enableInboxIntegration', e.target.checked)}
                className="w-4 h-4 text-blue-600 rounded"
              />
              <span className="text-gray-700 flex items-center">
                <InboxIcon className="w-5 h-5" />
                <span className="ml-2">ดึงข้อความเข้า Inbox</span>
              </span>
            </label>

            {pair.enableInboxIntegration && (
              <div className="ml-7 space-y-3 p-3 bg-blue-50 rounded-lg">
                <textarea
                  value={pair.inboxResponse || ''}
                  onChange={(e) => handlePairSettingChange('inboxResponse', e.target.value)}
                  className="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-none"
                  rows={2}
                  placeholder="ข้อความที่จะส่งใน Inbox (ไม่ใส่ = ใช้ข้อความจากคำตอบ)"
                />

                {/* Inbox Image Upload - Separate from comment response images */}
                <div className="flex items-center space-x-2">
                  <input
                    ref={inboxImageRef}
                    type="file"
                    accept="image/*"
                    className="hidden"
                    onChange={handleInboxImageUpload}
                  />
                  <button
                    onClick={() => inboxImageRef.current?.click()}
                    className="flex items-center px-3 py-1 text-sm border border-gray-300 rounded-lg hover:bg-gray-50 bg-white"
                  >
                    <PhotoIcon className="w-5 h-5" />
                    <span className="ml-2">{pair.inboxImage ? 'เปลี่ยนรูปใน Inbox' : 'เพิ่มรูปใน Inbox'}</span>
                  </button>

                  {pair.inboxImage && (
                    <div className="flex items-center space-x-2">
                      <img
                        src={pair.inboxImage}
                        alt="Inbox"
                        className="h-10 w-10 object-cover rounded"
                      />
                      <button
                        onClick={() => handlePairSettingChange('inboxImage', undefined)}
                        className="text-red-500 text-sm hover:underline"
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
