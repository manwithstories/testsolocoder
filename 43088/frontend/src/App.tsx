import { Routes, Route } from 'react-router-dom'
import Layout from '@/components/layout/Layout'
import PodcastList from '@/features/podcasts/PodcastList'
import PodcastDetail from '@/features/podcasts/PodcastDetail'
import EpisodeDetail from '@/features/episodes/EpisodeDetail'
import HistoryPage from '@/features/history/HistoryPage'
import NotesPage from '@/features/notes/NotesPage'
import PlaylistList from '@/features/playlists/PlaylistList'
import PlaylistDetail from '@/features/playlists/PlaylistDetail'
import StatsPage from '@/features/stats/StatsPage'
import ImportExportPage from '@/features/settings/ImportExportPage'
import AudioPlayer from '@/features/player/AudioPlayer'

function App() {
  return (
    <div className="min-h-screen bg-gray-50 pb-24">
      <Layout>
        <Routes>
          <Route path="/" element={<PodcastList />} />
          <Route path="/podcasts/:id" element={<PodcastDetail />} />
          <Route path="/episodes/:id" element={<EpisodeDetail />} />
          <Route path="/history" element={<HistoryPage />} />
          <Route path="/notes" element={<NotesPage />} />
          <Route path="/playlists" element={<PlaylistList />} />
          <Route path="/playlists/:id" element={<PlaylistDetail />} />
          <Route path="/stats" element={<StatsPage />} />
          <Route path="/import-export" element={<ImportExportPage />} />
        </Routes>
      </Layout>
      <AudioPlayer />
    </div>
  )
}

export default App
