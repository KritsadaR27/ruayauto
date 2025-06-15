'use client'

import React, { useState, useEffect, useRef } from 'react'
import FacebookAuth from './FacebookAuth'
import { Response } from '../types/rule'

// Icons
const PlusIcon = () => (
    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
    </svg>
)

const TrashIcon = () => (
    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
    </svg>
)

const PhotoIcon = () => (
    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
    </svg>
)

const FolderIcon = () => (
    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
    </svg>
)

const WeightIcon = () => (
    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 6l3 1m0 0l-3 9a5.002 5.002 0 006.001 0M6 7l3 9M6 7l6-2m6 2l3-1m-3 1l-3 9a5.002 5.002 0 006.001 0M18 7l3 9m-3-9l-6-2m0-2v2m0 16V5m0 16H9m3 0h3" />
    </svg>
)

const CloudIcon = () => (
    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10" />
    </svg>
)

interface MediaFile {
    id: string
    name: string
    url: string
    type: 'image' | 'video'
    uploadedAt: string
}

interface SimpleResponseManagerProps {
    ruleId: number
    responses: Response[]
    onUpdateResponses: (responses: Response[]) => void
}

export default function SimpleResponseManager({ ruleId, responses, onUpdateResponses }: SimpleResponseManagerProps) {
    const [localResponses, setLocalResponses] = useState<Response[]>(responses)
    const [mediaLibrary, setMediaLibrary] = useState<MediaFile[]>([])
    const [showMediaLibrary, setShowMediaLibrary] = useState<number | null>(null)
    const [uploading, setUploading] = useState<number | null>(null)
    const fileInputRefs = useRef<{ [key: number]: HTMLInputElement | null }>({})

    // Sync with parent component
    useEffect(() => {
        setLocalResponses(responses)
    }, [responses])

    // Load media library
    useEffect(() => {
        loadMediaLibrary()
    }, [])

    const loadMediaLibrary = async () => {
        try {
            const response = await fetch('/api/media-library')
            if (response.ok) {
                const data = await response.json()
                setMediaLibrary(data.files || [])
            }
        } catch (error) {
            console.error('Failed to load media library:', error)
        }
    }

    const addResponse = () => {
        const newResponse: Response = {
            text: '',
            weight: 10,
            is_active: true // เปิดใช้งานตั้งแต่แรก
        }
        const updatedResponses = [...localResponses, newResponse]
        setLocalResponses(updatedResponses)
        onUpdateResponses(updatedResponses)
    }

    const updateResponse = (index: number, field: keyof Response, value: any) => {
        const updatedResponses = localResponses.map((response, i) =>
            i === index ? { ...response, [field]: value } : response
        )
        setLocalResponses(updatedResponses)
        onUpdateResponses(updatedResponses)
    }

    const deleteResponse = (index: number) => {
        const updatedResponses = localResponses.filter((_, i) => i !== index)
        setLocalResponses(updatedResponses)
        onUpdateResponses(updatedResponses)
    }

    const handleFileUpload = async (index: number, file: File) => {
        if (!file) return

        setUploading(index)
        const formData = new FormData()
        formData.append('file', file)
        formData.append('ruleId', ruleId.toString())

        try {
            const response = await fetch('/api/upload-media', {
                method: 'POST',
                body: formData
            })

            if (response.ok) {
                const data = await response.json()
                updateResponse(index, 'image', data.url)
                // Reload media library to include new file
                loadMediaLibrary()
            } else {
                alert('การอัปโหลดไฟล์ล้มเหลว')
            }
        } catch (error) {
            console.error('Upload failed:', error)
            alert('การอัปโหลดไฟล์ล้มเหลว')
        } finally {
            setUploading(index)
        }
    }

    const selectFromLibrary = (index: number, mediaFile: MediaFile) => {
        updateResponse(index, 'image', mediaFile.url)
        setShowMediaLibrary(null)
    }

    const getTotalWeight = () => {
        return localResponses.reduce((sum, response) =>
            (response.is_active ?? true) ? sum + (response.weight || 10) : sum, 0
        )
    }

    const getWeightPercentage = (weight: number | undefined) => {
        const totalWeight = getTotalWeight()
        return totalWeight > 0 ? (((weight || 0) / totalWeight) * 100).toFixed(1) : '0'
    }

    return (
        <div className="space-y-4">
            <div className="flex items-center justify-between">
                <h3 className="text-lg font-semibold text-gray-800">คำตอบ ({localResponses.length})</h3>
                <button
                    onClick={addResponse}
                    className="flex items-center px-3 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
                >
                    <PlusIcon />
                    <span className="ml-2">เพิ่มคำตอบ</span>
                </button>
            </div>

            {localResponses.length === 0 && (
                <div className="text-center py-8 text-gray-500">
                    <p>ยังไม่มีคำตอบ</p>
                    <p className="text-sm">คลิก "เพิ่มคำตอบ" เพื่อเริ่มต้น</p>
                </div>
            )}

            {localResponses.map((response, index) => (
                <div
                    key={index}
                    className={`border rounded-lg p-4 space-y-4 ${(response.is_active ?? true) ? 'border-gray-200 bg-white' : 'border-gray-300 bg-gray-50'
                        }`}
                >
                    {/* Response Header */}
                    <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-3">
                            <label className="flex items-center cursor-pointer">
                                <input
                                    type="checkbox"
                                    checked={response.is_active ?? true}
                                    onChange={(e) => updateResponse(index, 'is_active', e.target.checked)}
                                    className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                />
                                <span className="ml-2 font-medium text-gray-700">
                                    คำตอบที่ {index + 1}
                                </span>
                            </label>

                            {(response.is_active ?? true) && (
                                <div className="flex items-center space-x-2 text-sm text-gray-600">
                                    <WeightIcon />
                                    <span>น้ำหนัก: {response.weight || 10}</span>
                                    <span className="text-blue-600">({getWeightPercentage(response.weight)}%)</span>
                                </div>
                            )}
                        </div>

                        <button
                            onClick={() => deleteResponse(index)}
                            className="text-red-500 hover:text-red-700 p-1"
                            title="ลบคำตอบ"
                        >
                            <TrashIcon />
                        </button>
                    </div>

                    {/* Text Response */}
                    <div className="space-y-2">
                        <label className="block text-sm font-medium text-gray-700">ข้อความตอบ</label>
                        <textarea
                            value={response.text}
                            onChange={(e) => updateResponse(index, 'text', e.target.value)}
                            placeholder="พิมพ์ข้อความตอบที่นี่..."
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-none"
                            rows={3}
                            disabled={!(response.is_active ?? true)}
                        />
                    </div>

                    {/* Image Upload Section */}
                    <div className="space-y-3">
                        <label className="block text-sm font-medium text-gray-700">รูปภาพแนบ (ไม่จำเป็น)</label>

                        {/* Current Image Preview */}
                        {response.image && (
                            <div className="relative inline-block">
                                <img
                                    src={response.image}
                                    alt="Response preview"
                                    className="h-20 w-20 object-cover rounded-lg border"
                                />
                                <button
                                    onClick={() => updateResponse(index, 'image', undefined)}
                                    className="absolute -top-2 -right-2 bg-red-500 text-white rounded-full w-6 h-6 flex items-center justify-center text-xs hover:bg-red-600"
                                >
                                    ✕
                                </button>
                            </div>
                        )}

                        {/* Upload Options */}
                        <div className="flex flex-wrap gap-2">
                            {/* Upload New File */}
                            <input
                                ref={(el) => { fileInputRefs.current[index] = el }}
                                type="file"
                                accept="image/*,video/*"
                                onChange={(e) => {
                                    const file = e.target.files?.[0]
                                    if (file) handleFileUpload(index, file)
                                }}
                                className="hidden"
                                disabled={!(response.is_active ?? true)}
                            />

                            <button
                                onClick={() => fileInputRefs.current[index]?.click()}
                                disabled={!(response.is_active ?? true) || uploading === index}
                                className="inline-flex items-center px-3 py-1.5 text-sm bg-blue-50 text-blue-700 border border-blue-200 rounded-md hover:bg-blue-100 disabled:bg-gray-100 disabled:text-gray-400 disabled:cursor-not-allowed transition-colors"
                            >
                                {uploading === index ? (
                                    <>
                                        <CloudIcon />
                                        <span className="ml-1.5">กำลังอัปโหลด...</span>
                                    </>
                                ) : (
                                    <>
                                        <PhotoIcon />
                                        <span className="ml-1.5">อัปโหลด</span>
                                    </>
                                )}
                            </button>

                            {/* Select from Library */}
                            <button
                                onClick={() => setShowMediaLibrary(showMediaLibrary === index ? null : index)}
                                disabled={!(response.is_active ?? true)}
                                className="inline-flex items-center px-3 py-1.5 text-sm bg-gray-50 text-gray-700 border border-gray-200 rounded-md hover:bg-gray-100 disabled:bg-gray-100 disabled:text-gray-400 disabled:cursor-not-allowed transition-colors"
                            >
                                <FolderIcon />
                                <span className="ml-1.5">ไฟล์เก่า</span>
                            </button>
                        </div>

                        {/* Media Library Modal */}
                        {showMediaLibrary === index && (
                            <div className="border border-gray-300 rounded-lg p-4 bg-gray-50 max-h-60 overflow-y-auto">
                                <div className="flex items-center justify-between mb-3">
                                    <h4 className="font-medium text-gray-700">ไฟล์ที่เคยอัปโหลด</h4>
                                    <button
                                        onClick={() => setShowMediaLibrary(null)}
                                        className="text-gray-500 hover:text-gray-700"
                                    >
                                        ✕
                                    </button>
                                </div>

                                {mediaLibrary.length === 0 ? (
                                    <p className="text-center text-gray-500 py-4">ยังไม่มีไฟล์</p>
                                ) : (
                                    <div className="grid grid-cols-4 gap-2">
                                        {mediaLibrary.map((file) => (
                                            <div
                                                key={file.id}
                                                onClick={() => selectFromLibrary(index, file)}
                                                className="cursor-pointer group relative"
                                            >
                                                <img
                                                    src={file.url}
                                                    alt={file.name}
                                                    className="w-full h-16 object-cover rounded border group-hover:opacity-80 transition-opacity"
                                                />
                                                <div className="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-20 rounded transition-all" />
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </div>
                        )}
                    </div>

                    {/* Weight Control */}
                    {(response.is_active ?? true) && (
                        <div className="space-y-2">
                            <label className="block text-sm font-medium text-gray-700">
                                น้ำหนักในการสุ่มตอบ (ยิ่งสูงยิ่งโอกาสมาก)
                            </label>
                            <div className="flex items-center space-x-3">
                                <input
                                    type="range"
                                    min="1"
                                    max="100"
                                    value={response.weight || 10}
                                    onChange={(e) => updateResponse(index, 'weight', parseInt(e.target.value))}
                                    className="flex-1"
                                />
                                <input
                                    type="number"
                                    min="1"
                                    max="100"
                                    value={response.weight || 10}
                                    onChange={(e) => updateResponse(index, 'weight', parseInt(e.target.value) || 1)}
                                    className="w-16 px-2 py-1 border border-gray-300 rounded text-sm text-center"
                                />
                            </div>
                        </div>
                    )}
                </div>
            ))}

            {/* Weight Summary */}
            {localResponses.some(r => (r.is_active ?? true)) && (
                <div className="bg-blue-50 border border-blue-200 rounded-lg p-3">
                    <div className="text-sm text-blue-800">
                        <strong>สรุปน้ำหนัก:</strong> รวม {getTotalWeight()} คะแนน จาก {localResponses.filter(r => (r.is_active ?? true)).length} คำตอบที่เปิดใช้งาน
                    </div>
                </div>
            )}
        </div>
    )
}
