import React, { useState } from 'react';
import { v4 as uuidv4 } from 'uuid';
import type { Tag, Card } from '../../types';
import { validateTag } from '../../utils/validators';
import { getRandomTagColor } from '../../utils/tagColors';
import { useToast } from '../common/Toast';
import ConfirmModal from '../common/ConfirmModal';

interface TagManagerProps {
  tags: Tag[];
  cards: Card[];
  selectedTagIds: string[];
  onAddTag: (tag: Tag) => void;
  onUpdateTag: (tag: Tag) => void;
  onDeleteTag: (id: string) => void;
  onMergeTags: (sourceId: string, targetId: string) => void;
  onToggleTagSelection: (id: string) => void;
}

const TagManager: React.FC<TagManagerProps> = ({
  tags,
  cards,
  selectedTagIds,
  onAddTag,
  onUpdateTag,
  onDeleteTag,
  onMergeTags,
  onToggleTagSelection,
}) => {
  const [name, setName] = useState('');
  const [parentId, setParentId] = useState<string | null>(null);
  const [editingTag, setEditingTag] = useState<Tag | null>(null);
  const [errors, setErrors] = useState<string[]>([]);
  const [deleteConfirm, setDeleteConfirm] = useState<string | null>(null);
  const [mergeMode, setMergeMode] = useState(false);
  const [mergeSource, setMergeSource] = useState<string | null>(null);
  const { showToast } = useToast();

  const rootTags = tags.filter((t) => t.parentId === null);

  const getChildTags = (parentId: string): Tag[] => {
    return tags.filter((t) => t.parentId === parentId);
  };

  const getTagCardCount = (tagId: string): number => {
    return cards.filter((c) => c.tags.includes(tagId)).length;
  };

  const resetForm = () => {
    setName('');
    setParentId(null);
    setEditingTag(null);
    setErrors([]);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const trimmedName = name.trim();

    const validation = validateTag(
      { name: trimmedName, parentId },
      tags,
      editingTag?.id
    );

    if (!validation.isValid) {
      setErrors(validation.errors);
      showToast(validation.errors[0], 'error');
      return;
    }

    if (editingTag) {
      onUpdateTag({ ...editingTag, name: trimmedName, parentId });
      showToast('标签更新成功', 'success');
    } else {
      const newTag: Tag = {
        id: uuidv4(),
        name: trimmedName,
        parentId,
        color: getRandomTagColor(),
      };
      onAddTag(newTag);
      showToast('标签创建成功', 'success');
    }

    resetForm();
  };

  const handleEdit = (tag: Tag) => {
    setEditingTag(tag);
    setName(tag.name);
    setParentId(tag.parentId);
    setErrors([]);
  };

  const handleDelete = (id: string) => {
    setDeleteConfirm(id);
  };

  const confirmDelete = () => {
    if (deleteConfirm) {
      onDeleteTag(deleteConfirm);
      setDeleteConfirm(null);
      showToast('标签删除成功', 'success');
    }
  };

  const handleMergeStart = (tagId: string) => {
    setMergeSource(tagId);
    setMergeMode(true);
  };

  const handleMergeComplete = (targetId: string) => {
    if (mergeSource && mergeSource !== targetId) {
      onMergeTags(mergeSource, targetId);
      showToast('标签合并成功', 'success');
    }
    setMergeMode(false);
    setMergeSource(null);
  };

  const renderTagTree = (tagList: Tag[], depth = 0) => {
    return tagList.map((tag) => {
      const children = getChildTags(tag.id);
      const cardCount = getTagCardCount(tag.id);
      const isSelected = selectedTagIds.includes(tag.id);

      return (
        <div key={tag.id} className="tag-tree-item" style={{ paddingLeft: depth * 20 }}>
          <div
            className={`tag-row ${isSelected ? 'selected' : ''} ${mergeMode && mergeSource === tag.id ? 'merge-source' : ''}`}
          >
            <span
              className="tag-color-dot"
              style={{ backgroundColor: tag.color }}
              onClick={() => onToggleTagSelection(tag.id)}
            />
            <span className="tag-name" onClick={() => onToggleTagSelection(tag.id)}>
              {tag.name}
            </span>
            <span className="tag-count">({cardCount})</span>

            {!mergeMode ? (
              <div className="tag-actions">
                <button className="icon-btn" onClick={() => handleEdit(tag)}>✏️</button>
                <button className="icon-btn" onClick={() => handleMergeStart(tag.id)}>🔗</button>
                <button className="icon-btn" onClick={() => handleDelete(tag.id)}>🗑️</button>
              </div>
            ) : (
              mergeSource !== tag.id && (
                <button
                  className="btn btn-primary btn-sm"
                  onClick={() => handleMergeComplete(tag.id)}
                >
                  合并到此
                </button>
              )
            )}
          </div>
          {children.length > 0 && renderTagTree(children, depth + 1)}
        </div>
      );
    });
  };

  return (
    <div className="tag-manager">
      <div className="list-header">
        <h3>
          标签管理
          <span className="count-badge">{tags.length}</span>
        </h3>
        {mergeMode && (
          <button className="btn btn-secondary" onClick={() => { setMergeMode(false); setMergeSource(null); }}>
            取消合并
          </button>
        )}
      </div>

      <div className="tag-form-section">
        <h4>{editingTag ? '编辑标签' : '新建标签'}</h4>
        <form onSubmit={handleSubmit}>
          {errors.length > 0 && (
            <div className="form-errors">
              {errors.map((error, index) => (
                <p key={index} className="error-message">{error}</p>
              ))}
            </div>
          )}
          <div className="form-group">
            <label htmlFor="tagName">标签名称</label>
            <input
              type="text"
              id="tagName"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="输入标签名称..."
              maxLength={50}
            />
          </div>
          <div className="form-group">
            <label htmlFor="parentTag">父标签（可选）</label>
            <select
              id="parentTag"
              value={parentId || ''}
              onChange={(e) => setParentId(e.target.value || null)}
            >
              <option value="">无（顶级标签）</option>
              {tags
                .filter((t) => t.id !== editingTag?.id)
                .map((tag) => (
                  <option key={tag.id} value={tag.id}>
                    {tag.name}
                  </option>
                ))}
            </select>
          </div>
          <div className="form-actions">
            {editingTag && (
              <button type="button" className="btn btn-secondary" onClick={resetForm}>
                取消
              </button>
            )}
            <button type="submit" className="btn btn-primary">
              {editingTag ? '保存' : '创建'}
            </button>
          </div>
        </form>
      </div>

      <div className="tag-tree">
        <h4>标签列表</h4>
        {tags.length === 0 ? (
          <p className="text-muted">还没有标签</p>
        ) : (
          renderTagTree(rootTags)
        )}
      </div>

      <ConfirmModal
        isOpen={!!deleteConfirm}
        title="删除标签"
        message="确定要删除这个标签吗？所有卡片上的此标签也会被移除。"
        confirmText="删除"
        variant="danger"
        onConfirm={confirmDelete}
        onCancel={() => setDeleteConfirm(null)}
      />
    </div>
  );
};

export default TagManager;
