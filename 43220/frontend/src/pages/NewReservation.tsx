import { useState, useMemo } from 'react'
import { Card, Form, Input, Select, DatePicker, InputNumber, Button, message, Row, Col, Statistic } from 'antd'
import { ArrowLeftOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { useQuery, useMutation } from '@tanstack/react-query'
import { petApi, packageApi, reservationApi } from '@/services/api'
import dayjs from 'dayjs'

const { RangePicker } = DatePicker

export default function NewReservation() {
  const navigate = useNavigate()
  const [form] = Form.useForm()
  const [selectedStore, setSelectedStore] = useState<string>('')
  const [selectedPet, setSelectedPet] = useState<string>('')
  const [selectedPackage, setSelectedPackage] = useState<string>('')
  const [dateRange, setDateRange] = useState<any>(null)

  const { data: petsData } = useQuery({
    queryKey: ['pets', 'all'],
    queryFn: () => petApi.list({ page_size: 100 }),
  })

  const { data: packagesData } = useQuery({
    queryKey: ['packages', selectedStore],
    queryFn: () => packageApi.list({ store_id: selectedStore, page_size: 100 }),
    enabled: !!selectedStore,
  })

  const pets = petsData?.data?.items || []
  const packages = packagesData?.data?.items || []

  const selectedPkg = useMemo(
    () => packages.find((p: any) => p.id === selectedPackage),
    [packages, selectedPackage]
  )

  const totalDays = useMemo(() => {
    if (dateRange && dateRange.length === 2) {
      return dayjs(dateRange[1]).diff(dayjs(dateRange[0]), 'day') + 1
    }
    return 0
  }, [dateRange])

  const totalAmount = useMemo(() => {
    if (selectedPkg && totalDays > 0) {
      return totalDays * selectedPkg.price_per_day
    }
    return 0
  }, [selectedPkg, totalDays])

  const createMutation = useMutation({
    mutationFn: (values: any) => reservationApi.create(values),
    onSuccess: () => {
      message.success('预约创建成功')
      navigate('/reservations')
    },
    onError: (err: any) => message.error(err.message || '创建失败'),
  })

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (!dateRange || dateRange.length !== 2) {
        message.error('请选择日期范围')
        return
      }
      createMutation.mutate({
        ...values,
        check_in_date: dateRange[0].toISOString(),
        check_out_date: dateRange[1].toISOString(),
      })
    } catch {}
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-4">
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/reservations')}>
          返回
        </Button>
        <h2 className="text-xl font-semibold m-0">新建寄养预约</h2>
      </div>

      <Row gutter={24}>
        <Col xs={24} lg={16}>
          <Card title="预约信息">
            <Form form={form} layout="vertical">
              <Form.Item name="pet_id" label="选择宠物" rules={[{ required: true, message: '请选择宠物' }]}>
                <Select
                  placeholder="请选择宠物"
                  onChange={(val) => setSelectedPet(val)}
                  options={pets.map((p: any) => ({
                    value: p.id,
                    label: `${p.name} (${p.species}${p.breed ? ' - ' + p.breed : ''})`,
                  }))}
                />
              </Form.Item>

              <Form.Item name="store_id" label="选择门店" rules={[{ required: true, message: '请选择门店' }]}>
                <Select
                  placeholder="请选择门店"
                  showSearch
                  onChange={(val) => {
                    setSelectedStore(val)
                    setSelectedPackage('')
                    form.setFieldsValue({ package_id: undefined })
                  }}
                  options={[
                    { value: '00000000-0000-0000-0000-000000000000', label: '示例门店（请输入实际门店ID）' },
                  ]}
                />
              </Form.Item>

              <Form.Item name="package_id" label="选择套餐" rules={[{ required: true, message: '请选择套餐' }]}>
                <Select
                  placeholder="请先选择门店"
                  disabled={!selectedStore}
                  onChange={setSelectedPackage}
                  options={packages.map((p: any) => ({
                    value: p.id,
                    label: `${p.name} - ¥${p.price_per_day}/天 (${p.type === 'daycare' ? '日托' : '寄养'})`,
                  }))}
                />
              </Form.Item>

              <Form.Item label="入住日期" required>
                <RangePicker
                  className="w-full"
                  value={dateRange}
                  onChange={setDateRange}
                  disabledDate={(current: any) => current && current < dayjs().startOf('day')}
                />
              </Form.Item>

              <Form.Item name="remark" label="备注">
                <Input.TextArea rows={3} placeholder="特殊需求或备注" />
              </Form.Item>

              <Form.Item>
                <Button
                  type="primary"
                  size="large"
                  block
                  loading={createMutation.isPending}
                  onClick={handleSubmit}
                >
                  提交预约
                </Button>
              </Form.Item>
            </Form>
          </Card>
        </Col>

        <Col xs={24} lg={8}>
          <Card title="费用明细">
            {selectedPkg ? (
              <div className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-gray-500">套餐名称</span>
                  <span>{selectedPkg.name}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-500">每日价格</span>
                  <span>¥{selectedPkg.price_per_day.toFixed(2)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-500">寄养天数</span>
                  <span>{totalDays} 天</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-500">预付定金 (30%)</span>
                  <span>¥{(totalAmount * 0.3).toFixed(2)}</span>
                </div>
                <div className="border-t pt-3 mt-3">
                  <div className="flex justify-between text-lg font-bold">
                    <span>总计</span>
                    <span className="text-sky-600">¥{totalAmount.toFixed(2)}</span>
                  </div>
                </div>
              </div>
            ) : (
              <div className="text-center text-gray-400 py-8">
                请选择套餐查看费用明细
              </div>
            )}
          </Card>

          <Card title="温馨提示" className="mt-4">
            <ul className="text-sm text-gray-500 space-y-2 m-0 pl-5">
              <li>预约时需确保宠物疫苗在有效期内</li>
              <li>提交预约后需门店确认</li>
              <li>预付定金为总金额的30%</li>
              <li>取消预约请提前联系门店</li>
            </ul>
          </Card>
        </Col>
      </Row>
    </div>
  )
}
