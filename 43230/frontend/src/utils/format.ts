import { OrderStatus } from '@/types';

export const formatFileSize = (bytes: number): string => {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
  if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
  return (bytes / (1024 * 1024 * 1024)).toFixed(2) + ' GB';
};

export const formatPrice = (price: number): string => {
  return '¥' + price.toFixed(2);
};

export const formatDuration = (hours: number): string => {
  if (hours < 1) {
    return Math.round(hours * 60) + ' 分钟';
  }
  const h = Math.floor(hours);
  const m = Math.round((hours - h) * 60);
  return `${h} 小时 ${m} 分钟`;
};

export const getOrderStatusText = (status: OrderStatus): string => {
  const statusMap: Record<OrderStatus, string> = {
    pending: '待付款',
    paid: '已付款',
    printing: '打印中',
    quality_check: '质检中',
    shipped: '已发货',
    delivered: '已送达',
    completed: '已完成',
    cancelled: '已取消',
    refunded: '已退款',
  };
  return statusMap[status] || status;
};

export const getOrderStatusColor = (status: OrderStatus): string => {
  const colorMap: Record<OrderStatus, string> = {
    pending: 'orange',
    paid: 'blue',
    printing: 'cyan',
    quality_check: 'purple',
    shipped: 'geekblue',
    delivered: 'gold',
    completed: 'green',
    cancelled: 'red',
    refunded: 'red',
  };
  return colorMap[status] || 'default';
};

export const getRoleText = (role: string): string => {
  const roleMap: Record<string, string> = {
    designer: '建模师',
    printer: '打印商',
    customer: '客户',
    admin: '管理员',
  };
  return roleMap[role] || role;
};

export const getMaterialTypeText = (type: string): string => {
  const materialMap: Record<string, string> = {
    pla: 'PLA',
    abs: 'ABS',
    petg: 'PETG',
    tpu: 'TPU',
    resin: '光固化树脂',
    nylon: '尼龙',
    pc: 'PC',
  };
  return materialMap[type] || type;
};

export const getQualityText = (quality: string): string => {
  const qualityMap: Record<string, string> = {
    draft: '草稿级',
    standard: '标准级',
    high: '高精度',
    ultra: '超精细',
  };
  return qualityMap[quality] || quality;
};

export const generateModelThumbnail = (modelId: string): string => {
  return `https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=3d%20printing%20model%20blueprint%20technical%20render&image_size=square`;
};
