<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import Header from '../components/Header.svelte';
  import { containerApi } from '../services/api';
  import { pageHeaderStore } from '../stores/pageHeaderStore';
  import type { Container } from '../types';
  
  let containers: Container[] = [];
  let loading = true;
  let error = '';
  let refreshing = false;
  let displayMode: 'compact' | 'standard' = 'standard';
  let actionError = '';
  let contentHeaderRef: HTMLElement;
  
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
  
  // Intersection Observer for scroll-based header
  let observer: IntersectionObserver | null = null;
  
  // Load display mode from localStorage
  onMount(() => {
    const savedMode = localStorage.getItem('displayMode');
    if (savedMode === 'compact' || savedMode === 'standard') {
      displayMode = savedMode;
    }
    
    // Set up page header
    pageHeaderStore.setTitle('容器列表');
    pageHeaderStore.setShowDisplayModeToggle(true);
    pageHeaderStore.setShowRefreshButton(true);
    pageHeaderStore.setDisplayMode(displayMode);
    pageHeaderStore.setOnToggleDisplayMode(toggleDisplayMode);
    pageHeaderStore.setOnRefresh(handleRefresh);
    
    // Set up intersection observer to detect when content header scrolls out of view
    // Header height is approximately 68px, using 64px as margin to trigger slightly before fully hidden
    const HEADER_HEIGHT_OFFSET = 64;
    if (contentHeaderRef) {
      observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            // When content header is not visible (scrolled past), show in main header
            const isScrolled = !entry.isIntersecting;
            pageHeaderStore.setIsScrolled(isScrolled);
            pageHeaderStore.setContentHeaderVisible(entry.isIntersecting);
          });
        },
        { 
          threshold: 0,
          rootMargin: `-${HEADER_HEIGHT_OFFSET}px 0px 0px 0px`
        }
      );
      observer.observe(contentHeaderRef);
    }
    
    loadContainers();
  });
  
  onDestroy(() => {
    // Clean up page header when leaving the page
    pageHeaderStore.reset();
    
    // Clean up observer
    if (observer) {
      observer.disconnect();
    }
  });
  
  function toggleDisplayMode() {
    displayMode = displayMode === 'compact' ? 'standard' : 'compact';
    localStorage.setItem('displayMode', displayMode);
    pageHeaderStore.setDisplayMode(displayMode);
  }
  
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
      pageHeaderStore.setRefreshing(false);
    }
  }
  
  async function handleAction(containerId: string, action: 'start' | 'stop' | 'restart', isSelf: boolean) {
    // Prevent stop/restart on self container
    if (isSelf && (action === 'stop' || action === 'restart')) {
      actionError = '无法停止或重启运行本应用的容器';
      setTimeout(() => { actionError = ''; }, 3000);
      return;
    }
    
    try {
      actionError = '';
      await containerApi.controlContainer({ containerId, action });
      await loadContainers();
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '未知错误';
      if (errorMessage.includes('cannot stop or restart')) {
        actionError = '无法停止或重启运行本应用的容器';
      } else {
        actionError = `操作失败: ${action}`;
      }
      console.error('Container action failed:', err);
      setTimeout(() => { actionError = ''; }, 3000);
    }
  }
  
  async function handleRefresh() {
    refreshing = true;
    pageHeaderStore.setRefreshing(true);
    await loadContainers();
  }
</script>

