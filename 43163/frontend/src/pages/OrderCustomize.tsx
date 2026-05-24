import React, { useState, useEffect } from 'react';
import api from '../api/client';

export default function OrderCustomize() {
  const [templates, setTemplates] = useState<any[]>([]);
  const [customers, setCustomers] = useState<any[]>([]);
  const [form, setForm] = useState<any>({
    customer_id: '', urgent: false, remark: '',
    items: [{ template_id: '', material_id: '', process_ids: [] as number[], quantity: 1, file_asset_id: null, specification: '' }],
  });
  const [previewPrice, setPreviewPrice] = useState<number | null>(null);
  const [uploadedFile, setUploadedFile] = useState<any>(null);

  useEffect(() => {
    api.get('/templates').then(({ data }) => setTemplates(data));
    api.get('/customers').then(({ data }) => setCustomers(data));
  }, []);

  const getTemplate = (id: number) => templates.find(t => t.id === id);
  const getMaterial = (t: any, id: number) => t?.materials?.find((m: any) => m.id === id);

  const calcPrice = async () => {
    if (!form.items[0].template_id || !form.items[0].material_id || !form.items[0].quantity) return;
    const total = await Promise.all(form.items.map(async (item: any) => {
      const { data } = await api.post('/pricing/calculate', {
        template_id: item.template_id,
        material_id: item.material_id,
        process_ids: item.process_ids,
        quantity: item.quantity,
        customer_level: customers.find(c => c.id === form.customer_id)?.level || 'normal',
        urgent: form.urgent,
      });
      return data.total_price;
    }));
    setPreviewPrice(total.reduce((a, b) => a + b, 0));
  };

  useEffect(() => { calcPrice(); }, [form.items, form.urgent, form.customer_id]);

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    const fd = new FormData();
    fd.append('file', file);
    try {
      const { data } = await api.post('/files/upload', fd, { headers: { 'Content-Type': 'multipart/form-data' } });
      setUploadedFile(data);
      setForm((f: any) => ({ ...f, items: f.items.map((it: any, i: number) => i === 0 ? { ...it, file_asset_id: data.id } : it) }));
    } catch (err: any) {
      alert(err.response?.data?.error || '上传失败');
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.customer_id) { alert('请选择客户'); return; }
    try {
      await api.post('/orders', form);
      alert('订单创建成功！');
      window.location.href = '/orders';
    } catch (err: any) {
      alert(err.response?.data?.error || '创建失败');
    }
  };

  return (
    <div>
      <div style={{ fontSize: 18, fontWeight: 600, marginBottom: 16 }}>在线定制下单</div>

      <form onSubmit={handleSubmit} style={{ background: '#fff', padding: 24, borderRadius: 8 }}>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: 16, marginBottom: 20 }}>
          <div>
            <label style={labelStyle}>客户 *</label>
            <select style={inputStyle} value={form.customer_id} onChange={(e) => setForm({ ...form, customer_id: parseInt(e.target.value) })} required>
              <option value="">请选择客户</option>
              {customers.map(c => <option key={c.id} value={c.id}>{c.name} ({c.level})</option>)}
            </select>
            {form.customer_id && (() => {
              const cust = customers.find(c => c.id === form.customer_id);
              if (cust) {
                const available = cust.credit_limit - cust.balance;
                return (
                  <div style={{ marginTop: 6, fontSize: 12, color: '#6b7280' }}>
                    信用额度 ¥{cust.credit_limit.toFixed(2)} · 已用 ¥{cust.balance.toFixed(2)} · 可用 ¥{available.toFixed(2)}
                  </div>
                );
              }
              return null;
            })()}
          </div>
          <div>
            <label style={labelStyle}>加急订单</label>
            <label style={{ display: 'flex', alignItems: 'center', marginTop: 6, fontSize: 14 }}>
              <input type="checkbox" checked={form.urgent} onChange={(e) => setForm({ ...form, urgent: e.target.checked })} style={{ marginRight: 6 }} />
              加急处理（价格×1.3）
            </label>
          </div>
          <div>
            <label style={labelStyle}>备注</label>
            <input style={inputStyle} value={form.remark} onChange={(e) => setForm({ ...form, remark: e.target.value })} />
          </div>
        </div>

        {form.items.map((item: any, idx: number) => {
          const tpl = getTemplate(item.template_id);
          return (
            <div key={idx} style={{ border: '1px solid #e5e7eb', borderRadius: 8, padding: 16, marginBottom: 12 }}>
              <div style={{ fontWeight: 600, marginBottom: 12 }}>商品 {idx + 1}</div>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: 12 }}>
                <div>
                  <label style={labelStyle}>模板 *</label>
                  <select style={inputStyle} value={item.template_id} onChange={(e) => setForm((f: any) => ({
                    ...f, items: f.items.map((it: any, i: number) => i === idx ? { ...it, template_id: parseInt(e.target.value), material_id: '', process_ids: [] } : it)
                  }))} required>
                    <option value="">选择模板</option>
                    {templates.map(t => <option key={t.id} value={t.id}>{t.name} ({t.category})</option>)}
                  </select>
                </div>
                <div>
                  <label style={labelStyle}>材质 *</label>
                  <select style={inputStyle} value={item.material_id} onChange={(e) => setForm((f: any) => ({
                    ...f, items: f.items.map((it: any, i: number) => i === idx ? { ...it, material_id: parseInt(e.target.value) } : it)
                  }))} disabled={!tpl} required>
                    <option value="">选择材质</option>
                    {tpl?.materials?.map((m: any) => <option key={m.id} value={m.id}>{m.name} (¥{m.base_price})</option>)}
                  </select>
                </div>
                <div>
                  <label style={labelStyle}>数量 *</label>
                  <input style={inputStyle} type="number" min="1" value={item.quantity} onChange={(e) => setForm((f: any) => ({
                    ...f, items: f.items.map((it: any, i: number) => i === idx ? { ...it, quantity: parseInt(e.target.value) } : it)
                  }))} required />
                </div>
              </div>

              {tpl && tpl.processes?.length > 0 && (
                <div style={{ marginTop: 12 }}>
                  <label style={labelStyle}>工艺选项</label>
                  <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8 }}>
                    {tpl.processes.map((p: any) => (
                      <label key={p.id} style={{ display: 'flex', alignItems: 'center', padding: '6px 10px', border: '1px solid #e5e7eb', borderRadius: 6, fontSize: 13, cursor: 'pointer', background: item.process_ids.includes(p.id) ? '#dbeafe' : '#fff' }}>
                        <input type="checkbox" checked={item.process_ids.includes(p.id)} onChange={(e) => {
                          const ids = e.target.checked ? [...item.process_ids, p.id] : item.process_ids.filter((id: number) => id !== p.id);
                          setForm((f: any) => ({ ...f, items: f.items.map((it: any, i: number) => i === idx ? { ...it, process_ids: ids } : it) }));
                        }} style={{ marginRight: 4 }} />
                        {p.name} {p.extra_price > 0 && <span style={{ color: '#6b7280' }}>+¥{p.extra_price}</span>}
                      </label>
                    ))}
                  </div>
                </div>
              )}

              <div style={{ marginTop: 12 }}>
                <label style={labelStyle}>设计文件</label>
                <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  <input type="file" onChange={handleFileUpload} accept=".pdf,.ai,.eps,.psd,.jpg,.jpeg,.png,.cdr" />
                  {uploadedFile && <span style={{ color: '#10b981', fontSize: 13 }}>✓ {uploadedFile.file_name}</span>}
                </div>
                <div style={{ fontSize: 12, color: '#9ca3af', marginTop: 4 }}>支持 PDF, AI, EPS, PSD, JPG, PNG, CDR 格式</div>
              </div>

              <div style={{ marginTop: 12 }}>
                <label style={labelStyle}>规格说明</label>
                <input style={inputStyle} placeholder="如：双面印刷、覆膜等" value={item.specification} onChange={(e) => setForm((f: any) => ({
                  ...f, items: f.items.map((it: any, i: number) => i === idx ? { ...it, specification: e.target.value } : it)
                }))} />
              </div>

              {form.items.length > 1 && (
                <button type="button" onClick={() => setForm((f: any) => ({ ...f, items: f.items.filter((_: any, i: number) => i !== idx) }))} style={{ marginTop: 8, ...linkBtn, color: '#ef4444' }}>
                  移除
                </button>
              )}
            </div>
          );
        })}

        <button type="button" onClick={() => setForm((f: any) => ({
          ...f, items: [...f.items, { template_id: '', material_id: '', process_ids: [], quantity: 1, file_asset_id: null, specification: '' }]
        }))} style={linkBtn}>+ 添加商品</button>

        <div style={{ marginTop: 20, padding: 16, borderRadius: 8, display: 'flex', justifyContent: 'space-between', alignItems: 'center',
          background: previewPrice !== null && form.customer_id && (() => {
            const cust = customers.find(c => c.id === form.customer_id);
            return cust && previewPrice > (cust.credit_limit - cust.balance);
          })() ? '#fef2f2' : '#f0f9ff'
        }}>
          <div>
            <span style={{ fontSize: 15 }}>预估总价:</span>
            {previewPrice !== null && form.customer_id && (() => {
              const cust = customers.find(c => c.id === form.customer_id);
              if (cust && previewPrice > (cust.credit_limit - cust.balance)) {
                return (
                  <span style={{ fontSize: 13, color: '#ef4444', marginLeft: 12 }}>
                    ⚠ 超出可用信用额度 ¥{(cust.credit_limit - cust.balance).toFixed(2)}
                  </span>
                );
              }
              return null;
            })()}
          </div>
          <span style={{ fontSize: 24, fontWeight: 700, color: previewPrice !== null && form.customer_id && (() => {
            const cust = customers.find(c => c.id === form.customer_id);
            return cust && previewPrice > (cust.credit_limit - cust.balance);
          })() ? '#ef4444' : '#3b82f6' }}>
            {previewPrice !== null ? `¥${previewPrice.toFixed(2)}` : '请选择模板和材质'}
          </span>
        </div>

        <div style={{ marginTop: 16, display: 'flex', gap: 12, justifyContent: 'flex-end' }}>
          <a href="/orders" style={secondaryBtn}>取消</a>
          <button type="submit" style={primaryBtn}>提交订单</button>
        </div>
      </form>
    </div>
  );
}

const labelStyle: React.CSSProperties = { display: 'block', marginBottom: 4, fontSize: 13, color: '#374151', fontWeight: 500 };
const inputStyle: React.CSSProperties = { width: '100%', padding: '8px 10px', border: '1px solid #d1d5db', borderRadius: 6, fontSize: 14, outline: 'none', boxSizing: 'border-box' };
const primaryBtn: React.CSSProperties = { padding: '10px 20px', background: '#3b82f6', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer', fontSize: 14 };
const secondaryBtn: React.CSSProperties = { padding: '10px 20px', background: '#fff', color: '#374151', border: '1px solid #d1d5db', borderRadius: 6, cursor: 'pointer', fontSize: 14, textDecoration: 'none' };
const linkBtn: React.CSSProperties = { background: 'none', border: 'none', color: '#3b82f6', cursor: 'pointer', fontSize: 13 };
