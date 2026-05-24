import React from 'react';
import { Spin } from 'antd';

interface LoadingProps {
  tip?: string;
  size?: 'small' | 'default' | 'large';
}

export const Loading: React.FC<LoadingProps> = ({ tip = '加载中...', size = 'large' }) => (
  <div className="flex items-center justify-center py-12">
    <Spin size={size} tip={tip} />
  </div>
);

export const PageLoading: React.FC = () => (
  <div className="min-h-screen flex items-center justify-center">
    <Spin size="large" tip="页面加载中..." />
  </div>
);
