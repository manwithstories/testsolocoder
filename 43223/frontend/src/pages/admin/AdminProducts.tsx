import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { productApi } from '@/api/product'
import { Product, ProductStatus } from '@/types'

export default function AdminProducts() {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)

  useEffect(() => {
    loadProducts()
  }, [page])

  const loadProducts = async () => {
    setLoading(true)
    try {
      const res = await productApi.list({ page, page_size: 10 })
      setProducts(res.data?.items || [])
      setTotal(res.data?.total || 0)
    } catch {
      setProducts([])
    } finally {
      setLoading(false)
    }
  }

  const handleStatusChange = async (id: number, status: ProductStatus) => {
    try {
      await productApi.updateStatus(id, status)
      alert('状态更新成功')
      loadProducts()
    } catch (err: any) {
      alert(err.message || '更新失败')
    }
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

  if (loading) {
    return <div className="text-center py-12">加载中...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-800">商品管理</h1>
        <Link to="/roaster/products" className="btn btn-primary">+ 新建商品</Link>
      </div>

      <div className="card overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">商品</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">产地</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">价格</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">库存</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">评分</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">状态</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">操作</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {products.map((product) => (
              <tr key={product.id}>
                <td className="px-4 py-3">
                  <div className="flex items-center gap-3">
                    <div className="w-12 h-12 bg-gray-100 rounded-lg overflow-hidden flex-shrink-0">
                      {product.images?.[0] ? (
                        <img src={product.images[0].url} alt="" className="w-full h-full object-cover" />
                      ) : (
                        <div className="w-full h-full flex items-center justify-center text-coffee-300">☕</div>
                      )}
                    </div>
                    <Link to={`/products/${product.id}`} className="font-medium hover:text-coffee-600">
                      {product.name}
                    </Link>
                  </div>
                </td>
                <td className="px-4 py-3 text-sm">{product.origin}</td>
                <td className="px-4 py-3 font-medium text-coffee-600">¥{product.price}</td>
                <td className="px-4 py-3 text-sm">{product.stock}</td>
                <td className="px-4 py-3 text-sm">{product.cupping_score.toFixed(1)}</td>
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
                  <button
                    onClick={() => handleDelete(product.id)}
                    className="text-red-500 hover:underline text-sm"
                  >
                    删除
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {total > 10 && (
        <div className="flex justify-center items-center gap-2">
          <button
            onClick={() => setPage(Math.max(1, page - 1))}
            disabled={page === 1}
            className="btn btn-secondary disabled:opacity-50"
          >
            上一页
          </button>
          <span className="text-gray-600">{page} / {Math.ceil(total / 10)}</span>
          <button
            onClick={() => setPage(page + 1)}
            disabled={page >= Math.ceil(total / 10)}
            className="btn btn-secondary disabled:opacity-50"
          >
            下一页
          </button>
        </div>
      )}
    </div>
  )
}
