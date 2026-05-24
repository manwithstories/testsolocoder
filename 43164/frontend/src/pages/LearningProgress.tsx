import { useState, useEffect } from 'react'
import { studentApi, notesApi, homeworkApi } from '@/services/api'
import { LearningGoal, LessonNote, Homework, Milestone } from '@/types'
import { Target, BookOpen, FileText, Award, Plus, CheckCircle } from 'lucide-react'

export default function LearningProgress() {
  const [goals, setGoals] = useState<LearningGoal[]>([])
  const [notes, setNotes] = useState<LessonNote[]>([])
  const [homework, setHomework] = useState<Homework[]>([])
  const [milestones, setMilestones] = useState<Milestone[]>([])
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState('goals')

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      setLoading(true)
      const [profileRes, notesRes, homeworkRes, milestonesRes] = await Promise.all([
        studentApi.getProfile(),
        notesApi.getAll(),
        homeworkApi.getAll(),
        studentApi.getMilestones(),
      ])
      setGoals(profileRes.data?.learningGoals || [])
      setNotes(notesRes.data || [])
      setHomework(homeworkRes.data || [])
      setMilestones(milestonesRes.data || [])
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      setLoading(false)
    }
  }

  const tabs = [
    { id: 'goals', label: '学习目标', icon: Target },
    { id: 'notes', label: '课程笔记', icon: BookOpen },
    { id: 'homework', label: '作业', icon: FileText },
    { id: 'milestones', label: '里程碑', icon: Award },
  ]

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">学习进度</h1>
        <p className="text-gray-500">追踪您的学习目标和成果</p>
      </div>

      <div className="flex gap-2 border-b border-gray-200">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`flex items-center gap-2 px-4 py-2 text-sm font-medium border-b-2 -mb-px transition-colors ${
              activeTab === tab.id
                ? 'border-primary-600 text-primary-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            <tab.icon className="h-4 w-4" />
            {tab.label}
          </button>
        ))}
      </div>

      {activeTab === 'goals' && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h2 className="text-lg font-semibold text-gray-900">学习目标</h2>
            <button className="btn-primary flex items-center gap-1">
              <Plus className="h-4 w-4" />
              添加目标
            </button>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {goals.map((goal) => (
              <div key={goal.id} className="card">
                <div className="flex items-center justify-between mb-3">
                  <h3 className="font-semibold text-gray-900">{goal.title}</h3>
                  {goal.isAchieved && (
                    <span className="badge bg-green-100 text-green-800">已达成</span>
                  )}
                </div>
                {goal.description && (
                  <p className="text-sm text-gray-500 mb-3">{goal.description}</p>
                )}
                <div className="flex items-center gap-2 mb-2">
                  <div className="flex-1 bg-gray-200 rounded-full h-2">
                    <div
                      className="bg-primary-600 h-2 rounded-full"
                      style={{ width: `${Math.min((goal.currentScore / (goal.targetScore || 100)) * 100, 100)}%` }}
                    />
                  </div>
                  <span className="text-sm font-medium text-gray-600">
                    {goal.currentScore}/{goal.targetScore}
                  </span>
                </div>
                {goal.deadline && (
                  <p className="text-xs text-gray-400">
                    截止日期: {new Date(goal.deadline).toLocaleDateString('zh-CN')}
                  </p>
                )}
              </div>
            ))}
            {goals.length === 0 && (
              <div className="col-span-2 text-center py-12">
                <Target className="h-12 w-12 text-gray-300 mx-auto mb-4" />
                <p className="text-gray-500">暂无学习目标</p>
              </div>
            )}
          </div>
        </div>
      )}

      {activeTab === 'notes' && (
        <div className="space-y-4">
          <h2 className="text-lg font-semibold text-gray-900">课程笔记</h2>
          <div className="space-y-4">
            {notes.map((note) => (
              <div key={note.id} className="card">
                <div className="flex items-center justify-between mb-3">
                  <h3 className="font-semibold text-gray-900">{note.title}</h3>
                  <span className="text-sm text-gray-400">
                    {new Date(note.createdAt).toLocaleDateString('zh-CN')}
                  </span>
                </div>
                <p className="text-gray-600 whitespace-pre-wrap">{note.content}</p>
                {note.tags && (
                  <div className="flex flex-wrap gap-1 mt-3">
                    {note.tags.split(',').map((tag, index) => (
                      <span key={index} className="badge bg-gray-100 text-gray-600">
                        {tag.trim()}
                      </span>
                    ))}
                  </div>
                )}
              </div>
            ))}
            {notes.length === 0 && (
              <div className="text-center py-12">
                <BookOpen className="h-12 w-12 text-gray-300 mx-auto mb-4" />
                <p className="text-gray-500">暂无课程笔记</p>
              </div>
            )}
          </div>
        </div>
      )}

      {activeTab === 'homework' && (
        <div className="space-y-4">
          <h2 className="text-lg font-semibold text-gray-900">作业</h2>
          <div className="space-y-4">
            {homework.map((hw) => (
              <div key={hw.id} className="card">
                <div className="flex items-center justify-between mb-3">
                  <h3 className="font-semibold text-gray-900">{hw.title}</h3>
                  <span className={`badge ${
                    hw.status === 'graded' ? 'bg-green-100 text-green-800' :
                    hw.status === 'submitted' ? 'bg-blue-100 text-blue-800' :
                    'bg-yellow-100 text-yellow-800'
                  }`}>
                    {hw.status === 'graded' ? '已批改' :
                     hw.status === 'submitted' ? '已提交' : '待提交'}
                  </span>
                </div>
                <p className="text-gray-600 mb-2">{hw.description}</p>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-gray-400">
                    截止日期: {new Date(hw.dueDate).toLocaleDateString('zh-CN')}
                  </span>
                  {hw.submission?.score !== undefined && (
                    <span className="font-medium text-green-600">
                      得分: {hw.submission.score}/{hw.maxScore}
                    </span>
                  )}
                </div>
                {hw.submission?.feedback && (
                  <div className="mt-3 p-3 bg-gray-50 rounded-lg">
                    <p className="text-sm text-gray-600">{hw.submission.feedback}</p>
                  </div>
                )}
              </div>
            ))}
            {homework.length === 0 && (
              <div className="text-center py-12">
                <FileText className="h-12 w-12 text-gray-300 mx-auto mb-4" />
                <p className="text-gray-500">暂无作业</p>
              </div>
            )}
          </div>
        </div>
      )}

      {activeTab === 'milestones' && (
        <div className="space-y-4">
          <h2 className="text-lg font-semibold text-gray-900">学习里程碑</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {milestones.map((milestone) => (
              <div
                key={milestone.id}
                className={`card text-center ${milestone.isAchieved ? 'bg-green-50' : ''}`}
              >
                <div className={`w-16 h-16 rounded-full mx-auto mb-3 flex items-center justify-center ${
                  milestone.isAchieved ? 'bg-green-100' : 'bg-gray-100'
                }`}>
                  {milestone.isAchieved ? (
                    <CheckCircle className="h-8 w-8 text-green-600" />
                  ) : (
                    <Award className="h-8 w-8 text-gray-400" />
                  )}
                </div>
                <h3 className="font-semibold text-gray-900 mb-1">{milestone.title}</h3>
                {milestone.description && (
                  <p className="text-sm text-gray-500">{milestone.description}</p>
                )}
              </div>
            ))}
            {milestones.length === 0 && (
              <div className="col-span-4 text-center py-12">
                <Award className="h-12 w-12 text-gray-300 mx-auto mb-4" />
                <p className="text-gray-500">暂无里程碑</p>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}
