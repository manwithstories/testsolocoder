import request from './request'

export interface Family {
  id: number
  name: string
  description: string
  ownerId: number
  owner?: any
  members?: FamilyMember[]
  createdAt: string
  updatedAt: string
}

export interface FamilyMember {
  id: number
  familyId: number
  userId: number
  role: string
  status: number
  user: any
}

export interface Invitation {
  id: number
  familyId: number
  family: Family
  email: string
  role: string
  status: string
  createdAt: string
  expiresAt: string
}

export function listFamilies(): Promise<Family[]> {
  return request.get('/families')
}

export function createFamily(data: { name: string; description?: string }): Promise<Family> {
  return request.post('/families', data)
}

export function getFamily(id: number): Promise<Family> {
  return request.get(`/families/${id}`)
}

export function updateFamily(id: number, data: { name?: string; description?: string }): Promise<Family> {
  return request.put(`/families/${id}`, data)
}

export function deleteFamily(id: number): Promise<void> {
  return request.delete(`/families/${id}`)
}

export function inviteMember(familyId: number, data: { email: string; role: string }): Promise<Invitation> {
  return request.post(`/families/${familyId}/invite`, data)
}

export function removeMember(familyId: number, memberId: number): Promise<void> {
  return request.delete(`/families/${familyId}/members/${memberId}`)
}

export function updateMemberRole(familyId: number, memberId: number, data: { role: string }): Promise<void> {
  return request.put(`/families/${familyId}/members/${memberId}/role`, data)
}

export function listInvitations(): Promise<Invitation[]> {
  return request.get('/invitations')
}

export function acceptInvitation(id: number): Promise<void> {
  return request.post(`/invitations/${id}/accept`)
}

export function rejectInvitation(id: number): Promise<void> {
  return request.post(`/invitations/${id}/reject`)
}
