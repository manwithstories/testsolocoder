import { useState } from 'react'
import { Card, Form, Input, InputNumber, Select, Button, message } from 'antd'
import { useNavigate } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { serviceApi } from '@/services/service'
import { CreateServiceRequest } from '@/types'

const createServiceSchema = z.object({
  title: z.string().min(1, '请输入服务标题').max(200, '标题最多200个字符'),
  description: z.string().max(1000, '描述最多1000个字符').optional(),
  service_type: z.enum(['legal', 'counseling', 'financial', 'other'], { required_error: '请选择服务类型' }),
  price: z.number().min(0.01, '价格必须大于0'),
  duration_minutes: z.number().min(15, '时长至少15分钟').max(480, '时长最多480分钟'),
  tags: z.string().max(500, '标签最多500个字符').optional(),
})

type CreateServiceFormData = z.infer<typeof createServiceSchema>

export function CreateService() {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateServiceFormData>({
    resolver: zodResolver(createServiceSchema),
    defaultValues: {
      service_type: 'legal',
      duration_minutes: 60,
    },
  })

  const createMutation = useMutation({
    mutationFn: (data: CreateServiceFormData) => serviceApi.create(data),
    onSuccess: (data) => {
      message.success('服务创建成功')
      navigate(`/services/${data.id}`)
    },
    onError: (error: any) => {
      message.error(error.message || '创建失败')
    },
    onSettled: () => {
      setLoading(false)
    },
  })

  const onSubmit = (data: CreateServiceFormData) => {
    setLoading(true)
    createMutation.mutate(data)
  }

  return (
    <div className="page-container">
      <h2 style={{ marginBottom: 24 }}>发布服务</h2>
      <Card>
        <Form layout="vertical" onFinish={handleSubmit(onSubmit)}>
          <Form.Item
            label="服务标题"
            validateStatus={errors.title ? 'error' : ''}
            help={errors.title?.message}
          >
            <Input placeholder="请输入服务标题" {...register('title')} />
          </Form.Item>

          <Form.Item
            label="服务描述"
            validateStatus={errors.description ? 'error' : ''}
            help={errors.description?.message}
          >
            <Input.TextArea
              rows={4}
              placeholder="请输入服务描述"
              {...register('description')}
            />
          </Form.Item>

          <Form.Item
            label="服务类型"
            validateStatus={errors.service_type ? 'error' : ''}
            help={errors.service_type?.message}
          >
            <Select
              {...register('service_type')}
              options={[
                { value: 'legal', label: '法律咨询' },
                { value: 'counseling', label: '心理咨询' },
                { value: 'financial', label: '财务咨询' },
                { value: 'other', label: '其他服务' },
              ]}
            />
          </Form.Item>

          <Form.Item
            label="服务价格（元）"
            validateStatus={errors.price ? 'error' : ''}
            help={errors.price?.message}
          >
            <InputNumber
              style={{ width: '100%' }}
              min={0.01}
              step={0.01}
              placeholder="请输入服务价格"
              {...register('price')}
            />
          </Form.Item>

          <Form.Item
            label="服务时长（分钟）"
            validateStatus={errors.duration_minutes ? 'error' : ''}
            help={errors.duration_minutes?.message}
          >
            <InputNumber
              style={{ width: '100%' }}
              min={15}
              max={480}
              step={15}
              placeholder="请输入服务时长"
              {...register('duration_minutes')}
            />
          </Form.Item>

          <Form.Item
            label="标签（用逗号分隔）"
            validateStatus={errors.tags ? 'error' : ''}
            help={errors.tags?.message}
          >
            <Input placeholder="如：合同审查,劳动纠纷,民事诉讼" {...register('tags')} />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block>
              发布服务
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}
