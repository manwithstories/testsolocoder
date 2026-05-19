import React, { useState } from 'react';
import { BookOpen, Highlighter, Bookmark, X } from 'lucide-react';
import type { NoteFormData, NoteType } from '../../types/note';
import { Input } from '../ui/Input';
import { Button } from '../ui/Button';
import { Textarea } from '../ui/Textarea';
import { Select } from '../ui/Select';
import { Modal } from '../ui/Modal';
import { HIGHLIGHT_COLORS } from '../../constants/config';
import { useToast } from '../ui/Toast';
import { useNoteStore } from '../../store/useNoteStore';
import { useBookStore } from '../../store/useBookStore';

interface NoteFormProps {
  bookId?: string;
  onClose: () => void;
}

export const NoteForm: React.FC<NoteFormProps> = ({ bookId, onClose }) => {
  const { addNote } = useNoteStore();
  const { books } = useBookStore();
  const { showToast } = useToast();
  const [selectedBookId, setSelectedBookId] = useState(bookId || (books.length > 0 ? books[0].id : ''));
  const [formData, setFormData] = useState<NoteFormData>({
    page: 0,
    chapter: '',
    content: '',
    type: 'note',
    highlightColor: HIGHLIGHT_COLORS[0],
  });

  const bookOptions = books.map((b) => ({ value: b.id, label: b.title }));

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (formData.type !== 'bookmark' && !formData.content.trim()) {
      showToast('error', '请输入内容');
      return;
    }
    if (!selectedBookId) {
      showToast('error', '请选择书籍');
      return;
    }
    try {
      addNote(selectedBookId, formData);
      showToast('success', formData.type === 'bookmark' ? '书签已保存' : '笔记已保存');
      onClose();
    } catch (error) {
      showToast('error', '保存失败，请重试');
    }
  };

  const typeOptions: Array<{ value: NoteType; label: string }> = [
    { value: 'note', label: '笔记' },
    { value: 'highlight', label: '高亮' },
    { value: 'bookmark', label: '书签' },
  ];

  const typeIcons = {
    note: BookOpen,
    highlight: Highlighter,
    bookmark: Bookmark,
  };

  const Icon = typeIcons[formData.type];

  const getModalTitle = () => {
    switch (formData.type) {
      case 'note':
        return '添加笔记';
      case 'highlight':
        return '添加高亮';
      case 'bookmark':
        return '添加书签';
      default:
        return '添加笔记';
    }
  };

  return (
    <Modal isOpen={true} onClose={onClose} title={getModalTitle()} size="md">
      <form onSubmit={handleSubmit} className="space-y-4">
        {!bookId && (
          <Select
            label="选择书籍"
            value={selectedBookId}
            onChange={(e) => setSelectedBookId(e.target.value)}
            options={bookOptions}
          />
        )}

        <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
          <div className="p-2 bg-white rounded-lg shadow-sm">
            <Icon className="w-5 h-5 text-blue-600" />
          </div>
          <div className="flex-1">
            <Select
              label="类型"
              value={formData.type}
              onChange={(e) => setFormData({ ...formData, type: e.target.value as NoteType })}
              options={typeOptions}
            />
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <Input
            label="页码"
            type="number"
            min="0"
            value={formData.page || ''}
            onChange={(e) => setFormData({ ...formData, page: parseInt(e.target.value) || 0 })}
            placeholder="可选"
          />
          <Input
            label="章节"
            value={formData.chapter}
            onChange={(e) => setFormData({ ...formData, chapter: e.target.value })}
            placeholder="可选"
          />
        </div>

        {formData.type === 'highlight' && (
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">高亮颜色</label>
            <div className="flex gap-2">
              {HIGHLIGHT_COLORS.map((color) => (
                <button
                  key={color}
                  type="button"
                  onClick={() => setFormData({ ...formData, highlightColor: color })}
                  className={`w-8 h-8 rounded-full border-2 transition-all ${
                    formData.highlightColor === color
                      ? 'border-blue-600 scale-110'
                      : 'border-transparent'
                  }`}
                  style={{ backgroundColor: color }}
                />
              ))}
            </div>
          </div>
        )}

        <Textarea
          label={formData.type === 'bookmark' ? '描述（可选）' : '内容 *'}
          value={formData.content}
          onChange={(e) => setFormData({ ...formData, content: e.target.value })}
          placeholder={
            formData.type === 'note'
              ? '记录你的想法和心得...'
              : formData.type === 'highlight'
              ? '选中的文本内容...'
              : '添加一些备注信息...'
          }
          rows={4}
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
