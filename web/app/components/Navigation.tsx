'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'

const Navigation = () => {
    const pathname = usePathname()

    const navItems = [
        { href: '/', label: '🏠 หน้าหลัก', description: 'จัดการคำตอบอัตโนมัติ' },
        { href: '/facebook-auth', label: '🔐 Facebook Auth', description: 'จัดการการเชื่อมต่อ Facebook' },
    ]

    return (
        <nav className="bg-white shadow-lg rounded-xl mb-6">
            <div className="px-6 py-4">
                <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                        <div className="w-8 h-8 bg-gradient-to-r from-blue-500 to-indigo-500 rounded-lg flex items-center justify-center">
                            <span className="text-white font-bold text-sm">R</span>
                        </div>
                        <span className="text-xl font-bold text-gray-800">ruayChatBot</span>
                    </div>
                    <div className="flex space-x-1">
                        {navItems.map((item) => (
                            <Link key={item.href} href={item.href}>
                                <div className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${pathname === item.href
                                        ? 'bg-blue-100 text-blue-700 border border-blue-200'
                                        : 'text-gray-600 hover:text-blue-600 hover:bg-blue-50'
                                    }`}>
                                    <div className="flex items-center space-x-2">
                                        <span>{item.label}</span>
                                    </div>
                                    <div className="text-xs text-gray-500 mt-1">
                                        {item.description}
                                    </div>
                                </div>
                            </Link>
                        ))}
                    </div>
                </div>
            </div>
        </nav>
    )
}

export default Navigation
