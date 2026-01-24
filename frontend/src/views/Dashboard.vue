<script setup lang="ts">
import { ref, h } from 'vue'
import { 
  NTag, NDataTable, DataTableColumns, NCard, NStatistic, 
  NIcon, NButton, NText, NSkeleton, NEmpty, NGrid, NGi 
} from 'naive-ui'
import { 
  RadioOutline, TimeOutline, StatsChartOutline, 
  PlayOutline, StopOutline 
} from '@vicons/ionicons5'
import { useTradingStore } from '../stores/trading'

const store = useTradingStore()
const loading = ref(false)

const columns: DataTableColumns = [
  { 
    title: '时间戳', 
    key: 'time',
    render: (row) => h('span', { class: 'timestamp' }, row.time as string)
  },
  { 
    title: '交易对', 
    key: 'symbol',
    render: (row) => h(NTag, { type: 'info', size: 'small' }, { default: () => row.symbol })
  },
  { 
    title: '操作', 
    key: 'side',
    render: (row) => h(NTag, { 
      type: row.side === 'BUY' ? 'success' : 'error',
      size: 'small' 
    }, { default: () => row.side === 'BUY' ? '买入' : '卖出' })
  },
  { 
    title: '价格', 
    key: 'price', 
    align: 'right',
    render: (row) => h('span', { class: 'price' }, row.price as string)
  }
]
</script>

<template>
  <div class="dashboard">
    <header class="dashboard-header">
      <div>
        <h1>控制面板</h1>
        <p class="subtitle">实时监控您的自动化交易信号</p>
      </div>
      <div class="header-actions">
        <n-button 
          :type="store.isRunning ? 'error' : 'success'"
          @click="store.toggleBot"
          size="large"
        >
          <template #icon>
            <n-icon>
              <StopOutline v-if="store.isRunning" />
              <PlayOutline v-else />
            </n-icon>
          </template>
          {{ store.isRunning ? '停止机器人' : '启动机器人' }}
        </n-button>
      </div>
    </header>
    
    <n-grid :x-gap="24" :y-gap="24" cols="1 s:2 m:3" responsive="screen" class="stats-grid">
      <n-gi>
        <n-card hoverable>
          <n-statistic label="机器人状态">
            <template #prefix>
              <n-icon><RadioOutline /></n-icon>
            </template>
            <n-tag :type="store.isRunning ? 'success' : 'error'">
              {{ store.isRunning ? '运行中' : '已离线' }}
            </n-tag>
          </n-statistic>
          <template #footer>
            <n-text depth="3">
              {{ store.isRunning ? '机器人正在监听信号' : '系统当前处于空闲状态' }}
            </n-text>
          </template>
        </n-card>
      </n-gi>
      
      <n-gi>
        <n-card hoverable>
          <n-statistic label="最后更新">
            <template #prefix>
              <n-icon><TimeOutline /></n-icon>
            </template>
            <span class="stat-value">{{ store.lastUpdate || '从未' }}</span>
          </n-statistic>
          <template #footer>
            <n-text depth="3">最新信号时间戳</n-text>
          </template>
        </n-card>
      </n-gi>
      
      <n-gi>
        <n-card hoverable>
          <n-statistic label="今日信号">
            <template #prefix>
              <n-icon><StatsChartOutline /></n-icon>
            </template>
            <span class="stat-value">{{ store.signals.length }}</span>
          </n-statistic>
          <template #footer>
            <n-text depth="3">累计接收信号总数</n-text>
          </template>
        </n-card>
      </n-gi>
    </n-grid>

    <n-card title="最近活动" class="signals-section" :segmented="{ content: true }">
      <template #header-extra>
        <n-tag type="primary" size="small" round>
          共 {{ store.signals.length }} 条
        </n-tag>
      </template>

      <div class="table-container">
        <n-skeleton v-if="loading" text :repeat="5" />
        <template v-else>
          <n-data-table
            v-if="store.signals.length"
            :columns="columns"
            :data="store.signals"
            :bordered="false"
            :single-line="false"
          />
          <n-empty v-else description="等待信号中..." />
        </template>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.dashboard {
  animation: slideUp 0.4s ease-out;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 2rem;
  gap: 1rem;
}

.subtitle {
  color: var(--text-secondary);
  font-size: 0.95rem;
  margin-top: 0.25rem;
}

.stats-grid {
  margin-bottom: 2rem;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 700;
}

.timestamp {
  font-family: monospace;
  color: var(--text-secondary);
}

.price {
  font-weight: 600;
  font-family: monospace;
}

.table-container {
  min-height: 200px;
}

@media (max-width: 640px) {
  .dashboard-header {
    flex-direction: column;
    align-items: flex-start;
  }
  .header-actions {
    width: 100%;
  }
  .header-actions :deep(.n-button) {
    width: 100%;
  }
}
</style>
