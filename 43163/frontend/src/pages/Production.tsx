import React, { useState, useEffect } from 'react';
import api from '../api/client';

export default function Production() {
  const [lines, setLines] = useState<any[]>([]);
  const [schedules, setSchedules] = useState<any[]>([]);
  const [orders, setOrders] = useState<any[]>([]);
  const [showLineModal, setShowLineModal] = useState(false);
  const [showSchedModal, setShowSchedModal] = useState(false);
  const [newLine, setNewLine] = useState({ name: '', code: '', capacity: 1000 });
  const [newSched, setNewSched] = useState<any>({ order_id: '', line_id: '', planned_qty: 100, start_date: '', end_date: '' });

  const load = () => {
    api.get('/production/lines').then(({ data }) => setLines(data));
    api.get('/production/schedules').then(({ data }) => setSchedules(data));
    api.get('/orders').then(({ data }) => setOrders(data.filter((o: any) => o.status !== 'shipped' && o.status !== 'cancelled')));
  };

  useEffect(() => { load(); }, []);

  const createLine = async (e: React.FormEvent) => {
    e.preventDefault();
    await api.post('/production/lines', newLine);
    setShowLineModal(false);
    setNewLine({ name: '', code: '', capacity: 1000 });
    load();
  };

  const createSched = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.post('/production/schedules', newSched);
      setShowSchedModal(false);
      setNewSched({ order_id: '', line_id: '', planned_qty: 100, start_date: '', end_date: '' });
      load();
    } catch (err: any) {
      alert(err.response?.data?.error || '排程失败');
    }
  };

  const updateProgress = async (id: number, qty: number) => {
    await api.put(`/production/schedules/${id}/progress`, { produced_qty: qty });
    load();
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <div style={{ fontSize: 18, fontWeight: 600 }}>生产排程管理</div>
        <div style={{ display: 'flex', gap: 8 }}>
          <button onClick={() => setShowLineModal(true)} style={primaryBtn}>+ 新建产线</button>
          <button onClick={() => setShowSchedModal(true)} style={primaryBtn}>+ 新建排程</button>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 2fr', gap: 16 }}>
        <div>
          <h4 style={{ marginBottom: 12 }}>产线负荷</h4>
          <div style={{ background: '#fff', borderRadius: 8, padding: 16 }}>
            {lines.map(l => (
              <div key={l.id} style={{ marginBottom: 16 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 4 }}>
                  <span style={{ fontSize: 14, fontWeight: 500 }}>{l.name}</span>
                  <span style={{ fontSize: 12, color: '#6b7280' }}>{l.workload}/{l.capacity}</span>
                </div>
                <div style={{ height: 8, background: '#e5e7eb', borderRadius: 4, overflow: 'hidden' }}>
                  <div style={{
                    height: '100%',
                    width: `${Math.min(100, l.capacity > 0 ? (l.workload / l.capacity) * 100 : 0)}%`,
                    background: l.workload > l.capacity ? '#ef4444' : l.workload > l.capacity * 0.8 ? '#f59e0b' : '#10b981',
                  }} />
                </div>
              </div>
            ))}
            {lines.length === 0 && <div style={{ textAlign: 'center', color: '#9ca3af', padding: 20 }}>暂无产线</div>}
          </div>
        </div>

        <div>
          <h4 style={{ marginBottom: 12 }}>生产排程</h4>
          <div style={{ background: '#fff', borderRadius: 8, overflow: 'hidden' }}>
            <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: 14 }}>
              <thead style={{ background: '#f9fafb' }}>
                <tr>
                  {['订单号', '产线', '计划/已生产', '进度', '开始日期', '结束日期', '状态', '操作'].map(h => (
                    <th key={h} style={{ textAlign: 'left', padding: '10px', color: '#6b7280', fontWeight: 500 }}>{h}</th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {schedules.map(s => {
                  const pct = s.planned_qty > 0 ? (s.produced_qty / s.planned_qty) * 100 : 0;
                  return (
                    <tr key={s.id} style={{ borderBottom: '1px solid #f3f4f6' }}>
                      <td style={{ padding: '10px' }}>{s.order?.order_no}</td>
                      <td style={{ padding: '10px' }}>{s.line?.name}</td>
                      <td style={{ padding: '10px' }}>{s.planned_qty}/{s.produced_qty}</td>
                      <td style={{ padding: '10px' }}>
                        <div style={{ width: 100, height: 6, background: '#e5e7eb', borderRadius: 3 }}>
                          <div style={{ width: `${pct}%`, height: '100%', background: pct >= 100 ? '#10b981' : '#3b82f6', borderRadius: 3 }} />
                        </div>
                      </td>
                      <td style={{ padding: '10px' }}>{new Date(s.start_date).toLocaleDateString()}</td>
                      <td style={{ padding: '10px' }}>{new Date(s.end_date).toLocaleDateString()}</td>
                      <td style={{ padding: '10px' }}>
                        <span style={{ padding: '2px 8px', borderRadius: 4, fontSize: 12, background: s.status === 'completed' ? '#dcfce7' : '#dbeafe', color: s.status === 'completed' ? '#16a34a' : '#2563eb' }}>
                          {s.status === 'completed' ? '已完成' : '生产中'}
                        </span>
                      </td>
                      <td style={{ padding: '10px' }}>
                        {s.status !== 'completed' && (
                          <button onClick={() => updateProgress(s.id, s.planned_qty)} style={linkBtn}>完成</button>
                        )}
                      </td>
                    </tr>
                  );
                })}
                {schedules.length === 0 && <tr><td colSpan={8} style={{ textAlign: 'center', padding: 40, color: '#9ca3af' }}>暂无排程</td></tr>}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {showLineModal && (
        <div style={modalBg}>
          <form onSubmit={createLine} style={modalBox}>
            <h3 style={{ marginBottom: 16 }}>新建产线</h3>
            <div style={{ marginBottom: 12 }}><label style={labelStyle}>产线名称</label><input style={inputStyle} value={newLine.name} onChange={e => setNewLine({ ...newLine, name: e.target.value })} required /></div>
            <div style={{ marginBottom: 12 }}><label style={labelStyle}>产线代码</label><input style={inputStyle} value={newLine.code} onChange={e => setNewLine({ ...newLine, code: e.target.value })} required /></div>
            <div style={{ marginBottom: 12 }}><label style={labelStyle}>日产能</label><input style={inputStyle} type="number" value={newLine.capacity} onChange={e => setNewLine({ ...newLine, capacity: parseInt(e.target.value) })} required /></div>
            <div style={{ display: 'flex', gap: 12, justifyContent: 'flex-end', marginTop: 16 }}>
              <button type="button" onClick={() => setShowLineModal(false)} style={secondaryBtn}>取消</button>
              <button type="submit" style={primaryBtn}>保存</button>
            </div>
          </form>
        </div>
      )}

      {showSchedModal && (
        <div style={modalBg}>
          <form onSubmit={createSched} style={modalBox}>
            <h3 style={{ marginBottom: 16 }}>新建排程</h3>
            <div style={{ marginBottom: 12 }}><label style={labelStyle}>订单</label>
              <select style={inputStyle} value={newSched.order_id} onChange={e => setNewSched({ ...newSched, order_id: parseInt(e.target.value) })} required>
                <option value="">选择订单</option>
                {orders.map((o: any) => <option key={o.id} value={o.id}>{o.order_no}</option>)}
              </select>
            </div>
            <div style={{ marginBottom: 12 }}><label style={labelStyle}>产线</label>
              <select style={inputStyle} value={newSched.line_id} onChange={e => setNewSched({ ...newSched, line_id: parseInt(e.target.value) })} required>
                <option value="">选择产线</option>
                {lines.map(l => <option key={l.id} value={l.id}>{l.name} ({l.workload}/{l.capacity})</option>)}
              </select>
            </div>
            <div style={{ marginBottom: 12 }}><label style={labelStyle}>计划数量</label><input style={inputStyle} type="number" value={newSched.planned_qty} onChange={e => setNewSched({ ...newSched, planned_qty: parseInt(e.target.value) })} required /></div>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12, marginBottom: 12 }}>
              <div><label style={labelStyle}>开始日期</label><input style={inputStyle} type="date" value={newSched.start_date} onChange={e => setNewSched({ ...newSched, start_date: e.target.value })} required /></div>
              <div><label style={labelStyle}>结束日期</label><input style={inputStyle} type="date" value={newSched.end_date} onChange={e => setNewSched({ ...newSched, end_date: e.target.value })} required /></div>
            </div>
            <div style={{ display: 'flex', gap: 12, justifyContent: 'flex-end', marginTop: 16 }}>
              <button type="button" onClick={() => setShowSchedModal(false)} style={secondaryBtn}>取消</button>
              <button type="submit" style={primaryBtn}>保存</button>
            </div>
          </form>
        </div>
      )}
    </div>
  );
}

const modalBg: React.CSSProperties = { position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.5)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000 };
const modalBox: React.CSSProperties = { background: '#fff', borderRadius: 8, padding: 24, width: 450 };
const labelStyle: React.CSSProperties = { display: 'block', marginBottom: 4, fontSize: 13, color: '#374151', fontWeight: 500 };
const inputStyle: React.CSSProperties = { width: '100%', padding: '8px 10px', border: '1px solid #d1d5db', borderRadius: 6, fontSize: 14, outline: 'none', boxSizing: 'border-box' };
const primaryBtn: React.CSSProperties = { padding: '8px 16px', background: '#3b82f6', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer', fontSize: 14 };
const secondaryBtn: React.CSSProperties = { padding: '8px 16px', background: '#fff', color: '#374151', border: '1px solid #d1d5db', borderRadius: 6, cursor: 'pointer', fontSize: 14 };
const linkBtn: React.CSSProperties = { background: 'none', border: 'none', color: '#3b82f6', cursor: 'pointer', fontSize: 13 };
