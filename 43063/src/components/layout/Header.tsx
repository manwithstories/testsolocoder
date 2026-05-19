import React, { useState } from 'react';
import type { Card, Tag } from '../../types';
import { validateSearchQuery } from '../../utils/validators';
import { exportToJSON, exportToMarkdown, downloadFile, clearState } from '../../utils/storage';
import { useToast } from '../common/Toast';
import ConfirmModal from '../common/ConfirmModal';

interface HeaderProps {
  searchQuery: string;
  theme: 'light' | 'dark';
  cards: Card[];
  tags: Tag[];
  onSearchChange: (query: string) => void;
  onToggleTheme: () => void;
  onClearData: () => void;
}

const Header: React.FC<HeaderProps> = ({
  searchQuery,
  theme,
  cards,
  tags,
  onSearchChange,
  onToggleTheme,
  onClearData,
}) => {
  const [showExportMenu, setShowExportMenu] = useState(false);
  const [clearConfirm, setClearConfirm] = useState(false);
  const { showToast } = useToast();

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    const validation = validateSearchQuery(value);
    if (!validation.isValid) {
      showToast(validation.errors[0], 'error');
      return;
    }
    onSearchChange(value);
  };

  const handleExportJSON = () => {
    const json = exportToJSON(cards, tags);
    const filename = `knowledge-network-${new Date().toISOString().slice(0, 10)}.json`;
    downloadFile(json, filename, 'application/json');
    setShowExportMenu(false);
  };

  const handleExportMarkdown = () => {
    const markdown = exportToMarkdown(cards, tags);
    const filename = `knowledge-network-${new Date().toISOString().slice(0, 10)}.md`;
    downloadFile(markdown, filename, 'text/markdown');
    setShowExportMenu(false);
  };

  const handleClearData = () => {
    try {
      clearState();
      onClearData();
      setClearConfirm(false);
      showToast('数据已清除', 'success');
    } catch (error) {
      showToast('清除数据失败', 'error');
    }
  };

  return (
    <header className="header">
      <div className="search-bar">
        <span className="search-icon">🔍</span>
        <input
          type="text"
          placeholder="搜索标题、内容、标签..."
          value={searchQuery}
          onChange={handleSearchChange}
        />
        {searchQuery && (
          <button className="icon-btn" onClick={() => onSearchChange('')}>
            ✕
          </button>
        )}
      </div>

      <div className="header-actions">
        <div className="export-menu">
          <button
            className="btn btn-secondary"
            onClick={() => setShowExportMenu(!showExportMenu)}
          >
            📤 导出
          </button>
          {showExportMenu && (
            <div className="dropdown-menu">
              <button className="dropdown-item" onClick={handleExportJSON}>
                导出为 JSON
              </button>
              <button className="dropdown-item" onClick={handleExportMarkdown}>
                导出为 Markdown
              </button>
            </div>
          )}
        </div>

        <button className="btn btn-secondary" onClick={() => setClearConfirm(true)}>
          🗑️ 清除数据
        </button>

        <button className="btn btn-secondary" onClick={onToggleTheme}>
          {theme === 'light' ? '🌙 暗色' : '☀️ 亮色'}
        </button>
      </div>

      <ConfirmModal
        isOpen={clearConfirm}
        title="清除所有数据"
        message="确定要清除所有数据吗？此操作无法撤销，所有卡片和标签都将被删除。"
        confirmText="确认清除"
        variant="danger"
        onConfirm={handleClearData}
        onCancel={() => setClearConfirm(false)}
      />
    </header>
  );
};

export default Header;
