export const validateDuration = (duration: number): { valid: boolean; error?: string } => {
  if (duration === null || duration === undefined) {
    return { valid: false, error: '时长不能为空' }
  }
  if (typeof duration !== 'number' || isNaN(duration)) {
    return { valid: false, error: '时长必须是数字' }
  }
  if (duration < 0) {
    return { valid: false, error: '时长不能为负数' }
  }
  if (duration === 0) {
    return { valid: false, error: '时长不能为0' }
  }
  return { valid: true }
}

export const validateCategoryName = (name: string, existingNames: string[], excludeId?: string): { valid: boolean; error?: string } => {
  if (!name || !name.trim()) {
    return { valid: false, error: '分类名称不能为空' }
  }
  if (name.trim().length > 20) {
    return { valid: false, error: '分类名称不能超过20个字符' }
  }
  const trimmedName = name.trim().toLowerCase()
  const duplicate = existingNames.some(n => n.toLowerCase() === trimmedName && n !== excludeId)
  if (duplicate) {
    return { valid: false, error: '分类名称已存在' }
  }
  return { valid: true }
}

export const validateDateRange = (startDate: string, endDate: string): { valid: boolean; error?: string } => {
  if (!startDate || !endDate) {
    return { valid: false, error: '请选择完整的日期范围' }
  }
  const start = new Date(startDate)
  const end = new Date(endDate)
  if (isNaN(start.getTime()) || isNaN(end.getTime())) {
    return { valid: false, error: '日期格式不正确' }
  }
  if (start > end) {
    return { valid: false, error: '开始日期不能晚于结束日期' }
  }
  return { valid: true }
}

export const validateColor = (color: string): { valid: boolean; error?: string } => {
  const hexRegex = /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/
  if (!hexRegex.test(color)) {
    return { valid: false, error: '颜色格式不正确' }
  }
  return { valid: true }
}

export const validateNumberRange = (value: number, min: number, max: number): { valid: boolean; error?: string } => {
  if (typeof value !== 'number' || isNaN(value)) {
    return { valid: false, error: '必须是数字' }
  }
  if (value < min || value > max) {
    return { valid: false, error: `值必须在 ${min} 到 ${max} 之间` }
  }
  return { valid: true }
}
