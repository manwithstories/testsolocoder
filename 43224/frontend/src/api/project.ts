import request from '@/utils/request'

export interface CreateProjectData {
  title: string
  description?: string
  source_lang: string
  target_lang: string
  expertise_tag_ids?: number[]
  word_count: number
  urgency: string
  deadline: string
}

export function listProjects(params?: any) {
  return request.get('/projects', { params })
}

export function getProject(id: number) {
  return request.get(`/projects/${id}`)
}

export function createProject(data: CreateProjectData) {
  return request.post('/projects', data)
}

export function approveProject(id: number) {
  return request.put(`/projects/${id}/approve`)
}

export function assignTranslator(id: number, translatorId: number) {
  return request.put(`/projects/${id}/assign`, { translator_id: translatorId })
}

export function startProject(id: number) {
  return request.put(`/projects/${id}/start`)
}

export function submitForReview(id: number) {
  return request.put(`/projects/${id}/submit`)
}

export function completeProject(id: number) {
  return request.put(`/projects/${id}/complete`)
}

export function cancelProject(id: number) {
  return request.put(`/projects/${id}/cancel`)
}

export function addProjectComment(id: number, content: string) {
  return request.post(`/projects/${id}/comments`, { content })
}

export function recommendTranslators(id: number) {
  return request.get(`/projects/${id}/recommend-translators`)
}

export function autoAssignTranslator(id: number) {
  return request.post(`/projects/${id}/auto-assign`)
}

export function listLanguagePairs() {
  return request.get('/language-pairs')
}

export function createLanguagePair(data: any) {
  return request.post('/language-pairs', data)
}

export function listExpertiseTags() {
  return request.get('/expertise-tags')
}

export function createExpertiseTag(data: any) {
  return request.post('/expertise-tags', data)
}
