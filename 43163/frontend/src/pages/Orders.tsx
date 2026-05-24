import React, { useState, useEffect } from 'react';
import api from '../api/client';

const statuses = ['created', 'reviewing', 'reviewed', 'producing', 'produced', 'shipped', 'cancelled'];

const statusLabels: Record<string, string> = {
  created: '待处理', reviewing: '审核中', reviewed: '已审核',
  producing: '生产中', produced: '已生产', shipped: '已发货', cancelled: '已取消',
};
const statusColors: Record<string, string> = {
  created: '#f59e0b', reviewing: '#3b82f6', reviewed: '#6366f1',
  producing: '#6366f1', produced: '#14b8a6', shipped: '#10b981', cancelled: '#ef4444',
};

export default function Orders() {
  const [orders, setOrders] = useState<any[]>([]);
  const [filter, setFilter] = useState({ status: '', customer_id: '', urgent: '' });
  const [loading, setLoading] = useState(true);

  const load = () => {
    const params = new URLSearchParams();
    if (filter.status) params.set('status', filter.status);
    if (filter.customer_id) params.set('customer_id', filter.customer_id);
    if (filter.urgent) params.set('urgent', filter.urgent);
    api.get('/orders', { params }).then(({ data }) => {
      setOrders(data);
      setLoading(false);
    });
  };

  useEffect(() => { load(); }, [filter]);

  const updateStatus = (id: number, status: string) => {
    api.put(`/orders/${id}/status`, { status }).then(load);
  };

  const splitOrder = (id: number, itemIds: number[]) => {
    api.post(`/orders/${id}/split`, { item_ids: itemIds }).then(({ data }) => {
      alert(`已拆分为新订单: ${data.order_no}`);
      load();
    });
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <div style={{ fontSize: 18, fontWeight: 600 }}>订单管理</div>
        <div style={{ display: 'flex', gap: 8 }}>
          <select style={selectStyle} value={filter.status} onChange={(e) => setFilter({ ...filter, status: e.target.value })}>
            <option value="">全部状态</option>
            {statuses.map(s => <option key={s} value={s}>{statusLabels[s]}</option>)}
          </select>
          <select style={selectStyle} value={filter.urgent} onChange={(e) => setFilter({ ...filter, urgent: e.target.value })}>
            <option value="">全部</option>
            <option value="true">加急单</option>
          </select>
          <a href="/orders/new" style={primaryBtn}>+ 新建订单</a>
        </div>
      </div>

      {loading ? <div>加载中...</div> : (
        <div style={{ background: '#fff', borderRadius: 8, overflow: 'hidden' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: 14 }}>
            <thead style={{ background: '#f9fafb' }}>
              <tr>
                {['订单号', '客户', '金额', '状态', '加急', '创建时间', '操作'].map(h => (
                  <th key={h} style={{ textAlign: 'left', padding: '12px 12px', color: '#6b7280', fontWeight: 500 }}>{h}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {orders.map((o) => (
                <React.Fragment key={o.id}>
                  <tr style={{ borderBottom: '1px solid #f3f4f6', cursor: 'pointer' }} onClick={() => {
                    const el = document.getElementById(`order-${o.id}`);
                    if (el) el.style.display = el.style.display === 'none' ? 'table-row' : 'none';
                  }}>
                    <td style={{ padding: '12px' }}>{o.order_no}</td>
                    <td style={{ padding: '12px' }}>{o.customer?.name || '-'}</td>
                    <td style={{ padding: '12px' }}>¥{o.final_price.toFixed(2)}</td>
                    <td style={{ padding: '12px' }}>
                      <span style={{ padding: '3px 10px', borderRadius: 4, fontSize: 12, background: statusColors[o.status] + '20', color: statusColors[o.status] }}>{statusLabels[o.status]}</span>
                    </td>
                    <td style={{ padding: '12px' }}>{o.urgent ? '⚡' : '-'}</td>
                    <td style={{ padding: '12px', color: '#6b7280' }}>{new Date(o.created_at).toLocaleString()}</td>
                    <td style={{ padding: '12px' }}>
                      <select style={selectStyle} value={o.status} onChange={(e) => { e.stopPropagation(); updateStatus(o.id, e.target.value); }} onClick={(e) => e.stopPropagation()}>
                        {statuses.map(s => <option key={s} value={s}>{statusLabels[s]}</option>)}
                      </select>
                    </td>
                  </tr>
                  <tr id={`order-${o.id}`} style={{ display: 'none', background: '#fafafa' }}>
                    <td colSpan={7} style={{ padding: 12 }}>
                      <div style={{ fontSize: 13, color: '#6b7280', marginBottom: 8 }}>订单明细:</div>
                      {o.items?.map((item: any) => (
                        <div key={item.id} style={{ display: 'flex', justifyContent: 'space-between', padding: '6px 0', borderBottom: '1px solid #eee' }}>
                          <span>{item.template?.name || '未知模板'} × {item.quantity}</span>
                          <span>¥{item.sub_total.toFixed(2)}</span>
                        </div>
                      ))}
                      {o.items?.length > 1 && (
                        <button onClick={(e) => { e.stopPropagation(); if (confirm('拆出选中项为新订单？')) splitOrder(o.id, o.items.map((i: any) => i.id).slice(0, 1)); }} style={{ marginTop: 8, ...linkBtn }}>
                          拆单
                        </button>
                      )}
                    </td>
                  </tr>
                </React.Fragment>
              ))}
              {orders.length === 0 && <tr><td colSpan={7} style={{ textAlign: 'center', padding: 40, color: '#9ca3af' }}>暂无订单</td></tr>}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

const selectStyle: React.CSSProperties = { padding: '6px 12px', border: '1px solid #d1d5db', borderRadius: 6, fontSize: 13, background: '#fff' };
const primaryBtn: React.CSSProperties = { padding: '8px 16px', background: '#3b82f6', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer', fontSize: 14, textDecoration: 'none' };
const linkBtn: React.CSSProperties = { background: 'none', border: 'none', color: '#3b82f6', cursor: 'pointer', fontSize: 13 };
