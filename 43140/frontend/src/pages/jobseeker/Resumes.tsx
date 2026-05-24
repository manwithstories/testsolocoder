import { useState, useEffect } from 'react'
import { Plus, Edit2, Trash2, FileText, Upload, Star, Check } from 'lucide-react'
import { resumeApi, Resume, Education, WorkExperience, Skill } from '@/api/jobs'
import dayjs from 'dayjs'

interface ResumeFormData {
  title: string
  full_name: string
  email: string
  phone: string
  location: string
  summary: string
  is_default: boolean
  education_list: { school: string; degree: string; major: string; start_date: string; end_date: string; description: string }[]
  work_experiences: { company: string; position: string; start_date: string; end_date: string; description: string }[]
  skills: { name: string; level: string }[]
}

export default function JobSeekerResumes() {
  const [resumes, setResumes] = useState<Resume[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [editingResume, setEditingResume] = useState<Resume | null>(null)

  useEffect(() => {
    loadResumes()
  }, [])

  const loadResumes = async () => {
    try {
      const response = await resumeApi.list()
      setResumes(response.data)
    } catch (err) {
      console.error('Failed to load resumes:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (id: number) => {
    if (confirm('Are you sure you want to delete this resume?')) {
      try {
        await resumeApi.delete(id)
        loadResumes()
      } catch (err) {
        alert('Failed to delete resume')
      }
    }
  }

  const handleSetDefault = async (id: number) => {
    try {
      await resumeApi.setDefault(id)
      loadResumes()
    } catch (err) {
      alert('Failed to update resume')
    }
  }

  const handleFileUpload = async (resumeId: number, file: File) => {
    try {
      const formData = new FormData()
      formData.append('file', file)
      await resumeApi.uploadFile(resumeId, formData)
      loadResumes()
    } catch (err) {
      alert('Failed to upload file')
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">My Resumes</h1>
        <button onClick={() => setShowModal(true)} className="btn-primary flex items-center gap-2">
          <Plus className="w-5 h-5" />
          Add Resume
        </button>
      </div>

      {loading ? (
        <div className="text-center py-12">
          <div className="animate-spin w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full mx-auto"></div>
        </div>
      ) : resumes.length === 0 ? (
        <div className="text-center py-12">
          <FileText className="w-16 h-16 text-gray-300 mx-auto" />
          <p className="mt-4 text-gray-500">No resumes yet. Create your first resume!</p>
        </div>
      ) : (
        <div className="grid gap-4">
          {resumes.map((resume) => (
            <div key={resume.id} className="card">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-2">
                    <h3 className="text-lg font-semibold">{resume.title}</h3>
                    {resume.is_default && (
                      <span className="badge bg-green-100 text-green-800 flex items-center gap-1">
                        <Star className="w-3 h-3 fill-current" />
                        Default
                      </span>
                    )}
                  </div>
                  <p className="text-gray-600">{resume.full_name}</p>
                  {resume.email && <p className="text-sm text-gray-500">{resume.email}</p>}
                  <div className="flex flex-wrap gap-2 mt-2">
                    {resume.skills?.slice(0, 5).map((skill, index) => (
                      <span
                        key={index}
                        className="px-2 py-1 bg-gray-100 text-gray-600 rounded-full text-xs"
                      >
                        {skill.name}
                      </span>
                    ))}
                  </div>
                  {resume.file_name && (
                    <p className="text-sm text-primary-600 mt-2 flex items-center gap-1">
                      <FileText className="w-4 h-4" />
                      {resume.file_name}
                    </p>
                  )}
                  <p className="text-xs text-gray-400 mt-2">
                    Updated {dayjs(resume.updated_at).fromNow()}
                  </p>
                </div>
                <div className="flex items-center gap-2">
                  <label className="btn-secondary cursor-pointer flex items-center gap-2">
                    <Upload className="w-4 h-4" />
                    Upload
                    <input
                      type="file"
                      accept=".pdf"
                      className="hidden"
                      onChange={(e) => {
                        const file = e.target.files?.[0]
                        if (file) handleFileUpload(resume.id, file)
                      }}
                    />
                  </label>
                  {!resume.is_default && (
                    <button
                      onClick={() => handleSetDefault(resume.id)}
                      className="p-2 text-gray-600 hover:text-yellow-600"
                      title="Set as default"
                    >
                      <Star className="w-5 h-5" />
                    </button>
                  )}
                  <button
                    onClick={() => {
                      setEditingResume(resume)
                      setShowModal(true)
                    }}
                    className="p-2 text-gray-600 hover:text-primary-600"
                  >
                    <Edit2 className="w-5 h-5" />
                  </button>
                  <button
                    onClick={() => handleDelete(resume.id)}
                    className="p-2 text-gray-600 hover:text-red-600"
                  >
                    <Trash2 className="w-5 h-5" />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {showModal && (
        <ResumeModal
          resume={editingResume}
          onClose={() => {
            setShowModal(false)
            setEditingResume(null)
          }}
          onSave={() => {
            setShowModal(false)
            setEditingResume(null)
            loadResumes()
          }}
        />
      )}
    </div>
  )
}

function ResumeModal({
  resume,
  onClose,
  onSave,
}: {
  resume: Resume | null
  onClose: () => void
  onSave: () => void
}) {
  const [formData, setFormData] = useState<ResumeFormData>({
    title: resume?.title || '',
    full_name: resume?.full_name || '',
    email: resume?.email || '',
    phone: resume?.phone || '',
    location: resume?.location || '',
    summary: resume?.summary || '',
    is_default: resume?.is_default || false,
    education_list: resume?.education_list?.map((e) => ({
      school: e.school,
      degree: e.degree || '',
      major: e.major || '',
      start_date: e.start_date || '',
      end_date: e.end_date || '',
      description: e.description || '',
    })) || [],
    work_experiences: resume?.work_experiences?.map((w) => ({
      company: w.company,
      position: w.position,
      start_date: w.start_date || '',
      end_date: w.end_date || '',
      description: w.description || '',
    })) || [],
    skills: resume?.skills?.map((s) => ({ name: s.name, level: s.level || '' })) || [],
  })
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!formData.title || !formData.full_name) {
      alert('Please fill in required fields')
      return
    }
    setSubmitting(true)
    try {
      if (resume) {
        await resumeApi.update(resume.id, formData as any)
      } else {
        await resumeApi.create(formData as any)
      }
      onSave()
    } catch (err: any) {
      alert(err.message || 'Failed to save resume')
    } finally {
      setSubmitting(false)
    }
  }

  const addEducation = () => {
    setFormData({
      ...formData,
      education_list: [
        ...formData.education_list,
        { school: '', degree: '', major: '', start_date: '', end_date: '', description: '' },
      ],
    })
  }

  const removeEducation = (index: number) => {
    setFormData({
      ...formData,
      education_list: formData.education_list.filter((_, i) => i !== index),
    })
  }

  const addWorkExperience = () => {
    setFormData({
      ...formData,
      work_experiences: [
        ...formData.work_experiences,
        { company: '', position: '', start_date: '', end_date: '', description: '' },
      ],
    })
  }

  const removeWorkExperience = (index: number) => {
    setFormData({
      ...formData,
      work_experiences: formData.work_experiences.filter((_, i) => i !== index),
    })
  }

  const addSkill = () => {
    setFormData({
      ...formData,
      skills: [...formData.skills, { name: '', level: '' }],
    })
  }

  const removeSkill = (index: number) => {
    setFormData({
      ...formData,
      skills: formData.skills.filter((_, i) => i !== index),
    })
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div className="p-6">
          <h2 className="text-2xl font-bold mb-6">{resume ? 'Edit Resume' : 'Create Resume'}</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Title *</label>
                <input
                  type="text"
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                  className="input-field"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Full Name *</label>
                <input
                  type="text"
                  value={formData.full_name}
                  onChange={(e) => setFormData({ ...formData, full_name: e.target.value })}
                  className="input-field"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
                <input
                  type="email"
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  className="input-field"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Phone</label>
                <input
                  type="tel"
                  value={formData.phone}
                  onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                  className="input-field"
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Location</label>
              <input
                type="text"
                value={formData.location}
                onChange={(e) => setFormData({ ...formData, location: e.target.value })}
                className="input-field"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Summary</label>
              <textarea
                value={formData.summary}
                onChange={(e) => setFormData({ ...formData, summary: e.target.value })}
                rows={3}
                className="input-field"
              />
            </div>

            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="is_default"
                checked={formData.is_default}
                onChange={(e) => setFormData({ ...formData, is_default: e.target.checked })}
                className="w-4 h-4"
              />
              <label htmlFor="is_default" className="text-sm">
                Set as default resume
              </label>
            </div>

            <div className="border-t pt-4">
              <div className="flex items-center justify-between mb-2">
                <h3 className="font-semibold">Education</h3>
                <button type="button" onClick={addEducation} className="text-primary-600 text-sm">
                  + Add Education
                </button>
              </div>
              {formData.education_list.map((edu, index) => (
                <div key={index} className="bg-gray-50 rounded-lg p-4 mb-2">
                  <div className="flex justify-end mb-2">
                    <button type="button" onClick={() => removeEducation(index)} className="text-red-500 text-sm">
                      Remove
                    </button>
                  </div>
                  <div className="grid grid-cols-2 gap-2">
                    <input
                      type="text"
                      placeholder="School *"
                      value={edu.school}
                      onChange={(e) => {
                        const newEdu = [...formData.education_list]
                        newEdu[index].school = e.target.value
                        setFormData({ ...formData, education_list: newEdu })
                      }}
                      className="input-field"
                      required
                    />
                    <input
                      type="text"
                      placeholder="Degree"
                      value={edu.degree}
                      onChange={(e) => {
                        const newEdu = [...formData.education_list]
                        newEdu[index].degree = e.target.value
                        setFormData({ ...formData, education_list: newEdu })
                      }}
                      className="input-field"
                    />
                    <input
                      type="text"
                      placeholder="Major"
                      value={edu.major}
                      onChange={(e) => {
                        const newEdu = [...formData.education_list]
                        newEdu[index].major = e.target.value
                        setFormData({ ...formData, education_list: newEdu })
                      }}
                      className="input-field"
                    />
                    <div className="grid grid-cols-2 gap-2">
                      <input
                        type="text"
                        placeholder="Start Date"
                        value={edu.start_date}
                        onChange={(e) => {
                          const newEdu = [...formData.education_list]
                          newEdu[index].start_date = e.target.value
                          setFormData({ ...formData, education_list: newEdu })
                        }}
                        className="input-field"
                      />
                      <input
                        type="text"
                        placeholder="End Date"
                        value={edu.end_date}
                        onChange={(e) => {
                          const newEdu = [...formData.education_list]
                          newEdu[index].end_date = e.target.value
                          setFormData({ ...formData, education_list: newEdu })
                        }}
                        className="input-field"
                      />
                    </div>
                  </div>
                </div>
              ))}
            </div>

            <div className="border-t pt-4">
              <div className="flex items-center justify-between mb-2">
                <h3 className="font-semibold">Work Experience</h3>
                <button type="button" onClick={addWorkExperience} className="text-primary-600 text-sm">
                  + Add Experience
                </button>
              </div>
              {formData.work_experiences.map((work, index) => (
                <div key={index} className="bg-gray-50 rounded-lg p-4 mb-2">
                  <div className="flex justify-end mb-2">
                    <button type="button" onClick={() => removeWorkExperience(index)} className="text-red-500 text-sm">
                      Remove
                    </button>
                  </div>
                  <div className="grid grid-cols-2 gap-2">
                    <input
                      type="text"
                      placeholder="Company *"
                      value={work.company}
                      onChange={(e) => {
                        const newWork = [...formData.work_experiences]
                        newWork[index].company = e.target.value
                        setFormData({ ...formData, work_experiences: newWork })
                      }}
                      className="input-field"
                      required
                    />
                    <input
                      type="text"
                      placeholder="Position *"
                      value={work.position}
                      onChange={(e) => {
                        const newWork = [...formData.work_experiences]
                        newWork[index].position = e.target.value
                        setFormData({ ...formData, work_experiences: newWork })
                      }}
                      className="input-field"
                      required
                    />
                    <input
                      type="text"
                      placeholder="Start Date"
                      value={work.start_date}
                      onChange={(e) => {
                        const newWork = [...formData.work_experiences]
                        newWork[index].start_date = e.target.value
                        setFormData({ ...formData, work_experiences: newWork })
                      }}
                      className="input-field"
                    />
                    <input
                      type="text"
                      placeholder="End Date"
                      value={work.end_date}
                      onChange={(e) => {
                        const newWork = [...formData.work_experiences]
                        newWork[index].end_date = e.target.value
                        setFormData({ ...formData, work_experiences: newWork })
                      }}
                      className="input-field"
                    />
                  </div>
                </div>
              ))}
            </div>

            <div className="border-t pt-4">
              <div className="flex items-center justify-between mb-2">
                <h3 className="font-semibold">Skills</h3>
                <button type="button" onClick={addSkill} className="text-primary-600 text-sm">
                  + Add Skill
                </button>
              </div>
              <div className="flex flex-wrap gap-2">
                {formData.skills.map((skill, index) => (
                  <div key={index} className="flex items-center gap-1">
                    <input
                      type="text"
                      placeholder="Skill name"
                      value={skill.name}
                      onChange={(e) => {
                        const newSkills = [...formData.skills]
                        newSkills[index].name = e.target.value
                        setFormData({ ...formData, skills: newSkills })
                      }}
                      className="input-field w-32"
                    />
                    <button type="button" onClick={() => removeSkill(index)} className="text-red-500">
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                ))}
              </div>
            </div>

            <div className="flex gap-3 pt-4">
              <button type="button" onClick={onClose} className="btn-secondary flex-1">
                Cancel
              </button>
              <button type="submit" disabled={submitting} className="btn-primary flex-1">
                {submitting ? 'Saving...' : resume ? 'Update' : 'Create'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}
