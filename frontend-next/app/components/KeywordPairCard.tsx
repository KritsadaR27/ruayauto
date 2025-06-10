'use client'

import { useState, useRef, useEffect } from 'react'
import { 
  PlusIcon, 
  TrashIcon,
  PhotoIcon,
  EyeSlashIcon,
  InboxIcon,
  CogIcon,
  ChevronUpIcon
} from '@heroicons/react/24/outline'
import type { Pair, Response } from '../types/keyword'

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
  const [isAddingKeyword, setIsAddingKeyword] = useState(false)
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false)
  const [showSettings, setShowSettings] = useState(false)
  const [isEditingTitle, setIsEditingTitle] = useState(false)
  const [editingTitle, setEditingTitle] = useState(pair.title || `รายการที่ ${index + 1}`)
  const [hasManuallyEditedTitle, setHasManuallyEditedTitle] = useState(pair.hasManuallyEditedTitle || false)
  const keywordInputRef = useRef<HTMLInputElement>(null)
  const titleInputRef = useRef<HTMLInputElement>(null)
  const responseImageRefs = useRef<{ [key: number]: HTMLInputElement | null }>({})
  const inboxImageRef = useRef<HTMLInputElement>(null)
  const keywordRefs = useRef<{ [key: string]: HTMLInputElement | null }>({})

  // Auto-focus on title input when editing
  useEffect(() => {
    if (isEditingTitle && titleInputRef.current) {
      titleInputRef.current.focus()
      titleInputRef.current.select()
    }
  }, [isEditingTitle])

  // Auto-focus on keyword input when adding
  useEffect(() => {
    if (isAddingKeyword && keywordInputRef.current) {
      keywordInputRef.current.focus()
    }
  }, [isAddingKeyword])

  // Auto-update title based on first keyword if title hasn't been manually edited
  useEffect(() => {
    const firstKeyword = editingKeywords.find(k => k.trim() !== '')
    const defaultTitle = `รายการที่ ${index + 1}`
    
    // Only auto-update title if:
    // 1. There's a first keyword
    // 2. User hasn't manually edited the title
    // 3. User is not currently editing the title
    if (firstKeyword && !hasManuallyEditedTitle && !isEditingTitle) {
      const newTitle = `${firstKeyword} (รายการที่ ${index + 1})`
      setEditingTitle(newTitle)
    } else if (!firstKeyword && !hasManuallyEditedTitle && !isEditingTitle) {
      // Reset to default if no keywords and no manual edit
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
      responses: editingResponses.filter(r => r.text.trim() !== '')
    }
    onUpdate(index, updatedPair)
  }, [editingKeywords, editingResponses, editingTitle, hasManuallyEditedTitle])

  const handleKeywordChange = (kidx: number, value: string) => {
    setEditingKeywords(prev => prev.map((k, i) => i === kidx ? value : k))
  }

  const addKeyword = () => {
    if (newKeyword.trim() && !editingKeywords.includes(newKeyword.trim())) {
      setEditingKeywords([...editingKeywords, newKeyword.trim()])
      setNewKeyword('')
      if (keywordInputRef.current) {
        keywordInputRef.current.focus()
      }
    }
  }

  const handleKeywordInputKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' || e.key === 'Tab') {
      e.preventDefault()
      if (newKeyword.trim()) {
        addKeyword()
      }
    } else if (e.key === 'Escape') {
      setIsAddingKeyword(false)
      setNewKeyword('')
    }
  }

  const handleKeywordInputBlur = (e: React.FocusEvent) => {
    const relatedTarget = e.relatedTarget as HTMLElement
    if (relatedTarget && (relatedTarget.closest('.keyword-add-button') || relatedTarget === keywordInputRef.current)) {
      return
    }
    
    if (newKeyword.trim()) {
      addKeyword()
    }
    setIsAddingKeyword(false)
  }

  const startAddingKeyword = () => {
    setIsAddingKeyword(true)
  }

  const finishAddingKeywords = () => {
    if (newKeyword.trim()) {
      addKeyword()
    }
    setIsAddingKeyword(false)
  }

  const deleteKeyword = (kidx: number) => {
    setEditingKeywords(prev => prev.filter((_, i) => i !== kidx))
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

  const handleDeleteConfirm = () => {
    onDelete(index)
  }

  const handleDeleteCancel = () => {
    setShowDeleteConfirm(false)
  }

  const handlePairSettingChange = (key: keyof Pair, value: any) => {
    const updatedPair: Pair = {
      ...pair,
      title: editingTitle,
      keywords: editingKeywords.filter(k => k.trim() !== ''),
      responses: editingResponses.filter(r => r.text.trim() !== ''),
      [key]: value
    }
    onUpdate(index, updatedPair)
  }

  const handleTitleChange = (value: string) => {
    setEditingTitle(value)
    setHasManuallyEditedTitle(true) // Mark as manually edited
  }

  const handleTitleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      setIsEditingTitle(false)
    } else if (e.key === 'Escape') {
      // Reset to previous state without marking as manually edited
      setEditingTitle(pair.title || `รายการที่ ${index + 1}`)
      setIsEditingTitle(false)
    }
  }

  const handleTitleBlur = () => {
    setIsEditingTitle(false)
  }

  const handleResponseImageUpload = (ridx: number, event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (e) => {
        const base64 = e.target?.result as string
        setEditingResponses(prev => prev.map((r, i) => 
          i === ridx ? { ...r, image: base64 } : r
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
        const base64 = e.target?.result as string
        handlePairSettingChange('inboxImage', base64)
      }
      reader.readAsDataURL(file)
    }
  }

  return (
    <div className="bg-white rounded-2xl shadow-xl border-2 border-gray-100 mb-6 overflow-hidden hover:shadow-2xl hover:border-gray-200 transition-all duration-300 hover:-translate-y-1" data-pair-index={index}>
      <div className="px-8 py-6">
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center space-x-2 flex-1">
            <span className="text-lg">📝</span>
            {isEditingTitle ? (
              <input
                ref={titleInputRef}
                type="text"
                value={editingTitle}
                onChange={(e) => handleTitleChange(e.target.value)}
                onKeyDown={handleTitleKeyDown}
                onBlur={handleTitleBlur}
                className="text-lg font-semibold text-gray-900 bg-transparent border-b-2 border-indigo-300 focus:border-indigo-500 focus:outline-none px-1 py-1 min-w-[200px]"
                placeholder="ชื่อรายการ"
              />
            ) : (
              <h3 
                className="text-lg font-semibold text-gray-900 cursor-pointer hover:text-indigo-600 transition-colors px-1 py-1"
                onClick={() => setIsEditingTitle(true)}
                title="คลิกเพื่อแก้ไขชื่อ"
              >
                {editingTitle}
              </h3>
            )}
          </div>
          <div className="flex items-center space-x-2">
            <button
              onClick={() => setShowSettings(!showSettings)}
              className="p-2 text-gray-600 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors"
              title="ตั้งค่ารายการนี้"
            >
              <CogIcon className="h-4 w-4" />
            </button>
            {!showDeleteConfirm ? (
              <button 
                onClick={() => setShowDeleteConfirm(true)}
                className="p-2 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors"
                title="ลบรายการนี้"
              >
                <TrashIcon className="h-4 w-4" />
              </button>
            ) : (
              <div className="flex items-center space-x-2 bg-red-50 border border-red-200 rounded-lg px-3 py-2">
                <span className="text-red-700 text-sm font-medium">ลบรายการนี้?</span>
                <button 
                  onClick={handleDeleteConfirm}
                  className="px-2 py-1 bg-red-600 text-white text-xs rounded hover:bg-red-700 transition-colors"
                >
                  ลบ
                </button>
                <button 
                  onClick={handleDeleteCancel}
                  className="px-2 py-1 bg-gray-300 text-gray-700 text-xs rounded hover:bg-gray-400 transition-colors"
                >
                  ยกเลิก
                </button>
              </div>
            )}
          </div>
        </div>

        <div className="space-y-6">
          {/* Advanced Settings Panel - Hidden by default */}
          {showSettings && (
            <div className="bg-gray-50 border border-gray-200 rounded-lg p-6 space-y-6">
              <div className="flex items-center justify-between">
                <h4 className="text-sm font-semibold text-gray-900 flex items-center">
                  ⚙️ ตั้งค่าขั้นสูง
                </h4>
                <button
                  onClick={() => setShowSettings(false)}
                  className="text-gray-400 hover:text-gray-600"
                >
                  <ChevronUpIcon className="h-4 w-4" />
                </button>
              </div>

              <div className="text-sm text-gray-500 text-center py-4">
                🔧 ตั้งค่าขั้นสูงจะถูกเพิ่มในอนาคต
              </div>
            </div>
          )}

          {/* Keywords Section */}
          <div className="space-y-4">
            <div className="flex items-center">
              <h4 className="text-sm font-semibold text-gray-900 flex items-center">
                🔍 คีย์เวิร์ด
              </h4>
            </div>
            
            <div className="flex flex-wrap gap-2 items-center">
              {editingKeywords.map((keyword, kidx) => (
                <div key={kidx} className="bg-indigo-100 text-indigo-800 px-3 py-1 rounded-full text-sm font-medium border-2 border-indigo-200 hover:border-indigo-300 transition-colors flex items-center">
                  <input
                    ref={(el) => {
                      keywordRefs.current[`${index}-${kidx}`] = el
                    }}
                    className="bg-transparent border-none outline-none text-sm font-medium text-indigo-800 placeholder-indigo-400 min-w-[60px] max-w-[120px]"
                    value={keyword}
                    onChange={(e) => handleKeywordChange(kidx, e.target.value)}
                    placeholder="คีย์เวิร์ด"
                  />
                  {editingKeywords.length > 1 && (
                    <button 
                      className="ml-2 text-indigo-400 hover:text-red-500 transition-colors" 
                      onClick={() => deleteKeyword(kidx)}
                    >
                      ×
                    </button>
                  )}
                </div>
              ))}

              {!isAddingKeyword && editingKeywords.length === 0 && (
                <div className="text-gray-400 text-sm italic">
                  ยังไม่มีคีย์เวิร์ด คลิก "เพิ่มคีย์เวิร์ด" เพื่อเริ่มต้น
                </div>
              )}

              {isAddingKeyword && (
                <div className="flex items-center gap-2">
                  <div className="bg-indigo-100 text-indigo-800 px-3 py-1 rounded-full text-sm font-medium border-2 border-indigo-200 hover:border-indigo-300 transition-colors">
                    <input
                      ref={keywordInputRef}
                      type="text"
                      value={newKeyword}
                      onChange={(e) => setNewKeyword(e.target.value)}
                      onKeyDown={handleKeywordInputKeyDown}
                      onBlur={handleKeywordInputBlur}
                      placeholder="พิมพ์คีย์เวิร์ด..."
                      className="bg-transparent border-none outline-none text-sm font-medium text-indigo-800 placeholder-indigo-400 min-w-[80px] max-w-[120px]"
                    />
                  </div>
                  <button 
                    onClick={finishAddingKeywords}
                    className="keyword-add-button btn-small-green text-xs px-2 py-1"
                    title="เสร็จแล้ว (Enter)"
                  >
                    ✓ เสร็จ
                  </button>
                </div>
              )}
            </div>

            {isAddingKeyword && (
              <div className="text-xs text-gray-500 mt-2">
                💡 พิมพ์คีย์เวิร์ดแล้วกด Enter หรือ Tab เพื่อเพิ่ม • กด Escape เพื่อยกเลิก • คลิกที่อื่นเพื่อเสร็จสิ้น
              </div>
            )}

            {!isAddingKeyword && (
              <div className="mt-3 pt-3 border-t border-gray-100">
                <button 
                  onClick={startAddingKeyword}
                  className="btn-small-indigo add-keyword-button"
                >
                  <PlusIcon className="h-3 w-3 mr-1" />
                  เพิ่มคีย์เวิร์ด
                </button>
              </div>
            )}
          </div>

          {/* Responses Section */}
          <div className="space-y-4">
            <div className="flex items-center">
              <h4 className="text-sm font-semibold text-gray-900 flex items-center">
                💬 คำตอบ
              </h4>
            </div>
            
            <div className="space-y-4">
              {editingResponses.map((response, ridx) => (
                <div key={ridx} className="bg-gray-50 border border-gray-200 rounded-lg p-4 space-y-3">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-xs font-medium text-purple-600 bg-purple-100 px-2 py-1 rounded-full">
                      คำตอบที่ {ridx + 1}
                    </span>
                    {editingResponses.length > 1 && (
                      <button 
                        className="p-1 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors" 
                        onClick={() => deleteResponse(ridx)}
                      >
                        <TrashIcon className="h-3.5 w-3.5" />
                      </button>
                    )}
                  </div>

                  <div className="flex items-start space-x-2">
                    <textarea
                      className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 focus:shadow-lg resize-none text-sm transition-all duration-200"
                      value={response.text}
                      onChange={(e) => handleResponseChange(ridx, e.target.value)}
                      placeholder="ข้อความตอบกลับ"
                      rows={3}
                    />
                  </div>

                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      <input
                        ref={(el) => {
                          responseImageRefs.current[ridx] = el
                        }}
                        type="file"
                        accept="image/*"
                        onChange={(e) => handleResponseImageUpload(ridx, e)}
                        className="hidden"
                      />
                      <button
                        type="button"
                        onClick={() => responseImageRefs.current[ridx]?.click()}
                        className="btn-small-indigo"
                      >
                        <PhotoIcon className="h-3 w-3 mr-1" />
                        {response.image ? 'เปลี่ยนรูป' : 'เพิ่มรูป'}
                      </button>
                    </div>

                    {response.image && (
                      <div className="flex items-center space-x-2">
                        <img
                          src={response.image}
                          alt={`Response ${ridx + 1} image`}
                          className="h-12 w-12 object-cover rounded border border-gray-200"
                        />
                        <button
                          type="button"
                          onClick={() => {
                            setEditingResponses(prev => prev.map((r, i) => 
                              i === ridx ? { ...r, image: undefined } : r
                            ))
                          }}
                          className="text-red-500 hover:text-red-700 text-xs"
                        >
                          ลบรูป
                        </button>
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>

            <div className="mt-3 pt-3 border-t border-gray-100">
              <button 
                onClick={addResponse}
                className="btn-small-purple"
              >
                <PlusIcon className="h-3 w-3 mr-1" />
                เพิ่มคำตอบ
              </button>
            </div>
          </div>

          {/* Feature Options - Bottom Section */}
          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 space-y-4">
            <h4 className="text-sm font-semibold text-blue-900 flex items-center">
              ⚡ ฟีเจอร์เสริม
            </h4>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {/* Hide Comments */}
              <label className="flex items-start space-x-3 cursor-pointer bg-white rounded-lg p-3 border border-blue-100 hover:border-blue-300 transition-colors">
                <input 
                  type="checkbox" 
                  className="mt-1 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  checked={pair.hideCommentsAfterReply || false} 
                  onChange={(e) => handlePairSettingChange('hideCommentsAfterReply', e.target.checked)} 
                />
                <div>
                  <span className="text-sm font-medium text-gray-900">
                    <EyeSlashIcon className="h-4 w-4 inline mr-1 text-blue-600" />
                    ซ่อนคอมเมนต์หลังตอบ
                  </span>
                  <p className="text-xs text-gray-500">
                    ซ่อนคอมเมนต์อัตโนมัติหลังจากตอบกลับ
                  </p>
                </div>
              </label>

              {/* Inbox Integration */}
              <label className="flex items-start space-x-3 cursor-pointer bg-white rounded-lg p-3 border border-blue-100 hover:border-blue-300 transition-colors">
                <input 
                  type="checkbox" 
                  className="mt-1 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  checked={pair.enableInboxIntegration || false} 
                  onChange={(e) => handlePairSettingChange('enableInboxIntegration', e.target.checked)} 
                />
                <div className="flex-1">
                  <span className="text-sm font-medium text-gray-900">
                    <InboxIcon className="h-4 w-4 inline mr-1 text-blue-600" />
                    ดึงข้อความเข้า Inbox
                  </span>
                  <p className="text-xs text-gray-500">
                    ดึงข้อความคอมเมนต์เข้า Facebook Inbox
                  </p>
                </div>
              </label>
            </div>

            {/* Inbox Settings - Show when enabled */}
            {pair.enableInboxIntegration && (
              <div className="bg-white border border-blue-200 rounded-lg p-4 space-y-3">
                <h5 className="text-sm font-medium text-blue-900">📧 ตั้งค่า Inbox</h5>
                
                <div>
                  <label className="block text-xs font-medium text-gray-600 mb-1">
                    ข้อความใน Inbox:
                  </label>
                  <textarea
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-sm"
                    value={pair.inboxResponse || ''}
                    onChange={(e) => handlePairSettingChange('inboxResponse', e.target.value)}
                    placeholder="ข้อความที่จะส่งใน Inbox หลังจากตอบในคอมเมนต์"
                    rows={2}
                  />
                </div>
                
                <div>
                  <label className="block text-xs font-medium text-gray-600 mb-2">
                    <PhotoIcon className="h-4 w-4 inline mr-1" />
                    รูปภาพแนบใน Inbox:
                  </label>
                  <div className="flex items-center space-x-2">
                    <input
                      ref={inboxImageRef}
                      type="file"
                      accept="image/*"
                      onChange={handleInboxImageUpload}
                      className="hidden"
                    />
                    <button
                      type="button"
                      onClick={() => inboxImageRef.current?.click()}
                      className="inline-flex items-center px-2 py-1 text-xs font-medium rounded-md transition-all duration-200 bg-gradient-to-r from-blue-500 to-blue-600 text-white border-transparent hover:from-blue-600 hover:to-blue-700 hover:shadow-md hover:scale-105 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-1 active:scale-95"
                    >
                      <PhotoIcon className="h-3 w-3 mr-1" />
                      เลือกรูป
                    </button>
                    {pair.inboxImage && (
                      <div className="flex items-center space-x-2">
                        <img
                          src={pair.inboxImage}
                          alt="Inbox preview"
                          className="h-8 w-8 object-cover rounded border"
                        />
                        <button
                          type="button"
                          onClick={() => handlePairSettingChange('inboxImage', undefined)}
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
      </div>
    </div>
  )
}
