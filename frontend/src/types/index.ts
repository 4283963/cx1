export interface Room {
  id: string
  name: string
  floor: number
  created_at: string
  updated_at: string
}

export interface Device {
  id: string
  room_id: string
  name: string
  type: string
  status: boolean
  created_at: string
  updated_at: string
}

export interface LinkageRule {
  id: string
  room_id: string
  name: string
  description: string
  trigger_type: string
  trigger_value: string
  action_type: string
  action_value: string
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface EnvironmentData {
  id: string
  room_id: string
  timestamp: string
  temp: number
  humidity: number
  pm25: number
  formaldehyde: number
}

export interface RoomEnvironmentData {
  room_id: string
  room_name: string
  timestamp: string
  temp: number
  humidity: number
  pm25: number
  formaldehyde: number
  temp_status: 'good' | 'warning' | 'danger'
  humidity_status: 'good' | 'warning' | 'danger'
  pm25_status: 'good' | 'warning' | 'danger'
  formaldehyde_status: 'good' | 'warning' | 'danger'
}

export interface WSMessage<T = any> {
  type: string
  data: T
}

export interface ApiResponse<T> {
  data: T
  message?: string
  error?: string
}

export const triggerTypeOptions = [
  { value: 'temp', label: '温度' },
  { value: 'humidity', label: '湿度' },
  { value: 'pm25', label: 'PM2.5' },
  { value: 'formaldehyde', label: '甲醛' },
  { value: 'time', label: '定时' },
  { value: 'manual', label: '手动触发' },
]

export const actionTypeOptions = [
  { value: 'light', label: '控制灯光' },
  { value: 'ac', label: '控制空调' },
  { value: 'purifier', label: '控制净化器' },
  { value: 'notify', label: '发送通知' },
]
