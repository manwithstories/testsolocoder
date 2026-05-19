import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { projectAPI, clientAPI } from '../api';
import type { Project, Client, ProjectStatus, Milestone } from '../types';
import { Plus, Edit2, Trash2, X, Calendar, DollarSign, Flag, CheckCircle2 } from 'lucide-react';
import { format } from 'date-fns';

const statusColors: Record<ProjectStatus, string> = {
  draft: 'bg-gray-100 text-gray-800',
  active: 'bg-blue-100 text-blue-800',
  completed: 'bg-green-100 text-green-800',
  archived: 'bg-purple-100 text-purple-800',
};

const statusLabels: Record<ProjectStatus, string> = {
  draft: 'Draft',
  active: 'Active',
  completed: 'Completed',
  archived: 'Archived',
};

const validTransitions: Record<ProjectStatus, ProjectStatus[]> = {
  draft: ['active'],
  active: ['completed', 'draft'],
  completed: ['archived', 'active'],
  archived: ['completed'],
};

const projectSchema = yup.object({
  name: yup.string().required('Name is required'),
  client_id: yup.number().required('Client is required'),
  description: yup.string(),
  hourly_rate: yup.number().min(0, 'Rate must be positive'),
  deadline: yup.string(),
  budget: yup.number().min(0, 'Budget must be positive'),
});

const milestoneSchema = yup.object({
  title: yup.string().required('Title is required'),
  description: yup.string(),
  due_date: yup.string(),
});

type ProjectFormData = yup.InferType<typeof projectSchema>;
type MilestoneFormData = yup.InferType<typeof milestoneSchema>;

interface TempMilestone {
  id?: number;
  title: string;
  description: string;
  due_date: string;
}

