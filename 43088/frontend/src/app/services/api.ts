import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import type {
  Podcast,
  PodcastStats,
  Episode,
  PlaybackProgress,
  Note,
  Playlist,
  PlaylistItem,
  ListeningHistory,
  ListeningStats,
} from '@/types'

interface PaginatedResponse<T> {
  data: T[]
  meta: {
    total: number
    page: number
    per_page: number
    total_pages: number
  }
}

export const api = createApi({
  reducerPath: 'api',
  baseQuery: fetchBaseQuery({ baseUrl: '/api' }),
  tagTypes: ['Podcast', 'Episode', 'Progress', 'Note', 'Playlist', 'History'],
  endpoints: (builder) => ({
    getPodcasts: builder.query<
      PaginatedResponse<Podcast>,
      { page?: number; perPage?: number; search?: string }
    >({
      query: ({ page = 1, perPage = 10, search = '' }) =>
        `podcasts?page=${page}&per_page=${perPage}&search=${encodeURIComponent(search)}`,
      providesTags: ['Podcast'],
    }),

    getPodcast: builder.query<
      { podcast: Podcast; stats: PodcastStats },
      string
    >({
      query: (id) => `podcasts/${id}`,
      providesTags: ['Podcast'],
    }),

    addPodcast: builder.mutation<Podcast, { feed_url: string }>({
      query: (body) => ({
        url: 'podcasts',
        method: 'POST',
        body,
      }),
      invalidatesTags: ['Podcast'],
    }),

    updatePodcast: builder.mutation<
      Podcast,
      { id: string; data: Partial<Podcast> }
    >({
      query: ({ id, data }) => ({
        url: `podcasts/${id}`,
        method: 'PUT',
        body: data,
      }),
      invalidatesTags: ['Podcast'],
    }),

    deletePodcast: builder.mutation<void, string>({
      query: (id) => ({
        url: `podcasts/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['Podcast', 'Episode'],
    }),

    refreshPodcast: builder.mutation<
      { podcast: Podcast; new_episodes: number },
      string
    >({
      query: (id) => ({
        url: `podcasts/${id}/refresh`,
        method: 'POST',
      }),
      invalidatesTags: ['Podcast', 'Episode'],
    }),

    getEpisodes: builder.query<
      PaginatedResponse<Episode>,
      {
        page?: number
        perPage?: number
        search?: string
        podcast_id?: string
      }
    >({
      query: ({ page = 1, perPage = 20, search = '', podcast_id = '' }) => {
        let url = `episodes?page=${page}&per_page=${perPage}&search=${encodeURIComponent(search)}`
        if (podcast_id) url += `&podcast_id=${podcast_id}`
        return url
      },
      providesTags: ['Episode'],
    }),

    getEpisode: builder.query<
      { episode: Episode; progress: PlaybackProgress },
      string
    >({
      query: (id) => `episodes/${id}`,
      providesTags: ['Episode', 'Progress'],
    }),

    updatePlaybackProgress: builder.mutation<
      PlaybackProgress,
      { id: string; current_time: number; duration: number }
    >({
      query: ({ id, ...data }) => ({
        url: `episodes/${id}/progress`,
        method: 'POST',
        body: data,
      }),
      invalidatesTags: ['Progress'],
    }),

    getPlaybackProgress: builder.query<PlaybackProgress, string>({
      query: (id) => `episodes/${id}/progress`,
      providesTags: ['Progress'],
    }),

    markAsCompleted: builder.mutation<{ message: string }, string>({
      query: (id) => ({
        url: `episodes/${id}/complete`,
        method: 'POST',
      }),
      invalidatesTags: ['Progress'],
    }),

    addListeningHistory: builder.mutation<
      ListeningHistory,
      {
        id: string
        start_time: string
        end_time: string
        duration: number
        completion: number
      }
    >({
      query: ({ id, ...data }) => ({
        url: `episodes/${id}/history`,
        method: 'POST',
        body: data,
      }),
      invalidatesTags: ['History'],
    }),

    getListeningHistory: builder.query<
      PaginatedResponse<ListeningHistory>,
      {
        page?: number
        perPage?: number
        podcast_id?: string
        start_date?: string
        end_date?: string
        completed?: boolean
      }
    >({
      query: ({
        page = 1,
        perPage = 20,
        podcast_id = '',
        start_date = '',
        end_date = '',
        completed,
      }) => {
        let url = `history?page=${page}&per_page=${perPage}`
        if (podcast_id) url += `&podcast_id=${podcast_id}`
        if (start_date) url += `&start_date=${start_date}`
        if (end_date) url += `&end_date=${end_date}`
        if (completed !== undefined) url += `&completed=${completed}`
        return url
      },
      providesTags: ['History'],
    }),

    getNotes: builder.query<Note[], { episode_id: string; search?: string; tag?: string }>({
      query: ({ episode_id, search = '', tag = '' }) =>
        `notes/episode/${episode_id}?search=${encodeURIComponent(search)}&tag=${encodeURIComponent(tag)}`,
      providesTags: ['Note'],
    }),

    searchNotes: builder.query<Note[], string>({
      query: (q) => `notes/search?q=${encodeURIComponent(q)}`,
      providesTags: ['Note'],
    }),

    addNote: builder.mutation<
      Note,
      { episode_id: string; timestamp: number; content: string; tags: string[] }
    >({
      query: ({ episode_id, ...data }) => ({
        url: `notes/episode/${episode_id}`,
        method: 'POST',
        body: data,
      }),
      invalidatesTags: ['Note'],
    }),

    updateNote: builder.mutation<
      Note,
      { id: string; content: string; tags: string[] }
    >({
      query: ({ id, ...data }) => ({
        url: `notes/${id}`,
        method: 'PUT',
        body: data,
      }),
      invalidatesTags: ['Note'],
    }),

    deleteNote: builder.mutation<void, string>({
      query: (id) => ({
        url: `notes/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['Note'],
    }),

    getPlaylists: builder.query<Playlist[], void>({
      query: () => 'playlists',
      providesTags: ['Playlist'],
    }),

    getPlaylist: builder.query<Playlist, string>({
      query: (id) => `playlists/${id}`,
      providesTags: ['Playlist'],
    }),

    createPlaylist: builder.mutation<
      Playlist,
      { name: string; description?: string; cover_image?: string }
    >({
      query: (data) => ({
        url: 'playlists',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: ['Playlist'],
    }),

    updatePlaylist: builder.mutation<
      Playlist,
      { id: string; data: Partial<Playlist> }
    >({
      query: ({ id, data }) => ({
        url: `playlists/${id}`,
        method: 'PUT',
        body: data,
      }),
      invalidatesTags: ['Playlist'],
    }),

    deletePlaylist: builder.mutation<void, string>({
      query: (id) => ({
        url: `playlists/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['Playlist'],
    }),

    addEpisodeToPlaylist: builder.mutation<
      PlaylistItem,
      { playlist_id: string; episode_id: string }
    >({
      query: ({ playlist_id, episode_id }) => ({
        url: `playlists/${playlist_id}/episodes`,
        method: 'POST',
        body: { episode_id },
      }),
      invalidatesTags: ['Playlist'],
    }),

    removeEpisodeFromPlaylist: builder.mutation<
      void,
      { playlist_id: string; item_id: string }
    >({
      query: ({ playlist_id, item_id }) => ({
        url: `playlists/${playlist_id}/episodes/${item_id}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['Playlist'],
    }),

    reorderPlaylistItems: builder.mutation<
      { message: string },
      { playlist_id: string; item_ids: string[] }
    >({
      query: ({ playlist_id, item_ids }) => ({
        url: `playlists/${playlist_id}/reorder`,
        method: 'POST',
        body: { item_ids },
      }),
      invalidatesTags: ['Playlist'],
    }),

    getListeningStats: builder.query<ListeningStats[], number>({
      query: (days = 30) => `stats/listening?days=${days}`,
    }),

    getPodcastDistribution: builder.query<any[], void>({
      query: () => 'stats/distribution',
    }),

    getCompletionStats: builder.query<any, void>({
      query: () => 'stats/completion',
    }),

    getListeningHabits: builder.query<any, void>({
      query: () => 'stats/habits',
    }),

    getNewEpisodesCount: builder.query<{ new_episodes_count: number }, void>({
      query: () => 'episodes/new-count',
    }),

    exportOPML: builder.mutation<Blob, void>({
      query: () => ({
        url: 'export/opml',
        responseHandler: (response) => response.blob(),
      }),
    }),

    exportHistoryCSV: builder.mutation<Blob, void>({
      query: () => ({
        url: 'export/history/csv',
        responseHandler: (response) => response.blob(),
      }),
    }),

    exportNotesCSV: builder.mutation<Blob, void>({
      query: () => ({
        url: 'export/notes/csv',
        responseHandler: (response) => response.blob(),
      }),
    }),

    importOPML: builder.mutation<{ imported_count: number; skipped_count: number; message: string }, FormData>({
      query: (formData) => ({
        url: 'import/opml',
        method: 'POST',
        body: formData,
      }),
      invalidatesTags: ['Podcast'],
    }),
  }),
})

export const {
  useGetPodcastsQuery,
  useGetPodcastQuery,
  useAddPodcastMutation,
  useUpdatePodcastMutation,
  useDeletePodcastMutation,
  useRefreshPodcastMutation,
  useGetEpisodesQuery,
  useGetEpisodeQuery,
  useUpdatePlaybackProgressMutation,
  useGetPlaybackProgressQuery,
  useMarkAsCompletedMutation,
  useAddListeningHistoryMutation,
  useGetListeningHistoryQuery,
  useGetNotesQuery,
  useSearchNotesQuery,
  useAddNoteMutation,
  useUpdateNoteMutation,
  useDeleteNoteMutation,
  useGetPlaylistsQuery,
  useGetPlaylistQuery,
  useCreatePlaylistMutation,
  useUpdatePlaylistMutation,
  useDeletePlaylistMutation,
  useAddEpisodeToPlaylistMutation,
  useRemoveEpisodeFromPlaylistMutation,
  useReorderPlaylistItemsMutation,
  useGetListeningStatsQuery,
  useGetPodcastDistributionQuery,
  useGetCompletionStatsQuery,
  useGetListeningHabitsQuery,
  useGetNewEpisodesCountQuery,
  useExportOPMLMutation,
  useExportHistoryCSVMutation,
  useExportNotesCSVMutation,
  useImportOPMLMutation,
} = api
