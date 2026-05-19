import React from 'react';
import { useReadingStore } from '../store/useReadingStore';
import { useBookStore } from '../store/useBookStore';
import { StatsCard } from '../components/dashboard/StatsCard';
import { ReadingChart } from '../components/dashboard/ReadingChart';
import { ReadingCalendar } from '../components/dashboard/ReadingCalendar';
import { Book, Clock, Target, TrendingUp, BookOpen } from 'lucide-react';

const Dashboard: React.FC = () => {
  const { getReadingStats } = useReadingStore();
  const { books } = useBookStore();
  const stats = getReadingStats();

  const formatDuration = (minutes: number): string => {
    if (minutes < 60) return `${minutes} 分钟`;
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return mins > 0 ? `${hours} 小时 ${mins} 分钟` : `${hours} 小时`;
  };

  const readingBooks = books.filter((b) => b.status === 'reading').length;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">阅读仪表盘</h1>
          <p className="text-gray-500 mt-1">追踪你的阅读进度和统计数据</p>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatsCard
          title="已完成书籍"
          value={stats.totalBooksRead.toString()}
          icon={<Book className="w-6 h-6" />}
          color="blue"
          description={`共 ${books.length} 本书`}
        />
        <StatsCard
          title="总阅读时长"
          value={formatDuration(stats.totalReadingTime)}
          icon={<Clock className="w-6 h-6" />}
          color="green"
          description={`${stats.totalPagesRead} 页`}
        />
        <StatsCard
          title="平均阅读速度"
          value={`${stats.averageReadingSpeed} 页/时`}
          icon={<TrendingUp className="w-6 h-6" />}
          color="purple"
          description="每小时阅读页数"
        />
        <StatsCard
          title="连续阅读"
          value={`${stats.currentStreak} 天`}
          icon={<Target className="w-6 h-6" />}
          color="orange"
          description={`最长 ${stats.longestStreak} 天`}
        />
      </div>

      {readingBooks > 0 && (
        <div className="bg-white rounded-xl shadow-sm border p-6">
          <div className="flex items-center gap-2 mb-4">
            <BookOpen className="w-5 h-5 text-blue-600" />
            <h2 className="text-xl font-semibold text-gray-900">正在阅读</h2>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {books
              .filter((b) => b.status === 'reading')
              .slice(0, 3)
              .map((book) => (
                <div key={book.id} className="flex gap-4 p-4 bg-gray-50 rounded-lg">
                  {book.coverUrl ? (
                    <img
                      src={book.coverUrl}
                      alt={book.title}
                      className="w-16 h-24 object-cover rounded shadow-sm"
                    />
                  ) : (
                    <div className="w-16 h-24 bg-gray-200 rounded flex items-center justify-center">
                      <Book className="w-8 h-8 text-gray-400" />
                    </div>
                  )}
                  <div className="flex-1 min-w-0">
                    <h3 className="font-medium text-gray-900 truncate">{book.title}</h3>
                    <p className="text-sm text-gray-500 truncate">{book.author}</p>
                    <div className="mt-2">
                      <div className="flex justify-between text-xs text-gray-500 mb-1">
                        <span>进度</span>
                        <span>{book.currentPage}/{book.totalPages} 页</span>
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2">
                        <div
                          className="bg-blue-600 h-2 rounded-full transition-all"
                          style={{
                            width: `${Math.min((book.currentPage / book.totalPages) * 100, 100)}%`,
                          }}
                        />
                      </div>
                    </div>
                  </div>
                </div>
              ))}
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <ReadingChart weeklyStats={stats.weeklyStats} monthlyStats={stats.monthlyStats} />
        <ReadingCalendar dailyStats={stats.dailyStats} />
      </div>
    </div>
  );
};

export default Dashboard;
