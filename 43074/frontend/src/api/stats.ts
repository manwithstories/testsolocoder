import request from './request'
import type { OverviewStats, MonthlyStats, HeatmapData, DurationStats, CategoryStats, TagStats } from '@/types'

export const getOverview = () => request.get<any, OverviewStats>('/stats/overview')

export const getYearlyTrend = (year: number) => request.get<any, {
  year: number
  monthly: MonthlyStats[]
}>('/stats/yearly-trend', { params: { year } })

export const getReadingHeatmap = (year: number) => request.get<any, {
  year: number
  data: HeatmapData[]
}>('/stats/heatmap', { params: { year } })

export const getDurationDistribution = () => request.get<any, DurationStats[]>('/stats/duration')

export const getCategoryStats = () => request.get<any, CategoryStats[]>('/stats/categories')

export const getTagStats = () => request.get<any, TagStats[]>('/stats/tags')
