import request, { PageData, PageParams } from '@/utils/request'

export interface Work {
  id: number
  user_id: number
  artist_id: number
  title: string
  artist_name: string
  album_id?: number
  type: string
  genre: string
  sub_genre?: string
  language?: string
  duration: number
  description?: string
  lyrics?: string
  composer?: string
  lyricist?: string
  arranger?: string
  producer?: string
  cover_url: string
  audio_url: string
  audio_format: string
  status: number
  published_at?: string
  play_count: number
  like_count: number
  comment_count: number
  share_count: number
  is_public: boolean
  explicit: boolean
  tags?: Tag[]
  copyright?: Copyright
  user?: {
    id: number
    username: string
    nickname: string
    avatar: string
  }
  created_at: string
  updated_at: string
}

export interface Tag {
  id: number
  name: string
  type: string
  usage_count: number
}

export interface Copyright {
  id: number
  work_id: number
  copyright_type: string
  owner: string
  license_type: string
  royalties_rate: number
  start_date: string
  end_date?: string
  is_exclusive: boolean
  territory: string
  contract_no: string
  notes?: string
}

export interface Album {
  id: number
  user_id: number
  artist_id: number
  title: string
  artist_name: string
  description?: string
  cover_url: string
  release_date: string
  type: string
  genre: string
  record_label?: string
  upc?: string
  status: number
  published_at?: string
  play_count: number
  like_count: number
  comment_count: number
  work_count: number
  works?: Work[]
  created_at: string
  updated_at: string
}

export interface UploadWorkParams {
  title: string
  artist_name: string
  genre?: string
  description?: string
  album_id?: number
  type?: string
  explicit?: boolean
  tags?: string[]
}

export const workApi = {
  list: (params: PageParams) => request.get<PageData<Work>>('/works', params),
  search: (params: PageParams & { keyword?: string; status?: number; artist_id?: number }) => 
    request.get<PageData<Work>>('/works', params),
  getById: (id: number) => request.get<Work>(`/works/${id}`),
  getByArtist: (artistId: number, params: PageParams & { status?: number }) => 
    request.get<PageData<Work>>(`/works/artist/${artistId}`, params),
  upload: (formData: FormData) => request.post('/works/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  }),
  update: (id: number, data: Partial<Work>) => request.put(`/works/${id}`, data),
  delete: (id: number) => request.delete(`/works/${id}`),
  batchPublish: (workIds: number[]) => request.post('/works/batch-publish', { work_ids: workIds }),
  recordPlay: (id: number, duration: number) => request.post(`/works/${id}/play`, { duration }),
  approve: (id: number) => request.put(`/works/${id}/approve`),
  reject: (id: number, data: { reason: string }) => request.put(`/works/${id}/reject`, data),
  updateStatus: (id: number, data: { status: number }) => request.put(`/works/${id}/status`, data),
  listTags: (keyword?: string) => request.get<Tag[]>('/tags', { keyword }),
  listAlbums: (params: PageParams) => request.get<PageData<Album>>('/albums', params),
  getAlbums: (artistId: number) => request.get<Album[]>(`/albums/artist/${artistId}`),
  getAlbumById: (id: number) => request.get<Album>(`/albums/${id}`),
  createAlbum: (data: Partial<Album> & { work_ids?: number[] }) => request.post<Album>('/albums', data),
  updateAlbum: (id: number, data: Partial<Album> & { work_ids?: number[] }) => request.put(`/albums/${id}`, data),
  deleteAlbum: (id: number) => request.delete(`/albums/${id}`),
  addWorkToAlbum: (albumId: number, workId: number) => request.post(`/albums/${albumId}/works`, { work_id: workId })
}
