import dayjs from 'dayjs'

export const formatDate = (date: string | Date | null, format = 'YYYY-MM-DD HH:mm:ss') => {
  if (!date) return '-'
  return dayjs(date).format(format)
}

export const formatPrice = (price: number) => {
  return `¥${price.toFixed(2)}`
}

export const formatDuration = (minutes: number) => {
  if (minutes < 60) return `${minutes}分钟`
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  return mins > 0 ? `${hours}小时${mins}分钟` : `${hours}小时`
}

export const getOrderStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待接单',
    confirmed: '已确认',
    in_service: '服务中',
    completed: '已完成',
    cancelled: '已取消',
    disputed: '争议中',
  }
  return statusMap[status] || status
}

export const getOrderStatusColor = (status: string) => {
  const colorMap: Record<string, string> = {
    pending: 'orange',
    confirmed: 'blue',
    in_service: 'cyan',
    completed: 'green',
    cancelled: 'red',
    disputed: 'purple',
  }
  return colorMap[status] || 'default'
}

export const getProviderStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝',
    suspended: '已暂停',
  }
  return statusMap[status] || status
}

export const getInvitationStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待响应',
    accepted: '已接受',
    rejected: '已拒绝',
    expired: '已过期',
  }
  return statusMap[status] || status
}

export const getMessageTypeText = (type: string) => {
  const typeMap: Record<string, string> = {
    system: '系统通知',
    order: '订单通知',
    invitation: '邀请通知',
    review: '评价提醒',
    complaint: '投诉通知',
    withdraw: '提现通知',
  }
  return typeMap[type] || type
}

export const getBillTypeText = (type: string) => {
  const typeMap: Record<string, string> = {
    income: '收入',
    withdraw: '提现',
    penalty: '违约金',
    refund: '退款',
    commission: '佣金',
  }
  return typeMap[type] || type
}

export const getWithdrawStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝',
    completed: '已完成',
    failed: '失败',
  }
  return statusMap[status] || status
}

export const getComplaintStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待处理',
    processing: '处理中',
    resolved: '已解决',
    rejected: '已驳回',
  }
  return statusMap[status] || status
}

export const getRoleText = (role: string) => {
  const roleMap: Record<string, string> = {
    customer: '客户',
    service_provider: '服务人员',
    admin: '管理员',
  }
  return roleMap[role] || role
}

export const getCertificationStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝',
    none: '未提交',
  }
  return statusMap[status] || status
}

export const formatAddress = (address: { province: string; city: string; district: string; address: string }) => {
  return `${address.province}${address.city}${address.district}${address.address}`
}
