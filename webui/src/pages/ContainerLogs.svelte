<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { push } from 'svelte-spa-router';
  import Header from '../components/Header.svelte';
  import { get } from 'svelte/store';
  import { authStore } from '../stores/authStore';
  
  export let params: { id?: string } = {};
  
  let containerId = params.id || '';
  let logs: string[] = [];
  let error = '';
  let connecting = true;
  let ws: WebSocket | null = null;
  let autoScroll = true;
  let logContainer: HTMLElement;
  
  onMount(() => {
    if (!containerId) {
      error = '未指定容器 ID';
      connecting = false;
      return;
    }
    
    connectWebSocket();
  });
  
  onDestroy(() => {
    if (ws) {
      ws.close();
    }
  });
  
  function connectWebSocket() {
    const authState = get(authStore);
    const token = authState.token;
    if (!token) {
      error = '未授权，请重新登录';
      connecting = false;
      return;
    }
    
    // Get WebSocket URL from current location
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    // Pass token as query parameter since WebSocket doesn't support custom headers in browsers
    const wsUrl = `${protocol}//${host}/api/containers/${containerId}/logs/stream?token=${encodeURIComponent(token)}`;
    
    try {
      // Create WebSocket connection with token in query parameter
      ws = new WebSocket(wsUrl);
      
      ws.onopen = () => {
        connecting = false;
        error = '';
      };
      
      ws.onmessage = (event) => {
        logs = [...logs, event.data];
        if (autoScroll && logContainer) {
          setTimeout(() => {
            logContainer.scrollTop = logContainer.scrollHeight;
          }, 0);
        }
      };
      
      ws.onerror = (event) => {
        error = 'WebSocket 连接错误';
        connecting = false;
        console.error('WebSocket error:', event);
      };
      
      ws.onclose = (event) => {
        connecting = false;
        if (event.code !== 1000) {
          error = `连接已关闭 (代码: ${event.code})`;
        }
      };
    } catch (err) {
      error = '无法建立 WebSocket 连接';
      connecting = false;
      console.error('WebSocket connection failed:', err);
    }
  }
  
  function handleBack() {
    push('/');
  }
  
  function clearLogs() {
    logs = [];
  }
  
  function toggleAutoScroll() {
    autoScroll = !autoScroll;
    if (autoScroll && logContainer) {
      logContainer.scrollTop = logContainer.scrollHeight;
    }
  }
  
  function handleScroll() {
    if (!logContainer) return;
    const isAtBottom = logContainer.scrollHeight - logContainer.scrollTop <= logContainer.clientHeight + 50;
    if (autoScroll && !isAtBottom) {
      autoScroll = false;
    }
  }
</script>

