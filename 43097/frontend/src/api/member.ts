import request from '@/utils/request'
import { Member, PageParams, PageResult, MemberLevel } from '@/types'

export const getMemberList = (params: PageParams & { level?: MemberLevel; status?: boolean }) => {
  return request.get<PageResult<Member>>('/members', { params })
}

export const getMemberById = (id: number) => {
  return request.get<Member>(`/members/${id}`)
}

export const getMemberByPhone = (phone: string) => {
  return request.get<Member>(`/members/phone/${phone}`)
}

export const getMemberByNo = (memberNo: string) => {
  return request.get<Member>(`/members/no/${memberNo}`)
}

export const createMember = (data: Omit<Member, 'id' | 'memberNo' | 'totalSpent' | 'totalStays' | 'createdAt' | 'updatedAt'>) => {
  return request.post<Member>('/members', data)
}

export const updateMember = (id: number, data: Partial<Member>) => {
  return request.put<Member>(`/members/${id}`, data)
}

export const updateMemberStatus = (id: number, status: boolean) => {
  return request.patch(`/members/${id}/status`, { status })
}

export const updateMemberLevelById = (id: number, level: MemberLevel) => {
  return request.patch(`/members/${id}/level`, { level })
}

export const rechargeMemberBalance = (id: number, amount: number) => {
  return request.post(`/members/${id}/recharge`, { amount })
}

export const consumeMemberBalance = (id: number, amount: number) => {
  return request.post(`/members/${id}/consume`, { amount })
}

export const addMemberPoints = (id: number, points: number) => {
  return request.post(`/members/${id}/points/add`, { points })
}

export const redeemMemberPoints = (id: number, points: number) => {
  return request.post(`/members/${id}/points/redeem`, { points })
}

export const getMemberConsumptionHistory = (id: number, params?: PageParams) => {
  return request.get(`/members/${id}/history`, { params })
}

export interface MemberLevelData {
  id?: number
  name: string
  discountRate: number
  pointsRate: number
  minPoints: number
  maxPoints: number
  createdAt?: string
  updatedAt?: string
}

export const getMemberLevelList = (params?: PageParams) => {
  return request.get<PageResult<MemberLevelData>>('/members/levels', { params })
}

export const getMemberLevelById = (id: number) => {
  return request.get<MemberLevelData>(`/members/levels/${id}`)
}

export const createMemberLevel = (data: Omit<MemberLevelData, 'id' | 'createdAt' | 'updatedAt'>) => {
  return request.post<MemberLevelData>('/members/levels', data)
}

export const updateMemberLevel = (id: number, data: Partial<MemberLevelData>) => {
  return request.put<MemberLevelData>(`/members/levels/${id}`, data)
}

export const deleteMemberLevel = (id: number) => {
  return request.delete(`/members/levels/${id}`)
}
