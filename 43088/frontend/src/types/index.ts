export interface Podcast {
  id: string
  title: string
  description: string
  feed_url: string
  website: string
  author: string
  cover_image: string
  language: string
  category: string
  last_checked: string
  last_updated: string
  created_at: string
  updated_at: string
}

export interface PodcastStats {
  podcast_id: string
  total_episodes: number
  unplayed_count: number
  completed_count: number
  total_listened_seconds: number
}

export interface Episode {
  id: string
  podcast_id: string
  title: string
  description: string
  guid: string
  audio_url: string
  audio_type: string
  duration: number
  pub_date: string
  episode_type: string
  season_number: number
  episode_number: number
  is_new: boolean
  created_at: string
  updated_at: string
  podcast?: Podcast
}

export interface PlaybackProgress {
  id: string
  episode_id: string
  current_time: number
  completed: boolean
  completed_at: string
  play_count: number
  last_played_at: string
  created_at: string
  updated_at: string
}

export interface Note {
  id: string
  episode_id: string
  timestamp: number
  content: string
  tags: string[]
  created_at: string
  updated_at: string
  episode?: Episode
}

export interface Playlist {
  id: string
  name: string
  description: string
  cover_image: string
  created_at: string
  updated_at: string
  items?: PlaylistItem[]
}

export interface PlaylistItem {
  id: string
  playlist_id: string
  episode_id: string
  position: number
  created_at: string
  updated_at: string
  episode?: Episode
}

export interface ListeningHistory {
  id: string
  episode_id: string
  start_time: string
  end_time: string
  duration: number
  completion: number
  created_at: string
  updated_at: string
  episode?: Episode
}

export interface ListeningStats {
  date: string
  total_duration: number
  episode_count: number
  completion_rate: number
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
  meta?: {
    total: number
    page: number
    per_page: number
    total_pages: number
  }
}

export interface PlayerState {
  currentEpisode: Episode | null
  isPlaying: boolean
  currentTime: number
  duration: number
  volume: number
  playbackRate: number
  isMuted: boolean
  skipSilence: boolean
  silenceThreshold: number
  silenceMinDuration: number
}
