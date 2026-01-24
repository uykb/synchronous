<script setup lang="ts">
import { ref } from 'vue'
import { NCard, NButton, NTag, NIcon, NSpace, NPopconfirm, NText, NModal, NForm, NFormItem, NInput, NSwitch } from 'naive-ui'
import { CreateOutline, FlaskOutline, TrashOutline, CheckmarkCircleOutline, CloseCircleOutline } from '@vicons/ionicons5'
import api from '../../api/client'
import { useNotify } from '../../composables/useNotify'

const props = defineProps<{
  exchange: { id: string, name: string, icon: string }
  apiKeyHint: string
  testnet?: boolean
}>()

const emit = defineEmits(['deleted', 'updated'])
const { success, error } = useNotify()

const showEditModal = ref(false)
const testing = ref(false)
const testResult = ref<'success' | 'error' | null>(null)
const deleting = ref(false)

const editForm = ref({
  api_key: '',
  api_secret: '',
  passphrase: '',
  testnet: props.testnet ?? false
})

async function testConnection() {
  testing.value = true
  testResult.value = null
  try {
    const res = await api.post(`/exchanges/${props.exchange.id}/test`)
    testResult.value = res.data.success ? 'success' : 'error'
    if (res.data.success) {
      success('连接测试成功')
    } else {
      error('连接测试失败')
    }
  } catch (err: any) {
    testResult.value = 'error'
    error('测试失败: ' + (err.response?.data?.error || err.message))
  } finally {
    testing.value = false
  }
}

async function saveEdit() {
  if (!editForm.value.api_key && !editForm.value.api_secret && !editForm.value.passphrase && editForm.value.testnet === props.testnet) {
    error('请至少修改一个字段')
    return
  }

  try {
    await api.put(`/exchanges/${props.exchange.id}`, editForm.value)
    success('配置已更新')
    showEditModal.value = false
    emit('updated')
  } catch (err: any) {
    error('更新失败: ' + (err.response?.data?.error || err.message))
  }
}

async function deleteExchange() {
  deleting.value = true
  try {
    await api.delete(`/exchanges/${props.exchange.id}`)
    success(`${props.exchange.name} 已删除`)
    emit('deleted')
  } catch (err: any) {
    error('删除失败: ' + (err.response?.data?.error || err.message))
  } finally {
    deleting.value = false
  }
}

function openEditModal() {
  editForm.value = { api_key: '', api_secret: '', passphrase: '', testnet: props.testnet ?? false }
  showEditModal.value = true
}
</script>

<template>
  <n-card class="exchange-card">
    <template #header>
      <n-space align="center" :size="12">
        <img :src="exchange.icon" :alt="exchange.name" class="exchange-logo" />
        <span class="exchange-name">{{ exchange.name }}</span>
        <n-tag type="success" size="small" round>已配置</n-tag>
      </n-space>
    </template>

    <n-space vertical :size="12">
      <div class="info-row">
        <n-text depth="3">API Key</n-text>
        <n-text code>{{ apiKeyHint || '••••••••' }}</n-text>
      </div>
      
      <div v-if="exchange.id === 'binance'" class="info-row">
        <n-text depth="3">测试网</n-text>
        <n-tag :type="testnet ? 'warning' : 'default'" size="small">
          {{ testnet ? '已启用' : '未启用' }}
        </n-tag>
      </div>
    </n-space>

    <template #action>
      <n-space justify="end" :size="8">
        <n-button size="small" @click="openEditModal">
          <template #icon><n-icon><CreateOutline /></n-icon></template>
          编辑
        </n-button>
        <n-button size="small" :loading="testing" @click="testConnection">
          <template #icon>
            <n-icon>
              <CheckmarkCircleOutline v-if="testResult === 'success'" />
              <CloseCircleOutline v-else-if="testResult === 'error'" />
              <FlaskOutline v-else />
            </n-icon>
          </template>
          测试
        </n-button>
        <n-popconfirm @positive-click="deleteExchange">
          <template #trigger>
            <n-button size="small" type="error" :loading="deleting">
              <template #icon><n-icon><TrashOutline /></n-icon></template>
              删除
            </n-button>
          </template>
          确定要删除 {{ exchange.name }} 配置吗？
        </n-popconfirm>
      </n-space>
    </template>
  </n-card>

  <!-- Edit Modal -->
  <n-modal v-model:show="showEditModal" preset="card" title="编辑配置" style="max-width: 450px">
    <n-text depth="3" class="mb-4 block">留空的字段将保持原有值不变</n-text>
    
    <n-form label-placement="top">
      <n-form-item label="新 API Key">
        <n-input v-model:value="editForm.api_key" type="password" show-password-on="click" placeholder="留空保持不变" />
      </n-form-item>
      <n-form-item label="新 Secret Key">
        <n-input v-model:value="editForm.api_secret" type="password" show-password-on="click" placeholder="留空保持不变" />
      </n-form-item>
      <n-form-item v-if="exchange.id === 'okx'" label="新 Passphrase">
        <n-input v-model:value="editForm.passphrase" type="password" show-password-on="click" placeholder="留空保持不变" />
      </n-form-item>
      <n-form-item v-if="exchange.id === 'binance'" label="测试网模式">
        <n-switch v-model:value="editForm.testnet" />
      </n-form-item>
    </n-form>

    <template #action>
      <n-space justify="end">
        <n-button @click="showEditModal = false">取消</n-button>
        <n-button type="primary" @click="saveEdit">保存</n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<style scoped>
.exchange-card { height: 100%; }
.exchange-logo { width: 24px; height: 24px; object-fit: contain; }
.exchange-name { font-weight: 600; font-size: 1.1rem; }
.info-row { display: flex; justify-content: space-between; align-items: center; }
.mb-4 { margin-bottom: 1rem; }
.block { display: block; }
</style>
