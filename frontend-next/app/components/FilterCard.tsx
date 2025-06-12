'use client'

import { useState } from 'react'
import { ShieldCheckIcon, AtSymbolIcon, FaceSmileIcon } from '@heroicons/react/24/outline'

interface FilterSettings {
  skipMentions: boolean
  skipStickers: boolean
}

interface Props {
  settings: FilterSettings
  onUpdate: (settings: FilterSettings) => void
}

export default function FilterCard({ settings, onUpdate }: Props) {
  const handleSettingChange = (key: keyof FilterSettings, value: boolean) => {
    const updatedSettings = {
      ...settings,
      [key]: value
    }
    onUpdate(updatedSettings)
  }

  return (
    <div className="modern-card mb-6 overflow-hidden bg-gradient-to-br from-orange-50 to-red-50 border-2 border-orange-300 hover:border-orange-400 transition-all duration-300 hover:-translate-y-1 fade-in-up">
      <div className="px-8 py-6">
        <div className="flex items-center mb-6">
          <h3 className="text-lg font-semibold text-gray-900 flex items-center">
            🛡️ ตัวกรองคอมเมนต์
          </h3>
        </div>

        <div className="space-y-4">
          <p className="text-sm text-gray-600 mb-4">
            ตั้งค่าการกรองคอมเมนต์ที่ไม่ต้องการตอบกลับ
          </p>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {/* Skip Mentions */}
            <label className="flex items-start space-x-3 cursor-pointer bg-orange-50 rounded-lg p-4 border border-orange-200 hover:border-orange-300 transition-colors">
              <input 
                type="checkbox" 
                className="mt-1 rounded border-gray-300 text-orange-600 focus:ring-orange-500"
                checked={settings.skipMentions || false} 
                onChange={(e) => handleSettingChange('skipMentions', e.target.checked)} 
              />
              <div>
                <span className="text-sm font-medium text-gray-900">
                  <AtSymbolIcon className="h-4 w-4 inline mr-1 text-orange-600" />
                  ไม่ตอบคอมเมนต์ที่มีการแท็ก
                </span>
                <p className="text-xs text-gray-500">
                  ข้ามคอมเมนต์ที่มี @ mention
                </p>
              </div>
            </label>

            {/* Skip Stickers */}
            <label className="flex items-start space-x-3 cursor-pointer bg-yellow-50 rounded-lg p-4 border border-yellow-200 hover:border-yellow-300 transition-colors">
              <input 
                type="checkbox" 
                className="mt-1 rounded border-gray-300 text-yellow-600 focus:ring-yellow-500"
                checked={settings.skipStickers || false} 
                onChange={(e) => handleSettingChange('skipStickers', e.target.checked)} 
              />
              <div>
                <span className="text-sm font-medium text-gray-900">
                  <FaceSmileIcon className="h-4 w-4 inline mr-1 text-yellow-600" />
                  ไม่ตอบคอมเมนต์ที่มีสติกเกอร์
                </span>
                <p className="text-xs text-gray-500">
                  ข้ามคอมเมนต์ที่มี emoji/sticker
                </p>
              </div>
            </label>
          </div>

          <div className="bg-gray-50 border border-gray-200 rounded-lg p-4 mt-4">
            <div className="flex items-center space-x-2 text-sm text-gray-600">
              <ShieldCheckIcon className="h-4 w-4 text-gray-500" />
              <span>การตั้งค่าเหล่านี้จะใช้กับทุกรายการคีย์เวิร์ด</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
