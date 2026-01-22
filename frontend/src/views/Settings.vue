<template>
  <div class="settings">
    <header class="settings-header">
      <h1>Settings</h1>
      <p class="subtitle">Configure your exchange credentials and synchronization rules</p>
    </header>
    
    <div class="settings-layout">
      <section class="settings-section">
        <div class="section-title">
          <span class="icon">🔑</span>
          <h2>Exchanges Configuration</h2>
        </div>
        
        <div class="settings-grid">
          <div class="exchange-card">
            <div class="exchange-header">
              <img src="https://cryptologos.cc/logos/binance-coin-bnb-logo.png" alt="Binance" class="exchange-logo" />
              <h3>Binance</h3>
            </div>
            <div class="exchange-body">
              <div class="form-group">
                <label>API Key</label>
                <input type="password" v-model="config.binance.api_key" placeholder="Enter API Key" />
              </div>
              <div class="form-group">
                <label>Secret Key</label>
                <input type="password" v-model="config.binance.api_secret" placeholder="Enter Secret Key" />
              </div>
              <label class="checkbox-label">
                <input type="checkbox" v-model="config.binance.testnet" />
                <span class="checkbox-text">Enable Testnet Mode</span>
              </label>
            </div>
          </div>

          <div class="exchange-card">
            <div class="exchange-header">
              <img src="https://cryptologos.cc/logos/okx-okb-logo.png" alt="OKX" class="exchange-logo" />
              <h3>OKX</h3>
            </div>
            <div class="exchange-body">
              <div class="form-group">
                <label>API Key</label>
                <input type="password" v-model="config.okx.api_key" placeholder="Enter API Key" />
              </div>
              <div class="form-group">
                <label>Secret Key</label>
                <input type="password" v-model="config.okx.api_secret" placeholder="Enter Secret Key" />
              </div>
              <div class="form-group">
                <label>Passphrase</label>
                <input type="password" v-model="config.okx.passphrase" placeholder="Enter Passphrase" />
              </div>
            </div>
          </div>
        </div>

        <div class="actions-footer">
          <button class="primary-btn" @click="saveConfig">
            <span>💾</span> Save Changes
          </button>
          <button class="secondary-btn" @click="restartBot">
            <span>🔄</span> Restart Bot
          </button>
        </div>
      </section>

      <section class="settings-section">
        <div class="section-title">
          <span class="icon">🔄</span>
          <h2>Sync Rules</h2>
        </div>
        
        <div class="sync-items-grid">
          <div v-for="item in syncItems" :key="item.id" class="sync-card">
            <div class="sync-header">
              <div class="sync-title">
                <h3>{{ item.name }}</h3>
                <span class="symbol-tag">{{ item.symbol }}</span>
              </div>
              <button class="delete-btn" @click="removeSyncItem(item.id)" title="Remove Rule">×</button>
            </div>
            <div class="sync-details">
              <div class="detail-row">
                <span class="label">Source</span>
                <span class="value">{{ item.source }}</span>
              </div>
              <div class="detail-row">
                <span class="label">Targets</span>
                <div class="targets-list">
                  <span v-for="t in item.targets" :key="t" class="target-tag">{{ t }}</span>
                </div>
              </div>
            </div>
          </div>

          <div class="add-sync-card" @click="showAddModal = true">
            <div class="add-plus">+</div>
            <p>Add New Sync Rule</p>
          </div>
        </div>
      </section>
    </div>

    <!-- Enhanced Modal -->
    <transition name="modal-fade">
      <div v-if="showAddModal" class="modal-overlay" @click.self="showAddModal = false">
        <div class="modal-content">
          <div class="modal-header">
            <h3>New Sync Rule</h3>
            <button class="close-btn" @click="showAddModal = false">×</button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label>Rule Name</label>
              <input v-model="newItem.name" placeholder="e.g. BTC Arbitrage" />
            </div>
            <div class="form-group">
              <label>Trading Pair</label>
              <input v-model="newItem.symbol" placeholder="e.g. BTC-USDT" />
            </div>
            <div class="form-group">
              <label>Source Exchange</label>
              <select v-model="newItem.source">
                <option value="binance">Binance</option>
                <option value="okx">OKX</option>
              </select>
            </div>
            <div class="form-group">
              <label>Target Exchanges</label>
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
            <button class="secondary-btn" @click="showAddModal = false">Cancel</button>
            <button class="primary-btn" @click="addSyncItem">Create Rule</button>
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
  align-items: center;
  gap: 1rem;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--border-color);
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

