import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Plus, Trash2, Edit, ListMusic } from 'lucide-react'
import { useForm } from 'react-hook-form'
import * as yup from 'yup'
import { yupResolver } from '@hookform/resolvers/yup'
import {
  useGetPlaylistsQuery,
  useCreatePlaylistMutation,
  useDeletePlaylistMutation,
} from '@/app/services/api'

const schema = yup.object({
  name: yup.string().required('请输入播放列表名称'),
  description: yup.string(),
})

type FormData = yup.InferType<typeof schema>

export default function PlaylistList() {
  const [showModal, setShowModal] = useState(false)
  const { data: playlists, isLoading, refetch } = useGetPlaylistsQuery()
  const [createPlaylist, { isLoading: creating }] = useCreatePlaylistMutation()
  const [deletePlaylist] = useDeletePlaylistMutation()

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
      await createPlaylist(data).unwrap()
      setShowModal(false)
      reset()
      refetch()
    } catch (err) {
      console.error('Failed to create playlist:', err)
    }
  }

  const handleDelete = async (id: string) => {
    if (window.confirm('确定要删除这个播放列表吗？')) {
      try {
        await deletePlaylist(id).unwrap()
        refetch()
      } catch (err) {
        console.error('Failed to delete playlist:', err)
      }
    }
  }

  if (isLoading) {
    return <div className="text-center py-12">加载中...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">播放列表</h1>
        <button
          onClick={() => setShowModal(true)}
          className="btn btn-primary flex items-center gap-2"
        >
          <Plus className="w-5 h-5" />
          创建播放列表
        </button>
      </div>

      {playlists && playlists.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {playlists.map((playlist) => (
            <div key={playlist.id} className="card group">
              <div className="p-4">
                <div className="flex items-start justify-between">
                  <div className="flex items-center gap-3">
                    <div className="w-12 h-12 bg-indigo-100 rounded-lg flex items-center justify-center">
                      <ListMusic className="w-6 h-6 text-indigo-600" />
                    </div>
                    <div>
                      <Link
                        to={`/playlists/${playlist.id}`}
                        className="font-semibold text-gray-900 hover:text-indigo-600"
                      >
                        {playlist.name}
                      </Link>
                      <p className="text-sm text-gray-500">
                        {playlist.items?.length || 0} 个剧集
                      </p>
                    </div>
                  </div>
                  <div className="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                    <button className="p-2 hover:bg-gray-100 rounded-lg">
                      <Edit className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => handleDelete(playlist.id)}
                      className="p-2 hover:bg-red-50 text-red-500 rounded-lg"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </div>
                {playlist.description && (
                  <p className="text-sm text-gray-600 mt-3 line-clamp-2">
                    {playlist.description}
                  </p>
                )}
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="text-center py-12 card">
          <p className="text-gray-500">暂无播放列表</p>
          <button
            onClick={() => setShowModal(true)}
            className="btn btn-primary mt-4"
          >
            创建第一个播放列表
          </button>
        </div>
      )}

      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 w-full max-w-md mx-4">
            <h2 className="text-xl font-bold mb-4">创建播放列表</h2>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  名称
                </label>
                <input
                  type="text"
                  {...register('name')}
                  placeholder="输入播放列表名称"
                  className="input"
                />
                {errors.name && (
                  <p className="text-red-500 text-sm mt-1">{errors.name.message}</p>
                )}
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  描述（可选）
                </label>
                <textarea
                  {...register('description')}
                  rows={3}
                  placeholder="输入播放列表描述"
                  className="input resize-none"
                />
              </div>
              <div className="flex gap-3 justify-end">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="btn btn-secondary"
                >
                  取消
                </button>
                <button type="submit" className="btn btn-primary" disabled={creating}>
                  {creating ? '创建中...' : '创建'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
