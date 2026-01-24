<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { NCard, NInput, NButton, NSpace, NText } from 'naive-ui'
import { useNotify } from '../composables/useNotify'

const router = useRouter()
const authStore = useAuthStore()
const { error } = useNotify()

const code = ref('')
const loading = ref(false)

async function handleLogin() {
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
  <div class="login-wrapper">
    <n-card class="login-card">
      <template #header>
        <div class="login-header">
          <img src="https://cryptologos.cc/logos/bitcoin-btc-logo.svg" alt="Logo" class="logo" />
        </div>
      </template>
      
      <n-space vertical size="large">
        <div class="input-group">
          <n-text depth="3" class="label">TOTP 验证码</n-text>
          <n-input 
            v-model:value="code"
            placeholder="000000"
            maxlength="6"
            size="large"
            class="code-input"
            @keyup.enter="handleLogin"
          />
          <n-text depth="3" class="helper">请输入您身份验证器应用中的 6 位代码。</n-text>
        </div>
        
        <n-button 
          type="primary" 
          block 
          size="large"
          :loading="loading"
          :disabled="code.length !== 6"
          @click="handleLogin"
        >
          {{ loading ? '验证中...' : '验证并进入' }}
        </n-button>
      </n-space>
    </n-card>
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
  width: 100%;
  max-width: 400px;
}

.login-header {
  text-align: center;
  padding: 1rem 0;
}

.logo {
  width: 64px;
  height: 64px;
}

.input-group {
  text-align: center;
}

.label {
  display: block;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  margin-bottom: 0.75rem;
  letter-spacing: 0.05em;
}

.code-input :deep(input) {
  text-align: center;
  font-size: 1.75rem !important;
  font-weight: 700;
  letter-spacing: 0.5rem;
}

.helper {
  display: block;
  margin-top: 0.75rem;
  font-size: 0.8125rem;
}
</style>
