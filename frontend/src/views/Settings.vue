<template>
  <div class="settings">
    <h1>Settings</h1>
    
    <section class="settings-section">
      <h2>Exchanges Configuration</h2>
      <div class="settings-card">
        <div class="exchange-group">
          <h3>Binance</h3>
          <div class="setting-item">
            <label>API Key</label>
            <input type="password" v-model="config.binance.api_key" />
          </div>
          <div class="setting-item">
            <label>Secret Key</label>
            <input type="password" v-model="config.binance.api_secret" />
          </div>
          <div class="setting-item">
            <label><input type="checkbox" v-model="config.binance.testnet" /> Testnet</label>
          </div>
        </div>

        <div class="exchange-group">
          <h3>OKX</h3>
          <div class="setting-item">
            <label>API Key</label>
            <input type="password" v-model="config.okx.api_key" />
          </div>
          <div class="setting-item">
            <label>Secret Key</label>
            <input type="password" v-model="config.okx.api_secret" />
          </div>
          <div class="setting-item">
            <label>Passphrase</label>
            <input type="password" v-model="config.okx.passphrase" />
          </div>
        </div>
        
        <button class="primary-btn" @click="saveConfig">Save Exchanges</button>
      </div>
    </section>

    <section class="settings-section">
      <h2>Sync Items</h2>
      <div class="sync-items-grid">
        <div v-for="item in syncItems" :key="item.id" class="sync-card">
          <div class="sync-header">
            <h3>{{ item.name }}</h3>
            <button class="delete-btn" @click="removeSyncItem(item.id)">×</button>
          </div>
          <div class="sync-details">
            <p><strong>Source:</strong> {{ item.source }}</p>
            <p><strong>Targets:</strong> {{ item.targets.join(', ') }}</p>
            <p><strong>Symbol:</strong> {{ item.symbol }}</p>
          </div>
        </div>

        <div class="sync-card add-card" @click="showAddModal = true">
          <div class="add-icon">+</div>
          <p>Add New Sync Item</p>
        </div>
      </div>
    </section>

    <!-- Simple Modal for adding sync item -->
    <div v-if="showAddModal" class="modal-overlay">
      <div class="modal">
        <h3>Add Sync Item</h3>
        <input v-model="newItem.name" placeholder="Item Name (e.g. Binance to OKX)" />
        <input v-model="newItem.symbol" placeholder="Symbol (e.g. BTC-USDT)" />
        <select v-model="newItem.source">
          <option value="binance">Binance</option>
        </select>
        <div class="checkbox-group">
          <label><input type="checkbox" value="okx" v-model="newItem.targets" /> OKX</label>
          <label><input type="checkbox" value="bybit" v-model="newItem.targets" /> Bybit</label>
        </div>
        <div class="modal-actions">
          <button @click="showAddModal = false">Cancel</button>
          <button class="primary-btn" @click="addSyncItem">Add</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import client from '../api/client'
import { v4 as uuidv4 } from 'uuid'

const config = ref({
  binance: { api_key: '', api_secret: '', testnet: false },
  okx: { api_key: '', api_secret: '', passphrase: '' },
  bybit: { api_key: '', api_secret: '' },
  sync: { symbol: '', position_ratio: 1.0, max_position: 1.0, stop_loss_ratio: 0.05 }
})

const syncItems = ref<any[]>([])
const showAddModal = ref(false)
const newItem = ref({
  name: '',
  source: 'binance',
  targets: [] as string[],
  symbol: 'BTC-USDT',
  enabled: true
})

const loadData = async () => {
  try {
    const [configRes, itemsRes] = await Promise.all([
      client.get('/config'),
      client.get('/sync-items')
    ])
    config.value = configRes.data
    syncItems.value = itemsRes.data || []
  } catch (e) {
    console.error('Failed to load data', e)
  }
}

const saveConfig = async () => {
  try {
    await client.put('/config', config.value)
    alert('Config saved')
  } catch (e) {
    alert('Failed to save config')
  }
}

const addSyncItem = async () => {
  try {
    const item = { ...newItem.value, id: uuidv4() }
    await client.post('/sync-items', item)
    syncItems.value.push(item)
    showAddModal.value = false
    newItem.value = { name: '', source: 'binance', targets: [], symbol: 'BTC-USDT', enabled: true }
  } catch (e) {
    alert('Failed to add item')
  }
}

const removeSyncItem = async (id: string) => {
  if (!confirm('Are you sure?')) return
  try {
    await client.delete(`/sync-items/${id}`)
    syncItems.value = syncItems.value.filter(i => i.id !== id)
  } catch (e) {
    alert('Failed to delete')
  }
}

onMounted(loadData)
</script>

<style scoped>
.settings { padding: 20px; }
.settings-section { margin-bottom: 40px; }
.settings-card { background: #242424; padding: 20px; border-radius: 8px; }
.exchange-group { margin-bottom: 30px; border-bottom: 1px solid #333; padding-bottom: 20px; }
.setting-item { margin-bottom: 15px; }
label { display: block; margin-bottom: 5px; color: #aaa; }
input[type="password"], input[type="text"], select { width: 100%; padding: 8px; background: #1a1a1a; border: 1px solid #444; color: white; border-radius: 4px; }

.sync-items-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(250px, 1fr)); gap: 20px; }
.sync-card { background: #242424; padding: 15px; border-radius: 8px; border: 1px solid #444; }
.sync-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 10px; }
.delete-btn { background: none; border: none; color: #ff4d4f; font-size: 1.5rem; cursor: pointer; padding: 0 5px; }
.add-card { border: 2px dashed #444; display: flex; flex-direction: column; align-items: center; justify-content: center; cursor: pointer; color: #888; min-height: 120px; }
.add-card:hover { border-color: #646cff; color: #646cff; }
.add-icon { font-size: 2rem; margin-bottom: 5px; }

.primary-btn { background: #646cff; color: white; border: none; padding: 10px 20px; border-radius: 4px; cursor: pointer; }
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.8); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal { background: #242424; padding: 30px; border-radius: 8px; width: 400px; }
.modal input, .modal select { margin-bottom: 15px; }
.checkbox-group { margin-bottom: 20px; display: flex; gap: 15px; }
.modal-actions { display: flex; justify-content: flex-end; gap: 10px; }
</style>
