import { defineStore } from 'pinia'
import { ref, onUnmounted } from 'vue'
import type { RoomEnvironmentData } from '@/types'

export const useEnvironmentStore = defineStore('environment', () => {
  const roomData = ref<Map<string, RoomEnvironmentData>>(new Map())
  const connectionStatus = ref<'connecting' | 'connected' | 'disconnected'>('disconnected')
  const lastUpdateTime = ref<string>('')
  let ws: WebSocket | null = null

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
          lastUpdateTime.value = new Date().toLocaleTimeString()
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

  onUnmounted(() => {
    disconnect()
  })

  return {
    roomData,
    connectionStatus,
    lastUpdateTime,
    connect,
    disconnect,
    getAllRoomsData,
    getRoomData,
  }
})
