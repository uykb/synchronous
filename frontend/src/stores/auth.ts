import { defineStore } from 'pinia'
import client from '../api/client'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    isConfigured: false,
    loading: false
  }),
  getters: {
    isAuthenticated: (state) => !!state.token
  },
  actions: {
    async checkStatus() {
      try {
        const { data } = await client.get('/status')
        this.isConfigured = data.is_configured
      } catch (error) {
        console.error('Failed to check status', error)
      }
    },
    async setup(password: string) {
      const { data } = await client.post('/auth/setup', { password })
      return data // contains totp_url and totp_secret
    },
    async login(code: string) {
      const { data } = await client.post('/auth/verify', { code })
      this.token = data.token
      localStorage.setItem('token', data.token)
    },
    logout() {
      this.token = ''
      localStorage.removeItem('token')
    }
  }
})
