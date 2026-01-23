<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const code = ref('')
const loading = ref(false)

async function handleLogin() {
  if (code.value.length !== 6) return
  
  loading.value = true
  try {
    await authStore.login(code.value)
    router.push('/')
  } catch (error) {
    console.error('Login failed', error)
    alert('Invalid verification code')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-wrapper">
    <div class="login-card">
      <div class="login-header">
        <div class="logo">
          <img src="https://cryptologos.cc/logos/bitcoin-btc-logo.svg" alt="BTC" />
        </div>
        <h1>Secure Access</h1>
        <p>Verification required to continue</p>
      </div>
      
      <div class="login-body">
        <div class="input-group">
          <label for="totp">TOTP Verification Code</label>
          <input 
            id="totp"
            type="text" 
            v-model="code" 
            placeholder="000000" 
            maxlength="6" 
            @keyup.enter="handleLogin"
            autocomplete="one-time-code"
            inputmode="numeric"
          />
          <p class="helper-text">Enter the 6-digit code from your Authenticator app.</p>
        </div>
        
        <button @click="handleLogin" :disabled="loading || code.length !== 6" class="login-btn">
          <span v-if="loading" class="spinner"></span>
          {{ loading ? 'Verifying...' : 'Verify & Enter' }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-wrapper {
  min-height: calc(100vh - 200px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.login-card {
  background: var(--surface-color);
  width: 100%;
  max-width: 400px;
  padding: 2.5rem;
  border-radius: var(--radius);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-lg);
}

.login-header {
  text-align: center;
  margin-bottom: 2rem;
}

.logo {
  background: transparent;
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1.5rem;
}

.logo img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.login-header h1 {
  font-size: 1.5rem;
  margin-bottom: 0.5rem;
}

.login-header p {
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.input-group {
  margin-bottom: 2rem;
}

.input-group label {
  display: block;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--text-muted);
  margin-bottom: 0.75rem;
  letter-spacing: 0.05em;
}

input {
  width: 100%;
  padding: 1rem;
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  color: var(--text-primary);
  text-align: center;
  font-size: 1.75rem;
  font-weight: 700;
  letter-spacing: 0.5rem;
  transition: var(--transition);
}

input:focus {
  outline: none;
  border-color: var(--primary-color);
  background: rgba(15, 23, 42, 0.8);
  box-shadow: 0 0 0 4px rgba(56, 189, 248, 0.1);
}

.helper-text {
  margin-top: 0.75rem;
  font-size: 0.8125rem;
  color: var(--text-muted);
  text-align: center;
}

.login-btn {
  width: 100%;
  padding: 1rem;
  font-size: 1rem;
  font-weight: 700;
}

.spinner {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(0,0,0,0.2);
  border-top-color: #000;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-right: 0.5rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>

