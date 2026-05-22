import { useEffect, useState, useRef } from 'react'
import { useParams, Link } from 'react-router-dom'
import { MapContainer, TileLayer, Marker, Popup, Polyline, useMap } from 'react-leaflet'
import L from 'leaflet'
import { activityAPI, planAPI } from '@/services/api'
import type { PlanDetail, Activity } from '@/types'
import dayjs from 'dayjs'
import 'leaflet/dist/leaflet.css'

delete (L.Icon.Default.prototype as any)._getIconUrl
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png',
  iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png',
  shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png',
})

const activityTypeColors: Record<string, string> = {
  sightseeing: '#3B82F6',
  transport: '#EAB308',
  accommodation: '#A855F7',
  food: '#F97316',
  other: '#6B7280',
}

const activityTypeLabels: Record<string, string> = {
  sightseeing: '景点',
  transport: '交通',
  accommodation: '住宿',
  food: '餐饮',
  other: '其他',
}

function createCustomMarker(color: string, number: number) {
  return L.divIcon({
    className: 'custom-marker',
    html: `<div style="background-color: ${color}; width: 36px; height: 36px; border-radius: 50%; display: flex; align-items: center; justify-content: center; color: white; font-weight: bold; font-size: 14px; border: 3px solid white; box-shadow: 0 2px 8px rgba(0,0,0,0.3);">${number}</div>`,
    iconSize: [36, 36],
    iconAnchor: [18, 18],
    popupAnchor: [0, -18],
  })
}

function MapController({ center, zoom }: { center: [number, number]; zoom: number }) {
  const map = useMap()
  useEffect(() => {
    map.setView(center, zoom)
  }, [center, zoom, map])
  return null
}

