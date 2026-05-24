import React, { useState, useEffect } from 'react';
import api from '../api/client';

interface Template {
  id: number;
  name: string;
  category: string;
  width_mm: number;
  height_mm: number;
  description: string;
  thumbnail: string;
  active: boolean;
  materials: any[];
  processes: any[];
  options: any[];
}

const categories = ['名片', '传单', '画册', '海报', '包装盒', '不干胶', '其他'];

export default function Templates() {
  const [templates, setTemplates] = useState<Template[]>([]);
  const [showModal, setShowModal] = useState(false);
  const [editing, setEditing] = useState<Template | null>(null);
  const [loading, setLoading] = useState(true);

  const load = () => {
    api.get('/templates').then(({ data }) => {
      setTemplates(data);
      setLoading(false);
    });
  };

  useEffect(() => { load(); }, []);

  const handleDelete = (id: number) => {
    if (confirm('确定删除此模板？')) {
      api.delete(`/templates/${id}`).then(load);
    }
  };

  const toggleActive = (t: Template) => {
    api.put(`/templates/${t.id}`, { ...t, active: !t.active }).then(load);
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <div style={{ fontSize: 18, fontWeight: 600 }}>商品模板管理</div>
        <button onClick={() => { setEditing(null); setShowModal(true); }} style={primaryBtn}>
          + 新建模板
        </button>
      </div>

      {loading ? <div>加载中...</div> : (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 16 }}>
          {templates.map((t) => (
            <div key={t.id} style={{
              background: '#fff', borderRadius: 8, overflow: 'hidden',
              border: '1px solid #e5e7eb',
            }}>
              <div style={{
                height: 140, background: `linear-gradient(135deg, ${categoryColor(t.category)}40, ${categoryColor(t.category)}10)`,
                display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 48,
              }}>
                {categoryIcon(t.category)}
              </div>
              <div style={{ padding: 14 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 6 }}>
                  <span style={{ fontWeight: 600, fontSize: 15 }}>{t.name}</span>
                  <span style={{ fontSize: 12, padding: '2px 6px', borderRadius: 4, background: t.active ? '#dcfce7' : '#fef3c7', color: t.active ? '#16a34a' : '#d97706' }}>
                    {t.active ? '启用' : '停用'}
                  </span>
                </div>
                <div style={{ fontSize: 13, color: '#6b7280', marginBottom: 4 }}>{t.category} · {t.width_mm}×{t.height_mm}mm</div>
                <div style={{ fontSize: 12, color: '#9ca3af', marginBottom: 12 }}>
                  {t.materials?.length} 种材质 · {t.processes?.length} 种工艺
                </div>
                <div style={{ display: 'flex', gap: 8 }}>
                  <button onClick={() => { setEditing(t); setShowModal(true); }} style={linkBtn}>编辑</button>
                  <button onClick={() => toggleActive(t)} style={linkBtn}>{t.active ? '停用' : '启用'}</button>
                  <button onClick={() => handleDelete(t.id)} style={{ ...linkBtn, color: '#ef4444' }}>删除</button>
                </div>
              </div>
            </div>
          ))}
          {templates.length === 0 && <div style={{ gridColumn: '1/-1', textAlign: 'center', padding: 40, color: '#9ca3af' }}>暂无模板</div>}
        </div>
      )}

      {showModal && (
        <TemplateForm template={editing} onClose={() => setShowModal(false)} onSave={load} />
      )}
    </div>
  );
}

