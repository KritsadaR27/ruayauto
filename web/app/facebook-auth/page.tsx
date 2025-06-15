'use client'

import FacebookAuth from '../components/FacebookAuth'
import Navigation from '../components/Navigation'

export default function FacebookAuthPage() {
    const handlePagesUpdate = (pages: any[]) => {
        console.log('Facebook pages updated:', pages)
        // Handle pages update here
    }

    return (
        <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 p-4">
            <div className="max-w-4xl mx-auto">
                {/* Navigation */}
                <Navigation />

                {/* Header */}
                <div className="bg-white rounded-2xl shadow-xl mb-8 overflow-hidden">
                    <div className="bg-gradient-to-r from-blue-600 to-indigo-600 p-8">
                        <div className="flex items-center justify-between">
                            <div>
                                <h1 className="text-3xl font-bold text-white mb-2">
                                    Facebook Integration Center 🔐
                                </h1>
                                <p className="text-blue-100">
                                    จัดการการเชื่อมต่อ Facebook OAuth และ Pages Management
                                </p>
                            </div>
                            <a
                                href="/"
                                className="bg-white/20 hover:bg-white/30 text-white px-4 py-2 rounded-lg font-medium transition-colors"
                            >
                                ← กลับหน้าหลัก
                            </a>
                        </div>
                    </div>
                </div>

                {/* Facebook Auth Component */}
                <FacebookAuth onPagesUpdate={handlePagesUpdate} />

                {/* Features Info */}
                <div className="mt-8 bg-white rounded-xl shadow-lg p-6">
                    <h2 className="text-xl font-bold text-gray-800 mb-4">ฟีเจอร์ที่รองรับ</h2>
                    <div className="grid md:grid-cols-2 gap-6">
                        <div className="space-y-3">
                            <h3 className="font-semibold text-gray-700">OAuth Authentication</h3>
                            <ul className="text-sm text-gray-600 space-y-1">
                                <li>• เชื่อมต่อกับ Facebook Account</li>
                                <li>• จัดการ Access Tokens</li>
                                <li>• Auto Token Refresh</li>
                                <li>• Secure Session Management</li>
                            </ul>
                        </div>
                        <div className="space-y-3">
                            <h3 className="font-semibold text-gray-700">Page Management</h3>
                            <ul className="text-sm text-gray-600 space-y-1">
                                <li>• เชื่อมต่อ Facebook Pages</li>
                                <li>• ตรวจสอบสถานะการเชื่อมต่อ</li>
                                <li>• จัดการ Page Permissions</li>
                                <li>• Monitor Page Activities</li>
                            </ul>
                        </div>
                    </div>
                </div>

                {/* Status Info */}
                <div className="mt-6 p-4 bg-green-50 rounded-lg border border-green-200">
                    <p className="text-sm text-green-800">
                        ✅ <strong>Production Mode:</strong> ระบบใช้ Facebook API จริง ไม่มี Mock Data แล้ว
                        <br />
                        📋 <strong>Requirements:</strong> ต้องตั้งค่า FACEBOOK_APP_ID และ FACEBOOK_APP_SECRET ใน environment variables
                        <br />
                        🔗 <strong>Backend Integration:</strong> เชื่อมต่อกับ Go backend API สำหรับจัดการ OAuth และ Pages
                    </p>
                </div>
            </div>
        </div>
    )
}
