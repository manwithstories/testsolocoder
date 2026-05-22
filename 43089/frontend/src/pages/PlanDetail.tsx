import { useEffect, useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { planAPI } from '@/services/api'
import type { PlanDetail } from '@/types'
import dayjs from 'dayjs'

export default function PlanDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [plan, setPlan] = useState<PlanDetail | null>(null)
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState('overview')

  useEffect(() => {
    if (id) {
      loadPlan()
    }
  }, [id])

  const loadPlan = async () => {
    try {
      setLoading(true)
      const data = await planAPI.getPlan(id!)
      setPlan(data)
    } catch (error) {
      console.error('Failed to load plan:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleExportJSON = async () => {
    try {
      const blob = await planAPI.exportJSON(id!)
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `${plan?.title}_travel_plan.json`
      a.click()
      URL.revokeObjectURL(url)
    } catch (error: any) {
      alert(error.message || '导出失败')
    }
  }

  const handleExportPDF = async () => {
    try {
      const blob = await planAPI.exportPDF(id!)
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `${plan?.title}_travel_plan.pdf`
      a.click()
      URL.revokeObjectURL(url)
    } catch (error: any) {
      alert(error.message || '导出失败')
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  if (!plan) {
    return (
      <div className="p-6 text-center">
        <p className="text-gray-500">计划不存在</p>
        <Link to="/plans" className="text-primary-600 hover:underline">
          返回计划列表
        </Link>
      </div>
    )
  }

  const tabs = [
    { id: 'overview', label: '概览', icon: '📋' },
    { id: 'activities', label: '行程安排', icon: '📅' },
    { id: 'budget', label: '预算追踪', icon: '💰' },
    { id: 'files', label: '文件附件', icon: '📁' },
    { id: 'checklist', label: '清单管理', icon: '✅' },
    { id: 'map', label: '地图', icon: '🗺️' },
  ]

  return (
    <div className="p-6">
      <div className="mb-6">
        <Link to="/plans" className="text-primary-600 hover:underline text-sm mb-2 inline-block">
          ← 返回计划列表
        </Link>
        <div className="flex items-start justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">{plan.title}</h1>
            <div className="flex items-center gap-4 mt-2 text-gray-500">
              <span>📍 {plan.destination}</span>
              <span>
                📅 {dayjs(plan.start_date).format('YYYY/MM/DD')} - {dayjs(plan.end_date).format('YYYY/MM/DD')}
              </span>
              <span>💰 {plan.budget.toLocaleString()} {plan.currency}</span>
            </div>
          </div>
          <div className="flex gap-2">
            <button
              onClick={handleExportJSON}
              className="px-3 py-2 text-sm border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
            >
              📤 导出JSON
            </button>
            <button
              onClick={handleExportPDF}
              className="px-3 py-2 text-sm border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
            >
              📄 导出PDF
            </button>
          </div>
        </div>
      </div>

      <div className="bg-white rounded-xl shadow-sm mb-6">
        <div className="flex border-b overflow-x-auto">
          {tabs.map((tab) => (
            <Link
              key={tab.id}
              to={tab.id === 'overview' ? `/plans/${id}` : `/plans/${id}/${tab.id}`}
              onClick={(e) => {
                e.preventDefault()
                setActiveTab(tab.id)
                navigate(tab.id === 'overview' ? `/plans/${id}` : `/plans/${id}/${tab.id}`)
              }}
              className={`px-6 py-3 text-sm font-medium whitespace-nowrap transition-colors ${
                activeTab === tab.id ||
                (tab.id === 'overview' &&
                  !['activities', 'budget', 'files', 'checklist', 'map'].includes(
                    window.location.pathname.split('/').pop() || ''
                  ))
                  ? 'text-primary-600 border-b-2 border-primary-600'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              <span className="mr-2">{tab.icon}</span>
              {tab.label}
            </Link>
          ))}
        </div>

        <div className="p-6">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <div className="bg-blue-50 rounded-xl p-4">
              <p className="text-sm text-blue-600 mb-1">总预算</p>
              <p className="text-2xl font-bold text-blue-700">
                {plan.budget.toLocaleString()} {plan.currency}
              </p>
            </div>
            <div className="bg-green-50 rounded-xl p-4">
              <p className="text-sm text-green-600 mb-1">已花费</p>
              <p className="text-2xl font-bold text-green-700">
                {plan.total_spent.toLocaleString()} {plan.currency}
              </p>
            </div>
            <div className="bg-purple-50 rounded-xl p-4">
              <p className="text-sm text-purple-600 mb-1">活动数量</p>
              <p className="text-2xl font-bold text-purple-700">{plan.activity_count}</p>
            </div>
          </div>

          <div className="mb-8">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">计划描述</h3>
            <p className="text-gray-600">{plan.description || '暂无描述'}</p>
          </div>

          <div>
            <h3 className="text-lg font-semibold text-gray-900 mb-4">参与人员 ({plan.participant_count})</h3>
            <div className="flex flex-wrap gap-3">
              <div className="flex items-center bg-gray-100 rounded-full px-4 py-2">
                <div className="w-8 h-8 bg-primary-500 text-white rounded-full flex items-center justify-center text-sm font-medium mr-2">
                  {plan.owner?.first_name?.charAt(0) || plan.owner?.username?.charAt(0) || 'O'}
                </div>
                <span className="text-sm">
                  {plan.owner?.first_name || plan.owner?.username} (所有者)
                </span>
              </div>
              {plan.participants?.map((p) => (
                <div key={p.id} className="flex items-center bg-gray-100 rounded-full px-4 py-2">
                  <div className="w-8 h-8 bg-gray-500 text-white rounded-full flex items-center justify-center text-sm font-medium mr-2">
                    {p.user?.first_name?.charAt(0) || p.user?.username?.charAt(0) || 'U'}
                  </div>
                  <span className="text-sm">
                    {p.user?.first_name || p.user?.username} ({p.role})
                  </span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
