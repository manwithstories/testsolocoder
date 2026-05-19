import type { Card, Tag, ValidationResult } from '../types';

export const validateCard = (card: Omit<Card, 'id' | 'createdAt' | 'updatedAt' | 'review' | 'order'>): ValidationResult => {
  const errors: string[] = [];

  if (!card.title.trim()) {
    errors.push('卡片标题不能为空');
  } else if (card.title.length > 200) {
    errors.push('卡片标题不能超过200个字符');
  }

  if (!card.content.trim()) {
    errors.push('卡片内容不能为空');
  }

  card.tags.forEach((tag) => {
    if (!tag.trim()) {
      errors.push('标签名称不能为空');
    }
  });

  const duplicateTags = card.tags.filter((tag, index) => card.tags.indexOf(tag) !== index);
  if (duplicateTags.length > 0) {
    errors.push(`存在重复标签: ${duplicateTags.join(', ')}`);
  }

  return {
    isValid: errors.length === 0,
    errors,
  };
};

export const validateTag = (tag: Omit<Tag, 'id' | 'color'>, existingTags: Tag[], tagId?: string): ValidationResult => {
  const errors: string[] = [];

  if (!tag.name.trim()) {
    errors.push('标签名称不能为空');
  } else if (tag.name.length > 50) {
    errors.push('标签名称不能超过50个字符');
  }

  const duplicateTag = existingTags.find(
    (t) => t.name.toLowerCase() === tag.name.trim().toLowerCase() && t.id !== tagId
  );
  if (duplicateTag) {
    errors.push('标签名称已存在');
  }

  if (tag.parentId) {
    const parentTag = existingTags.find((t) => t.id === tag.parentId);
    if (!parentTag) {
      errors.push('父标签不存在');
    }
    if (tagId && tag.parentId === tagId) {
      errors.push('不能将自己设为父标签');
    }
  }

  return {
    isValid: errors.length === 0,
    errors,
  };
};

export const validateSearchQuery = (query: string): ValidationResult => {
  const errors: string[] = [];

  if (query.length > 100) {
    errors.push('搜索关键词不能超过100个字符');
  }

  return {
    isValid: errors.length === 0,
    errors,
  };
};
