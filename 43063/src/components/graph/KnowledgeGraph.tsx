import React, { useRef, useCallback, useState, useEffect } from 'react';
import ForceGraph2D from 'react-force-graph-2d';
import type { Card, Tag } from '../../types';
import { useToast } from '../common/Toast';

interface KnowledgeGraphProps {
  cards: Card[];
  tags: Tag[];
  selectedCardId: string | null;
  onSelectCard: (id: string | null) => void;
}

interface GraphNode {
  id: string;
  name: string;
  val: number;
  color: string;
  type: 'card' | 'tag';
  card?: Card;
}

interface GraphLink {
  source: string;
  target: string;
  value: number;
}

const KnowledgeGraph: React.FC<KnowledgeGraphProps> = ({
  cards,
  tags,
  selectedCardId,
  onSelectCard,
}) => {
  const graphRef = useRef<any>();
  const { showToast } = useToast();
  const [dimensions, setDimensions] = useState({ width: 800, height: 600 });
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const updateDimensions = () => {
      if (containerRef.current) {
        setDimensions({
          width: containerRef.current.clientWidth,
          height: Math.max(600, window.innerHeight - 200),
        });
      }
    };
    updateDimensions();
    window.addEventListener('resize', updateDimensions);
    return () => window.removeEventListener('resize', updateDimensions);
  }, []);

  const { nodes, links } = React.useMemo(() => {
    const graphNodes: GraphNode[] = [];
    const graphLinks: GraphLink[] = [];
    const tagCardMap = new Map<string, string[]>();

    cards.forEach((card) => {
      const linkCount = card.linkedCardIds.length;
      const size = Math.max(5, Math.min(15, 8 + linkCount));

      let nodeColor = '#4A90D9';
      if (card.tags.length > 0) {
        const firstTag = tags.find((t) => t.id === card.tags[0]);
        if (firstTag) {
          nodeColor = firstTag.color;
        }
      }

      graphNodes.push({
        id: card.id,
        name: card.title,
        val: size,
        color: nodeColor,
        type: 'card',
        card,
      });

      card.tags.forEach((tagId) => {
        if (!tagCardMap.has(tagId)) {
          tagCardMap.set(tagId, []);
        }
        tagCardMap.get(tagId)!.push(card.id);
      });

      card.linkedCardIds.forEach((linkedId) => {
        if (card.id < linkedId) {
          graphLinks.push({
            source: card.id,
            target: linkedId,
            value: 2,
          });
        }
      });
    });

    tagCardMap.forEach((cardIds, tagId) => {
      const tag = tags.find((t) => t.id === tagId);
      if (tag && cardIds.length > 1) {
        graphNodes.push({
          id: `tag-${tagId}`,
          name: tag.name,
          val: Math.max(10, cardIds.length * 2),
          color: tag.color,
          type: 'tag',
        });

        cardIds.forEach((cardId) => {
          graphLinks.push({
            source: `tag-${tagId}`,
            target: cardId,
            value: 1,
          });
        });
      }
    });

    return { nodes: graphNodes, links: graphLinks };
  }, [cards, tags]);

  const handleNodeClick = useCallback(
    (node: GraphNode) => {
      if (node.type === 'card' && node.card) {
        onSelectCard(node.card.id);
      }
    },
    [onSelectCard]
  );

  const handleNodeHover = useCallback((node: GraphNode | null) => {
    if (node) {
      document.body.style.cursor = 'pointer';
    } else {
      document.body.style.cursor = 'default';
    }
  }, []);

  const handleExportImage = () => {
    if (graphRef.current) {
      const canvas = containerRef.current?.querySelector('canvas');
      if (canvas) {
        const link = document.createElement('a');
        link.download = 'knowledge-graph.png';
        link.href = canvas.toDataURL('image/png');
        link.click();
        showToast('图谱已导出', 'success');
      }
    }
  };

  return (
    <div className="knowledge-graph">
      <div className="list-header">
        <h3>
          知识图谱
          <span className="count-badge">{cards.length} 张卡片</span>
        </h3>
        <button className="btn btn-secondary" onClick={handleExportImage}>
          导出图片
        </button>
      </div>

      <div className="graph-container" ref={containerRef}>
        {cards.length === 0 ? (
          <div className="empty-state">
            <p>还没有卡片，无法生成知识图谱</p>
          </div>
        ) : (
          <ForceGraph2D
            ref={graphRef}
            graphData={{ nodes, links }}
            width={dimensions.width}
            height={dimensions.height}
            nodeColor={(node) => {
              const n = node as GraphNode;
              if (n.id === selectedCardId) {
                return '#FFD700';
              }
              return n.color;
            }}
            nodeLabel={(node) => (node as GraphNode).name}
            nodeVal={(node) => (node as GraphNode).val}
            linkColor={() => 'rgba(100, 100, 100, 0.4)'}
            linkWidth={(link) => (link as GraphLink).value}
            onNodeClick={handleNodeClick}
            onNodeHover={handleNodeHover}
            enableNodeDrag={true}
            backgroundColor="transparent"
            linkDirectionalParticles={0}
            cooldownTicks={100}
          />
        )}
      </div>

      <div className="graph-legend">
        <div className="legend-item">
          <span className="legend-dot" style={{ backgroundColor: '#4A90D9' }} />
          <span>知识卡片</span>
        </div>
        <div className="legend-item">
          <span className="legend-dot" style={{ backgroundColor: '#FFD700' }} />
          <span>选中卡片</span>
        </div>
        <div className="legend-item">
          <span className="legend-dot" style={{ backgroundColor: '#96CEB4' }} />
          <span>标签分组</span>
        </div>
      </div>
    </div>
  );
};

export default KnowledgeGraph;
