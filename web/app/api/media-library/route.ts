import { NextRequest, NextResponse } from 'next/server'

// Google Cloud Storage integration (optional)
let storage: any
let bucket: any

try {
  // Only initialize if environment variables are set
  if (process.env.GOOGLE_CLOUD_PROJECT_ID && process.env.GOOGLE_CLOUD_STORAGE_BUCKET) {
    const { Storage } = require('@google-cloud/storage')
    storage = new Storage({
      projectId: process.env.GOOGLE_CLOUD_PROJECT_ID,
      // In production, use service account key file
      // keyFilename: process.env.GOOGLE_CLOUD_KEY_FILE,
    })
    bucket = storage.bucket(process.env.GOOGLE_CLOUD_STORAGE_BUCKET)
  }
} catch (error) {
  console.log('Google Cloud Storage not configured, using mock data fallback')
}

// Enhanced mock media library with real sample images
const mockMediaLibrary = [
  {
    id: '1',
    name: 'sample1.jpg',
    url: 'https://images.unsplash.com/photo-1517077304055-6e89abbf09b0?w=300&h=200&fit=crop',
    type: 'image' as const,
    uploadedAt: '2024-01-01T00:00:00Z'
  },
  {
    id: '2', 
    name: 'sample2.jpg',
    url: 'https://images.unsplash.com/photo-1501594907352-04cda38ebc29?w=300&h=200&fit=crop',
    type: 'image' as const,
    uploadedAt: '2024-01-02T00:00:00Z'
  },
  {
    id: '3',
    name: 'sample3.jpg', 
    url: 'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=300&h=200&fit=crop',
    type: 'image' as const,
    uploadedAt: '2024-01-03T00:00:00Z'
  },
  {
    id: '4',
    name: 'sample4.jpg', 
    url: 'https://images.unsplash.com/photo-1519904981063-b0cf448d479e?w=300&h=200&fit=crop',
    type: 'image' as const,
    uploadedAt: '2024-01-04T00:00:00Z'
  },
  {
    id: '5',
    name: 'sample5.jpg', 
    url: 'https://images.unsplash.com/photo-1516641051054-9df6a1aad654?w=300&h=200&fit=crop',
    type: 'image' as const,
    uploadedAt: '2024-01-05T00:00:00Z'
  }
]

export async function GET() {
  try {
    if (bucket) {
      // Use Google Cloud Storage
      const [files] = await bucket.getFiles({
        prefix: 'uploads/',
      })

      const mediaFiles = files.map((file: any) => ({
        id: file.name,
        name: file.name.split('/').pop(),
        url: `https://storage.googleapis.com/${bucket.name}/${file.name}`,
        type: file.metadata.contentType?.startsWith('video/') ? 'video' : 'image',
        uploadedAt: file.metadata.timeCreated
      }))

      return NextResponse.json({
        success: true,
        files: mediaFiles,
        source: 'google-cloud-storage'
      })
    } else {
      // Fallback to mock data for development
      return NextResponse.json({
        success: true,
        files: mockMediaLibrary,
        source: 'mock-data'
      })
    }
  } catch (error) {
    console.error('Failed to fetch media library:', error)
    
    // Fallback to mock data on error
    return NextResponse.json({
      success: true,
      files: mockMediaLibrary,
      source: 'fallback-mock-data'
    })
  }
}
