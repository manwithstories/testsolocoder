import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useBookStore } from '../store/useBookStore';
import { useSettingsStore } from '../store/useSettingsStore';
import { BookCard } from '../components/books/BookCard';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { Select } from '../components/ui/Select';
import { Badge } from '../components/ui/Badge';
import Empty from '../components/Empty';
import {
  Plus,
  Search,
  Grid,
  List,
  Heart,
  BookMarked,
  Filter,
  X,
} from 'lucide-react';
import { BOOK_STATUS, STATUS_COLORS } from '../constants/categories';
import type { BookStatus } from '../types/book';
import type { ViewMode } from '../types/settings';

const Shelf: React.FC = () => {
  const navigate = useNavigate();
  const {
    books,
    searchQuery,
    filterStatus,
    filterCategory,
    showFavoritesOnly,
    showReadingQueueOnly,
    sortBy,
    sortOrder,
    setSearchQuery,
    setFilterStatus,
    setFilterCategory,
    setSortBy,
    setSortOrder,
    toggleFavoritesOnly,
    toggleReadingQueueOnly,
    getFilteredBooks,
    getCategories,
  } = useBookStore();
  const { defaultView } = useSettingsStore();
  const [viewMode, setViewMode] = useState<ViewMode>(defaultView);

  const filteredBooks = getFilteredBooks();
  const categories = getCategories();

  const hasActiveFilters =
    filterStatus !== 'all' ||
    filterCategory !== null ||
    showFavoritesOnly ||
    showReadingQueueOnly ||
    searchQuery !== '';

  const clearFilters = () => {
    setSearchQuery('');
    setFilterStatus('all');
    setFilterCategory(null);
    if (showFavoritesOnly) toggleFavoritesOnly();
    if (showReadingQueueOnly) toggleReadingQueueOnly();
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">我的书架</h1>
          <p className="text-gray-500 mt-1">
            共 {books.length} 本书，{filteredBooks.length} 本符合筛选条件
          </p>
        </div>
        <Button onClick={() => navigate('/add-book')}>
          <Plus className="w-4 h-4 mr-2" />
          添加书籍
        </Button>
      </div>

      <div className="bg-white rounded-xl shadow-sm border p-4 space-y-4">
        <div className="flex flex-wrap gap-4 items-center">
          <div className="flex-1 min-w-[200px] max-w-md">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <Input
                placeholder="搜索书名或作者..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-9"
              />
            </div>
          </div>

          <div className="flex items-center gap-2">
            <span className="text-sm text-gray-500">视图:</span>
            <button
              onClick={() => setViewMode('grid')}
              className={`p-2 rounded-lg transition-colors ${
                viewMode === 'grid'
                  ? 'bg-blue-100 text-blue-600'
                  : 'text-gray-400 hover:bg-gray-100'
              }`}
            >
              <Grid className="w-4 h-4" />
            </button>
            <button
              onClick={() => setViewMode('list')}
              className={`p-2 rounded-lg transition-colors ${
                viewMode === 'list'
                  ? 'bg-blue-100 text-blue-600'
                  : 'text-gray-400 hover:bg-gray-100'
              }`}
            >
              <List className="w-4 h-4" />
            </button>
          </div>
        </div>

        <div className="flex flex-wrap gap-3 items-center">
          <div className="flex items-center gap-2">
            <Filter className="w-4 h-4 text-gray-400" />
            <span className="text-sm text-gray-500">筛选:</span>
          </div>

          <button
            onClick={toggleFavoritesOnly}
            className={`flex items-center gap-1 px-3 py-1.5 rounded-full text-sm transition-colors ${
              showFavoritesOnly
                ? 'bg-red-100 text-red-600'
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
            }`}
          >
            <Heart className="w-3.5 h-3.5" />
            收藏夹
          </button>

          <button
            onClick={toggleReadingQueueOnly}
            className={`flex items-center gap-1 px-3 py-1.5 rounded-full text-sm transition-colors ${
              showReadingQueueOnly
                ? 'bg-blue-100 text-blue-600'
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
            }`}
          >
            <BookMarked className="w-3.5 h-3.5" />
            阅读队列
          </button>

          <Select
            value={filterStatus}
            onChange={(e) => setFilterStatus(e.target.value as BookStatus | 'all')}
            className="w-32"
            options={[
              { value: 'all', label: '全部状态' },
              ...BOOK_STATUS.map((status) => ({ value: status.value, label: status.label })),
            ]}
          />

          {categories.length > 0 && (
            <Select
              value={filterCategory || ''}
              onChange={(e) => setFilterCategory(e.target.value || null)}
              className="w-32"
              options={[
                { value: '', label: '全部分类' },
                ...categories.map((cat) => ({ value: cat, label: cat })),
              ]}
            />
          )}

          <Select
            value={`${sortBy}-${sortOrder}`}
            onChange={(e) => {
              const [by, order] = e.target.value.split('-');
              setSortBy(by as any);
              setSortOrder(order as 'asc' | 'desc');
            }}
            className="w-40"
            options={[
              { value: 'createdAt-desc', label: '最新添加' },
              { value: 'createdAt-asc', label: '最早添加' },
              { value: 'title-asc', label: '书名 A-Z' },
              { value: 'title-desc', label: '书名 Z-A' },
              { value: 'author-asc', label: '作者 A-Z' },
              { value: 'author-desc', label: '作者 Z-A' },
              { value: 'currentPage-desc', label: '阅读进度' },
            ]}
          />

          {hasActiveFilters && (
            <button
              onClick={clearFilters}
              className="flex items-center gap-1 px-3 py-1.5 text-sm text-gray-500 hover:text-gray-700 transition-colors"
            >
              <X className="w-3.5 h-3.5" />
              清除筛选
            </button>
          )}
        </div>

        {filterStatus !== 'all' && (
          <div className="pt-2 border-t flex gap-2 flex-wrap">
            <span className="text-sm text-gray-500">状态:</span>
            <Badge variant={STATUS_COLORS[filterStatus as BookStatus] as 'default' | 'success' | 'warning' | 'danger' | 'info'}>
              {BOOK_STATUS.find((s) => s.value === filterStatus)?.label}
            </Badge>
          </div>
        )}
      </div>

      {filteredBooks.length === 0 ? (
        <Empty
          title={hasActiveFilters ? '没有符合条件的书籍' : '书架还是空的'}
          description={
            hasActiveFilters
              ? '尝试调整筛选条件'
              : '点击右上角的"添加书籍"按钮开始添加你的第一本书'
          }
          action={
            !hasActiveFilters && (
              <Button onClick={() => navigate('/add-book')}>
                <Plus className="w-4 h-4 mr-2" />
                添加书籍
              </Button>
            )
          }
        />
      ) : (
        <div
          className={
            viewMode === 'grid'
              ? 'grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6'
              : 'space-y-4'
          }
        >
          {filteredBooks.map((book) => (
            <BookCard
              key={book.id}
              book={book}
              viewMode={viewMode}
            />
          ))}
        </div>
      )}
    </div>
  );
};

export default Shelf;
