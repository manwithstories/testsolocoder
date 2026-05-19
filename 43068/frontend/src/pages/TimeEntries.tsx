import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { timeEntryAPI, projectAPI } from '../api';
import type { TimeEntry, Project } from '../types';
import { Plus, Edit2, Trash2, X, Play, Square, Clock } from 'lucide-react';
import { format } from 'date-fns';

const schema = yup.object({
  project_id: yup.number().required('Project is required'),
  date: yup.string().required('Date is required'),
  hours: yup.number().positive('Hours must be positive').required('Hours is required'),
  description: yup.string(),
  billable: yup.boolean(),
});

type FormData = yup.InferType<typeof schema>;

export default function TimeEntries() {
  const [timeEntries, setTimeEntries] = useState<TimeEntry[]>([]);
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingEntry, setEditingEntry] = useState<TimeEntry | null>(null);
  const [activeTimer, setActiveTimer] = useState<TimeEntry | null>(null);
  const [elapsedTime, setElapsedTime] = useState(0);
  const [error, setError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<FormData>({
    resolver: yupResolver(schema),
    defaultValues: {
      date: format(new Date(), 'yyyy-MM-dd'),
      billable: true,
    },
  });

  const fetchData = async () => {
    try {
      const [entriesRes, projectsRes, timerRes] = await Promise.all([
        timeEntryAPI.list({ per_page: 100 }),
        projectAPI.list({ per_page: 100 }),
        timeEntryAPI.getActiveTimer(),
      ]);
      if (entriesRes.data.success && entriesRes.data.data) {
        setTimeEntries(entriesRes.data.data);
      }
      if (projectsRes.data.success && projectsRes.data.data) {
        setProjects(projectsRes.data.data);
      }
      if (timerRes.data.success && timerRes.data.data) {
        setActiveTimer(timerRes.data.data);
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load data');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  useEffect(() => {
    let interval: number | null = null;
    if (activeTimer && activeTimer.start_time) {
      interval = window.setInterval(() => {
        const start = new Date(activeTimer.start_time!).getTime();
        setElapsedTime(Math.floor((Date.now() - start) / 1000));
      }, 1000);
    }
    return () => {
      if (interval) clearInterval(interval);
    };
  }, [activeTimer]);

  const formatElapsedTime = (seconds: number) => {
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = seconds % 60;
    return `${h.toString().padStart(2, '0')}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
  };

  const onSubmit = async (data: FormData) => {
    setError(null);
    try {
      if (editingEntry) {
        await timeEntryAPI.update(editingEntry.id, data);
      } else {
        await timeEntryAPI.create(data);
      }
      setShowModal(false);
      setEditingEntry(null);
      reset();
      fetchData();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Operation failed');
    }
  };

  const handleStartTimer = async (projectId: number) => {
    try {
      const response = await timeEntryAPI.startTimer({ project_id: projectId });
      if (response.data.success && response.data.data) {
        setActiveTimer(response.data.data);
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to start timer');
    }
  };

  const handleStopTimer = async () => {
    if (!activeTimer) return;
    try {
      await timeEntryAPI.stopTimer(activeTimer.id);
      setActiveTimer(null);
      setElapsedTime(0);
      fetchData();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to stop timer');
    }
  };

  const handleEdit = (entry: TimeEntry) => {
    setEditingEntry(entry);
    reset({
      project_id: entry.project_id,
      date: entry.date.split('T')[0],
      hours: entry.hours,
      description: entry.description,
      billable: entry.billable,
    });
    setShowModal(true);
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this time entry?')) return;
    try {
      await timeEntryAPI.delete(id);
      fetchData();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to delete time entry');
    }
  };

  const openModal = () => {
    setEditingEntry(null);
    reset({
      date: format(new Date(), 'yyyy-MM-dd'),
      billable: true,
    });
    setShowModal(true);
  };

  const totalHoursToday = timeEntries
    .filter((e) => e.date.split('T')[0] === format(new Date(), 'yyyy-MM-dd'))
    .reduce((sum, e) => sum + e.hours, 0);

  if (loading) {
    return <div className="text-center py-12">Loading...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold text-gray-900">Time Entries</h1>
        <div className="flex space-x-3">
          {activeTimer ? (
            <div className="flex items-center px-4 py-2 bg-red-100 text-red-700 rounded-lg">
              <Clock className="w-5 h-5 mr-2 animate-pulse" />
              <span className="mr-3 font-mono text-lg">{formatElapsedTime(elapsedTime)}</span>
              <button
                onClick={handleStopTimer}
                className="flex items-center px-3 py-1 bg-red-600 text-white rounded hover:bg-red-700"
              >
                <Square className="w-4 h-4 mr-1" />
                Stop
              </button>
            </div>
          ) : (
            <select
              onChange={(e) => e.target.value && handleStartTimer(Number(e.target.value))}
              className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-indigo-500"
              defaultValue=""
            >
              <option value="" disabled>
                Start Timer...
              </option>
              {projects.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name}
                </option>
              ))}
            </select>
          )}
          <button
            onClick={openModal}
            className="flex items-center px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
          >
            <Plus className="w-5 h-5 mr-2" />
            Log Time
          </button>
        </div>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-400 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      )}

      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-lg font-semibold text-gray-900">Today's Summary</h2>
          <span className="text-2xl font-bold text-indigo-600">
            {totalHoursToday.toFixed(1)} / 24 hrs
          </span>
        </div>
        <div className="w-full bg-gray-200 rounded-full h-3">
          <div
            className="bg-indigo-600 h-3 rounded-full transition-all"
            style={{ width: `${Math.min((totalHoursToday / 24) * 100, 100)}%` }}
          />
        </div>
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Date
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Project
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Description
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Hours
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {timeEntries.map((entry) => (
              <tr key={entry.id}>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {format(new Date(entry.date), 'MMM dd, yyyy')}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {entry.project?.name || 'Unknown'}
                </td>
                <td className="px-6 py-4 text-sm text-gray-500 max-w-xs truncate">
                  {entry.description || '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  {entry.hours.toFixed(1)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span
                    className={`px-2 py-1 text-xs rounded-full ${
                      entry.billable
                        ? 'bg-green-100 text-green-800'
                        : 'bg-gray-100 text-gray-800'
                    }`}
                  >
                    {entry.billable ? 'Billable' : 'Non-billable'}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <button
                    onClick={() => handleEdit(entry)}
                    className="text-indigo-600 hover:text-indigo-900 mr-3"
                  >
                    <Edit2 className="w-4 h-4 inline" />
                  </button>
                  <button
                    onClick={() => handleDelete(entry.id)}
                    className="text-red-600 hover:text-red-900"
                  >
                    <Trash2 className="w-4 h-4 inline" />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {timeEntries.length === 0 && (
          <div className="text-center py-12 text-gray-500">
            No time entries yet. Start tracking your time!
          </div>
        )}
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-xl w-full max-w-md mx-4 p-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-bold text-gray-900">
                {editingEntry ? 'Edit Time Entry' : 'Log Time'}
              </h2>
              <button
                onClick={() => setShowModal(false)}
                className="text-gray-500 hover:text-gray-700"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Project *
                </label>
                <select
                  {...register('project_id')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                >
                  <option value="">Select a project</option>
                  {projects.map((p) => (
                    <option key={p.id} value={p.id}>
                      {p.name}
                    </option>
                  ))}
                </select>
                {errors.project_id && (
                  <p className="text-red-500 text-xs mt-1">{errors.project_id.message}</p>
                )}
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Date *
                </label>
                <input
                  type="date"
                  {...register('date')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                />
                {errors.date && (
                  <p className="text-red-500 text-xs mt-1">{errors.date.message}</p>
                )}
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Hours *
                </label>
                <input
                  type="number"
                  step="0.25"
                  min="0.25"
                  {...register('hours')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                />
                {errors.hours && (
                  <p className="text-red-500 text-xs mt-1">{errors.hours.message}</p>
                )}
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Description
                </label>
                <textarea
                  {...register('description')}
                  rows={2}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  placeholder="What did you work on?"
                />
              </div>
              <div className="flex items-center">
                <input
                  type="checkbox"
                  id="billable"
                  {...register('billable')}
                  className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                />
                <label htmlFor="billable" className="ml-2 text-sm text-gray-700">
                  Billable
                </label>
              </div>
              <div className="flex space-x-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="flex-1 px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="flex-1 px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700"
                >
                  {editingEntry ? 'Update' : 'Save'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
