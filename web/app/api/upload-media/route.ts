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
  console.log('Google Cloud Storage not configured, using mock upload')
}

export async function POST(request: NextRequest) {
  try {
    const formData = await request.formData()
    const file = formData.get('file') as File
    const ruleId = formData.get('ruleId') as string

    if (!file) {
      return NextResponse.json(
        { success: false, error: 'No file provided' },
        { status: 400 }
      )
    }

    // Validate file type
    const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp', 'video/mp4', 'video/webm']
    if (!allowedTypes.includes(file.type)) {
      return NextResponse.json(
        { success: false, error: 'Invalid file type. Allowed: JPEG, PNG, GIF, WebP, MP4, WebM' },
        { status: 400 }
      )
    }

    // Validate file size (10MB limit)
    const maxSize = 10 * 1024 * 1024 // 10MB
    if (file.size > maxSize) {
      return NextResponse.json(
        { success: false, error: 'File too large. Maximum size is 10MB' },
        { status: 400 }
      )
    }

    let finalUrl: string

    if (bucket) {
      // Upload to Google Cloud Storage
      const fileName = `uploads/${ruleId}/${Date.now()}-${file.name}`
      const gcsFile = bucket.file(fileName)
      
      const buffer = Buffer.from(await file.arrayBuffer())
      
      await gcsFile.save(buffer, {
        metadata: {
          contentType: file.type,
          metadata: {
            ruleId: ruleId,
            originalName: file.name,
            uploadedAt: new Date().toISOString()
          }
        }
      })

      // Make file publicly accessible
      await gcsFile.makePublic()
      
      finalUrl = `https://storage.googleapis.com/${bucket.name}/${fileName}`
      
      return NextResponse.json({
        success: true,
        url: finalUrl,
        filename: file.name,
        size: file.size,
        type: file.type,
        source: 'google-cloud-storage'
      })
    } else {
      // Mock upload for development
      finalUrl = `https://images.unsplash.com/photo-${Date.now()}?w=300&h=200&fit=crop&mock=true`
      
      return NextResponse.json({
        success: true,
        url: finalUrl,
        filename: file.name,
        size: file.size,
        type: file.type,
        source: 'mock-upload'
      })
    }

  } catch (error) {
    console.error('Upload failed:', error)
    return NextResponse.json(
      { success: false, error: 'Upload failed: ' + (error as Error).message },
      { status: 500 }
    )
  }
}
