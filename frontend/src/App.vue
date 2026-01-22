<script setup lang="ts">
import { useAuthStore } from './stores/auth'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()

const logout = () => {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <nav class="nav-bar">
    <div class="nav-brand">CryptoSyncBot</div>
    <div class="nav-links" v-if="auth.isAuthenticated">
      <router-link to="/">Dashboard</router-link>
      <router-link to="/settings">Settings</router-link>
      <a href="#" @click.prevent="logout">Logout</a>
    </div>
  </nav>
  
  <main class="content">
    <router-view></router-view>
  </main>
</template>

<style>
#app {
  max-width: 1280px;
  margin: 0 auto;
  font-family: Inter, system-ui, Avenir, Helvetica, Arial, sans-serif;
}

.nav-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background: #1a1a1a;
  border-bottom: 1px solid #333;
}

.nav-brand {
  font-size: 1.5rem;
  font-weight: bold;
  color: #646cff;
}

.nav-links a {
  margin-left: 20px;
  text-decoration: none;
  color: #888;
  transition: color 0.3s;
}

.nav-links a:hover, .nav-links a.router-link-active {
  color: #646cff;
}

.content {
  padding: 20px;
}
</style>
