import React, { useState, useEffect } from 'react';
import api from '../api/client';

export default function Invoices() {
  const [invoices, setInvoices] = useState<any[]>([]);
  const [customers, setCustomers] = useState<any[]>([]);
  const [showModal, setShowModal] = useState(false);
  const [form, setForm] = useState({ customer_id: '', period_start: '', period_end: '' });
  const [filterCustomer, setFilterCustomer] = useState('');

  const load = () => {
    const params = new URLSearchParams();
    if (filterCustomer) params.set('customer_id', filterCustomer);
    api.get('/invoices', { params }).then(({ data }) => setInvoices(data));
    api.get('/customers').then(({ data }) => setCustomers(data));
  };

  useEffect(() => { load(); }, [filterCustomer]);

  const generate = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.post('/invoices', form);
      setShowModal(false);
      setForm({ customer_id: '', period_start: '', period_end: '' });
      load();
    } catch (err: any) {
      alert(err.response?.data?.error || '生成失败');
    }
  };

  const exportCSV = (inv: any) => {
    let csv = '订单号,金额\n';
    inv.items?.forEach((it: any) => { csv += `${it.order_id},${it.amount}\n`; });
    csv += `\n总计,${inv.total_amount}`;
    const blob = new Blob([csv], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url; a.download = `${inv.invoice_no}.csv`; a.click();
    URL.revokeObjectURL(url);
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <div style={{ fontSize: 18, fontWeight: 600 }}>财务对账</div>
        <div style={{ display: 'flex', gap: 8 }}>
          <select style={selectStyle} value={filterCustomer} onChange={(e) => setFilterCustomer(e.target.value)}>
            <option value="">全部客户</option>
            {customers.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
          </select>
          <button onClick={() => setShowModal(true)} style={primaryBtn}>+ 生成对账单</button>
        </div>
      </div>

      <div style={{ background: '#fff', borderRadius: 8, overflow: 'hidden' }}>
        <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: 14 }}>
          <thead style={{ background: '#f9fafb' }}>
            <tr>
              {['对账单号', '客户', '期间', '金额', '已付', '状态', '操作'].map(h => (
                <th key={h} style={{ textAlign: 'left', padding: '12px', color: '#6b7280', fontWeight: 500 }}>{h}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {invoices.map(inv => (
              <tr key={inv.id} style={{ borderBottom: '1px solid #f3f4f6' }}>
                <td style={{ padding: '12px', fontWeight: 500 }}>{inv.invoice_no}</td>
                <td style={{ padding: '12px' }}>{inv.customer?.name}</td>
                <td style={{ padding: '12px' }}>{new Date(inv.period_start).toLocaleDateString()} ~ {new Date(inv.period_end).toLocaleDateString()}</td>
                <td style={{ padding: '12px' }}>¥{inv.total_amount.toFixed(2)}</td>
                <td style={{ padding: '12px' }}>¥{inv.paid_amount.toFixed(2)}</td>
                <td style={{ padding: '12px' }}>
                  <span style={{ padding: '2px 8px', borderRadius: 4, fontSize: 12, background: inv.status === 'paid' ? '#dcfce7' : '#fef3c7', color: inv.status === 'paid' ? '#16a34a' : '#d97706' }}>
                    {inv.status === 'paid' ? '已结清' : '未结清'}
                  </span>
                </td>
                <td style={{ padding: '12px' }}>
                  <button onClick={() => exportCSV(inv)} style={linkBtn}>导出</button>
                </td>
              </tr>
            ))}
            {invoices.length === 0 && <tr><td colSpan={7} style={{ textAlign: 'center', padding: 40, color: '#9ca3af' }}>暂无对账单</td></tr>}
          </tbody>
        </table>
      </div>

      {showModal && (
        <div style={modalBg}>
          <form onSubmit={generate} style={modalBox}>
            <h3 style={{ marginBottom: 16 }}>生成对账单</h3>
            <div style={{ marginBottom: 12 }}>
              <label style={labelStyle}>客户 *</label>
              <select style={inputStyle} value={form.customer_id} onChange={e => setForm({ ...form, customer_id: e.target.value })} required>
                <option value="">选择客户</option>
                {customers.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
              </select>
            </div>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12, marginBottom: 12 }}>
              <div><label style={labelStyle}>开始日期 *</label><input style={inputStyle} type="date" value={form.period_start} onChange={e => setForm({ ...form, period_start: e.target.value })} required /></div>
              <div><label style={labelStyle}>结束日期 *</label><input style={inputStyle} type="date" value={form.period_end} onChange={e => setForm({ ...form, period_end: e.target.value })} required /></div>
            </div>
            <div style={{ display: 'flex', gap: 12, justifyContent: 'flex-end', marginTop: 16 }}>
              <button type="button" onClick={() => setShowModal(false)} style={secondaryBtn}>取消</button>
              <button type="submit" style={primaryBtn}>生成</button>
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
const selectStyle: React.CSSProperties = { padding: '6px 12px', border: '1px solid #d1d5db', borderRadius: 6, fontSize: 13, background: '#fff' };
const primaryBtn: React.CSSProperties = { padding: '8px 16px', background: '#3b82f6', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer', fontSize: 14 };
const secondaryBtn: React.CSSProperties = { padding: '8px 16px', background: '#fff', color: '#374151', border: '1px solid #d1d5db', borderRadius: 6, cursor: 'pointer', fontSize: 14 };
const linkBtn: React.CSSProperties = { background: 'none', border: 'none', color: '#3b82f6', cursor: 'pointer', fontSize: 13 };
