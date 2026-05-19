import type { ReviewInfo } from '../types';

const MIN_EASE_FACTOR = 1.3;
const DEFAULT_EASE_FACTOR = 2.5;
const MIN_INTERVAL = 1;

export const createInitialReviewInfo = (): ReviewInfo => {
  const now = Date.now();
  return {
    nextReviewAt: now,
    interval: MIN_INTERVAL,
    easeFactor: DEFAULT_EASE_FACTOR,
    repetitions: 0,
    lastReviewedAt: null,
  };
};

export const calculateNextReview = (
  current: ReviewInfo,
  quality: number
): ReviewInfo => {
  let { interval, easeFactor, repetitions } = current;
  const now = Date.now();

  if (quality < 3) {
    repetitions = 0;
    interval = MIN_INTERVAL;
  } else {
    if (repetitions === 0) {
      interval = 1;
    } else if (repetitions === 1) {
      interval = 6;
    } else {
      interval = Math.round(interval * easeFactor);
    }
    repetitions += 1;
  }

  easeFactor = easeFactor + (0.1 - (5 - quality) * (0.08 + (5 - quality) * 0.02));
  if (easeFactor < MIN_EASE_FACTOR) {
    easeFactor = MIN_EASE_FACTOR;
  }

  const nextReviewAt = now + interval * 24 * 60 * 60 * 1000;

  return {
    nextReviewAt,
    interval,
    easeFactor,
    repetitions,
    lastReviewedAt: now,
  };
};

export const isOverdue = (review: ReviewInfo): boolean => {
  return Date.now() > review.nextReviewAt;
};

export const getReviewStatus = (review: ReviewInfo): 'overdue' | 'due' | 'upcoming' => {
  const now = Date.now();
  const oneDay = 24 * 60 * 60 * 1000;

  if (now > review.nextReviewAt) {
    return 'overdue';
  } else if (now + oneDay > review.nextReviewAt) {
    return 'due';
  }
  return 'upcoming';
};

export const formatNextReview = (timestamp: number): string => {
  const now = Date.now();
  const diff = timestamp - now;
  const days = Math.ceil(diff / (24 * 60 * 60 * 1000));

  if (days < 0) {
    return `已逾期 ${Math.abs(days)} 天`;
  } else if (days === 0) {
    return '今天';
  } else if (days === 1) {
    return '明天';
  } else {
    return `${days} 天后`;
  }
};
