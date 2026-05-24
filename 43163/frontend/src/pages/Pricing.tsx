import React, { useState, useEffect } from 'react';
import api from '../api/client';

export default function Pricing() {
  const [rules, setRules] = useState<any[]>([]);
  const [showModal, setShowModal] = useState(false);
  const [editing, setEditing] = useState<any>(null);
  const [templates, setTemplates] = useState<any[]>([]);
  const [calcForm, setCalcForm] = useState<any>({ template_id: '', material_id: '', quantity: 1, urgent: false });
  const [calcResult, setCalcResult] = useState<number | null>(null);

  const load = () => {
    api.get('/pricing/rules').then(({ data }) => setRules(data));
    api.get('/templates').then(({ data }) => setTemplates(data));
  };

  useEffect(() => { load(); }, []);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editing) {
        await api.put(`/pricing/rules/${editing.id}`, editing);
      } else {
        await api.post('/pricing/rules', editing);
      }
      setShowModal(false);
      load();
    } catch (err: any) {
      alert(err.response?.data?.error || '保存失败');
    }
  };

  const handleDelete = (id: number) => {
    if (confirm('确定删除此规则？')) {
      api.delete(`/pricing/rules/${id}`).then(load);
    }
  };

  const handleCalc = async () => {
    if (!calcForm.template_id || !calcForm.material_id) return;
    const { data } = await api.post('/pricing/calculate', calcForm);
    setCalcResult(data.total_price);
  };

  const tpl = templates.find(t => t.id === calcForm.template_id);

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <div style={{ fontSize: 18, fontWeight: 600 }}>智能报价引擎</div>
        <button onClick={() => { setEditing(null); setShowModal(true); }} style={primaryBtn}>+ 新建规则</button>
      </div>

      <div style={{ background: '#fff', padding: 20, borderRadius: 8, marginBottom: 20 }}>
        <h4 style={{ marginBottom: 12 }}>快速报价</h4>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr 1fr auto', gap: 12, alignItems: 'end' }}>
          <div>
            <label style={labelStyle}>模板</label>
            <select style={inputStyle} value={calcForm.template_id} onChange={(e) => setCalcForm({ ...calcForm, template_id: parseInt(e.target.value), material_id: '' })}>
              <option value="">选择模板</option>
              {templates.map(t => <option key={t.id} value={t.id}>{t.name}</option>)}
            </select>
          </div>
          <div>
            <label style={labelStyle}>材质</label>
            <select style={inputStyle} value={calcForm.material_id} onChange={(e) => setCalcForm({ ...calcForm, material_id: parseInt(e.target.value) })} disabled={!tpl}>
              <option value="">选择材质</option>
              {tpl?.materials?.map((m: any) => <option key={m.id} value={m.id}>{m.name}</option>)}
            </select>
          </div>
          <div>
            <label style={labelStyle}>数量</label>
            <input style={inputStyle} type="number" min="1" value={calcForm.quantity} onChange={(e) => setCalcForm({ ...calcForm, quantity: parseInt(e.target.value) })} />
          </div>
          <div>
            <label style={{ ...labelStyle, display: 'flex', alignItems: 'center', marginTop: 4 }}>
              <input type="checkbox" checked={calcForm.urgent} onChange={(e) => setCalcForm({ ...calcForm, urgent: e.target.checked })} style={{ marginRight: 4 }} />
              加急
            </label>
          </div>
          <button onClick={handleCalc} style={primaryBtn}>计算</button>
        </div>
        {calcResult !== null && (
          <div style={{ marginTop: 16, fontSize: 18 }}>
            预估价格: <span style={{ color: '#3b82f6', fontWeight: 600 }}>¥{calcResult.toFixed(2)}</span>
          </div>
        )}
      </div>

      <div style={{ background: '#fff', borderRadius: 8, overflow: 'hidden' }}>
        <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: 14 }}>
          <thead style={{ background: '#f9fafb' }}>
            <tr>
              {['分类', '最小数量', '最大数量', '单价', '折扣(%)', '客户等级', '操作'].map(h => (
                <th key={h} style={{ textAlign: 'left', padding: '12px', color: '#6b7280', fontWeight: 500 }}>{h}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rules.map(r => (
              <tr key={r.id} style={{ borderBottom: '1px solid #f3f4f6' }}>
                <td style={{ padding: '12px' }}>{r.category}</td>
                <td style={{ padding: '12px' }}>{r.min_qty}</td>
                <td style={{ padding: '12px' }}>{r.max_qty}</td>
                <td style={{ padding: '12px' }}>¥{r.unit_price.toFixed(2)}</td>
                <td style={{ padding: '12px' }}>{r.discount}%</td>
                <td style={{ padding: '12px' }}>{r.customer_level}</td>
                <td style={{ padding: '12px' }}>
                  <button onClick={() => { setEditing(r); setShowModal(true); }} style={linkBtn}>编辑</button>
                  <button onClick={() => handleDelete(r.id)} style={{ ...linkBtn, color: '#ef4444' }}>删除</button>
                </td>
              </tr>
            ))}
            {rules.length === 0 && <tr><td colSpan={7} style={{ textAlign: 'center', padding: 40, color: '#9ca3af' }}>暂无规则</td></tr>}
          </tbody>
        </table>
      </div>

      {showModal && (
        <div style={{ position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.5)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000 }}>
          <form onSubmit={handleSave} style={{ background: '#fff', borderRadius: 8, padding: 24, width: 500 }}>
            <h3 style={{ marginBottom: 16 }}>{editing ? '编辑规则' : '新建规则'}</h3>
            {[
              { k: 'category', l: '分类', t: 'text' },
              { k: 'min_qty', l: '最小数量', t: 'number' },
              { k: 'max_qty', l: '最大数量', t: 'number' },
              { k: 'unit_price', l: '单价', t: 'number' },
              { k: 'discount', l: '折扣(%)', t: 'number' },
              { k: 'customer_level', l: '客户等级', t: 'text' },
              { k: 'description', l: '描述', t: 'text' },
            ].map(f => (
              <div key={f.k} style={{ marginBottom: 12 }}>
                <label style={labelStyle}>{f.l}</label>
                <input style={inputStyle} type={f.t} value={editing?.[f.k] || ''} onChange={(e) => setEditing({ ...editing, [f.k]: f.t === 'number' ? parseFloat(e.target.value) : e.target.value })} />
              </div>
            ))}
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

const labelStyle: React.CSSProperties = { display: 'block', marginBottom: 4, fontSize: 13, color: '#374151', fontWeight: 500 };
const inputStyle: React.CSSProperties = { width: '100%', padding: '8px 10px', border: '1px solid #d1d5db', borderRadius: 6, fontSize: 14, outline: 'none', boxSizing: 'border-box' };
const primaryBtn: React.CSSProperties = { padding: '8px 16px', background: '#3b82f6', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer', fontSize: 14 };
const secondaryBtn: React.CSSProperties = { padding: '8px 16px', background: '#fff', color: '#374151', border: '1px solid #d1d5db', borderRadius: 6, cursor: 'pointer', fontSize: 14 };
const linkBtn: React.CSSProperties = { background: 'none', border: 'none', color: '#3b82f6', cursor: 'pointer', fontSize: 13, marginRight: 8 };
