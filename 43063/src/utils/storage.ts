import type { AppState, Card, Tag } from '../types';
import { logAction, logError } from './logger';

const STORAGE_KEY = 'knowledge-network-state';

export const saveState = (state: AppState): void => {
  try {
    const serialized = JSON.stringify(state);
    localStorage.setItem(STORAGE_KEY, serialized);
    logAction('STATE_SAVED', { cardsCount: state.cards.length, tagsCount: state.tags.length });
  } catch (error) {
    logError('STATE_SAVE_FAILED', error);
    throw new Error('保存数据失败，请检查浏览器存储设置');
  }
};

export const loadState = (): AppState | null => {
  try {
    const serialized = localStorage.getItem(STORAGE_KEY);
    if (serialized === null) {
      return null;
    }
    const state = JSON.parse(serialized) as AppState;
    logAction('STATE_LOADED', { cardsCount: state.cards.length, tagsCount: state.tags.length });
    return state;
  } catch (error) {
    logError('STATE_LOAD_FAILED', error);
    return null;
  }
};

export const clearState = (): void => {
  try {
    localStorage.removeItem(STORAGE_KEY);
    logAction('STATE_CLEARED');
  } catch (error) {
    logError('STATE_CLEAR_FAILED', error);
    throw new Error('清除数据失败');
  }
};

export const exportToJSON = (cards: Card[], tags: Tag[]): string => {
  const exportData = {
    version: '1.0',
    exportedAt: new Date().toISOString(),
    cards,
    tags,
  };
  logAction('EXPORT_JSON', { cardsCount: cards.length, tagsCount: tags.length });
  return JSON.stringify(exportData, null, 2);
};

export const exportToMarkdown = (cards: Card[], tags: Tag[]): string => {
  const tagMap = new Map(tags.map((t) => [t.id, t]));
  let markdown = `# 知识卡片导出\n\n`;
  markdown += `导出时间: ${new Date().toLocaleString()}\n\n`;
  markdown += `卡片总数: ${cards.length}\n\n`;
  markdown += `---\n\n`;

  cards.forEach((card) => {
    markdown += `## ${card.title}\n\n`;
    const tagNames = card.tags
      .map((tagId) => tagMap.get(tagId)?.name || tagId)
      .join(', ');
    if (tagNames) {
      markdown += `**标签**: ${tagNames}\n\n`;
    }
    markdown += `**创建时间**: ${new Date(card.createdAt).toLocaleString()}\n\n`;
    if (card.linkedCardIds.length > 0) {
      markdown += `**关联卡片**: ${card.linkedCardIds.length} 张\n\n`;
    }
    markdown += `${card.content}\n\n`;
    markdown += `---\n\n`;
  });

  logAction('EXPORT_MARKDOWN', { cardsCount: cards.length });
  return markdown;
};

export const downloadFile = (content: string, filename: string, mimeType: string): void => {
  const blob = new Blob([content], { type: mimeType });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
  logAction('FILE_DOWNLOADED', { filename, mimeType });
};
