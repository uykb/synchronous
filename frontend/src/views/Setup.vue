<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import QRCode from 'qrcode'

const router = useRouter()
const authStore = useAuthStore()

const password = ref('')
const setupData = ref<any>(null)
const code = ref('')
const loading = ref(false)
const qrCanvas = ref<HTMLCanvasElement | null>(null)

async function handleSetup() {
  if (!password.value) return
  loading.value = true
  try {
    const data = await authStore.setup(password.value)
    setupData.value = data
    await nextTick()
    if (qrCanvas.value && data.totp_url) {
      await QRCode.toCanvas(qrCanvas.value, data.totp_url, {
        width: 200,
        margin: 2,
        color: {
          dark: '#0f172a',
          light: '#ffffff'
        }
      })
    }
  } catch (error) {
    console.error('Setup failed', error)
    alert('初始化设置失败')
  } finally {
    loading.value = false
  }
}

async function handleVerify() {
  if (code.value.length !== 6) return
  
  loading.value = true
  try {
    await authStore.login(code.value)
    router.push('/')
  } catch (error) {
    console.error('Verification failed', error)
    alert('验证码无效')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="setup-wrapper">
    <div class="setup-card">
      <div class="setup-header">
        <div class="logo">
          <img src="https://cryptologos.cc/logos/bitcoin-btc-logo.svg" alt="BTC" />
        </div>
        <h1>首次运行设置</h1>
        <p>初始化您的 CryptoSyncBot 实例</p>
      </div>

      <div v-if="!setupData" class="setup-body">
        <div class="form-group">
          <label>管理员密码</label>
          <input 
            type="password" 
            v-model="password" 
            placeholder="请选择一个强密码" 
            class="setup-input"
          />
          <p class="helper-text">此密码将用于访问控制面板。</p>
        </div>
        <button @click="handleSetup" :disabled="loading || !password" class="primary-btn full-width">
          {{ loading ? '正在初始化...' : '初始化机器人' }}
        </button>
      </div>
      
      <div v-else class="totp-setup">
        <div class="step-badge">安全步骤</div>
        <h3>双重身份验证 (2FA)</h3>
        <p class="instruction">请使用 Google Authenticator 或 Authy 扫描此二维码。</p>
        
        <div class="qr-container">
          <canvas ref="qrCanvas"></canvas>
        </div>
        
        <div class="secret-box">
          <label>手动输入密钥</label>
          <code>{{ setupData.totp_secret }}</code>
        </div>

        <div class="verify-section">
          <div class="form-group">
            <label>验证码</label>
            <input 
              type="text" 
              v-model="code" 
              placeholder="000 000" 
              maxlength="6"
              class="verify-input"
            />
          </div>
          <button @click="handleVerify" :disabled="code.length !== 6" class="primary-btn full-width">
            验证并完成设置
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.setup-wrapper {
  min-height: calc(100vh - 200px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.5rem;
}

.setup-card {
  background: var(--surface-color);
  width: 100%;
  max-width: 450px;
  padding: 2.5rem;
  border-radius: var(--radius);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-lg);
}

.setup-header {
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

.setup-header h1 {
  font-size: 1.5rem;
  margin-bottom: 0.5rem;
}

.setup-header p {
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.form-group {
  margin-bottom: 1.5rem;
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

.setup-input, .verify-input {
  width: 100%;
  padding: 0.875rem 1rem;
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-primary);
  font-size: 1rem;
  transition: var(--transition);
}

.verify-input {
  text-align: center;
  font-size: 1.5rem;
  letter-spacing: 0.25rem;
  font-weight: 700;
}

.setup-input:focus, .verify-input:focus {
  outline: none;
  border-color: var(--primary-color);
}

.helper-text {
  margin-top: 0.5rem;
  font-size: 0.8125rem;
  color: var(--text-muted);
}

.full-width {
  width: 100%;
}

.totp-setup {
  text-align: center;
}

.step-badge {
  display: inline-block;
  background: rgba(129, 140, 248, 0.1);
  color: var(--secondary-color);
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 700;
  margin-bottom: 1rem;
}

.instruction {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-bottom: 1.5rem;
}

.qr-container {
  background: white;
  padding: 1rem;
  border-radius: 12px;
  display: inline-block;
  margin-bottom: 1.5rem;
}

.secret-box {
  background: rgba(0,0,0,0.2);
  padding: 1rem;
  border-radius: 8px;
  margin-bottom: 2rem;
  border: 1px solid var(--border-color);
}

.secret-box label {
  display: block;
  font-size: 0.625rem;
  text-transform: uppercase;
  color: var(--text-muted);
  margin-bottom: 0.5rem;
}

.secret-box code {
  color: var(--primary-color);
  font-weight: 700;
  font-family: monospace;
  font-size: 1.125rem;
}

.verify-section {
  text-align: left;
  border-top: 1px solid var(--border-color);
  padding-top: 1.5rem;
}
</style>

