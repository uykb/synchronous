<script setup lang="ts">
import { ref } from 'vue'
import { useTradingStore } from '../stores/trading'

const tradingStore = useTradingStore()

const config = ref({
  binance: {
    api_key: '',
    api_secret: '',
    testnet: false
  },
  okx: {
    api_key: '',
    api_secret: '',
    passphrase: ''
  }
})

const activeExchanges = ref([])
const showExchangeSelectModal = ref(false)

function removeExchange(exchange: string) {
  activeExchanges.value = activeExchanges.value.filter(e => e !== exchange)
}

function addExchange(exchange: string) {
  if (!activeExchanges.value.includes(exchange)) {
    activeExchanges.value.push(exchange)
  }
  showExchangeSelectModal.value = false
}

const syncItems = ref([
  { id: 1, name: 'BTC Arbitrage', symbol: 'BTC-USDT', source: 'binance', targets: ['okx', 'bybit'] },
  { id: 2, name: 'ETH Sync', symbol: 'ETH-USDT', source: 'okx', targets: ['bybit'] }
])

const showAddModal = ref(false)
const newItem = ref({
  name: '',
  symbol: '',
  source: 'binance',
  targets: [] as string[]
})

function saveConfig() {
  console.log('Saving config:', config.value)
  alert('配置保存成功！')
}

function restartBot() {
  console.log('Restarting bot...')
  if (tradingStore.isRunning) {
    tradingStore.toggleBot() // stop
    setTimeout(() => tradingStore.toggleBot(), 1000) // start
  }
  alert('机器人重启命令已发送。')
}

function removeSyncItem(id: number) {
  syncItems.value = syncItems.value.filter(item => item.id !== id)
}

function addSyncItem() {
  if (!newItem.value.name || !newItem.value.symbol) return
  
  syncItems.value.push({
    ...newItem.value,
    id: Date.now(),
    targets: [...newItem.value.targets]
  })
  
  showAddModal.value = false
  newItem.value = {
    name: '',
    symbol: '',
    source: 'binance',
    targets: []
  }
}
</script>

