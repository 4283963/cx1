<template>
  <div class="min-h-screen bg-gradient-to-br from-dark-400 via-dark-300 to-dark-400 text-white overflow-hidden">
    <header class="px-8 py-4 border-b border-white/10">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold bg-gradient-to-r from-blue-400 to-cyan-400 bg-clip-text text-transparent">
            🏠 智能家居环境监控中心
          </h1>
          <p class="text-gray-400 text-sm mt-1">全屋环境数据实时监控大屏</p>
        </div>
        <div class="flex items-center gap-6">
          <div
            v-if="envStore.severelyPollutedRooms.size > 0"
            class="flex items-center gap-2 px-4 py-2 bg-red-500/20 border border-red-500/50 rounded-lg animate-pulse"
          >
            <span class="text-red-500 text-lg">⚠️</span>
            <span class="text-red-400 text-sm font-medium">
              {{ envStore.severelyPollutedRooms.size }} 个房间严重污染
            </span>
          </div>
          <div
            v-if="envStore.globalForceMode"
            class="flex items-center gap-2 px-4 py-2 bg-red-600/30 border-2 border-red-500 rounded-lg force-mode-glow"
          >
            <span class="text-red-400 text-lg">🚨</span>
            <span class="text-red-300 text-sm font-bold">
              全屋新风强启中
            </span>
          </div>
          <button
            v-if="hasDangerRoom && !envStore.globalForceMode"
            @click="handleForceMode"
            class="px-6 py-3 bg-gradient-to-r from-red-600 to-red-700 hover:from-red-500 hover:to-red-600 text-white font-bold rounded-lg shadow-lg shadow-red-500/30 transition-all hover:scale-105 active:scale-95 force-mode-glow"
          >
            🚨 一键开启全屋新风
          </button>
          <button
            v-if="envStore.globalForceMode"
            @click="handleForceMode(false)"
            class="px-6 py-3 bg-gradient-to-r from-gray-600 to-gray-700 hover:from-gray-500 hover:to-gray-600 text-white font-bold rounded-lg transition-all hover:scale-105 active:scale-95"
          >
            关闭强启模式
          </button>
          <div class="flex items-center gap-2">
            <span
              :class="[
                'w-3 h-3 rounded-full animate-pulse',
                envStore.connectionStatus === 'connected' ? 'bg-green-500' :
                envStore.connectionStatus === 'connecting' ? 'bg-yellow-500' : 'bg-red-500'
              ]"
            ></span>
            <span class="text-sm text-gray-300">
              {{ envStore.connectionStatus === 'connected' ? '实时连接中' :
                 envStore.connectionStatus === 'connecting' ? '连接中...' : '连接断开' }}
            </span>
          </div>
          <div class="text-right">
            <div class="text-2xl font-mono text-cyan-400">{{ currentTime }}</div>
            <div class="text-xs text-gray-500">最后更新: {{ envStore.lastUpdateTime || '--:--:--' }}</div>
          </div>
        </div>
      </div>
    </header>

    <main class="p-6 h-[calc(100vh-80px)]">
      <div class="grid grid-cols-12 gap-6 h-full">
        <div class="col-span-8 flex flex-col gap-6">
          <div class="grid grid-cols-5 gap-4">
            <div
              v-for="room in roomsData"
              :key="room.room_id"
              @click="selectedRoom = room"
              :class="[
                'bg-dark-200/50 backdrop-blur rounded-xl p-4 border transition-all cursor-pointer hover:scale-105',
                envStore.isRoomSeverelyPolluted(room.room_id)
                  ? 'border-red-500 border-2 danger-border-flash danger-pulse bg-red-900/20'
                  : selectedRoom?.room_id === room.room_id
                    ? 'border-cyan-500/50 shadow-lg shadow-cyan-500/20'
                    : 'border-white/10 hover:border-white/30'
              ]"
            >
              <div class="flex items-center justify-between mb-3">
                <h3 class="font-semibold text-base">{{ room.room_name }}</h3>
                <span
                  :class="[
                    'text-xs px-2 py-0.5 rounded-full',
                    getOverallStatus(room) === 'good' ? 'bg-green-500/20 text-green-400' :
                    getOverallStatus(room) === 'warning' ? 'bg-yellow-500/20 text-yellow-400' :
                    'bg-red-500/20 text-red-400'
                  ]"
                >
                  {{ getOverallStatus(room) === 'good' ? '优' : getOverallStatus(room) === 'warning' ? '良' : '差' }}
                </span>
              </div>
              <div class="space-y-2 text-sm">
                <div class="flex items-center justify-between">
                  <span class="text-gray-400">🌡️ 温度</span>
                  <span :class="getValueColor(room.temp_status)">{{ room.temp.toFixed(1) }}°C</span>
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-gray-400">💧 湿度</span>
                  <span :class="getValueColor(room.humidity_status)">{{ room.humidity.toFixed(1) }}%</span>
                </div>
                <div class="flex items-center justify-between">
                  <span :class="[
                    'flex items-center gap-1',
                    envStore.isRoomSeverelyPolluted(room.room_id) ? 'text-red-400 font-bold' : 'text-gray-400'
                  ]">
                    <span>🌫️</span>
                    <span>PM2.5</span>
                    <span
                      v-if="envStore.isRoomSeverelyPolluted(room.room_id)"
                      class="text-xs px-1.5 py-0.5 bg-red-500 text-white rounded animate-pulse"
                    >
                      严重超标
                    </span>
                  </span>
                  <span :class="[
                    getValueColor(room.pm25_status),
                    envStore.isRoomSeverelyPolluted(room.room_id) ? 'font-bold text-lg animate-pulse' : ''
                  ]">
                    {{ room.pm25.toFixed(0) }} μg/m³
                  </span>
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-gray-400">🧪 甲醛</span>
                  <span :class="getValueColor(room.formaldehyde_status)">{{ room.formaldehyde.toFixed(3) }} mg/m³</span>
                </div>
              </div>
            </div>
          </div>

          <div class="flex-1 bg-dark-200/50 backdrop-blur rounded-xl p-6 border border-white/10">
            <div class="flex items-center justify-between mb-4">
              <h2 class="text-lg font-semibold flex items-center gap-2">
                <span class="w-1 h-6 bg-cyan-500 rounded"></span>
                {{ selectedRoom?.room_name || '请选择房间' }} - 环境趋势
              </h2>
              <div class="flex gap-2">
                <button
                  v-for="metric in metricOptions"
                  :key="metric.value"
                  @click="selectedMetric = metric.value"
                  :class="[
                    'px-3 py-1 rounded-lg text-sm transition-all',
                    selectedMetric === metric.value
                      ? 'bg-cyan-500/30 text-cyan-400 border border-cyan-500/50'
                      : 'bg-white/5 text-gray-400 hover:bg-white/10'
                  ]"
                >
                  {{ metric.label }}
                </button>
              </div>
            </div>
            <div ref="chartRef" class="w-full h-[calc(100%-60px)]"></div>
          </div>
        </div>

        <div class="col-span-4 flex flex-col gap-6">
          <div class="grid grid-cols-2 gap-4">
            <div class="bg-gradient-to-br from-blue-500/20 to-blue-600/10 rounded-xl p-4 border border-blue-500/30">
              <div class="text-gray-400 text-sm mb-1">平均温度</div>
              <div class="text-3xl font-bold text-blue-400 font-mono">
                {{ avgTemp.toFixed(1) }}<span class="text-lg">°C</span>
              </div>
              <div class="text-xs text-gray-500 mt-2">全屋均值</div>
            </div>
            <div class="bg-gradient-to-br from-cyan-500/20 to-cyan-600/10 rounded-xl p-4 border border-cyan-500/30">
              <div class="text-gray-400 text-sm mb-1">平均湿度</div>
              <div class="text-3xl font-bold text-cyan-400 font-mono">
                {{ avgHumidity.toFixed(1) }}<span class="text-lg">%</span>
              </div>
              <div class="text-xs text-gray-500 mt-2">全屋均值</div>
            </div>
            <div class="bg-gradient-to-br from-purple-500/20 to-purple-600/10 rounded-xl p-4 border border-purple-500/30">
              <div class="text-gray-400 text-sm mb-1">平均PM2.5</div>
              <div class="text-3xl font-bold text-purple-400 font-mono">
                {{ avgPm25.toFixed(0) }}<span class="text-sm"> μg/m³</span>
              </div>
              <div :class="['text-xs mt-2', avgPm25 <= 35 ? 'text-green-400' : avgPm25 <= 75 ? 'text-yellow-400' : 'text-red-400']">
                {{ avgPm25 <= 35 ? '空气质量优' : avgPm25 <= 75 ? '空气质量良' : '空气质量差' }}
              </div>
            </div>
            <div class="bg-gradient-to-br from-rose-500/20 to-rose-600/10 rounded-xl p-4 border border-rose-500/30">
              <div class="text-gray-400 text-sm mb-1">平均甲醛</div>
              <div class="text-3xl font-bold text-rose-400 font-mono">
                {{ avgFormaldehyde.toFixed(3) }}<span class="text-sm"> mg/m³</span>
              </div>
              <div :class="['text-xs mt-2', avgFormaldehyde <= 0.08 ? 'text-green-400' : avgFormaldehyde <= 0.12 ? 'text-yellow-400' : 'text-red-400']">
                {{ avgFormaldehyde <= 0.08 ? '安全范围' : avgFormaldehyde <= 0.12 ? '轻度超标' : '严重超标' }}
              </div>
            </div>
          </div>

          <div class="flex-1 bg-dark-200/50 backdrop-blur rounded-xl p-6 border border-white/10">
            <h2 class="text-lg font-semibold flex items-center gap-2 mb-4">
              <span class="w-1 h-6 bg-purple-500 rounded"></span>
              全屋空气质量
            </h2>
            <div ref="gaugeRef" class="w-full h-[200px] mb-6"></div>
            <div class="space-y-3">
              <div
                v-for="room in roomsData"
                :key="'air-' + room.room_id"
                class="flex items-center justify-between bg-dark-100/50 rounded-lg p-3"
              >
                <span class="font-medium">{{ room.room_name }}</span>
                <div class="flex items-center gap-3">
                  <div class="w-32 h-2 bg-dark-300 rounded-full overflow-hidden">
                    <div
                      :class="[
                        'h-full transition-all duration-500',
                        room.pm25_status === 'good' ? 'bg-green-500' :
                        room.pm25_status === 'warning' ? 'bg-yellow-500' : 'bg-red-500'
                      ]"
                      :style="{ width: Math.min(room.pm25 / 2, 100) + '%' }"
                    ></div>
                  </div>
                  <span
                    :class="[
                      'text-sm font-mono w-16 text-right',
                      getValueColor(room.pm25_status)
                    ]"
                  >
                    {{ room.pm25.toFixed(0) }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <div class="bg-dark-200/50 backdrop-blur rounded-xl p-6 border border-white/10">
            <h2 class="text-lg font-semibold flex items-center gap-2 mb-4">
              <span class="w-1 h-6 bg-green-500 rounded"></span>
              系统状态
            </h2>
            <div class="grid grid-cols-2 gap-4">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-green-500/20 rounded-lg flex items-center justify-center">
                  <span class="text-xl">📡</span>
                </div>
                <div>
                  <div class="text-sm text-gray-400">在线网关</div>
                  <div class="font-bold text-green-400">{{ roomsData.length }} / {{ roomsData.length }}</div>
                </div>
              </div>
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-blue-500/20 rounded-lg flex items-center justify-center">
                  <span class="text-xl">📊</span>
                </div>
                <div>
                  <div class="text-sm text-gray-400">数据频率</div>
                  <div class="font-bold text-blue-400">1秒/次</div>
                </div>
              </div>
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-purple-500/20 rounded-lg flex items-center justify-center">
                  <span class="text-xl">🔗</span>
                </div>
                <div>
                  <div class="text-sm text-gray-400">联动规则</div>
                  <div class="font-bold text-purple-400">{{ ruleCount }} 条</div>
                </div>
              </div>
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-cyan-500/20 rounded-lg flex items-center justify-center">
                  <span class="text-xl">⚡</span>
                </div>
                <div>
                  <div class="text-sm text-gray-400">运行时间</div>
                  <div class="font-bold text-cyan-400">{{ uptime }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import * as echarts from 'echarts'
import { useEnvironmentStore } from '@/stores/useEnvironmentStore'
import { linkageRuleApi } from '@/api'
import type { RoomEnvironmentData } from '@/types'

const envStore = useEnvironmentStore()
const chartRef = ref<HTMLElement>()
const gaugeRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null
let gaugeChart: echarts.ECharts | null = null

const selectedRoom = ref<RoomEnvironmentData | null>(null)
const selectedMetric = ref('temp')
const currentTime = ref('')
const ruleCount = ref(0)
const startTime = ref(Date.now())

const metricOptions = [
  { value: 'temp', label: '温度' },
  { value: 'humidity', label: '湿度' },
  { value: 'pm25', label: 'PM2.5' },
  { value: 'formaldehyde', label: '甲醛' },
]

const historyData = ref<Map<string, number[]>>(new Map([
  ['temp', []],
  ['humidity', []],
  ['pm25', []],
  ['formaldehyde', []],
  ['time', []],
]))

const roomsData = computed(() => envStore.getAllRoomsData())

const avgTemp = computed(() => {
  const data = roomsData.value
  if (data.length === 0) return 0
  return data.reduce((sum, d) => sum + d.temp, 0) / data.length
})

const avgHumidity = computed(() => {
  const data = roomsData.value
  if (data.length === 0) return 0
  return data.reduce((sum, d) => sum + d.humidity, 0) / data.length
})

const avgPm25 = computed(() => {
  const data = roomsData.value
  if (data.length === 0) return 0
  return data.reduce((sum, d) => sum + d.pm25, 0) / data.length
})

const avgFormaldehyde = computed(() => {
  const data = roomsData.value
  if (data.length === 0) return 0
  return data.reduce((sum, d) => sum + d.formaldehyde, 0) / data.length
})

const hasDangerRoom = computed(() => {
  return envStore.severelyPollutedRooms.size > 0
})

const uptime = computed(() => {
  const diff = Date.now() - startTime.value
  const hours = Math.floor(diff / 3600000)
  const minutes = Math.floor((diff % 3600000) / 60000)
  const seconds = Math.floor((diff % 60000) / 1000)
  return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
})

const getOverallStatus = (room: RoomEnvironmentData): 'good' | 'warning' | 'danger' => {
  const statuses = [room.temp_status, room.humidity_status, room.pm25_status, room.formaldehyde_status]
  if (statuses.some(s => s === 'danger')) return 'danger'
  if (statuses.some(s => s === 'warning')) return 'warning'
  return 'good'
}

const getValueColor = (status: string) => {
  return status === 'good' ? 'text-green-400' : status === 'warning' ? 'text-yellow-400' : 'text-red-400'
}

const handleForceMode = async (enable: boolean = true) => {
  if (enable) {
    const confirmed = window.confirm(
      `⚠️ 即将开启全屋新风强启模式\n\n` +
      `检测到 ${envStore.severelyPollutedRooms.size} 个房间 PM2.5 连续严重超标\n` +
      `开启后将强制启动所有房间的新风系统和空气净化器\n\n` +
      `是否确认开启？`
    )
    if (!confirmed) return
  }

  envStore.setGlobalForceMode(enable)

  if (enable) {
    console.log('🚨 全屋新风强启模式已开启')
    alert('✅ 全屋新风强启模式已开启！\n\n所有房间的新风系统和空气净化器已强制启动。')
  } else {
    console.log('✅ 全屋新风强启模式已关闭')
  }
}

const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', { hour12: false })
}

const initChart = () => {
  if (!chartRef.value) return
  chart = echarts.init(chartRef.value)
  updateChart()
}

const updateChart = () => {
  if (!chart) return

  const metric = selectedMetric.value
  const metricLabel = metricOptions.find(o => o.value === metric)?.label || ''
  const data = historyData.value.get(metric) || []
  const timeLabels = historyData.value.get('time') || []

  const colors: Record<string, string[]> = {
    temp: ['#60a5fa', '#3b82f6'],
    humidity: ['#22d3ee', '#06b6d4'],
    pm25: ['#a78bfa', '#8b5cf6'],
    formaldehyde: ['#fb7185', '#f43f5e'],
  }

  const color = colors[metric] || ['#60a5fa', '#3b82f6']

  chart.setOption({
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(15, 23, 42, 0.9)',
      borderColor: 'rgba(255,255,255,0.1)',
      textStyle: { color: '#fff' },
      formatter: (params: any) => {
        const p = params[0]
        return `${p.axisValue}<br/>${metricLabel}: ${p.value}`
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      top: '10%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: timeLabels,
      axisLine: { lineStyle: { color: 'rgba(255,255,255,0.1)' } },
      axisLabel: { color: 'rgba(255,255,255,0.5)', fontSize: 11 }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      splitLine: { lineStyle: { color: 'rgba(255,255,255,0.05)' } },
      axisLabel: { color: 'rgba(255,255,255,0.5)', fontSize: 11 }
    },
    series: [{
      name: metricLabel,
      type: 'line',
      smooth: true,
      symbol: 'none',
      sampling: 'lttb',
      data: data,
      lineStyle: { color: color[0], width: 2 },
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: color[0] + '40' },
          { offset: 1, color: color[0] + '00' }
        ])
      }
    }]
  }, true)
}

