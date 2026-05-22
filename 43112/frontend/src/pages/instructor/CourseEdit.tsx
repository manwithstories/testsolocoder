import React, { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import {
  Card, Form, Input, Select, Button, Upload, message, Tabs,
  List, Modal, InputNumber, Space, Divider, Row, Col,
} from 'antd'
import {
  PlusOutlined, DeleteOutlined, UploadOutlined,
  PlayCircleOutlined, FileTextOutlined, QuestionOutlined,
} from '@ant-design/icons'
import { courseApi, chapterApi, lessonApi, uploadApi } from '@/services'
import { Course, Chapter, Lesson } from '@/types'

const { TextArea } = Input
const { Option } = Select

const InstructorCourseEdit: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [form] = Form.useForm()
  const [chapterForm] = Form.useForm()
  const [lessonForm] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [course, setCourse] = useState<Course | null>(null)
  const [chapters, setChapters] = useState<Chapter[]>([])
  const [chapterModalVisible, setChapterModalVisible] = useState(false)
  const [lessonModalVisible, setLessonModalVisible] = useState(false)
  const [selectedChapter, setSelectedChapter] = useState<Chapter | null>(null)
  const [uploading, setUploading] = useState(false)

  const isEdit = !!id

  const loadCourse = async () => {
    if (!id) return
    setLoading(true)
    try {
      const res = await courseApi.get(id)
      if (res.code === 0 && res.data) {
        setCourse(res.data)
        setChapters(res.data.chapters || [])
        form.setFieldsValue({
          title: res.data.title,
          subtitle: res.data.subtitle,
          description: res.data.description,
          cover: res.data.cover,
          category: res.data.category,
          level: res.data.level,
          price: res.data.price,
          original_price: res.data.original_price,
          is_free: res.data.is_free,
          tags: res.data.tags,
        })
      }
    } catch (error) {
      console.error('Failed to load course:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (isEdit) loadCourse()
  }, [id])

  const handleCoverUpload = async (file: File) => {
    setUploading(true)
    try {
      const res = await uploadApi.upload(file, 'image')
      if (res.code === 0 && res.data) {
        form.setFieldsValue({ cover: res.data.url })
        message.success('封面上传成功')
      }
    } catch (error: any) {
      message.error(error.message || '上传失败')
    } finally {
      setUploading(false)
    }
    return false
  }

  const handleSubmit = async (values: any) => {
    setLoading(true)
    try {
      if (isEdit && course) {
        await courseApi.update(course.id, values)
        message.success('课程更新成功')
      } else {
        const res = await courseApi.create(values)
        if (res.code === 0 && res.data) {
          message.success('课程创建成功')
          navigate(`/instructor/courses/${res.data.id}/edit`)
          return
        }
      }
    } catch (error: any) {
      message.error(error.message || '保存失败')
    } finally {
      setLoading(false)
    }
  }

  const handleAddChapter = async () => {
    if (!course) {
      message.warning('请先保存课程基本信息')
      return
    }
    try {
      const values = await chapterForm.validateFields()
      const res = await chapterApi.create(course.id, values)
      if (res.code === 0 && res.data) {
        message.success('章节添加成功')
        setChapterModalVisible(false)
        chapterForm.resetFields()
        loadCourse()
      }
    } catch (error: any) {
      if (error.errorFields) return
      message.error(error.message || '添加失败')
    }
  }

  const handleAddLesson = async () => {
    if (!selectedChapter) return
    try {
      const values = await lessonForm.validateFields()
      const res = await lessonApi.create(selectedChapter.id, values)
      if (res.code === 0 && res.data) {
        message.success('课时添加成功')
        setLessonModalVisible(false)
        lessonForm.resetFields()
        loadCourse()
      }
    } catch (error: any) {
      if (error.errorFields) return
      message.error(error.message || '添加失败')
    }
  }

  const handleVideoUpload = async (file: File, onProgress: (p: number) => void) => {
    try {
      const res = await uploadApi.upload(file, 'video', onProgress)
      if (res.code === 0 && res.data) {
        lessonForm.setFieldsValue({ video_url: res.data.url })
        message.success('视频上传成功')
      }
    } catch (error: any) {
      message.error(error.message || '上传失败')
    }
    return false
  }

  const handleDocUpload = async (file: File) => {
    try {
      const res = await uploadApi.upload(file, 'document')
      if (res.code === 0 && res.data) {
        lessonForm.setFieldsValue({
          doc_url: res.data.url,
          doc_name: res.data.file_name,
        })
        message.success('文档上传成功')
      }
    } catch (error: any) {
      message.error(error.message || '上传失败')
    }
    return false
  }

  const basicInfoTab = {
    key: 'basic',
    label: '基本信息',
    children: (
      <Card>
        <Form form={form} onFinish={handleSubmit} layout="vertical">
          <Row gutter={16}>
            <Col xs={24} md={16}>
              <Form.Item name="title" label="课程标题" rules={[{ required: true }]}>
                <Input placeholder="请输入课程标题" />
              </Form.Item>
              <Form.Item name="subtitle" label="课程副标题">
                <Input placeholder="请输入课程副标题" />
              </Form.Item>
              <Form.Item name="description" label="课程描述" rules={[{ required: true }]}>
                <TextArea rows={6} placeholder="请输入课程描述" />
              </Form.Item>
            </Col>
            <Col xs={24} md={8}>
              <Form.Item name="cover" label="课程封面">
                <Input placeholder="封面URL" />
              </Form.Item>
              <Upload beforeUpload={handleCoverUpload} showUploadList={false}>
                <Button icon={<UploadOutlined />} loading={uploading}>上传封面</Button>
              </Upload>
              {form.getFieldValue('cover') && (
                <img
                  src={form.getFieldValue('cover')}
                  alt="cover"
                  style={{ width: '100%', marginTop: 8, borderRadius: 4 }}
                />
              )}
            </Col>
          </Row>
          <Row gutter={16}>
            <Col xs={24} md={8}>
              <Form.Item name="category" label="分类" rules={[{ required: true }]}>
                <Input placeholder="如：编程、设计、营销" />
              </Form.Item>
            </Col>
            <Col xs={24} md={8}>
              <Form.Item name="level" label="难度" rules={[{ required: true }]} initialValue="beginner">
                <Select>
                  <Option value="beginner">入门</Option>
                  <Option value="intermediate">中级</Option>
                  <Option value="advanced">高级</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col xs={24} md={8}>
              <Form.Item name="tags" label="标签">
                <Input placeholder="多个标签用逗号分隔" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col xs={24} md={8}>
              <Form.Item name="is_free" label="是否免费" valuePropName="checked" initialValue={false}>
                <input type="checkbox" />
              </Form.Item>
            </Col>
            <Col xs={24} md={8}>
              <Form.Item name="price" label="价格">
                <InputNumber style={{ width: '100%' }} min={0} precision={2} />
              </Form.Item>
            </Col>
            <Col xs={24} md={8}>
              <Form.Item name="original_price" label="原价">
                <InputNumber style={{ width: '100%' }} min={0} precision={2} />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              {isEdit ? '保存修改' : '创建课程'}
            </Button>
          </Form.Item>
        </Form>
      </Card>
    ),
  }

  const chapterTab = {
    key: 'chapters',
    label: '章节管理',
    children: (
      <Card>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
          <h3>课程章节</h3>
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => setChapterModalVisible(true)}
          >
            添加章节
          </Button>
        </div>
        {chapters.length === 0 ? (
          <div style={{ textAlign: 'center', padding: 48, color: '#999' }}>暂无章节，点击上方按钮添加</div>
        ) : (
          chapters.map((chapter, idx) => (
            <Card
              key={chapter.id}
              size="small"
              title={`${idx + 1}. ${chapter.title}`}
              style={{ marginBottom: 8 }}
              extra={
                <Space>
                  <Button
                    type="link"
                    icon={<PlusOutlined />}
                    onClick={() => {
                      setSelectedChapter(chapter)
                      setLessonModalVisible(true)
                    }}
                  >
                    添加课时
                  </Button>
                </Space>
              }
            >
              <List
                size="small"
                dataSource={chapter.lessons || []}
                renderItem={(lesson) => (
                  <List.Item>
                    <List.Item.Meta
                      avatar={
                        lesson.type === 'video' ? <PlayCircleOutlined /> :
                        lesson.type === 'document' ? <FileTextOutlined /> :
                        <QuestionOutlined />
                      }
                      title={lesson.title}
                      description={
                        lesson.type === 'video'
                          ? `视频 · ${Math.floor(lesson.video_length / 60)}分${lesson.video_length % 60}秒`
                          : lesson.type === 'document'
                          ? `文档 · ${lesson.doc_name}`
                          : '测验'
                      }
                    />
                  </List.Item>
                )}
              />
            </Card>
          ))
        )}
      </Card>
    ),
  }

  return (
    <div>
      <h2>{isEdit ? '编辑课程' : '创建课程'}</h2>
      <Tabs
        defaultActiveKey="basic"
        items={isEdit ? [basicInfoTab, chapterTab] : [basicInfoTab]}
      />

      <Modal
        title="添加章节"
        open={chapterModalVisible}
        onCancel={() => setChapterModalVisible(false)}
        onOk={handleAddChapter}
        okText="添加"
      >
        <Form form={chapterForm} layout="vertical">
          <Form.Item name="title" label="章节标题" rules={[{ required: true }]}>
            <Input placeholder="请输入章节标题" />
          </Form.Item>
          <Form.Item name="position" label="排序" initialValue={0}>
            <InputNumber min={0} />
          </Form.Item>
          <Form.Item name="is_free" label="免费试看" valuePropName="checked" initialValue={false}>
            <input type="checkbox" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title={`添加课时 - ${selectedChapter?.title}`}
        open={lessonModalVisible}
        onCancel={() => setLessonModalVisible(false)}
        onOk={handleAddLesson}
        okText="添加"
        width={600}
      >
        <Form form={lessonForm} layout="vertical">
          <Form.Item name="title" label="课时标题" rules={[{ required: true }]}>
            <Input placeholder="请输入课时标题" />
          </Form.Item>
          <Form.Item name="type" label="类型" rules={[{ required: true }]} initialValue="video">
            <Select>
              <Option value="video">视频</Option>
              <Option value="document">文档</Option>
              <Option value="quiz">测验</Option>
            </Select>
          </Form.Item>
          <Form.Item name="video_url" label="视频URL">
            <Input placeholder="视频地址" />
          </Form.Item>
          <Upload beforeUpload={(file) => handleVideoUpload(file, () => {})} showUploadList={false}>
            <Button icon={<UploadOutlined />}>上传视频</Button>
          </Upload>
          <Form.Item name="video_length" label="视频时长(秒)">
            <InputNumber min={0} />
          </Form.Item>
          <Form.Item name="doc_url" label="文档URL">
            <Input placeholder="文档地址" />
          </Form.Item>
          <Upload beforeUpload={handleDocUpload} showUploadList={false}>
            <Button icon={<UploadOutlined />}>上传文档</Button>
          </Upload>
          <Form.Item name="position" label="排序" initialValue={0}>
            <InputNumber min={0} />
          </Form.Item>
          <Form.Item name="is_free" label="免费试看" valuePropName="checked" initialValue={false}>
            <input type="checkbox" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default InstructorCourseEdit
