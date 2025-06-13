'use client'

import { useState, useEffect } from 'react'
import { PlusIcon, BookmarkIcon } from '@heroicons/react/24/outline'
import KeywordPairCard from './components/KeywordPairCard'
import FallbackCard from './components/FallbackCard'
import StatusCard from './components/StatusCard'
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
      responses: [{ text: '' }] // Start with one empty response object
    }
    addPair(newPair)
    
    // Auto-focus on the new pair's keyword input after a short delay
    setTimeout(() => {
      const newPairIndex = data.pairs.length
      const newCard = document.querySelector(`[data-pair-index="${newPairIndex}"]`)
      if (newCard) {
        const addKeywordButton = newCard.querySelector('.add-keyword-button') as HTMLButtonElement
        if (addKeywordButton) {
          addKeywordButton.click()
        }
      }
    }, 100)
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 p-4">
      <div className="max-w-4xl mx-auto">
        {/* Header Section */}
        <div className="bg-white rounded-2xl shadow-xl border-2 border-gray-100 mb-8 overflow-hidden hover:shadow-2xl hover:border-gray-200 transition-all duration-300">
          <div className="bg-gradient-to-r from-indigo-600 to-purple-600 px-8 py-6">
            <div className="flex items-center justify-between">
              <div>
                <h1 className="text-2xl font-bold text-white mb-2">
                  🤖 ระบบตอบกลับอัตโนมัติ Facebook 🟡 Checking
                </h1>
              </div>
              <div className="hidden md:flex items-center space-x-4">
                <div className="bg-white/10 rounded-lg px-4 py-2">
                  <span className="text-white text-sm font-medium">
                    {data.pairs.length} รายการ
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Fallback Response Card */}
        <div className="mb-8">
          <FallbackCard
            settings={data.fallbackSettings || { enabled: false, responses: [{ text: '' }] }}
            onUpdate={updateFallbackSettings}
          />
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
        <div className="bg-white rounded-2xl shadow-xl border-2 border-gray-100 px-8 py-6">
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