export default function MapView() {
  const { id } = useParams<{ id: string }>()
  const [plan, setPlan] = useState<PlanDetail | null>(null)
  const [activities, setActivities] = useState<Activity[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedIndex, setSelectedIndex] = useState<number | null>(null)
  const mapRef = useRef<L.Map | null>(null)

  useEffect(() => {
    if (id) {
      loadData()
    }
  }, [id])

  const loadData = async () => {
    try {
      setLoading(true)
      const [planData, activitiesData] = await Promise.all([
        planAPI.getPlan(id!),
        activityAPI.getActivities(id!),
      ])
      setPlan(planData)
      setActivities(activitiesData as Activity[])
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      setLoading(false)
    }
  }

  const getActivitiesWithLocation = () => {
    return activities.filter((a) => a.location && a.latitude && a.longitude)
  }

  const getMapBounds = () => {
    const locations = getActivitiesWithLocation()
    if (locations.length === 0) return null

    const lats = locations.map((a) => a.latitude)
    const lngs = locations.map((a) => a.longitude)

    return {
      minLat: Math.min(...lats),
      maxLat: Math.max(...lats),
      minLng: Math.min(...lngs),
      maxLng: Math.max(...lngs),
    }
  }

  const getCenterAndZoom = () => {
    const bounds = getMapBounds()
    if (!bounds) {
      return { center: [39.9042, 116.4074] as [number, number], zoom: 4 }
    }

    const centerLat = (bounds.minLat + bounds.maxLat) / 2
    const centerLng = (bounds.minLng + bounds.maxLng) / 2

    const latDiff = bounds.maxLat - bounds.minLat
    const lngDiff = bounds.maxLng - bounds.minLng
    const maxDiff = Math.max(latDiff, lngDiff)

    let zoom = 12
    if (maxDiff > 10) zoom = 4
    else if (maxDiff > 5) zoom = 6
    else if (maxDiff > 2) zoom = 8
    else if (maxDiff > 1) zoom = 10
    else if (maxDiff > 0.5) zoom = 11
    else if (maxDiff > 0.2) zoom = 12
    else if (maxDiff > 0.1) zoom = 13
    else zoom = 14

    return { center: [centerLat, centerLng] as [number, number], zoom }
  }

  const getSortedActivities = () => {
    return getActivitiesWithLocation().sort((a, b) => {
      const dateCompare = a.date.localeCompare(b.date)
      if (dateCompare !== 0) return dateCompare
      return (a.start_time || '').localeCompare(b.start_time || '')
    })
  }

  const handleMarkerClick = (index: number) => {
    setSelectedIndex(index)
  }

  const handleListClick = (activity: Activity) => {
    const index = getSortedActivities().findIndex((a) => a.id === activity.id)
    setSelectedIndex(index)
    if (mapRef.current && activity.latitude && activity.longitude) {
      mapRef.current.setView([activity.latitude, activity.longitude], 15)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  const sortedActivities = getSortedActivities()
  const { center, zoom } = getCenterAndZoom()
  const polylineCoords = sortedActivities
    .filter((a) => a.latitude && a.longitude)
    .map((a) => [a.latitude, a.longitude] as [number, number])

  return (
    <div className="p-6">
      <div className="mb-6">
        <Link to={`/plans/${id}`} className="text-primary-600 hover:underline text-sm mb-2 inline-block">
          ← 返回计划详情
        </Link>
        <h1 className="text-2xl font-bold text-gray-900">地图视图</h1>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <div className="bg-white rounded-xl shadow-sm overflow-hidden">
            <div className="p-4 border-b border-gray-200">
              <h2 className="text-lg font-semibold">行程路线图 - {plan?.destination}</h2>
            </div>

            {sortedActivities.length === 0 ? (
              <div className="text-center py-16 text-gray-500 bg-gray-50">
                <p className="text-4xl mb-4">🗺️</p>
                <p className="mb-2">暂无带有位置信息的活动</p>
                <p className="text-sm">在添加活动时填写地点和坐标以在地图上显示</p>
              </div>
            ) : (
              <div style={{ height: '500px', width: '100%' }}>
                <MapContainer
                  center={center}
                  zoom={zoom}
                  style={{ height: '100%', width: '100%' }}
                  ref={mapRef as any}
                >
                  <MapController center={center} zoom={zoom} />
                  <TileLayer
                    attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                  />

                  {polylineCoords.length > 1 && (
                    <Polyline
                      positions={polylineCoords}
                      color="#667eea"
                      weight={3}
                      opacity={0.7}
                      dashArray="10,10"
                    />
                  )}

                  {sortedActivities.map((activity, index) => {
                    if (!activity.latitude || !activity.longitude) return null
                    const color = activityTypeColors[activity.type] || '#6B7280'
                    return (
                      <Marker
                        key={activity.id}
                        position={[activity.latitude, activity.longitude]}
                        icon={createCustomMarker(color, index + 1)}
                        eventHandlers={{
                          click: () => handleMarkerClick(index),
                        }}
                      >
                        <Popup>
                          <div className="min-w-[200px]">
                            <div className="flex items-center gap-2 mb-2">
                              <span
                                className="px-2 py-0.5 rounded text-xs text-white font-medium"
                                style={{ backgroundColor: color }}
                              >
                                {activityTypeLabels[activity.type]}
                              </span>
                              <span className="text-xs text-gray-500">
                                #{index + 1}
                              </span>
                            </div>
                            <h4 className="font-medium text-gray-900 mb-1">{activity.title}</h4>
                            <p className="text-sm text-gray-600 mb-1">📍 {activity.location}</p>
                            <p className="text-xs text-gray-400">
                              {dayjs(activity.date).format('YYYY/MM/DD')}
                              {activity.start_time && ` ${activity.start_time}`}
                            </p>
                            {activity.cost > 0 && (
                              <p className="text-sm text-orange-600 mt-1">
                                💰 {activity.cost.toLocaleString()} {activity.currency || 'CNY'}
                              </p>
                            )}
                          </div>
                        </Popup>
                      </Marker>
                    )
                  })}
                </MapContainer>
              </div>
            )}

            <div className="p-4 bg-gray-50 border-t border-gray-200">
              <p className="text-sm text-gray-600">
                📍 共 {sortedActivities.length} 个地点标记在地图上
                {plan?.destination && ` · 目的地: ${plan.destination}`}
              </p>
            </div>
          </div>
        </div>

        <div className="lg:col-span-1">
          <div className="bg-white rounded-xl shadow-sm p-6">
            <h2 className="text-lg font-semibold mb-4">地点列表</h2>
            {sortedActivities.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <p>暂无地点数据</p>
              </div>
            ) : (
              <div className="space-y-3 max-h-[500px] overflow-y-auto">
                {sortedActivities.map((activity, index) => {
                  const color = activityTypeColors[activity.type] || '#6B7280'
                  return (
                    <div
                      key={activity.id}
                      className={`p-3 rounded-lg border cursor-pointer transition-all ${
                        selectedIndex === index
                          ? 'border-primary-500 bg-primary-50 shadow-md'
                          : 'border-gray-200 hover:bg-gray-50'
                      }`}
                      onClick={() => handleListClick(activity)}
                    >
                      <div className="flex items-start gap-3">
                        <div
                          className="w-8 h-8 rounded-full flex items-center justify-center text-white text-sm font-bold flex-shrink-0"
                          style={{ backgroundColor: color }}
                        >
                          {index + 1}
                        </div>
                        <div className="flex-1 min-w-0">
                          <h4 className="font-medium text-gray-900 truncate">{activity.title}</h4>
                          <p className="text-xs text-gray-500 truncate">📍 {activity.location}</p>
                          <p className="text-xs text-gray-400 mt-1">
                            {dayjs(activity.date).format('MM/DD')}
                            {activity.start_time && ` ${activity.start_time}`}
                          </p>
                        </div>
                      </div>
                    </div>
                  )
                })}
              </div>
            )}
          </div>

          <div className="bg-white rounded-xl shadow-sm p-6 mt-6">
            <h3 className="text-sm font-semibold mb-3">图例</h3>
            <div className="space-y-2">
              {Object.entries(activityTypeLabels).map(([type, label]) => (
                <div key={type} className="flex items-center gap-2">
                  <div
                    className="w-4 h-4 rounded-full"
                    style={{ backgroundColor: activityTypeColors[type] }}
                  ></div>
                  <span className="text-sm text-gray-600">{label}</span>
                </div>
              ))}
            </div>
          </div>

          <div className="bg-white rounded-xl shadow-sm p-6 mt-6">
            <h3 className="text-sm font-semibold mb-3">地图说明</h3>
            <ul className="text-xs text-gray-500 space-y-2">
              <li>🔍 滚轮缩放地图</li>
              <li>🖱️ 拖拽移动地图</li>
              <li>📍 点击标记查看详情</li>
              <li>📋 点击列表定位到标记</li>
            </ul>
          </div>
        </div>
      </div>

      <div className="mt-6 bg-white rounded-xl shadow-sm p-6">
        <h2 className="text-lg font-semibold mb-4">所有活动的位置信息</h2>
        {activities.length === 0 ? (
          <p className="text-gray-500">暂无活动</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-200">
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">序号</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">活动</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">类型</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">地点</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">坐标</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">日期</th>
                </tr>
              </thead>
              <tbody>
                {activities
                  .sort((a, b) => {
                    const dateCompare = a.date.localeCompare(b.date)
                    if (dateCompare !== 0) return dateCompare
                    return (a.start_time || '').localeCompare(b.start_time || '')
                  })
                  .map((activity, index) => (
                    <tr key={activity.id} className="border-b border-gray-100 hover:bg-gray-50">
                      <td className="py-3 px-4 text-sm">{index + 1}</td>
                      <td className="py-3 px-4 font-medium">{activity.title}</td>
                      <td className="py-3 px-4">
                        <span
                          className="px-2 py-0.5 rounded text-xs text-white font-medium"
                          style={{ backgroundColor: activityTypeColors[activity.type] }}
                        >
                          {activityTypeLabels[activity.type]}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-sm text-gray-600">
                        {activity.location || '-'}
                      </td>
                      <td className="py-3 px-4 text-sm text-gray-500 font-mono">
                        {activity.latitude && activity.longitude
                          ? `${activity.latitude.toFixed(4)}, ${activity.longitude.toFixed(4)}`
                          : '-'}
                      </td>
                      <td className="py-3 px-4 text-sm text-gray-600">
                        {dayjs(activity.date).format('YYYY/MM/DD')}
                      </td>
                    </tr>
                  ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}
