import { defineStore } from 'pinia'

export const useTradingStore = defineStore('trading', {
  state: () => ({
    signals: [] as any[],
    isRunning: false,
    lastUpdate: null as string | null,
  }),
  actions: {
    addSignal(signal: any) {
      this.signals.unshift(signal)
      if (this.signals.length > 50) {
        this.signals.pop()
      }
      this.lastUpdate = new Date().toLocaleTimeString()
    },
    toggleBot() {
      this.isRunning = !this.isRunning
    }
  }
})
