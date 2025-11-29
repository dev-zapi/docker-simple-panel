<script lang="ts">
  import { onMount } from 'svelte';
  import Header from '../components/Header.svelte';
  import { containerApi } from '../services/api';
  import type { Container } from '../types';
  
  let containers: Container[] = [];
  let loading = true;
  let error = '';
  let refreshing = false;
  
  const stateEmojis: Record<string, string> = {
    created: '🆕',
    running: '🟢',
    exited: '🔴',
    paused: '🟡',
    restarting: '🔄',
    removing: '🗑️',
    dead: '💀'
  };
  
  const healthEmojis: Record<string, string> = {
    healthy: '✅',
    unhealthy: '❌',
    starting: '🔄',
    none: ''
  };
  
  async function loadContainers() {
    try {
      error = '';
      containers = await containerApi.getContainers();
    } catch (err) {
      error = '获取容器列表失败';
      console.error('Failed to load containers:', err);
    } finally {
      loading = false;
      refreshing = false;
    }
  }
  
  async function handleAction(containerId: string, action: 'start' | 'stop' | 'restart') {
    try {
      await containerApi.controlContainer({ containerId, action });
      await loadContainers();
    } catch (err) {
      error = `操作失败: ${action}`;
      console.error('Container action failed:', err);
    }
  }
  
  async function handleRefresh() {
    refreshing = true;
    await loadContainers();
  }
  
  onMount(() => {
    loadContainers();
  });
</script>

<div class="home-container">
  <Header />
  
  <main class="main-content">
    <div class="content-header">
      <h2>容器列表</h2>
      <button class="refresh-button" on:click={handleRefresh} disabled={refreshing}>
        <span class="refresh-icon" class:spinning={refreshing}>🔄</span>
        刷新
      </button>
    </div>
    
    {#if error}
      <div class="error-banner">
        {error}
      </div>
    {/if}
    
    {#if loading}
      <div class="loading">
        <div class="spinner"></div>
        <p>加载中...</p>
      </div>
    {:else if containers.length === 0}
      <div class="empty-state">
        <div class="empty-icon">📦</div>
        <p>暂无容器</p>
      </div>
    {:else}
      <div class="container-list">
        {#each containers as container (container.id)}
          <div class="container-item">
            <div class="container-info">
              <div class="container-name">
                <span class="name-text">{container.name}</span>
              </div>
              <div class="container-image">{container.image}</div>
              <div class="container-meta">
                <span class="status">
                  <span class="status-emoji">{stateEmojis[container.state] || '⚪'}</span>
                  {container.status}
                </span>
                {#if container.health && container.health !== 'none'}
                  <span class="health">
                    <span class="health-emoji">{healthEmojis[container.health]}</span>
                    {container.health}
                  </span>
                {/if}
              </div>
            </div>
            
            <div class="container-actions">
              {#if container.state === 'running'}
                <button class="action-btn stop" on:click={() => handleAction(container.id, 'stop')}>
                  ⏸️ 停止
                </button>
                <button class="action-btn restart" on:click={() => handleAction(container.id, 'restart')}>
                  🔄 重启
                </button>
              {:else if ['exited', 'created', 'dead'].includes(container.state)}
                <button class="action-btn start" on:click={() => handleAction(container.id, 'start')}>
                  ▶️ 启动
                </button>
              {:else}
                <button class="action-btn restart" on:click={() => handleAction(container.id, 'restart')}>
                  🔄 重启
                </button>
              {/if}
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </main>
</div>

<style>
  .home-container {
    min-height: 100vh;
    background: #f5f5f5;
  }
  
  .main-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }
  
  .content-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
  }
  
  .content-header h2 {
    font-size: 1.75rem;
    font-weight: 700;
    color: #333;
    margin: 0;
  }
  
  .refresh-button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: white;
    border: 2px solid #e0e0e0;
    padding: 0.5rem 1rem;
    border-radius: 8px;
    cursor: pointer;
    font-size: 0.95rem;
    transition: all 0.2s;
  }
  
  .refresh-button:hover:not(:disabled) {
    border-color: #667eea;
    color: #667eea;
  }
  
  .refresh-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  .refresh-icon {
    display: inline-block;
    transition: transform 0.3s;
  }
  
  .refresh-icon.spinning {
    animation: spin 1s linear infinite;
  }
  
  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }
  
  .error-banner {
    background: #fee;
    border: 1px solid #fcc;
    color: #c33;
    padding: 1rem;
    border-radius: 8px;
    margin-bottom: 1.5rem;
  }
  
  .loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 4rem 2rem;
    color: #666;
  }
  
  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid #f3f3f3;
    border-top: 4px solid #667eea;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 1rem;
  }
  
  .empty-state {
    text-align: center;
    padding: 4rem 2rem;
    color: #999;
  }
  
  .empty-icon {
    font-size: 4rem;
    margin-bottom: 1rem;
  }
  
  .container-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .container-item {
    background: white;
    border-radius: 12px;
    padding: 1.25rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    transition: box-shadow 0.2s;
  }
  
  .container-item:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
  
  .container-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .container-name {
    font-size: 1.1rem;
    font-weight: 600;
    color: #333;
  }
  
  .container-image {
    font-size: 0.9rem;
    color: #666;
    font-family: monospace;
  }
  
  .container-meta {
    display: flex;
    gap: 1rem;
    font-size: 0.9rem;
  }
  
  .status, .health {
    display: flex;
    align-items: center;
    gap: 0.25rem;
  }
  
  .status {
    font-weight: 500;
    text-transform: capitalize;
  }
  
  .container-actions {
    display: flex;
    gap: 0.5rem;
    flex-shrink: 0;
  }
  
  .action-btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-size: 0.9rem;
    font-weight: 500;
    transition: all 0.2s;
  }
  
  .action-btn.start {
    background: #27ae60;
    color: white;
  }
  
  .action-btn.start:hover {
    background: #229954;
  }
  
  .action-btn.stop {
    background: #e74c3c;
    color: white;
  }
  
  .action-btn.stop:hover {
    background: #c0392b;
  }
  
  .action-btn.restart {
    background: #f39c12;
    color: white;
  }
  
  .action-btn.restart:hover {
    background: #e67e22;
  }
</style>
