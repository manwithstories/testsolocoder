import React, { useState, useEffect } from 'react';
import { Row, Col, Input, Select, Button, Space, Pagination, Modal, Form, Upload, message, Card, Tag } from 'antd';
import { SearchOutlined, PlusOutlined, ScanOutlined, UploadOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { textbookApi, categoryApi, orderApi } from '../services/api';
import { Textbook, Category } from '../types';
import { TextbookCard } from '../components/TextbookCard';
import { Loading } from '../components/Loading';
import { useAuthStore } from '../context/authStore';

const { Option } = Select;

export const TextbookListPage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [textbooks, setTextbooks] = useState<Textbook[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(12);
  const [keyword, setKeyword] = useState('');
  const [categoryId, setCategoryId] = useState<string | undefined>();
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [isbnModalVisible, setIsbnModalVisible] = useState(false);
  const navigate = useNavigate();
  const { isAuthenticated } = useAuthStore();
  const [form] = Form.useForm();

  useEffect(() => {
    loadCategories();
    loadTextbooks();
  }, [page, keyword, categoryId]);

  const loadCategories = async () => {
    try {
      const response: any = await categoryApi.getAll();
      setCategories(response.data || []);
    } catch (error) {
      console.error('Failed to load categories:', error);
    }
  };

  const loadTextbooks = async () => {
    setLoading(true);
    try {
      const params: any = { page, page_size: pageSize, status: 'available' };
      if (keyword) params.keyword = keyword;
      if (categoryId) params.category_id = categoryId;
      const response: any = await textbookApi.getAll(params);
      setTextbooks(response.data || []);
      setTotal(response.total || 0);
    } catch (error) {
      console.error('Failed to load textbooks:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (value: string) => {
    setKeyword(value);
    setPage(1);
  };

  const handleCategoryChange = (value: string | undefined) => {
    setCategoryId(value);
    setPage(1);
  };

  const handleBuy = (textbook: Textbook) => {
    if (!isAuthenticated) {
      message.warning('请先登录');
      navigate('/login');
      return;
    }
    navigate(`/textbooks/${textbook.id}`);
  };

  const handleCreate = async (values: any) => {
    try {
      await textbookApi.create(values);
      message.success('教材发布成功');
      setCreateModalVisible(false);
      form.resetFields();
      loadTextbooks();
    } catch (error: any) {
      message.error(error.message || '发布失败');
    }
  };

  const handleISBNSearch = async (isbn: string) => {
    try {
      const response: any = await textbookApi.searchByISBN(isbn);
      if (response.data) {
        navigate(`/textbooks/${response.data.id}`);
      }
    } catch (error: any) {
      message.error(error.message || '未找到该ISBN的教材');
    }
  };

  const handleUpload = async (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    try {
      const response: any = await textbookApi.uploadCoverImage(formData);
      return response.data.image_url || response.data.avatar_url;
    } catch (error) {
      message.error('上传失败');
      return '';
    }
  };

  if (loading && textbooks.length === 0) return <Loading />;

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">教材市场</h1>
        <Space>
          <Button icon={<ScanOutlined />} onClick={() => setIsbnModalVisible(true)}>
            ISBN搜索
          </Button>
          {isAuthenticated && (
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => setCreateModalVisible(true)}
            >
              发布教材
            </Button>
          )}
        </Space>
      </div>

      <Card className="mb-6">
        <Space direction="vertical" size="middle" className="w-full">
          <Input.Search
            placeholder="搜索教材标题、ISBN、作者、课程..."
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
      ) : textbooks.length === 0 ? (
        <div className="text-center py-16">
          <p className="text-gray-500 text-lg">暂无教材</p>
        </div>
      ) : (
        <>
          <Row gutter={[16, 16]}>
            {textbooks.map((textbook) => (
              <Col xs={12} sm={12} md={8} lg={6} key={textbook.id}>
                <TextbookCard textbook={textbook} onBuy={handleBuy} />
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
        title="发布教材"
        open={createModalVisible}
        onCancel={() => setCreateModalVisible(false)}
        footer={null}
        width={600}
      >
        <Form form={form} layout="vertical" onFinish={handleCreate}>
          <Form.Item name="title" label="教材名称" rules={[{ required: true }]}>
            <Input placeholder="请输入教材名称" />
          </Form.Item>
          <Form.Item name="isbn" label="ISBN">
            <Input placeholder="请输入ISBN" />
          </Form.Item>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="author" label="作者">
                <Input placeholder="请输入作者" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="edition" label="版本">
                <Input placeholder="请输入版本" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="publisher" label="出版社">
                <Input placeholder="请输入出版社" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="course_name" label="课程名称">
                <Input placeholder="请输入课程名称" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="price" label="售价" rules={[{ required: true }]}>
                <Input type="number" placeholder="请输入售价" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="original_price" label="原价">
                <Input type="number" placeholder="请输入原价" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="condition" label="成色" rules={[{ required: true }]}>
                <Select placeholder="请选择成色">
                  <Option value="new">全新</Option>
                  <Option value="like_new">九成新</Option>
                  <Option value="good">良好</Option>
                  <Option value="fair">一般</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="category_id" label="分类">
                <Select placeholder="请选择分类" allowClear>
                  {categories.map((cat) => (
                    <Option key={cat.id} value={cat.id}>
                      {cat.name}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} placeholder="请输入教材描述" />
          </Form.Item>
          <Form.Item name="cover_image" label="封面图片">
            <Upload
              listType="picture-card"
              maxCount={1}
              customRequest={async ({ file, onSuccess }: any) => {
                const url = await handleUpload(file);
                if (url && onSuccess) onSuccess(url);
              }}
            >
              <div>
                <UploadOutlined />
                <div>上传封面</div>
              </div>
            </Upload>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              发布
            </Button>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="ISBN搜索"
        open={isbnModalVisible}
        onCancel={() => setIsbnModalVisible(false)}
        footer={null}
      >
        <Input.Search
          placeholder="请输入ISBN号"
          enterButton="搜索"
          size="large"
          onSearch={(value) => {
            handleISBNSearch(value);
            setIsbnModalVisible(false);
          }}
        />
      </Modal>
    </div>
  );
};

export const TextbookDetailPage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [textbook, setTextbook] = useState<Textbook | null>(null);
  const [buyModalVisible, setBuyModalVisible] = useState(false);
  const navigate = useNavigate();
  const { isAuthenticated } = useAuthStore();
  const id = window.location.pathname.split('/').pop();

  useEffect(() => {
    if (id) loadTextbook(id);
  }, [id]);

  const loadTextbook = async (textbookId: string) => {
    setLoading(true);
    try {
      const response: any = await textbookApi.getById(textbookId);
      setTextbook(response.data);
    } catch (error) {
      console.error('Failed to load textbook:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleBuy = () => {
    if (!isAuthenticated) {
      message.warning('请先登录');
      navigate('/login');
      return;
    }
    setBuyModalVisible(true);
  };

  const handleConfirmBuy = async () => {
    try {
      const orderData = {
        items: [{ textbook_id: textbook?.id, quantity: 1, price: textbook?.price }],
        seller_id: textbook?.seller_id,
      };
      await orderApi.create(orderData);
      message.success('订单创建成功');
      setBuyModalVisible(false);
      navigate('/orders');
    } catch (error: any) {
      message.error(error.message || '购买失败');
    }
  };

  if (loading) return <Loading />;
  if (!textbook) return <div className="text-center py-16">教材不存在</div>;

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <Row gutter={24}>
        <Col xs={24} md={10}>
          <Card>
            <div className="aspect-square bg-gray-100 flex items-center justify-center">
              {textbook.cover_image ? (
                <img
                  src={textbook.cover_image}
                  alt={textbook.title}
                  className="max-w-full max-h-full object-contain"
                />
              ) : (
                <div className="text-8xl">📚</div>
              )}
            </div>
          </Card>
        </Col>
        <Col xs={24} md={14}>
          <Card>
            <h1 className="text-2xl font-bold mb-4">{textbook.title}</h1>
            <div className="space-y-3 mb-6">
              {textbook.author && <p><span className="text-gray-500">作者：</span>{textbook.author}</p>}
              {textbook.isbn && <p><span className="text-gray-500">ISBN：</span>{textbook.isbn}</p>}
              {textbook.edition && <p><span className="text-gray-500">版本：</span>{textbook.edition}</p>}
              {textbook.publisher && <p><span className="text-gray-500">出版社：</span>{textbook.publisher}</p>}
              {textbook.course_name && <p><span className="text-gray-500">课程：</span>{textbook.course_name}</p>}
            </div>
            <div className="flex items-center gap-4 mb-6">
              <Tag color="green" className="text-lg px-4 py-1">
                {textbook.condition === 'new' ? '全新' : textbook.condition === 'like_new' ? '九成新' : textbook.condition === 'good' ? '良好' : '一般'}
              </Tag>
              <div className="text-3xl font-bold text-red-500">¥{textbook.price}</div>
              {textbook.original_price > 0 && (
                <div className="text-gray-400 line-through">¥{textbook.original_price}</div>
              )}
            </div>
            {textbook.description && (
              <div className="mb-6">
                <h3 className="font-semibold mb-2">商品描述</h3>
                <p className="text-gray-600 whitespace-pre-wrap">{textbook.description}</p>
              </div>
            )}
            {textbook.seller && (
              <div className="mb-6">
                <h3 className="font-semibold mb-2">卖家信息</h3>
                <div className="flex items-center gap-3">
                  <img
                    src={textbook.seller.avatar || 'https://api.dicebear.com/7.x/identicon/svg?seed=' + textbook.seller.username}
                    alt={textbook.seller.username}
                    className="w-12 h-12 rounded-full"
                  />
                  <div>
                    <p className="font-medium">{textbook.seller.username}</p>
                    <p className="text-sm text-gray-500">评分: {textbook.seller.rating}</p>
                  </div>
                </div>
              </div>
            )}
            {textbook.status === 'available' && (
              <Button type="primary" size="large" block onClick={handleBuy}>
                立即购买
              </Button>
            )}
            {textbook.status !== 'available' && (
              <Button disabled size="large" block>
                {textbook.status === 'sold' ? '已售出' : '已预定'}
              </Button>
            )}
          </Card>
        </Col>
      </Row>

      <Modal
        title="确认购买"
        open={buyModalVisible}
        onOk={handleConfirmBuy}
        onCancel={() => setBuyModalVisible(false)}
        okText="确认下单"
        cancelText="取消"
      >
        <p>您即将购买《{textbook.title}》</p>
        <p className="text-red-500 text-xl font-bold mt-2">¥{textbook.price}</p>
      </Modal>
    </div>
  );
};
