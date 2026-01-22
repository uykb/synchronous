<template>
  <div class="setup-wrapper">
    <div class="setup-card">
      <div class="setup-header">
        <div class="logo">CSB</div>
        <h1>First Time Setup</h1>
        <p>Initialize your CryptoSyncBot instance</p>
      </div>

      <div v-if="!setupData" class="setup-body">
        <div class="form-group">
          <label>Administrator Password</label>
          <input 
            type="password" 
            v-model="password" 
            placeholder="Choose a strong password" 
            class="setup-input"
          />
          <p class="helper-text">This password will be used for dashboard access.</p>
        </div>
        <button @click="handleSetup" :disabled="loading || !password" class="primary-btn full-width">
          {{ loading ? 'Initializing...' : 'Initialize Bot' }}
        </button>
      </div>
      
      <div v-else class="totp-setup">
        <div class="step-badge">Security Step</div>
        <h3>Two-Factor Authentication</h3>
        <p class="instruction">Scan this QR code with Google Authenticator or Authy.</p>
        
        <div class="qr-container">
          <canvas ref="qrCanvas"></canvas>
        </div>
        
        <div class="secret-box">
          <label>Manual Entry Key</label>
          <code>{{ setupData.totp_secret }}</code>
        </div>

        <div class="verify-section">
          <div class="form-group">
            <label>Verification Code</label>
            <input 
              type="text" 
              v-model="code" 
              placeholder="000 000" 
              maxlength="6"
              class="verify-input"
            />
          </div>
          <button @click="handleVerify" :disabled="code.length !== 6" class="primary-btn full-width">
            Verify & Finish Setup
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
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  font-weight: 900;
  color: #000;
  margin: 0 auto 1rem;
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

