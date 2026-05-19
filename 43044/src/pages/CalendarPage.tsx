import { useState } from 'react'
import { ChevronLeft, ChevronRight, Calendar as CalendarIcon } from 'lucide-react'
import { useAppStore } from '@/store'
import { Button } from '@/components/ui/Button'
import { Badge } from '@/components/ui/Badge'
import { Modal } from '@/components/ui/Modal'
import {
  getMonthDays,
  getWeekDays,
  isCurrentMonth,
  isTodayDate,
  isSameDayDate,
  nextMonth,
  prevMonth,
  nextWeek,
  prevWeek,
  getMonthLabel,
  getWeekLabel,
  formatDateReadable,
} from '@/utils/date'
import { getTasksByDate, getPriorityColor, isOverdue } from '@/utils/progress'
import type { Task } from '@/types'

type ViewMode = 'month' | 'week'

export function CalendarPage() {
  const { tasks } = useAppStore()
  const [currentDate, setCurrentDate] = useState(new Date())
  const [viewMode, setViewMode] = useState<ViewMode>('month')
  const [selectedDate, setSelectedDate] = useState<Date | null>(null)
  const [isDayModalOpen, setIsDayModalOpen] = useState(false)
  
  const days = viewMode === 'month' ? getMonthDays(currentDate) : getWeekDays(currentDate)
  const label = viewMode === 'month' ? getMonthLabel(currentDate) : getWeekLabel(currentDate)
  const weekDays = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
  
  const handlePrev = () => {
    setCurrentDate(viewMode === 'month' ? prevMonth(currentDate) : prevWeek(currentDate))
  }
  
  const handleNext = () => {
    setCurrentDate(viewMode === 'month' ? nextMonth(currentDate) : nextWeek(currentDate))
  }
  
  const handleDayClick = (day: Date) => {
    setSelectedDate(day)
    setIsDayModalOpen(true)
  }
  
  const dayTasks = selectedDate ? getTasksByDate(tasks, selectedDate) : []
  
  const getDayTasks = (day: Date): Task[] => {
    return getTasksByDate(tasks, day)
  }
  
  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-800">日历视图</h1>
          <p className="text-gray-500 mt-1">查看任务截止日期和安排</p>
        </div>
        <div className="flex items-center gap-2">
          <div className="flex bg-gray-100 rounded-lg p-1">
            <button
              onClick={() => setViewMode('month')}
              className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                viewMode === 'month'
                  ? 'bg-white text-blue-600 shadow-sm'
                  : 'text-gray-600 hover:text-gray-800'
              }`}
            >
              月视图
            </button>
            <button
              onClick={() => setViewMode('week')}
              className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                viewMode === 'week'
                  ? 'bg-white text-blue-600 shadow-sm'
                  : 'text-gray-600 hover:text-gray-800'
              }`}
            >
              周视图
            </button>
          </div>
        </div>
      </div>
      
      <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <div className="flex items-center justify-between p-4 border-b border-gray-200">
          <Button variant="ghost" size="sm" onClick={handlePrev}>
            <ChevronLeft className="w-5 h-5" />
          </Button>
          <h2 className="text-lg font-semibold text-gray-800">{label}</h2>
          <Button variant="ghost" size="sm" onClick={handleNext}>
            <ChevronRight className="w-5 h-5" />
          </Button>
        </div>
        
        <div className="grid grid-cols-7 border-b border-gray-200">
          {weekDays.map((day) => (
            <div
              key={day}
              className="py-3 text-center text-sm font-medium text-gray-600"
            >
              {day}
            </div>
          ))}
        </div>
        
        <div className="grid grid-cols-7">
          {days.map((day, index) => {
            const dayTasksList = getDayTasks(day)
            const isCurrentMonthDay = viewMode === 'month' ? isCurrentMonth(day, currentDate) : true
            const isToday = isTodayDate(day)
            const isSelected = selectedDate && isSameDayDate(day, selectedDate)
            const hasOverdue = dayTasksList.some((t) => isOverdue(t))
            
            return (
              <div
                key={index}
                onClick={() => handleDayClick(day)}
                className={`min-h-[120px] p-2 border-b border-r border-gray-100 cursor-pointer transition-colors hover:bg-gray-50 ${
                  !isCurrentMonthDay ? 'bg-gray-50' : ''
                } ${isSelected ? 'bg-blue-50' : ''} ${
                  index % 7 === 6 ? 'border-r-0' : ''
                }`}
              >
                <div className="flex items-center justify-between mb-1">
                  <span
                    className={`text-sm font-medium ${
                      isToday
                        ? 'bg-blue-600 text-white w-7 h-7 rounded-full flex items-center justify-center'
                        : isCurrentMonthDay
                        ? 'text-gray-800'
                        : 'text-gray-400'
                    }`}
                  >
                    {day.getDate()}
                  </span>
                  {hasOverdue && (
                    <span className="w-2 h-2 bg-red-500 rounded-full" />
                  )}
                </div>
                
                <div className="space-y-1">
                  {dayTasksList.slice(0, 3).map((task) => (
                    <div
                      key={task.id}
                      className={`text-xs px-1.5 py-0.5 rounded truncate ${
                        task.status === 'completed'
                          ? 'bg-green-100 text-green-700 line-through'
                          : isOverdue(task)
                          ? 'bg-red-100 text-red-700'
                          : 'bg-blue-100 text-blue-700'
                      }`}
                      title={task.title}
                    >
                      {task.title}
                    </div>
                  ))}
                  {dayTasksList.length > 3 && (
                    <div className="text-xs text-gray-500 text-center">
                      +{dayTasksList.length - 3} 更多
                    </div>
                  )}
                </div>
              </div>
            )
          })}
        </div>
      </div>
      
      <Modal
        isOpen={isDayModalOpen}
        onClose={() => setIsDayModalOpen(false)}
        title={selectedDate ? formatDateReadable(selectedDate) : ''}
      >
        {dayTasks.length === 0 ? (
          <div className="text-center py-8">
            <CalendarIcon className="w-12 h-12 text-gray-300 mx-auto mb-3" />
            <p className="text-gray-500">当天没有任务</p>
          </div>
        ) : (
          <div className="space-y-3">
            {dayTasks.map((task) => (
              <div
                key={task.id}
                className={`p-3 rounded-lg border ${
                  task.status === 'completed'
                    ? 'bg-green-50 border-green-200'
                    : isOverdue(task)
                    ? 'bg-red-50 border-red-200'
                    : 'bg-gray-50 border-gray-200'
                }`}
              >
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <h4 className={`font-medium ${
                      task.status === 'completed' ? 'line-through text-gray-500' : 'text-gray-800'
                    }`}>
                      {task.title}
                    </h4>
                    {task.description && (
                      <p className="text-sm text-gray-500 mt-1">{task.description}</p>
                    )}
                  </div>
                  <div className="flex flex-col items-end gap-1">
                    <Badge
                      text={task.status === 'completed' ? '已完成' : isOverdue(task) ? '已逾期' : '待完成'}
                      color={
                        task.status === 'completed'
                          ? 'bg-green-500 text-white'
                          : isOverdue(task)
                          ? 'bg-red-500 text-white'
                          : 'bg-blue-500 text-white'
                      }
                      size="sm"
                    />
                    <Badge
                      text={
                        task.priority === 'urgent' ? '紧急' :
                        task.priority === 'high' ? '高' :
                        task.priority === 'medium' ? '中' : '低'
                      }
                      color={getPriorityColor(task.priority)}
                      size="sm"
                    />
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </Modal>
      
      <div className="mt-6 flex flex-wrap gap-4 text-sm">
        <div className="flex items-center gap-2">
          <span className="w-3 h-3 rounded-full bg-blue-100 border border-blue-300" />
          <span className="text-gray-600">待完成</span>
        </div>
        <div className="flex items-center gap-2">
          <span className="w-3 h-3 rounded-full bg-green-100 border border-green-300" />
          <span className="text-gray-600">已完成</span>
        </div>
        <div className="flex items-center gap-2">
          <span className="w-3 h-3 rounded-full bg-red-100 border border-red-300" />
          <span className="text-gray-600">已逾期</span>
        </div>
        <div className="flex items-center gap-2">
          <span className="w-3 h-3 rounded-full bg-blue-600" />
          <span className="text-gray-600">今天</span>
        </div>
      </div>
    </div>
  )
}
