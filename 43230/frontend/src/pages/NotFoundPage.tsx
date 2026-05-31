import React from 'react';
import { Result, Button, Space } from 'antd';
import { Link } from 'react-router-dom';
import { HomeOutlined, ArrowLeftOutlined } from '@ant-design/icons';

const NotFoundPage: React.FC = () => {
  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-50">
      <Result
        status="404"
        title="404"
        subTitle="抱歉，您访问的页面不存在"
        extra={
          <Space>
            <Link to="/">
              <Button type="primary" icon={<HomeOutlined />}>
                返回首页
              </Button>
            </Link>
            <Button onClick={() => window.history.back()} icon={<ArrowLeftOutlined />}>
              返回上一页
            </Button>
          </Space>
        }
      />
    </div>
  );
};

export default NotFoundPage;
