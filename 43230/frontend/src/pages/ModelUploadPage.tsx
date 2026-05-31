import React, { useState } from 'react';
import {
  Card,
  Form,
  Input,
  InputNumber,
  Select,
  Button,
  Space,
  Typography,
  Upload,
  message,
  Progress,
  Row,
  Col,
  Tag,
  Divider,
  Steps,
} from 'antd';
import {
  UploadOutlined,
  SaveOutlined,
  CheckCircleOutlined,
  InboxOutlined,
  FileTextOutlined,
  DollarOutlined,
  TagOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { modelApi } from '@/services/api';
import { useChunkedUpload } from '@/hooks/useChunkedUpload';
import { formatFileSize } from '@/utils/format';

const { Title, Text, Paragraph } = Typography;
const { TextArea } = Input;
const { Option } = Select;
const { Step } = Steps;
const { Dragger } = Upload;

const ModelUploadPage: React.FC = () => {
  const [form] = Form.useForm();
  const navigate = useNavigate();
  const [currentStep, setCurrentStep] = useState(0);
  const [modelId, setModelId] = useState<string | null>(null);
  const [modelFile, setModelFile] = useState<File | null>(null);
  const [thumbnailFile, setThumbnailFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const { uploadFile } = useChunkedUpload();

  const categories = [
    '建筑',
    '人物',
    '交通工具',
    '家居',
    '电子产品',
    '玩具',
    '艺术',
    '其他',
  ];

  const handleBasicInfoSubmit = async () => {
    try {
      const values = await form.validateFields();
      const response = await modelApi.create({
        ...values,
        tags: values.tags || [],
        recommended_materials: values.recommended_materials || [],
      });
      setModelId(response.data.id);
      setCurrentStep(1);
      message.success('基本信息保存成功');
    } catch (error: any) {
      message.error(error.response?.data?.error || '保存失败');
    }
  };

  const handleFileUpload = async () => {
    if (!modelFile || !modelId) {
      message.error('请先选择模型文件');
      return;
    }

    setUploading(true);
    try {
      await modelApi.uploadFile(modelId, modelFile);
      message.success('模型文件上传成功');
    } catch (error: any) {
      message.error(error.response?.data?.error || '上传失败');
      setUploading(false);
      return;
    }

    if (thumbnailFile) {
      try {
        await modelApi.uploadThumbnail(modelId, thumbnailFile);
        message.success('缩略图上传成功');
      } catch (error: any) {
        message.error(error.response?.data?.error || '缩略图上传失败');
      }
    }

    setUploading(false);
    setCurrentStep(2);
  };

  const handleValidateFile = async (file: File) => {
    const maxSize = 500 * 1024 * 1024; // 500MB
    if (file.size > maxSize) {
      message.error(`文件大小不能超过500MB，当前文件：${formatFileSize(file.size)}`);
      return Upload.LIST_IGNORE;
    }

    const ext = file.name.split('.').pop()?.toLowerCase();
    const allowedExts = ['stl', 'obj'];
    if (!ext || !allowedExts.includes(ext)) {
      message.error('只支持STL和OBJ格式的3D模型文件');
      return Upload.LIST_IGNORE;
    }

    try {
      const response = await modelApi.validateFile(file);
      if (response.data.valid) {
        message.success(`文件格式校验通过，类型：${response.data.file_type}`);
        return false;
      } else {
        message.error('文件格式校验失败');
        return Upload.LIST_IGNORE;
      }
    } catch (error) {
      message.error('文件校验失败');
      return Upload.LIST_IGNORE;
    }
  };

  const handleComplete = async () => {
    if (!modelId) return;

    try {
      await modelApi.update(modelId, { status: 'published' });
      message.success('模型发布成功');
      navigate('/my-models');
    } catch (error: any) {
      message.error(error.response?.data?.error || '发布失败');
    }
  };

  const steps = [
    {
      title: '填写基本信息',
      icon: <FileTextOutlined />,
    },
    {
      title: '上传模型文件',
      icon: <UploadOutlined />,
    },
    {
      title: '完成发布',
      icon: <CheckCircleOutlined />,
    },
  ];

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <Card>
        <Title level={3} className="!mb-6">
          <UploadOutlined className="mr-2" />
          上传3D模型
        </Title>

        <Steps current={currentStep} items={steps} className="mb-8" />

        {currentStep === 0 && (
          <Form form={form} layout="vertical">
            <Row gutter={16}>
              <Col xs={24} md={16}>
                <Form.Item
                  name="title"
                  label="模型名称"
                  rules={[{ required: true, message: '请输入模型名称' }]}
                >
                  <Input placeholder="请输入模型名称，如：3D打印蝙蝠侠雕像" />
                </Form.Item>
              </Col>
              <Col xs={24} md={8}>
                <Form.Item
                  name="category"
                  label="分类"
                  rules={[{ required: true, message: '请选择分类' }]}
                >
                  <Select placeholder="请选择分类">
                    {categories.map((cat) => (
                      <Option key={cat} value={cat}>
                        {cat}
                      </Option>
                    ))}
                  </Select>
                </Form.Item>
              </Col>
            </Row>

            <Form.Item
              name="description"
              label="模型描述"
              rules={[{ required: true, message: '请输入模型描述' }]}
            >
              <TextArea
                rows={4}
                placeholder="请详细描述模型的特点、用途、打印建议等信息"
              />
            </Form.Item>

            <Row gutter={16}>
              <Col xs={24} md={8}>
                <Form.Item
                  name="price"
                  label="售价 (¥)"
                  rules={[
                    { required: true, message: '请输入售价' },
                    { type: 'number', min: 0, message: '售价不能为负数' },
                  ]}
                >
                  <InputNumber
                    className="w-full"
                    placeholder="0.00"
                    min={0}
                    step={0.01}
                    prefix={<DollarOutlined />}
                  />
                </Form.Item>
              </Col>
              <Col xs={24} md={8}>
                <Form.Item
                  name="license_type"
                  label="授权方式"
                  rules={[{ required: true, message: '请选择授权方式' }]}
                >
                  <Select placeholder="请选择">
                    <Option value="per_purchase">按件购买</Option>
                    <Option value="subscription">订阅下载</Option>
                  </Select>
                </Form.Item>
              </Col>
              <Col xs={24} md={8}>
                <Form.Item
                  name="subscription_price"
                  label="订阅价格 (¥/月)"
                  dependencies={['license_type']}
                  rules={[
                    ({ getFieldValue }) => ({
                      validator(_, value) {
                        if (getFieldValue('license_type') === 'subscription' && !value) {
                          return Promise.reject(new Error('请输入订阅价格'));
                        }
                        return Promise.resolve();
                      },
                    }),
                  ]}
                >
                  <InputNumber
                    className="w-full"
                    placeholder="0.00"
                    min={0}
                    step={0.01}
                  />
                </Form.Item>
              </Col>
            </Row>

            <Form.Item name="tags" label="标签">
              <Select mode="tags" placeholder="输入标签后按回车添加">
                {['建筑', '艺术', '玩具', '实用', '装饰'].map((tag) => (
                  <Option key={tag} value={tag}>
                    {tag}
                  </Option>
                ))}
              </Select>
            </Form.Item>

            <Divider orientation="left">
              <TagOutlined className="mr-1" /> 技术参数
            </Divider>

            <Row gutter={16}>
              <Col xs={24} md={8}>
                <Form.Item name="volume" label="体积 (cm³)">
                  <InputNumber className="w-full" placeholder="0.00" min={0} step={0.01} />
                </Form.Item>
              </Col>
              <Col xs={24} md={8}>
                <Form.Item name="bounding_box" label="外边框尺寸">
                  <Input placeholder="如：100x100x150mm" />
                </Form.Item>
              </Col>
              <Col xs={24} md={8}>
                <Form.Item name="print_time_hours" label="预估打印时间 (小时)">
                  <InputNumber className="w-full" placeholder="0.0" min={0} step={0.1} />
                </Form.Item>
              </Col>
            </Row>

            <Row gutter={16}>
              <Col xs={24} md={12}>
                <Form.Item name="polygon_count" label="多边形数量">
                  <InputNumber className="w-full" placeholder="0" min={0} />
                </Form.Item>
              </Col>
              <Col xs={24} md={12}>
                <Form.Item name="recommended_materials" label="推荐打印材料">
                  <Select mode="multiple" placeholder="请选择推荐材料">
                    <Option value="pla">PLA</Option>
                    <Option value="abs">ABS</Option>
                    <Option value="petg">PETG</Option>
                    <Option value="tpu">TPU</Option>
                    <Option value="resin">光固化树脂</Option>
                    <Option value="nylon">尼龙</Option>
                  </Select>
                </Form.Item>
              </Col>
            </Row>

            <div className="flex justify-end mt-6">
              <Button type="primary" size="large" onClick={handleBasicInfoSubmit}>
                下一步 <UploadOutlined />
              </Button>
            </div>
          </Form>
        )}

        {currentStep === 1 && (
          <div className="space-y-6">
            <Card title="上传模型文件" type="inner">
              <Dragger
                accept=".stl,.obj"
                multiple={false}
                maxCount={1}
                customRequest={({ file }) => {
                  setModelFile(file as File);
                  return Promise.resolve();
                }}
                beforeUpload={handleValidateFile}
                fileList={modelFile ? [{ name: modelFile.name, size: modelFile.size, uid: '1' }] : []}
                showUploadList={false}
              >
                <p className="ant-upload-drag-icon">
                  <InboxOutlined />
                </p>
                <p className="ant-upload-text">点击或拖拽模型文件到此处上传</p>
                <p className="ant-upload-hint">
                  支持 STL、OBJ 格式，最大 500MB。系统会自动校验文件格式。
                </p>
              </Dragger>
              {modelFile && (
                <div className="mt-4 p-4 bg-green-50 rounded-lg">
                  <Space>
                    <CheckCircleOutlined className="text-green-500" />
                    <span>已选择: {modelFile.name}</span>
                    <Tag color="blue">{formatFileSize(modelFile.size)}</Tag>
                  </Space>
                </div>
              )}
            </Card>

            <Card title="上传缩略图（可选）" type="inner">
              <Dragger
                accept=".jpg,.jpeg,.png"
                multiple={false}
                maxCount={1}
                customRequest={({ file }) => {
                  setThumbnailFile(file as File);
                  return Promise.resolve();
                }}
                beforeUpload={(file) => {
                  if (file.size > 5 * 1024 * 1024) {
                    message.error('缩略图大小不能超过5MB');
                    return Upload.LIST_IGNORE;
                  }
                  return false;
                }}
                fileList={
                  thumbnailFile
                    ? [{ name: thumbnailFile.name, size: thumbnailFile.size, uid: '1' }]
                    : []
                }
                showUploadList={false}
              >
                <p className="ant-upload-drag-icon">
                  <UploadOutlined />
                </p>
                <p className="ant-upload-text">点击或拖拽缩略图到此处上传</p>
                <p className="ant-upload-hint">
                  支持 JPG、PNG 格式，建议 1024x1024 分辨率，最大 5MB。
                </p>
              </Dragger>
              {thumbnailFile && (
                <div className="mt-4 p-4 bg-green-50 rounded-lg">
                  <Space>
                    <CheckCircleOutlined className="text-green-500" />
                    <span>已选择: {thumbnailFile.name}</span>
                    <Tag color="blue">{formatFileSize(thumbnailFile.size)}</Tag>
                  </Space>
                </div>
              )}
            </Card>

            {uploading && (
              <Card>
                <Progress percent={uploadProgress} status="active" />
                <Text type="secondary" className="block mt-2 text-center">
                  上传中，请不要关闭页面...
                </Text>
              </Card>
            )}

            <div className="flex justify-between mt-6">
              <Button size="large" onClick={() => setCurrentStep(0)}>
                上一步
              </Button>
              <Button
                type="primary"
                size="large"
                onClick={handleFileUpload}
                loading={uploading}
                disabled={!modelFile}
              >
                <UploadOutlined /> 上传并继续
              </Button>
            </div>
          </div>
        )}

        {currentStep === 2 && (
          <div className="text-center py-12">
            <div className="w-24 h-24 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6">
              <CheckCircleOutlined className="text-6xl text-green-500" />
            </div>
            <Title level={3}>模型上传完成！</Title>
            <Paragraph type="secondary" className="mb-8">
              您的模型已成功上传并发布，其他用户现在可以浏览和购买您的模型了。
            </Paragraph>

            <Space size="large">
              <Button size="large" onClick={handleComplete}>
                <SaveOutlined /> 发布模型
              </Button>
              <Button type="primary" size="large" onClick={() => navigate('/my-models')}>
                查看我的模型
              </Button>
            </Space>
          </div>
        )}
      </Card>
    </div>
  );
};

export default ModelUploadPage;
