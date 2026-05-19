const TAG_COLORS = [
  '#FF6B6B',
  '#4ECDC4',
  '#45B7D1',
  '#96CEB4',
  '#FFEAA7',
  '#DDA0DD',
  '#98D8C8',
  '#F7DC6F',
  '#BB8FCE',
  '#85C1E9',
  '#F8B500',
  '#00CED1',
  '#FF8C00',
  '#9370DB',
  '#3CB371',
];

export const getRandomTagColor = (): string => {
  return TAG_COLORS[Math.floor(Math.random() * TAG_COLORS.length)];
};

export const getTagColorByIndex = (index: number): string => {
  return TAG_COLORS[index % TAG_COLORS.length];
};

export const generateConsistentColor = (seed: string): string => {
  let hash = 0;
  for (let i = 0; i < seed.length; i++) {
    hash = seed.charCodeAt(i) + ((hash << 5) - hash);
  }
  const index = Math.abs(hash % TAG_COLORS.length);
  return TAG_COLORS[index];
};
