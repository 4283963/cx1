<template>
  <div class="min-h-screen bg-gray-50">
    <div class="bg-white shadow-sm sticky top-0 z-10">
      <div class="max-w-lg mx-auto px-4 py-4 flex items-center">
        <button @click="goBack" class="mr-3 text-gray-600 hover:text-gray-800">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
          </svg>
        </button>
        <h1 class="text-xl font-bold text-gray-800">{{ isEdit ? '编辑联动规则' : '新建联动规则' }}</h1>
      </div>
    </div>

    <div class="max-w-lg mx-auto px-4 py-6">
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-sm font-medium text-gray-700 mb-2">规则名称</label>
          <input
            v-model="form.name"
            type="text"
            placeholder="请输入规则名称"
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none transition-all"
            required
          />
        </div>

        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-sm font-medium text-gray-700 mb-2">适用房间</label>
          <select
            v-model="form.room_id"
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none transition-all bg-white"
            required
          >
            <option value="" disabled>请选择房间</option>
            <option v-for="room in rooms" :key="room.id" :value="room.id">
              {{ room.name }}
            </option>
          </select>
        </div>

        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-sm font-medium text-gray-700 mb-2">规则描述</label>
          <textarea
            v-model="form.description"
            placeholder="请输入规则描述（可选）"
            rows="2"
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none transition-all resize-none"
          ></textarea>
        </div>

        <div class="bg-white rounded-xl p-4 shadow-sm">
          <h3 class="text-sm font-medium text-gray-700 mb-4">触发条件</h3>
          <div class="space-y-4">
            <div>
              <label class="block text-xs text-gray-500 mb-1">触发类型</label>
              <select
                v-model="form.trigger_type"
                class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none transition-all bg-white"
                required
              >
                <option value="" disabled>请选择触发类型</option>
                <option v-for="opt in triggerTypeOptions" :key="opt.value" :value="opt.value">
                  {{ opt.label }}
                </option>
              </select>
            </div>
            <div>
              <label class="block text-xs text-gray-500 mb-1">触发值</label>
              <input
                v-model="form.trigger_value"
                type="text"
                :placeholder="getTriggerPlaceholder()"
                class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none transition-all"
                required
              />
            </div>
          </div>
        </div>

        <div class="bg-white rounded-xl p-4 shadow-sm">
          <h3 class="text-sm font-medium text-gray-700 mb-4">执行动作</h3>
          <div class="space-y-4">
            <div>
              <label class="block text-xs text-gray-500 mb-1">动作类型</label>
              <select
                v-model="form.action_type"
                class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none transition-all bg-white"
                required
              >
                <option value="" disabled>请选择动作类型</option>
                <option v-for="opt in actionTypeOptions" :key="opt.value" :value="opt.value">
                  {{ opt.label }}
                </option>
              </select>
            </div>
            <div>
              <label class="block text-xs text-gray-500 mb-1">动作值</label>
              <input
                v-model="form.action_value"
                type="text"
                :placeholder="getActionPlaceholder()"
                class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none transition-all"
                required
              />
            </div>
          </div>
        </div>

        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="flex items-center justify-between">
            <div>
              <h3 class="text-sm font-medium text-gray-700">启用规则</h3>
              <p class="text-xs text-gray-500 mt-1">关闭后规则将暂时不生效</p>
            </div>
            <button
              type="button"
              @click="form.enabled = !form.enabled"
              :class="[
                'w-14 h-8 rounded-full transition-all relative',
                form.enabled ? 'bg-primary-500' : 'bg-gray-300'
              ]"
            >
              <span
                :class="[
                  'absolute top-1 w-6 h-6 bg-white rounded-full transition-all',
                  form.enabled ? 'left-7' : 'left-1'
                ]"
              ></span>
            </button>
          </div>
        </div>

        <div class="flex gap-3 pt-4">
          <button
            type="button"
            @click="goBack"
            class="flex-1 py-4 bg-gray-100 text-gray-700 rounded-xl font-medium hover:bg-gray-200 transition-colors"
          >
            取消
          </button>
          <button
            type="submit"
            :disabled="submitting"
            class="flex-1 py-4 bg-primary-500 text-white rounded-xl font-medium hover:bg-primary-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ submitting ? '保存中...' : '保存' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { roomApi, linkageRuleApi } from '@/api'
import type { Room, LinkageRule } from '@/types'
import { triggerTypeOptions, actionTypeOptions } from '@/types'

const router = useRouter()
const route = useRoute()
const rooms = ref<Room[]>([])
const submitting = ref(false)

const isEdit = computed(() => !!route.params.id)
const ruleId = computed(() => route.params.id as string)

const form = ref({
  room_id: '',
  name: '',
  description: '',
  trigger_type: '',
  trigger_value: '',
  action_type: '',
  action_value: '',
  enabled: true
})

const getTriggerPlaceholder = () => {
  switch (form.value.trigger_type) {
    case 'temp': return '例如：>28'
    case 'humidity': return '例如：>70'
    case 'pm25': return '例如：>35'
    case 'formaldehyde': return '例如：>0.08'
    case 'time': return '例如：08:00'
    default: return '请输入触发值'
  }
}

const getActionPlaceholder = () => {
  switch (form.value.action_type) {
    case 'light':
    case 'ac':
    case 'purifier': return 'on 或 off'
    case 'notify': return '通知消息内容'
    default: return '请输入动作值'
  }
}

const loadRooms = async () => {
  try {
    const res = await roomApi.getAll()
    rooms.value = res.data
  } catch (e) {
    console.error('Load rooms error:', e)
  }
}

const loadRule = async () => {
  try {
    const res = await linkageRuleApi.getById(ruleId.value)
    const rule = res.data
    form.value = {
      room_id: rule.room_id,
      name: rule.name,
      description: rule.description,
      trigger_type: rule.trigger_type,
      trigger_value: rule.trigger_value,
      action_type: rule.action_type,
      action_value: rule.action_value,
      enabled: rule.enabled
    }
  } catch (e) {
    console.error('Load rule error:', e)
  }
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    if (isEdit.value) {
      await linkageRuleApi.update(ruleId.value, form.value as Partial<LinkageRule>)
    } else {
      await linkageRuleApi.create(form.value as Omit<LinkageRule, 'id' | 'created_at' | 'updated_at'>)
    }
    router.back()
  } catch (e) {
    console.error('Submit error:', e)
    alert('保存失败，请重试')
  } finally {
    submitting.value = false
  }
}

const goBack = () => {
  router.back()
}

onMounted(() => {
  loadRooms()
  if (isEdit.value) {
    loadRule()
  }
})
</script>
