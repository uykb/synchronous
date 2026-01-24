<script setup lang="ts">
import { ref, nextTick, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import QRCode from 'qrcode'
import { 
  NCard, NInput, NButton, NSpace, NText, NSpin, NTag 
} from 'naive-ui'
import { useClipboard } from '../utils/clipboard'
import { useNotify } from '../composables/useNotify'

const router = useRouter()
const authStore = useAuthStore()
const { copied, copyToClipboard } = useClipboard()
const { error } = useNotify()

const setupData = ref<any>(null)
const code = ref('')
const loading = ref(false)
const qrCanvas = ref<HTMLCanvasElement | null>(null)

async function handleSetup() {
  loading.value = true
  try {
    const data = await authStore.setup('')
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
  } catch (e) {
    error('初始化设置失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  handleSetup()
})

async function handleVerify() {
  if (code.value.length !== 6) return
  
  loading.value = true
  try {
    await authStore.login(code.value)
    router.push('/')
  } catch (e) {
    error('验证码无效')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="setup-wrapper">
    <n-card class="setup-card">
      <template #header>
        <div class="setup-header">
          <img src="https://cryptologos.cc/logos/bitcoin-btc-logo.svg" alt="Logo" class="logo" />
          <h1>首次运行设置</h1>
          <n-text depth="3">初始化您的 CryptoSyncBot 实例</n-text>
        </div>
      </template>

      <n-spin :show="loading">
        <div v-if="!setupData" class="setup-body centered">
          <div v-if="loading">
            <n-text depth="3">正在生成安全密钥...</n-text>
          </div>
          <div v-else>
            <p>初始化失败，请重试</p>
            <n-button type="primary" class="mt-4" @click="handleSetup">
              生成 TOTP 密钥
            </n-button>
          </div>
        </div>
        
        <div v-else class="totp-setup">
          <n-tag type="info" :bordered="false" class="step-badge">安全步骤</n-tag>
          <h3>双重身份验证 (2FA)</h3>
          <p class="instruction">请使用 Google Authenticator 或 Authy 扫描此二维码。</p>
          
          <div class="qr-container">
            <canvas ref="qrCanvas"></canvas>
          </div>
          
          <div class="secret-box">
            <n-text depth="3" class="label">手动输入密钥</n-text>
            <div class="secret-value-wrapper">
              <code>{{ setupData.totp_secret }}</code>
              <n-button 
                size="small" 
                :type="copied ? 'success' : 'default'"
                @click="copyToClipboard(setupData.totp_secret)"
              >
                {{ copied ? '已复制' : '复制' }}
              </n-button>
            </div>
          </div>

          <div class="verify-section">
            <n-space vertical size="large">
              <div class="input-group">
                <n-text depth="3" class="label">验证码</n-text>
                <n-input 
                  v-model:value="code"
                  placeholder="000 000"
                  maxlength="6"
                  size="large"
                  class="verify-input"
                  @keyup.enter="handleVerify"
                />
              </div>
              <n-button 
                type="primary" 
                block 
                size="large"
                :loading="loading"
                :disabled="code.length !== 6"
                @click="handleVerify"
              >
                验证并完成设置
              </n-button>
            </n-space>
          </div>
        </div>
      </n-spin>
    </n-card>
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
  width: 100%;
  max-width: 450px;
}

.setup-header {
  text-align: center;
  padding: 1rem 0;
}

.logo {
  width: 64px;
  height: 64px;
  margin-bottom: 1rem;
}

.centered {
  text-align: center;
  padding: 2rem 0;
}

.mt-4 {
  margin-top: 1rem;
}

.totp-setup {
  text-align: center;
}

.step-badge {
  margin-bottom: 1rem;
}

.instruction {
  font-size: 0.875rem;
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
  text-align: left;
}

.label {
  display: block;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  margin-bottom: 0.5rem;
  letter-spacing: 0.05em;
}

.secret-value-wrapper {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  justify-content: space-between;
}

.secret-box code {
  color: var(--primary-color);
  font-weight: 700;
  font-family: monospace;
  font-size: 1.125rem;
  word-break: break-all;
  white-space: pre-wrap;
  display: block;
  flex: 1;
}

.verify-section {
  text-align: left;
  border-top: 1px solid var(--border-color);
  padding-top: 1.5rem;
}

.verify-input :deep(input) {
  text-align: center;
  font-size: 1.5rem !important;
  font-weight: 700;
  letter-spacing: 0.25rem;
}
</style>
