import { NextRequest } from 'next/server'

export async function DELETE(
    request: NextRequest,
    { params }: { params: Promise<{ id: string }> }
) {
    try {
        const { id } = await params
        const response = await fetch(`http://chatbot:8090/api/rules/${id}`, {
            method: 'DELETE',
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
            { error: 'Failed to delete rule', message: error instanceof Error ? error.message : 'Unknown error' },
            { status: 500 }
        )
    }
}

export async function PUT(
    request: NextRequest,
    { params }: { params: Promise<{ id: string }> }
) {
    try {
        const { id } = await params
        const body = await request.json()
        const response = await fetch(`http://chatbot:8090/api/rules/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        })
        if (!response.ok) throw new Error('Failed to update backend API')
        const data = await response.json()
        return Response.json({
            success: true,
            data: data,
            timestamp: new Date().toISOString()
        })
    } catch (error) {
        console.error('API Error:', error)
        return Response.json(
            { error: 'Failed to update rule', message: error instanceof Error ? error.message : 'Unknown error' },
            { status: 500 }
        )
    }
}
