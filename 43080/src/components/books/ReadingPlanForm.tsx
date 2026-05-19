import React, { useState, useMemo } from 'react';
import { Calendar, Target, Clock } from 'lucide-react';
import type { ReadingPlanFormData } from '../../types/reading';
import { Input } from '../ui/Input';
import { Button } from '../ui/Button';
import { Modal } from '../ui/Modal';
import { useToast } from '../ui/Toast';
import { useBookStore } from '../../store/useBookStore';
import { useReadingStore } from '../../store/useReadingStore';
import { validateReadingPlan } from '../../utils/validation';
import { getDateInputValue, formatDate } from '../../utils/date';
import { addDays } from 'date-fns';

interface ReadingPlanFormProps {
  bookId: string;
  onClose: () => void;
}

export const ReadingPlanForm: React.FC<ReadingPlanFormProps> = ({
  bookId,
  onClose,
}) => {
  const book = useBookStore((state) => state.getBookById(bookId));
  const addPlan = useReadingStore((state) => state.addPlan);
  const { showToast } = useToast();

  const [dailyPages, setDailyPages] = useState(20);
  const [startDate, setStartDate] = useState(getDateInputValue());

  const remainingPages = book ? book.totalPages - book.currentPage : 0;

  const estimatedDays = useMemo(() => {
    if (dailyPages <= 0) return 0;
    return Math.ceil(remainingPages / dailyPages);
  }, [dailyPages, remainingPages]);

  const estimatedCompletion = useMemo(() => {
    if (estimatedDays <= 0) return startDate;
    return formatDate(addDays(new Date(startDate), estimatedDays - 1));
  }, [startDate, estimatedDays]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!book) {
      showToast('error', '书籍不存在');
      return;
    }

    const validation = validateReadingPlan(dailyPages, startDate, estimatedCompletion);
    if (!validation.valid) {
      showToast('error', validation.error || '输入无效');
      return;
    }

    try {
      addPlan(bookId, {
        dailyPages,
        startDate,
        targetDate: estimatedCompletion,
      });
      showToast('success', '阅读计划已保存');
      onClose();
    } catch (error) {
      showToast('error', error instanceof Error ? error.message : '保存失败');
    }
  };

  if (!book) return null;

  return (
    <Modal
      isOpen={true}
      onClose={onClose}
      title="创建阅读计划"
      description={`为《${book.title}》制定阅读计划`}
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="bg-blue-50 rounded-lg p-4">
          <div className="flex items-center gap-2 text-blue-700 mb-2">
            <Target className="w-5 h-5" />
            <span className="font-medium">阅读目标</span>
          </div>
          <p className="text-sm text-blue-600">
            剩余 <span className="font-semibold">{remainingPages}</span> 页需要阅读
          </p>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <Input
            label="每日计划页数"
            type="number"
            min="1"
            value={dailyPages}
            onChange={(e) => setDailyPages(parseInt(e.target.value) || 1)}
            leftIcon={<Clock className="w-4 h-4" />}
          />
          <Input
            label="开始日期"
            type="date"
            value={startDate}
            onChange={(e) => setStartDate(e.target.value)}
            leftIcon={<Calendar className="w-4 h-4" />}
            min={getDateInputValue()}
          />
        </div>

        <div className="bg-gray-50 rounded-lg p-4 space-y-3">
          <h4 className="font-medium text-gray-700">计划预览</h4>
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span className="text-gray-500">预计天数：</span>
              <span className="font-semibold text-gray-900 ml-1">{estimatedDays} 天</span>
            </div>
            <div>
              <span className="text-gray-500">预计完成：</span>
              <span className="font-semibold text-blue-600 ml-1">{estimatedCompletion}</span>
            </div>
            <div className="col-span-2">
              <span className="text-gray-500">阅读进度：</span>
              <div className="w-full bg-gray-200 rounded-full h-2 mt-1">
                <div
                  className="bg-blue-600 h-2 rounded-full transition-all"
                  style={{ width: `${(book.currentPage / book.totalPages) * 100}%` }}
                />
              </div>
              <div className="flex justify-between text-xs text-gray-400 mt-1">
                <span>已读 {book.currentPage} 页</span>
                <span>目标 {book.totalPages} 页</span>
              </div>
            </div>
          </div>
        </div>

        <div className="flex justify-end gap-3 pt-2">
          <Button type="button" variant="outline" onClick={onClose}>
            取消
          </Button>
          <Button type="submit">
            保存计划
          </Button>
        </div>
      </form>
    </Modal>
  );
};
