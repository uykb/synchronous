<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useTradingStore } from '../stores/trading'

const store = useTradingStore()

let interval: any = null

onMounted(() => {
  // Simulate incoming signals if the bot is running
  interval = setInterval(() => {
    if (store.isRunning) {
      const symbols = ['BTC/USDT', 'ETH/USDT', 'SOL/USDT', 'BNB/USDT']
      const sides = ['BUY', 'SELL']
      const signal = {
        time: new Date().toLocaleTimeString(),
        symbol: symbols[Math.floor(Math.random() * symbols.length)],
        side: sides[Math.floor(Math.random() * sides.length)],
        price: (Math.random() * 50000 + 1000).toFixed(2)
      }
      store.addSignal(signal)
    }
  }, 5000)
})

onUnmounted(() => {
  if (interval) clearInterval(interval)
})
</script>

<template>
  <div class="dashboard">
    <header class="dashboard-header">
      <div>
        <h1>Dashboard</h1>
        <p class="subtitle">Monitor your automated trading signals in real-time</p>
      </div>
      <div class="header-actions">
        <button 
          @click="store.toggleBot" 
          :class="['bot-toggle-btn', { 'is-running': store.isRunning }]"
        >
          <span class="btn-icon">{{ store.isRunning ? '⏹' : '▶' }}</span>
          {{ store.isRunning ? 'Stop Bot' : 'Start Bot' }}
        </button>
      </div>
    </header>
    
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-header">
          <span class="stat-icon status">📡</span>
          <h3>Bot Status</h3>
        </div>
        <div class="stat-value">
          <span :class="['status-badge', store.isRunning ? 'status-active' : 'status-inactive']">
            {{ store.isRunning ? 'Running' : 'Offline' }}
          </span>
        </div>
        <div class="stat-footer">
          {{ store.isRunning ? 'Bot is actively listening' : 'System is currently idle' }}
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-header">
          <span class="stat-icon time">🕒</span>
          <h3>Last Update</h3>
        </div>
        <div class="stat-value">{{ store.lastUpdate || 'Never' }}</div>
        <div class="stat-footer">Latest signal timestamp</div>
      </div>
      
      <div class="stat-card">
        <div class="stat-header">
          <span class="stat-icon signals">📊</span>
          <h3>Signals Today</h3>
        </div>
        <div class="stat-value">{{ store.signals.length }}</div>
        <div class="stat-footer">Total signals received</div>
      </div>
    </div>

    <div class="signals-section card">
      <div class="section-header">
        <h2>Recent Activity</h2>
        <div class="badge">{{ store.signals.length }} Total</div>
      </div>
      
      <div class="table-container">
        <table v-if="store.signals.length">
          <thead>
            <tr>
              <th>Timestamp</th>
              <th>Pair</th>
              <th>Action</th>
              <th class="text-right">Price</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(sig, index) in store.signals" :key="index">
              <td class="timestamp">{{ sig.time }}</td>
              <td class="symbol">
                <span class="symbol-tag">{{ sig.symbol }}</span>
              </td>
              <td>
                <span :class="['side-badge', sig.side === 'BUY' ? 'side-buy' : 'side-sell']">
                  {{ sig.side }}
                </span>
              </td>
              <td class="price text-right">{{ sig.price }}</td>
            </tr>
          </tbody>
        </table>
        <div v-else class="empty-state">
          <div class="empty-icon">📭</div>
          <p>Waiting for signals...</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard {
  animation: slideUp 0.4s ease-out;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 2rem;
  gap: 1rem;
}

.subtitle {
  color: var(--text-secondary);
  font-size: 0.95rem;
  margin-top: 0.25rem;
}

.bot-toggle-btn {
  background: var(--success);
  color: #fff;
  gap: 0.5rem;
}

.bot-toggle-btn.is-running {
  background: var(--danger);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: var(--surface-color);
  border-radius: var(--radius);
  border: 1px solid var(--border-color);
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  transition: var(--transition);
}

.stat-card:hover {
  transform: translateY(-4px);
  border-color: rgba(56, 189, 248, 0.2);
}

.stat-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.stat-header h3 {
  font-size: 0.875rem;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.stat-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  font-size: 1rem;
  background: rgba(255, 255, 255, 0.05);
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--text-primary);
}

.stat-footer {
  font-size: 0.8125rem;
  color: var(--text-muted);
}

.status-badge {
  display: inline-flex;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 600;
}

.status-active {
  background: rgba(16, 185, 129, 0.1);
  color: var(--success);
}

.status-inactive {
  background: rgba(239, 68, 68, 0.1);
  color: var(--danger);
}

.signals-section {
  padding: 0;
  overflow: hidden;
}

.section-header {
  padding: 1.5rem;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.badge {
  background: var(--surface-hover);
  padding: 0.25rem 0.625rem;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--primary-color);
}

.table-container {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th {
  background: rgba(255, 255, 255, 0.02);
  padding: 1rem 1.5rem;
  text-align: left;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--text-muted);
  border-bottom: 1px solid var(--border-color);
}

td {
  padding: 1rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
  font-size: 0.9375rem;
}

tr:hover td {
  background: rgba(255, 255, 255, 0.01);
}

.timestamp {
  color: var(--text-secondary);
  font-family: monospace;
}

.symbol-tag {
  background: rgba(56, 189, 248, 0.1);
  color: var(--primary-color);
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-weight: 600;
  font-size: 0.875rem;
}

.side-badge {
  display: inline-flex;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 700;
}

.side-buy {
  background: rgba(16, 185, 129, 0.1);
  color: var(--success);
}

.side-sell {
  background: rgba(239, 68, 68, 0.1);
  color: var(--danger);
}

.price {
  font-weight: 600;
  font-family: monospace;
}

.text-right {
  text-align: right;
}

.empty-state {
  padding: 4rem 2rem;
  text-align: center;
  color: var(--text-muted);
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

@media (max-width: 640px) {
  .dashboard-header {
    flex-direction: column;
    align-items: flex-start;
  }
  .bot-toggle-btn {
    width: 100%;
  }
}
</style>

