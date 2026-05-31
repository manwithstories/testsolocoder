import { useState } from 'react'
import { Card, List, Button, Modal, Form, Input, message } from 'antd'
import { EditOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { notificationApi } from '@/services/notification'
import { NotificationTemplate } from '@/types'

export function TemplateManagement() {
  const queryClient = useQueryClient()
  const [editVisible, setEditVisible] = useState(false)
  const [editingTemplate, setEditingTemplate] = useState<NotificationTemplate | null>(null)
  const [form] = Form.useForm()

  const { data, isLoading } = useQuery({
    queryKey: ['notification-templates'],
    queryFn: () => notificationApi.getTemplates({ page: 1, page_size: 50 }),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, title, content }: { id: string; title: string; content: string }) =>
      notificationApi.updateTemplate(id, { title, content }),
    onSuccess: () => {
      message.success('模板更新成功')
      setEditVisible(false)
      queryClient.invalidateQueries({ queryKey: ['notification-templates'] })
    },
    onError: (error: any) => {
      message.error(error.message || '更新失败')
    },
  })

  const handleEdit = (template: NotificationTemplate) => {
    setEditingTemplate(template)
    form.setFieldsValue({
      title: template.title,
      content: template.content,
    })
    setEditVisible(true)
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingTemplate) {
        updateMutation.mutate({
          id: editingTemplate.id,
          title: values.title,
          content: values.content,
        })
      }
    } catch (error) {
      console.error('Validation failed:', error)
    }
  }

  const templateTypeMap: Record<string, string> = {
    appointment_success: '预约成功通知',
    appointment_cancel: '预约取消通知',
    appointment_remind: '预约提醒通知',
    payment_success: '支付成功通知',
    payment_refund: '退款通知',
    review_reply: '评价回复通知',
    system: '系统通知',
  }

  return (
    <div className="page-container">
      <h2 style={{ marginBottom: 24 }}>通知模板管理</h2>

      <Card>
        <List
          dataSource={data?.items || []}
          loading={isLoading}
          renderItem={(template: NotificationTemplate) => (
            <List.Item
              actions={[
                <Button
                  key="edit"
                  type="link"
                  icon={<EditOutlined />}
                  onClick={() => handleEdit(template)}
                >
                  编辑
                </Button>,
              ]}
            >
              <List.Item.Meta
                title={
                  <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                    <span style={{ fontSize: 16, fontWeight: 500 }}>
                      {templateTypeMap[template.type] || template.type}
                    </span>
                    <span style={{ color: '#999', fontSize: 12 }}>({template.type})</span>
                  </div>
                }
                description={
                  <div>
                    <div style={{ marginBottom: 4 }}>
                      <strong>标题：</strong>{template.title}
                    </div>
                    <div>
                      <strong>内容：</strong>
                      <div style={{ marginTop: 4, padding: 8, background: '#f5f5f5', borderRadius: 4 }}>
                        {template.content}
                      </div>
                    </div>
                  </div>
                }
              />
            </List.Item>
          )}
        />
      </Card>

      <Modal
        title="编辑通知模板"
        open={editVisible}
        onOk={handleSubmit}
        onCancel={() => setEditVisible(false)}
        confirmLoading={updateMutation.isPending}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="title"
            label="模板标题"
            rules={[{ required: true, message: '请输入标题' }]}
          >
            <Input placeholder="请输入模板标题" />
          </Form.Item>
          <Form.Item
            name="content"
            label="模板内容"
            rules={[{ required: true, message: '请输入内容' }]}
          >
            <Input.TextArea
              rows={6}
              placeholder="请输入模板内容，支持变量占位符"
            />
          </Form.Item>
          <div style={{ color: '#999', fontSize: 12 }}>
            提示：模板内容支持变量占位符，如 {`{{appointment_id}}`}、{`{{date}}`} 等
          </div>
        </Form>
      </Modal>
    </div>
  )
}
