import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Calendar, Filter, Play } from 'lucide-react'
import { useDispatch } from 'react-redux'
import {
  useGetListeningHistoryQuery,
  useGetPodcastsQuery,
} from '@/app/services/api'
import { setCurrentEpisode, setPlaying } from '@/features/player/playerSlice'
import Pagination from '@/components/common/Pagination'
import { formatDateTime, formatDuration } from '@/utils/format'

export default function HistoryPage() {
  const dispatch = useDispatch()
  const [page, setPage] = useState(1)
  const [podcastId, setPodcastId] = useState('')
  const [startDate, setStartDate] = useState('')
  const [endDate, setEndDate] = useState('')
  const [completed, setCompleted] = useState<string>('')
  const perPage = 20

  const { data: historyData, isLoading } = useGetListeningHistoryQuery({
    page,
    perPage,
    podcast_id: podcastId || undefined,
    start_date: startDate || undefined,
    end_date: endDate || undefined,
    completed: completed ? completed === 'true' : undefined,
  })
  const { data: podcasts } = useGetPodcastsQuery({ page: 1, perPage: 100 })

  const handlePlay = (episode: any) => {
    dispatch(setCurrentEpisode(episode))
    dispatch(setPlaying(true))
  }

  const clearFilters = () => {
    setPodcastId('')
    setStartDate('')
    setEndDate('')
    setCompleted('')
    setPage(1)
  }

  if (isLoading) {
    return <div className="text-center py-12">加载中...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">收听历史</h1>
      </div>

      <div className="card p-4">
        <div className="flex items-center gap-2 mb-4">
          <Filter className="w-5 h-5 text-gray-500" />
          <span className="font-medium text-gray-700">筛选条件</span>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div>
            <label className="block text-sm text-gray-600 mb-1">播客</label>
            <select
              value={podcastId}
              onChange={(e) => {
                setPodcastId(e.target.value)
                setPage(1)
              }}
              className="input"
            >
              <option value="">全部播客</option>
              {podcasts?.data?.map((podcast) => (
                <option key={podcast.id} value={podcast.id}>
                  {podcast.title}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm text-gray-600 mb-1">开始日期</label>
            <input
              type="date"
              value={startDate}
              onChange={(e) => {
                setStartDate(e.target.value)
                setPage(1)
              }}
              className="input"
            />
          </div>
          <div>
            <label className="block text-sm text-gray-600 mb-1">结束日期</label>
            <input
              type="date"
              value={endDate}
              onChange={(e) => {
                setEndDate(e.target.value)
                setPage(1)
              }}
              className="input"
            />
          </div>
          <div>
            <label className="block text-sm text-gray-600 mb-1">完成状态</label>
            <select
              value={completed}
              onChange={(e) => {
                setCompleted(e.target.value)
                setPage(1)
              }}
              className="input"
            >
              <option value="">全部</option>
              <option value="true">已完成</option>
              <option value="false">未完成</option>
            </select>
          </div>
        </div>
        {(podcastId || startDate || endDate || completed) && (
          <button
            onClick={clearFilters}
            className="mt-4 text-sm text-indigo-600 hover:underline"
          >
            清除筛选
          </button>
        )}
      </div>

      {historyData?.data && historyData.data.length > 0 ? (
        <>
          <div className="space-y-3">
            {historyData.data.map((history) => (
              <div key={history.id} className="card p-4">
                <div className="flex items-center gap-4">
                  <button
                    onClick={() => handlePlay(history.episode)}
                    className="p-3 bg-indigo-600 text-white rounded-full hover:bg-indigo-700 transition-colors"
                  >
                    <Play className="w-5 h-5 ml-0.5" />
                  </button>
                  <div className="flex-1 min-w-0">
                    <Link
                      to={`/episodes/${history.episode_id}`}
                      className="font-semibold text-gray-900 hover:text-indigo-600"
                    >
                      {history.episode?.title}
                    </Link>
                    <p className="text-sm text-gray-500">
                      {history.episode?.podcast?.title}
                    </p>
                    <div className="flex items-center gap-4 mt-2 text-xs text-gray-400">
                      <span className="flex items-center gap-1">
                        <Calendar className="w-3 h-3" />
                        {formatDateTime(history.start_time)}
                      </span>
                      <span>收听时长: {formatDuration(history.duration)}</span>
                      <span>
                        完成度: {Math.round(history.completion * 100)}%
                      </span>
                    </div>
                  </div>
                  <div className="text-right">
                    <div
                      className={`badge ${
                        history.completion >= 0.9
                          ? 'bg-green-100 text-green-700'
                          : 'bg-yellow-100 text-yellow-700'
                      }`}
                    >
                      {history.completion >= 0.9 ? '已完成' : '进行中'}
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>

          {historyData.meta && historyData.meta.total_pages > 1 && (
            <Pagination
              currentPage={page}
              totalPages={historyData.meta.total_pages}
              onPageChange={setPage}
            />
          )}
        </>
      ) : (
        <div className="text-center py-12 card">
          <p className="text-gray-500">暂无收听记录</p>
        </div>
      )}
    </div>
  )
}
