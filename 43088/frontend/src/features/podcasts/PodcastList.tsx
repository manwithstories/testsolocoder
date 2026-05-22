import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Plus, Search, RefreshCw, Trash2, Edit } from 'lucide-react'
import { useForm } from 'react-hook-form'
import * as yup from 'yup'
import { yupResolver } from '@hookform/resolvers/yup'
import {
  useGetPodcastsQuery,
  useAddPodcastMutation,
  useDeletePodcastMutation,
  useRefreshPodcastMutation,
} from '@/app/services/api'
import Pagination from '@/components/common/Pagination'
import { formatDate } from '@/utils/format'

const schema = yup.object({
  feed_url: yup.string().url('请输入有效的URL').required('请输入RSS订阅链接'),
})

type FormData = yup.InferType<typeof schema>

export default function PodcastList() {
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [showModal, setShowModal] = useState(false)
  const perPage = 12

  const { data, isLoading, refetch } = useGetPodcastsQuery({ page, perPage, search })
  const [addPodcast, { isLoading: adding }] = useAddPodcastMutation()
  const [deletePodcast] = useDeletePodcastMutation()
  const [refreshPodcast] = useRefreshPodcastMutation()

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<FormData>({
    resolver: yupResolver(schema),
  })

  const onSubmit = async (data: FormData) => {
    try {
      await addPodcast(data).unwrap()
      setShowModal(false)
      reset()
      refetch()
    } catch (err) {
      console.error('Failed to add podcast:', err)
    }
  }

  const handleDelete = async (id: string) => {
    if (window.confirm('确定要删除这个播客吗？')) {
      try {
        await deletePodcast(id).unwrap()
        refetch()
      } catch (err) {
        console.error('Failed to delete podcast:', err)
      }
    }
  }

  const handleRefresh = async (id: string) => {
    try {
      await refreshPodcast(id).unwrap()
      refetch()
    } catch (err) {
      console.error('Failed to refresh podcast:', err)
    }
  }

  if (isLoading) {
    return <div className="text-center py-12">加载中...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">播客订阅</h1>
        <button
          onClick={() => setShowModal(true)}
          className="btn btn-primary flex items-center gap-2"
        >
          <Plus className="w-5 h-5" />
          添加播客
        </button>
      </div>

      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
        <input
          type="text"
          placeholder="搜索播客..."
          value={search}
          onChange={(e) => {
            setSearch(e.target.value)
            setPage(1)
          }}
          className="input pl-10"
        />
      </div>

      {data?.data && data.data.length > 0 ? (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            {data.data.map((podcast) => (
              <div key={podcast.id} className="card group">
                <div className="relative">
                  <img
                    src={podcast.cover_image || '/placeholder.png'}
                    alt={podcast.title}
                    className="w-full aspect-square object-cover"
                  />
                  <div className="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
                    <Link
                      to={`/podcasts/${podcast.id}`}
                      className="p-3 bg-white rounded-full text-gray-900 hover:bg-gray-100"
                    >
                      <Edit className="w-5 h-5" />
                    </Link>
                    <button
                      onClick={() => handleRefresh(podcast.id)}
                      className="p-3 bg-white rounded-full text-gray-900 hover:bg-gray-100"
                    >
                      <RefreshCw className="w-5 h-5" />
                    </button>
                    <button
                      onClick={() => handleDelete(podcast.id)}
                      className="p-3 bg-white rounded-full text-red-600 hover:bg-gray-100"
                    >
                      <Trash2 className="w-5 h-5" />
                    </button>
                  </div>
                </div>
                <div className="p-4">
                  <Link
                    to={`/podcasts/${podcast.id}`}
                    className="font-semibold text-gray-900 hover:text-indigo-600 line-clamp-1"
                  >
                    {podcast.title}
                  </Link>
                  <p className="text-sm text-gray-500 line-clamp-1">{podcast.author}</p>
                  <p className="text-xs text-gray-400 mt-1">
                    更新于 {formatDate(podcast.last_updated)}
                  </p>
                </div>
              </div>
            ))}
          </div>

          {data.meta && data.meta.total_pages > 1 && (
            <Pagination
              currentPage={page}
              totalPages={data.meta.total_pages}
              onPageChange={setPage}
            />
          )}
        </>
      ) : (
        <div className="text-center py-12 card">
          <p className="text-gray-500">暂无播客订阅</p>
          <button
            onClick={() => setShowModal(true)}
            className="btn btn-primary mt-4"
          >
            添加第一个播客
          </button>
        </div>
      )}

      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 w-full max-w-md mx-4">
            <h2 className="text-xl font-bold mb-4">添加播客</h2>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  RSS 订阅链接
                </label>
                <input
                  type="url"
                  {...register('feed_url')}
                  placeholder="https://example.com/feed.xml"
                  className="input"
                />
                {errors.feed_url && (
                  <p className="text-red-500 text-sm mt-1">{errors.feed_url.message}</p>
                )}
              </div>
              <div className="flex gap-3 justify-end">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="btn btn-secondary"
                >
                  取消
                </button>
                <button type="submit" className="btn btn-primary" disabled={adding}>
                  {adding ? '添加中...' : '添加'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
