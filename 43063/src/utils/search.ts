import type { Card, Tag } from '../types';

export interface SearchResult {
  card: Card;
  score: number;
  matches: {
    title: number[];
    content: number[];
    tags: number[];
  };
}

export const searchCards = (
  cards: Card[],
  tags: Tag[],
  query: string
): SearchResult[] => {
  if (!query.trim()) {
    return [];
  }

  const lowerQuery = query.toLowerCase();
  const tagMap = new Map(tags.map((t) => [t.id, t.name.toLowerCase()]));

  return cards
    .map((card): SearchResult | null => {
      const lowerTitle = card.title.toLowerCase();
      const lowerContent = card.content.toLowerCase();
      const cardTagNames = card.tags
        .map((tagId) => tagMap.get(tagId) || '')
        .join(' ');

      const titleMatches = findAllMatches(lowerTitle, lowerQuery);
      const contentMatches = findAllMatches(lowerContent, lowerQuery);
      const tagMatches = findAllMatches(cardTagNames, lowerQuery);

      const hasMatch =
        titleMatches.length > 0 || contentMatches.length > 0 || tagMatches.length > 0;

      if (!hasMatch) {
        return null;
      }

      const score =
        titleMatches.length * 3 + contentMatches.length * 1 + tagMatches.length * 2;

      return {
        card,
        score,
        matches: {
          title: titleMatches,
          content: contentMatches,
          tags: tagMatches,
        },
      };
    })
    .filter((result): result is SearchResult => result !== null)
    .sort((a, b) => b.score - a.score);
};

const findAllMatches = (text: string, query: string): number[] => {
  const indices: number[] = [];
  let index = text.indexOf(query);
  while (index !== -1) {
    indices.push(index);
    index = text.indexOf(query, index + 1);
  }
  return indices;
};

export interface HighlightSegment {
  text: string;
  isMatch: boolean;
}

export const getHighlightSegments = (text: string, query: string): HighlightSegment[] => {
  if (!query.trim()) {
    return [{ text, isMatch: false }];
  }

  const lowerQuery = query.toLowerCase();
  const lowerText = text.toLowerCase();
  const segments: HighlightSegment[] = [];
  let lastIndex = 0;
  let index = lowerText.indexOf(lowerQuery);

  while (index !== -1) {
    if (index > lastIndex) {
      segments.push({
        text: text.slice(lastIndex, index),
        isMatch: false,
      });
    }
    segments.push({
      text: text.slice(index, index + query.length),
      isMatch: true,
    });
    lastIndex = index + query.length;
    index = lowerText.indexOf(lowerQuery, lastIndex);
  }

  if (lastIndex < text.length) {
    segments.push({
      text: text.slice(lastIndex),
      isMatch: false,
    });
  }

  return segments.length > 0 ? segments : [{ text, isMatch: false }];
};