<div class="logs-container">
  <Header />
  
  <main class="main-content">
    <div class="logs-header">
      <div class="header-left">
        <button class="back-button" on:click={handleBack}>
          ← 返回
        </button>
        <h2>容器日志</h2>
        <span class="container-id">{containerId}</span>
      </div>
      <div class="header-actions">
        <button 
          class="action-button" 
          class:active={autoScroll}
          on:click={toggleAutoScroll}
          title={autoScroll ? '关闭自动滚动' : '开启自动滚动'}
        >
          {autoScroll ? '🔽 自动滚动' : '⏸️ 暂停滚动'}
        </button>
        <button class="action-button" on:click={clearLogs}>
          🗑️ 清空日志
        </button>
      </div>
    </div>
    
    {#if error}
      <div class="error-banner">
        {error}
      </div>
    {/if}
    
    {#if connecting}
      <div class="loading">
        <div class="spinner"></div>
        <p>正在连接...</p>
      </div>
    {:else}
      <div 
        class="logs-viewer" 
        bind:this={logContainer}
        on:scroll={handleScroll}
      >
        {#if logs.length === 0}
          <div class="empty-state">
            <p>暂无日志数据</p>
            <p class="hint">日志将从最近 30 分钟开始显示</p>
          </div>
        {:else}
          {#each logs as log, index (index)}
            <div class="log-line">
              {log}
            </div>
          {/each}
        {/if}
      </div>
    {/if}
  </main>
</div>

<style>
  .logs-container {
    min-height: 100vh;
    background: var(--color-background, #f5f5f4);
    display: flex;
    flex-direction: column;
  }
  
  .main-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    max-width: 1400px;
    width: 100%;
    margin: 0 auto;
    padding: 2rem;
    box-sizing: border-box;
  }
  
  .logs-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    gap: 1rem;
    flex-wrap: wrap;
  }
  
  .header-left {
    display: flex;
    align-items: center;
    gap: 1rem;
    flex: 1;
    min-width: 0;
  }
  
  .back-button {
    padding: 0.5rem 1rem;
    background: var(--color-surface, #e7e5e4);
    border: 1px solid rgba(0, 0, 0, 0.1);
    border-radius: var(--radius, 0.25rem);
    cursor: pointer;
    font-size: 0.9rem;
    transition: all 0.2s;
    color: var(--color-text, #0a0a0a);
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .back-button:hover {
    background: var(--color-background, #f5f5f4);
    border-color: var(--color-primary, #171717);
  }
  
  .logs-header h2 {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--color-text, #0a0a0a);
    margin: 0;
    font-family: var(--font-heading, "Playfair Display", serif);
  }
  
  .container-id {
    font-family: monospace;
    font-size: 0.9rem;
    color: var(--color-muted, #78716c);
    padding: 0.25rem 0.5rem;
    background: var(--color-surface, #e7e5e4);
    border-radius: var(--radius, 0.25rem);
  }
  
  .header-actions {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
  }
  
  .action-button {
    padding: 0.5rem 1rem;
    background: var(--color-surface, #e7e5e4);
    border: 1px solid rgba(0, 0, 0, 0.1);
    border-radius: var(--radius, 0.25rem);
    cursor: pointer;
    font-size: 0.9rem;
    transition: all 0.2s;
    color: var(--color-text, #0a0a0a);
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .action-button:hover {
    background: var(--color-background, #f5f5f4);
    border-color: var(--color-primary, #171717);
  }
  
  .action-button.active {
    background: var(--color-primary, #171717);
    color: var(--color-background, #f5f5f4);
    border-color: var(--color-primary, #171717);
  }
  
  .error-banner {
    background: rgba(153, 27, 27, 0.1);
    border: 1px solid var(--color-error, #991b1b);
    color: var(--color-error, #991b1b);
    padding: 1rem;
    border-radius: var(--radius, 0.25rem);
    margin-bottom: 1.5rem;
  }
  
  .loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 4rem 2rem;
    color: var(--color-muted, #78716c);
  }
  
  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid var(--color-surface, #e7e5e4);
    border-top: 4px solid var(--color-primary, #171717);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 1rem;
  }
  
  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }
  
  .logs-viewer {
    flex: 1;
    background: #1e1e1e;
    color: #d4d4d4;
    padding: 1rem;
    border-radius: var(--radius, 0.25rem);
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
    font-size: 0.85rem;
    line-height: 1.5;
    overflow-y: auto;
    overflow-x: auto;
    max-height: calc(100vh - 250px);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
  
  .log-line {
    white-space: pre-wrap;
    word-break: break-all;
    padding: 0.125rem 0;
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  }
  
  .log-line:last-child {
    border-bottom: none;
  }
  
  .empty-state {
    text-align: center;
    padding: 4rem 2rem;
    color: #888;
  }
  
  .empty-state p {
    margin: 0.5rem 0;
  }
  
  .empty-state .hint {
    font-size: 0.9rem;
    color: #666;
  }
  
  /* Mobile responsive */
  @media (max-width: 768px) {
    .main-content {
      padding: 1rem;
    }
    
    .logs-header {
      flex-direction: column;
      align-items: stretch;
    }
    
    .header-left {
      flex-wrap: wrap;
    }
    
    .logs-header h2 {
      font-size: 1.25rem;
      width: 100%;
    }
    
    .container-id {
      font-size: 0.8rem;
    }
    
    .header-actions {
      width: 100%;
      justify-content: stretch;
    }
    
    .action-button {
      flex: 1;
      font-size: 0.85rem;
    }
    
    .logs-viewer {
      font-size: 0.75rem;
      max-height: calc(100vh - 300px);
    }
  }
</style>
