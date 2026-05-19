import React, { useState } from 'react';
import { BookOpen, Search, Upload } from 'lucide-react';
import type { BookFormData, Book } from '../../types/book';
import { Input } from '../ui/Input';
import { Button } from '../ui/Button';
import { Textarea } from '../ui/Textarea';
import { TagInput } from '../ui/TagInput';
import { DEFAULT_CATEGORIES } from '../../constants/categories';
import { fetchBookByISBN } from '../../services/isbnService';
import { useToast } from '../ui/Toast';
import { Skeleton } from '../ui/Skeleton';

interface BookFormProps {
  initialData?: BookFormData;
  onSubmit: (data: BookFormData) => void;
  onCancel: () => void;
  isLoading?: boolean;
  isEdit?: boolean;
}

export const BookForm: React.FC<BookFormProps> = ({ initialData, onSubmit, onCancel, isLoading, isEdit }) => {
  const [formData, setFormData] = useState<BookFormData>({
    title: initialData?.title || '',
    author: initialData?.author || '',
    isbn: initialData?.isbn || '',
    coverUrl: initialData?.coverUrl || '',
    totalPages: initialData?.totalPages || 0,
    categories: initialData?.categories || [],
  });
  const [isbnInput, setIsbnInput] = useState('');
  const [isFetching, setIsFetching] = useState(false);
  const { showToast } = useToast();

  const handleISBNLookup = async () => {
    if (!isbnInput.trim()) {
      showToast('warning', '请输入ISBN号');
      return;
    }

    setIsFetching(true);
    try {
      const result = await fetchBookByISBN(isbnInput);
      if (result) {
        setFormData({
          title: result.title,
          author: result.author,
          isbn: result.isbn,
          coverUrl: result.coverUrl || '',
          totalPages: result.totalPages || 0,
          categories: result.categories || [],
        });
        showToast('success', '书籍信息获取成功');
      } else {
        showToast('warning', '未找到该ISBN对应的书籍信息');
      }
    } catch (error) {
      showToast('error', error instanceof Error ? error.message : '获取书籍信息失败');
    } finally {
      setIsFetching(false);
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.title.trim()) {
      showToast('error', '请输入书名');
      return;
    }
    if (!formData.author.trim()) {
      showToast('error', '请输入作者');
      return;
    }
    if (formData.totalPages <= 0) {
      showToast('error', '请输入有效的总页数');
      return;
    }
    onSubmit(formData);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="bg-blue-50 rounded-xl p-4">
        <h3 className="font-medium text-blue-900 mb-3 flex items-center gap-2">
          <Search className="w-5 h-5" />
          ISBN 快速查询
        </h3>
        <div className="flex gap-2">
          <Input
            placeholder="输入ISBN号，如 9787111213826"
            value={isbnInput}
            onChange={(e) => setIsbnInput(e.target.value)}
            className="flex-1"
          />
          <Button
            type="button"
            onClick={handleISBNLookup}
            isLoading={isFetching}
            disabled={isFetching}
          >
            查询
          </Button>
        </div>
        <p className="text-xs text-blue-600 mt-2">
          提示：输入ISBN号可自动获取书籍信息
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Input
          label="书名 *"
          value={formData.title}
          onChange={(e) => setFormData({ ...formData, title: e.target.value })}
          placeholder="请输入书名"
        />
        <Input
          label="作者 *"
          value={formData.author}
          onChange={(e) => setFormData({ ...formData, author: e.target.value })}
          placeholder="请输入作者"
        />
        <Input
          label="ISBN"
          value={formData.isbn}
          onChange={(e) => setFormData({ ...formData, isbn: e.target.value })}
          placeholder="可选"
        />
        <Input
          label="总页数 *"
          type="number"
          min="1"
          value={formData.totalPages || ''}
          onChange={(e) => setFormData({ ...formData, totalPages: parseInt(e.target.value) || 0 })}
          placeholder="请输入总页数"
        />
      </div>

      <div>
        <Input
          label="封面图片URL"
          value={formData.coverUrl}
          onChange={(e) => setFormData({ ...formData, coverUrl: e.target.value })}
          placeholder="https://... (可选)"
          leftIcon={<Upload className="w-4 h-4" />}
        />
        {formData.coverUrl && (
          <div className="mt-3">
            <p className="text-sm text-gray-500 mb-2">预览：</p>
            <div className="w-24 h-36 bg-gray-100 rounded-lg overflow-hidden">
              <img
                src={formData.coverUrl}
                alt="封面预览"
                className="w-full h-full object-cover"
                onError={(e) => {
                  (e.target as HTMLImageElement).style.display = 'none';
                }}
              />
            </div>
          </div>
        )}
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1.5">
          分类标签
        </label>
        <TagInput
          tags={formData.categories}
          onChange={(tags) => setFormData({ ...formData, categories: tags })}
          suggestions={Array.from(DEFAULT_CATEGORIES)}
          placeholder="输入分类名称后按回车添加"
        />
      </div>

      <div className="flex justify-end gap-3 pt-4 border-t border-gray-100">
        <Button type="button" variant="outline" onClick={onCancel}>
          取消
        </Button>
        <Button type="submit" isLoading={isLoading}>
          {isEdit ? '保存修改' : '添加书籍'}
        </Button>
      </div>
    </form>
  );
};
