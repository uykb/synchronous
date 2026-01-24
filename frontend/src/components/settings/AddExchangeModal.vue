<script setup lang="ts">
import { NModal, NGrid, NGridItem, NTag, NText } from 'naive-ui'
import { EXCHANGES } from '../../utils/exchange'

const props = defineProps<{
  show: boolean
  activeExchanges: string[]
}>()

const emit = defineEmits(['update:show', 'add'])

const handleAdd = (id: string) => {
  if (!props.activeExchanges.includes(id)) {
    emit('add', id)
  }
}
</script>

<template>
  <n-modal 
    :show="show" 
    @update:show="$emit('update:show', $event)"
    preset="card"
    title="添加交易所"
    style="max-width: 500px"
  >
    <n-grid :cols="2" :x-gap="16" :y-gap="16">
      <n-grid-item v-for="ex in EXCHANGES" :key="ex.id">
        <div 
          class="exchange-option" 
          :class="{ disabled: activeExchanges.includes(ex.id) }"
          @click="handleAdd(ex.id)"
        >
          <img :src="ex.icon" :alt="ex.name" class="option-icon" />
          <n-text strong>{{ ex.name }}</n-text>
          <n-tag v-if="activeExchanges.includes(ex.id)" size="small" type="success" class="added-badge">
            已添加
          </n-tag>
        </div>
      </n-grid-item>
    </n-grid>
  </n-modal>
</template>

<style scoped>
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
  position: relative;
}

.exchange-option:hover:not(.disabled) {
  border-color: var(--primary-color);
  transform: translateY(-4px);
  background: rgba(56, 189, 248, 0.05);
}

.exchange-option.disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.option-icon {
  width: 48px;
  height: 48px;
  margin-bottom: 12px;
  object-fit: contain;
}

.added-badge {
  position: absolute;
  top: 8px;
  right: 8px;
}
</style>
