import React from 'react';
import { Card, Tag, Rate, Button } from 'antd';
import { ShoppingCartOutlined, EyeOutlined, UserOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { Textbook } from '../types';

interface TextbookCardProps {
  textbook: Textbook;
  onBuy?: (textbook: Textbook) => void;
}

export const TextbookCard: React.FC<TextbookCardProps> = ({ textbook, onBuy }) => {
  const navigate = useNavigate();

  const getConditionColor = (condition: string) => {
    const colors: Record<string, string> = {
      new: 'green',
      like_new: 'cyan',
      good: 'blue',
      fair: 'orange',
    };
    return colors[condition] || 'default';
  };

  const getConditionText = (condition: string) => {
    const texts: Record<string, string> = {
      new: '全新',
      like_new: '九成新',
      good: '良好',
      fair: '一般',
    };
    return texts[condition] || condition;
  };

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      available: 'green',
      reserved: 'orange',
      sold: 'red',
    };
    return colors[status] || 'default';
  };

  const getStatusText = (status: string) => {
    const texts: Record<string, string> = {
      available: '在售',
      reserved: '已预定',
      sold: '已售出',
    };
    return texts[status] || status;
  };

  return (
    <Card
      hoverable
      className="card-hover h-full"
      cover={
        <div className="h-48 bg-gray-100 flex items-center justify-center overflow-hidden">
          {textbook.cover_image ? (
            <img
              src={textbook.cover_image}
              alt={textbook.title}
              className="h-full w-full object-cover"
            />
          ) : (
            <div className="text-gray-400 text-4xl">📚</div>
          )}
        </div>
      }
      actions={[
        <Button
          type="text"
          icon={<EyeOutlined />}
          onClick={() => navigate(`/textbooks/${textbook.id}`)}
        >
          详情
        </Button>,
        textbook.status === 'available' && (
          <Button
            type="primary"
            icon={<ShoppingCartOutlined />}
            onClick={() => onBuy?.(textbook)}
          >
            购买
          </Button>
        ),
      ].filter(Boolean)}
    >
      <div className="space-y-2">
        <div className="flex justify-between items-start">
          <h3 className="font-semibold text-base truncate flex-1 pr-2">{textbook.title}</h3>
          <Tag color={getStatusColor(textbook.status)}>{getStatusText(textbook.status)}</Tag>
        </div>
        {textbook.author && <p className="text-gray-500 text-sm">{textbook.author}</p>}
        {textbook.course_name && (
          <p className="text-gray-500 text-sm">课程: {textbook.course_name}</p>
        )}
        <div className="flex items-center gap-2">
          <Tag color={getConditionColor(textbook.condition)}>{getConditionText(textbook.condition)}</Tag>
        </div>
        <div className="flex justify-between items-center pt-2">
          <span className="text-xl font-bold text-red-500">¥{textbook.price}</span>
          {textbook.original_price > 0 && (
            <span className="text-gray-400 line-through text-sm">¥{textbook.original_price}</span>
          )}
        </div>
        {textbook.seller && (
          <div className="flex items-center gap-1 text-gray-500 text-sm">
            <UserOutlined />
            <span>{textbook.seller.username}</span>
            <Rate disabled allowHalf value={textbook.seller.rating} className="text-xs" />
          </div>
        )}
      </div>
    </Card>
  );
};
