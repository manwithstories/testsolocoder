import React, { useState } from 'react';
import { Bookmark } from 'lucide-react';
import type { BookmarkFormData } from '../../types/note';
import { Input } from '../ui/Input';
import { Button } from '../ui/Button';
import { Textarea } from '../ui/Textarea';
import { Modal } from '../ui/Modal';
import { useToast } from '../ui/Toast';
import { useNoteStore } from '../../store/useNoteStore';

interface BookmarkFormProps {
  bookId: string;
  onClose: () => void;
}

export const BookmarkForm: React.FC<BookmarkFormProps> = ({ bookId, onClose }) => {
  const { addBookmark } = useNoteStore();
  const { showToast } = useToast();
  const [formData, setFormData] = useState<BookmarkFormData & { description: string }>({
    page: 1,
    label: '',
    description: '',
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.page || formData.page <= 0) {
      showToast('error', '请输入有效的页码');
      return;
    }
    try {
      addBookmark(bookId, {
        page: formData.page,
        label: formData.label,
        description: formData.description,
      });
      showToast('success', '书签已添加');
      onClose();
    } catch (error) {
      showToast('error', '保存失败，请重试');
    }
  };

  return (
    <Modal isOpen={true} onClose={onClose} title="添加书签" size="md">
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="flex items-center gap-3 p-3 bg-blue-50 rounded-lg">
          <div className="p-2 bg-white rounded-lg shadow-sm">
            <Bookmark className="w-5 h-5 text-blue-600" />
          </div>
          <div>
            <p className="font-medium text-gray-900">新书签</p>
            <p className="text-sm text-gray-500">标记重要的页面以便日后查看</p>
          </div>
        </div>

        <Input
          label="页码 *"
          type="number"
          min="1"
          value={formData.page || ''}
          onChange={(e) => setFormData({ ...formData, page: parseInt(e.target.value) || 0 })}
          placeholder="输入页码"
        />

        <Input
          label="标签"
          value={formData.label}
          onChange={(e) => setFormData({ ...formData, label: e.target.value })}
          placeholder="例如：第3章 精彩片段"
        />

        <Textarea
          label="描述"
          value={formData.description}
          onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          placeholder="添加一些备注信息..."
          rows={3}
        />

        <div className="flex justify-end gap-3 pt-2">
          <Button type="button" variant="outline" onClick={onClose}>
            取消
          </Button>
          <Button type="submit">
            保存
          </Button>
        </div>
      </form>
    </Modal>
  );
};
