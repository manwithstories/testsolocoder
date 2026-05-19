import React, { useState } from 'react';
import { Clock, BookOpen } from 'lucide-react';
import type { ReadingSessionFormData } from '../../types/reading';
import { Input } from '../ui/Input';
import { Button } from '../ui/Button';
import { Textarea } from '../ui/Textarea';
import { Modal } from '../ui/Modal';
import { useToast } from '../ui/Toast';
import { useBookStore } from '../../store/useBookStore';
import { useReadingStore } from '../../store/useReadingStore';
import { getDateInputValue } from '../../utils/date';
import { validateReadingSession } from '../../utils/validation';

interface ReadingSessionFormProps {
  bookId: string;
  onClose: () => void;
}

export const ReadingSessionForm: React.FC<ReadingSessionFormProps> = ({
  bookId,
  onClose,
}) => {
  const book = useBookStore((state) => state.getBookById(bookId));
  const addSession = useReadingStore((state) => state.addSession);
  const { showToast } = useToast();

  const [formData, setFormData] = useState<ReadingSessionFormData>({
    startPage: book?.currentPage || 0,
    endPage: Math.min((book?.currentPage || 0) + 20, book?.totalPages || 0),
    duration: 30,
    date: getDateInputValue(),
    notes: '',
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!book) {
      showToast('error', '书籍不存在');
      return;
    }

    const validation = validateReadingSession(
      formData.startPage,
      formData.endPage,
      book.totalPages
    );

    if (!validation.valid) {
      showToast('error', validation.error || '输入无效');
      return;
    }

    if (formData.duration <= 0) {
      showToast('error', '阅读时长必须大于0');
      return;
    }

    try {
      addSession(bookId, formData);
      showToast('success', '阅读记录已保存');
      onClose();
    } catch (error) {
      showToast('error', error instanceof Error ? error.message : '保存失败');
    }
  };

  if (!book) return null;

  const pagesRead = formData.endPage - formData.startPage;

  return (
    <Modal
      isOpen={true}
      onClose={onClose}
      title="记录阅读"
      description={`记录《${book.title}》的阅读进度`}
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="bg-blue-50 rounded-lg p-3">
          <p className="text-sm text-blue-700">
            <BookOpen className="w-4 h-4 inline mr-1" />
            总页数：{book.totalPages} 页，已读：{book.currentPage} 页
          </p>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <Input
            label="起始页"
            type="number"
            min="0"
            max={book.totalPages}
            value={formData.startPage}
            onChange={(e) => setFormData({ ...formData, startPage: parseInt(e.target.value) || 0 })}
          />
          <Input
            label="结束页"
            type="number"
            min="1"
            max={book.totalPages}
            value={formData.endPage}
            onChange={(e) => setFormData({ ...formData, endPage: parseInt(e.target.value) || 0 })}
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <Input
            label="阅读时长（分钟）"
            type="number"
            min="1"
            value={formData.duration}
            onChange={(e) => setFormData({ ...formData, duration: parseInt(e.target.value) || 0 })}
            leftIcon={<Clock className="w-4 h-4" />}
          />
          <Input
            label="日期"
            type="date"
            value={formData.date}
            onChange={(e) => setFormData({ ...formData, date: e.target.value })}
          />
        </div>

        <div className="flex items-center gap-2 text-sm text-gray-600 bg-gray-50 rounded-lg p-3">
          <span>本次阅读：</span>
          <span className="font-semibold text-blue-600">{pagesRead} 页</span>
          {formData.duration > 0 && (
            <>
              <span>·</span>
              <span>阅读速度：{Math.round(pagesRead / (formData.duration / 60))} 页/小时</span>
            </>
          )}
        </div>

        <Textarea
          label="阅读笔记（可选）"
          value={formData.notes}
          onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
          placeholder="记录这次阅读的心得和想法..."
          rows={3}
        />

        <div className="flex justify-end gap-3 pt-2">
          <Button type="button" variant="outline" onClick={onClose}>
            取消
          </Button>
          <Button type="submit">
            保存记录
          </Button>
        </div>
      </form>
    </Modal>
  );
};
