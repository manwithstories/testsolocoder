import React, { useState } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import type { Card, Tag } from '../../types';
import { formatNextReview, getReviewStatus } from '../../utils/spacedRepetition';
import { useToast } from '../common/Toast';

interface CardDetailProps {
  card: Card;
  tags: Tag[];
  allCards: Card[];
  onClose: () => void;
  onEdit: () => void;
  onDelete: () => void;
  onLinkCard: (cardId: string) => void;
  onUnlinkCard: (cardId: string) => void;
  onSelectCard: (cardId: string) => void;
  onReview: (quality: number) => void;
}

const CardDetail: React.FC<CardDetailProps> = ({
  card,
  tags,
  allCards,
  onClose,
  onEdit,
  onDelete,
  onLinkCard,
  onUnlinkCard,
  onSelectCard,
  onReview,
}) => {
  const [showReview, setShowReview] = useState(false);
  const [linkSearch, setLinkSearch] = useState('');
  const { showToast } = useToast();

  const cardTags = card.tags.map((tagId) => tags.find((t) => t.id === tagId)).filter(Boolean) as Tag[];
  const linkedCards = card.linkedCardIds
    .map((id) => allCards.find((c) => c.id === id))
    .filter(Boolean) as Card[];

  const availableToLink = allCards.filter(
    (c) => c.id !== card.id && !card.linkedCardIds.includes(c.id)
  );

  const filteredToLink = availableToLink.filter((c) =>
    c.title.toLowerCase().includes(linkSearch.toLowerCase())
  );

  const reviewStatus = getReviewStatus(card.review);

  const handleReview = (quality: number) => {
    onReview(quality);
    setShowReview(false);
    const messages = [
      '完全忘记了',
      '记错了',
      '想起来了',
      '记得',
      '很简单',
    ];
    showToast(`复习完成: ${messages[quality]}`, 'success');
  };

  return (
    <div className="card-detail">
      <div className="detail-header">
        <button className="icon-btn" onClick={onClose}>
          ← 返回
        </button>
        <div className="detail-actions">
          <button className="btn btn-secondary" onClick={onEdit}>
            编辑
          </button>
          <button className="btn btn-danger" onClick={onDelete}>
            删除
          </button>
        </div>
      </div>

      <div className="detail-content">
        <h2>{card.title}</h2>

        <div className="detail-meta">
          <div className="detail-tags">
            {cardTags.map((tag) => (
              <span key={tag.id} className="tag-badge" style={{ backgroundColor: tag.color }}>
                {tag.name}
              </span>
            ))}
          </div>
          <div className="detail-dates">
            <span>创建: {new Date(card.createdAt).toLocaleString()}</span>
            <span>更新: {new Date(card.updatedAt).toLocaleString()}</span>
          </div>
        </div>

        <div className={`review-status review-${reviewStatus}`}>
          <span>📅 下次复习: {formatNextReview(card.review.nextReviewAt)}</span>
          <button className="btn btn-primary btn-sm" onClick={() => setShowReview(!showReview)}>
            开始复习
          </button>
        </div>

        {showReview && (
          <div className="review-panel">
            <p>你对这张卡片的记忆程度如何？</p>
            <div className="review-buttons">
              <button className="btn btn-danger" onClick={() => handleReview(0)}>0 - 完全忘记</button>
              <button className="btn btn-warning" onClick={() => handleReview(1)}>1 - 记错了</button>
              <button className="btn btn-warning" onClick={() => handleReview(2)}>2 - 想起来了</button>
              <button className="btn btn-info" onClick={() => handleReview(3)}>3 - 记得</button>
              <button className="btn btn-success" onClick={() => handleReview(4)}>4 - 很简单</button>
            </div>
          </div>
        )}

        <div className="markdown-content">
          <ReactMarkdown remarkPlugins={[remarkGfm]}>
            {card.content}
          </ReactMarkdown>
        </div>

        <div className="linked-cards-section">
          <h3>关联卡片 ({linkedCards.length})</h3>
          {linkedCards.length > 0 ? (
            <div className="linked-cards-list">
              {linkedCards.map((linkedCard) => (
                <div key={linkedCard.id} className="linked-card-item">
                  <span onClick={() => onSelectCard(linkedCard.id)} className="linked-card-title">
                    {linkedCard.title}
                  </span>
                  <button
                    className="icon-btn"
                    onClick={() => onUnlinkCard(linkedCard.id)}
                    title="取消关联"
                  >
                    ✕
                  </button>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-muted">还没有关联卡片</p>
          )}

          <div className="add-link-section">
            <h4>添加关联</h4>
            <input
              type="text"
              placeholder="搜索卡片..."
              value={linkSearch}
              onChange={(e) => setLinkSearch(e.target.value)}
            />
            {linkSearch && filteredToLink.length > 0 && (
              <div className="link-search-results">
                {filteredToLink.slice(0, 5).map((c) => (
                  <div key={c.id} className="link-search-item" onClick={() => {
                    onLinkCard(c.id);
                    setLinkSearch('');
                    showToast('关联已添加', 'success');
                  }}>
                    {c.title}
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default CardDetail;
