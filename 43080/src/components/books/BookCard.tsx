import React from 'react';
import { Link } from 'react-router-dom';
import { Heart, ListPlus, MoreVertical, BookOpen } from 'lucide-react';
import type { Book } from '../../types/book';
import { Card } from '../ui/Card';
import { Badge } from '../ui/Badge';
import { ProgressBar } from '../ui/ProgressBar';
import { Button } from '../ui/Button';
import { STATUS_LABELS, STATUS_COLORS } from '../../constants/categories';
import { useBookStore } from '../../store/useBookStore';
import { useReadingStore } from '../../store/useReadingStore';
import { useToast } from '../ui/Toast';

interface BookCardProps {
  book: Book;
  viewMode?: 'grid' | 'list';
}

export const BookCard: React.FC<BookCardProps> = ({ book, viewMode = 'grid' }) => {
  const { toggleFavorite, toggleReadingQueue, deleteBook } = useBookStore();
  const { calculateProgress } = useReadingStore();
  const { showToast } = useToast();
  const [showMenu, setShowMenu] = React.useState(false);

  const progress = calculateProgress(book.id);

  const handleDelete = () => {
    if (window.confirm('确定要删除这本书吗？相关的阅读记录和笔记也会被删除。')) {
      deleteBook(book.id);
      showToast('success', '书籍已删除');
    }
  };

  if (viewMode === 'list') {
    return (
      <Card hoverable className="p-4">
        <div className="flex items-center gap-4">
          <div className="w-16 h-24 bg-gray-200 rounded-lg overflow-hidden flex-shrink-0">
            {book.coverUrl ? (
              <img src={book.coverUrl} alt={book.title} className="w-full h-full object-cover" />
            ) : (
              <div className="w-full h-full flex items-center justify-center text-gray-400">
                <BookOpen className="w-8 h-8" />
              </div>
            )}
          </div>
          <div className="flex-1 min-w-0">
            <Link to={`/books/${book.id}`} className="hover:text-blue-600">
              <h3 className="font-semibold text-gray-900 truncate">{book.title}</h3>
            </Link>
            <p className="text-sm text-gray-500 truncate">{book.author}</p>
            <div className="flex items-center gap-2 mt-2">
              <Badge variant={book.status === 'completed' ? 'success' : book.status === 'reading' ? 'info' : 'default'} size="sm">
                {STATUS_LABELS[book.status]}
              </Badge>
              <span className="text-xs text-gray-400">
                {book.currentPage} / {book.totalPages} 页
              </span>
            </div>
            <ProgressBar value={progress} size="sm" className="mt-2" />
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={(e) => {
                e.preventDefault();
                toggleFavorite(book.id);
              }}
            >
              <Heart className={book.isFavorite ? 'fill-red-500 text-red-500' : ''} />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={(e) => {
                e.preventDefault();
                toggleReadingQueue(book.id);
              }}
            >
              <ListPlus className={book.inReadingQueue ? 'text-blue-600' : ''} />
            </Button>
            <div className="relative">
              <Button variant="ghost" size="sm" onClick={() => setShowMenu(!showMenu)}>
                <MoreVertical />
              </Button>
              {showMenu && (
                <div className="absolute right-0 top-full mt-1 bg-white border border-gray-200 rounded-lg shadow-lg py-1 z-10 min-w-[120px]">
                  <Link
                    to={`/books/${book.id}`}
                    className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-50"
                    onClick={() => setShowMenu(false)}
                  >
                    查看详情
                  </Link>
                  <button
                    onClick={() => {
                      handleDelete();
                      setShowMenu(false);
                    }}
                    className="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-gray-50"
                  >
                    删除
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      </Card>
    );
  }

  return (
    <Card hoverable className="group">
      <div className="relative">
        <div className="aspect-[3/4] bg-gray-100 overflow-hidden">
          {book.coverUrl ? (
            <img
              src={book.coverUrl}
              alt={book.title}
              className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center text-gray-400 bg-gradient-to-br from-gray-100 to-gray-200">
              <BookOpen className="w-16 h-16" />
            </div>
          )}
        </div>
        <div className="absolute top-2 right-2 flex gap-1">
          {book.isFavorite && (
            <div className="bg-white/90 rounded-full p-1.5">
              <Heart className="w-4 h-4 fill-red-500 text-red-500" />
            </div>
          )}
          {book.inReadingQueue && (
            <div className="bg-white/90 rounded-full p-1.5">
              <ListPlus className="w-4 h-4 text-blue-600" />
            </div>
          )}
        </div>
      </div>
      <div className="p-4">
        <Link to={`/books/${book.id}`}>
          <h3 className="font-semibold text-gray-900 truncate hover:text-blue-600 transition-colors">
            {book.title}
          </h3>
        </Link>
        <p className="text-sm text-gray-500 truncate mt-1">{book.author}</p>
        <div className="mt-3">
          <div className="flex items-center justify-between text-xs text-gray-500 mb-1">
            <span>阅读进度</span>
            <span className={STATUS_COLORS[book.status]}>{STATUS_LABELS[book.status]}</span>
          </div>
          <ProgressBar value={progress} size="sm" />
          <p className="text-xs text-gray-400 mt-1">
            {book.currentPage} / {book.totalPages} 页
          </p>
        </div>
      </div>
    </Card>
  );
};
