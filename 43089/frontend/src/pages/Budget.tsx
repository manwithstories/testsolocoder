import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { activityAPI, planAPI } from '@/services/api'
import type { BudgetSummary, PlanDetail, Activity } from '@/types'
import dayjs from 'dayjs'

const categoryColors: Record<string, string> = {
  sightseeing: 'bg-blue-500',
  transport: 'bg-yellow-500',
  accommodation: 'bg-purple-500',
  food: 'bg-orange-500',
  other: 'bg-gray-500',
}

const categoryLabels: Record<string, string> = {
  sightseeing: '景点',
  transport: '交通',
  accommodation: '住宿',
  food: '餐饮',
  other: '其他',
}

export default function Budget() {
  const { id } = useParams<{ id: string }>()
  const [plan, setPlan] = useState<PlanDetail | null>(null)
  const [budgetSummary, setBudgetSummary] = useState<BudgetSummary | null>(null)
  const [activities, setActivities] = useState<Activity[]>([])
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState<'overview' | 'details'>('overview')

  useEffect(() => {
    if (id) {
      loadData()
    }
  }, [id])

  const loadData = async () => {
    try {
      setLoading(true)
      const [planData, summaryData, activitiesData] = await Promise.all([
        planAPI.getPlan(id!),
        activityAPI.getBudgetSummary(id!),
        activityAPI.getActivities(id!),
      ])
      setPlan(planData)
      setBudgetSummary(summaryData)
      setActivities(activitiesData as Activity[])
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      setLoading(false)
    }
  }

  const getBudgetProgressColor = (usage: number) => {
    if (usage >= 100) return 'bg-red-500'
    if (usage >= 80) return 'bg-yellow-500'
    return 'bg-green-500'
  }

  const getExpensesByCategory = () => {
    const expenses: Record<string, Activity[]> = {}
    activities.forEach((activity) => {
      if (!expenses[activity.type]) {
        expenses[activity.type] = []
      }
      expenses[activity.type].push(activity)
    })
    return expenses
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  if (!plan || !budgetSummary) {
    return (
      <div className="p-6 text-center">
        <p className="text-gray-500">数据加载失败</p>
        <Link to={`/plans/${id}`} className="text-primary-600 hover:underline">
          返回计划详情
        </Link>
      </div>
    )
  }

  const expensesByCategory = getExpensesByCategory()
  const paidActivities = activities.filter((a) => a.cost > 0).sort((a, b) => b.cost - a.cost)

  return (
    <div className="p-6">
      <div className="mb-6">
        <Link to={`/plans/${id}`} className="text-primary-600 hover:underline text-sm mb-2 inline-block">
          ← 返回计划详情
        </Link>
        <h1 className="text-2xl font-bold text-gray-900">预算追踪</h1>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <div className="bg-white rounded-xl shadow-sm p-4">
          <p className="text-sm text-gray-500 mb-1">总预算</p>
          <p className="text-2xl font-bold text-gray-900">
            {plan.budget.toLocaleString()} {plan.currency}
          </p>
        </div>
        <div className="bg-white rounded-xl shadow-sm p-4">
          <p className="text-sm text-gray-500 mb-1">已花费</p>
          <p className="text-2xl font-bold text-orange-600">
            {budgetSummary.total_spent.toLocaleString()} {budgetSummary.plan_currency}
          </p>
        </div>
        <div className="bg-white rounded-xl shadow-sm p-4">
          <p className="text-sm text-gray-500 mb-1">剩余预算</p>
          <p className={`text-2xl font-bold ${budgetSummary.budget_remaining >= 0 ? 'text-green-600' : 'text-red-600'}`}>
            {budgetSummary.budget_remaining.toLocaleString()} {budgetSummary.plan_currency}
          </p>
        </div>
        <div className="bg-white rounded-xl shadow-sm p-4">
          <p className="text-sm text-gray-500 mb-1">使用率</p>
          <p className="text-2xl font-bold text-gray-900">{budgetSummary.budget_usage.toFixed(1)}%</p>
        </div>
      </div>

      <div className="bg-white rounded-xl shadow-sm mb-6">
        <div className="p-6">
          <h2 className="text-lg font-semibold mb-4">预算使用情况</h2>
          <div className="w-full bg-gray-200 rounded-full h-4 mb-2">
            <div
              className={`h-4 rounded-full transition-all ${getBudgetProgressColor(budgetSummary.budget_usage)}`}
              style={{ width: `${Math.min(budgetSummary.budget_usage, 100)}%` }}
            ></div>
          </div>
          <div className="flex justify-between text-sm text-gray-500">
            <span>0%</span>
            <span>
              {budgetSummary.total_spent.toLocaleString()} / {plan.budget.toLocaleString()} {plan.currency}
            </span>
            <span>100%</span>
          </div>
          {budgetSummary.budget_usage >= 100 && (
            <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
              ⚠️ 警告：已超出预算！请控制后续支出。
            </div>
          )}
          {budgetSummary.budget_usage >= 80 && budgetSummary.budget_usage < 100 && (
            <div className="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded-lg text-yellow-700 text-sm">
              ⚡ 注意：预算使用率已超过80%，请合理安排后续支出。
            </div>
          )}
        </div>
      </div>

      <div className="flex gap-4 mb-6">
        <button
          onClick={() => setActiveTab('overview')}
          className={`px-4 py-2 rounded-lg font-medium transition-colors ${
            activeTab === 'overview'
              ? 'bg-primary-600 text-white'
              : 'bg-white border border-gray-200 hover:bg-gray-50'
          }`}
        >
          按类别统计
        </button>
        <button
          onClick={() => setActiveTab('details')}
          className={`px-4 py-2 rounded-lg font-medium transition-colors ${
            activeTab === 'details'
              ? 'bg-primary-600 text-white'
              : 'bg-white border border-gray-200 hover:bg-gray-50'
          }`}
        >
          支出明细
        </button>
      </div>

      {activeTab === 'overview' && (
        <div className="bg-white rounded-xl shadow-sm">
          <div className="p-6">
            <h2 className="text-lg font-semibold mb-4">各类别支出</h2>
            {budgetSummary.by_category.length === 0 ? (
              <div className="text-center py-12 text-gray-500">
                <p className="text-4xl mb-2">💰</p>
                <p>暂无支出记录</p>
              </div>
            ) : (
              <div className="space-y-4">
                {budgetSummary.by_category.map((item) => {
                  const percentage = budgetSummary.total_spent > 0 
                    ? (item.total / budgetSummary.total_spent) * 100 
                    : 0
                  return (
                    <div key={item.type} className="border border-gray-200 rounded-lg p-4">
                      <div className="flex items-center justify-between mb-2">
                        <div className="flex items-center gap-3">
                          <div className={`w-3 h-3 rounded-full ${categoryColors[item.type] || 'bg-gray-500'}`}></div>
                          <span className="font-medium">{categoryLabels[item.type] || item.type}</span>
                          <span className="text-sm text-gray-500">({item.count} 项)</span>
                        </div>
                        <span className="font-bold">
                          {item.total.toLocaleString()} {budgetSummary.plan_currency}
                        </span>
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2">
                        <div
                          className={`h-2 rounded-full ${categoryColors[item.type] || 'bg-gray-500'}`}
                          style={{ width: `${percentage}%` }}
                        ></div>
                      </div>
                      <p className="text-right text-xs text-gray-500 mt-1">{percentage.toFixed(1)}%</p>
                    </div>
                  )
                })}
              </div>
            )}
          </div>
        </div>
      )}

      {activeTab === 'details' && (
        <div className="bg-white rounded-xl shadow-sm">
          <div className="p-6">
            <h2 className="text-lg font-semibold mb-4">支出明细</h2>
            {paidActivities.length === 0 ? (
              <div className="text-center py-12 text-gray-500">
                <p className="text-4xl mb-2">📝</p>
                <p>暂无支出记录</p>
              </div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-gray-200">
                      <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">日期</th>
                      <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">活动名称</th>
                      <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">类别</th>
                      <th className="text-right py-3 px-4 text-sm font-medium text-gray-500">金额</th>
                    </tr>
                  </thead>
                  <tbody>
                    {paidActivities.map((activity) => (
                      <tr key={activity.id} className="border-b border-gray-100 hover:bg-gray-50">
                        <td className="py-3 px-4 text-sm">
                          {dayjs(activity.date).format('YYYY/MM/DD')}
                        </td>
                        <td className="py-3 px-4">
                          <p className="font-medium">{activity.title}</p>
                          {activity.location && (
                            <p className="text-xs text-gray-500">📍 {activity.location}</p>
                          )}
                        </td>
                        <td className="py-3 px-4">
                          <span className={`px-2 py-1 rounded text-xs font-medium ${categoryColors[activity.type]} text-white`}>
                            {categoryLabels[activity.type] || activity.type}
                          </span>
                        </td>
                        <td className="py-3 px-4 text-right font-medium">
                          {activity.cost.toLocaleString()} {activity.currency || plan.currency}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                  <tfoot>
                    <tr className="bg-gray-50">
                      <td colSpan={3} className="py-3 px-4 text-right font-medium">
                        总计：
                      </td>
                      <td className="py-3 px-4 text-right font-bold text-lg">
                        {budgetSummary.total_spent.toLocaleString()} {budgetSummary.plan_currency}
                      </td>
                    </tr>
                  </tfoot>
                </table>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}
