<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { 
  NSpace, NGrid, NGridItem, NButton, NIcon, NText, 
  NSkeleton, NCard, NInput, NInputGroup
} from 'naive-ui'
import { 
  CopyOutline, CheckmarkOutline, KeyOutline, 
  SyncOutline, AddOutline, RefreshOutline, SaveOutline
} from '@vicons/ionicons5'
import axios from 'axios'
import { useTradingStore } from '../stores/trading'
import api from '../api/client'
import { useNotify } from '../composables/useNotify'
import { useClipboard } from '../utils/clipboard'
import { EXCHANGES, getExchangeById } from '../utils/exchange'

// Components
import ExchangeCard from '../components/settings/ExchangeCard.vue'
import SyncRuleCard from '../components/settings/SyncRuleCard.vue'
import AddExchangeModal from '../components/settings/AddExchangeModal.vue'
import AddSyncRuleModal from '../components/settings/AddSyncRuleModal.vue'

interface SyncItem {
  id: string
  name: string
  symbol: string
  source: string
  targets: string[]
  enabled: boolean
}

const tradingStore = useTradingStore()
const { success, error } = useNotify()
const { copied, copyToClipboard } = useClipboard()

const serverIp = ref('正在获取...')
const loading = ref(true)
const activeExchanges = ref<string[]>([])
const showExchangeSelectModal = ref(false)
const showAddModal = ref(false)
const syncItems = ref<SyncItem[]>([])

const config = ref({
  binance: { api_key: '', api_secret: '', testnet: false },
  okx: { api_key: '', api_secret: '', passphrase: '' },
  bybit: { api_key: '', api_secret: '' },
  sync: { enabled: true, check_interval_ms: 5000 }
})

const configuredExchanges = computed(() => {
  return activeExchanges.value
    .map(id => getExchangeById(id))
    .filter((ex): ex is NonNullable<typeof ex> => !!ex)
})

async function fetchIp() {
  try {
    const response = await axios.get('/api/system/ip')
    serverIp.value = response.data.ip
  } catch (err) {
    serverIp.value = '获取失败'
  }
}

const fetchConfig = async () => {
  try {
    const res = await api.get('/config')
    if (res.data.binance) config.value.binance.testnet = res.data.binance.testnet
    if (res.data.sync) config.value.sync = res.data.sync
    
    activeExchanges.value = []
    if (res.data.binance?.enabled) activeExchanges.value.push('binance')
    if (res.data.okx?.enabled) activeExchanges.value.push('okx')
    if (res.data.bybit?.enabled) activeExchanges.value.push('bybit')
  } catch (err) {
    error('获取配置失败')
  }
}

const fetchSyncItems = async () => {
  try {
    const res = await api.get('/sync-items')
    syncItems.value = res.data || []
  } catch (err) {
    error('获取同步规则失败')
  }
}

onMounted(async () => {
  loading.value = true
  await Promise.all([fetchIp(), fetchConfig(), fetchSyncItems()])
  loading.value = false
})

function copyIp() {
  if (serverIp.value === '正在获取...' || serverIp.value === '获取失败') return
  copyToClipboard(serverIp.value).then(res => {
    if (res) success('IP 已复制到剪贴板')
  })
}

const isExchangeConfigured = (id: string) => {
  const c = config.value[id as keyof typeof config.value] as any
  if (!c) return false
  if (id === 'okx') return !!c.api_key && !!c.api_secret && !!c.passphrase
  return !!c.api_key && !!c.api_secret
}

function removeExchange(id: string) {
  activeExchanges.value = activeExchanges.value.filter(e => e !== id)
}

function addExchange(id: string) {
  if (!activeExchanges.value.includes(id)) activeExchanges.value.push(id)
  showExchangeSelectModal.value = false
}

const saveConfig = async () => {
  try {
    await api.put('/config', { ...config.value })
    success('配置已保存！需要重启服务以应用新配置。')
  } catch (err: any) {
    error('保存失败: ' + err.message)
  }
}

function restartBot() {
  if (tradingStore.isRunning) {
    tradingStore.toggleBot()
    setTimeout(() => tradingStore.toggleBot(), 1000)
  }
  success('机器人重启命令已发送。')
}

const removeSyncItem = async (id: string) => {
  try {
    await api.delete(`/sync-items/${id}`)
    syncItems.value = syncItems.value.filter(item => item.id !== id)
    success('同步规则已删除')
  } catch (err: any) {
    error('删除失败: ' + err.message)
  }
}

const handleAddSyncRule = async (newRule: any) => {
  const item = { ...newRule, id: Date.now().toString() }
  try {
    const res = await api.post('/sync-items', item)
    syncItems.value.push(res.data)
    showAddModal.value = false
    success('同步规则已添加')
  } catch (err: any) {
    error('添加失败: ' + err.message)
  }
}
</script>

