import React from 'react';
import { getHighlightSegments } from '../../utils/search';

interface HighlightTextProps {
  text: string;
  query: string;
}

const HighlightText: React.FC<HighlightTextProps> = ({ text, query }) => {
  const segments = getHighlightSegments(text, query);

  return (
    <>
      {segments.map((segment, index) =>
        segment.isMatch ? (
          <mark key={index} className="search-highlight">
            {segment.text}
          </mark>
        ) : (
          <span key={index}>{segment.text}</span>
        )
      )}
    </>
  );
};

export default HighlightText;
