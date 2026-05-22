import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Search, Play } from 'lucide-react'
import { useDispatch } from 'react-redux'
import { useSearchNotesQuery } from '@/app/services/api'
import { setCurrentEpisode, setPlaying } from '@/features/player/playerSlice'
import { formatDateTime, formatTime } from '@/utils/format'

export default function NotesPage() {
  const dispatch = useDispatch()
  const [searchQuery, setSearchQuery] = useState('')

  const { data: notes, isLoading } = useSearchNotesQuery(searchQuery, {
    skip: !searchQuery,
  })

  const handlePlay = (episode: any, timestamp: number) => {
    dispatch(setCurrentEpisode(episode))
    dispatch(setPlaying(true))
    setTimeout(() => {
      const audio = document.querySelector('audio')
      if (audio) {
        audio.currentTime = timestamp
      }
    }, 500)
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">笔记搜索</h1>
      </div>

      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
        <input
          type="text"
          placeholder="搜索笔记内容..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="input pl-10"
        />
      </div>

      {searchQuery && isLoading && (
        <div className="text-center py-12">搜索中...</div>
      )}

      {searchQuery && !isLoading && notes && notes.length > 0 ? (
        <div className="space-y-3">
          {notes.map((note) => (
            <div key={note.id} className="card p-4">
              <div className="flex items-start gap-4">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-3 mb-2">
                    <Link
                      to={`/podcasts/${note.episode?.podcast_id}`}
                      className="text-sm text-indigo-600 hover:underline"
                    >
                      {note.episode?.podcast?.title}
                    </Link>
                    <span className="text-sm text-gray-500">›</span>
                    <Link
                      to={`/episodes/${note.episode_id}`}
                      className="text-sm font-medium text-gray-900 hover:underline"
                    >
                      {note.episode?.title}
                    </Link>
                  </div>
                  <div className="flex items-center gap-3 mb-2">
                    <button
                      onClick={() => handlePlay(note.episode, note.timestamp)}
                      className="inline-flex items-center gap-1 text-sm text-indigo-600 hover:underline"
                    >
                      <Play className="w-3 h-3" />
                      {formatTime(note.timestamp)}
                    </button>
                    {note.tags?.map((tag) => (
                      <span key={tag} className="badge bg-gray-100 text-gray-600 text-xs">
                        {tag}
                      </span>
                    ))}
                  </div>
                  <p className="text-gray-700">{note.content}</p>
                  <p className="text-xs text-gray-400 mt-2">
                    {formatDateTime(note.created_at)}
                  </p>
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : searchQuery && !isLoading ? (
        <div className="text-center py-12 card">
          <p className="text-gray-500">未找到匹配的笔记</p>
        </div>
      ) : null}

      {!searchQuery && (
        <div className="text-center py-12 card">
          <p className="text-gray-500">输入关键词搜索笔记</p>
        </div>
      )}
    </div>
  )
}
