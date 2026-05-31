import React, { useEffect, useState } from 'react';
import {
  Table,
  Tag,
  Button,
  Space,
  Card,
  Typography,
  Modal,
  Form,
  Input,
  InputNumber,
  Select,
  message,
  Statistic,
  Row,
  Col,
  Progress,
  Empty,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  PrinterOutlined,
  CheckCircleOutlined,
  WarningOutlined,
  DashboardOutlined,
  PlayCircleOutlined,
  PauseCircleOutlined,
  ToolOutlined,
} from '@ant-design/icons';
import { printerApi } from '@/services/api';
import { PrinterDevice, MaterialInventory, PrintSchedule } from '@/types';
import { formatDate, getRelativeTime } from '@/utils/date';

const { Title, Text } = Typography;
const { Option } = Select;

const statusMap: Record<string, { text: string; color: string; icon: React.ReactNode }> = {
  idle: { text: '空闲', color: 'success', icon: <CheckCircleOutlined /> },
  printing: { text: '打印中', color: 'processing', icon: <PlayCircleOutlined /> },
  maintenance: { text: '维护中', color: 'warning', icon: <ToolOutlined /> },
  offline: { text: '离线', color: 'default', icon: <PauseCircleOutlined /> },
};

const PrinterDevicesPage: React.FC = () => {
  const [devices, setDevices] = useState<PrinterDevice[]>([]);
  const [inventory, setInventory] = useState<MaterialInventory[]>([]);
  const [schedules, setSchedules] = useState<PrintSchedule[]>([]);
  const [loading, setLoading] = useState(false);
  const [deviceModalVisible, setDeviceModalVisible] = useState(false);
  const [inventoryModalVisible, setInventoryModalVisible] = useState(false);
  const [editingDevice, setEditingDevice] = useState<PrinterDevice | null>(null);
  const [deviceForm] = Form.useForm();
  const [inventoryForm] = Form.useForm();
  const [materials, setMaterials] = useState<any[]>([]);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [devicesRes, inventoryRes, materialsRes] = await Promise.all([
        printerApi.getDevices(),
        printerApi.getInventory(),
        printerApi.getMaterials(),
      ]);

      setDevices(devicesRes.data || []);
      setInventory(inventoryRes.data || []);
      setMaterials(materialsRes.data || []);
    } catch (error) {
      console.error('Failed to fetch data:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleAddDevice = () => {
    setEditingDevice(null);
    deviceForm.resetFields();
    setDeviceModalVisible(true);
  };

  const handleEditDevice = (device: PrinterDevice) => {
    setEditingDevice(device);
    deviceForm.setFieldsValue(device);
    setDeviceModalVisible(true);
  };

  const handleSaveDevice = async () => {
    try {
      const values = await deviceForm.validateFields();
      if (editingDevice) {
        await printerApi.updateDevice(editingDevice.id, values);
        message.success('设备更新成功');
      } else {
        await printerApi.createDevice(values);
        message.success('设备添加成功');
      }
      setDeviceModalVisible(false);
      fetchData();
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const handleDeleteDevice = (id: string) => {
    Modal.confirm({
      title: '确认删除设备',
      content: '删除后设备数据将无法恢复，确定要删除吗？',
      okType: 'danger',
      onOk: async () => {
        try {
          await printerApi.deleteDevice(id);
          message.success('设备删除成功');
          fetchData();
        } catch (error: any) {
          message.error(error.response?.data?.error || '删除失败');
        }
      },
    });
  };

  const handleUpdateStatus = async (id: string, status: string) => {
    try {
      await printerApi.updateDevice(id, { status });
      message.success('状态更新成功');
      fetchData();
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const handleAddInventory = () => {
    inventoryForm.resetFields();
    setInventoryModalVisible(true);
  };

  const handleSaveInventory = async () => {
    try {
      const values = await inventoryForm.validateFields();
      await printerApi.createInventory(values);
      message.success('库存添加成功');
      setInventoryModalVisible(false);
      fetchData();
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const handleUpdateInventory = async (id: string, quantity: number) => {
    try {
      await printerApi.updateInventory(id, { quantity_grams: quantity });
      message.success('库存更新成功');
      fetchData();
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const deviceColumns = [
    {
      title: '设备名称',
      dataIndex: 'name',
      key: 'name',
      render: (name: string, record: PrinterDevice) => (
        <div>
          <div className="font-medium">{name}</div>
          <div className="text-sm text-gray-500">
            {record.manufacturer} {record.model}
          </div>
        </div>
      ),
    },
    {
      title: '最大打印尺寸',
      dataIndex: 'max_print_size',
      key: 'max_print_size',
    },
    {
      title: '累计打印',
      key: 'stats',
      render: (_: any, record: PrinterDevice) => (
        <div>
          <div>{record.total_print_jobs} 件</div>
          <div className="text-sm text-gray-500">
            {record.total_print_hours.toFixed(1)} 小时
          </div>
        </div>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const info = statusMap[status] || statusMap.offline;
        return (
          <Tag color={info.color} icon={info.icon}>
            {info.text}
          </Tag>
        );
      },
    },
    {
      title: '上次维护',
      dataIndex: 'last_maintenance',
      key: 'last_maintenance',
      render: (date: string) => (date ? getRelativeTime(date) : '暂无记录'),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: PrinterDevice) => (
        <Space size="small">
          <Button type="link" size="small" onClick={() => handleEditDevice(record)}>
            <EditOutlined /> 编辑
          </Button>
          {record.status === 'idle' && (
            <Button
              type="link"
              size="small"
              onClick={() => handleUpdateStatus(record.id, 'maintenance')}
            >
              <ToolOutlined /> 维护
            </Button>
          )}
          {record.status === 'maintenance' && (
            <Button
              type="link"
              size="small"
              type="primary"
              onClick={() => handleUpdateStatus(record.id, 'idle')}
            >
              <CheckCircleOutlined /> 完成维护
            </Button>
          )}
          <Button
            type="link"
            size="small"
            danger
            onClick={() => handleDeleteDevice(record.id)}
          >
            <DeleteOutlined /> 删除
          </Button>
        </Space>
      ),
    },
  ];

  const inventoryColumns = [
    {
      title: '材料',
      key: 'material',
      render: (_: any, record: MaterialInventory) => (
        <div>
          <div className="font-medium">{record.material?.name}</div>
          <div className="text-sm text-gray-500">颜色: {record.color}</div>
        </div>
      ),
    },
    {
      title: '当前库存',
      dataIndex: 'quantity_grams',
      key: 'quantity_grams',
      render: (qty: number, record: MaterialInventory) => {
        const percent = (qty / (record.reorder_level * 3)) * 100;
        const isLow = qty <= record.reorder_level;
        return (
          <div>
            <div className={isLow ? 'text-red-500 font-medium' : ''}>
              {qty.toFixed(0)}g
              {isLow && <WarningOutlined className="ml-2 text-red-500" />}
            </div>
            <Progress
              percent={Math.min(percent, 100)}
              size="small"
              strokeColor={isLow ? 'red' : 'green'}
              showInfo={false}
              className="mt-1"
            />
          </div>
        );
      },
    },
    {
      title: '预警值',
      dataIndex: 'reorder_level',
      key: 'reorder_level',
      render: (val: number) => `${val.toFixed(0)}g`,
    },
    {
      title: '最后更新',
      dataIndex: 'last_updated',
      key: 'last_updated',
      render: (date: string) => getRelativeTime(date),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: MaterialInventory) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            onClick={() => {
              Modal.confirm({
                title: '调整库存',
                content: (
                  <InputNumber
                    defaultValue={record.quantity_grams}
                    min={0}
                    id="inventory-adjust"
                    className="w-full"
                    addonAfter="g"
                  />
                ),
                onOk: async () => {
                  const input = document.getElementById(
                    'inventory-adjust'
                  ) as HTMLInputElement;
                  const val = parseFloat(input?.value || '0');
                  await handleUpdateInventory(record.id, val);
                },
              });
            }}
          >
            调整
          </Button>
          <Button
            type="link"
            size="small"
            danger
            onClick={() => {
              Modal.confirm({
                title: '确认删除',
                onOk: async () => {
                  try {
                    await printerApi.deleteInventory(record.id);
                    message.success('删除成功');
                    fetchData();
                  } catch (error) {
                    message.error('删除失败');
                  }
                },
              });
            }}
          >
            <DeleteOutlined />
          </Button>
        </Space>
      ),
    },
  ];

  const stats = {
    totalDevices: devices.length,
    idleDevices: devices.filter((d) => d.status === 'idle').length,
    printingDevices: devices.filter((d) => d.status === 'printing').length,
    lowInventory: inventory.filter((i) => i.quantity_grams <= i.reorder_level).length,
  };

  return (
    <div className="space-y-6">
      {/* 统计卡片 */}
      <Row gutter={16}>
        <Col xs={12} md={6}>
          <Card>
            <Statistic
              title="设备总数"
              value={stats.totalDevices}
              prefix={<PrinterOutlined />}
              valueStyle={{ color: '#3b82f6' }}
            />
          </Card>
        </Col>
        <Col xs={12} md={6}>
          <Card>
            <Statistic
              title="空闲设备"
              value={stats.idleDevices}
              prefix={<CheckCircleOutlined />}
              valueStyle={{ color: '#10b981' }}
            />
          </Card>
        </Col>
        <Col xs={12} md={6}>
          <Card>
            <Statistic
              title="正在打印"
              value={stats.printingDevices}
              prefix={<PlayCircleOutlined />}
              valueStyle={{ color: '#3b82f6' }}
            />
          </Card>
        </Col>
        <Col xs={12} md={6}>
          <Card>
            <Statistic
              title="库存预警"
              value={stats.lowInventory}
              prefix={<WarningOutlined />}
              valueStyle={{ color: stats.lowInventory > 0 ? '#ef4444' : '#10b981' }}
            />
          </Card>
        </Col>
      </Row>

      {/* 设备管理 */}
      <Card
        title={
          <Space>
            <PrinterOutlined />
            设备管理
          </Space>
        }
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAddDevice}>
            添加设备
          </Button>
        }
      >
        <Table
          rowKey="id"
          columns={deviceColumns}
          dataSource={devices}
          loading={loading}
          locale={{ emptyText: <Empty description="暂无设备" /> }}
          pagination={{ pageSize: 10 }}
        />
      </Card>

      {/* 材料库存 */}
      <Card
        title={
          <Space>
            <DashboardOutlined />
            材料库存
          </Space>
        }
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAddInventory}>
            添加库存
          </Button>
        }
      >
        <Table
          rowKey="id"
          columns={inventoryColumns}
          dataSource={inventory}
          loading={loading}
          locale={{ emptyText: <Empty description="暂无库存" /> }}
          pagination={{ pageSize: 10 }}
        />
      </Card>

      {/* 添加/编辑设备弹窗 */}
      <Modal
        title={editingDevice ? '编辑设备' : '添加设备'}
        open={deviceModalVisible}
        onCancel={() => setDeviceModalVisible(false)}
        onOk={handleSaveDevice}
        okText={editingDevice ? '保存' : '添加'}
      >
        <Form form={deviceForm} layout="vertical">
          <Form.Item
            name="name"
            label="设备名称"
            rules={[{ required: true, message: '请输入设备名称' }]}
          >
            <Input placeholder="如：FDM打印机-01" />
          </Form.Item>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="manufacturer" label="制造商">
                <Input placeholder="如：Creality" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="model" label="型号">
                <Input placeholder="如：Ender-3 Pro" />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item
            name="max_print_size"
            label="最大打印尺寸"
            rules={[{ required: true, message: '请输入打印尺寸' }]}
          >
            <Input placeholder="如：220x220x250mm" />
          </Form.Item>
          <Form.Item name="max_print_volume" label="最大打印体积 (L)">
            <InputNumber className="w-full" placeholder="如：12.1" />
          </Form.Item>
          <Form.Item
            name="supported_materials"
            label="支持材料"
            rules={[{ required: true, message: '请选择支持材料' }]}
          >
            <Select mode="multiple" placeholder="请选择支持的材料">
              <Option value="pla">PLA</Option>
              <Option value="abs">ABS</Option>
              <Option value="petg">PETG</Option>
              <Option value="tpu">TPU</Option>
              <Option value="nylon">尼龙</Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="supported_qualities"
            label="支持精度"
            rules={[{ required: true, message: '请选择支持精度' }]}
          >
            <Select mode="multiple" placeholder="请选择支持的打印精度">
              <Option value="draft">草稿级</Option>
              <Option value="standard">标准级</Option>
              <Option value="high">高精度</Option>
              <Option value="ultra">超精细</Option>
            </Select>
          </Form.Item>
          <Form.Item name="ip_address" label="IP地址">
            <Input placeholder="如：192.168.1.100" />
          </Form.Item>
        </Form>
      </Modal>

      {/* 添加库存弹窗 */}
      <Modal
        title="添加库存"
        open={inventoryModalVisible}
        onCancel={() => setInventoryModalVisible(false)}
        onOk={handleSaveInventory}
        okText="添加"
      >
        <Form form={inventoryForm} layout="vertical">
          <Form.Item
            name="material_id"
            label="材料类型"
            rules={[{ required: true, message: '请选择材料' }]}
          >
            <Select placeholder="请选择材料">
              {materials.map((m) => (
                <Option key={m.id} value={m.id}>
                  {m.name}
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="color"
            label="颜色"
            rules={[{ required: true, message: '请选择颜色' }]}
          >
            <Select placeholder="请选择颜色">
              <Option value="white">白色</Option>
              <Option value="black">黑色</Option>
              <Option value="gray">灰色</Option>
              <Option value="red">红色</Option>
              <Option value="blue">蓝色</Option>
              <Option value="green">绿色</Option>
              <Option value="yellow">黄色</Option>
              <Option value="transparent">透明</Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="quantity_grams"
            label="库存数量 (g)"
            rules={[{ required: true, message: '请输入库存数量' }]}
          >
            <InputNumber className="w-full" min={0} placeholder="如：1000" />
          </Form.Item>
          <Form.Item name="reorder_level" label="预警值 (g)">
            <InputNumber className="w-full" min={0} placeholder="如：200" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default PrinterDevicesPage;