<div class="home-container">
  <Header />
  
  <main class="main-content">
    <div class="content-header" bind:this={contentHeaderRef}>
      <h2>容器列表</h2>
      <div class="header-actions">
        <button 
          class="mode-toggle" 
          onclick={toggleDisplayMode} 
          title={displayMode === 'compact' ? '切换到标准模式' : '切换到紧凑模式'}
          aria-label={displayMode === 'compact' ? '切换到标准模式' : '切换到紧凑模式'}
        >
          {#if displayMode === 'compact'}
            <span class="mode-icon">📋</span>
            <span class="mode-text">标准</span>
          {:else}
            <span class="mode-icon">📑</span>
            <span class="mode-text">紧凑</span>
          {/if}
        </button>
        <button class="refresh-button" onclick={handleRefresh} disabled={refreshing}>
          <span class="refresh-icon" class:spinning={refreshing}>🔄</span>
          刷新
        </button>
      </div>
    </div>
    
    {#if error}
      <div class="error-banner">
        {error}
      </div>
    {/if}
    
    {#if actionError}
      <div class="error-banner action-error">
        {actionError}
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
      <div class="container-list" class:compact={displayMode === 'compact'}>
        {#each containers as container (container.id)}
          <div class="container-item" class:is-self={container.is_self}>
            {#if displayMode === 'compact'}
              <!-- Compact mode: single line -->
              <div class="container-compact">
                <span class="compact-status">
                  <span class="status-emoji">{stateEmojis[container.state] || '⚪'}</span>
                  {#if container.health && container.health !== 'none'}
                    <span class="health-emoji">{healthEmojis[container.health]}</span>
                  {/if}
                </span>
                <span class="compact-name" title={container.name}>{container.name}</span>
                {#if container.is_self}
                  <span class="self-badge">本应用</span>
                {/if}
                <span class="compact-image" title={container.image}>{container.image}</span>
                <span class="compact-state">{container.status}</span>
                <div class="compact-actions">
                  {#if container.state === 'running'}
                    <button 
                      class="action-btn-compact stop" 
                      onclick={() => handleAction(container.id, 'stop', container.is_self ?? false)}
                      disabled={container.is_self}
                      title={container.is_self ? '无法停止本应用容器' : '停止'}
                    >
                      ⏸️
                    </button>
                    <button 
                      class="action-btn-compact restart" 
                      onclick={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                      disabled={container.is_self}
                      title={container.is_self ? '无法重启本应用容器' : '重启'}
                    >
                      🔄
                    </button>
                  {:else if ['exited', 'created', 'dead'].includes(container.state)}
                    <button 
                      class="action-btn-compact start" 
                      onclick={() => handleAction(container.id, 'start', container.is_self ?? false)}
                      title="启动"
                    >
                      ▶️
                    </button>
                  {:else}
                    <button 
                      class="action-btn-compact restart" 
                      onclick={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                      disabled={container.is_self}
                      title={container.is_self ? '无法重启本应用容器' : '重启'}
                    >
                      🔄
                    </button>
                  {/if}
                </div>
              </div>
            {:else}
              <!-- Standard mode: multi-line card -->
              <div class="container-info">
                <div class="container-name">
                  <span class="name-text">{container.name}</span>
                  {#if container.is_self}
                    <span class="self-badge">本应用</span>
                  {/if}
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
                  <button 
                    class="action-btn stop" 
                    onclick={() => handleAction(container.id, 'stop', container.is_self ?? false)}
                    disabled={container.is_self}
                    title={container.is_self ? '无法停止本应用容器' : ''}
                  >
                    ⏸️ 停止
                  </button>
                  <button 
                    class="action-btn restart" 
                    onclick={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                    disabled={container.is_self}
                    title={container.is_self ? '无法重启本应用容器' : ''}
                  >
                    🔄 重启
                  </button>
                {:else if ['exited', 'created', 'dead'].includes(container.state)}
                  <button 
                    class="action-btn start" 
                    onclick={() => handleAction(container.id, 'start', container.is_self ?? false)}
                  >
                    ▶️ 启动
                  </button>
                {:else}
                  <button 
                    class="action-btn restart" 
                    onclick={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                    disabled={container.is_self}
                    title={container.is_self ? '无法重启本应用容器' : ''}
                  >
                    🔄 重启
                  </button>
                {/if}
              </div>
            {/if}
          </div>
        {/each}
      </div>
    {/if}
  </main>
</div>

<style>
  .home-container {
    min-height: 100vh;
    background: var(--color-background, #f5f5f4);
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
    color: var(--color-text, #0a0a0a);
    margin: 0;
    font-family: var(--font-heading, "Playfair Display", serif);
  }
  
  .header-actions {
    display: flex;
    gap: 0.75rem;
  }
  
  .mode-toggle,
  .refresh-button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: var(--color-surface, #e7e5e4);
    border: 1px solid rgba(0, 0, 0, 0.1);
    padding: 0.5rem 1rem;
    border-radius: var(--radius, 0.25rem);
    cursor: pointer;
    font-size: 0.95rem;
    transition: all 0.2s;
    color: var(--color-text, #0a0a0a);
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .mode-toggle:hover,
  .refresh-button:hover:not(:disabled) {
    background: var(--color-background, #f5f5f4);
    border-color: var(--color-primary, #171717);
  }
  
  .refresh-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  .mode-icon,
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
    background: rgba(153, 27, 27, 0.1);
    border: 1px solid var(--color-error, #991b1b);
    color: var(--color-error, #991b1b);
    padding: 1rem;
    border-radius: var(--radius, 0.25rem);
    margin-bottom: 1.5rem;
  }
  
  .action-error {
    background: rgba(180, 83, 9, 0.1);
    border: 1px solid var(--color-warning, #b45309);
    color: var(--color-warning, #b45309);
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
  
  .empty-state {
    text-align: center;
    padding: 4rem 2rem;
    color: var(--color-muted, #78716c);
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
  
  .container-list.compact {
    gap: 0.5rem;
  }
  
  .container-item {
    background: var(--color-surface, #e7e5e4);
    border-radius: var(--radius, 0.25rem);
    padding: 1.25rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    transition: box-shadow 0.2s;
  }
  
  .container-item:hover {
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  }
  
  .container-item.is-self {
    border-left: 4px solid var(--color-accent, #991b1b);
  }
  
  .container-list.compact .container-item {
    padding: 0.5rem 1rem;
  }
  
  /* Standard mode styles */
  .container-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .container-name {
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--color-text, #0a0a0a);
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-family: var(--font-heading, "Playfair Display", serif);
  }
  
  .self-badge {
    background: var(--color-accent, #991b1b);
    color: white;
    font-size: 0.7rem;
    padding: 0.15rem 0.5rem;
    border-radius: var(--radius, 0.25rem);
    font-weight: 500;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .container-image {
    font-size: 0.9rem;
    color: var(--color-muted, #78716c);
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
    color: var(--color-text, #0a0a0a);
  }
  
  .container-actions {
    display: flex;
    gap: 0.5rem;
    flex-shrink: 0;
  }
  
  .action-btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: var(--radius, 0.25rem);
    cursor: pointer;
    font-size: 0.9rem;
    font-weight: 500;
    transition: all 0.2s;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .action-btn.start {
    background: var(--color-success, #15803d);
    color: white;
  }
  
  .action-btn.start:hover {
    background: #166534;
  }
  
  .action-btn.stop {
    background: var(--color-error, #991b1b);
    color: white;
  }
  
  .action-btn.stop:hover:not(:disabled) {
    background: #7f1d1d;
  }
  
  .action-btn.restart {
    background: var(--color-warning, #b45309);
    color: white;
  }
  
  .action-btn.restart:hover:not(:disabled) {
    background: #92400e;
  }
  
  .action-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  /* Compact mode styles */
  .container-compact {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    width: 100%;
    font-size: 0.9rem;
  }
  
  .compact-status {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    flex-shrink: 0;
    width: 40px;
  }
  
  .compact-name {
    font-weight: 600;
    color: var(--color-text, #0a0a0a);
    flex-shrink: 0;
    max-width: 200px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-family: var(--font-heading, "Playfair Display", serif);
  }
  
  .compact-image {
    flex: 1;
    color: var(--color-muted, #78716c);
    font-family: monospace;
    font-size: 0.85rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  
  .compact-state {
    flex-shrink: 0;
    color: var(--color-secondary, #525252);
    font-size: 0.85rem;
    min-width: 100px;
    text-align: right;
  }
  
  .compact-actions {
    display: flex;
    gap: 0.25rem;
    flex-shrink: 0;
  }
  
  .action-btn-compact {
    padding: 0.25rem 0.5rem;
    border: none;
    border-radius: var(--radius, 0.25rem);
    cursor: pointer;
    font-size: 0.85rem;
    transition: all 0.2s;
    background: transparent;
  }
  
  .action-btn-compact:hover:not(:disabled) {
    background: rgba(0, 0, 0, 0.1);
  }
  
  .action-btn-compact:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }
</style>
