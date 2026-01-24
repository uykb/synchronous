<script setup lang="ts">
import { computed } from 'vue'
import { NCard, NInput, NForm, NFormItem, NSwitch, NButton, NTag, NIcon, NSpace } from 'naive-ui'
import { CloseOutline } from '@vicons/ionicons5'

const props = defineProps<{
  exchange: { id: string, name: string, icon: string }
  modelValue: any
  isConfigured: boolean
}>()

const emit = defineEmits(['update:modelValue', 'remove'])

const updateField = (field: string, value: any) => {
  const newValue = { ...props.modelValue, [field]: value }
  emit('update:modelValue', newValue)
}

const showPassphrase = computed(() => props.exchange.id === 'okx')
const showTestnet = computed(() => props.exchange.id === 'binance')
</script>

<template>
  <n-card class="exchange-card" :bordered="true">
    <template #header>
      <n-space align="center" :size="12">
        <img :src="exchange.icon" :alt="exchange.name" class="exchange-logo" />
        <span class="exchange-name">{{ exchange.name }}</span>
        <n-tag :type="isConfigured ? 'success' : 'warning'" size="small" round>
          {{ isConfigured ? '已配置' : '未完成' }}
        </n-tag>
      </n-space>
    </template>
    
    <template #header-extra>
      <n-button quaternary circle size="small" @click="$emit('remove')">
        <template #icon>
          <n-icon><CloseOutline /></n-icon>
        </template>
      </n-button>
    </template>
    
    <n-form label-placement="top" :show-feedback="false">
      <n-form-item label="API Key">
        <n-input 
          type="password" 
          show-password-on="click"
          :value="modelValue.api_key"
          @update:value="updateField('api_key', $event)"
          placeholder="请输入 API Key" 
        />
      </n-form-item>
      
      <n-form-item label="Secret Key" class="mt-4">
        <n-input 
          type="password" 
          show-password-on="click"
          :value="modelValue.api_secret"
          @update:value="updateField('api_secret', $event)"
          placeholder="请输入 Secret Key" 
        />
      </n-form-item>

      <n-form-item v-if="showPassphrase" label="Passphrase (密码)" class="mt-4">
        <n-input 
          type="password" 
          show-password-on="click"
          :value="modelValue.passphrase"
          @update:value="updateField('passphrase', $event)"
          placeholder="请输入 Passphrase" 
        />
      </n-form-item>

      <n-form-item v-if="showTestnet" label="启用测试网模式" class="mt-4">
        <n-switch 
          :value="modelValue.testnet"
          @update:value="updateField('testnet', $event)"
        />
      </n-form-item>
    </n-form>
  </n-card>
</template>

<style scoped>
.exchange-card {
  height: 100%;
  transition: all 0.3s ease;
}

.exchange-card:hover {
  border-color: var(--primary-color);
  transform: translateY(-2px);
}

.exchange-logo {
  width: 24px;
  height: 24px;
  object-fit: contain;
}

.exchange-name {
  font-weight: 600;
  font-size: 1.1rem;
}

.mt-4 {
  margin-top: 1rem;
}
</style>
