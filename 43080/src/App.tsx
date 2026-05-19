import { useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { Layout } from './components/layout/Layout';
import { ToastProvider } from './components/ui/Toast';
import { useBookStore } from './store/useBookStore';
import { useReadingStore } from './store/useReadingStore';
import { useNoteStore } from './store/useNoteStore';
import { useSettingsStore } from './store/useSettingsStore';
import Dashboard from './pages/Dashboard';
import Shelf from './pages/Shelf';
import BookDetail from './pages/BookDetail';
import AddBook from './pages/AddBook';
import Notes from './pages/Notes';
import Settings from './pages/Settings';

function App() {
  const loadBooks = useBookStore((state) => state.loadBooks);
  const loadReadingData = useReadingStore((state) => state.loadReadingData);
  const loadNotes = useNoteStore((state) => state.loadNotes);
  const loadSettings = useSettingsStore((state) => state.loadSettings);

  useEffect(() => {
    loadSettings();
    loadBooks();
    loadReadingData();
    loadNotes();
  }, [loadSettings, loadBooks, loadReadingData, loadNotes]);

  return (
    <ToastProvider>
      <Router>
        <Layout>
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/dashboard" element={<Navigate to="/" replace />} />
            <Route path="/shelf" element={<Shelf />} />
            <Route path="/books/:id" element={<BookDetail />} />
            <Route path="/add-book" element={<AddBook />} />
            <Route path="/edit-book/:id" element={<AddBook />} />
            <Route path="/notes" element={<Notes />} />
            <Route path="/settings" element={<Settings />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </Layout>
      </Router>
    </ToastProvider>
  );
}

export default App;
