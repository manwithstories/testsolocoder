import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { productApi } from '@/api/product'
import { Product, CreateProductRequest, RoastLevel, ProcessMethod, ProductStatus } from '@/types'
import { useAuthStore } from '@/store/auth'

export default function RoasterProducts() {
  const { user } = useAuthStore()
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [editingProduct, setEditingProduct] = useState<Product | null>(null)
  const [formData, setFormData] = useState<CreateProductRequest>({
    name: '',
    origin: '',
    farm: '',
    variety: '',
    altitude: '',
    process_method: 'washed',
    roast_level: 'medium',
    flavor_notes: '',
    cupping_score: 0,
    description: '',
    price: 0,
    weight: 250,
    stock: 0,
    status: 'draft',
  })

  useEffect(() => {
    loadProducts()
  }, [user?.id])

  const loadProducts = async () => {
    if (!user?.id) return
    setLoading(true)
    try {
      const res = await productApi.list({
        roaster_id: String(user.id),
        page_size: 50,
        status: '',
      })
      setProducts(res.data?.items || [])
    } catch {
      setProducts([])
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (editingProduct) {
        await productApi.update(editingProduct.id, formData)
        alert('更新成功')
      } else {
        await productApi.create(formData)
        alert('创建成功')
      }
      setShowForm(false)
      setEditingProduct(null)
      loadProducts()
    } catch (err: any) {
      alert(err.message || '保存失败')
    }
  }

  const handleEdit = (product: Product) => {
    setEditingProduct(product)
    setFormData({
      name: product.name,
      origin: product.origin,
      farm: product.farm || '',
      variety: product.variety || '',
      altitude: product.altitude || '',
      process_method: product.process_method,
      roast_level: product.roast_level,
      flavor_notes: product.flavor_notes || '',
      cupping_score: product.cupping_score,
      description: product.description || '',
      price: product.price,
      weight: product.weight,
      stock: product.stock,
      status: product.status,
    })
    setShowForm(true)
  }

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除该商品吗？')) return
    try {
      await productApi.delete(id)
      alert('删除成功')
      loadProducts()
    } catch (err: any) {
      alert(err.message || '删除失败')
    }
  }

  const handleStatusChange = async (id: number, status: ProductStatus) => {
    try {
      await productApi.updateStatus(id, status)
      loadProducts()
    } catch (err: any) {
      alert(err.message || '更新失败')
    }
  }

  if (loading) {
    return <div className="text-center py-12">加载中...</div>
  }

  if (showForm) {
    return (
      <div className="max-w-2xl mx-auto space-y-6">
        <h1 className="text-2xl font-bold text-gray-800">
          {editingProduct ? '编辑商品' : '新建商品'}
        </h1>
        <div className="card p-6">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="label">商品名称 *</label>
              <input
                type="text"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                className="input"
                required
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="label">产地 *</label>
                <input
                  type="text"
                  value={formData.origin}
                  onChange={(e) => setFormData({ ...formData, origin: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="label">庄园</label>
                <input
                  type="text"
                  value={formData.farm}
                  onChange={(e) => setFormData({ ...formData, farm: e.target.value })}
                  className="input"
                />
              </div>
              <div>
                <label className="label">品种</label>
                <input
                  type="text"
                  value={formData.variety}
                  onChange={(e) => setFormData({ ...formData, variety: e.target.value })}
                  className="input"
                />
              </div>
              <div>
                <label className="label">海拔</label>
                <input
                  type="text"
                  value={formData.altitude}
                  onChange={(e) => setFormData({ ...formData, altitude: e.target.value })}
                  className="input"
                />
              </div>
              <div>
                <label className="label">处理法 *</label>
                <select
                  value={formData.process_method}
                  onChange={(e) => setFormData({ ...formData, process_method: e.target.value as ProcessMethod })}
                  className="input"
                  required
                >
                  <option value="washed">水洗</option>
                  <option value="natural">日晒</option>
                  <option value="honey">蜜处理</option>
                  <option value="anaerobic">厌氧</option>
                  <option value="wet_hulled">湿刨</option>
                </select>
              </div>
              <div>
                <label className="label">烘焙度 *</label>
                <select
                  value={formData.roast_level}
                  onChange={(e) => setFormData({ ...formData, roast_level: e.target.value as RoastLevel })}
                  className="input"
                  required
                >
                  <option value="light">浅烘</option>
                  <option value="medium">中烘</option>
                  <option value="medium_dark">中深烘</option>
                  <option value="dark">深烘</option>
                </select>
              </div>
              <div>
                <label className="label">价格 (元) *</label>
                <input
                  type="number"
                  step="0.01"
                  value={formData.price}
                  onChange={(e) => setFormData({ ...formData, price: Number(e.target.value) })}
                  className="input"
                  required
                  min="0"
                />
              </div>
              <div>
                <label className="label">重量 (g) *</label>
                <input
                  type="number"
                  value={formData.weight}
                  onChange={(e) => setFormData({ ...formData, weight: Number(e.target.value) })}
                  className="input"
                  required
                  min="1"
                />
              </div>
              <div>
                <label className="label">库存</label>
                <input
                  type="number"
                  value={formData.stock}
                  onChange={(e) => setFormData({ ...formData, stock: Number(e.target.value) })}
                  className="input"
                  min="0"
                />
              </div>
              <div>
                <label className="label">杯测评分</label>
                <input
                  type="number"
                  step="0.1"
                  value={formData.cupping_score}
                  onChange={(e) => setFormData({ ...formData, cupping_score: Number(e.target.value) })}
                  className="input"
                  min="0"
                  max="10"
                />
              </div>
            </div>
            <div>
              <label className="label">风味描述</label>
              <textarea
                value={formData.flavor_notes}
                onChange={(e) => setFormData({ ...formData, flavor_notes: e.target.value })}
                className="input"
                rows={2}
              />
            </div>
            <div>
              <label className="label">商品描述</label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                className="input"
                rows={4}
              />
            </div>
            <div>
              <label className="label">状态</label>
              <select
                value={formData.status}
                onChange={(e) => setFormData({ ...formData, status: e.target.value as ProductStatus })}
                className="input"
              >
                <option value="draft">草稿</option>
                <option value="on_sale">上架</option>
                <option value="offline">下架</option>
              </select>
            </div>
            <div className="flex gap-4">
              <button type="button" onClick={() => { setShowForm(false); setEditingProduct(null) }} className="btn btn-secondary flex-1">
                取消
              </button>
              <button type="submit" className="btn btn-primary flex-1">
                {editingProduct ? '更新' : '创建'}
              </button>
            </div>
          </form>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-800">商品管理</h1>
        <button onClick={() => { setEditingProduct(null); setShowForm(true) }} className="btn btn-primary">
          + 新建商品
        </button>
      </div>

      {products.length === 0 ? (
        <div className="text-center py-12">
          <div className="text-6xl mb-4">☕</div>
          <p className="text-gray-500 mb-4">暂无商品</p>
          <button onClick={() => setShowForm(true)} className="btn btn-primary">
            创建第一个商品
          </button>
        </div>
      ) : (
        <div className="card overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">商品</th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">产地</th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">价格</th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">库存</th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">状态</th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">操作</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {products.map((product) => (
                <tr key={product.id}>
                  <td className="px-4 py-3">
                    <Link to={`/products/${product.id}`} className="font-medium hover:text-coffee-600">
                      {product.name}
                    </Link>
                  </td>
                  <td className="px-4 py-3 text-sm">{product.origin}</td>
                  <td className="px-4 py-3 font-medium text-coffee-600">¥{product.price}</td>
                  <td className="px-4 py-3 text-sm">{product.stock}</td>
                  <td className="px-4 py-3">
                    <select
                      value={product.status}
                      onChange={(e) => handleStatusChange(product.id, e.target.value as ProductStatus)}
                      className="input text-sm py-1"
                    >
                      <option value="on_sale">上架</option>
                      <option value="offline">下架</option>
                      <option value="draft">草稿</option>
                    </select>
                  </td>
                  <td className="px-4 py-3">
                    <button onClick={() => handleEdit(product)} className="text-coffee-600 hover:underline text-sm mr-2">
                      编辑
                    </button>
                    <button onClick={() => handleDelete(product.id)} className="text-red-500 hover:underline text-sm">
                      删除
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
