import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useDispatch, useSelector } from 'react-redux'
import { ArrowLeft, Play, Plus, Tag, Trash2, Edit, Download } from 'lucide-react'
import { useForm } from 'react-hook-form'
import * as yup from 'yup'
import { yupResolver } from '@hookform/resolvers/yup'
import {
  useGetEpisodeQuery,
  useGetNotesQuery,
  useAddNoteMutation,
  useDeleteNoteMutation,
  useUpdateNoteMutation,
  useGetPlaylistsQuery,
  useAddEpisodeToPlaylistMutation,
} from '@/app/services/api'
import { setCurrentEpisode as setPlayerEpisode, setPlaying } from '@/features/player/playerSlice'
import { RootState } from '@/app/store'
import { formatDate, formatDuration, formatTime } from '@/utils/format'

const noteSchema = yup.object({
  timestamp: yup.number().min(0).required(),
  content: yup.string().required('请输入笔记内容'),
  tags: yup.array().of(yup.string()).default([]),
})

type NoteFormData = yup.InferType<typeof noteSchema>

export default function EpisodeDetail() {
  const { id } = useParams<{ id: string }>()
  const dispatch = useDispatch()
  const { currentTime } = useSelector((state: RootState) => state.player)
  
  const [editingNote, setEditingNote] = useState<string | null>(null)
  const [showAddToPlaylist, setShowAddToPlaylist] = useState(false)
  const [tagInput, setTagInput] = useState('')

  const { data: episodeData, isLoading } = useGetEpisodeQuery(id!)
  const { data: notes, refetch: refetchNotes } = useGetNotesQuery({ episode_id: id! })
  const { data: playlists } = useGetPlaylistsQuery()
  const [addNote] = useAddNoteMutation()
  const [deleteNote] = useDeleteNoteMutation()
  const [updateNote] = useUpdateNoteMutation()
  const [addToPlaylist] = useAddEpisodeToPlaylistMutation()

  const { register, handleSubmit, setValue, reset, getValues, formState: { errors } } = useForm<NoteFormData>({
    resolver: yupResolver(noteSchema),
    defaultValues: {
      timestamp: 0,
      content: '',
      tags: [],
    },
  })

  const handlePlay = () => {
    if (episodeData?.episode) {
      dispatch(setPlayerEpisode(episodeData.episode))
      dispatch(setPlaying(true))
    }
  }

  const onSubmitNote = async (data: NoteFormData) => {
    if (!id) return
    try {
      const tags = (data.tags || []).filter(Boolean) as string[]
      if (editingNote) {
        await updateNote({ id: editingNote, content: data.content, tags }).unwrap()
        setEditingNote(null)
      } else {
        await addNote({
          episode_id: id,
          timestamp: data.timestamp,
          content: data.content,
          tags,
        }).unwrap()
      }
      reset()
      setTagInput('')
      refetchNotes()
    } catch (err) {
      console.error('Failed to save note:', err)
    }
  }

  const handleAddTag = () => {
    if (!tagInput.trim()) return
    const currentTags = getValues('tags') || []
    if (!currentTags.includes(tagInput.trim())) {
      setValue('tags', [...currentTags, tagInput.trim()])
    }
    setTagInput('')
  }

  const handleRemoveTag = (tag: string) => {
    const currentTags = getValues('tags') || []
    setValue('tags', currentTags.filter((t: string | undefined) => t !== tag))
  }

  const handleEditNote = (note: any) => {
    setEditingNote(note.id)
    setValue('content', note.content)
    setValue('tags', note.tags)
    setValue('timestamp', note.timestamp)
  }

  const handleDeleteNote = async (noteId: string) => {
    if (window.confirm('确定要删除这条笔记吗？')) {
      try {
        await deleteNote(noteId).unwrap()
        refetchNotes()
      } catch (err) {
        console.error('Failed to delete note:', err)
      }
    }
  }

  const handleAddToPlaylist = async (playlistId: string) => {
    if (!id) return
    try {
      await addToPlaylist({ playlist_id: playlistId, episode_id: id }).unwrap()
      setShowAddToPlaylist(false)
    } catch (err) {
      console.error('Failed to add to playlist:', err)
    }
  }

  const handleExportNotes = () => {
    window.open(`/api/export/notes/markdown?episode_id=${id}`, '_blank')
  }

  if (isLoading) {
    return <div className="text-center py-12">加载中...</div>
  }

  const episode = episodeData?.episode
  const progress = episodeData?.progress

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link to={episode?.podcast_id ? `/podcasts/${episode.podcast_id}` : '/'} className="p-2 hover:bg-gray-100 rounded-lg">
          <ArrowLeft className="w-5 h-5" />
        </Link>
        <h1 className="text-2xl font-bold text-gray-900">剧集详情</h1>
      </div>

      {episode && (
        <div className="card p-6">
          <div className="flex flex-col md:flex-row gap-6">
            <div className="relative">
              <img
                src={episode.podcast?.cover_image || '/placeholder.png'}
                alt={episode.title}
                className="w-48 h-48 rounded-xl object-cover"
              />
              <button
                onClick={handlePlay}
                className="absolute inset-0 flex items-center justify-center bg-black/30 rounded-xl opacity-0 hover:opacity-100 transition-opacity"
              >
                <div className="p-4 bg-white rounded-full">
                  <Play className="w-8 h-8 text-indigo-600 ml-1" />
                </div>
              </button>
            </div>
            <div className="flex-1">
              <div className="flex items-start justify-between">
                <div>
                  <h2 className="text-2xl font-bold text-gray-900">{episode.title}</h2>
                  <Link
                    to={`/podcasts/${episode.podcast_id}`}
                    className="text-indigo-600 hover:underline"
                  >
                    {episode.podcast?.title}
                  </Link>
                </div>
                <div className="flex gap-2">
                  <button
                    onClick={handlePlay}
                    className="btn btn-primary flex items-center gap-2"
                  >
                    <Play className="w-4 h-4" />
                    播放
                  </button>
                  <button
                    onClick={() => setShowAddToPlaylist(true)}
                    className="btn btn-secondary flex items-center gap-2"
                  >
                    <Plus className="w-4 h-4" />
                    加入列表
                  </button>
                </div>
              </div>

              <div className="flex items-center gap-4 mt-4 text-sm text-gray-500">
                <span>{formatDate(episode.pub_date)}</span>
                <span>{formatDuration(episode.duration)}</span>
                {episode.season_number > 0 && (
                  <span>S{episode.season_number} E{episode.episode_number}</span>
                )}
              </div>

              {progress && progress.current_time > 0 && (
                <div className="mt-4 p-4 bg-indigo-50 rounded-lg">
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-indigo-600">播放进度</span>
                    <span className="text-indigo-600">
                      {formatTime(progress.current_time)} / {formatDuration(episode.duration)}
                    </span>
                  </div>
                  <div className="mt-2 h-2 bg-indigo-200 rounded-full overflow-hidden">
                    <div
                      className="h-full bg-indigo-600"
                      style={{ width: `${(progress.current_time / episode.duration) * 100}%` }}
                    />
                  </div>
                </div>
              )}

              <div className="mt-6 text-gray-600 whitespace-pre-wrap">
                {episode.description}
              </div>
            </div>
          </div>
        </div>
      )}

      <div className="flex items-center justify-between">
        <h3 className="text-xl font-semibold text-gray-900">笔记</h3>
        {notes && notes.length > 0 && (
          <button
            onClick={handleExportNotes}
            className="btn btn-secondary flex items-center gap-2 text-sm"
          >
            <Download className="w-4 h-4" />
            导出 Markdown
          </button>
        )}
      </div>

      <div className="card p-4">
        <form onSubmit={handleSubmit(onSubmitNote)} className="space-y-4">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <label className="text-sm text-gray-600">时间戳:</label>
              <input
                type="number"
                step="0.1"
                {...register('timestamp')}
                className="input w-32 text-sm"
                placeholder="0"
              />
              <button
                type="button"
                onClick={() => setValue('timestamp', currentTime)}
                className="text-sm text-indigo-600 hover:underline"
              >
                使用当前播放时间
              </button>
            </div>
          </div>

          <textarea
            {...register('content')}
            rows={3}
            placeholder="添加笔记..."
            className="input resize-none"
          />
          {errors.content && (
            <p className="text-red-500 text-sm">{errors.content.message}</p>
          )}

          <div>
            <div className="flex items-center gap-2 mb-2">
              <Tag className="w-4 h-4 text-gray-400" />
              <input
                type="text"
                value={tagInput}
                onChange={(e) => setTagInput(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), handleAddTag())}
                placeholder="添加标签，按回车确认"
                className="input flex-1 text-sm"
              />
              <button type="button" onClick={handleAddTag} className="btn btn-secondary text-sm">
                添加
              </button>
            </div>
            <div className="flex flex-wrap gap-2">
              {(getValues('tags') || []).filter(Boolean).map((tag) => (
                <span
                  key={tag}
                  className="badge bg-indigo-100 text-indigo-700 flex items-center gap-1"
                >
                  {tag}
                  <button
                    type="button"
                    onClick={() => handleRemoveTag(tag as string)}
                    className="hover:text-indigo-900"
                  >
                    ×
                  </button>
                </span>
              ))}
            </div>
          </div>

          <div className="flex gap-2 justify-end">
            {editingNote && (
              <button
                type="button"
                onClick={() => {
                  setEditingNote(null)
                  reset()
                  setTagInput('')
                }}
                className="btn btn-secondary"
              >
                取消
              </button>
            )}
            <button type="submit" className="btn btn-primary">
              {editingNote ? '更新笔记' : '添加笔记'}
            </button>
          </div>
        </form>
      </div>

      {notes && notes.length > 0 ? (
        <div className="space-y-3">
          {notes.map((note) => (
            <div key={note.id} className="card p-4">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <span className="badge bg-indigo-100 text-indigo-700">
                      {formatTime(note.timestamp)}
                    </span>
                    {note.tags?.map((tag) => (
                      <span key={tag} className="badge bg-gray-100 text-gray-600">
                        {tag}
                      </span>
                    ))}
                  </div>
                  <p className="text-gray-700 whitespace-pre-wrap">{note.content}</p>
                </div>
                <div className="flex gap-2 ml-4">
                  <button
                    onClick={() => handleEditNote(note)}
                    className="p-2 hover:bg-gray-100 rounded-lg"
                  >
                    <Edit className="w-4 h-4" />
                  </button>
                  <button
                    onClick={() => handleDeleteNote(note.id)}
                    className="p-2 hover:bg-red-50 text-red-500 rounded-lg"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="text-center py-8 card">
          <p className="text-gray-500">暂无笔记</p>
        </div>
      )}

      {showAddToPlaylist && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 w-full max-w-md mx-4">
            <h3 className="text-lg font-semibold mb-4">添加到播放列表</h3>
            {playlists && playlists.length > 0 ? (
              <div className="space-y-2 max-h-64 overflow-y-auto">
                {playlists.map((playlist) => (
                  <button
                    key={playlist.id}
                    onClick={() => handleAddToPlaylist(playlist.id)}
                    className="w-full text-left p-3 hover:bg-gray-50 rounded-lg transition-colors"
                  >
                    {playlist.name}
                  </button>
                ))}
              </div>
            ) : (
              <p className="text-gray-500 text-center py-4">暂无播放列表</p>
            )}
            <button
              onClick={() => setShowAddToPlaylist(false)}
              className="mt-4 w-full btn btn-secondary"
            >
              取消
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
