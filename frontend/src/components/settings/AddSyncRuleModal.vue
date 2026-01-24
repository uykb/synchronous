<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { NModal, NForm, NFormItem, NInput, NButton, NSpace, NText, NAlert } from 'naive-ui'

const props = defineProps<{
  show: boolean
  configuredExchanges: Array<{ id: string, name: string, icon: string }>
}>()

const emit = defineEmits(['update:show', 'add'])

const formValue = ref({
  name: '',
  symbol: '',
  source: '',
  targets: [] as string[],
  enabled: true
})

// Initialize source when modal opens
watch(() => props.show, (val) => {
  if (val && !formValue.value.source && props.configuredExchanges.length > 0) {
    formValue.value.source = props.configuredExchanges[0].id
  }
})

const filteredTargets = computed(() => {
  return props.configuredExchanges.filter(ex => ex.id !== formValue.value.source)
})

const toggleTarget = (id: string) => {
  const index = formValue.value.targets.indexOf(id)
  if (index === -1) {
    formValue.value.targets.push(id)
  } else {
    formValue.value.targets.splice(index, 1)
  }
}

const handleAdd = () => {
  if (!formValue.value.name || !formValue.value.symbol || !formValue.value.source || formValue.value.targets.length === 0) {
    return
  }
  emit('add', { ...formValue.value })
  // Reset form
  formValue.value = { name: '', symbol: '', source: '', targets: [], enabled: true }
}
</script>

<template>
  <n-modal 
    :show="show" 
    @update:show="$emit('update:show', $event)"
    preset="card"
    title="新建同步规则"
    style="max-width: 500px"
  >
    <div v-if="configuredExchanges.length === 0" class="mb-4">
      <n-alert type="warning" title="需要配置交易所">
        请先在主界面添加并配置至少两个交易所（一个来源，一个目标）。
      </n-alert>
    </div>

    <n-form :model="formValue" label-placement="top">
      <n-form-item label="规则名称" path="name">
        <n-input v-model:value="formValue.name" placeholder="例如：BTC 套利" />
      </n-form-item>
      
      <n-form-item label="交易对" path="symbol">
        <n-input v-model:value="formValue.symbol" placeholder="例如：BTC-USDT" />
      </n-form-item>
      
      <n-form-item label="来源交易所">
        <n-space>
          <n-button 
            v-for="ex in configuredExchanges" 
            :key="ex.id"
            :type="formValue.source === ex.id ? 'primary' : 'default'"
            round
            size="small"
            @click="formValue.source = ex.id"
          >
            {{ ex.name }}
          </n-button>
        </n-space>
      </n-form-item>
      
      <n-form-item label="目标交易所">
        <n-space>
          <n-button 
            v-for="ex in filteredTargets" 
            :key="ex.id"
            :type="formValue.targets.includes(ex.id) ? 'info' : 'default'"
            round
            size="small"
            @click="toggleTarget(ex.id)"
          >
            {{ ex.name }}
          </n-button>
        </n-space>
      </n-form-item>
    </n-form>

    <template #footer>
      <n-space justify="end">
        <n-button @click="$emit('update:show', false)">取消</n-button>
        <n-button 
          type="primary" 
          @click="handleAdd"
          :disabled="!formValue.name || !formValue.symbol || !formValue.source || formValue.targets.length === 0"
        >
          创建规则
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<style scoped>
.mb-4 {
  margin-bottom: 1rem;
}
</style>