<template>
  <div class="settings">
    <header class="settings-header">
      <h1>设置</h1>
      <p class="subtitle">配置您的交易所凭据和同步规则</p>
    </header>
    
    <div class="settings-layout">
      <section class="settings-section">
        <div class="section-title">
          <span class="icon">🔑</span>
          <h2>交易所配置</h2>
        </div>
        
        <div class="settings-grid">
          <div v-if="activeExchanges.includes('binance')" class="exchange-card">
            <div class="exchange-header">
              <div class="exchange-info">
                <img src="https://cryptologos.cc/logos/binance-coin-bnb-logo.svg?v=040" alt="Binance" class="exchange-logo" />
                <h3>Binance</h3>
              </div>
              <button class="remove-exchange-btn" @click="removeExchange('binance')" title="移除交易所">×</button>
            </div>
            <div class="exchange-body">
              <div class="form-group">
                <label>API Key</label>
                <input type="password" v-model="config.binance.api_key" placeholder="请输入 API Key" />
              </div>
              <div class="form-group">
                <label>Secret Key</label>
                <input type="password" v-model="config.binance.api_secret" placeholder="请输入 Secret Key" />
              </div>
              <label class="checkbox-label">
                <input type="checkbox" v-model="config.binance.testnet" />
                <span class="checkbox-text">启用测试网模式</span>
              </label>
            </div>
          </div>

          <div v-if="activeExchanges.includes('okx')" class="exchange-card">
            <div class="exchange-header">
              <div class="exchange-info">
                <img src="https://logo.svgcdn.com/token-branded/okx.svg" alt="OKX" class="exchange-logo" />
                <h3>OKX</h3>
              </div>
              <button class="remove-exchange-btn" @click="removeExchange('okx')" title="移除交易所">×</button>
            </div>
            <div class="exchange-body">
              <div class="form-group">
                <label>API Key</label>
                <input type="password" v-model="config.okx.api_key" placeholder="请输入 API Key" />
              </div>
              <div class="form-group">
                <label>Secret Key</label>
                <input type="password" v-model="config.okx.api_secret" placeholder="请输入 Secret Key" />
              </div>
              <div class="form-group">
                <label>Passphrase (密码)</label>
                <input type="password" v-model="config.okx.passphrase" placeholder="请输入 Passphrase" />
              </div>
            </div>
          </div>

          <div class="add-exchange-card" @click="showExchangeSelectModal = true">
            <div class="add-plus">+</div>
            <p>添加交易所</p>
          </div>
        </div>

        <div class="actions-footer">
          <button class="primary-btn" @click="saveConfig">
            <span>💾</span> 保存更改
          </button>
          <button class="secondary-btn" @click="restartBot">
            <span>🔄</span> 重启机器人
          </button>
        </div>
      </section>

      <section class="settings-section">
        <div class="section-title">
          <span class="icon">🔄</span>
          <h2>同步规则</h2>
        </div>
        
        <div class="sync-items-grid">
          <div v-for="item in syncItems" :key="item.id" class="sync-card">
            <div class="sync-header">
              <div class="sync-title">
                <h3>{{ item.name }}</h3>
                <span class="symbol-tag">{{ item.symbol }}</span>
              </div>
              <button class="delete-btn" @click="removeSyncItem(item.id)" title="移除规则">×</button>
            </div>
            <div class="sync-details">
              <div class="detail-row">
                <span class="label">来源</span>
                <span class="value">{{ item.source }}</span>
              </div>
              <div class="detail-row">
                <span class="label">目标</span>
                <div class="targets-list">
                  <span v-for="t in item.targets" :key="t" class="target-tag">{{ t }}</span>
                </div>
              </div>
            </div>
          </div>

          <div class="add-sync-card" @click="showAddModal = true">
            <div class="add-plus">+</div>
            <p>添加新同步规则</p>
          </div>
        </div>
      </section>
    </div>

    <!-- Exchange Selection Modal -->
    <transition name="modal-fade">
      <div v-if="showExchangeSelectModal" class="modal-overlay" @click.self="showExchangeSelectModal = false">
        <div class="modal-content">
          <div class="modal-header">
            <h3>添加交易所</h3>
            <button class="close-btn" @click="showExchangeSelectModal = false">×</button>
          </div>
          <div class="modal-body">
            <div class="exchange-selector-grid">
              <div 
                class="selector-item" 
                :class="{ disabled: activeExchanges.includes('binance') }"
                @click="!activeExchanges.includes('binance') && addExchange('binance')"
              >
                <img src="https://cryptologos.cc/logos/binance-coin-bnb-logo.svg?v=040" alt="Binance" />
                <span>Binance</span>
                <div v-if="activeExchanges.includes('binance')" class="badge">已添加</div>
              </div>
              <div 
                class="selector-item" 
                :class="{ disabled: activeExchanges.includes('okx') }"
                @click="!activeExchanges.includes('okx') && addExchange('okx')"
              >
                <img src="https://logo.svgcdn.com/token-branded/okx.svg" alt="OKX" />
                <span>OKX</span>
                <div v-if="activeExchanges.includes('okx')" class="badge">已添加</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- Enhanced Modal -->
    <transition name="modal-fade">
      <div v-if="showAddModal" class="modal-overlay" @click.self="showAddModal = false">
        <div class="modal-content">
          <div class="modal-header">
            <h3>新建同步规则</h3>
            <button class="close-btn" @click="showAddModal = false">×</button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label>规则名称</label>
              <input v-model="newItem.name" placeholder="例如：BTC 套利" />
            </div>
            <div class="form-group">
              <label>交易对</label>
              <input v-model="newItem.symbol" placeholder="例如：BTC-USDT" />
            </div>
            <div class="form-group">
              <label>来源交易所</label>
              <select v-model="newItem.source">
                <option value="binance">Binance</option>
                <option value="okx">OKX</option>
              </select>
            </div>
            <div class="form-group">
              <label>目标交易所</label>
              <div class="checkbox-row">
                <label class="custom-checkbox">
                  <input type="checkbox" value="okx" v-model="newItem.targets" />
                  <span class="checkmark"></span> OKX
                </label>
                <label class="custom-checkbox">
                  <input type="checkbox" value="bybit" v-model="newItem.targets" />
                  <span class="checkmark"></span> Bybit
                </label>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button class="secondary-btn" @click="showAddModal = false">取消</button>
            <button class="primary-btn" @click="addSyncItem">创建规则</button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<style scoped>
.settings {
  animation: fadeIn 0.4s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.settings-header {
  margin-bottom: 2.5rem;
}

.subtitle {
  color: var(--text-secondary);
  margin-top: 0.5rem;
}

.settings-section {
  margin-bottom: 3.5rem;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
}

.section-title h2 {
  font-size: 1.25rem;
}

.icon {
  font-size: 1.5rem;
}

.settings-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(340px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.exchange-card {
  background: var(--surface-color);
  border-radius: var(--radius);
  border: 1px solid var(--border-color);
  padding: 1.5rem;
}

.exchange-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--border-color);
}

.exchange-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.remove-exchange-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 1.5rem;
  cursor: pointer;
  line-height: 1;
  transition: var(--transition);
}

