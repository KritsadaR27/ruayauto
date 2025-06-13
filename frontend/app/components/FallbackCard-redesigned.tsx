'use client'

import { useState, useRef, useEffect } from 'react'
import { 
  PlusIcon, 
  TrashIcon,
  PhotoIcon,
  ExclamationTriangleIcon
} from '@heroicons/react/24/outline'
import type { Response } from '../types/keyword'

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

interface FallbackSettings {
  enabled: boolean
  responses: Response[]
  hideCommentsAfterReply?: boolean
  enableInboxIntegration?: boolean
  inboxResponse?: string
  inboxImage?: string
}

interface Props {
  settings: FallbackSettings
  onUpdate: (settings: FallbackSettings) => void
}

export default function FallbackCard({ settings, onUpdate }: Props) {
  const [editingResponses, setEditingResponses] = useState<Response[]>(
    settings.responses && settings.responses.length > 0 
      ? settings.responses 
      : [{ text: 'ขอบคุณสำหรับข้อความครับ จะรีบตอบกลับโดยเร็วที่สุด' }]
  )
  const responseImageRefs = useRef<{ [key: number]: HTMLInputElement | null }>({})

  // Auto-save when data changes
  useEffect(() => {
    const updatedSettings: FallbackSettings = {
      ...settings,
      responses: editingResponses.filter(r => r.text.trim() !== '')
    }
    onUpdate(updatedSettings)
  }, [editingResponses])

  const handleEnabledChange = (enabled: boolean) => {
    onUpdate({
      ...settings,
      enabled
    })
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

  return (
    <div className="bg-white rounded-xl shadow-lg p-6">
      <div className="flex items-center justify-between mb-4">
        <div>
          <h3 className="text-lg font-semibold text-gray-900">คำตอบเมื่อไม่ตรงกับคีย์เวิร์ดใดๆ</h3>
          <p className="text-sm text-gray-500">ใช้เมื่อคอมเมนต์ไม่ตรงกับกฎที่ตั้งไว้</p>
        </div>
        <button
          onClick={() => handleEnabledChange(!settings.enabled)}
          className={`w-12 h-6 rounded-full transition-colors duration-200 ${
            settings.enabled ? 'bg-green-500' : 'bg-gray-300'
          }`}
        >
          <div className={`w-5 h-5 bg-white rounded-full shadow-md transform transition-transform duration-200 ${
            settings.enabled ? 'translate-x-6' : 'translate-x-0.5'
          }`} />
        </button>
      </div>
      
      {settings.enabled && (
        <div className="space-y-4">
          {/* Info Banner */}
          <div className="bg-blue-50 border border-blue-200 rounded-lg p-3">
            <div className="flex items-center space-x-2 text-sm text-blue-800">
              <ExclamationTriangleIcon className="h-4 w-4 flex-shrink-0" />
              <span>คำตอบนี้จะถูกส่งเมื่อไม่พบคีย์เวิร์ดที่ตรงกับคอมเมนต์</span>
            </div>
          </div>

          {/* Responses Section */}
          <div>
            <div className="flex items-center justify-between mb-3">
              <label className="block text-sm font-medium text-gray-700">
                คำตอบสำรอง
              </label>
              <div className="flex items-center text-sm text-gray-500">
                <ShuffleIcon />
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
                        placeholder="ข้อความตอบกลับสำรอง"
                      />
                    </div>
                    {editingResponses.length > 1 && (
                      <button
                        onClick={() => deleteResponse(idx)}
                        className="p-2 text-red-500 hover:bg-red-50 rounded-lg"
                      >
                        <TrashIcon className="w-5 h-5" />
                      </button>
                    )}
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
              เพิ่มคำตอบสำรอง (เพื่อให้ระบบสุ่มเลือก)
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
