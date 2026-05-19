import React, { useState, useMemo } from 'react';
import { DndProvider } from 'react-dnd';
import { HTML5Backend } from 'react-dnd-html5-backend';
import CardItem from './CardItem';
import CardForm from './CardForm';
import CardDetail from './CardDetail';
import ConfirmModal from '../common/ConfirmModal';
import type { Card, Tag } from '../../types';
import { searchCards } from '../../utils/search';

interface CardListProps {
  cards: Card[];
  tags: Tag[];
  selectedTagIds: string[];
  searchQuery: string;
  selectedCardId: string | null;
  onSelectCard: (id: string | null) => void;
  onAddCard: (card: Card) => void;
  onUpdateCard: (card: Card) => void;
  onDeleteCard: (id: string) => void;
  onLinkCards: (cardId1: string, cardId2: string) => void;
  onUnlinkCards: (cardId1: string, cardId2: string) => void;
  onReviewCard: (cardId: string, quality: number) => void;
  onReorderCards: (sourceId: string, targetId: string) => void;
}

const CardList: React.FC<CardListProps> = ({
  cards,
  tags,
  selectedTagIds,
  searchQuery,
  selectedCardId,
  onSelectCard,
  onAddCard,
  onUpdateCard,
  onDeleteCard,
  onLinkCards,
  onUnlinkCards,
  onReviewCard,
  onReorderCards,
}) => {
  const [showForm, setShowForm] = useState(false);
  const [editingCard, setEditingCard] = useState<Card | null>(null);
  const [deleteConfirm, setDeleteConfirm] = useState<string | null>(null);

  const filteredCards = useMemo(() => {
    let result = [...cards];

    if (selectedTagIds.length > 0) {
      result = result.filter((card) =>
        selectedTagIds.some((tagId) => card.tags.includes(tagId))
      );
    }

    if (searchQuery.trim()) {
      const searchResults = searchCards(result, tags, searchQuery);
      result = searchResults.map((r) => r.card);
    }

    result.sort((a, b) => a.order - b.order);

    return result;
  }, [cards, tags, selectedTagIds, searchQuery]);

  const selectedCard = selectedCardId ? cards.find((c) => c.id === selectedCardId) : null;

  const handleEdit = (card: Card) => {
    setEditingCard(card);
    setShowForm(true);
  };

  const handleDelete = (id: string) => {
    setDeleteConfirm(id);
  };

  const confirmDelete = () => {
    if (deleteConfirm) {
      onDeleteCard(deleteConfirm);
      setDeleteConfirm(null);
    }
  };

  const handleFormSubmit = (card: Card) => {
    if (editingCard) {
      onUpdateCard(card);
    } else {
      onAddCard(card);
    }
    setShowForm(false);
    setEditingCard(null);
  };

  const handleFormCancel = () => {
    setShowForm(false);
    setEditingCard(null);
  };

  if (selectedCard) {
    return (
      <CardDetail
        card={selectedCard}
        tags={tags}
        allCards={cards}
        onClose={() => onSelectCard(null)}
        onEdit={() => {
          onSelectCard(null);
          handleEdit(selectedCard);
        }}
        onDelete={() => handleDelete(selectedCard.id)}
        onLinkCard={(cardId) => onLinkCards(selectedCard.id, cardId)}
        onUnlinkCard={(cardId) => onUnlinkCards(selectedCard.id, cardId)}
        onSelectCard={onSelectCard}
        onReview={(quality) => onReviewCard(selectedCard.id, quality)}
      />
    );
  }

  if (showForm) {
    return (
      <CardForm
        card={editingCard}
        tags={tags}
        onSubmit={handleFormSubmit}
        onCancel={handleFormCancel}
      />
    );
  }

  return (
    <DndProvider backend={HTML5Backend}>
      <div className="card-list">
        <div className="list-header">
          <h3>
            知识卡片
            <span className="count-badge">{filteredCards.length}</span>
          </h3>
          <button className="btn btn-primary" onClick={() => setShowForm(true)}>
            + 新建卡片
          </button>
        </div>

        {filteredCards.length === 0 ? (
          <div className="empty-state">
            <p>
              {searchQuery
                ? '没有找到匹配的卡片'
                : selectedTagIds.length > 0
                ? '该标签下没有卡片'
                : '还没有卡片，点击上方按钮创建第一张'}
            </p>
          </div>
        ) : (
          <div className="cards-grid">
            {filteredCards.map((card) => (
              <CardItem
                key={card.id}
                card={card}
                tags={tags}
                searchQuery={searchQuery}
                isSelected={card.id === selectedCardId}
                onSelect={() => onSelectCard(card.id)}
                onEdit={() => handleEdit(card)}
                onDelete={() => handleDelete(card.id)}
                onMoveCard={onReorderCards}
              />
            ))}
          </div>
        )}

        <ConfirmModal
          isOpen={!!deleteConfirm}
          title="删除卡片"
          message="确定要删除这张卡片吗？此操作无法撤销。"
          confirmText="删除"
          variant="danger"
          onConfirm={confirmDelete}
          onCancel={() => setDeleteConfirm(null)}
        />
      </div>
    </DndProvider>
  );
};

export default CardList;