export default function Projects() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [clients, setClients] = useState<Client[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingProject, setEditingProject] = useState<Project | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState<ProjectStatus | 'all'>('all');

  const [showMilestoneModal, setShowMilestoneModal] = useState(false);
  const [selectedProject, setSelectedProject] = useState<Project | null>(null);
  const [editingMilestone, setEditingMilestone] = useState<Milestone | null>(null);
  const [tempMilestones, setTempMilestones] = useState<TempMilestone[]>([]);

  const {
    register: registerProject,
    handleSubmit: handleProjectSubmit,
    reset: resetProject,
    formState: { errors: projectErrors },
  } = useForm<ProjectFormData>({
    resolver: yupResolver(projectSchema),
  });

  const {
    register: registerMilestone,
    handleSubmit: handleMilestoneSubmit,
    reset: resetMilestone,
    formState: { errors: milestoneErrors },
  } = useForm<MilestoneFormData>({
    resolver: yupResolver(milestoneSchema),
  });

  const fetchData = async () => {
    try {
      const [projectsRes, clientsRes] = await Promise.all([
        projectAPI.list({ per_page: 100 }),
        clientAPI.list({ per_page: 100 }),
      ]);
      if (projectsRes.data.success && projectsRes.data.data) {
        setProjects(projectsRes.data.data);
      }
      if (clientsRes.data.success && clientsRes.data.data) {
        setClients(clientsRes.data.data);
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

  const onSubmitProject = async (data: ProjectFormData) => {
    setError(null);
    try {
      const payload: any = {
        ...data,
        milestones: tempMilestones.map(m => ({
          title: m.title,
          description: m.description,
          due_date: m.due_date ? new Date(m.due_date).toISOString() : null,
        })),
      };
      if (data.deadline) {
        payload.deadline = new Date(data.deadline).toISOString();
      }
      if (editingProject) {
        await projectAPI.update(editingProject.id, payload);
      } else {
        await projectAPI.create(payload);
      }
      setShowModal(false);
      setEditingProject(null);
      setTempMilestones([]);
      resetProject();
      fetchData();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Operation failed');
    }
  };

  const handleStatusChange = async (projectId: number, newStatus: ProjectStatus) => {
    try {
      await projectAPI.update(projectId, { status: newStatus });
      fetchData();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to update status');
    }
  };

  const handleEdit = (project: Project) => {
    setEditingProject(project);
    setTempMilestones(project.milestones?.map(m => ({
      id: m.id,
      title: m.title,
      description: m.description,
      due_date: m.due_date ? m.due_date.split('T')[0] : '',
    })) || []);
    resetProject({
      name: project.name,
      client_id: project.client_id,
      description: project.description,
      hourly_rate: project.hourly_rate,
      deadline: project.deadline ? project.deadline.split('T')[0] : '',
      budget: project.budget,
    });
    setShowModal(true);
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this project?')) return;
    try {
      await projectAPI.delete(id);
      fetchData();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to delete project');
    }
  };

  const openModal = () => {
    setEditingProject(null);
    setTempMilestones([]);
    resetProject();
    setShowModal(true);
  };

  const openMilestoneModal = (project: Project) => {
    setSelectedProject(project);
    setEditingMilestone(null);
    resetMilestone();
    setShowMilestoneModal(true);
  };

  const handleAddMilestone = (data: MilestoneFormData) => {
    if (editingMilestone) {
      setTempMilestones(prev => prev.map(m =>
        m.id === editingMilestone.id
          ? { ...m, title: data.title, description: data.description || '', due_date: data.due_date || '' }
          : m
      ));
    } else {
      setTempMilestones(prev => [...prev, {
        title: data.title,
        description: data.description || '',
        due_date: data.due_date || '',
      }]);
    }
    resetMilestone();
    setEditingMilestone(null);
  };

  const handleEditMilestone = (milestone: Milestone) => {
    setEditingMilestone(milestone);
    resetMilestone({
      title: milestone.title,
      description: milestone.description,
      due_date: milestone.due_date ? milestone.due_date.split('T')[0] : '',
    });
  };

  const handleRemoveMilestone = (index: number) => {
    setTempMilestones(prev => prev.filter((_, i) => i !== index));
  };

  const handleSaveMilestones = async () => {
    if (!selectedProject) return;
    setError(null);
    try {
      const existingMilestoneIds = selectedProject.milestones?.map(m => m.id) || [];
      const newMilestones = tempMilestones.filter(m => !existingMilestoneIds.includes(m.id!));
      const updatedMilestones = tempMilestones.filter(m => existingMilestoneIds.includes(m.id!));

      for (const m of newMilestones) {
        await projectAPI.addMilestone(selectedProject.id, {
          title: m.title,
          description: m.description,
          due_date: m.due_date ? new Date(m.due_date).toISOString() : null,
        });
      }

      for (const m of updatedMilestones) {
        if (m.id) {
          await projectAPI.updateMilestone(m.id, {
            title: m.title,
            description: m.description,
            due_date: m.due_date ? new Date(m.due_date).toISOString() : null,
          });
        }
      }

      setShowMilestoneModal(false);
      setSelectedProject(null);
      setTempMilestones([]);
      setEditingMilestone(null);
      fetchData();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to save milestones');
    }
  };

  const handleDeleteMilestone = async (milestoneId: number) => {
    if (!confirm('Are you sure you want to delete this milestone?')) return;
    try {
      await projectAPI.deleteMilestone(milestoneId);
      if (selectedProject) {
        const updated = { ...selectedProject };
        updated.milestones = updated.milestones?.filter(m => m.id !== milestoneId);
        setSelectedProject(updated);
      }
      fetchData();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to delete milestone');
    }
  };

  const handleToggleMilestoneComplete = async (milestone: Milestone) => {
    try {
      await projectAPI.updateMilestone(milestone.id, {
        completed: !milestone.completed,
      });
      fetchData();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to update milestone');
    }
  };

  const filteredProjects = filter === 'all'
    ? projects
    : projects.filter((p) => p.status === filter);

  if (loading) {
    return <div className="text-center py-12">Loading...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold text-gray-900">Projects</h1>
        <button
          onClick={openModal}
          className="flex items-center px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
        >
          <Plus className="w-5 h-5 mr-2" />
          Add Project
        </button>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-400 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      )}

      <div className="flex space-x-2">
        {(['all', 'draft', 'active', 'completed', 'archived'] as const).map((s) => (
          <button
            key={s}
            onClick={() => setFilter(s)}
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
              filter === s
                ? 'bg-indigo-600 text-white'
                : 'bg-white text-gray-600 hover:bg-gray-100'
            }`}
          >
            {s === 'all' ? 'All' : statusLabels[s]}
          </button>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {filteredProjects.map((project) => (
          <div key={project.id} className="bg-white rounded-lg shadow p-6">
            <div className="flex justify-between items-start mb-4">
              <div>
                <h3 className="text-xl font-semibold text-gray-900">{project.name}</h3>
                <p className="text-sm text-gray-500">
                  Client: {project.client?.name || 'Unknown'}
                </p>
              </div>
              <div className="flex items-center space-x-2">
                <span
                  className={`px-3 py-1 rounded-full text-xs font-medium ${statusColors[project.status]}`}
                >
                  {statusLabels[project.status]}
                </span>
                <button
                  onClick={() => openMilestoneModal(project)}
                  className="p-2 text-gray-500 hover:text-indigo-600 hover:bg-gray-100 rounded"
                  title="Manage milestones"
                >
                  <Flag className="w-4 h-4" />
                </button>
                <button
                  onClick={() => handleEdit(project)}
                  className="p-2 text-gray-500 hover:text-indigo-600 hover:bg-gray-100 rounded"
                >
                  <Edit2 className="w-4 h-4" />
                </button>
                <button
                  onClick={() => handleDelete(project.id)}
                  className="p-2 text-gray-500 hover:text-red-600 hover:bg-gray-100 rounded"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>

            {project.description && (
              <p className="text-gray-600 text-sm mb-4">{project.description}</p>
            )}

            <div className="grid grid-cols-2 gap-4 mb-4 text-sm">
              <div className="flex items-center text-gray-600">
                <DollarSign className="w-4 h-4 mr-2" />
                ${project.hourly_rate.toFixed(2)}/hr
              </div>
              {project.budget > 0 && (
                <div className="text-gray-600">
                  Budget: ${project.budget.toFixed(2)}
                </div>
              )}
              {project.deadline && (
                <div className="flex items-center text-gray-600 col-span-2">
                  <Calendar className="w-4 h-4 mr-2" />
                  Deadline: {format(new Date(project.deadline), 'MMM dd, yyyy')}
                </div>
              )}
            </div>

            <div className="border-t pt-4">
              <p className="text-xs text-gray-500 mb-2">Change Status:</p>
              <div className="flex flex-wrap gap-2">
                {validTransitions[project.status].map((nextStatus) => (
                  <button
                    key={nextStatus}
                    onClick={() => handleStatusChange(project.id, nextStatus)}
                    className={`px-3 py-1 text-xs rounded-full ${statusColors[nextStatus]} hover:opacity-80`}
                  >
                    → {statusLabels[nextStatus]}
                  </button>
                ))}
              </div>
            </div>

            {project.milestones && project.milestones.length > 0 && (
              <div className="border-t pt-4 mt-4">
                <p className="text-sm font-medium text-gray-700 mb-2">
                  Milestones ({project.milestones.filter(m => m.completed).length}/{project.milestones.length})
                </p>
                <div className="space-y-2">
                  {project.milestones.map((m) => (
                    <div key={m.id} className="flex items-center text-sm">
                      <button
                        onClick={() => handleToggleMilestoneComplete(m)}
                        className="mr-2"
                      >
                        {m.completed ? (
                          <CheckCircle2 className="w-5 h-5 text-green-500" />
                        ) : (
                          <div className="w-5 h-5 border-2 border-gray-300 rounded-full" />
                        )}
                      </button>
                      <span className={m.completed ? 'line-through text-gray-400' : 'text-gray-700'}>
                        {m.title}
                      </span>
                      {m.due_date && (
                        <span className="ml-auto text-xs text-gray-500">
                          {format(new Date(m.due_date), 'MMM dd')}
                        </span>
                      )}
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        ))}
      </div>

      {filteredProjects.length === 0 && (
        <div className="text-center py-12 bg-white rounded-lg shadow">
          <p className="text-gray-500">No projects found.</p>
        </div>
      )}

      {/* Project Create/Edit Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-xl w-full max-w-2xl mx-4 p-6 max-h-[90vh] overflow-y-auto">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-bold text-gray-900">
                {editingProject ? 'Edit Project' : 'Add Project'}
              </h2>
              <button
                onClick={() => { setShowModal(false); setTempMilestones([]); }}
                className="text-gray-500 hover:text-gray-700"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            <form onSubmit={handleProjectSubmit(onSubmitProject)} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Name *
                </label>
                <input
                  type="text"
                  {...registerProject('name')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                />
                {projectErrors.name && (
                  <p className="text-red-500 text-xs mt-1">{projectErrors.name.message}</p>
                )}
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Client *
                </label>
                <select
                  {...registerProject('client_id')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                >
                  <option value="">Select a client</option>
                  {clients.map((client) => (
                    <option key={client.id} value={client.id}>
                      {client.name}
                    </option>
                  ))}
                </select>
                {projectErrors.client_id && (
                  <p className="text-red-500 text-xs mt-1">{projectErrors.client_id.message}</p>
                )}
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Description
                </label>
                <textarea
                  {...registerProject('description')}
                  rows={2}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Hourly Rate ($)
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    {...registerProject('hourly_rate')}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Budget ($)
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    {...registerProject('budget')}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  />
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Deadline
                </label>
                <input
                  type="date"
                  {...registerProject('deadline')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                />
              </div>

              {/* Milestones section in project form */}
              <div className="border-t pt-4">
                <div className="flex justify-between items-center mb-3">
                  <h3 className="text-sm font-medium text-gray-700">Milestones</h3>
                </div>

                <div className="space-y-2 mb-3">
                  {tempMilestones.map((m, index) => (
                    <div key={index} className="flex items-center p-2 bg-gray-50 rounded">
                      <div className="flex-1">
                        <p className="text-sm font-medium text-gray-900">{m.title}</p>
                        {m.description && (
                          <p className="text-xs text-gray-500">{m.description}</p>
                        )}
                      </div>
                      {m.due_date && (
                        <span className="text-xs text-gray-500 mr-3">
                          {m.due_date}
                        </span>
                      )}
                      <button
                        type="button"
                        onClick={() => handleRemoveMilestone(index)}
                        className="text-red-500 hover:text-red-700"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  ))}
                </div>

                <div className="flex gap-2">
                  <input
                    type="text"
                    placeholder="Milestone title..."
                    {...registerMilestone('title')}
                    className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  />
                  <input
                    type="date"
                    {...registerMilestone('due_date')}
                    className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  />
                  <button
                    type="button"
                    onClick={handleMilestoneSubmit(handleAddMilestone)}
                    className="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300"
                  >
                    <Plus className="w-4 h-4" />
                  </button>
                </div>
                {milestoneErrors.title && (
                  <p className="text-red-500 text-xs mt-1">{milestoneErrors.title.message}</p>
                )}
              </div>

              <div className="flex space-x-3 pt-4">
                <button
                  type="button"
                  onClick={() => { setShowModal(false); setTempMilestones([]); }}
                  className="flex-1 px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="flex-1 px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700"
                >
                  {editingProject ? 'Update' : 'Create'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Milestone Management Modal */}
      {showMilestoneModal && selectedProject && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-xl w-full max-w-lg mx-4 p-6">
            <div className="flex justify-between items-center mb-4">
              <div>
                <h2 className="text-xl font-bold text-gray-900">Manage Milestones</h2>
                <p className="text-sm text-gray-500">{selectedProject.name}</p>
              </div>
              <button
                onClick={() => {
                  setShowMilestoneModal(false);
                  setSelectedProject(null);
                  setTempMilestones([]);
                  setEditingMilestone(null);
                }}
                className="text-gray-500 hover:text-gray-700"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-3 mb-4 max-h-64 overflow-y-auto">
              {selectedProject.milestones?.map((m) => (
                <div key={m.id} className="flex items-center p-3 bg-gray-50 rounded-lg">
                  <input
                    type="checkbox"
                    checked={m.completed}
                    onChange={() => handleToggleMilestoneComplete(m)}
                    className="h-5 w-5 text-indigo-600 border-gray-300 rounded mr-3"
                  />
                  <div className="flex-1">
                    <p className={`font-medium ${m.completed ? 'line-through text-gray-400' : 'text-gray-900'}`}>
                      {m.title}
                    </p>
                    {m.description && (
                      <p className="text-xs text-gray-500">{m.description}</p>
                    )}
                    {m.due_date && (
                      <p className="text-xs text-gray-400 mt-1">
                        Due: {format(new Date(m.due_date), 'MMM dd, yyyy')}
                      </p>
                    )}
                  </div>
                  <div className="flex space-x-2">
                    <button
                      onClick={() => handleEditMilestone(m)}
                      className="p-1 text-gray-500 hover:text-indigo-600"
                    >
                      <Edit2 className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => handleDeleteMilestone(m.id)}
                      className="p-1 text-gray-500 hover:text-red-600"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              ))}
              {(!selectedProject.milestones || selectedProject.milestones.length === 0) && (
                <p className="text-center text-gray-500 py-4">No milestones yet</p>
              )}
            </div>

            <div className="border-t pt-4">
              <h3 className="text-sm font-medium text-gray-700 mb-3">
                {editingMilestone ? 'Edit Milestone' : 'Add New Milestone'}
              </h3>
              <form
                onSubmit={handleMilestoneSubmit(async (data) => {
                  if (editingMilestone) {
                    await projectAPI.updateMilestone(editingMilestone.id, {
                      title: data.title,
                      description: data.description || '',
                      due_date: data.due_date ? new Date(data.due_date).toISOString() : null,
                    });
                    fetchData();
                    setEditingMilestone(null);
                    resetMilestone();
                  } else {
                    await projectAPI.addMilestone(selectedProject.id, {
                      title: data.title,
                      description: data.description || '',
                      due_date: data.due_date ? new Date(data.due_date).toISOString() : null,
                    });
                    fetchData();
                    resetMilestone();
                  }
                })}
                className="space-y-3"
              >
                <div>
                  <input
                    type="text"
                    placeholder="Milestone title"
                    {...registerMilestone('title')}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  />
                  {milestoneErrors.title && (
                    <p className="text-red-500 text-xs mt-1">{milestoneErrors.title.message}</p>
                  )}
                </div>
                <div>
                  <input
                    type="text"
                    placeholder="Description (optional)"
                    {...registerMilestone('description')}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  />
                </div>
                <div>
                  <label className="block text-sm text-gray-600 mb-1">Due Date</label>
                  <input
                    type="date"
                    {...registerMilestone('due_date')}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  />
                </div>
                <div className="flex space-x-2">
                  {editingMilestone && (
                    <button
                      type="button"
                      onClick={() => { setEditingMilestone(null); resetMilestone(); }}
                      className="flex-1 px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
                    >
                      Cancel Edit
                    </button>
                  )}
                  <button
                    type="submit"
                    className={`${editingMilestone ? 'flex-1' : 'w-full'} px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700`}
                  >
                    {editingMilestone ? 'Update' : 'Add Milestone'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
