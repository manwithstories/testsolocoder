import React, { useState, useEffect } from 'react';
import { v4 as uuidv4 } from 'uuid';
import type { Card, Tag } from '../../types';
import { validateCard } from '../../utils/validators';
import { useToast } from '../common/Toast';
import { createInitialReviewInfo } from '../../utils/spacedRepetition';

interface CardFormProps {
  card?: Card | null;
  tags: Tag[];
  onSubmit: (card: Card) => void;
  onCancel: () => void;
}

const CardForm: React.FC<CardFormProps> = ({ card, tags, onSubmit, onCancel }) => {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [errors, setErrors] = useState<string[]>([]);
  const { showToast } = useToast();

  useEffect(() => {
    if (card) {
      setTitle(card.title);
      setContent(card.content);
      setSelectedTags(card.tags);
    } else {
      setTitle('');
      setContent('');
      setSelectedTags([]);
    }
    setErrors([]);
  }, [card]);

  const handleTagToggle = (tagId: string) => {
    setSelectedTags((prev) =>
      prev.includes(tagId) ? prev.filter((id) => id !== tagId) : [...prev, tagId]
    );
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const trimmedTitle = title.trim();
    const trimmedContent = content.trim();

    const validation = validateCard({
      title: trimmedTitle,
      content: trimmedContent,
      tags: selectedTags,
      linkedCardIds: card?.linkedCardIds || [],
    });

    if (!validation.isValid) {
      setErrors(validation.errors);
      showToast(validation.errors[0], 'error');
      return;
    }

    const now = Date.now();
    const newCard: Card = card
      ? {
          ...card,
          title: trimmedTitle,
          content: trimmedContent,
          tags: selectedTags,
          updatedAt: now,
        }
      : {
          id: uuidv4(),
          title: trimmedTitle,
          content: trimmedContent,
          tags: selectedTags,
          createdAt: now,
          updatedAt: now,
          linkedCardIds: [],
          review: createInitialReviewInfo(),
          order: 0,
        };

    onSubmit(newCard);
    showToast(card ? '卡片更新成功' : '卡片创建成功', 'success');
  };

  return (
    <div className="card-form">
      <h3>{card ? '编辑卡片' : '新建卡片'}</h3>
      <form onSubmit={handleSubmit}>
        {errors.length > 0 && (
          <div className="form-errors">
            {errors.map((error, index) => (
              <p key={index} className="error-message">
                {error}
              </p>
            ))}
          </div>
        )}

        <div className="form-group">
          <label htmlFor="title">标题 *</label>
          <input
            type="text"
            id="title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="输入卡片标题..."
            maxLength={200}
          />
          <span className="char-count">{title.length}/200</span>
        </div>

        <div className="form-group">
          <label htmlFor="content">内容 * (支持 Markdown)</label>
          <textarea
            id="content"
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder="输入卡片内容，支持 Markdown 格式..."
            rows={10}
          />
        </div>

        <div className="form-group">
          <label>标签</label>
          <div className="tag-selector">
            {tags.map((tag) => (
              <button
                key={tag.id}
                type="button"
                className={`tag-pill ${selectedTags.includes(tag.id) ? 'selected' : ''}`}
                style={{ backgroundColor: selectedTags.includes(tag.id) ? tag.color : 'transparent', borderColor: tag.color }}
                onClick={() => handleTagToggle(tag.id)}
              >
                {tag.name}
              </button>
            ))}
            {tags.length === 0 && <p className="text-muted">还没有标签，请到标签页面创建</p>}
          </div>
        </div>

        <div className="form-actions">
          <button type="button" className="btn btn-secondary" onClick={onCancel}>
            取消
          </button>
          <button type="submit" className="btn btn-primary">
            {card ? '保存' : '创建'}
          </button>
        </div>
      </form>
    </div>
  );
};

export default CardForm;
