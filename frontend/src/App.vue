<script setup lang="ts">
import { useAuthStore } from './stores/auth'
import { useRouter } from 'vue-router'
import { 
  NConfigProvider, 
  NMessageProvider, 
  NNotificationProvider, 
  NDialogProvider, 
  NLoadingBarProvider, 
  darkTheme 
} from 'naive-ui'

const auth = useAuthStore()
const router = useRouter()

const logout = () => {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <n-config-provider :theme="darkTheme">
    <n-message-provider>
      <n-notification-provider>
        <n-dialog-provider>
          <n-loading-bar-provider>
            <div class="app-layout">
              <nav class="nav-bar">
                <div class="nav-container">
                  <div class="nav-brand">
                    <div class="brand-logo">
                      <img src="https://cryptologos.cc/logos/bitcoin-btc-logo.svg" alt="BTC" />
                    </div>
                    <span>CryptoSyncBot</span>
                  </div>
                  <div class="nav-links" v-if="auth.isAuthenticated">
                    <router-link to="/" class="nav-link">
                      <span class="link-icon">📊</span> 控制面板
                    </router-link>
                    <router-link to="/settings" class="nav-link">
                      <span class="link-icon">⚙️</span> 设置
                    </router-link>
                    <button @click="logout" class="logout-btn">
                      退出登录
                    </button>
                  </div>
                </div>
              </nav>
              
              <main class="content-wrapper">
                <router-view v-slot="{ Component }">
                  <transition name="fade" mode="out-in">
                    <component :is="Component" />
                  </transition>
                </router-view>
              </main>
            </div>
          </n-loading-bar-provider>
        </n-dialog-provider>
      </n-notification-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<style>
.app-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.nav-bar {
  position: sticky;
  top: 0;
  z-index: 100;
  background: rgba(15, 23, 42, 0.8);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--border-color);
  height: 72px;
  display: flex;
  align-items: center;
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  padding: 0 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.nav-brand {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 1.25rem;
  font-weight: 800;
  color: var(--text-primary);
  letter-spacing: -0.025em;
}

.brand-logo {
  background: transparent;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.brand-logo img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.nav-link {
  padding: 0.5rem 1rem;
  border-radius: 8px;
  color: var(--text-secondary);
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: var(--transition);
}

.nav-link:hover {
  color: var(--text-primary);
  background: var(--surface-color);
}

.nav-link.router-link-active {
  color: var(--primary-color);
  background: rgba(56, 189, 248, 0.1);
}

.logout-btn {
  background: transparent;
  color: var(--danger);
  border: 1px solid rgba(239, 68, 68, 0.2);
  padding: 0.4rem 1rem;
  font-size: 0.875rem;
}

.logout-btn:hover {
  background: rgba(239, 68, 68, 0.1);
  border-color: var(--danger);
  color: var(--danger);
}

.content-wrapper {
  flex: 1;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  padding: 2rem 1.5rem;
}

/* Page Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@media (max-width: 640px) {
  .nav-container {
    padding: 0 1rem;
  }
  .nav-brand span {
    display: none;
  }
  .nav-link span.link-icon {
    font-size: 1.25rem;
  }
  .nav-link {
    padding: 0.5rem;
  }
}
</style>