.remove-exchange-btn:hover {
  color: var(--danger);
}

.add-exchange-card {
  border: 2px dashed var(--border-color);
  border-radius: var(--radius);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: var(--text-muted);
  min-height: 200px;
  transition: var(--transition);
}

.add-exchange-card:hover {
  border-color: var(--primary-color);
  color: var(--primary-color);
  background: rgba(56, 189, 248, 0.02);
}

.exchange-selector-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

.selector-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  background: var(--surface-hover);
  border-radius: var(--radius);
  border: 1px solid var(--border-color);
  cursor: pointer;
  transition: var(--transition);
  position: relative;
}

.selector-item:hover:not(.disabled) {
  border-color: var(--primary-color);
  transform: translateY(-4px);
  background: rgba(56, 189, 248, 0.05);
}

.selector-item img {
  width: 48px;
  height: 48px;
  margin-bottom: 1rem;
}

.selector-item span {
  font-weight: 600;
  font-size: 1.1rem;
}

.selector-item.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.selector-item .badge {
  position: absolute;
  top: 0.75rem;
  right: 0.75rem;
  background: var(--success);
  color: white;
  font-size: 0.7rem;
  padding: 0.2rem 0.5rem;
  border-radius: 1rem;
  font-weight: 700;
}

.exchange-logo {
  width: 32px;
  height: 32px;
  object-fit: contain;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-group label {
  display: block;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--text-muted);
  margin-bottom: 0.5rem;
  letter-spacing: 0.05em;
}

input[type="password"], 
input[type="text"], 
select {
  width: 100%;
  padding: 0.75rem 1rem;
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-primary);
  transition: var(--transition);
}

input:focus, select:focus {
  outline: none;
  border-color: var(--primary-color);
  background: rgba(15, 23, 42, 0.8);
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  cursor: pointer;
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.actions-footer {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

.sync-items-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1.5rem;
}

.sync-card {
  background: var(--surface-color);
  border-radius: var(--radius);
  border: 1px solid var(--border-color);
  padding: 1.25rem;
  transition: var(--transition);
}

.sync-card:hover {
  border-color: var(--primary-color);
  transform: translateY(-2px);
}

.sync-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.25rem;
}

.sync-title h3 {
  font-size: 1rem;
  margin-bottom: 0.25rem;
}

.symbol-tag {
  font-size: 0.75rem;
  color: var(--primary-color);
  font-weight: 700;
  font-family: monospace;
}

.delete-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 1.5rem;
  cursor: pointer;
  line-height: 1;
  transition: var(--transition);
}

.delete-btn:hover {
  color: var(--danger);
}

.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
  font-size: 0.875rem;
}

.detail-row .label {
  color: var(--text-muted);
}

.detail-row .value {
  color: var(--text-primary);
  font-weight: 600;
  text-transform: capitalize;
}

.targets-list {
  display: flex;
  gap: 0.5rem;
}

.target-tag {
  background: rgba(129, 140, 248, 0.1);
  color: var(--secondary-color);
  padding: 0.125rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 700;
}

.add-sync-card {
  border: 2px dashed var(--border-color);
  border-radius: var(--radius);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: var(--text-muted);
  min-height: 140px;
  transition: var(--transition);
}

.add-sync-card:hover {
  border-color: var(--primary-color);
  color: var(--primary-color);
  background: rgba(56, 189, 248, 0.02);
}

.add-plus {
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

/* Modal Styling */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(15, 23, 42, 0.9);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1.5rem;
}

.modal-content {
  background: var(--surface-color);
  width: 100%;
  max-width: 500px;
  border-radius: var(--radius);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-lg);
  overflow: hidden;
}

.modal-header {
  padding: 1.5rem;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-body {
  padding: 1.5rem;
}

.modal-footer {
  padding: 1.25rem 1.5rem;
  background: rgba(255, 255, 255, 0.02);
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
}

.checkbox-row {
  display: flex;
  gap: 1.5rem;
  margin-top: 0.5rem;
}

.custom-checkbox {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  font-size: 0.875rem;
}

.close-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 1.5rem;
  cursor: pointer;
}

.secondary-btn {
  background: var(--surface-hover);
  color: var(--text-primary);
}

.modal-fade-enter-active, .modal-fade-leave-active {
  transition: opacity 0.3s ease;
}

.modal-fade-enter-from, .modal-fade-leave-to {
  opacity: 0;
}

@media (max-width: 640px) {
  .actions-footer {
    flex-direction: column;
  }
  .actions-footer button {
    width: 100%;
  }
}
</style>

