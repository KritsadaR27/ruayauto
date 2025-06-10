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
  const [isEditing, setIsEditing] = useState(true) // Always in editing mode
  const [editingKeywords, setEditingKeywords] = useState<string[]>(pair.keywords)
  const [editingResponses, setEditingResponses] = useState<Response[]>(
    Array.isArray(pair.responses) && typeof pair.responses[0] === 'string' 
      ? (pair.responses as string[]).map(text => ({ text })) 
      : (pair.responses as Response[]) || [{ text: '' }]
  )
  const [newKeyword, setNewKeyword] = useState('')
  const [isAddingKeyword, setIsAddingKeyword] = useState(false)
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false)
  const [showSettings, setShowSettings] = useState(false)
  const keywordInputRef = useRef<HTMLInputElement>(null)
  const responseImageRefs = useRef<{ [key: number]: HTMLInputElement | null }>({})
  const inboxImageRef = useRef<HTMLInputElement>(null)
  const keywordRefs = useRef<{ [key: string]: HTMLInputElement | null }>({})

  // Auto-focus on keyword input when adding
  useEffect(() => {
    if (isAddingKeyword && keywordInputRef.current) {
      keywordInputRef.current.focus()
    }
  }, [isAddingKeyword])

  // Auto-save when data changes
  useEffect(() => {
    const updatedPair: Pair = {
      ...pair,
      keywords: editingKeywords.filter(k => k.trim() !== ''),
      responses: editingResponses.filter(r => r.text.trim() !== '')
    }
    onUpdate(index, updatedPair)
  }, [editingKeywords, editingResponses])

  const handleKeywordChange = (kidx: number, value: string) => {
    setEditingKeywords(prev => prev.map((k, i) => i === kidx ? value : k))
  }

  const addKeyword = () => {
    if (newKeyword.trim() && !editingKeywords.includes(newKeyword.trim())) {
      setEditingKeywords([...editingKeywords, newKeyword.trim()])
      setNewKeyword('')
      // Keep the input active for continuous adding
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
    // Only close if not clicking on the add button or the input itself
    const relatedTarget = e.relatedTarget as HTMLElement
    if (relatedTarget && (relatedTarget.closest('.keyword-add-button') || relatedTarget === keywordInputRef.current)) {
      return
    }
    
    // Add current keyword if not empty before closing
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
      keywords: editingKeywords.filter(k => k.trim() !== ''),
      responses: editingResponses.filter(r => r.text.trim() !== ''),
      [key]: value
    }
    onUpdate(index, updatedPair)
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
    <div className="card" data-pair-index={index}>
      <div className="px-8 py-6">
        <div className="flex items-center justify-between mb-6">
          <h3 className="text-lg font-semibold text-gray-900">
            📝 รายการที่ {index + 1}
          </h3>
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
          {/* Settings Panel */}
          {showSettings && (
            <div className="bg-gray-50 border border-gray-200 rounded-lg p-6 space-y-6">
              <div className="flex items-center justify-between">
                <h4 className="text-sm font-semibold text-gray-900 flex items-center">
                  ⚙️ ตั้งค่ารายการนี้
                </h4>
                <button
                  onClick={() => setShowSettings(false)}
                  className="text-gray-400 hover:text-gray-600"
                >
                  <ChevronUpIcon className="h-4 w-4" />
                </button>
              </div>

              <div className="grid grid-cols-1 gap-6">
                {/* Advanced Features */}
                <div className="space-y-4">
                  <label className="flex items-start space-x-3 cursor-pointer">
                    <input 
                      type="checkbox" 
                      className="mt-1 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
                      checked={pair.hideCommentsAfterReply || false} 
                      onChange={(e) => handlePairSettingChange('hideCommentsAfterReply', e.target.checked)} 
                    />
                    <div>
                      <span className="text-sm font-medium text-gray-900">
                        <EyeSlashIcon className="h-4 w-4 inline mr-1" />
                        ซ่อนคอมเมนต์หลังตอบ
                      </span>
                      <p className="text-xs text-gray-500">
                        ซ่อนคอมเมนต์อัตโนมัติหลังจากตอบกลับรายการนี้
                      </p>
                    </div>
                  </label>

                  <label className="flex items-start space-x-3 cursor-pointer">
                    <input 
                      type="checkbox" 
                      className="mt-1 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
                      checked={pair.enableInboxIntegration || false} 
                      onChange={(e) => handlePairSettingChange('enableInboxIntegration', e.target.checked)} 
                    />
                    <div className="flex-1">
                      <span className="text-sm font-medium text-gray-900">
                        <InboxIcon className="h-4 w-4 inline mr-1" />
                        ดึงข้อความเข้า Inbox
                      </span>
                      <p className="text-xs text-gray-500 mb-3">
                        ดึงข้อความคอมเมนต์เข้า Facebook Inbox สำหรับรายการนี้
                      </p>
                      
                      {pair.enableInboxIntegration && (
                        <div className="space-y-3 mt-3">
                          <div>
                            <label className="block text-xs font-medium text-gray-600 mb-1">
                              ข้อความใน Inbox:
                            </label>
                            <textarea
                              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-sm"
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
                                className="btn-small-indigo"
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
                  </label>
                </div>
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
              {/* Existing Keywords */}
              {editingKeywords.map((keyword, kidx) => (
                <div key={kidx} className="badge-keyword flex items-center">
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

              {/* Show placeholder if no keywords exist */}
              {!isAddingKeyword && editingKeywords.length === 0 && (
                <div className="text-gray-400 text-sm italic">
                  ยังไม่มีคีย์เวิร์ด คลิก "เพิ่มคีย์เวิร์ด" เพื่อเริ่มต้น
                </div>
              )}

              {/* Inline Keyword Input */}
              {isAddingKeyword && (
                <div className="flex items-center gap-2">
                  <div className="badge-keyword">
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

            {/* Add Keyword Button - Moved to bottom */}
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
                  <div className="flex items-start space-x-2">
                    <textarea
                      className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 focus:shadow-lg resize-none text-sm transition-all duration-200"
                      value={response.text}
                      onChange={(e) => handleResponseChange(ridx, e.target.value)}
                      placeholder="ข้อความตอบกลับ"
                      rows={3}
                    />
                    {editingResponses.length > 1 && (
                      <button 
                        className="mt-1 p-1.5 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors" 
                        onClick={() => deleteResponse(ridx)}
                      >
                        <TrashIcon className="h-3.5 w-3.5" />
                      </button>
                    )}
                  </div>

                  {/* Image section for each response */}
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

            {/* Add Response Button - Moved to bottom */}
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
        </div>
      </div>
    </div>
  )
}
