import React, { useState, useEffect } from 'react';
import api from '../api/client';

interface Stats {
  total_orders: number;
  total_revenue: number;
  pending_orders: number;
  producing_orders: number;
  completed_orders: number;
  urgent_orders: number;
  total_customers: number;
  capacity_usage: number;
  recent_orders: any[];
}

const statCards = [
  { key: 'total_orders', label: '订单总数', color: '#3b82f6', icon: '📋' },
  { key: 'total_revenue', label: '总营收', color: '#10b981', icon: '💰', format: (v: number) => `¥${v.toFixed(2)}` },
  { key: 'pending_orders', label: '待处理', color: '#f59e0b', icon: '⏳' },
  { key: 'producing_orders', label: '生产中', color: '#6366f1', icon: '🏭' },
  { key: 'completed_orders', label: '已完成', color: '#14b8a6', icon: '✅' },
  { key: 'urgent_orders', label: '加急单', color: '#ef4444', icon: '⚡' },
  { key: 'total_customers', label: '客户数', color: '#8b5cf6', icon: '👥' },
  { key: 'capacity_usage', label: '产能利用率', color: '#ec4899', icon: '📈', format: (v: number) => `${v.toFixed(1)}%` },
];

export default function Dashboard() {
  const [stats, setStats] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get('/dashboard/stats').then(({ data }) => {
      setStats(data);
      setLoading(false);
    });
  }, []);

  if (loading) return <div>加载中...</div>;
  if (!stats) return <div>无数据</div>;

  return (
    <div>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 16, marginBottom: 24 }}>
        {statCards.map((card) => {
          const value = (stats as any)[card.key];
          const display = card.format ? card.format(value) : value;
          return (
            <div key={card.key} style={{
              background: '#fff',
              padding: 20,
              borderRadius: 8,
              border: `1px solid ${card.color}20`,
              borderLeft: `4px solid ${card.color}`,
            }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                  <div style={{ fontSize: 13, color: '#6b7280', marginBottom: 6 }}>{card.label}</div>
                  <div style={{ fontSize: 24, fontWeight: 600, color: card.color }}>{display}</div>
                </div>
                <div style={{ fontSize: 28 }}>{card.icon}</div>
              </div>
            </div>
          );
        })}
      </div>

      <div style={{ background: '#fff', padding: 20, borderRadius: 8 }}>
        <h3 style={{ marginBottom: 16, fontSize: 16, fontWeight: 600 }}>最近订单</h3>
        <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: 14 }}>
          <thead>
            <tr style={{ borderBottom: '1px solid #e5e7eb' }}>
              {['订单号', '客户', '金额', '状态', '加急', '创建时间'].map(h => (
                <th key={h} style={{ textAlign: 'left', padding: '10px 8px', color: '#6b7280', fontWeight: 500 }}>{h}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {stats.recent_orders.map((o: any) => (
              <tr key={o.id} style={{ borderBottom: '1px solid #f3f4f6' }}>
                <td style={{ padding: '10px 8px' }}>{o.order_no}</td>
                <td style={{ padding: '10px 8px' }}>{o.customer?.name || '-'}</td>
                <td style={{ padding: '10px 8px' }}>¥{o.final_price.toFixed(2)}</td>
                <td style={{ padding: '10px 8px' }}>
                  <span style={{
                    padding: '2px 8px',
                    borderRadius: 4,
                    fontSize: 12,
                    background: statusColor(o.status) + '20',
                    color: statusColor(o.status),
                  }}>{statusLabel(o.status)}</span>
                </td>
                <td style={{ padding: '10px 8px' }}>{o.urgent ? '⚡ 加急' : '-'}</td>
                <td style={{ padding: '10px 8px', color: '#6b7280' }}>{new Date(o.created_at).toLocaleString()}</td>
              </tr>
            ))}
            {stats.recent_orders.length === 0 && (
              <tr><td colSpan={6} style={{ textAlign: 'center', padding: 20, color: '#9ca3af' }}>暂无订单</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}

function statusColor(s: string): string {
  const map: Record<string, string> = {
    created: '#f59e0b', reviewing: '#3b82f6', reviewed: '#6366f1',
    producing: '#6366f1', produced: '#14b8a6', shipped: '#10b981', cancelled: '#ef4444',
  };
  return map[s] || '#6b7280';
}
function statusLabel(s: string): string {
  const map: Record<string, string> = {
    created: '待处理', reviewing: '审核中', reviewed: '已审核',
    producing: '生产中', produced: '已生产', shipped: '已发货', cancelled: '已取消',
  };
  return map[s] || s;
}
