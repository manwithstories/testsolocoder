interface ProgressBarProps {
  progress: number
  label?: string
  showPercentage?: boolean
  size?: 'sm' | 'md' | 'lg'
}

export default function ProgressBar({
  progress,
  label,
  showPercentage = true,
  size = 'md',
}: ProgressBarProps) {
  const heightClass = {
    sm: 'h-1',
    md: 'h-2',
    lg: 'h-3',
  }[size]

  const clampedProgress = Math.max(0, Math.min(100, progress))

  return (
    <div className="w-full">
      {(label || showPercentage) && (
        <div className="flex justify-between items-center mb-1">
          {label && <span className="text-sm text-gray-600">{label}</span>}
          {showPercentage && (
            <span className="text-sm text-gray-600">{clampedProgress.toFixed(0)}%</span>
          )}
        </div>
      )}
      <div className={`w-full bg-gray-200 rounded-full overflow-hidden ${heightClass}`}>
        <div
          className="bg-indigo-600 rounded-full transition-all duration-300 ease-out h-full"
          style={{ width: `${clampedProgress}%` }}
        />
      </div>
    </div>
  )
}
