<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { NModal, NForm, NFormItem, NInput, NButton, NSpace, NSwitch, NGrid, NGridItem, NText } from 'naive-ui'
import { EXCHANGES, type ExchangeInfo } from '../../utils/exchange'
import api from '../../api/client'
import { useNotify } from '../../composables/useNotify'

const props = defineProps<{
  show: boolean
  activeExchanges: string[]
}>()

const emit = defineEmits(['update:show', 'added'])
const { success, error } = useNotify()

const step = ref(1)
const selectedExchange = ref<ExchangeInfo | null>(null)
const loading = ref(false)

const formData = ref({
  api_key: '',
  api_secret: '',
  passphrase: '',
  testnet: false
})

const availableExchanges = computed(() => 
  EXCHANGES.filter(ex => !props.activeExchanges.includes(ex.id))
)

watch(() => props.show, (val) => {
  if (val) {
    step.value = 1
    selectedExchange.value = null
    formData.value = { api_key: '', api_secret: '', passphrase: '', testnet: false }
  }
})

function selectExchange(ex: ExchangeInfo) {
  selectedExchange.value = ex
  step.value = 2
}

async function saveConfig() {
  if (!selectedExchange.value) return
  if (!formData.value.api_key || !formData.value.api_secret) {
    error('请填写 API Key 和 Secret Key')
    return
  }
  if (selectedExchange.value.id === 'okx' && !formData.value.passphrase) {
    error('OKX 需要填写 Passphrase')
    return
  }

  loading.value = true
  try {
    await api.put(`/exchanges/${selectedExchange.value.id}`, formData.value)
    success(`${selectedExchange.value.name} 配置已保存`)
    emit('added', selectedExchange.value.id)
    emit('update:show', false)
  } catch (err: any) {
    error('保存失败: ' + (err.response?.data?.error || err.message))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <n-modal 
    :show="show" 
    @update:show="$emit('update:show', $event)"
    preset="card"
    :title="step === 1 ? '添加交易所' : `配置 ${selectedExchange?.name}`"
    style="max-width: 500px"
    :mask-closable="!loading"
  >
    <!-- Step 1: Select Exchange -->
    <template v-if="step === 1">
      <n-grid :cols="2" :x-gap="16" :y-gap="16">
        <n-grid-item v-for="ex in availableExchanges" :key="ex.id">
          <div class="exchange-option" @click="selectExchange(ex)">
            <img :src="ex.icon" :alt="ex.name" class="option-icon" />
            <n-text strong>{{ ex.name }}</n-text>
          </div>
        </n-grid-item>
      </n-grid>
      
      <n-text v-if="availableExchanges.length === 0" depth="3" class="empty-hint">
        所有交易所都已添加
      </n-text>
    </template>

    <!-- Step 2: Configure Exchange -->
    <template v-else-if="selectedExchange">
      <n-space vertical :size="4" class="mb-4">
        <n-space align="center" :size="8">
          <img :src="selectedExchange.icon" :alt="selectedExchange.name" class="config-icon" />
          <n-text strong class="text-lg">{{ selectedExchange.name }}</n-text>
        </n-space>
        <n-text depth="3" class="text-sm">请输入您的 API 凭据</n-text>
      </n-space>

      <n-form label-placement="top">
        <n-form-item label="API Key" required>
          <n-input 
            v-model:value="formData.api_key"
            type="password"
            show-password-on="click"
            placeholder="请输入 API Key"
          />
        </n-form-item>
        
        <n-form-item label="Secret Key" required>
          <n-input 
            v-model:value="formData.api_secret"
            type="password"
            show-password-on="click"
            placeholder="请输入 Secret Key"
          />
        </n-form-item>

        <n-form-item v-if="selectedExchange.id === 'okx'" label="Passphrase" required>
          <n-input 
            v-model:value="formData.passphrase"
            type="password"
            show-password-on="click"
            placeholder="请输入 Passphrase"
          />
        </n-form-item>

        <n-form-item v-if="selectedExchange.id === 'binance'" label="测试网模式">
          <n-switch v-model:value="formData.testnet" />
        </n-form-item>
      </n-form>

      <n-space justify="space-between" class="mt-4">
        <n-button @click="step = 1" :disabled="loading">返回</n-button>
        <n-button type="primary" @click="saveConfig" :loading="loading">
          保存配置
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<style scoped>
.config-icon {
  width: 32px;
  height: 32px;
}
.text-lg { font-size: 1.125rem; }
.text-sm { font-size: 0.875rem; }
.mb-4 { margin-bottom: 1rem; }
.mt-4 { margin-top: 1rem; }
.empty-hint {
  text-align: center;
  padding: 2rem;
  display: block;
}
.exchange-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background: var(--n-color-embedded);
  border-radius: 12px;
  border: 1px solid var(--n-border-color);
  cursor: pointer;
  transition: all 0.3s ease;
}
.exchange-option:hover {
  border-color: var(--primary-color);
  transform: translateY(-4px);
  background: rgba(56, 189, 248, 0.05);
}
.option-icon {
  width: 48px;
  height: 48px;
  margin-bottom: 12px;
  object-fit: contain;
}
</style>
