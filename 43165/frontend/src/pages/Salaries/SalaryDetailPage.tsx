import { useState, useEffect } from 'react';
import { Card, Descriptions, Tag, Button, List, Spin, message } from 'antd';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeftOutlined, ExportOutlined } from '@ant-design/icons';
import { salaryApi } from '../../services/api';
import { SalaryRecord, SalaryDetail } from '../../types';
import dayjs from 'dayjs';

export const SalaryDetailPage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [salary, setSalary] = useState<SalaryRecord | null>(null);
  const [details, setDetails] = useState<SalaryDetail[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchDetail();
  }, [id]);

  const fetchDetail = async () => {
    try {
      const res = await salaryApi.getSalary(id!);
      if (res.data.code === 200) {
        setSalary(res.data.data.salary);
        setDetails(res.data.data.details || []);
      }
    } catch (error) {
      message.error('获取薪资详情失败');
    } finally {
      setLoading(false);
    }
  };

  const handleExport = async () => {
    try {
      const res = await salaryApi.exportSalary(id!);
      const blob = new Blob([res.data]);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `salary_${id}.pdf`;
      a.click();
      window.URL.revokeObjectURL(url);
      message.success('导出成功');
    } catch (error) {
      message.error('导出失败');
    }
  };

  if (loading) {
    return <div className="flex justify-center items-center h-96"><Spin size="large" /></div>;
  }

  if (!salary) {
    return <div>薪资记录不存在</div>;
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate(-1)}>
          返回
        </Button>
        <Button icon={<ExportOutlined />} onClick={handleExport}>
          导出PDF
        </Button>
      </div>

      <Card title="薪资详情" className="mb-4">
        <Descriptions bordered column={2}>
          <Descriptions.Item label="临时工">{salary.temporary?.real_name}</Descriptions.Item>
          <Descriptions.Item label="岗位">{salary.job_posting?.position}</Descriptions.Item>
          <Descriptions.Item label="薪资周期">
            {dayjs(salary.period_start).format('YYYY-MM-DD')} ~ {dayjs(salary.period_end).format('YYYY-MM-DD')}
          </Descriptions.Item>
          <Descriptions.Item label="状态">
            <Tag color={salary.status === 'paid' ? 'green' : 'orange'}>
              {salary.status === 'paid' ? '已支付' : '待支付'}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="总工时">{salary.total_hours.toFixed(1)} 小时</Descriptions.Item>
          <Descriptions.Item label="基本工资">¥{salary.base_salary.toFixed(2)}</Descriptions.Item>
          <Descriptions.Item label="加班工时">{salary.overtime_hours.toFixed(1)} 小时</Descriptions.Item>
          <Descriptions.Item label="加班工资">¥{salary.overtime_pay.toFixed(2)}</Descriptions.Item>
          <Descriptions.Item label="扣款">¥{salary.deductions.toFixed(2)}</Descriptions.Item>
          <Descriptions.Item label="实发工资">
            <span className="text-green-600 font-bold text-lg">¥{salary.total_salary.toFixed(2)}</span>
          </Descriptions.Item>
          {salary.status === 'paid' && (
            <>
              <Descriptions.Item label="支付时间">
                {salary.payment_at ? dayjs(salary.payment_at).format('YYYY-MM-DD HH:mm') : '-'}
              </Descriptions.Item>
              <Descriptions.Item label="交易号">{salary.transaction_id}</Descriptions.Item>
            </>
          )}
        </Descriptions>
      </Card>

      <Card title="薪资明细">
        {details.length === 0 ? (
          <div className="text-center py-8 text-gray-500">暂无明细</div>
        ) : (
          <List
            itemLayout="horizontal"
            dataSource={details}
            renderItem={(item) => (
              <List.Item>
                <List.Item.Meta
                  title={item.description}
                  description={
                    <div>
                      <span>日期: {dayjs(item.date).format('YYYY-MM-DD')}</span>
                      <span className="mx-2">|</span>
                      <span>工时: {item.work_hours.toFixed(1)} 小时</span>
                      <span className="mx-2">|</span>
                      <span>时薪: ¥{item.hourly_rate.toFixed(2)}</span>
                    </div>
                  }
                />
                <div className="text-right">
                  <Tag color={item.type === 'overtime' ? 'orange' : 'blue'}>
                    {item.type === 'overtime' ? '加班' : '正常'}
                  </Tag>
                  <span className="font-semibold text-green-600 ml-2">¥{item.amount.toFixed(2)}</span>
                </div>
              </List.Item>
            )}
          />
        )}
      </Card>
    </div>
  );
};