const initGauge = () => {
  if (!gaugeRef.value) return
  gaugeChart = echarts.init(gaugeRef.value)
  updateGauge()
}

const updateGauge = () => {
  if (!gaugeChart) return

  const aqi = Math.min(avgPm25.value * 2, 500)
  let level = '优'
  let color = '#10b981'
  if (aqi > 300) { level = '严重污染'; color = '#7c3aed' }
  else if (aqi > 200) { level = '重度污染'; color = '#dc2626' }
  else if (aqi > 150) { level = '中度污染'; color = '#ea580c' }
  else if (aqi > 100) { level = '轻度污染'; color = '#eab308' }
  else if (aqi > 50) { level = '良'; color = '#84cc16' }

  gaugeChart.setOption({
    backgroundColor: 'transparent',
    series: [{
      type: 'gauge',
      startAngle: 180,
      endAngle: 0,
      min: 0,
      max: 500,
      splitNumber: 5,
      center: ['50%', '70%'],
      radius: '100%',
      axisLine: {
        lineStyle: {
          width: 15,
          color: [
            [0.1, '#10b981'],
            [0.2, '#84cc16'],
            [0.3, '#eab308'],
            [0.4, '#ea580c'],
            [0.6, '#dc2626'],
            [1, '#7c3aed']
          ]
        }
      },
      pointer: {
        icon: 'path://M12.8,0.7l12,40.1H0.7L12.8,0.7z',
        length: '60%',
        width: 10,
        itemStyle: { color: color }
      },
      axisTick: { show: false },
      splitLine: { show: false },
      axisLabel: { show: false },
      detail: {
        valueAnimation: true,
        formatter: '{value}',
        color: color,
        fontSize: 32,
        fontWeight: 'bold',
        offsetCenter: [0, '10%']
      },
      title: {
        offsetCenter: [0, '50%'],
        fontSize: 14,
        color: '#9ca3af'
      },
      data: [{ value: Math.round(aqi), name: level }]
    }]
  }, true)
}

