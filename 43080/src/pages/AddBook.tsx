import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useBookStore } from '../store/useBookStore';
import { useToast } from '../components/ui/Toast';
import { BookForm } from '../components/books/BookForm';
import Empty from '../components/Empty';
import { Button } from '../components/ui/Button';
import { ArrowLeft } from 'lucide-react';
import type { BookFormData } from '../types/book';

const AddBook: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { showToast } = useToast();
  const { addBook, updateBook, getBookById } = useBookStore();
  const [existingBook, setExistingBook] = useState<BookFormData | null>(null);
  const [loading, setLoading] = useState(true);

  const isEdit = Boolean(id);

  useEffect(() => {
    if (id) {
      const book = getBookById(id);
      if (book) {
        setExistingBook({
          title: book.title,
          author: book.author,
          isbn: book.isbn || '',
          coverUrl: book.coverUrl || '',
          totalPages: book.totalPages,
          categories: book.categories,
          description: book.description || '',
          publishDate: book.publishDate || '',
          publisher: book.publisher || '',
        });
      }
      setLoading(false);
    } else {
      setLoading(false);
    }
  }, [id, getBookById]);

  if (loading) {
    return <div className="flex items-center justify-center h-64">加载中...</div>;
  }

  if (isEdit && !existingBook) {
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

  const handleSubmit = (data: BookFormData) => {
    try {
      if (isEdit && id) {
        updateBook(id, data);
        showToast('success', '书籍信息已更新');
      } else {
        addBook(data);
        showToast('success', '书籍添加成功');
      }
      navigate('/shelf');
    } catch (error) {
      showToast('error', isEdit ? '更新失败，请重试' : '添加失败，请重试');
    }
  };

  const handleCancel = () => {
    navigate('/shelf');
  };

  return (
    <div className="max-w-3xl mx-auto">
      <div className="flex items-center gap-4 mb-6">
        <button
          onClick={() => navigate('/shelf')}
          className="p-2 rounded-lg hover:bg-gray-100 transition-colors"
        >
          <ArrowLeft className="w-5 h-5 text-gray-600" />
        </button>
        <div>
          <h1 className="text-3xl font-bold text-gray-900">
            {isEdit ? '编辑书籍' : '添加书籍'}
          </h1>
          <p className="text-gray-500 mt-1">
            {isEdit ? '修改书籍的详细信息' : '手动输入书籍信息或通过ISBN查询'}
          </p>
        </div>
      </div>

      <BookForm
        initialData={existingBook || undefined}
        onSubmit={handleSubmit}
        onCancel={handleCancel}
        isEdit={isEdit}
      />
    </div>
  );
};

export default AddBook;
