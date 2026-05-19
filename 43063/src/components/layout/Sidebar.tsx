import React from 'react';
import type { Tag } from '../../types';

interface SidebarProps {
  tags: Tag[];
  selectedTagIds: string[];
  currentView: string;
  onToggleTag: (id: string) => void;
  onSetView: (view: string) => void;
}

const Sidebar: React.FC<SidebarProps> = ({
  tags,
  selectedTagIds,
  currentView,
  onToggleTag,
  onSetView,
}) => {
  const navItems = [
    { id: 'cards', icon: '📚', label: '知识卡片' },
    { id: 'graph', icon: '🕸️', label: '知识图谱' },
    { id: 'review', icon: '📅', label: '复习系统' },
    { id: 'tags', icon: '🏷️', label: '标签管理' },
    { id: 'stats', icon: '📊', label: '统计分析' },
  ];

  const rootTags = tags.filter((t) => t.parentId === null);
  const getChildTags = (parentId: string) => tags.filter((t) => t.parentId === parentId);

  const renderTagTree = (tagList: Tag[], depth = 0) => {
    return tagList.map((tag) => {
      const children = getChildTags(tag.id);
      const isSelected = selectedTagIds.includes(tag.id);

      return (
        <React.Fragment key={tag.id}>
          <div
            className={`sidebar-tag ${isSelected ? 'selected' : ''}`}
            style={{ paddingLeft: 16 + depth * 16 }}
            onClick={() => onToggleTag(tag.id)}
          >
            <span
              className="tag-color-dot"
              style={{ backgroundColor: tag.color }}
            />
            <span className="sidebar-tag-name">{tag.name}</span>
          </div>
          {children.length > 0 && renderTagTree(children, depth + 1)}
        </React.Fragment>
      );
    });
  };

  return (
    <aside className="sidebar">
      <div className="sidebar-brand">
        <h1>🧠 知识网络</h1>
      </div>

      <nav className="sidebar-nav">
        {navItems.map((item) => (
          <button
            key={item.id}
            className={`nav-item ${currentView === item.id ? 'active' : ''}`}
            onClick={() => onSetView(item.id)}
          >
            <span className="nav-icon">{item.icon}</span>
            <span className="nav-label">{item.label}</span>
          </button>
        ))}
      </nav>

      <div className="sidebar-section">
        <h3 className="sidebar-section-title">标签筛选</h3>
        {tags.length === 0 ? (
          <p className="text-muted small">暂无标签</p>
        ) : (
          <div className="sidebar-tags">
            {renderTagTree(rootTags)}
          </div>
        )}
      </div>

      <div className="sidebar-footer">
        <p className="text-muted small">数据保存在本地浏览器</p>
      </div>
    </aside>
  );
};

export default Sidebar;
