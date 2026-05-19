import React, { useState } from 'react';
import { useNoteStore } from '../store/useNoteStore';
import { useBookStore } from '../store/useBookStore';
import { NoteCard } from '../components/notes/NoteCard';
import { NoteForm } from '../components/notes/NoteForm';
import Empty from '../components/Empty';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { Select } from '../components/ui/Select';
import { Search, Plus, FileText, Highlighter, Filter } from 'lucide-react';
import type { NoteType } from '../types/note';

const Notes: React.FC = () => {
  const {
    notes,
    filterBookId,
    filterType,
    searchQuery,
    setFilterBookId,
    setFilterType,
    setSearchQuery,
    getFilteredNotes,
  } = useNoteStore();
  const { books } = useBookStore();
  const [showNoteForm, setShowNoteForm] = useState(false);

  const filteredNotes = getFilteredNotes();

  const hasActiveFilters =
    filterBookId !== null || filterType !== 'all' || searchQuery !== '';

  const clearFilters = () => {
    setFilterBookId(null);
    setFilterType('all');
    setSearchQuery('');
  };

  const noteTypes: { value: NoteType | 'all'; label: string; icon: React.ReactNode }[] = [
    { value: 'all', label: '全部', icon: <Filter className="w-4 h-4" /> },
    { value: 'note', label: '笔记', icon: <FileText className="w-4 h-4" /> },
    { value: 'highlight', label: '高亮', icon: <Highlighter className="w-4 h-4" /> },
  ];

  const getBookTitle = (bookId: string): string => {
    return books.find((b) => b.id === bookId)?.title || '未知书籍';
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">笔记与高亮</h1>
          <p className="text-gray-500 mt-1">
            共 {notes.length} 条笔记，{filteredNotes.length} 条符合筛选条件
          </p>
        </div>
        {books.length > 0 && (
          <Button onClick={() => setShowNoteForm(true)}>
            <Plus className="w-4 h-4 mr-2" />
            添加笔记
          </Button>
        )}
      </div>

      <div className="bg-white rounded-xl shadow-sm border p-4 space-y-4">
        <div className="flex flex-wrap gap-4 items-center">
          <div className="flex-1 min-w-[200px] max-w-md">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <Input
                placeholder="搜索笔记内容..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-9"
              />
            </div>
          </div>

          <div className="flex items-center gap-2">
            {noteTypes.map((type) => (
              <button
                key={type.value}
                onClick={() => setFilterType(type.value)}
                className={`flex items-center gap-1.5 px-3 py-2 rounded-lg text-sm transition-colors ${
                  filterType === type.value
                    ? 'bg-blue-100 text-blue-600'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                {type.icon}
                {type.label}
              </button>
            ))}
          </div>
        </div>

        <div className="flex flex-wrap gap-3 items-center">
          <span className="text-sm text-gray-500">筛选书籍:</span>
          <Select
            value={filterBookId || ''}
            onChange={(e) => setFilterBookId(e.target.value || null)}
            className="w-48"
            options={[
              { value: '', label: '全部书籍' },
              ...books.map((book) => ({ value: book.id, label: book.title })),
            ]}
          />

          {hasActiveFilters && (
            <button
              onClick={clearFilters}
              className="text-sm text-gray-500 hover:text-gray-700 transition-colors"
            >
              清除筛选
            </button>
          )}
        </div>
      </div>

      {books.length === 0 ? (
        <Empty
          title="还没有添加书籍"
          description="先添加书籍，然后开始记录你的阅读笔记"
        />
      ) : filteredNotes.length === 0 ? (
        <Empty
          title={hasActiveFilters ? '没有符合条件的笔记' : '还没有笔记'}
          description={
            hasActiveFilters
              ? '尝试调整筛选条件'
              : '开始阅读并记录你的想法和心得'
          }
          action={
            !hasActiveFilters && (
              <Button onClick={() => setShowNoteForm(true)}>
                <Plus className="w-4 h-4 mr-2" />
                添加第一条笔记
              </Button>
            )
          }
        />
      ) : (
        <div className="space-y-3">
          {filteredNotes.map((note) => (
            <div key={note.id} className="relative">
              <div className="absolute -left-2 top-1/2 -translate-y-1/2 text-xs text-gray-400 bg-white px-1">
                {getBookTitle(note.bookId)}
              </div>
              <NoteCard note={note} showBookTitle />
            </div>
          ))}
        </div>
      )}

      {showNoteForm && (
        <NoteForm
          onClose={() => setShowNoteForm(false)}
        />
      )}
    </div>
  );
};

export default Notes;
