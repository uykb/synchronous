<template>
  <div class="login-container">
    <h1>Login</h1>
    <p>Enter the 6-digit code from your Authenticator app.</p>
    <div class="login-form">
      <input type="text" v-model="code" placeholder="6-digit TOTP Code" maxlength="6" @keyup.enter="handleLogin" />
      <button @click="handleLogin" :disabled="loading">Verify</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()
const code = ref('')
const loading = ref(false)

const handleLogin = async () => {
  if (code.value.length !== 6) return
  loading.value = true
  try {
    await auth.login(code.value)
    router.push('/')
  } catch (e) {
    alert('Invalid or expired code')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container { max-width: 400px; margin: 100px auto; padding: 20px; background: #242424; border-radius: 8px; text-align: center; }
input { width: 100%; padding: 12px; margin: 20px 0; background: #1a1a1a; border: 1px solid #444; color: white; text-align: center; font-size: 1.5rem; letter-spacing: 0.5rem; }
button { width: 100%; padding: 12px; background: #646cff; color: white; border: none; cursor: pointer; font-size: 1rem; }
</style>
