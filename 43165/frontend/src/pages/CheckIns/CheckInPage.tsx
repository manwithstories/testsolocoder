import { useState, useEffect } from 'react';
import { Card, Button, Space, message, Statistic, Row, Col, Tag, Modal, Input, Form, Select } from 'antd';
import { ScanOutlined, EnvironmentOutlined, UserOutlined, LogoutOutlined } from '@ant-design/icons';
import { checkInApi, scheduleApi } from '../../services/api';
import { useAuthStore } from '../../context/AuthContext';
import { Schedule, CheckIn } from '../../types';
import dayjs from 'dayjs';

export const CheckInPage = () => {
  const { user } = useAuthStore();
  const [loading, setLoading] = useState(false);
  const [todaySchedules, setTodaySchedules] = useState<Schedule[]>([]);
  const [currentCheckIn, setCurrentCheckIn] = useState<CheckIn | null>(null);
  const [modalVisible, setModalVisible] = useState(false);
  const [checkInType, setCheckInType] = useState<'qr' | 'location' | 'face'>('location');
  const [form] = Form.useForm();

  const fetchTodaySchedules = async () => {
    setLoading(true);
    try {
      const res = await scheduleApi.getMySchedules({
        date: dayjs().format('YYYY-MM-DD'),
        status: 'scheduled',
      });
      if (res.data.code === 200) {
        setTodaySchedules(res.data.data.data || []);
      }
    } catch (error) {
      message.error('获取今日排班失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchCurrentCheckIn = async () => {
    try {
      const res = await checkInApi.getCheckIns({
        status: 'checked_in',
        start_date: dayjs().format('YYYY-MM-DD'),
        end_date: dayjs().add(1, 'day').format('YYYY-MM-DD'),
      });
      if (res.data.code === 200) {
        const records = res.data.data.data || [];
        if (records.length > 0) {
          setCurrentCheckIn(records[0]);
        }
      }
    } catch (error) {
      console.error('获取当前签到失败');
    }
  };

  useEffect(() => {
    if (user?.role === 'temporary') {
      fetchTodaySchedules();
      fetchCurrentCheckIn();
    }
  }, [user]);

  const handleCheckIn = async (values: any) => {
    setLoading(true);
    try {
      const data = {
        ...values,
        check_in_type: checkInType,
      };
      const res = await checkInApi.checkIn(data);
      if (res.data.code === 201) {
        message.success('签到成功');
        setCurrentCheckIn(res.data.data);
        setModalVisible(false);
        form.resetFields();
        fetchCurrentCheckIn();
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '签到失败');
    } finally {
      setLoading(false);
    }
  };

  const handleCheckOut = async () => {
    if (!currentCheckIn) return;
    setLoading(true);
    try {
      const res = await checkInApi.checkOut({ check_in_id: currentCheckIn.id });
      if (res.data.code === 200) {
        message.success('签退成功');
        setCurrentCheckIn(null);
        fetchCurrentCheckIn();
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '签退失败');
    } finally {
      setLoading(false);
    }
  };

  const getCurrentLocation = () => {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          form.setFieldsValue({
            latitude: position.coords.latitude,
            longitude: position.coords.longitude,
            location: `当前位置 (${position.coords.latitude.toFixed(4)}, ${position.coords.longitude.toFixed(4)})`,
          });
        },
        () => {
          message.warning('无法获取位置信息');
        }
      );
    }
  };

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">签到打卡</h2>

      <Row gutter={[16, 16]}>
        <Col xs={24} md={12}>
          <Card>
            <Statistic
              title="当前状态"
              value={currentCheckIn ? '已签到' : '未签到'}
              valueStyle={{ color: currentCheckIn ? '#52c41a' : '#ff4d4f' }}
              prefix={currentCheckIn ? <UserOutlined /> : <LogoutOutlined />}
            />
            {currentCheckIn && (
              <div className="mt-4">
                <p className="text-gray-500">签到时间: {dayjs(currentCheckIn.check_in_time).format('YYYY-MM-DD HH:mm:ss')}</p>
                <p className="text-gray-500">签到方式: {currentCheckIn.check_in_type === 'qr' ? '扫码' : currentCheckIn.check_in_type === 'face' ? '人脸识别' : '定位'}</p>
              </div>
            )}
          </Card>
        </Col>
        <Col xs={24} md={12}>
          <Card title="今日排班">
            {todaySchedules.length === 0 ? (
              <div className="text-center py-8 text-gray-500">今日暂无排班</div>
            ) : (
              todaySchedules.map((schedule) => (
                <div key={schedule.id} className="flex justify-between items-center py-2 border-b last:border-b-0">
                  <div>
                    <p className="font-medium">{schedule.job_posting?.position}</p>
                    <p className="text-sm text-gray-500">{schedule.start_time} - {schedule.end_time}</p>
                  </div>
                  <Tag color={schedule.status === 'completed' ? 'green' : schedule.status === 'in_progress' ? 'orange' : 'blue'}>
                    {schedule.status === 'scheduled' ? '待签到' : schedule.status === 'in_progress' ? '进行中' : '已完成'}
                  </Tag>
                </div>
              ))
            )}
          </Card>
        </Col>
      </Row>

      <div className="mt-6 flex justify-center gap-4">
        {!currentCheckIn ? (
          <Button type="primary" size="large" icon={<ScanOutlined />} onClick={() => setModalVisible(true)}>
            立即签到
          </Button>
        ) : (
          <Button danger size="large" icon={<LogoutOutlined />} onClick={handleCheckOut} loading={loading}>
            签退
          </Button>
        )}
      </div>

      <Modal
        title="签到打卡"
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
        width={500}
      >
        <Form form={form} onFinish={handleCheckIn} layout="vertical">
          <Form.Item name="schedule_id" label="选择排班" rules={[{ required: true, message: '请选择排班' }]}>
            <Select placeholder="请选择今日排班">
              {todaySchedules.map((s) => (
                <Select.Option key={s.id} value={s.id}>
                  {s.job_posting?.position} - {s.start_time} ~ {s.end_time}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>

          <div className="mb-4">
            <label className="block mb-2 font-medium">签到方式</label>
            <Space>
              <Button
                type={checkInType === 'location' ? 'primary' : 'default'}
                icon={<EnvironmentOutlined />}
                onClick={() => setCheckInType('location')}
              >
                定位签到
              </Button>
              <Button
                type={checkInType === 'qr' ? 'primary' : 'default'}
                icon={<ScanOutlined />}
                onClick={() => setCheckInType('qr')}
              >
                扫码签到
              </Button>
              <Button
                type={checkInType === 'face' ? 'primary' : 'default'}
                icon={<UserOutlined />}
                onClick={() => setCheckInType('face')}
              >
                人脸识别
              </Button>
            </Space>
          </div>

          {checkInType === 'location' && (
            <>
              <Form.Item name="latitude" hidden>
                <Input />
              </Form.Item>
              <Form.Item name="longitude" hidden>
                <Input />
              </Form.Item>
              <Form.Item name="location" label="位置">
                <Input />
              </Form.Item>
              <Button type="dashed" onClick={getCurrentLocation} block>
                <EnvironmentOutlined /> 获取当前位置
              </Button>
            </>
          )}

          {checkInType === 'qr' && (
            <Form.Item name="qr_code" label="扫码内容" rules={[{ required: true, message: '请扫描二维码' }]}>
              <Input placeholder="请扫描或输入二维码内容" />
            </Form.Item>
          )}

          {checkInType === 'face' && (
            <div className="text-center py-4 bg-gray-100 rounded">
              <div className="text-4xl mb-2">📷</div>
              <p>请将面部对准摄像头</p>
              <Button type="primary" className="mt-2">开始人脸识别</Button>
            </div>
          )}

          <Form.Item className="mt-6">
            <Button type="primary" htmlType="submit" block loading={loading}>
              确认签到
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
