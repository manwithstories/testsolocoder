import React, { useState, useEffect } from 'react';
import api from '../api/client';

const levels = ['normal', 'silver', 'gold', 'platinum'];

export default function Customers() {
  const [customers, setCustomers] = useState<any[]>([]);
  const [showModal, setShowModal] = useState(false);
  const [editing, setEditing] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  const load = () => {
    api.get('/customers').then(({ data }) => { setCustomers(data); setLoading(false); });
  };

  useEffect(() => { load(); }, []);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editing?.id) {
        await api.put(`/customers/${editing.id}`, editing);
      } else {
        await api.post('/customers', editing);
      }
      setShowModal(false);
      load();
    } catch (err: any) {
      alert(err.response?.data?.error || '保存失败');
    }
  };

  const handleDelete = (id: number) => {
    if (confirm('确定删除此客户？')) {
      api.delete(`/customers/${id}`).then(load);
    }
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <div style={{ fontSize: 18, fontWeight: 600 }}>客户管理</div>
        <button onClick={() => { setEditing({ name: '', contact: '', phone: '', email: '', address: '', level: 'normal', credit_limit: 0 }); setShowModal(true); }} style={primaryBtn}>
          + 新建客户
        </button>
      </div>

      {loading ? <div>加载中...</div> : (
        <div style={{ background: '#fff', borderRadius: 8, overflow: 'hidden' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: 14 }}>
            <thead style={{ background: '#f9fafb' }}>
              <tr>
                {['客户名称', '联系人', '电话', '邮箱', '等级', '信用额度', '余额', '操作'].map(h => (
                  <th key={h} style={{ textAlign: 'left', padding: '12px', color: '#6b7280', fontWeight: 500 }}>{h}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {customers.map(c => (
                <tr key={c.id} style={{ borderBottom: '1px solid #f3f4f6' }}>
                  <td style={{ padding: '12px', fontWeight: 500 }}>{c.name}</td>
                  <td style={{ padding: '12px' }}>{c.contact || '-'}</td>
                  <td style={{ padding: '12px' }}>{c.phone || '-'}</td>
                  <td style={{ padding: '12px' }}>{c.email || '-'}</td>
                  <td style={{ padding: '12px' }}>
                    <span style={{ padding: '2px 8px', borderRadius: 4, fontSize: 12, background: levelColor(c.level) + '20', color: levelColor(c.level) }}>{c.level}</span>
                  </td>
                  <td style={{ padding: '12px' }}>¥{c.credit_limit.toFixed(2)}</td>
                  <td style={{ padding: '12px' }}>¥{c.balance.toFixed(2)}</td>
                  <td style={{ padding: '12px' }}>
                    <button onClick={() => { setEditing(c); setShowModal(true); }} style={linkBtn}>编辑</button>
                    <button onClick={() => handleDelete(c.id)} style={{ ...linkBtn, color: '#ef4444' }}>删除</button>
                  </td>
                </tr>
              ))}
              {customers.length === 0 && <tr><td colSpan={8} style={{ textAlign: 'center', padding: 40, color: '#9ca3af' }}>暂无客户</td></tr>}
            </tbody>
          </table>
        </div>
      )}

      {showModal && (
        <div style={modalBg}>
          <form onSubmit={handleSave} style={modalBox}>
            <h3 style={{ marginBottom: 16 }}>{editing?.id ? '编辑客户' : '新建客户'}</h3>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
              <div style={{ marginBottom: 12 }}><label style={labelStyle}>客户名称 *</label><input style={inputStyle} value={editing?.name || ''} onChange={e => setEditing({ ...editing, name: e.target.value })} required /></div>
              <div style={{ marginBottom: 12 }}><label style={labelStyle}>联系人</label><input style={inputStyle} value={editing?.contact || ''} onChange={e => setEditing({ ...editing, contact: e.target.value })} /></div>
              <div style={{ marginBottom: 12 }}><label style={labelStyle}>电话</label><input style={inputStyle} value={editing?.phone || ''} onChange={e => setEditing({ ...editing, phone: e.target.value })} /></div>
              <div style={{ marginBottom: 12 }}><label style={labelStyle}>邮箱</label><input style={inputStyle} type="email" value={editing?.email || ''} onChange={e => setEditing({ ...editing, email: e.target.value })} /></div>
              <div style={{ marginBottom: 12 }}>
                <label style={labelStyle}>等级</label>
                <select style={inputStyle} value={editing?.level || 'normal'} onChange={e => setEditing({ ...editing, level: e.target.value })}>
                  {levels.map(l => <option key={l} value={l}>{l}</option>)}
                </select>
              </div>
              <div style={{ marginBottom: 12 }}><label style={labelStyle}>信用额度</label><input style={inputStyle} type="number" step="0.01" value={editing?.credit_limit || 0} onChange={e => setEditing({ ...editing, credit_limit: parseFloat(e.target.value) })} /></div>
            </div>
            <div style={{ marginBottom: 12 }}><label style={labelStyle}>地址</label><input style={inputStyle} value={editing?.address || ''} onChange={e => setEditing({ ...editing, address: e.target.value })} /></div>
            <div style={{ display: 'flex', gap: 12, justifyContent: 'flex-end', marginTop: 16 }}>
              <button type="button" onClick={() => setShowModal(false)} style={secondaryBtn}>取消</button>
              <button type="submit" style={primaryBtn}>保存</button>
            </div>
          </form>
        </div>
      )}
    </div>
  );
}

function levelColor(l: string): string {
  const map: Record<string, string> = { normal: '#6b7280', silver: '#9ca3af', gold: '#f59e0b', platinum: '#8b5cf6' };
  return map[l] || '#6b7280';
}

const modalBg: React.CSSProperties = { position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.5)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000 };
const modalBox: React.CSSProperties = { background: '#fff', borderRadius: 8, padding: 24, width: 550, maxHeight: '90vh', overflow: 'auto' };
const labelStyle: React.CSSProperties = { display: 'block', marginBottom: 4, fontSize: 13, color: '#374151', fontWeight: 500 };
const inputStyle: React.CSSProperties = { width: '100%', padding: '8px 10px', border: '1px solid #d1d5db', borderRadius: 6, fontSize: 14, outline: 'none', boxSizing: 'border-box' };
const primaryBtn: React.CSSProperties = { padding: '8px 16px', background: '#3b82f6', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer', fontSize: 14 };
const secondaryBtn: React.CSSProperties = { padding: '8px 16px', background: '#fff', color: '#374151', border: '1px solid #d1d5db', borderRadius: 6, cursor: 'pointer', fontSize: 14 };
const linkBtn: React.CSSProperties = { background: 'none', border: 'none', color: '#3b82f6', cursor: 'pointer', fontSize: 13, marginRight: 8 };
