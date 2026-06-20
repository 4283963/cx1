import axios from 'axios'
import type { ApiResponse, Room, Device, LinkageRule, EnvironmentData } from '@/types'

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
})

request.interceptors.response.use(
  (response) => response.data,
  (error) => {
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

export const roomApi = {
  getAll: () => request.get<any, ApiResponse<Room[]>>('/rooms'),
  getById: (id: string) => request.get<any, ApiResponse<Room>>(`/rooms/${id}`),
  create: (data: { name: string; floor: number }) => request.post<any, ApiResponse<Room>>('/rooms', data),
  update: (id: string, data: { name: string; floor: number }) => request.put<any, ApiResponse<Room>>(`/rooms/${id}`, data),
  delete: (id: string) => request.delete<any, ApiResponse<void>>(`/rooms/${id}`),
}

export const deviceApi = {
  getAll: (roomId?: string) => request.get<any, ApiResponse<Device[]>>('/devices', { params: { room_id: roomId } }),
  toggleStatus: (id: string) => request.patch<any, ApiResponse<Device>>(`/devices/${id}/toggle`),
}

export const linkageRuleApi = {
  getAll: (roomId?: string) => request.get<any, ApiResponse<LinkageRule[]>>('/linkage-rules', { params: { room_id: roomId } }),
  getById: (id: string) => request.get<any, ApiResponse<LinkageRule>>(`/linkage-rules/${id}`),
  create: (data: Omit<LinkageRule, 'id' | 'created_at' | 'updated_at'>) => request.post<any, ApiResponse<LinkageRule>>('/linkage-rules', data),
  update: (id: string, data: Partial<LinkageRule>) => request.put<any, ApiResponse<LinkageRule>>(`/linkage-rules/${id}`, data),
  toggleEnabled: (id: string) => request.patch<any, ApiResponse<LinkageRule>>(`/linkage-rules/${id}/toggle`),
  delete: (id: string) => request.delete<any, ApiResponse<void>>(`/linkage-rules/${id}`),
}

export const environmentApi = {
  getLatestAll: () => request.get<any, ApiResponse<EnvironmentData[]>>('/environment/latest'),
  getLatestByRoom: (roomId: string) => request.get<any, ApiResponse<EnvironmentData>>(`/environment/latest/${roomId}`),
  getHistory: (roomId: string, params?: { start_time?: string; end_time?: string; limit?: number }) =>
    request.get<any, ApiResponse<EnvironmentData[]>>(`/environment/history/${roomId}`, { params }),
}

export interface SystemStatus {
  force_mode: boolean
}

export const systemApi = {
  getStatus: () => request.get<any, ApiResponse<SystemStatus>>('/system/status'),
  setForceMode: (enabled: boolean) =>
    request.post<any, ApiResponse<SystemStatus>>('/system/force-mode', { enabled }),
}
