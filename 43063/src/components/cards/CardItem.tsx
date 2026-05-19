import React from 'react';
import { useDrag, useDrop } from 'react-dnd';
import type { Card, Tag } from '../../types';
import { isOverdue, formatNextReview } from '../../utils/spacedRepetition';
import HighlightText from '../common/HighlightText';

interface CardItemProps {
  card: Card;
  tags: Tag[];
  searchQuery: string;
  isSelected: boolean;
  onSelect: () => void;
  onEdit: () => void;
  onDelete: () => void;
  onMoveCard: (sourceId: string, targetId: string) => void;
}

const CardItem: React.FC<CardItemProps> = ({
  card,
  tags,
  searchQuery,
  isSelected,
  onSelect,
  onEdit,
  onDelete,
  onMoveCard,
}) => {
  const [{ isDragging }, drag] = useDrag(() => ({
    type: 'CARD',
    item: { id: card.id },
    collect: (monitor) => ({
      isDragging: monitor.isDragging(),
    }),
  }));

  const [{ isOver }, drop] = useDrop(() => ({
    accept: 'CARD',
    drop: (item: { id: string }) => {
      if (item.id !== card.id) {
        onMoveCard(item.id, card.id);
      }
    },
    collect: (monitor) => ({
      isOver: monitor.isOver(),
    }),
  }));

  const cardTags = card.tags.map((tagId) => tags.find((t) => t.id === tagId)).filter(Boolean) as Tag[];
  const overdue = isOverdue(card.review);

  return (
    <div
      ref={(node) => drag(drop(node))}
      className={`card-item ${isSelected ? 'selected' : ''} ${overdue ? 'overdue' : ''} ${isDragging ? 'dragging' : ''} ${isOver ? 'drop-target' : ''}`}
      onClick={onSelect}
    >
      <div className="card-header">
        <h4 className="card-title"><HighlightText text={card.title} query={searchQuery} /></h4>
        <div className="card-actions">
          <button className="icon-btn" onClick={(e) => { e.stopPropagation(); onEdit(); }}>
            ✏️
          </button>
          <button className="icon-btn" onClick={(e) => { e.stopPropagation(); onDelete(); }}>
            🗑️
          </button>
        </div>
      </div>

      <div className="card-tags">
        {cardTags.map((tag) => (
          <span key={tag.id} className="tag-badge" style={{ backgroundColor: tag.color }}>
            <HighlightText text={tag.name} query={searchQuery} />
          </span>
        ))}
      </div>

      <p className="card-preview">
        <HighlightText text={card.content.slice(0, 100)} query={searchQuery} />
        {card.content.length > 100 && '...'}
      </p>

      <div className="card-footer">
        <span className="card-date">
          {new Date(card.updatedAt).toLocaleDateString()}
        </span>
        <span className={`card-review ${overdue ? 'overdue' : ''}`}>
          📅 {formatNextReview(card.review.nextReviewAt)}
        </span>
        {card.linkedCardIds.length > 0 && (
          <span className="card-links">🔗 {card.linkedCardIds.length}</span>
        )}
      </div>
    </div>
  );
};

export default CardItem;
