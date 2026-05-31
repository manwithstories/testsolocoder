import request from '@/utils/request'

export function listDocuments(projectId: number, params?: any) {
  return request.get(`/projects/${projectId}/documents`, { params })
}

export function uploadDocument(projectId: number, formData: FormData) {
  return request.post(`/projects/${projectId}/documents`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

export function downloadDocument(id: number) {
  return request.get(`/documents/${id}`, { responseType: 'blob' })
}

export function deleteDocument(id: number) {
  return request.delete(`/documents/${id}`)
}

export function getDocumentVersions(projectId: number, isSource: boolean = true) {
  return request.get(`/projects/${projectId}/documents/versions`, { params: { is_source: isSource } })
}

export function extractSegments(id: number, maxLen?: number) {
  return request.post(`/documents/${id}/extract-segments`, { max_len: maxLen })
}

export function listSegments(projectId: number, params?: any) {
  return request.get(`/projects/${projectId}/segments`, { params })
}

export function updateSegmentTranslation(id: number, translatedText: string) {
  return request.put(`/segments/${id}`, { translated_text: translatedText })
}

export function getMemorySuggestions(text: string, sourceLang: string, targetLang: string) {
  return request.get('/memory/suggestions', {
    params: { text, source_lang: sourceLang, target_lang: targetLang }
  })
}
