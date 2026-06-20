import { defineStore } from 'pinia'
import { ref, onUnmounted } from 'vue'
import type { RoomEnvironmentData } from '@/types'

export const useEnvironmentStore = defineStore('environment', () => {
  const roomData = ref<Map<string, RoomEnvironmentData>>(new Map())
  const connectionStatus = ref<'connecting' | 'connected' | 'disconnected'>('disconnected')
  const lastUpdateTime = ref<string>('')
  const globalForceMode = ref(false)
  const pm25ConsecutiveOverCount = ref<Map<string, number>>(new Map())
  const severelyPollutedRooms = ref<Set<string>>(new Set())
  let ws: WebSocket | null = null

  const PM25_DANGER_THRESHOLD = 150
  const CONSECUTIVE_COUNT_REQUIRED = 3

  const checkPm25ConsecutiveOver = (data: RoomEnvironmentData[]) => {
    data.forEach((item) => {
      const roomId = item.room_id
      const currentCount = pm25ConsecutiveOverCount.value.get(roomId) || 0

      if (item.pm25 > PM25_DANGER_THRESHOLD) {
        const newCount = currentCount + 1
        pm25ConsecutiveOverCount.value.set(roomId, newCount)

        if (newCount >= CONSECUTIVE_COUNT_REQUIRED) {
          if (!severelyPollutedRooms.value.has(roomId)) {
            severelyPollutedRooms.value.add(roomId)
            console.warn(`⚠️ 房间 [${item.room_name}] PM2.5 连续 ${CONSECUTIVE_COUNT_REQUIRED} 次超过 ${PM25_DANGER_THRESHOLD} μg/m³，已触发严重污染告警！`)
          }
        }
      } else {
        pm25ConsecutiveOverCount.value.set(roomId, 0)
        severelyPollutedRooms.value.delete(roomId)
      }
    })
  }

  const isRoomSeverelyPolluted = (roomId: string) => {
    return severelyPollutedRooms.value.has(roomId)
  }

  const setGlobalForceMode = (enabled: boolean) => {
    globalForceMode.value = enabled
    console.log(`全局强启模式: ${enabled ? '开启' : '关闭'}`)

    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'force_mode',
        data: { enabled }
      }))
    }
  }

  const connect = () => {
    if (ws && ws.readyState === WebSocket.OPEN) return

    connectionStatus.value = 'connecting'
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}/ws`

    ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      connectionStatus.value = 'connected'
      console.log('WebSocket connected')
    }

    ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)

        if (message.type === 'environment_data') {
          const data: RoomEnvironmentData[] = message.data
          data.forEach((item) => {
            roomData.value.set(item.room_id, item)
          })
          checkPm25ConsecutiveOver(data)
          lastUpdateTime.value = new Date().toLocaleTimeString()
        } else if (message.type === 'force_mode_update') {
          globalForceMode.value = message.data.enabled
          console.log(`收到服务器强启模式更新: ${message.data.enabled}`)
        } else if (message.type === 'force_mode_result') {
          console.log(`强启模式操作结果: ${message.data.success}`)
        }
      } catch (e) {
        console.error('WebSocket message parse error:', e)
      }
    }

    ws.onclose = () => {
      connectionStatus.value = 'disconnected'
      console.log('WebSocket disconnected, reconnecting...')
      setTimeout(connect, 3000)
    }

    ws.onerror = (error) => {
      console.error('WebSocket error:', error)
      connectionStatus.value = 'disconnected'
    }
  }

  const disconnect = () => {
    if (ws) {
      ws.close()
      ws = null
    }
  }

  const getAllRoomsData = () => {
    return Array.from(roomData.value.values())
  }

  const getRoomData = (roomId: string) => {
    return roomData.value.get(roomId)
  }

  const resetPollutedStatus = (roomId?: string) => {
    if (roomId) {
      severelyPollutedRooms.value.delete(roomId)
      pm25ConsecutiveOverCount.value.set(roomId, 0)
    } else {
      severelyPollutedRooms.value.clear()
      pm25ConsecutiveOverCount.value.clear()
    }
  }

  onUnmounted(() => {
    disconnect()
  })

  return {
    roomData,
    connectionStatus,
    lastUpdateTime,
    globalForceMode,
    severelyPollutedRooms,
    PM25_DANGER_THRESHOLD,
    CONSECUTIVE_COUNT_REQUIRED,
    connect,
    disconnect,
    getAllRoomsData,
    getRoomData,
    isRoomSeverelyPolluted,
    setGlobalForceMode,
    resetPollutedStatus,
  }
})