<template>
  <div class="settings">
    <header class="mb-8">
      <n-space justify="space-between" align="end">
        <div>
          <n-text tag="h1" class="text-3xl font-bold m-0">设置</n-text>
          <n-text depth="3">配置您的交易所凭据和同步规则</n-text>
        </div>
        
        <n-card size="small" class="ip-card">
          <n-space vertical :size="4">
            <n-text depth="3" class="text-xs font-bold uppercase tracking-wider">
              服务器 IP (用于添加白名单)
            </n-text>
            <n-input-group>
              <n-input :value="serverIp" readonly placeholder="获取中..." style="width: 160px" />
              <n-button @click="copyIp" :type="copied ? 'success' : 'default'">
                <template #icon>
                  <n-icon>
                    <CheckmarkOutline v-if="copied" />
                    <CopyOutline v-else />
                  </n-icon>
                </template>
              </n-button>
            </n-input-group>
          </n-space>
        </n-card>
      </n-space>
    </header>

    <div v-if="loading">
      <n-grid :cols="2" :x-gap="24" :y-gap="24">
        <n-grid-item v-for="i in 4" :key="i">
          <n-skeleton height="200px" border-radius="8px" />
        </n-grid-item>
      </n-grid>
    </div>

    <n-space v-else vertical :size="48">
      <!-- Exchange Section -->
      <section>
        <n-space align="center" :size="12" class="mb-4">
          <n-icon size="24" color="var(--primary-color)"><KeyOutline /></n-icon>
          <n-text tag="h2" class="text-xl font-bold m-0">交易所配置</n-text>
        </n-space>

        <n-grid cols="1 s:1 m:2 l:3" responsive="screen" :x-gap="20" :y-gap="20">
          <n-grid-item v-for="id in activeExchanges" :key="id">
            <ExchangeCard 
              v-if="getExchangeById(id)"
              :exchange="getExchangeById(id)!"
              v-model="config[id as keyof typeof config]"
              :is-configured="isExchangeConfigured(id)"
              @remove="removeExchange(id)"
            />
          </n-grid-item>
          
          <n-grid-item>
            <div class="add-card" @click="showExchangeSelectModal = true">
              <n-icon size="32"><AddOutline /></n-icon>
              <n-text depth="3">添加交易所</n-text>
            </div>
          </n-grid-item>
        </n-grid>

        <n-space justify="end" :size="16" class="mt-6">
          <n-button secondary @click="restartBot">
            <template #icon><n-icon><RefreshOutline /></n-icon></template>
            重启机器人
          </n-button>
          <n-button type="primary" @click="saveConfig">
            <template #icon><n-icon><SaveOutline /></n-icon></template>
            保存更改
          </n-button>
        </n-space>
      </section>

      <!-- Sync Rules Section -->
      <section>
        <n-space align="center" :size="12" class="mb-4">
          <n-icon size="24" color="var(--primary-color)"><SyncOutline /></n-icon>
          <n-text tag="h2" class="text-xl font-bold m-0">同步规则</n-text>
        </n-space>

        <n-grid cols="1 s:1 m:2 l:4" responsive="screen" :x-gap="20" :y-gap="20">
          <n-grid-item v-for="item in syncItems" :key="item.id">
            <SyncRuleCard 
              :rule="item" 
              @delete="removeSyncItem(item.id)" 
            />
          </n-grid-item>
          
          <n-grid-item>
            <div class="add-card small" @click="showAddModal = true">
              <n-icon size="24"><AddOutline /></n-icon>
              <n-text depth="3">添加同步规则</n-text>
            </div>
          </n-grid-item>
        </n-grid>
      </section>
    </n-space>

    <!-- Modals -->
    <AddExchangeModal 
      v-model:show="showExchangeSelectModal" 
      :active-exchanges="activeExchanges"
      @add="addExchange"
    />

    <AddSyncRuleModal 
      v-model:show="showAddModal"
      :configured-exchanges="configuredExchanges"
      @add="handleAddSyncRule"
    />
  </div>
</template>

<style scoped>
.mb-8 { margin-bottom: 2rem; }
.mb-4 { margin-bottom: 1rem; }
.mt-6 { margin-top: 1.5rem; }

.ip-card {
  background: rgba(56, 189, 248, 0.05);
  border: 1px solid rgba(56, 189, 248, 0.2);
}

.add-card {
  height: 100%;
  min-height: 200px;
  border: 2px dashed var(--n-border-color);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background: transparent;
}

.add-card.small {
  min-height: 140px;
}

.add-card:hover {
  border-color: var(--primary-color);
  background: rgba(56, 189, 248, 0.02);
}

.add-card:hover :deep(.n-text) {
  color: var(--primary-color);
}

.text-3xl { font-size: 1.875rem; line-height: 2.25rem; }
.text-xl { font-size: 1.25rem; line-height: 1.75rem; }
.text-xs { font-size: 0.75rem; line-height: 1rem; }
.font-bold { font-weight: 700; }
.m-0 { margin: 0; }
.uppercase { text-transform: uppercase; }
.tracking-wider { letter-spacing: 0.05em; }
</style>
