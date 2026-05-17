interface ProgressBarProps {
  percentage: number
  showLabel?: boolean
  size?: 'sm' | 'md' | 'lg'
  color?: string
}

export function ProgressBar({ percentage, showLabel = true, size = 'md', color }: ProgressBarProps) {
  const heightClass = size === 'sm' ? 'h-1' : size === 'lg' ? 'h-3' : 'h-2'
  const bgColor = color || (percentage === 100 ? 'bg-green-500' : percentage >= 70 ? 'bg-blue-500' : percentage >= 30 ? 'bg-yellow-500' : 'bg-red-500')
  
  return (
    <div className="w-full">
      <div className={`w-full bg-gray-200 rounded-full ${heightClass} overflow-hidden`}>
        <div
          className={`${heightClass} ${bgColor} rounded-full transition-all duration-300`}
          style={{ width: `${Math.min(100, Math.max(0, percentage))}%` }}
        />
      </div>
      {showLabel && (
        <div className="flex justify-between text-xs text-gray-500 mt-1">
          <span>进度</span>
          <span className="font-medium">{percentage}%</span>
        </div>
      )}
    </div>
  )
}
