<template>
  <div class="dashboard">
    <h1>Crypto Sync Bot Dashboard</h1>
    
    <div class="stats-grid">
      <div class="stat-card">
        <h3>Status</h3>
        <p :class="{ 'status-active': store.isRunning, 'status-inactive': !store.isRunning }">
          {{ store.isRunning ? 'Running' : 'Stopped' }}
        </p>
        <button @click="store.toggleBot">
          {{ store.isRunning ? 'Stop Bot' : 'Start Bot' }}
        </button>
      </div>
      
      <div class="stat-card">
        <h3>Last Update</h3>
        <p>{{ store.lastUpdate || 'Never' }}</p>
      </div>
      
      <div class="stat-card">
        <h3>Signals Today</h3>
        <p>{{ store.signals.length }}</p>
      </div>
    </div>

    <div class="signals-section">
      <h2>Recent Signals</h2>
      <table v-if="store.signals.length">
        <thead>
          <tr>
            <th>Time</th>
            <th>Symbol</th>
            <th>Side</th>
            <th>Price</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(sig, index) in store.signals" :key="index">
            <td>{{ sig.time }}</td>
            <td>{{ sig.symbol }}</td>
            <td :class="sig.side === 'BUY' ? 'side-buy' : 'side-sell'">{{ sig.side }}</td>
            <td>{{ sig.price }}</td>
          </tr>
        </tbody>
      </table>
      <p v-else>No signals received yet.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTradingStore } from '../stores/trading'

const store = useTradingStore()

// Mock some initial data if empty for demonstration
if (store.signals.length === 0) {
  store.addSignal({ time: new Date().toLocaleTimeString(), symbol: 'BTCUSDT', side: 'BUY', price: '42500.50' })
  store.addSignal({ time: new Date().toLocaleTimeString(), symbol: 'ETHUSDT', side: 'SELL', price: '2250.20' })
}
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  background: #242424;
  padding: 20px;
  border-radius: 8px;
  border: 1px solid #444;
}

.status-active { color: #42b883; font-weight: bold; }
.status-inactive { color: #ff4d4f; font-weight: bold; }

.signals-section {
  background: #242424;
  padding: 20px;
  border-radius: 8px;
  border: 1px solid #444;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th, td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #444;
}

.side-buy { color: #42b883; }
.side-sell { color: #ff4d4f; }

button {
  margin-top: 10px;
  cursor: pointer;
}
</style>
