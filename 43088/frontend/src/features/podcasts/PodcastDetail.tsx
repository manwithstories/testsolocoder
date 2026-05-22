import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useDispatch } from 'react-redux'
import { ArrowLeft, Play, RefreshCw, Search } from 'lucide-react'
import {
  useGetPodcastQuery,
  useGetEpisodesQuery,
  useRefreshPodcastMutation,
} from '@/app/services/api'
import { setCurrentEpisode, setPlaying } from '@/features/player/playerSlice'
import Pagination from '@/components/common/Pagination'
import { formatDate, formatDuration } from '@/utils/format'

export default function PodcastDetail() {
  const { id } = useParams<{ id: string }>()
  const dispatch = useDispatch()
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const perPage = 20

  const { data: podcastData, isLoading: loadingPodcast } = useGetPodcastQuery(id!)
  const { data: episodesData, isLoading: loadingEpisodes } = useGetEpisodesQuery({
    page,
    perPage,
    search,
    podcast_id: id,
  })
  const [refreshPodcast] = useRefreshPodcastMutation()

  const handlePlay = (episode: any) => {
    dispatch(setCurrentEpisode(episode))
    dispatch(setPlaying(true))
  }

  const handleRefresh = async () => {
    if (id) {
      try {
        await refreshPodcast(id).unwrap()
      } catch (err) {
        console.error('Failed to refresh podcast:', err)
      }
    }
  }

  if (loadingPodcast || loadingEpisodes) {
    return <div className="text-center py-12">加载中...</div>
  }

  const podcast = podcastData?.podcast
  const stats = podcastData?.stats

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link to="/" className="p-2 hover:bg-gray-100 rounded-lg">
          <ArrowLeft className="w-5 h-5" />
        </Link>
        <h1 className="text-2xl font-bold text-gray-900">播客详情</h1>
      </div>

      {podcast && (
        <div className="card p-6">
          <div className="flex flex-col md:flex-row gap-6">
            <img
              src={podcast.cover_image || '/placeholder.png'}
              alt={podcast.title}
              className="w-48 h-48 rounded-xl object-cover"
            />
            <div className="flex-1">
              <div className="flex items-start justify-between">
                <div>
                  <h2 className="text-2xl font-bold text-gray-900">{podcast.title}</h2>
                  <p className="text-gray-600 mt-1">{podcast.author}</p>
                  {podcast.category && (
                    <span className="badge bg-indigo-100 text-indigo-700 mt-2">
                      {podcast.category}
                    </span>
                  )}
                </div>
                <button
                  onClick={handleRefresh}
                  className="btn btn-secondary flex items-center gap-2"
                >
                  <RefreshCw className="w-4 h-4" />
                  刷新
                </button>
              </div>

              <p className="text-gray-600 mt-4 line-clamp-4">{podcast.description}</p>

              <div className="grid grid-cols-3 gap-4 mt-6">
                <div className="text-center p-4 bg-gray-50 rounded-lg">
                  <p className="text-2xl font-bold text-gray-900">{stats?.total_episodes || 0}</p>
                  <p className="text-sm text-gray-500">总集数</p>
                </div>
                <div className="text-center p-4 bg-gray-50 rounded-lg">
                  <p className="text-2xl font-bold text-orange-500">{stats?.unplayed_count || 0}</p>
                  <p className="text-sm text-gray-500">未播放</p>
                </div>
                <div className="text-center p-4 bg-gray-50 rounded-lg">
                  <p className="text-2xl font-bold text-green-500">{stats?.completed_count || 0}</p>
                  <p className="text-sm text-gray-500">已完成</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      <div className="flex items-center justify-between">
        <h3 className="text-xl font-semibold text-gray-900">剧集列表</h3>
        <div className="relative w-64">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input
            type="text"
            placeholder="搜索剧集..."
            value={search}
            onChange={(e) => {
              setSearch(e.target.value)
              setPage(1)
            }}
            className="input pl-9 text-sm"
          />
        </div>
      </div>

      {episodesData?.data && episodesData.data.length > 0 ? (
        <>
          <div className="space-y-3">
            {episodesData.data.map((episode) => (
              <div
                key={episode.id}
                className="card p-4 hover:shadow-md transition-shadow"
              >
                <div className="flex items-center gap-4">
                  <button
                    onClick={() => handlePlay(episode)}
                    className="p-3 bg-indigo-600 text-white rounded-full hover:bg-indigo-700 transition-colors"
                  >
                    <Play className="w-5 h-5 ml-0.5" />
                  </button>
                  <div className="flex-1 min-w-0">
                    <Link
                      to={`/episodes/${episode.id}`}
                      className="font-semibold text-gray-900 hover:text-indigo-600 line-clamp-1"
                    >
                      {episode.title}
                      {episode.is_new && (
                        <span className="ml-2 badge bg-red-100 text-red-700">NEW</span>
                      )}
                    </Link>
                    <p className="text-sm text-gray-500 line-clamp-1 mt-1">
                      {episode.description}
                    </p>
                    <div className="flex items-center gap-4 mt-2 text-xs text-gray-400">
                      <span>{formatDate(episode.pub_date)}</span>
                      <span>{formatDuration(episode.duration)}</span>
                      {episode.season_number > 0 && (
                        <span>
                          S{episode.season_number} E{episode.episode_number}
                        </span>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>

          {episodesData.meta && episodesData.meta.total_pages > 1 && (
            <Pagination
              currentPage={page}
              totalPages={episodesData.meta.total_pages}
              onPageChange={setPage}
            />
          )}
        </>
      ) : (
        <div className="text-center py-12 card">
          <p className="text-gray-500">暂无剧集</p>
        </div>
      )}
    </div>
  )
}
