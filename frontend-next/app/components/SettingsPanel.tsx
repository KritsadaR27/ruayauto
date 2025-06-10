'use client'

import { useState, useRef } from 'react'
import { 
  PlusIcon, 
  TrashIcon, 
  CogIcon,
  PhotoIcon,
  EyeSlashIcon,
  InboxIcon
} from '@heroicons/react/24/outline'
import type { Settings } from '../types/keyword'

interface Props {
  settings: Settings
  onUpdate: (settings: Settings) => void
}

export default function SettingsPanel({ settings, onUpdate }: Props) {
  const [localSettings, setLocalSettings] = useState<Settings>(settings)
  const [newDefaultResponse, setNewDefaultResponse] = useState('')
  const inboxImageRef = useRef<HTMLInputElement>(null)

  const handleSettingChange = (key: keyof Settings, value: any) => {
    const updated = { ...localSettings, [key]: value }
    setLocalSettings(updated)
    onUpdate(updated)
  }

  const handleDefaultResponseChange = (idx: number, value: string) => {
    const updatedResponses = localSettings.defaultResponses.map((r, i) => (i === idx ? value : r))
    const updated = { ...localSettings, defaultResponses: updatedResponses }
    setLocalSettings(updated)
    onUpdate(updated)
  }

  const addDefaultResponse = () => {
    if (newDefaultResponse.trim() && !localSettings.defaultResponses.includes(newDefaultResponse.trim())) {
      const updated = {
        ...localSettings,
        defaultResponses: [...localSettings.defaultResponses, newDefaultResponse.trim()]
      }
      setLocalSettings(updated)
      onUpdate(updated)
      setNewDefaultResponse('')
    }
  }

  const removeDefaultResponse = (index: number) => {
    const updated = {
      ...localSettings,
      defaultResponses: localSettings.defaultResponses.filter((_, i) => i !== index)
    }
    setLocalSettings(updated)
    onUpdate(updated)
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
    <div className="card mb-6">
      <div className="px-8 py-6 border-b border-gray-100">
        <h2 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
          ⚙️ การตั้งค่าระบบ
        </h2>
        
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Default Response Settings */}
          <div className="space-y-4">
            <div className="flex items-start space-x-3">
              <input 
                type="checkbox" 
                id="enableDefault"
                className="checkbox-input mt-1"
                checked={localSettings.enableDefault} 
                onChange={e => handleSettingChange('enableDefault', e.target.checked)} 
              />
              <div className="flex-1">
                <label htmlFor="enableDefault" className="text-sm font-medium text-gray-900 cursor-pointer">
                  คำตอบกรณีไม่ตรงคีย์เวิร์ด
                </label>
                <p className="text-xs text-gray-500 mt-1">
                  ระบบจะสุ่มเลือกคำตอบจากรายการด้านล่าง
                </p>
              </div>
            </div>
            
            {localSettings.enableDefault && (
              <div className="ml-7 space-y-3">
                {localSettings.defaultResponses.map((res, idx) => (
                  <div key={idx} className="flex items-start space-x-2">
                    <textarea
                      className="input-field flex-1"
                      value={res}
                      onChange={e => handleDefaultResponseChange(idx, e.target.value)}
                      placeholder="ข้อความตอบกลับเมื่อไม่ตรงคีย์เวิร์ด"
                      rows={2}
                    />
                    {localSettings.defaultResponses.length > 1 && (
                      <button 
                        className="mt-1 p-1.5 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors" 
                        onClick={() => removeDefaultResponse(idx)} 
                        type="button"
                      >
                        <TrashIcon className="h-4 w-4" />
                      </button>
                    )}
                  </div>
                ))}
                
                <div className="flex space-x-2">
                  <input
                    type="text"
                    value={newDefaultResponse}
                    onChange={(e) => setNewDefaultResponse(e.target.value)}
                    onKeyPress={(e) => e.key === 'Enter' && addDefaultResponse()}
                    placeholder="เพิ่มข้อความตอบกลับ..."
                    className="input-field flex-1"
                  />
                  <button 
                    className="btn-small-indigo" 
                    onClick={addDefaultResponse} 
                    type="button"
                  >
                    <PlusIcon className="h-3 w-3 mr-1" />
                    เพิ่มข้อความ
                  </button>
                </div>
              </div>
            )}
          </div>

          {/* Comment Filter Settings */}
          <div className="space-y-4">
            <h3 className="text-sm font-medium text-gray-900">
              🛡️ ตัวกรองคอมเมนต์
            </h3>
            <div className="space-y-3">
              <label className="flex items-center space-x-3 cursor-pointer group">
                <input 
                  type="checkbox" 
                  className="checkbox-input"
                  checked={localSettings.noTag} 
                  onChange={e => handleSettingChange('noTag', e.target.checked)} 
                />
                <div>
                  <span className="text-sm font-medium text-gray-900 group-hover:text-indigo-600 transition-colors">
                    ไม่ตอบคอมเมนต์ที่มีการแท็ก
                  </span>
                  <p className="text-xs text-gray-500">
                    ข้ามคอมเมนต์ที่มี @ mention
                  </p>
                </div>
              </label>
              <label className="flex items-center space-x-3 cursor-pointer group">
                <input 
                  type="checkbox" 
                  className="checkbox-input"
                  checked={localSettings.noSticker} 
                  onChange={e => handleSettingChange('noSticker', e.target.checked)} 
                />
                <div>
                  <span className="text-sm font-medium text-gray-900 group-hover:text-indigo-600 transition-colors">
                    ไม่ตอบคอมเมนต์ที่มีสติกเกอร์
                  </span>
                  <p className="text-xs text-gray-500">
                    ข้ามคอมเมนต์ที่มี emoji/sticker
                  </p>
                </div>
              </label>
            </div>
          </div>

          {/* Advanced Features */}
          <div className="space-y-4">
            <h3 className="text-sm font-medium text-gray-900">
              🚀 ฟีเจอร์เพิ่มเติม
            </h3>
            <div className="space-y-3">
              <label className="flex items-center space-x-3 cursor-pointer group">
                <input 
                  type="checkbox" 
                  className="checkbox-input"
                  checked={localSettings.hideCommentsAfterReply || false} 
                  onChange={e => handleSettingChange('hideCommentsAfterReply', e.target.checked)} 
                />
                <div>
                  <span className="text-sm font-medium text-gray-900 group-hover:text-indigo-600 transition-colors">
                    <EyeSlashIcon className="h-4 w-4 inline mr-1" />
                    ซ่อนคอมเมนต์หลังตอบ
                  </span>
                  <p className="text-xs text-gray-500">
                    ซ่อนคอมเมนต์อัตโนมัติหลังจากระบบตอบกลับแล้ว
                  </p>
                </div>
              </label>

              <label className="flex items-start space-x-3 cursor-pointer group">
                <input 
                  type="checkbox" 
                  className="checkbox-input mt-1"
                  checked={localSettings.enableInboxIntegration || false} 
                  onChange={e => handleSettingChange('enableInboxIntegration', e.target.checked)} 
                />
                <div className="flex-1">
                  <span className="text-sm font-medium text-gray-900 group-hover:text-indigo-600 transition-colors">
                    <InboxIcon className="h-4 w-4 inline mr-1" />
                    ดึงข้อความเข้า Inbox
                  </span>
                  <p className="text-xs text-gray-500 mb-3">
                    ดึงข้อความคอมเมนต์เข้า Facebook Inbox เพื่อการติดตามที่ดีขึ้น
                  </p>
                  
                  {localSettings.enableInboxIntegration && (
                    <div className="space-y-3 ml-0">
                      <div>
                        <label className="block text-xs font-medium text-gray-700 mb-1">
                          ข้อความใน Inbox หลังดึงเข้า:
                        </label>
                        <textarea
                          className="input-field w-full"
                          value={localSettings.inboxResponse || ''}
                          onChange={e => handleSettingChange('inboxResponse', e.target.value)}
                          placeholder="สวัสดีครับ เราได้ตอบในคอมเมนต์แล้ว หากมีคำถามเพิ่มเติมสามารถสอบถามที่นี่ได้เลยครับ"
                          rows={3}
                        />
                      </div>
                      
                      <div>
                        <label className="block text-xs font-medium text-gray-700 mb-2">
                          <PhotoIcon className="h-4 w-4 inline mr-1" />
                          รูปภาพแนบใน Inbox (ไม่บังคับ):
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
                          {localSettings.inboxImage && (
                            <div className="flex items-center space-x-2">
                              <img
                                src={localSettings.inboxImage}
                                alt="Inbox preview"
                                className="h-8 w-8 object-cover rounded border"
                              />
                              <button
                                type="button"
                                onClick={() => handleSettingChange('inboxImage', undefined)}
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
      </div>
    </div>
  )
}
