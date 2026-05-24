import request, { PageData, PageParams } from '@/utils/request'

export interface Comment {
  id: number
  user_id: number
  work_id?: number
  album_id?: number
  playlist_id?: number
  parent_id?: number
  content: string
  like_count: number
  reply_count: number
  is_pinned: boolean
  user?: {
    id: number
    username: string
    nickname: string
    avatar: string
  }
  replies?: Comment[]
  created_at: string
}

export interface Playlist {
  id: number
  user_id: number
  name: string
  description?: string
  cover_url: string
  is_public: boolean
  play_count: number
  like_count: number
  follow_count: number
  work_count: number
  user?: {
    id: number
    username: string
    nickname: string
    avatar: string
  }
  works?: any[]
  created_at: string
  updated_at: string
}

export interface Notification {
  id: number
  user_id: number
  type: string
  title: string
  content: string
  data?: string
  is_read: boolean
  created_at: string
}

export interface PlayRecord {
  id: number
  user_id: number
  work_id: number
  artist_id: number
  duration: number
  ip: string
  device: string
  source: string
  work?: any
  created_at: string
}

export const communityApi = {
  follow: (followingId: number) => request.post('/follow', { following_id: followingId }),
  unfollow: (followingId: number) => request.post('/unfollow', { following_id: followingId }),
  isFollowing: (userId: number) => request.get<{ is_following: boolean }>(`/follow/is-following/${userId}`),
  getFollowers: (userId: number, params: PageParams) => request.get<PageData<any>>(`/users/${userId}/followers`, params),
  getFollowings: (userId: number, params: PageParams) => request.get<PageData<any>>(`/users/${userId}/followings`, params),
  getFollowingArtists: (params: PageParams) => request.get<PageData<any>>('/follow/artists', params),
  getFollowingUsers: (params: PageParams) => request.get<PageData<any>>('/follow/users', params),
  getMyFollowers: (params: PageParams) => request.get<PageData<any>>('/follow/followers', params),
  getFollowerCount: (userId: number) => request.get<{ count: number }>(`/users/${userId}/follower-count`),
  getFollowingCount: (userId: number) => request.get<{ count: number }>(`/users/${userId}/following-count`),
  
  createComment: (data: { work_id?: number; album_id?: number; playlist_id?: number; parent_id?: number; content: string }) => 
    request.post<Comment>('/comments', data),
  deleteComment: (id: number) => request.delete(`/comments/${id}`),
  getComments: (params: { work_id?: number; album_id?: number; playlist_id?: number; page?: number; page_size?: number }) => 
    request.get<PageData<Comment>>('/comments', params),
  
  createPlaylist: (data: { name: string; description?: string; cover_url?: string; is_public?: boolean }) => 
    request.post<Playlist>('/playlists', data),
  updatePlaylist: (id: number, data: { name?: string; description?: string; cover_url?: string; is_public?: boolean }) => 
    request.put(`/playlists/${id}`, data),
  deletePlaylist: (id: number) => request.delete(`/playlists/${id}`),
  getPlaylistById: (id: number) => request.get<Playlist>(`/playlists/${id}`),
  getMyPlaylists: () => request.get<Playlist[]>('/my-playlists'),
  listPlaylists: (params: { user_id?: number; keyword?: string; page?: number; page_size?: number }) => 
    request.get<PageData<Playlist>>('/playlists', params),
  addWorkToPlaylist: (playlistId: number, workId: number) => 
    request.post('/playlists/works', { playlist_id: playlistId, work_id: workId }),
  removeWorkFromPlaylist: (playlistId: number, workId: number) => 
    request.delete('/playlists/works', { playlist_id: playlistId, work_id: workId }),
  
  getNotifications: (params: PageParams) => request.get<PageData<Notification>>('/notifications', params),
  markNotificationAsRead: (id: number) => request.put(`/notifications/${id}/read`),
  markAllNotificationsAsRead: () => request.put('/notifications/read-all'),
  getUnreadNotificationCount: () => request.get<{ count: number }>('/notifications/unread-count'),
  
  getPlayRecords: (params: PageParams) => request.get<PageData<PlayRecord>>('/play-records', params)
}
