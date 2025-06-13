import React, { useState, useRef, useEffect } from 'react'

// Simple test component for keyword editing
export default function KeywordEditTest() {
  const [keywords, setKeywords] = useState(['สวัสดี', 'ทดสอบ', 'แก้ไข'])
  const [editingIndex, setEditingIndex] = useState<number | null>(null)
  const inputRefs = useRef<{ [key: number]: HTMLInputElement | null }>({})

  const startEdit = (index: number) => {
    console.log('Start editing index:', index)
    setEditingIndex(index)
  }

  const finishEdit = () => {
    console.log('Finish editing')
    setEditingIndex(null)
  }

  const handleChange = (index: number, value: string) => {
    setKeywords(prev => prev.map((k, i) => i === index ? value : k))
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      finishEdit()
    }
  }

  // Auto-focus when editing starts
  useEffect(() => {
    if (editingIndex !== null) {
      const input = inputRefs.current[editingIndex]
      if (input) {
        setTimeout(() => {
          input.focus()
          input.select()
        }, 10)
      }
    }
  }, [editingIndex])

  return (
    <div className="p-8 bg-white">
      <h2 className="text-xl font-bold mb-4">Keyword Edit Test</h2>
      <div className="flex flex-wrap gap-2">
        {keywords.map((keyword, index) => (
          <div
            key={index}
            className="bg-blue-100 text-blue-800 rounded-full px-3 py-1 flex items-center"
          >
            {editingIndex === index ? (
              <input
                ref={(ref) => { inputRefs.current[index] = ref }}
                type="text"
                value={keyword}
                onChange={(e) => handleChange(index, e.target.value)}
                onKeyDown={handleKeyDown}
                onBlur={finishEdit}
                className="bg-transparent text-blue-800 text-sm focus:outline-none min-w-[60px]"
              />
            ) : (
              <span
                onClick={() => startEdit(index)}
                className="cursor-pointer hover:underline"
              >
                {keyword}
              </span>
            )}
            <button
              onClick={() => setKeywords(prev => prev.filter((_, i) => i !== index))}
              className="ml-2 text-red-500 hover:text-red-700"
            >
              ×
            </button>
          </div>
        ))}
      </div>
    </div>
  )
}
