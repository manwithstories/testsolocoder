import { useEffect, useState } from 'react';
import { dashboardAPI } from '../api';
import type { DashboardStats } from '../types';
import {
  DollarSign,
  Users,
  FolderKanban,
  Clock,
  AlertTriangle,
  FileText,
} from 'lucide-react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

export default function Dashboard() {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const response = await dashboardAPI.getStats();
        if (response.data.success && response.data.data) {
          setStats(response.data.data);
        }
      } catch (err: any) {
        setError(err.response?.data?.error || 'Failed to load stats');
      } finally {
        setLoading(false);
      }
    };
    fetchStats();
  }, []);

  if (loading) {
    return <div className="text-center py-12">Loading...</div>;
  }

  if (error) {
    return <div className="text-center py-12 text-red-500">{error}</div>;
  }

  if (!stats) {
    return <div className="text-center py-12">No data available</div>;
  }

  const statCards = [
    {
      label: 'Monthly Revenue',
      value: `$${stats.monthly_revenue.toFixed(2)}`,
      icon: DollarSign,
      color: 'bg-green-500',
    },
    {
      label: 'Total Clients',
      value: stats.total_clients.toString(),
      icon: Users,
      color: 'bg-blue-500',
    },
    {
      label: 'Active Projects',
      value: stats.active_projects.toString(),
      icon: FolderKanban,
      color: 'bg-purple-500',
    },
    {
      label: 'Hours This Month',
      value: stats.total_hours.toFixed(1),
      icon: Clock,
      color: 'bg-orange-500',
    },
    {
      label: 'Overdue Invoices',
      value: stats.overdue_invoices.toString(),
      icon: AlertTriangle,
      color: 'bg-red-500',
    },
    {
      label: 'Pending Invoices',
      value: stats.pending_invoices.toString(),
      icon: FileText,
      color: 'bg-yellow-500',
    },
  ];

  return (
    <div className="space-y-8">
      <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {statCards.map((card, index) => (
          <div key={index} className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center">
              <div className={`${card.color} rounded-lg p-3`}>
                <card.icon className="w-6 h-6 text-white" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-500">{card.label}</p>
                <p className="text-2xl font-bold text-gray-900">{card.value}</p>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Monthly Earnings</h2>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={stats.monthly_earnings}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="month" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="revenue" fill="#4f46e5" name="Revenue ($)" />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Project Progress</h2>
          <div className="space-y-4">
            {stats.project_progress.map((project) => (
              <div key={project.id}>
                <div className="flex justify-between mb-1">
                  <span className="text-sm font-medium text-gray-700">{project.name}</span>
                  <span className="text-sm text-gray-500">{project.progress.toFixed(0)}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-indigo-600 h-2 rounded-full transition-all"
                    style={{ width: `${Math.min(project.progress, 100)}%` }}
                  />
                </div>
                <div className="flex justify-between mt-1">
                  <span className="text-xs text-gray-500">
                    ${project.earned_amount.toFixed(0)} / ${project.budget.toFixed(0)}
                  </span>
                  <span className={`text-xs ${project.days_left < 0 ? 'text-red-500' : 'text-gray-500'}`}>
                    {project.days_left >= 0 ? `${project.days_left} days left` : `${Math.abs(project.days_left)} days overdue`}
                  </span>
                </div>
              </div>
            ))}
            {stats.project_progress.length === 0 && (
              <p className="text-gray-500 text-center py-4">No active projects</p>
            )}
          </div>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-xl font-semibold text-gray-900 mb-4">Overdue Reminders</h2>
        <div className="space-y-3">
          {stats.overdue_reminders.map((reminder) => (
            <div
              key={`${reminder.type}-${reminder.id}`}
              className="flex items-center justify-between p-3 bg-red-50 rounded-lg"
            >
              <div className="flex items-center">
                <AlertTriangle className="w-5 h-5 text-red-500 mr-3" />
                <div>
                  <p className="font-medium text-gray-900">{reminder.title}</p>
                  <p className="text-sm text-gray-500">
                    {reminder.days_overdue} days overdue
                    {reminder.amount && ` - $${reminder.amount.toFixed(2)}`}
                  </p>
                </div>
              </div>
              <span className="px-3 py-1 bg-red-100 text-red-700 rounded-full text-sm">
                {reminder.type}
              </span>
            </div>
          ))}
          {stats.overdue_reminders.length === 0 && (
            <p className="text-gray-500 text-center py-4">No overdue items</p>
          )}
        </div>
      </div>
    </div>
  );
}
