'use client'

import { useState, useEffect } from 'react'
import { PlusIcon, BookmarkIcon, ChartBarIcon, CheckCircleIcon, ChatBubbleLeftRightIcon } from '@heroicons/react/24/outline'
import KeywordPairCard from './components/KeywordPairCard-redesigned'
import FallbackCard from './components/FallbackCard-redesigned'
import StatusCard from './components/StatusCard'
import KeywordEditTest from './components/KeywordEditTest'
import { useKeywordData } from './hooks/useKeywordData'
import type { Pair, FilterSettings, FallbackSettings } from './types/keyword'

export default function HomePage() {
  const {
    data,
    loading,
    error,
    updatePair,
    addPair,
    deletePair,
    updateSettings,
    updateFilterSettings,
    updateFallbackSettings,
    saveData
  } = useKeywordData()

  const [healthStatus, setHealthStatus] = useState<{
    status: string
    database?: { is_healthy: boolean }
  } | null>(null)

  useEffect(() => {
    // Fetch health status
    const fetchHealth = async () => {
      try {
        const response = await fetch('/health')
        if (response.ok) {
          const data = await response.json()
          setHealthStatus(data)
        }
      } catch (error) {
        console.error('Failed to fetch health status:', error)
      }
    }

    fetchHealth()
    // Poll health status every 30 seconds
    const interval = setInterval(fetchHealth, 30000)
    return () => clearInterval(interval)
  }, [])

  const handleAddNewPair = () => {
    const newPair: Pair = {
      keywords: [],  // Start with empty keywords array
      responses: [{ text: '' }], // Start with one empty response object
      enabled: true, // Default to enabled
      expanded: true // Default to EXPANDED for new pairs
    }
    addPair(newPair)
    
    // Auto-focus on the new pair's keyword input after a short delay
    setTimeout(() => {
      const newPairIndex = data.pairs.length - 1 // Fix: Use length - 1 for correct index
      const newCard = document.querySelector(`[data-pair-index="${newPairIndex}"]`)
      if (newCard) {
        // Try to find and focus the keyword input directly
        const keywordInput = newCard.querySelector('input[placeholder*="คีย์เวิร์ด"]') as HTMLInputElement
        if (keywordInput) {
          keywordInput.focus()
          keywordInput.scrollIntoView({ behavior: 'smooth', block: 'center' })
        }
      }
    }, 300) // Increase delay to ensure DOM is fully updated
  }

  // Calculate quick stats
  const totalRules = data.pairs.length
  const enabledRules = data.pairs.filter(pair => pair.enabled !== false).length
  const totalKeywords = data.pairs.reduce((sum, pair) => sum + pair.keywords.length, 0)

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 p-4">
      <div className="max-w-4xl mx-auto">
        {/* Header Section with Gradient */}
        <div className="bg-white rounded-2xl shadow-xl border-0 mb-8 overflow-hidden">
          <div className="bg-gradient-to-r from-blue-600 to-indigo-600 px-8 py-6">
            <div className="flex items-center justify-between">
              <div>
                <h1 className="text-3xl font-bold text-white mb-2">
                  🤖 ระบบตอบกลับอัตโนมัติ Facebook
                </h1>
                <p className="text-blue-100 text-sm">
                  จัดการคีย์เวิร์ดและคำตอบสำหรับคอมเมนต์ Facebook
                </p>
              </div>
              {healthStatus && (
                <div className="hidden md:block">
                  <div className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${
                    healthStatus.status === 'healthy' 
                      ? 'bg-green-100 text-green-800' 
                      : 'bg-yellow-100 text-yellow-800'
                  }`}>
                    <div className={`w-2 h-2 rounded-full mr-2 ${
                      healthStatus.status === 'healthy' ? 'bg-green-500' : 'bg-yellow-500'
                    }`}></div>
                    {healthStatus.status === 'healthy' ? 'ออนไลน์' : 'กำลังตรวจสอบ'}
                  </div>
                </div>
              )}
            </div>
          </div>
          
          {/* Quick Stats */}
          <div className="px-8 py-6 bg-gray-50 border-t border-gray-100">
            <div className="grid grid-cols-3 gap-6">
              <div className="text-center">
                <div className="flex items-center justify-center w-12 h-12 bg-blue-100 rounded-lg mx-auto mb-2">
                  <ChartBarIcon className="w-6 h-6 text-blue-600" />
                </div>
                <div className="text-2xl font-bold text-gray-900">{totalRules}</div>
                <div className="text-sm text-gray-500">กฎทั้งหมด</div>
              </div>
              <div className="text-center">
                <div className="flex items-center justify-center w-12 h-12 bg-green-100 rounded-lg mx-auto mb-2">
                  <CheckCircleIcon className="w-6 h-6 text-green-600" />
                </div>
                <div className="text-2xl font-bold text-gray-900">{enabledRules}</div>
                <div className="text-sm text-gray-500">กฎที่เปิดใช้งาน</div>
              </div>
              <div className="text-center">
                <div className="flex items-center justify-center w-12 h-12 bg-purple-100 rounded-lg mx-auto mb-2">
                  <ChatBubbleLeftRightIcon className="w-6 h-6 text-purple-600" />
                </div>
                <div className="text-2xl font-bold text-gray-900">{totalKeywords}</div>
                <div className="text-sm text-gray-500">คีย์เวิร์ดทั้งหมด</div>
              </div>
            </div>
          </div>
        </div>

        {/* Fallback Response Card - Moved to top */}
        <div className="mb-8">
          <FallbackCard
            settings={data.fallbackSettings || { enabled: false, responses: [{ text: '' }] }}
            onUpdate={updateFallbackSettings}
          />
        </div>

        {/* DEBUG: Keyword Edit Test */}
        <div className="mb-8">
          <KeywordEditTest />
        </div>

        {/* Keyword-Response Pairs Section */}
        <div className="mb-8">
          <h2 className="text-xl font-bold text-gray-800 mb-6 flex items-center">
            💬 รายการคีย์เวิร์ด-คำตอบ
          </h2>
          
          <div className="space-y-4">
            {data.pairs.map((pair, idx) => (
              <KeywordPairCard
                key={idx}
                pair={pair}
                index={idx}
                onUpdate={updatePair}
                onDelete={deletePair}
              />
            ))}
          </div>
        </div>

        {/* Empty State */}
        {data.pairs.length === 0 && !loading && (
          <div className="text-center py-16 bg-white rounded-2xl shadow-xl border-2 border-gray-100 mb-8">
            <div className="w-16 h-16 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-full flex items-center justify-center mx-auto mb-4 shadow-lg">
              <span className="text-2xl">📝</span>
            </div>
            <h3 className="text-xl font-bold text-gray-900 mb-2">ยังไม่มีรายการคีย์เวิร์ด</h3>
            <p className="text-gray-500 mb-6 max-w-md mx-auto">สร้างรายการแรกเพื่อเริ่มใช้งานระบบตอบกลับอัตโนมัติ Facebook</p>
            <button
              onClick={handleAddNewPair}
              className="inline-flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-all duration-200 bg-gradient-to-r from-indigo-500 to-purple-500 text-white border-transparent hover:from-indigo-600 hover:to-purple-600 hover:shadow-lg hover:scale-105 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 active:scale-95"
            >
              <PlusIcon className="h-5 w-5 mr-2" />
              เพิ่มรายการแรก
            </button>
          </div>
        )}

        {/* Action Buttons */}
        <div className="bg-white rounded-2xl shadow-xl border-0 px-8 py-6 mb-8">
          <div className="flex flex-col lg:flex-row gap-6 items-center justify-between">
            <div className="flex items-center space-x-4">
              <div className="w-12 h-12 bg-gradient-to-br from-green-500 to-emerald-600 rounded-full flex items-center justify-center shadow-lg">
                <span className="text-xl">💡</span>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-900">
                  เคล็ดลับการใช้งาน
                </p>
                <p className="text-sm text-gray-600">
                  การเปลี่ยนแปลงจะถูกบันทึกอัตโนมัติ • สามารถเพิ่มรูปภาพและคำตอบหลายๆ อันได้
                </p>
              </div>
            </div>
            <div className="flex gap-3">
              <button
                onClick={handleAddNewPair}
                className="btn-modern btn-success"
                type="button"
              >
                <PlusIcon className="h-5 w-5 mr-2" />
                เพิ่มรายการใหม่
              </button>
            </div>
          </div>
        </div>

        {/* Advanced Settings */}
        <div className="bg-white rounded-2xl shadow-xl border-0 px-8 py-6">
          <div className="border-b border-gray-200 pb-4 mb-6">
            <h3 className="text-lg font-semibold text-gray-900">ตั้งค่าขั้นสูง</h3>
            <p className="text-sm text-gray-500 mt-1">ตั้งค่าเพิ่มเติมสำหรับการตอบกลับ</p>
          </div>
          
          <div className="space-y-4">
            <label className="flex items-center">
              <input
                type="checkbox"
                checked={data.filterSettings?.skipMentions || false}
                onChange={(e) => updateFilterSettings({ 
                  ...data.filterSettings, 
                  skipMentions: e.target.checked 
                })}
                className="rounded border-gray-300 text-indigo-600 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50"
              />
              <span className="ml-3 text-sm text-gray-700">
                ไม่ตอบคอมเมนต์ที่มีการแท็ก (@)
              </span>
            </label>
            
            <label className="flex items-center">
              <input
                type="checkbox"
                checked={data.filterSettings?.skipStickers || false}
                onChange={(e) => updateFilterSettings({ 
                  ...data.filterSettings, 
                  skipStickers: e.target.checked 
                })}
                className="rounded border-gray-300 text-indigo-600 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50"
              />
              <span className="ml-3 text-sm text-gray-700">
                ไม่ตอบคอมเมนต์ที่มีสติกเกอร์
              </span>
            </label>
          </div>
        </div>

        {/* Footer Note */}
        <div className="mt-8 text-center">
          <p className="text-sm text-gray-500">
            📝 หน้านี้สำหรับตั้งค่าตอบคอมเมนต์ Facebook • สำหรับตอบข้อความแชท (Facebook, LINE, TikTok) จะอยู่ในหน้าตั้งค่าแยกต่างหาก
          </p>
        </div>

        {/* Loading State */}
        {loading && data.pairs.length === 0 && (
          <div className="text-center py-12">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600 mx-auto"></div>
            <p className="text-gray-500 mt-4">กำลังโหลดข้อมูล...</p>
          </div>
        )}

        {/* Error Display */}
        {error && (
          <div className="mb-6 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
            <strong>เกิดข้อผิดพลาด:</strong> {error}
          </div>
        )}
      </div>
    </div>
  )
}
