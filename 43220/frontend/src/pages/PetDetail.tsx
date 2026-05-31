import { useState } from 'react'
import { Card, Descriptions, Button, Table, Modal, Form, Input, DatePicker, Tabs, Tag, message, Divider } from 'antd'
import { ArrowLeftOutlined, PlusOutlined, EditOutlined } from '@ant-design/icons'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { petApi, dailyRecordApi } from '@/services/api'
import dayjs from 'dayjs'

export default function PetDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [vaccineModal, setVaccineModal] = useState(false)
  const [dewormModal, setDewormModal] = useState(false)
  const [vaccineForm] = Form.useForm()
  const [dewormForm] = Form.useForm()

  const { data: petData, isLoading } = useQuery({
    queryKey: ['pet', id],
    queryFn: () => petApi.get(id!),
    enabled: !!id,
  })

  const { data: vaccinesData } = useQuery({
    queryKey: ['pet', id, 'vaccines'],
    queryFn: () => petApi.getVaccines(id!),
    enabled: !!id,
  })

  const { data: dewormsData } = useQuery({
    queryKey: ['pet', id, 'deworms'],
    queryFn: () => petApi.getDeworms(id!),
    enabled: !!id,
  })

  const { data: recordsData } = useQuery({
    queryKey: ['pet', id, 'records'],
    queryFn: () => dailyRecordApi.listByPet(id!),
    enabled: !!id,
  })

  const pet = petData?.data
  const vaccines = vaccinesData?.data || []
  const deworms = dewormsData?.data || []
  const records = recordsData?.data?.items || []

  const addVaccineMutation = useMutation({
    mutationFn: (values: any) =>
      petApi.addVaccine({
        ...values,
        pet_id: id,
        vaccinated_at: values.vaccinated_at.toISOString(),
        expire_at: values.expire_at.toISOString(),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['pet', id, 'vaccines'] })
      message.success('添加成功')
      setVaccineModal(false)
      vaccineForm.resetFields()
    },
    onError: (err: any) => message.error(err.message || '添加失败'),
  })

  const addDewormMutation = useMutation({
    mutationFn: (values: any) =>
      petApi.addDeworm({
        ...values,
        pet_id: id,
        dewormed_at: values.dewormed_at.toISOString(),
        expire_at: values.expire_at.toISOString(),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['pet', id, 'deworms'] })
      message.success('添加成功')
      setDewormModal(false)
      dewormForm.resetFields()
    },
    onError: (err: any) => message.error(err.message || '添加失败'),
  })

  const vaccineColumns = [
    { title: '疫苗名称', dataIndex: 'vaccine_name', key: 'vaccine_name' },
    {
      title: '接种日期',
      dataIndex: 'vaccinated_at',
      key: 'vaccinated_at',
      render: (d: string) => dayjs(d).format('YYYY-MM-DD'),
    },
    {
      title: '到期日期',
      dataIndex: 'expire_at',
      key: 'expire_at',
      render: (d: string) => {
        const expired = dayjs(d).isBefore(dayjs())
        return (
          <Tag color={expired ? 'red' : 'green'}>
            {dayjs(d).format('YYYY-MM-DD')} {expired ? '(已过期)' : '(有效)'}
          </Tag>
        )
      },
    },
    { title: '接种医院', dataIndex: 'hospital', key: 'hospital' },
  ]

  const dewormColumns = [
    { title: '驱虫类型', dataIndex: 'deworm_type', key: 'deworm_type' },
    {
      title: '驱虫日期',
      dataIndex: 'dewormed_at',
      key: 'dewormed_at',
      render: (d: string) => dayjs(d).format('YYYY-MM-DD'),
    },
    {
      title: '到期日期',
      dataIndex: 'expire_at',
      key: 'expire_at',
      render: (d: string) => {
        const expired = dayjs(d).isBefore(dayjs())
        return (
          <Tag color={expired ? 'red' : 'green'}>
            {dayjs(d).format('YYYY-MM-DD')} {expired ? '(已过期)' : '(有效)'}
          </Tag>
        )
      },
    },
    { title: '用药', dataIndex: 'medicine', key: 'medicine' },
  ]

  const recordColumns = [
    {
      title: '记录日期',
      dataIndex: 'record_date',
      key: 'record_date',
      render: (d: string) => dayjs(d).format('YYYY-MM-DD'),
    },
    { title: '饮食', dataIndex: 'feed_status', key: 'feed_status' },
    { title: '活动', dataIndex: 'activity', key: 'activity' },
    { title: '健康', dataIndex: 'health_status', key: 'health_status' },
    { title: '心情', dataIndex: 'mood', key: 'mood' },
  ]

  if (isLoading) {
    return <div className="text-center py-10">加载中...</div>
  }

  if (!pet) {
    return <div className="text-center py-10">宠物不存在</div>
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-4">
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/pets')}>
          返回
        </Button>
        <h2 className="text-xl font-semibold m-0">{pet.name} 的档案</h2>
      </div>

      <Card>
        <div className="flex gap-6">
          {pet.avatar_url ? (
            <img
              src={pet.avatar_url}
              alt={pet.name}
              className="w-32 h-32 rounded-xl object-cover"
            />
          ) : (
            <div className="w-32 h-32 rounded-xl bg-gray-200 flex items-center justify-center text-5xl">
              🐾
            </div>
          )}
          <div className="flex-1">
            <Descriptions column={2} size="small">
              <Descriptions.Item label="名字">{pet.name}</Descriptions.Item>
              <Descriptions.Item label="物种">{pet.species}</Descriptions.Item>
              <Descriptions.Item label="品种">{pet.breed || '-'}</Descriptions.Item>
              <Descriptions.Item label="性别">{pet.gender === 'male' ? '公' : '母'}</Descriptions.Item>
              <Descriptions.Item label="体重">{pet.weight} kg</Descriptions.Item>
              <Descriptions.Item label="毛色">{pet.color || '-'}</Descriptions.Item>
              <Descriptions.Item label="出生日期">
                {pet.birth_date ? dayjs(pet.birth_date).format('YYYY-MM-DD') : '-'}
              </Descriptions.Item>
              <Descriptions.Item label="过敏史">{pet.allergies || '-'}</Descriptions.Item>
              <Descriptions.Item label="饮食习惯" span={2}>{pet.diet_habit || '-'}</Descriptions.Item>
              <Descriptions.Item label="性格特点" span={2}>{pet.temperament || '-'}</Descriptions.Item>
            </Descriptions>
          </div>
        </div>
      </Card>

      <Card>
        <Tabs
          items={[
            {
              key: 'vaccines',
              label: `疫苗记录 (${vaccines.length})`,
              children: (
                <>
                  <div className="mb-3 flex justify-end">
                    <Button type="primary" icon={<PlusOutlined />} onClick={() => setVaccineModal(true)}>
                      添加疫苗
                    </Button>
                  </div>
                  <Table columns={vaccineColumns} dataSource={vaccines} rowKey="id" size="small" pagination={false} />
                </>
              ),
            },
            {
              key: 'deworms',
              label: `驱虫记录 (${deworms.length})`,
              children: (
                <>
                  <div className="mb-3 flex justify-end">
                    <Button type="primary" icon={<PlusOutlined />} onClick={() => setDewormModal(true)}>
                      添加驱虫
                    </Button>
                  </div>
                  <Table columns={dewormColumns} dataSource={deworms} rowKey="id" size="small" pagination={false} />
                </>
              ),
            },
            {
              key: 'records',
              label: `寄养记录 (${records.length})`,
              children: (
                <Table columns={recordColumns} dataSource={records} rowKey="id" size="small" pagination={false} />
              ),
            },
          ]}
        />
      </Card>

      <Modal
        title="添加疫苗记录"
        open={vaccineModal}
        onCancel={() => {
          setVaccineModal(false)
          vaccineForm.resetFields()
        }}
        onOk={() => vaccineForm.submit()}
        confirmLoading={addVaccineMutation.isPending}
      >
        <Form
          form={vaccineForm}
          layout="vertical"
          onFinish={(values) => addVaccineMutation.mutate(values)}
        >
          <Form.Item name="vaccine_name" label="疫苗名称" rules={[{ required: true }]}>
            <Input placeholder="如：狂犬疫苗" />
          </Form.Item>
          <Form.Item name="vaccinated_at" label="接种日期" rules={[{ required: true }]}>
            <DatePicker className="w-full" />
          </Form.Item>
          <Form.Item name="expire_at" label="到期日期" rules={[{ required: true }]}>
            <DatePicker className="w-full" />
          </Form.Item>
          <Form.Item name="hospital" label="接种医院">
            <Input placeholder="医院名称" />
          </Form.Item>
          <Form.Item name="proof_url" label="证明文件URL">
            <Input placeholder="证明链接" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="添加驱虫记录"
        open={dewormModal}
        onCancel={() => {
          setDewormModal(false)
          dewormForm.resetFields()
        }}
        onOk={() => dewormForm.submit()}
        confirmLoading={addDewormMutation.isPending}
      >
        <Form
          form={dewormForm}
          layout="vertical"
          onFinish={(values) => addDewormMutation.mutate(values)}
        >
          <Form.Item name="deworm_type" label="驱虫类型" rules={[{ required: true }]}>
            <Input placeholder="如：体内驱虫" />
          </Form.Item>
          <Form.Item name="dewormed_at" label="驱虫日期" rules={[{ required: true }]}>
            <DatePicker className="w-full" />
          </Form.Item>
          <Form.Item name="expire_at" label="到期日期" rules={[{ required: true }]}>
            <DatePicker className="w-full" />
          </Form.Item>
          <Form.Item name="medicine" label="用药">
            <Input placeholder="药品名称" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
