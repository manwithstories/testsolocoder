import { useQuery } from '@tanstack/react-query'
import { useState } from 'react'
import { Search, Filter, Leaf } from 'lucide-react'
import { plantApi } from '@/api'

export default function PlantsPage() {
  const [search, setSearch] = useState('')
  const [category, setCategory] = useState('')
  const [difficulty, setDifficulty] = useState('')

  const { data, isLoading } = useQuery({
    queryKey: ['plants', search, category, difficulty],
    queryFn: () =>
      plantApi.getAll({
        search: search || undefined,
        category: category || undefined,
        difficulty: difficulty || undefined,
        page_size: 50,
      }),
  })

  const plants = data?.data?.plants || []
  const categories = ['蔬菜', '水果', '草本', '花卉', '谷物', '豆类']
  const difficulties = ['简单', '中等', '困难']

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">植物数据库</h1>
        <p className="text-gray-500">浏览和搜索各类植物信息</p>
      </div>

      {/* Search and Filter */}
      <div className="card">
        <div className="card-body space-y-4">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              type="text"
              placeholder="搜索植物名称、描述..."
              className="input pl-10"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
          </div>
          <div className="flex flex-wrap gap-4">
            <div className="flex items-center gap-2">
              <Filter className="w-4 h-4 text-gray-400" />
              <span className="text-sm text-gray-500">分类：</span>
              <select
                className="input w-auto py-1 text-sm"
                value={category}
                onChange={(e) => setCategory(e.target.value)}
              >
                <option value="">全部</option>
                {categories.map((cat) => (
                  <option key={cat} value={cat}>
                    {cat}
                  </option>
                ))}
              </select>
            </div>
            <div className="flex items-center gap-2">
              <span className="text-sm text-gray-500">难度：</span>
              <select
                className="input w-auto py-1 text-sm"
                value={difficulty}
                onChange={(e) => setDifficulty(e.target.value)}
              >
                <option value="">全部</option>
                {difficulties.map((d) => (
                  <option key={d} value={d}>
                    {d}
                  </option>
                ))}
              </select>
            </div>
          </div>
        </div>
      </div>

      {/* Plant List */}
      {isLoading ? (
        <div className="text-center py-12 text-gray-500">加载中...</div>
      ) : plants.length === 0 ? (
        <div className="card text-center py-12">
          <Leaf className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <p className="text-gray-500">没有找到匹配的植物</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {plants.map((plant: any) => (
            <div key={plant.id} className="card hover:shadow-md transition-shadow group">
              <div className="aspect-video bg-garden-50 flex items-center justify-center overflow-hidden">
                {plant.image_url ? (
                  <img
                    src={plant.image_url}
                    alt={plant.name}
                    className="w-full h-full object-cover group-hover:scale-105 transition-transform"
                  />
                ) : (
                  <Leaf className="w-12 h-12 text-garden-300" />
                )}
              </div>
              <div className="card-body">
                <div className="flex items-start justify-between mb-2">
                  <h3 className="font-semibold text-gray-900">{plant.name}</h3>
                  {plant.difficulty && (
                    <span
                      className={`badge ${
                        plant.difficulty === '简单'
                          ? 'bg-green-100 text-green-700'
                          : plant.difficulty === '中等'
                          ? 'bg-amber-100 text-amber-700'
                          : 'bg-red-100 text-red-700'
                      }`}
                    >
                      {plant.difficulty}
                    </span>
                  )}
                </div>
                {plant.latin_name && (
                  <p className="text-sm text-gray-500 italic">{plant.latin_name}</p>
                )}
                {plant.category && (
                  <p className="text-sm text-garden-600 mt-1">{plant.category}</p>
                )}
                <div className="mt-3 space-y-1 text-sm text-gray-600">
                  {plant.growth_cycle > 0 && (
                    <p>生长周期：{plant.growth_cycle} 天</p>
                  )}
                  {plant.water_frequency && <p>浇水：{plant.water_frequency}</p>}
                  {plant.sunlight_need && <p>光照：{plant.sunlight_need}</p>}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
