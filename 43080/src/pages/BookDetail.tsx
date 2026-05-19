import React, { useState, useMemo } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useBookStore } from '../store/useBookStore';
import { useReadingStore } from '../store/useReadingStore';
import { useNoteStore } from '../store/useNoteStore';
import { useToast } from '../components/ui/Toast';
import { useConfirm } from '../hooks/useConfirm';
import { ConfirmDialog } from '../components/ui/ConfirmDialog';
import { Button } from '../components/ui/Button';
import { Badge } from '../components/ui/Badge';
import { ProgressBar } from '../components/ui/ProgressBar';
import { Tabs } from '../components/ui/Tabs';
import { ReadingPlanForm } from '../components/books/ReadingPlanForm';
import { ReadingSessionForm } from '../components/books/ReadingSessionForm';
import { NoteForm } from '../components/notes/NoteForm';
import { BookmarkForm } from '../components/notes/BookmarkForm';
import { NoteCard } from '../components/notes/NoteCard';
import Empty from '../components/Empty';
import {
  ArrowLeft,
  Heart,
  BookMarked,
  Edit,
  Trash2,
  Book,
  Calendar,
  Clock,
  Target,
  FileText,
  Bookmark,
  Highlighter,
  Plus,
  CheckCircle2,
} from 'lucide-react';
import { BOOK_STATUS, STATUS_COLORS } from '../constants/categories';
import { formatDate } from '../utils/date';
import type { BookStatus as BookStatusType } from '../types/book';

const BookDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { showToast } = useToast();
  const { confirmState, confirm, cancel } = useConfirm();
  const { getBookById, toggleFavorite, toggleReadingQueue, updateBookStatus, deleteBook } = useBookStore();
  const { getPlanByBookId, getSessionsByBookId, calculateProgress, calculateEstimatedDays } = useReadingStore();
  const { getNotesByBookId, getBookmarksByBookId, getHighlightsByBookId, deleteBookmark } = useNoteStore();

  const [showPlanForm, setShowPlanForm] = useState(false);
  const [showSessionForm, setShowSessionForm] = useState(false);
  const [showNoteForm, setShowNoteForm] = useState(false);
  const [showBookmarkForm, setShowBookmarkForm] = useState(false);

  const book = id ? getBookById(id) : undefined;
  const plan = id ? getPlanByBookId(id) : undefined;
  const sessions = id ? getSessionsByBookId(id) : [];
  const notes = id ? getNotesByBookId(id) : [];
  const bookmarks = id ? getBookmarksByBookId(id) : [];
  const highlights = id ? getHighlightsByBookId(id) : [];
  const progress = id ? calculateProgress(id) : 0;
  const estimatedDays = id ? calculateEstimatedDays(id) : null;

  const totalReadingTime = useMemo(
    () => sessions.reduce((sum, s) => sum + s.duration, 0),
    [sessions]
  );

  const totalPagesRead = useMemo(
    () => sessions.reduce((sum, s) => sum + (s.endPage - s.startPage), 0),
    [sessions]
  );

  if (!book) {
    return (
      <Empty
        title="书籍不存在"
        description="你访问的书籍可能已被删除或不存在"
        action={
          <Button onClick={() => navigate('/shelf')}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            返回书架
          </Button>
        }
      />
    );
  }

  const handleDelete = () => {
    confirm(
      '删除书籍',
      `确定要删除《${book.title}》吗？相关的阅读记录、笔记和书签也会被删除。此操作不可撤销。`,
      () => {
        deleteBook(book.id);
        showToast('success', '书籍已删除');
        navigate('/shelf');
      }
    );
  };

  const handleStatusChange = (status: BookStatusType) => {
    updateBookStatus(book.id, status);
    showToast('success', `已标记为${BOOK_STATUS.find((s) => s.value === status)?.label}`);
  };

  const formatDuration = (minutes: number): string => {
    if (minutes < 60) return `${minutes} 分钟`;
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return mins > 0 ? `${hours} 小时 ${mins} 分钟` : `${hours} 小时`;
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <button
          onClick={() => navigate('/shelf')}
          className="p-2 rounded-lg hover:bg-gray-100 transition-colors"
        >
          <ArrowLeft className="w-5 h-5 text-gray-600" />
        </button>
        <div className="flex-1">
          <h1 className="text-2xl font-bold text-gray-900">{book.title}</h1>
          <p className="text-gray-500">{book.author}</p>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => toggleFavorite(book.id)}
            className={`p-2 rounded-lg transition-colors ${
              book.isFavorite
                ? 'text-red-500 bg-red-50'
                : 'text-gray-400 hover:bg-gray-100'
            }`}
            title={book.isFavorite ? '取消收藏' : '添加收藏'}
          >
            <Heart className={`w-5 h-5 ${book.isFavorite ? 'fill-current' : ''}`} />
          </button>
          <button
            onClick={() => toggleReadingQueue(book.id)}
            className={`p-2 rounded-lg transition-colors ${
              book.inReadingQueue
                ? 'text-blue-500 bg-blue-50'
                : 'text-gray-400 hover:bg-gray-100'
            }`}
            title={book.inReadingQueue ? '移出阅读队列' : '加入阅读队列'}
          >
            <BookMarked className="w-5 h-5" />
          </button>
          <Button variant="outline" size="sm" onClick={() => navigate(`/edit-book/${book.id}`)}>
            <Edit className="w-4 h-4 mr-1" />
            编辑
          </Button>
          <Button variant="ghost" size="sm" onClick={handleDelete} className="text-red-500 hover:text-red-600">
            <Trash2 className="w-4 h-4" />
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-1">
          <div className="bg-white rounded-xl shadow-sm border p-6 sticky top-6">
            <div className="flex gap-6">
              {book.coverUrl ? (
                <img
                  src={book.coverUrl}
                  alt={book.title}
                  className="w-32 h-48 object-cover rounded-lg shadow-md"
                />
              ) : (
                <div className="w-32 h-48 bg-gray-200 rounded-lg flex items-center justify-center">
                  <Book className="w-12 h-12 text-gray-400" />
                </div>
              )}
              <div className="flex-1 space-y-3">
                <div>
                  <p className="text-sm text-gray-500">阅读状态</p>
                  <Badge color={STATUS_COLORS[book.status]} className="mt-1">
                    {BOOK_STATUS.find((s) => s.value === book.status)?.label}
                  </Badge>
                </div>
                <div>
                  <p className="text-sm text-gray-500">总页数</p>
                  <p className="font-medium">{book.totalPages} 页</p>
                </div>
                <div>
                  <p className="text-sm text-gray-500">已读</p>
                  <p className="font-medium">{book.currentPage} 页</p>
                </div>
              </div>
            </div>

            <div className="mt-6">
              <div className="flex justify-between text-sm mb-2">
                <span className="text-gray-500">阅读进度</span>
                <span className="font-medium">{progress}%</span>
              </div>
              <ProgressBar value={progress} className="h-3" />
            </div>

            {book.categories.length > 0 && (
              <div className="mt-6">
                <p className="text-sm text-gray-500 mb-2">分类标签</p>
                <div className="flex flex-wrap gap-2">
                  {book.categories.map((cat) => (
                    <Badge key={cat} variant="default">
                      {cat}
                    </Badge>
                  ))}
                </div>
              </div>
            )}

            {plan && (
              <div className="mt-6 pt-6 border-t">
                <h3 className="font-medium text-gray-900 mb-3 flex items-center gap-2">
                  <Target className="w-4 h-4 text-blue-500" />
                  阅读计划
                </h3>
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-gray-500">每日计划</span>
                    <span className="font-medium">{plan.dailyPages} 页</span>
                  </div>
                  {plan.targetDate && (
                    <div className="flex justify-between">
                      <span className="text-gray-500">目标日期</span>
                      <span className="font-medium">{formatDate(plan.targetDate)}</span>
                    </div>
                  )}
                  {estimatedDays !== null && estimatedDays > 0 && (
                    <div className="flex justify-between">
                      <span className="text-gray-500">预计还需</span>
                      <span className="font-medium">{estimatedDays} 天</span>
                    </div>
                  )}
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  className="w-full mt-3"
                  onClick={() => setShowPlanForm(true)}
                >
                  修改计划
                </Button>
              </div>
            )}

            <div className="mt-6 space-y-2">
              {!plan && (
                <Button className="w-full" onClick={() => setShowPlanForm(true)}>
                  <Target className="w-4 h-4 mr-2" />
                  设置阅读计划
                </Button>
              )}
              <Button variant="outline" className="w-full" onClick={() => setShowSessionForm(true)}>
                <Plus className="w-4 h-4 mr-2" />
                记录阅读
              </Button>
            </div>
          </div>
        </div>

        <div className="lg:col-span-2">
          <Tabs
            tabs={[
              {
                id: 'sessions',
                label: `阅读记录 (${sessions.length})`,
                content: (
                  <div className="space-y-4">
                    <div className="grid grid-cols-3 gap-4">
                      <div className="bg-white rounded-xl shadow-sm border p-4">
                        <div className="flex items-center gap-2 text-gray-500 mb-1">
                          <Clock className="w-4 h-4" />
                          <span className="text-sm">总阅读时长</span>
                        </div>
                        <p className="text-2xl font-bold text-gray-900">{formatDuration(totalReadingTime)}</p>
                      </div>
                      <div className="bg-white rounded-xl shadow-sm border p-4">
                        <div className="flex items-center gap-2 text-gray-500 mb-1">
                          <FileText className="w-4 h-4" />
                          <span className="text-sm">已读页数</span>
                        </div>
                        <p className="text-2xl font-bold text-gray-900">{totalPagesRead} 页</p>
                      </div>
                      <div className="bg-white rounded-xl shadow-sm border p-4">
                        <div className="flex items-center gap-2 text-gray-500 mb-1">
                          <Calendar className="w-4 h-4" />
                          <span className="text-sm">阅读次数</span>
                        </div>
                        <p className="text-2xl font-bold text-gray-900">{sessions.length} 次</p>
                      </div>
                    </div>

                    {book.status !== 'completed' && (
                      <div className="bg-white rounded-xl shadow-sm border p-4">
                        <h3 className="font-medium text-gray-900 mb-3">快速更新状态</h3>
                        <div className="flex flex-wrap gap-2">
                          {BOOK_STATUS.filter((s) => s.value !== book.status).map((status) => (
                            <Button
                              key={status.value}
                              variant="outline"
                              size="sm"
                              onClick={() => handleStatusChange(status.value as BookStatusType)}
                            >
                              {status.value === 'completed' && (
                                <CheckCircle2 className="w-4 h-4 mr-1 text-green-500" />
                              )}
                              标记为{status.label}
                            </Button>
                          ))}
                        </div>
                      </div>
                    )}

                    {sessions.length === 0 ? (
                      <Empty
                        title="还没有阅读记录"
                        description="开始阅读并记录你的阅读进度吧"
                        action={
                          <Button onClick={() => setShowSessionForm(true)}>
                            <Plus className="w-4 h-4 mr-2" />
                            添加第一条记录
                          </Button>
                        }
                      />
                    ) : (
                      <div className="space-y-3">
                        {sessions.map((session) => (
                          <div
                            key={session.id}
                            className="bg-white rounded-xl shadow-sm border p-4 hover:shadow-md transition-shadow"
                          >
                            <div className="flex justify-between items-start mb-2">
                              <div>
                                <p className="font-medium text-gray-900">
                                  第 {session.startPage} - {session.endPage} 页
                                </p>
                                <p className="text-sm text-gray-500">
                                  阅读了 {session.endPage - session.startPage} 页 · {formatDuration(session.duration)}
                                </p>
                              </div>
                              <p className="text-sm text-gray-400">{formatDate(session.date)}</p>
                            </div>
                            {session.notes && (
                              <p className="text-sm text-gray-600 mt-2 pt-2 border-t">{session.notes}</p>
                            )}
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                ),
              },
              {
                id: 'notes',
                label: `笔记 (${notes.length})`,
                content: (
                  <div className="space-y-4">
                    <div className="flex justify-between items-center">
                      <h3 className="font-medium text-gray-900">章节笔记</h3>
                      <Button size="sm" onClick={() => setShowNoteForm(true)}>
                        <Plus className="w-4 h-4 mr-1" />
                        添加笔记
                      </Button>
                    </div>
                    {notes.length === 0 ? (
                      <Empty
                        title="还没有笔记"
                        description="记录你的阅读心得和想法"
                        action={
                          <Button onClick={() => setShowNoteForm(true)}>
                            <Plus className="w-4 h-4 mr-2" />
                            添加第一条笔记
                          </Button>
                        }
                      />
                    ) : (
                      <div className="space-y-3">
                        {notes.map((note) => (
                          <NoteCard key={note.id} note={note} />
                        ))}
                      </div>
                    )}
                  </div>
                ),
              },
              {
                id: 'highlights',
                label: `高亮 (${highlights.length})`,
                content: (
                  <div className="space-y-4">
                    <h3 className="font-medium text-gray-900">文本高亮</h3>
                    {highlights.length === 0 ? (
                      <Empty
                        title="还没有高亮"
                        description="标记书中重要的段落和句子"
                      />
                    ) : (
                      <div className="space-y-3">
                        {highlights.map((highlight) => (
                          <NoteCard key={highlight.id} note={highlight} />
                        ))}
                      </div>
                    )}
                  </div>
                ),
              },
              {
                id: 'bookmarks',
                label: `书签 (${bookmarks.length})`,
                content: (
                  <div className="space-y-4">
                    <div className="flex justify-between items-center">
                      <h3 className="font-medium text-gray-900">书签</h3>
                      <Button size="sm" onClick={() => setShowBookmarkForm(true)}>
                        <Plus className="w-4 h-4 mr-1" />
                        添加书签
                      </Button>
                    </div>
                    {bookmarks.length === 0 ? (
                      <Empty
                        title="还没有书签"
                        description="添加书签到重要的页面"
                        action={
                          <Button onClick={() => setShowBookmarkForm(true)}>
                            <Plus className="w-4 h-4 mr-2" />
                            添加第一个书签
                          </Button>
                        }
                      />
                    ) : (
                      <div className="space-y-3">
                        {bookmarks.map((bookmark) => (
                          <div
                            key={bookmark.id}
                            className="bg-white rounded-xl shadow-sm border p-4 flex items-start gap-3"
                          >
                            <Bookmark className="w-5 h-5 text-blue-500 mt-0.5 flex-shrink-0" />
                            <div className="flex-1 min-w-0">
                              <div className="flex justify-between items-start gap-3">
                                <div className="min-w-0">
                                  <p className="font-medium text-gray-900">
                                    第 {bookmark.page} 页
                                    {bookmark.label && (
                                      <span className="text-gray-500 font-normal ml-2">
                                        — {bookmark.label}
                                      </span>
                                    )}
                                  </p>
                                  {bookmark.description && (
                                    <p className="text-sm text-gray-500 mt-1">{bookmark.description}</p>
                                  )}
                                </div>
                                <div className="flex items-center gap-2 flex-shrink-0">
                                  <p className="text-xs text-gray-400">
                                    {formatDate(bookmark.createdAt)}
                                  </p>
                                  <button
                                    onClick={() => {
                                      deleteBookmark(bookmark.id);
                                      showToast('success', '书签已删除');
                                    }}
                                    className="p-1 text-gray-400 hover:text-red-500 transition-colors"
                                    title="删除书签"
                                  >
                                    <Trash2 className="w-4 h-4" />
                                  </button>
                                </div>
                              </div>
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                ),
              },
            ]}
          />
        </div>
      </div>

      {showPlanForm && (
        <ReadingPlanForm
          bookId={book.id}
          onClose={() => setShowPlanForm(false)}
        />
      )}

      {showSessionForm && (
        <ReadingSessionForm
          bookId={book.id}
          onClose={() => setShowSessionForm(false)}
        />
      )}

      {showNoteForm && (
        <NoteForm
          bookId={book.id}
          onClose={() => setShowNoteForm(false)}
        />
      )}

      {showBookmarkForm && (
        <BookmarkForm
          bookId={book.id}
          onClose={() => setShowBookmarkForm(false)}
        />
      )}

      <ConfirmDialog
        isOpen={confirmState.isOpen}
        title={confirmState.title}
        message={confirmState.message}
        onConfirm={confirmState.onConfirm}
        onClose={cancel}
      />
    </div>
  );
};

export default BookDetail;
