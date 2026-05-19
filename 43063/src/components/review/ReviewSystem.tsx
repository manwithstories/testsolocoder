import React, { useState, useMemo } from 'react';
import type { Card, Tag } from '../../types';
import { getReviewStatus, formatNextReview, isOverdue } from '../../utils/spacedRepetition';
import { useToast } from '../common/Toast';

interface ReviewSystemProps {
  cards: Card[];
  tags: Tag[];
  onReviewCard: (cardId: string, quality: number) => void;
  onSelectCard: (cardId: string) => void;
}

const ReviewSystem: React.FC<ReviewSystemProps> = ({
  cards,
  tags,
  onReviewCard,
  onSelectCard,
}) => {
  const [showAnswer, setShowAnswer] = useState(false);
  const [currentIndex, setCurrentIndex] = useState(0);
  const { showToast } = useToast();

  const reviewCards = useMemo(() => {
    return cards
      .filter((card) => isOverdue(card.review))
      .sort((a, b) => a.review.nextReviewAt - b.review.nextReviewAt);
  }, [cards]);

  const dueCards = useMemo(() => {
    return cards.filter((card) => {
      const status = getReviewStatus(card.review);
      return status === 'due';
    });
  }, [cards]);

  const upcomingCards = useMemo(() => {
    return cards.filter((card) => {
      const status = getReviewStatus(card.review);
      return status === 'upcoming';
    });
  }, [cards]);

  const currentCard = reviewCards[currentIndex];

  const handleReview = (quality: number) => {
    if (currentCard) {
      onReviewCard(currentCard.id, quality);
      setShowAnswer(false);
      const messages = ['完全忘记了', '记错了', '想起来了', '记得', '很简单'];
      showToast(`复习完成: ${messages[quality]}`, 'success');

      if (currentIndex < reviewCards.length - 1) {
        setCurrentIndex(currentIndex + 1);
      } else {
        setCurrentIndex(0);
      }
    }
  };

  const getTagNames = (tagIds: string[]): string[] => {
    return tagIds.map((id) => tags.find((t) => t.id === id)?.name || '').filter(Boolean);
  };

  return (
    <div className="review-system">
      <div className="list-header">
        <h3>复习系统</h3>
      </div>

      <div className="review-stats">
        <div className="stat-card stat-overdue">
          <span className="stat-number">{reviewCards.length}</span>
          <span className="stat-label">待复习</span>
        </div>
        <div className="stat-card stat-due">
          <span className="stat-number">{dueCards.length}</span>
          <span className="stat-label">今天到期</span>
        </div>
        <div className="stat-card stat-upcoming">
          <span className="stat-number">{upcomingCards.length}</span>
          <span className="stat-label">即将到期</span>
        </div>
        <div className="stat-card stat-total">
          <span className="stat-number">{cards.length}</span>
          <span className="stat-label">总卡片</span>
        </div>
      </div>

      {reviewCards.length === 0 ? (
        <div className="empty-state">
          <p>🎉 太棒了！现在没有需要复习的卡片</p>
          <p className="text-muted">继续添加新知识，或者去浏览已有的卡片吧</p>
        </div>
      ) : (
        <div className="review-session">
          <div className="review-progress">
            进度: {currentIndex + 1} / {reviewCards.length}
          </div>

          {currentCard && (
            <div className="review-card">
              <div className="review-card-header">
                <h4>{currentCard.title}</h4>
                <div className="review-card-tags">
                  {getTagNames(currentCard.tags).map((tagName, index) => (
                    <span key={index} className="tag-badge">
                      {tagName}
                    </span>
                  ))}
                </div>
              </div>

              <div className="review-card-content">
                {!showAnswer ? (
                  <div className="question-preview">
                    <p>{currentCard.content.slice(0, 200)}...</p>
                    <button className="btn btn-primary" onClick={() => setShowAnswer(true)}>
                      显示答案
                    </button>
                  </div>
                ) : (
                  <div className="answer-content">
                    <p>{currentCard.content}</p>
                    <div className="review-buttons">
                      <p>你对这张卡片的记忆程度如何？</p>
                      <div className="button-row">
                        <button className="btn btn-danger" onClick={() => handleReview(0)}>
                          0<br />
                          <small>完全忘记</small>
                        </button>
                        <button className="btn btn-warning" onClick={() => handleReview(1)}>
                          1<br />
                          <small>记错了</small>
                        </button>
                        <button className="btn btn-warning" onClick={() => handleReview(2)}>
                          2<br />
                          <small>想起来</small>
                        </button>
                        <button className="btn btn-info" onClick={() => handleReview(3)}>
                          3<br />
                          <small>记得</small>
                        </button>
                        <button className="btn btn-success" onClick={() => handleReview(4)}>
                          4<br />
                          <small>很简单</small>
                        </button>
                      </div>
                    </div>
                  </div>
                )}
              </div>

              <div className="review-card-footer">
                <span className="text-muted">
                  逾期: {formatNextReview(currentCard.review.nextReviewAt)}
                </span>
                <button
                  className="btn btn-secondary btn-sm"
                  onClick={() => onSelectCard(currentCard.id)}
                >
                  查看详情
                </button>
              </div>
            </div>
          )}
        </div>
      )}

      <div className="review-list">
        <h4>待复习列表</h4>
        {reviewCards.length === 0 ? (
          <p className="text-muted">没有待复习的卡片</p>
        ) : (
          <div className="review-list-items">
            {reviewCards.map((card) => (
              <div
                key={card.id}
                className={`review-list-item ${card.id === currentCard?.id ? 'active' : ''}`}
                onClick={() => {
                  setCurrentIndex(reviewCards.indexOf(card));
                  setShowAnswer(false);
                }}
              >
                <span className="review-item-title">{card.title}</span>
                <span className="review-item-date overdue">
                  {formatNextReview(card.review.nextReviewAt)}
                </span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default ReviewSystem;
