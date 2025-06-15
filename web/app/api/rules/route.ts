import { NextRequest } from 'next/server'

export async function GET(request: NextRequest) {
    try {
        const { searchParams } = request.nextUrl
        const limit = searchParams.get('limit') || '10'
        const response = await fetch('http://chatbot:8090/api/rules', {
            cache: 'no-store'
        })
        if (!response.ok) throw new Error('Failed to fetch from backend API')
        const data = await response.json()
        return Response.json({
            success: true,
            data: data,
            timestamp: new Date().toISOString(),
            cached: false
        })
    } catch (error) {
        console.error('API Error:', error)
        return Response.json(
            { error: 'Failed to fetch rules', message: error instanceof Error ? error.message : 'Unknown error' },
            { status: 500 }
        )
    }
}

export async function POST(request: NextRequest) {
    try {
        const body = await request.json()
        const response = await fetch('http://chatbot:8090/api/rules', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        })
        if (!response.ok) throw new Error('Failed to save to backend API')
        const data = await response.json()
        return Response.json({
            success: true,
            data: data,
            timestamp: new Date().toISOString()
        })
    } catch (error) {
        console.error('API Error:', error)
        return Response.json(
            { error: 'Failed to save rules', message: error instanceof Error ? error.message : 'Unknown error' },
            { status: 500 }
        )
    }
}

export async function DELETE(request: NextRequest) {
    try {
        const body = await request.json()
        const response = await fetch('http://chatbot:8090/api/rules', {
            method: 'DELETE',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        })
        if (!response.ok) throw new Error('Failed to delete from backend API')
        const data = await response.json()
        return Response.json({
            success: true,
            data: data,
            timestamp: new Date().toISOString()
        })
    } catch (error) {
        console.error('API Error:', error)
        return Response.json(
            { error: 'Failed to delete rules', message: error instanceof Error ? error.message : 'Unknown error' },
            { status: 500 }
        )
    }
}
