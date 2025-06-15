import { NextRequest } from 'next/server'

export async function GET(
    request: NextRequest,
    { params }: { params: Promise<{ id: string }> }
) {
    try {
        const { id } = await params
        const response = await fetch(`http://chatbot:8090/api/rules/${id}/responses`, {
            cache: 'no-store'
        })

        if (!response.ok) {
            throw new Error('Failed to fetch rule responses from backend')
        }

        const data = await response.json()

        return Response.json({
            success: true,
            data: data.data || [],
            timestamp: new Date().toISOString()
        })
    } catch (error) {
        console.error('API Error:', error)
        return Response.json(
            {
                error: 'Failed to fetch rule responses',
                message: error instanceof Error ? error.message : 'Unknown error'
            },
            { status: 500 }
        )
    }
}

export async function POST(
    request: NextRequest,
    { params }: { params: Promise<{ id: string }> }
) {
    try {
        const { id } = await params
        const body = await request.json()

        const response = await fetch(`http://chatbot:8090/api/rules/${id}/responses`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        })

        if (!response.ok) {
            throw new Error('Failed to create rule response in backend')
        }

        const data = await response.json()

        return Response.json({
            success: true,
            data: data.data,
            timestamp: new Date().toISOString()
        })
    } catch (error) {
        console.error('API Error:', error)
        return Response.json(
            {
                error: 'Failed to create rule response',
                message: error instanceof Error ? error.message : 'Unknown error'
            },
            { status: 500 }
        )
    }
}
