<template>
  <div class="min-h-screen bg-gray-50">
    <div class="bg-white shadow-sm sticky top-0 z-10">
      <div class="max-w-lg mx-auto px-4 py-4">
        <h1 class="text-xl font-bold text-gray-800">联动规则配置</h1>
        <p class="text-sm text-gray-500 mt-1">管理各房间的智能联动规则</p>
      </div>
    </div>

    <div class="max-w-lg mx-auto px-4 py-4">
      <div class="mb-4">
        <div class="flex gap-2 overflow-x-auto scrollbar-hide pb-2">
          <button
            v-for="room in rooms"
            :key="room.id"
            @click="selectedRoomId = room.id"
            :class="[
              'px-4 py-2 rounded-full text-sm font-medium whitespace-nowrap transition-all',
              selectedRoomId === room.id
                ? 'bg-primary-500 text-white shadow-md'
                : 'bg-white text-gray-600 hover:bg-gray-100'
            ]"
          >
            {{ room.name }}
          </button>
        </div>
      </div>

      <div v-if="loading" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div>
      </div>

      <div v-else-if="filteredRules.length === 0" class="text-center py-16">
        <div class="text-6xl mb-4">🔗</div>
        <p class="text-gray-500">暂无联动规则</p>
        <p class="text-gray-400 text-sm mt-1">点击右下角按钮创建第一条规则</p>
      </div>

      <div v-else class="space-y-3">
        <div
          v-for="rule in filteredRules"
          :key="rule.id"
          class="bg-white rounded-xl shadow-sm p-4 transition-all hover:shadow-md"
        >
          <div class="flex items-start justify-between">
            <div class="flex-1">
              <div class="flex items-center gap-2">
                <h3 class="font-semibold text-gray-800">{{ rule.name }}</h3>
                <span
                  :class="[
                    'px-2 py-0.5 rounded-full text-xs font-medium',
                    rule.enabled ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-500'
                  ]"
                >
                  {{ rule.enabled ? '已启用' : '已禁用' }}
                </span>
              </div>
              <p class="text-sm text-gray-500 mt-1">{{ rule.description || '暂无描述' }}</p>
              <div class="flex items-center gap-2 mt-3 text-sm">
                <span class="bg-blue-50 text-blue-600 px-2 py-1 rounded">
                  {{ getTriggerLabel(rule.trigger_type) }}: {{ rule.trigger_value }}
                </span>
                <span class="text-gray-300">→</span>
                <span class="bg-purple-50 text-purple-600 px-2 py-1 rounded">
                  {{ getActionLabel(rule.action_type) }}: {{ rule.action_value }}
                </span>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <button
                @click.stop="toggleRule(rule.id)"
                :class="[
                  'w-12 h-6 rounded-full transition-all relative',
                  rule.enabled ? 'bg-primary-500' : 'bg-gray-300'
                ]"
              >
                <span
                  :class="[
                    'absolute top-1 w-4 h-4 bg-white rounded-full transition-all',
                    rule.enabled ? 'left-7' : 'left-1'
                  ]"
                ></span>
              </button>
            </div>
          </div>
          <div class="flex gap-2 mt-4 pt-3 border-t border-gray-100">
            <button
              @click="editRule(rule.id)"
              class="flex-1 py-2 text-sm text-primary-600 hover:bg-primary-50 rounded-lg transition-colors"
            >
              编辑
            </button>
            <button
              @click="deleteRule(rule.id)"
              class="flex-1 py-2 text-sm text-red-600 hover:bg-red-50 rounded-lg transition-colors"
            >
              删除
            </button>
          </div>
        </div>
      </div>
    </div>

    <button
      @click="createRule"
      class="fixed right-6 bottom-6 w-14 h-14 bg-primary-500 text-white rounded-full shadow-lg hover:bg-primary-600 active:scale-95 transition-all flex items-center justify-center text-2xl"
    >
      +
    </button>

    <div
      v-if="showDeleteConfirm"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 px-4"
    >
      <div class="bg-white rounded-2xl p-6 max-w-sm w-full">
        <h3 class="text-lg font-semibold text-gray-800 mb-2">确认删除</h3>
        <p class="text-gray-600 mb-6">确定要删除这条联动规则吗？此操作不可撤销。</p>
        <div class="flex gap-3">
          <button
            @click="showDeleteConfirm = false"
            class="flex-1 py-3 text-gray-600 bg-gray-100 rounded-xl hover:bg-gray-200 transition-colors font-medium"
          >
            取消
          </button>
          <button
            @click="confirmDelete"
            class="flex-1 py-3 text-white bg-red-500 rounded-xl hover:bg-red-600 transition-colors font-medium"
          >
            删除
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { roomApi, linkageRuleApi } from '@/api'
import type { Room, LinkageRule } from '@/types'
import { triggerTypeOptions, actionTypeOptions } from '@/types'

const router = useRouter()
const rooms = ref<Room[]>([])
const rules = ref<LinkageRule[]>([])
const selectedRoomId = ref<string>('')
const loading = ref(false)
const showDeleteConfirm = ref(false)
const deletingRuleId = ref<string>('')

const filteredRules = computed(() => {
  if (!selectedRoomId.value) return rules.value
  return rules.value.filter(r => r.room_id === selectedRoomId.value)
})

const getTriggerLabel = (type: string) => {
  return triggerTypeOptions.find(o => o.value === type)?.label || type
}

const getActionLabel = (type: string) => {
  return actionTypeOptions.find(o => o.value === type)?.label || type
}

const loadData = async () => {
  loading.value = true
  try {
    const [roomRes, ruleRes] = await Promise.all([
      roomApi.getAll(),
      linkageRuleApi.getAll()
    ])
    rooms.value = roomRes.data
    rules.value = ruleRes.data
    if (rooms.value.length > 0 && !selectedRoomId.value) {
      selectedRoomId.value = rooms.value[0].id
    }
  } catch (e) {
    console.error('Load data error:', e)
  } finally {
    loading.value = false
  }
}

const createRule = () => {
  router.push('/mobile/rules/new')
}

const editRule = (id: string) => {
  router.push(`/mobile/rules/${id}/edit`)
}

const toggleRule = async (id: string) => {
  try {
    await linkageRuleApi.toggleEnabled(id)
    const rule = rules.value.find(r => r.id === id)
    if (rule) rule.enabled = !rule.enabled
  } catch (e) {
    console.error('Toggle rule error:', e)
  }
}

const deleteRule = (id: string) => {
  deletingRuleId.value = id
  showDeleteConfirm.value = true
}

const confirmDelete = async () => {
  try {
    await linkageRuleApi.delete(deletingRuleId.value)
    rules.value = rules.value.filter(r => r.id !== deletingRuleId.value)
  } catch (e) {
    console.error('Delete rule error:', e)
  } finally {
    showDeleteConfirm.value = false
    deletingRuleId.value = ''
  }
}

onMounted(() => {
  loadData()
})
</script>
