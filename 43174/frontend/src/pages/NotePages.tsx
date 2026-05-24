import React, { useState, useEffect } from 'react';
import { Row, Col, Input, Select, Button, Space, Pagination, Modal, Form, Upload, message, Card, Tag, Rate } from 'antd';
import { SearchOutlined, PlusOutlined, UploadOutlined, DownloadOutlined, StarOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { noteApi, categoryApi } from '../services/api';
import { Note, Category } from '../types';
import { NoteCard } from '../components/NoteCard';
import { Loading } from '../components/Loading';
import { useAuthStore } from '../context/authStore';

const { Option } = Select;

export const NoteListPage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [notes, setNotes] = useState<Note[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(12);
  const [keyword, setKeyword] = useState('');
  const [subject] = useState<string | undefined>();
  const [categoryId, setCategoryId] = useState<string | undefined>();
  const [isFeatured, setIsFeatured] = useState(false);
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const { isAuthenticated } = useAuthStore();
  const [form] = Form.useForm();

  useEffect(() => {
    loadCategories();
    loadNotes();
  }, [page, keyword, subject, categoryId, isFeatured]);

  const loadCategories = async () => {
    try {
      const response: any = await categoryApi.getAll();
      setCategories(response.data || []);
    } catch (error) {
      console.error('Failed to load categories:', error);
    }
  };

  const loadNotes = async () => {
    setLoading(true);
    try {
      const params: any = { page, page_size: pageSize };
      if (keyword) params.keyword = keyword;
      if (subject) params.subject = subject;
      if (categoryId) params.category_id = categoryId;
      if (isFeatured) params.is_featured = true;
      const response: any = await noteApi.getAll(params);
      setNotes(response.data || []);
      setTotal(response.total || 0);
    } catch (error) {
      console.error('Failed to load notes:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (value: string) => {
    setKeyword(value);
    setPage(1);
  };

  const handleDownload = async (note: Note) => {
    try {
      await noteApi.incrementDownload(note.id);
      window.open(note.file_url, '_blank');
    } catch (error: any) {
      message.error(error.message || '下载失败');
    }
  };

  const handleCreate = async (values: any) => {
    try {
      await noteApi.create(values);
      message.success('笔记发布成功');
      setCreateModalVisible(false);
      form.resetFields();
      loadNotes();
    } catch (error: any) {
      message.error(error.message || '发布失败');
    }
  };

  const handleFileUpload = async (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    try {
      const response: any = await noteApi.uploadFile(formData);
      return response.data.file_url || '';
    } catch (error) {
      message.error('上传失败');
      return '';
    }
  };

  const handleCategoryChange = (value: string | undefined) => {
    setCategoryId(value);
    setPage(1);
  };

  if (loading && notes.length === 0) return <Loading />;

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">学习笔记</h1>
        <Space>
          <Button
            type={isFeatured ? 'primary' : 'default'}
            icon={<StarOutlined />}
            onClick={() => setIsFeatured(!isFeatured)}
          >
            精选
          </Button>
          {isAuthenticated && (
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => setCreateModalVisible(true)}
            >
              上传笔记
            </Button>
          )}
        </Space>
      </div>

      <Card className="mb-6">
        <Space direction="vertical" size="middle" className="w-full">
          <Input.Search
            placeholder="搜索笔记标题、描述..."
            allowClear
            enterButton={<SearchOutlined />}
            size="large"
            onSearch={handleSearch}
            className="w-full"
          />
          <Space wrap>
            <span className="text-gray-600">分类:</span>
            <Select
              placeholder="全部分类"
              allowClear
              style={{ width: 200 }}
              onChange={handleCategoryChange}
              value={categoryId}
            >
              {categories.map((cat) => (
                <Option key={cat.id} value={cat.id}>
                  {cat.name}
                </Option>
              ))}
            </Select>
          </Space>
        </Space>
      </Card>

      {loading ? (
        <Loading />
      ) : notes.length === 0 ? (
        <div className="text-center py-16">
          <p className="text-gray-500 text-lg">暂无笔记</p>
        </div>
      ) : (
        <>
          <Row gutter={[16, 16]}>
            {notes.map((note) => (
              <Col xs={24} sm={12} md={8} lg={6} key={note.id}>
                <NoteCard note={note} onDownload={handleDownload} />
              </Col>
            ))}
          </Row>
          <div className="flex justify-center mt-8">
            <Pagination
              current={page}
              pageSize={pageSize}
              total={total}
              onChange={setPage}
              showSizeChanger={false}
            />
          </div>
        </>
      )}

      <Modal
        title="上传笔记"
        open={createModalVisible}
        onCancel={() => setCreateModalVisible(false)}
        footer={null}
        width={600}
      >
        <Form form={form} layout="vertical" onFinish={handleCreate}>
          <Form.Item name="title" label="笔记标题" rules={[{ required: true }]}>
            <Input placeholder="请输入笔记标题" />
          </Form.Item>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="subject" label="科目">
                <Input placeholder="请输入科目" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="course_name" label="课程名称">
                <Input placeholder="请输入课程名称" />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="category_id" label="分类">
            <Select placeholder="请选择分类" allowClear>
              {categories.map((cat) => (
                <Option key={cat.id} value={cat.id}>
                  {cat.name}
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="description" label="笔记描述">
            <Input.TextArea rows={3} placeholder="请输入笔记描述" />
          </Form.Item>
          <Form.Item name="file_url" label="笔记文件" rules={[{ required: true }]}>
            <Upload
              maxCount={1}
              customRequest={async ({ file, onSuccess }: any) => {
                const url = await handleFileUpload(file);
                if (url && onSuccess) onSuccess(url);
              }}
            >
              <Button icon={<UploadOutlined />}>选择文件</Button>
            </Upload>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              上传
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export const NoteDetailPage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [note, setNote] = useState<Note | null>(null);
  const navigate = useNavigate();
  const { isAuthenticated } = useAuthStore();
  const id = window.location.pathname.split('/').pop();

  useEffect(() => {
    if (id) loadNote(id);
  }, [id]);

  const loadNote = async (noteId: string) => {
    setLoading(true);
    try {
      const response: any = await noteApi.getById(noteId);
      setNote(response.data);
    } catch (error) {
      console.error('Failed to load note:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDownload = async () => {
    if (!isAuthenticated) {
      message.warning('请先登录');
      navigate('/login');
      return;
    }
    if (note) {
      try {
        await noteApi.incrementDownload(note.id);
        window.open(note.file_url, '_blank');
      } catch (error: any) {
        message.error(error.message || '下载失败');
      }
    }
  };

  if (loading) return <Loading />;
  if (!note) return <div className="text-center py-16">笔记不存在</div>;

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <Card>
        <div className="flex justify-between items-start mb-6">
          <div>
            <h1 className="text-2xl font-bold mb-2">{note.title}</h1>
            {note.subject && <Tag color="blue" className="mb-2">{note.subject}</Tag>}
            {note.is_featured && <Tag color="gold" icon={<StarOutlined />}>精选笔记</Tag>}
          </div>
          <Button type="primary" icon={<DownloadOutlined />} onClick={handleDownload}>
            下载笔记
          </Button>
        </div>

        <div className="grid grid-cols-2 gap-4 mb-6">
          {note.course_name && (
            <div>
              <span className="text-gray-500">课程名称：</span>
              {note.course_name}
            </div>
          )}
          <div>
            <span className="text-gray-500">文件大小：</span>
            {(note.file_size / 1024 / 1024).toFixed(2)} MB
          </div>
          <div>
            <span className="text-gray-500">浏览次数：</span>
            {note.view_count}
          </div>
          <div>
            <span className="text-gray-500">下载次数：</span>
            {note.download_count}
          </div>
        </div>

        {note.description && (
          <div className="mb-6">
            <h3 className="font-semibold mb-2">笔记描述</h3>
            <p className="text-gray-600 whitespace-pre-wrap">{note.description}</p>
          </div>
        )}

        <div className="mb-6">
          <h3 className="font-semibold mb-2">评分</h3>
          <div className="flex items-center gap-2">
            <Rate disabled allowHalf value={note.rating} />
            <span className="text-gray-500">({note.rating_count} 人评价)</span>
          </div>
        </div>

        {note.uploader && (
          <div className="border-t pt-4">
            <h3 className="font-semibold mb-2">上传者</h3>
            <div className="flex items-center gap-3">
              <img
                src={note.uploader.avatar || 'https://api.dicebear.com/7.x/identicon/svg?seed=' + note.uploader.username}
                alt={note.uploader.username}
                className="w-12 h-12 rounded-full"
              />
              <div>
                <p className="font-medium">{note.uploader.username}</p>
                <p className="text-sm text-gray-500">评分: {note.uploader.rating}</p>
              </div>
            </div>
          </div>
        )}
      </Card>
    </div>
  );
};
