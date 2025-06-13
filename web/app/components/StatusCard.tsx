interface StatusCardProps {
  title: string
  status: 'healthy' | 'checking' | 'error' | 'ready'
  description: string
}

export default function StatusCard({ title, status, description }: StatusCardProps) {
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'healthy':
        return 'badge-success'
      case 'ready':
        return 'badge-info'
      case 'checking':
        return 'badge-warning'
      case 'error':
        return 'badge-error'
      default:
        return 'bg-gray-100 text-gray-800 border-gray-200'
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'healthy':
        return '✅'
      case 'ready':
        return '🚀'
      case 'checking':
        return '⏳'
      case 'error':
        return '❌'
      default:
        return '⚪'
    }
  }

  const getGradientIcon = (status: string) => {
    switch (status) {
      case 'healthy':
        return 'bg-gradient-to-br from-green-500 to-emerald-600'
      case 'ready':
        return 'bg-gradient-to-br from-blue-500 to-indigo-600'
      case 'checking':
        return 'bg-gradient-to-br from-amber-500 to-orange-600'
      case 'error':
        return 'bg-gradient-to-br from-red-500 to-rose-600'
      default:
        return 'bg-gradient-to-br from-gray-500 to-gray-600'
    }
  }

  return (
    <div className="modern-card p-6 fade-in-up">
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center space-x-3">
          <div className={`w-8 h-8 ${getGradientIcon(status)} rounded-xl flex items-center justify-center shadow-md`}>
            <span className="text-white text-sm">{getStatusIcon(status)}</span>
          </div>
          <h3 className="font-semibold text-gray-900">{title}</h3>
        </div>
      </div>
      
      <div className={`badge ${getStatusColor(status)} mb-3`}>
        {status.charAt(0).toUpperCase() + status.slice(1)}
      </div>
      
      <p className="text-sm text-gray-600 leading-relaxed">{description}</p>
    </div>
  )
}
