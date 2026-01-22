<template>
  <div class="setup-container">
    <h1>First Time Setup</h1>
    <div v-if="!setupData" class="setup-form">
      <p>Set an administrator password for your dashboard.</p>
      <input type="password" v-model="password" placeholder="New Password" />
      <button @click="handleSetup" :disabled="loading">Initialize Bot</button>
    </div>
    
    <div v-else class="totp-setup">
      <h3>Scan with Google Authenticator</h3>
      <div class="qr-code">
        <canvas ref="qrCanvas"></canvas>
      </div>
      <p>Secret: <code>{{ setupData.totp_secret }}</code></p>
      <div class="verify-step">
        <input type="text" v-model="code" placeholder="Enter 6-digit code" />
        <button @click="handleVerify">Verify & Login</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import QRCode from 'qrcode'

const auth = useAuthStore()
const router = useRouter()
const password = ref('')
const code = ref('')
const loading = ref(false)
const setupData = ref<any>(null)
const qrCanvas = ref<HTMLCanvasElement | null>(null)

const handleSetup = async () => {
  loading.value = true
  try {
    const data = await auth.setup(password.value)
    setupData.value = data
    await nextTick()
    if (qrCanvas.value && data.totp_url) {
      await QRCode.toCanvas(qrCanvas.value, data.totp_url)
    }
  } catch (e) {
    alert('Setup failed')
  } finally {
    loading.value = false
  }
}

const handleVerify = async () => {
  try {
    await auth.login(code.value)
    router.push('/')
  } catch (e) {
    alert('Invalid code')
  }
}
</script>

<style scoped>
.setup-container { max-width: 400px; margin: 50px auto; padding: 20px; background: #242424; border-radius: 8px; }
input { width: 100%; padding: 10px; margin-bottom: 20px; background: #1a1a1a; border: 1px solid #444; color: white; }
button { width: 100%; padding: 10px; background: #646cff; color: white; border: none; cursor: pointer; }
.qr-code { background: white; padding: 10px; margin: 20px 0; display: inline-block; }
code { background: #1a1a1a; padding: 2px 5px; }
</style>
