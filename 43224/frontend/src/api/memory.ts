import request from '@/utils/request'

export function listTranslationMemories(params?: any) {
  return request.get('/translation-memories', { params })
}

export function createTranslationMemory(data: any) {
  return request.post('/translation-memories', data)
}

export function updateTranslationMemory(id: number, data: any) {
  return request.put(`/translation-memories/${id}`, data)
}

export function deleteTranslationMemory(id: number) {
  return request.delete(`/translation-memories/${id}`)
}

export function listGlossaryTerms(params?: any) {
  return request.get('/glossary-terms', { params })
}

export function createGlossaryTerm(data: any) {
  return request.post('/glossary-terms', data)
}

export function updateGlossaryTerm(id: number, data: any) {
  return request.put(`/glossary-terms/${id}`, data)
}

export function deleteGlossaryTerm(id: number) {
  return request.delete(`/glossary-terms/${id}`)
}
