<script setup lang="ts">
import { NCard, NTag, NButton, NIcon, NSpace, NText } from 'naive-ui'
import { TrashOutline } from '@vicons/ionicons5'

defineProps<{
  rule: {
    id: string
    name: string
    symbol: string
    source: string
    targets: string[]
    enabled: boolean
  }
}>()

defineEmits(['delete'])
</script>

<template>
  <n-card class="sync-rule-card" hoverable>
    <template #header>
      <n-space vertical :size="0">
        <n-text strong class="rule-name">{{ rule.name }}</n-text>
        <n-text code depth="3" class="rule-symbol">{{ rule.symbol }}</n-text>
      </n-space>
    </template>
    
    <template #header-extra>
      <n-button quaternary circle type="error" size="small" @click="$emit('delete')">
        <template #icon>
          <n-icon><TrashOutline /></n-icon>
        </template>
      </n-button>
    </template>
    
    <n-space vertical :size="12">
      <div class="detail-item">
        <span class="label">来源</span>
        <n-tag size="small" :bordered="false" type="primary" class="value">{{ rule.source }}</n-tag>
      </div>
      <div class="detail-item">
        <span class="label">目标</span>
        <n-space :size="4">
          <n-tag 
            v-for="target in rule.targets" 
            :key="target" 
            size="small" 
            :bordered="false" 
            type="info"
          >
            {{ target }}
          </n-tag>
        </n-space>
      </div>
    </n-space>
  </n-card>
</template>

<style scoped>
.sync-rule-card {
  transition: all 0.3s ease;
}

.rule-name {
  font-size: 1rem;
}

.rule-symbol {
  font-size: 0.8rem;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.label {
  font-size: 0.85rem;
  color: var(--text-muted, #999);
}

.value {
  text-transform: capitalize;
}
</style>
