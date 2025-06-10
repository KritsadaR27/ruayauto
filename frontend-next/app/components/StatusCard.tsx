interface StatusCardProps {
  title: string
  status: 'healthy' | 'checking' | 'error' | 'ready'
  description: string
}

export default function StatusCard({ title, status, description }: StatusCardProps) {
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'healthy':
        return 'bg-green-100 text-green-800 border-green-200'
      case 'ready':
        return 'bg-blue-100 text-blue-800 border-blue-200'
      case 'checking':
        return 'bg-yellow-100 text-yellow-800 border-yellow-200'
      case 'error':
        return 'bg-red-100 text-red-800 border-red-200'
      default:
        return 'bg-gray-100 text-gray-800 border-gray-200'
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'healthy':
        return '✅'
      case 'ready':
        return '🟦'
      case 'checking':
        return '🟡'
      case 'error':
        return '❌'
      default:
        return '⚪'
    }
  }

  return (
    <div className="bg-white rounded-xl shadow-lg border-2 border-gray-100 p-4 hover:shadow-xl hover:border-gray-200 transition-all duration-300 hover:-translate-y-1">
      <div className="flex items-center justify-between mb-2">
        <h3 className="font-medium text-gray-900">{title}</h3>
        <span className="text-lg">{getStatusIcon(status)}</span>
      </div>
      <div className={`px-2 py-1 rounded-full text-xs font-medium border ${getStatusColor(status)} inline-block`}>
        {status.charAt(0).toUpperCase() + status.slice(1)}
      </div>
      <p className="text-sm text-gray-600 mt-2">{description}</p>
    </div>
  )
}