function TemplateForm({ template, onClose, onSave }: { template: Template | null; onClose: () => void; onSave: () => void }) {
  const [form, setForm] = useState<any>(template || {
    name: '', category: '名片', width_mm: 90, height_mm: 54, description: '',
    materials: [{ name: '', description: '', base_price: 0, unit: 'sheet' }],
    processes: [{ name: '', description: '', extra_price: 0 }],
    options: [],
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (template) {
        await api.put(`/templates/${template.id}`, form);
      } else {
        await api.post('/templates', form);
      }
      onSave();
      onClose();
    } catch (err: any) {
      alert(err.response?.data?.error || '保存失败');
    }
  };

  return (
    <div style={{
      position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.5)',
      display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000,
    }}>
      <form onSubmit={handleSubmit} style={{
        background: '#fff', borderRadius: 8, padding: 24, width: 600, maxHeight: '90vh', overflow: 'auto',
      }}>
        <h3 style={{ marginBottom: 16 }}>{template ? '编辑模板' : '新建模板'}</h3>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
          <Field label="模板名称" value={form.name} onChange={(v) => setForm({ ...form, name: v })} required />
          <div>
            <label style={labelStyle}>分类</label>
            <select style={inputStyle} value={form.category} onChange={(e) => setForm({ ...form, category: e.target.value })}>
              {categories.map(c => <option key={c} value={c}>{c}</option>)}
            </select>
          </div>
          <Field label="宽度(mm)" type="number" value={form.width_mm} onChange={(v) => setForm({ ...form, width_mm: parseFloat(v) })} required />
          <Field label="高度(mm)" type="number" value={form.height_mm} onChange={(v) => setForm({ ...form, height_mm: parseFloat(v) })} required />
        </div>
        <div style={{ marginTop: 12 }}>
          <label style={labelStyle}>描述</label>
          <textarea style={{ ...inputStyle, minHeight: 60 }} value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} />
        </div>

        <h4 style={{ marginTop: 16, marginBottom: 8 }}>材质选项</h4>
        {form.materials.map((m: any, i: number) => (
          <div key={i} style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr auto', gap: 8, marginBottom: 8 }}>
            <input style={inputStyle} placeholder="材质名称" value={m.name} onChange={(e) => {
              const mats = [...form.materials]; mats[i] = { ...mats[i], name: e.target.value }; setForm({ ...form, materials: mats });
            }} />
            <input style={inputStyle} placeholder="描述" value={m.description} onChange={(e) => {
              const mats = [...form.materials]; mats[i] = { ...mats[i], description: e.target.value }; setForm({ ...form, materials: mats });
            }} />
            <input style={inputStyle} type="number" step="0.01" placeholder="基础价" value={m.base_price} onChange={(e) => {
              const mats = [...form.materials]; mats[i] = { ...mats[i], base_price: parseFloat(e.target.value) }; setForm({ ...form, materials: mats });
            }} />
            {form.materials.length > 1 && <button type="button" onClick={() => setForm({ ...form, materials: form.materials.filter((_: any, j: number) => j !== i) })} style={linkBtn}>×</button>}
          </div>
        ))}
        <button type="button" style={linkBtn} onClick={() => setForm({ ...form, materials: [...form.materials, { name: '', description: '', base_price: 0, unit: 'sheet' }] })}>+ 添加材质</button>

        <h4 style={{ marginTop: 16, marginBottom: 8 }}>工艺选项</h4>
        {form.processes.map((p: any, i: number) => (
          <div key={i} style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr auto', gap: 8, marginBottom: 8 }}>
            <input style={inputStyle} placeholder="工艺名称" value={p.name} onChange={(e) => {
              const procs = [...form.processes]; procs[i] = { ...procs[i], name: e.target.value }; setForm({ ...form, processes: procs });
            }} />
            <input style={inputStyle} placeholder="描述" value={p.description} onChange={(e) => {
              const procs = [...form.processes]; procs[i] = { ...procs[i], description: e.target.value }; setForm({ ...form, processes: procs });
            }} />
            <input style={inputStyle} type="number" step="0.01" placeholder="加价" value={p.extra_price} onChange={(e) => {
              const procs = [...form.processes]; procs[i] = { ...procs[i], extra_price: parseFloat(e.target.value) }; setForm({ ...form, processes: procs });
            }} />
            {form.processes.length > 1 && <button type="button" onClick={() => setForm({ ...form, processes: form.processes.filter((_: any, j: number) => j !== i) })} style={linkBtn}>×</button>}
          </div>
        ))}
        <button type="button" style={linkBtn} onClick={() => setForm({ ...form, processes: [...form.processes, { name: '', description: '', extra_price: 0 }] })}>+ 添加工艺</button>

        <div style={{ marginTop: 20, display: 'flex', gap: 12, justifyContent: 'flex-end' }}>
          <button type="button" onClick={onClose} style={secondaryBtn}>取消</button>
          <button type="submit" style={primaryBtn}>保存</button>
        </div>
      </form>
    </div>
  );
}

function Field({ label, value, onChange, type = 'text', required }: { label: string; value: any; onChange: (v: string) => void; type?: string; required?: boolean }) {
  return (
    <div>
      <label style={labelStyle}>{label}{required && '*'}</label>
      <input style={inputStyle} type={type} value={value} onChange={(e) => onChange(e.target.value)} required={required} />
    </div>
  );
}

const labelStyle: React.CSSProperties = { display: 'block', marginBottom: 4, fontSize: 13, color: '#374151', fontWeight: 500 };
const inputStyle: React.CSSProperties = { width: '100%', padding: '8px 10px', border: '1px solid #d1d5db', borderRadius: 6, fontSize: 14, outline: 'none', boxSizing: 'border-box' };
const primaryBtn: React.CSSProperties = { padding: '8px 16px', background: '#3b82f6', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer', fontSize: 14 };
const secondaryBtn: React.CSSProperties = { padding: '8px 16px', background: '#fff', color: '#374151', border: '1px solid #d1d5db', borderRadius: 6, cursor: 'pointer', fontSize: 14 };
const linkBtn: React.CSSProperties = { background: 'none', border: 'none', color: '#3b82f6', cursor: 'pointer', fontSize: 13, padding: 4 };

function categoryColor(c: string): string {
  const map: Record<string, string> = { '名片': '#3b82f6', '传单': '#10b981', '画册': '#f59e0b', '海报': '#6366f1', '包装盒': '#ec4899', '不干胶': '#8b5cf6', '其他': '#6b7280' };
  return map[c] || '#6b7280';
}
function categoryIcon(c: string): string {
  const map: Record<string, string> = { '名片': '🪪', '传单': '📄', '画册': '📚', '海报': '🖼️', '包装盒': '📦', '不干胶': '🏷️', '其他': '📦' };
  return map[c] || '📦';
}
