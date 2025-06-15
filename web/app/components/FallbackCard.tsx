'use client'

import { useState, useRef, useEffect } from 'react'
import {
  PlusIcon,
  TrashIcon,
  PhotoIcon,
  ExclamationTriangleIcon,
  EyeSlashIcon,
  InboxIcon
} from '@heroicons/react/24/outline'
import type { Response } from '../types/keyword'
import SimpleResponseManager from './SimpleResponseManager'

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
      : [{ text: '' }]
  )
  const responseImageRefs = useRef<{ [key: number]: HTMLInputElement | null }>({})
  const inboxImageRef = useRef<HTMLInputElement>(null)

  // Auto-save when data changes
  useEffect(() => {
    const updatedSettings: FallbackSettings = {
      ...settings,
      responses: editingResponses.filter(r => r.text.trim() !== '')
    }
    onUpdate(updatedSettings)
  }, [editingResponses, onUpdate, settings])

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
    if (editingResponses.length > 1) {
      setEditingResponses(prev => prev.filter((_, i) => i !== ridx))
    }
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

  const handleEnabledChange = (enabled: boolean) => {
    const updatedSettings: FallbackSettings = {
      enabled,
      responses: editingResponses.filter(r => r.text.trim() !== ''),
      hideCommentsAfterReply: settings.hideCommentsAfterReply,
      enableInboxIntegration: settings.enableInboxIntegration,
      inboxResponse: settings.inboxResponse,
      inboxImage: settings.inboxImage
    }
    onUpdate(updatedSettings)
  }

  const handleSettingChange = (key: keyof FallbackSettings, value: any) => {
    const updatedSettings: FallbackSettings = {
      ...settings,
      responses: editingResponses.filter(r => r.text.trim() !== ''),
      [key]: value
    }
    onUpdate(updatedSettings)
  }

  const handleInboxImageUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (e) => {
        const base64 = e.target?.result as string
        handleSettingChange('inboxImage', base64)
      }
      reader.readAsDataURL(file)
    }
  }

  return (
    <div className="modern-card mb-6 overflow-hidden bg-gradient-to-br from-blue-50 to-sky-50 border-2 border-blue-300 fade-in-up">
      <div className="px-8 py-6">
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center space-x-3 flex-1">
            <div className="w-10 h-10 bg-blue-500 rounded-full flex items-center justify-center text-white text-lg font-bold shadow-lg">
              💬
            </div>
            <div>
              <h3 className="text-lg font-semibold text-blue-800">
                คำตอบกรณีไม่ตรงคีย์เวิร์ด
              </h3>
              <p className="text-sm text-blue-600">สำหรับคอมเมนต์ที่ไม่ตรงกับคีย์เวิร์ดใดๆ</p>
            </div>
          </div>
          <label className="flex items-center space-x-2 bg-white rounded-lg px-3 py-2 border border-blue-300 shadow-sm">
            <input
              type="checkbox"
              checked={settings.enabled}
              onChange={(e) => handleEnabledChange(e.target.checked)}
              className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
            />
            <span className="text-sm font-medium text-blue-700">เปิดใช้งาน</span>
          </label>
        </div>

        {settings.enabled ? (
          <div className="space-y-6">
            {/* Info Banner */}
            <div className="bg-blue-100 border-2 border-blue-300 rounded-xl p-4 shadow-inner">
              <div className="flex items-center space-x-3 text-sm text-blue-800">
                <ExclamationTriangleIcon className="h-5 w-5 flex-shrink-0" />
                <div>
                  <p className="font-medium">การทำงานของคำตอบ Fallback</p>
                  <p className="text-blue-700">คำตอบนี้จะถูกส่งเมื่อไม่พบคีย์เวิร์ดที่ตรงกับคอมเมนต์</p>
                </div>
              </div>
            </div>

            {/* Responses Section */}
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <h4 className="text-base font-semibold text-amber-800 flex items-center">
                  💬 คำตอบอัตโนมัติ
                </h4>
                <span className="text-xs text-amber-600 bg-amber-200 px-2 py-1 rounded-full">
                  {editingResponses.length} คำตอบ
                </span>
              </div>

              <div className="mt-2">
                <SimpleResponseManager
                  ruleId={0}
                  responses={editingResponses}
                  onUpdateResponses={(newResponses) => setEditingResponses(newResponses)}
                />
              </div>
            </div>

            {/* Feature Options - Bottom Section */}
            <div className="bg-amber-100 border-2 border-amber-300 rounded-lg p-4 space-y-4">
              <h4 className="text-sm font-semibold text-amber-800 flex items-center mb-4">
                ⚡ ฟีเจอร์เสริม
              </h4>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* Hide Comments */}
                <label className="flex items-start space-x-3 cursor-pointer bg-white rounded-lg p-4 border-2 border-amber-200 hover:border-amber-400 transition-all shadow-sm hover:shadow-md">
                  <input
                    type="checkbox"
                    className="mt-1 rounded border-gray-300 text-amber-600 focus:ring-amber-500"
                    checked={settings.hideCommentsAfterReply || false}
                    onChange={(e) => handleSettingChange('hideCommentsAfterReply', e.target.checked)}
                  />
                  <div>
                    <span className="text-sm font-medium text-gray-900 flex items-center">
                      <EyeSlashIcon className="h-4 w-4 mr-2 text-amber-600" />
                      ซ่อนคอมเมนต์หลังตอบ
                    </span>
                    <p className="text-xs text-gray-500 mt-1">
                      ซ่อนคอมเมนต์อัตโนมัติหลังจากตอบกลับ
                    </p>
                  </div>
                </label>

                {/* Inbox Integration */}
                <label className="flex items-start space-x-3 cursor-pointer bg-white rounded-lg p-4 border-2 border-amber-200 hover:border-amber-400 transition-all shadow-sm hover:shadow-md">
                  <input
                    type="checkbox"
                    className="mt-1 rounded border-gray-300 text-amber-600 focus:ring-amber-500"
                    checked={settings.enableInboxIntegration || false}
                    onChange={(e) => handleSettingChange('enableInboxIntegration', e.target.checked)}
                  />
                  <div className="flex-1">
                    <span className="text-sm font-medium text-gray-900 flex items-center">
                      <InboxIcon className="h-4 w-4 mr-2 text-amber-600" />
                      ดึงข้อความเข้า Inbox
                    </span>
                    <p className="text-xs text-gray-500 mt-1">
                      ดึงข้อความคอมเมนต์เข้า Facebook Inbox
                    </p>
                  </div>
                </label>
              </div>

              {/* Inbox Settings - Show when enabled */}
              {settings.enableInboxIntegration && (
                <div className="bg-white border-2 border-amber-300 rounded-lg p-4 space-y-3 mt-4 shadow-sm">
                  <h5 className="text-sm font-medium text-amber-800 flex items-center">
                    📧 ตั้งค่า Inbox
                  </h5>

                  <div>
                    <label className="block text-xs font-medium text-gray-600 mb-1">
                      ข้อความใน Inbox:
                    </label>
                    <textarea
                      className="w-full px-3 py-2 border border-amber-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500 text-sm"
                      value={settings.inboxResponse || ''}
                      onChange={(e) => handleSettingChange('inboxResponse', e.target.value)}
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
                        id="inbox-image-fallback"
                      />
                      <label
                        htmlFor="inbox-image-fallback"
                        className="inline-flex items-center px-3 py-1 text-sm bg-white hover:bg-gray-50 border border-amber-300 rounded-md cursor-pointer transition-colors"
                      >
                        <PhotoIcon className="h-4 w-4 mr-1" />
                        <span className="ml-1">{settings.inboxImage ? 'เปลี่ยนรูปใน Inbox' : 'เพิ่มรูปใน Inbox'}</span>
                      </label>
                      {settings.inboxImage && (
                        <div className="flex items-center space-x-2">
                          <img
                            src={settings.inboxImage}
                            alt="Inbox preview"
                            className="h-8 w-8 object-cover rounded border border-amber-400"
                          />
                          <button
                            type="button"
                            onClick={() => handleSettingChange('inboxImage', undefined)}
                            className="text-red-500 hover:text-red-700 text-xs"
                          >
                            ลบรูป
                          </button>
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        ) : (
          <div className="text-center py-8">
            <div className="text-amber-600 text-sm">
              เปิดใช้งานเพื่อตั้งค่าคำตอบกรณีไม่ตรงคีย์เวิร์ด
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
