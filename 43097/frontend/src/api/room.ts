import request from '@/utils/request'
import { Room, RoomType, PageParams, PageResult, RoomStatus } from '@/types'

export const getRoomList = (params: PageParams & { status?: RoomStatus; roomTypeId?: number; floor?: number }) => {
  return request.get<PageResult<Room>>('/rooms', { params })
}

export const getRoomById = (id: number) => {
  return request.get<Room>(`/rooms/${id}`)
}

export const createRoom = (data: Omit<Room, 'id' | 'createdAt' | 'updatedAt'>) => {
  return request.post<Room>('/rooms', data)
}

export const updateRoom = (id: number, data: Partial<Room>) => {
  return request.put<Room>(`/rooms/${id}`, data)
}

export const deleteRoom = (id: number) => {
  return request.delete(`/rooms/${id}`)
}

export const updateRoomStatus = (id: number, status: RoomStatus) => {
  return request.patch(`/rooms/${id}/status`, { status })
}

export const getRoomTypeList = (params?: PageParams) => {
  return request.get<PageResult<RoomType>>('/room-types', { params })
}

export const getAllRoomTypes = () => {
  return request.get<RoomType[]>('/room-types/all')
}

export const getRoomTypeById = (id: number) => {
  return request.get<RoomType>(`/room-types/${id}`)
}

export const createRoomType = (data: Omit<RoomType, 'id' | 'createdAt' | 'updatedAt'>) => {
  return request.post<RoomType>('/room-types', data)
}

export const updateRoomType = (id: number, data: Partial<RoomType>) => {
  return request.put<RoomType>(`/room-types/${id}`, data)
}

export const deleteRoomType = (id: number) => {
  return request.delete(`/room-types/${id}`)
}

export const getAvailableRooms = (checkInDate: string, checkOutDate: string, roomTypeId?: number) => {
  return request.get<Room[]>('/rooms/available', {
    params: { checkInDate, checkOutDate, roomTypeId }
  })
}
