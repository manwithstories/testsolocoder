import { useParams, Link } from 'react-router-dom'
import { useDispatch } from 'react-redux'
import { useState } from 'react'
import { ArrowLeft, Play, Trash2, GripVertical } from 'lucide-react'
import { DragDropContext, Droppable, Draggable, DropResult } from 'react-beautiful-dnd'
import {
  useGetPlaylistQuery,
  useRemoveEpisodeFromPlaylistMutation,
  useReorderPlaylistItemsMutation,
} from '@/app/services/api'
import { setCurrentEpisode, setPlaying } from '@/features/player/playerSlice'
import { formatDuration } from '@/utils/format'

export default function PlaylistDetail() {
  const { id } = useParams<{ id: string }>()
  const dispatch = useDispatch()

  const { data: playlist, isLoading, refetch } = useGetPlaylistQuery(id!)
  const [removeEpisode] = useRemoveEpisodeFromPlaylistMutation()
  const [reorderItems] = useReorderPlaylistItemsMutation()
  const [isReordering, setIsReordering] = useState(false)

  const handlePlay = (episode: any) => {
    dispatch(setCurrentEpisode(episode))
    dispatch(setPlaying(true))
  }

  const handleRemove = async (itemId: string) => {
    if (window.confirm('确定要从播放列表中移除这个剧集吗？')) {
      try {
        await removeEpisode({ playlist_id: id!, item_id: itemId }).unwrap()
        refetch()
      } catch (err) {
        console.error('Failed to remove episode:', err)
      }
    }
  }

  const handleDragEnd = async (result: DropResult) => {
    if (!result.destination || !playlist?.items) return

    const items = Array.from(playlist.items)
    const [reorderedItem] = items.splice(result.source.index, 1)
    items.splice(result.destination.index, 0, reorderedItem)

    const itemIds = items.map((item) => item.id)

    try {
      setIsReordering(true)
      await reorderItems({ playlist_id: id!, item_ids: itemIds }).unwrap()
      refetch()
    } catch (err) {
      console.error('Failed to reorder items:', err)
    } finally {
      setIsReordering(false)
    }
  }

  if (isLoading) {
    return <div className="text-center py-12">加载中...</div>
  }

  if (!playlist) {
    return <div className="text-center py-12">播放列表不存在</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link to="/playlists" className="p-2 hover:bg-gray-100 rounded-lg">
          <ArrowLeft className="w-5 h-5" />
        </Link>
        <div className="flex-1">
          <h1 className="text-2xl font-bold text-gray-900">{playlist.name}</h1>
          {isReordering && <p className="text-sm text-indigo-600">正在保存排序...</p>}
        </div>
      </div>

      <div className="card p-6">
        {playlist.description && (
          <p className="text-gray-600 mb-4">{playlist.description}</p>
        )}
        <p className="text-sm text-gray-500">{playlist.items?.length || 0} 个剧集</p>
      </div>

      {playlist.items && playlist.items.length > 0 ? (
        <DragDropContext onDragEnd={handleDragEnd}>
          <Droppable droppableId="playlist-items">
            {(provided, snapshot) => (
              <div
                {...provided.droppableProps}
                ref={provided.innerRef}
                className={`space-y-2 ${snapshot.isDraggingOver ? 'bg-indigo-50 rounded-xl p-2' : ''}`}
              >
                {(playlist.items || []).map((item, index) => (
                  <Draggable key={item.id} draggableId={item.id} index={index}>
                    {(provided, snapshot) => (
                      <div
                        ref={provided.innerRef}
                        {...provided.draggableProps}
                        className={`card p-4 flex items-center gap-4 transition-all ${
                          snapshot.isDragging ? 'shadow-lg ring-2 ring-indigo-500' : ''
                        }`}
                      >
                        <div
                          {...provided.dragHandleProps}
                          className="text-gray-400 cursor-move hover:text-gray-600"
                        >
                          <GripVertical className="w-5 h-5" />
                        </div>
                        <span className="w-8 text-center text-gray-500">{index + 1}</span>
                        <button
                          onClick={() => handlePlay(item.episode)}
                          className="p-2 bg-indigo-600 text-white rounded-full hover:bg-indigo-700 transition-colors"
                        >
                          <Play className="w-4 h-4 ml-0.5" />
                        </button>
                        <div className="flex-1 min-w-0">
                          <Link
                            to={`/episodes/${item.episode_id}`}
                            className="font-medium text-gray-900 hover:text-indigo-600 line-clamp-1"
                          >
                            {item.episode?.title}
                          </Link>
                          <p className="text-sm text-gray-500 line-clamp-1">
                            {item.episode?.podcast?.title}
                          </p>
                        </div>
                        <span className="text-sm text-gray-400">
                          {formatDuration(item.episode?.duration || 0)}
                        </span>
                        <button
                          onClick={() => handleRemove(item.id)}
                          className="p-2 hover:bg-red-50 text-red-500 rounded-lg"
                        >
                          <Trash2 className="w-4 h-4" />
                        </button>
                      </div>
                    )}
                  </Draggable>
                ))}
                {provided.placeholder}
              </div>
            )}
          </Droppable>
        </DragDropContext>
      ) : (
        <div className="text-center py-12 card">
          <p className="text-gray-500">播放列表为空</p>
          <p className="text-sm text-gray-400 mt-2">
            进入剧集页面，将剧集添加到此播放列表
          </p>
        </div>
      )}
    </div>
  )
}