const updateHistoryData = () => {
  if (!selectedRoom.value) return

  const room = selectedRoom.value
  const now = new Date()
  const timeStr = now.toLocaleTimeString('zh-CN', { hour12: false })

  const metrics = ['temp', 'humidity', 'pm25', 'formaldehyde'] as const
  metrics.forEach(m => {
    const arr = historyData.value.get(m) || []
    arr.push(room[m])
    if (arr.length > 60) arr.shift()
    historyData.value.set(m, arr)
  })

  const timeArr = historyData.value.get('time') || []
  timeArr.push(timeStr)
  if (timeArr.length > 60) timeArr.shift()
  historyData.value.set('time', timeArr)
}

watch(selectedRoom, () => {
  historyData.value = new Map([
    ['temp', []],
    ['humidity', []],
    ['pm25', []],
    ['formaldehyde', []],
    ['time', []],
  ])
  nextTick(updateChart)
})

watch(selectedMetric, updateChart)

watch(roomsData, () => {
  if (!selectedRoom.value && roomsData.value.length > 0) {
    selectedRoom.value = roomsData.value[0]
  }
  updateHistoryData()
  nextTick(() => {
    updateChart()
    updateGauge()
  })
}, { deep: true })

let timeInterval: number
let dataInterval: number

onMounted(async () => {
  envStore.connect()

  updateTime()
  timeInterval = window.setInterval(updateTime, 1000)

  try {
    const res = await linkageRuleApi.getAll()
    ruleCount.value = res.data.length
  } catch (e) {
    console.error('Load rule count error:', e)
  }

  await nextTick()
  initChart()
  initGauge()

  const handleResize = () => {
    chart?.resize()
    gaugeChart?.resize()
  }
  window.addEventListener('resize', handleResize)

  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
    chart?.dispose()
    gaugeChart?.dispose()
  })
})

onUnmounted(() => {
  clearInterval(timeInterval)
  clearInterval(dataInterval)
  envStore.disconnect()
})
</script>
