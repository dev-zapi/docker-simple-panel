<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import Header from '../components/Header.svelte';
  import { containerApi } from '../services/api';
  import type { Container } from '../types';
  
  let containers: Container[] = [];
  let loading = true;
  let error = '';
  let refreshing = false;
  let displayMode: 'compact' | 'standard' = 'standard';
  let groupMode: 'none' | 'compose' | 'label' | 'status' | 'health' = 'none';
  let actionError = '';
  let collapsedGroups: Set<string> = new Set();
  let selectedLabelKey: string = '';
  let availableLabelKeys: string[] = [];
  
  // Scroll-based header state
  let isScrolled = false;
  let contentHeaderRef: HTMLElement;
  let observer: IntersectionObserver | null = null;
  
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
  
  // Load display mode and group mode from localStorage
  onMount(() => {
    const savedMode = localStorage.getItem('displayMode');
    if (savedMode === 'compact' || savedMode === 'standard') {
      displayMode = savedMode;
    }
    const savedGroupMode = localStorage.getItem('groupMode');
    if (savedGroupMode === 'none' || savedGroupMode === 'compose' || savedGroupMode === 'label' || savedGroupMode === 'status' || savedGroupMode === 'health') {
      groupMode = savedGroupMode;
    }
    const savedSelectedLabelKey = localStorage.getItem('selectedLabelKey');
    if (savedSelectedLabelKey) {
      selectedLabelKey = savedSelectedLabelKey;
    }
    const savedCollapsedGroups = localStorage.getItem('collapsedGroups');
    if (savedCollapsedGroups) {
      try {
        collapsedGroups = new Set(JSON.parse(savedCollapsedGroups));
      } catch (e) {
        // Invalid JSON in localStorage, reset to empty set
        console.warn('Failed to parse collapsedGroups from localStorage:', e);
        collapsedGroups = new Set();
      }
    }
    loadContainers();
    
    // Set up intersection observer to detect when content header scrolls out of view
    const HEADER_HEIGHT = 68; // Header height in pixels
    if (contentHeaderRef) {
      observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            isScrolled = !entry.isIntersecting;
          });
        },
        { 
          threshold: 0,
          rootMargin: `-${HEADER_HEIGHT}px 0px 0px 0px`
        }
      );
      observer.observe(contentHeaderRef);
    }
  });
  
  onDestroy(() => {
    if (observer) {
      observer.disconnect();
    }
  });
  
  function toggleDisplayMode() {
    displayMode = displayMode === 'compact' ? 'standard' : 'compact';
    localStorage.setItem('displayMode', displayMode);
  }
  
  function handleGroupModeChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    groupMode = target.value as 'none' | 'compose' | 'label' | 'status' | 'health';
    localStorage.setItem('groupMode', groupMode);
  }
  
  function toggleGroupCollapse(groupName: string) {
    const newCollapsedGroups = new Set(collapsedGroups);
    if (newCollapsedGroups.has(groupName)) {
      newCollapsedGroups.delete(groupName);
    } else {
      newCollapsedGroups.add(groupName);
    }
    collapsedGroups = newCollapsedGroups;
    localStorage.setItem('collapsedGroups', JSON.stringify(Array.from(collapsedGroups)));
  }
  
  // Group containers by compose project
  function groupContainersByCompose(containers: Container[]) {
    const grouped = new Map<string, Container[]>();
    const ungrouped: Container[] = [];
    
    for (const container of containers) {
      if (container.compose_project) {
        const existing = grouped.get(container.compose_project) || [];
        existing.push(container);
        grouped.set(container.compose_project, existing);
      } else {
        ungrouped.push(container);
      }
    }
    
    return { grouped, ungrouped };
  }
  
  // Group containers by selected label
  function groupContainersByLabel(containers: Container[], labelKey: string) {
    const grouped = new Map<string, Container[]>();
    const ungrouped: Container[] = [];
    
    for (const container of containers) {
      const labelValue = container.labels?.[labelKey];
      if (labelValue) {
        const existing = grouped.get(labelValue) || [];
        existing.push(container);
        grouped.set(labelValue, existing);
      } else {
        ungrouped.push(container);
      }
    }
    
    return { grouped, ungrouped };
  }
  
  // Group containers by state (status)
  function groupContainersByStatus(containers: Container[]) {
    const grouped = new Map<string, Container[]>();
    
    for (const container of containers) {
      const state = container.state;
      const existing = grouped.get(state) || [];
      existing.push(container);
      grouped.set(state, existing);
    }
    
    return { grouped, ungrouped: [] };
  }
  
  // Group containers by health
  function groupContainersByHealth(containers: Container[]) {
    const grouped = new Map<string, Container[]>();
    
    for (const container of containers) {
      const health = container.health || 'none';
      const existing = grouped.get(health) || [];
      existing.push(container);
      grouped.set(health, existing);
    }
    
    return { grouped, ungrouped: [] };
  }
  
  // Extract all unique label keys from containers
  function extractLabelKeys(containers: Container[]): string[] {
    const keysSet = new Set<string>();
    for (const container of containers) {
      if (container.labels) {
        for (const key of Object.keys(container.labels)) {
          keysSet.add(key);
        }
      }
    }
    return Array.from(keysSet).sort();
  }
  
  // Update available label keys when containers change
  $: {
    availableLabelKeys = extractLabelKeys(containers);
    // If selected label key is not available anymore, reset it
    if (selectedLabelKey && !availableLabelKeys.includes(selectedLabelKey)) {
      selectedLabelKey = availableLabelKeys[0] || '';
      localStorage.setItem('selectedLabelKey', selectedLabelKey);
    }
    // If no label key is selected and label grouping is active, select the first one
    if (groupMode === 'label' && !selectedLabelKey && availableLabelKeys.length > 0) {
      selectedLabelKey = availableLabelKeys[0];
      localStorage.setItem('selectedLabelKey', selectedLabelKey);
    }
  }
  
  function handleLabelKeyChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    selectedLabelKey = target.value;
    localStorage.setItem('selectedLabelKey', selectedLabelKey);
  }
  
  function scrollToGroup(groupName: string) {
    const element = document.getElementById(`group-${groupName}`);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
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
    await loadContainers();
  }
</script>

<div class="home-container" class:scrolled={isScrolled}>
  <Header />
  
  <!-- Floating header that appears when scrolled -->
  <div class="floating-header" class:visible={isScrolled}>
    <h2>容器列表</h2>
    <div class="header-actions">
      <button 
        class="mode-toggle" 
        on:click={toggleDisplayMode} 
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
      <select 
        class="group-mode-select" 
        value={groupMode} 
        on:change={handleGroupModeChange}
        aria-label="选择分组方式"
      >
        <option value="none">不分组</option>
        <option value="compose">按 Compose 分组</option>
        <option value="label">按标签分组</option>
      </select>
      {#if groupMode === 'label' && availableLabelKeys.length > 0}
        <select 
          class="label-key-select" 
          value={selectedLabelKey} 
          on:change={handleLabelKeyChange}
          aria-label="选择标签"
        >
          {#each availableLabelKeys as labelKey}
            <option value={labelKey}>{labelKey}</option>
          {/each}
        </select>
      {/if}
      <button class="refresh-button" on:click={handleRefresh} disabled={refreshing}>
        <span class="refresh-icon" class:spinning={refreshing}>🔄</span>
        刷新
      </button>
    </div>
  </div>
  
  <main class="main-content">
    <div class="content-header" bind:this={contentHeaderRef}>
      <h2>容器列表</h2>
      <div class="header-actions">
        <button 
          class="mode-toggle" 
          on:click={toggleDisplayMode} 
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
        <select 
          class="group-mode-select" 
          value={groupMode} 
          on:change={handleGroupModeChange}
          aria-label="选择分组方式"
        >
          <option value="none">不分组</option>
          <option value="compose">按 Compose 分组</option>
          <option value="label">按标签分组</option>
          <option value="status">按状态分组</option>
          <option value="health">按健康状态分组</option>
        </select>
        {#if groupMode === 'label' && availableLabelKeys.length > 0}
          <select 
            class="label-key-select" 
            value={selectedLabelKey} 
            on:change={handleLabelKeyChange}
            aria-label="选择标签"
          >
            {#each availableLabelKeys as labelKey}
              <option value={labelKey}>{labelKey}</option>
            {/each}
          </select>
        {/if}
        <button class="refresh-button" on:click={handleRefresh} disabled={refreshing}>
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
      {#if groupMode === 'compose'}
        <!-- Grouped by compose project -->
        {@const { grouped, ungrouped } = groupContainersByCompose(containers)}
        
        <!-- Quick navigation sidebar -->
        {#if grouped.size > 0 || ungrouped.length > 0}
          <div class="quick-nav-sidebar">
            <div class="quick-nav-title">快速跳转</div>
            {#each Array.from(grouped.keys()) as projectName}
              <button class="quick-nav-item" on:click={() => scrollToGroup(projectName)}>
                {projectName}
              </button>
            {/each}
            {#if ungrouped.length > 0}
              <button class="quick-nav-item" on:click={() => scrollToGroup('_ungrouped_')}>
                独立容器
              </button>
            {/if}
          </div>
        {/if}
        
        {#if grouped.size > 0}
          {#each Array.from(grouped.entries()) as [projectName, projectContainers] (projectName)}
            <div class="compose-group" id="group-{projectName}">
              <button 
                class="compose-group-header" 
                class:compact={displayMode === 'compact'}
                on:click={() => toggleGroupCollapse(projectName)}
                aria-expanded={!collapsedGroups.has(projectName)}
                aria-label={`${projectName} compose group, ${projectContainers.length} containers`}
              >
                <span class="compose-icon">📚</span>
                <h3 class="compose-project-name">{projectName}</h3>
                <span class="compose-count">{projectContainers.length} 个容器</span>
                <span class="collapse-icon" aria-hidden="true">{collapsedGroups.has(projectName) ? '▶' : '▼'}</span>
              </button>
              {#if !collapsedGroups.has(projectName)}
              <div class="container-list" class:compact={displayMode === 'compact'}>
                {#each projectContainers as container (container.id)}
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
                        {#if container.compose_service}
                          <span class="compose-service-badge">{container.compose_service}</span>
                        {/if}
                        <span class="compact-image" title={container.image}>{container.image}</span>
                        <span class="compact-state">{container.status}</span>
                        <div class="compact-actions">
                          {#if container.state === 'running'}
                            <button 
                              class="action-btn-compact stop" 
                              on:click={() => handleAction(container.id, 'stop', container.is_self ?? false)}
                              disabled={container.is_self}
                              title={container.is_self ? '无法停止本应用容器' : '停止'}
                            >
                              ⏸️
                            </button>
                            <button 
                              class="action-btn-compact restart" 
                              on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                              disabled={container.is_self}
                              title={container.is_self ? '无法重启本应用容器' : '重启'}
                            >
                              🔄
                            </button>
                          {:else if ['exited', 'created', 'dead'].includes(container.state)}
                            <button 
                              class="action-btn-compact start" 
                              on:click={() => handleAction(container.id, 'start', container.is_self ?? false)}
                              title="启动"
                            >
                              ▶️
                            </button>
                          {:else}
                            <button 
                              class="action-btn-compact restart" 
                              on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                              disabled={container.is_self}
                              title={container.is_self ? '无法重启本应用容器' : '重启'}
                            >
                              🔄
                            </button>
                          {/if}
                          <a 
                            class="action-btn-compact logs" 
                            href={`#/logs/${container.id}`}
                            title="查看日志"
                          >
                            📋
                          </a>
                          <a 
                            class="action-btn-compact details" 
                            href={`#/container/${container.id}`}
                            title="查看详情"
                          >
                            ℹ️
                          </a>
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
                          {#if container.compose_service}
                            <span class="compose-service-badge">{container.compose_service}</span>
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
                            on:click={() => handleAction(container.id, 'stop', container.is_self ?? false)}
                            disabled={container.is_self}
                            title={container.is_self ? '无法停止本应用容器' : ''}
                          >
                            ⏸️ 停止
                          </button>
                          <button 
                            class="action-btn restart" 
                            on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                            disabled={container.is_self}
                            title={container.is_self ? '无法重启本应用容器' : ''}
                          >
                            🔄 重启
                          </button>
                        {:else if ['exited', 'created', 'dead'].includes(container.state)}
                          <button 
                            class="action-btn start" 
                            on:click={() => handleAction(container.id, 'start', container.is_self ?? false)}
                          >
                            ▶️ 启动
                          </button>
                        {:else}
                          <button 
                            class="action-btn restart" 
                            on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                            disabled={container.is_self}
                            title={container.is_self ? '无法重启本应用容器' : ''}
                          >
                            🔄 重启
                          </button>
                        {/if}
                        <a 
                          class="action-btn logs" 
                          href={`#/logs/${container.id}`}
                        >
                          📋 日志
                        </a>
                        <a 
                          class="action-btn details" 
                          href={`#/container/${container.id}`}
                        >
                          ℹ️ 详情
                        </a>
                      </div>
                    {/if}
                  </div>
                {/each}
              </div>
              {/if}
            </div>
          {/each}
        {/if}
        
        {#if ungrouped.length > 0}
          <div class="compose-group" id="group-_ungrouped_">
            <button 
              class="compose-group-header" 
              class:compact={displayMode === 'compact'}
              on:click={() => toggleGroupCollapse('_ungrouped_')}
              aria-expanded={!collapsedGroups.has('_ungrouped_')}
              aria-label={`独立容器 group, ${ungrouped.length} containers`}
            >
              <span class="compose-icon">📦</span>
              <h3 class="compose-project-name">独立容器</h3>
              <span class="compose-count">{ungrouped.length} 个容器</span>
              <span class="collapse-icon" aria-hidden="true">{collapsedGroups.has('_ungrouped_') ? '▶' : '▼'}</span>
            </button>
            {#if !collapsedGroups.has('_ungrouped_')}
            <div class="container-list" class:compact={displayMode === 'compact'}>
              {#each ungrouped as container (container.id)}
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
                            on:click={() => handleAction(container.id, 'stop', container.is_self ?? false)}
                            disabled={container.is_self}
                            title={container.is_self ? '无法停止本应用容器' : '停止'}
                          >
                            ⏸️
                          </button>
                          <button 
                            class="action-btn-compact restart" 
                            on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                            disabled={container.is_self}
                            title={container.is_self ? '无法重启本应用容器' : '重启'}
                          >
                            🔄
                          </button>
                        {:else if ['exited', 'created', 'dead'].includes(container.state)}
                          <button 
                            class="action-btn-compact start" 
                            on:click={() => handleAction(container.id, 'start', container.is_self ?? false)}
                            title="启动"
                          >
                            ▶️
                          </button>
                        {:else}
                          <button 
                            class="action-btn-compact restart" 
                            on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                            disabled={container.is_self}
                            title={container.is_self ? '无法重启本应用容器' : '重启'}
                          >
                            🔄
                          </button>
                        {/if}
                        <a 
                          class="action-btn-compact logs" 
                          href={`#/logs/${container.id}`}
                          title="查看日志"
                        >
                          📋
                        </a>
                        <a 
                          class="action-btn-compact details" 
                          href={`#/container/${container.id}`}
                          title="查看详情"
                        >
                          ℹ️
                        </a>
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
                          on:click={() => handleAction(container.id, 'stop', container.is_self ?? false)}
                          disabled={container.is_self}
                          title={container.is_self ? '无法停止本应用容器' : ''}
                        >
                          ⏸️ 停止
                        </button>
                        <button 
                          class="action-btn restart" 
                          on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                          disabled={container.is_self}
                          title={container.is_self ? '无法重启本应用容器' : ''}
                        >
                          🔄 重启
                        </button>
                      {:else if ['exited', 'created', 'dead'].includes(container.state)}
                        <button 
                          class="action-btn start" 
                          on:click={() => handleAction(container.id, 'start', container.is_self ?? false)}
                        >
                          ▶️ 启动
                        </button>
                      {:else}
                        <button 
                          class="action-btn restart" 
                          on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                          disabled={container.is_self}
                          title={container.is_self ? '无法重启本应用容器' : ''}
                        >
                          🔄 重启
                        </button>
                      {/if}
                      <a 
                        class="action-btn logs" 
                        href={`#/logs/${container.id}`}
                      >
                        📋 日志
                      </a>
                      <a 
                        class="action-btn details" 
                        href={`#/container/${container.id}`}
                      >
                        ℹ️ 详情
                      </a>
                    </div>
                  {/if}
                </div>
              {/each}
            </div>
            {/if}
          </div>
        {/if}
      {:else if groupMode === 'label' && selectedLabelKey}
        <!-- Grouped by selected label -->
        {@const { grouped, ungrouped } = groupContainersByLabel(containers, selectedLabelKey)}
        
        <!-- Quick navigation sidebar -->
        {#if grouped.size > 0 || ungrouped.length > 0}
          <div class="quick-nav-sidebar">
            <div class="quick-nav-title">快速跳转</div>
            {#each Array.from(grouped.keys()) as groupName}
              <button class="quick-nav-item" on:click={() => scrollToGroup(`label-${groupName}`)}>
                {groupName}
              </button>
            {/each}
            {#if ungrouped.length > 0}
              <button class="quick-nav-item" on:click={() => scrollToGroup('_ungrouped_label_')}>
                未分组
              </button>
            {/if}
          </div>
        {/if}
        
        {#if grouped.size > 0}
          {#each Array.from(grouped.entries()) as [labelValue, labelContainers] (labelValue)}
            <div class="compose-group" id="group-label-{labelValue}">
              <button 
                class="compose-group-header" 
                class:compact={displayMode === 'compact'}
                on:click={() => toggleGroupCollapse(`label-${labelValue}`)}
                aria-expanded={!collapsedGroups.has(`label-${labelValue}`)}
                aria-label={`${labelValue} label group, ${labelContainers.length} containers`}
              >
                <span class="compose-icon">🏷️</span>
                <h3 class="compose-project-name">{labelValue}</h3>
                <span class="compose-count">{labelContainers.length} 个容器</span>
                <span class="collapse-icon" aria-hidden="true">{collapsedGroups.has(`label-${labelValue}`) ? '▶' : '▼'}</span>
              </button>
              {#if !collapsedGroups.has(`label-${labelValue}`)}
              <div class="container-list" class:compact={displayMode === 'compact'}>
                {#each labelContainers as container (container.id)}
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
                        {#if container.compose_service}
                          <span class="compose-service-badge">{container.compose_service}</span>
                        {/if}
                        <span class="compact-image" title={container.image}>{container.image}</span>
                        <span class="compact-state">{container.status}</span>
                        <div class="compact-actions">
                          {#if container.state === 'running'}
                            <button 
                              class="action-btn-compact stop" 
                              on:click={() => handleAction(container.id, 'stop', container.is_self ?? false)}
                              disabled={container.is_self}
                              title={container.is_self ? '无法停止本应用容器' : '停止'}
                            >
                              ⏸️
                            </button>
                            <button 
                              class="action-btn-compact restart" 
                              on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                              disabled={container.is_self}
                              title={container.is_self ? '无法重启本应用容器' : '重启'}
                            >
                              🔄
                            </button>
                          {:else if ['exited', 'created', 'dead'].includes(container.state)}
                            <button 
                              class="action-btn-compact start" 
                              on:click={() => handleAction(container.id, 'start', container.is_self ?? false)}
                              title="启动"
                            >
                              ▶️
                            </button>
                          {:else}
                            <button 
                              class="action-btn-compact restart" 
                              on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                              disabled={container.is_self}
                              title={container.is_self ? '无法重启本应用容器' : '重启'}
                            >
                              🔄
                            </button>
                          {/if}
                          <a 
                            class="action-btn-compact logs" 
                            href={`#/logs/${container.id}`}
                            title="查看日志"
                          >
                            📋
                          </a>
                          <a 
                            class="action-btn-compact details" 
                            href={`#/container/${container.id}`}
                            title="查看详情"
                          >
                            ℹ️
                          </a>
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
                          {#if container.compose_service}
                            <span class="compose-service-badge">{container.compose_service}</span>
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
                            on:click={() => handleAction(container.id, 'stop', container.is_self ?? false)}
                            disabled={container.is_self}
                            title={container.is_self ? '无法停止本应用容器' : ''}
                          >
                            ⏸️ 停止
                          </button>
                          <button 
                            class="action-btn restart" 
                            on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                            disabled={container.is_self}
                            title={container.is_self ? '无法重启本应用容器' : ''}
                          >
                            🔄 重启
                          </button>
                        {:else if ['exited', 'created', 'dead'].includes(container.state)}
                          <button 
                            class="action-btn start" 
                            on:click={() => handleAction(container.id, 'start', container.is_self ?? false)}
                          >
                            ▶️ 启动
                          </button>
                        {:else}
                          <button 
                            class="action-btn restart" 
                            on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                            disabled={container.is_self}
                            title={container.is_self ? '无法重启本应用容器' : ''}
                          >
                            🔄 重启
                          </button>
                        {/if}
                        <a 
                          class="action-btn logs" 
                          href={`#/logs/${container.id}`}
                        >
                          📋 日志
                        </a>
                        <a 
                          class="action-btn details" 
                          href={`#/container/${container.id}`}
                        >
                          ℹ️ 详情
                        </a>
                      </div>
                    {/if}
                  </div>
                {/each}
              </div>
              {/if}
            </div>
          {/each}
        {/if}
        
        {#if ungrouped.length > 0}
          <div class="compose-group" id="group-_ungrouped_label_">
            <button 
              class="compose-group-header" 
              class:compact={displayMode === 'compact'}
              on:click={() => toggleGroupCollapse('_ungrouped_label_')}
              aria-expanded={!collapsedGroups.has('_ungrouped_label_')}
              aria-label={`未分组容器, ${ungrouped.length} containers`}
            >
              <span class="compose-icon">📦</span>
              <h3 class="compose-project-name">未分组容器</h3>
              <span class="compose-count">{ungrouped.length} 个容器</span>
              <span class="collapse-icon" aria-hidden="true">{collapsedGroups.has('_ungrouped_label_') ? '▶' : '▼'}</span>
            </button>
            {#if !collapsedGroups.has('_ungrouped_label_')}
            <div class="container-list" class:compact={displayMode === 'compact'}>
              {#each ungrouped as container (container.id)}
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
                            on:click={() => handleAction(container.id, 'stop', container.is_self ?? false)}
                            disabled={container.is_self}
                            title={container.is_self ? '无法停止本应用容器' : '停止'}
                          >
                            ⏸️
                          </button>
                          <button 
                            class="action-btn-compact restart" 
                            on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                            disabled={container.is_self}
                            title={container.is_self ? '无法重启本应用容器' : '重启'}
                          >
                            🔄
                          </button>
                        {:else if ['exited', 'created', 'dead'].includes(container.state)}
                          <button 
                            class="action-btn-compact start" 
                            on:click={() => handleAction(container.id, 'start', container.is_self ?? false)}
                            title="启动"
                          >
                            ▶️
                          </button>
                        {:else}
                          <button 
                            class="action-btn-compact restart" 
                            on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                            disabled={container.is_self}
                            title={container.is_self ? '无法重启本应用容器' : '重启'}
                          >
                            🔄
                          </button>
                        {/if}
                        <a 
                          class="action-btn-compact logs" 
                          href={`#/logs/${container.id}`}
                          title="查看日志"
                        >
                          📋
                        </a>
                        <a 
                          class="action-btn-compact details" 
                          href={`#/container/${container.id}`}
                          title="查看详情"
                        >
                          ℹ️
                        </a>
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
                          on:click={() => handleAction(container.id, 'stop', container.is_self ?? false)}
                          disabled={container.is_self}
                          title={container.is_self ? '无法停止本应用容器' : ''}
                        >
                          ⏸️ 停止
                        </button>
                        <button 
                          class="action-btn restart" 
                          on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                          disabled={container.is_self}
                          title={container.is_self ? '无法重启本应用容器' : ''}
                        >
                          🔄 重启
                        </button>
                      {:else if ['exited', 'created', 'dead'].includes(container.state)}
                        <button 
                          class="action-btn start" 
                          on:click={() => handleAction(container.id, 'start', container.is_self ?? false)}
                        >
                          ▶️ 启动
                        </button>
                      {:else}
                        <button 
                          class="action-btn restart" 
                          on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                          disabled={container.is_self}
                          title={container.is_self ? '无法重启本应用容器' : ''}
                        >
                          🔄 重启
                        </button>
                      {/if}
                      <a 
                        class="action-btn logs" 
                        href={`#/logs/${container.id}`}
                      >
                        📋 日志
                      </a>
                      <a 
                        class="action-btn details" 
                        href={`#/container/${container.id}`}
                      >
                        ℹ️ 详情
                      </a>
                    </div>
                  {/if}
                </div>
              {/each}
            </div>
            {/if}
          </div>
        {/if}
      {:else if groupMode === 'status'}
        <!-- Grouped by status -->
        {@const { grouped, ungrouped } = groupContainersByStatus(containers)}
        
        <!-- Quick navigation sidebar -->
        {#if grouped.size > 0}
          <div class="quick-nav-sidebar">
            <div class="quick-nav-title">快速跳转</div>
            {#each Array.from(grouped.keys()) as statusName}
              <button class="quick-nav-item" on:click={() => scrollToGroup(`status-${statusName}`)}>
                {stateEmojis[statusName] || '⚪'} {statusName}
              </button>
            {/each}
          </div>
        {/if}
        
        {#if grouped.size > 0}
          {#each Array.from(grouped.entries()) as [statusName, statusContainers] (statusName)}
            <div class="compose-group" id="group-status-{statusName}">
              <button 
                class="compose-group-header" 
                class:compact={displayMode === 'compact'}
                on:click={() => toggleGroupCollapse(`status-${statusName}`)}
                aria-expanded={!collapsedGroups.has(`status-${statusName}`)}
                aria-label={`${statusName} status group, ${statusContainers.length} containers`}
              >
                <span class="compose-icon">{stateEmojis[statusName] || '⚪'}</span>
                <h3 class="compose-project-name">{statusName}</h3>
                <span class="compose-count">{statusContainers.length} 个容器</span>
                <span class="collapse-icon" aria-hidden="true">{collapsedGroups.has(`status-${statusName}`) ? '▶' : '▼'}</span>
              </button>
              {#if !collapsedGroups.has(`status-${statusName}`)}
              <div class="container-list" class:compact={displayMode === 'compact'}>
                {#each statusContainers as container (container.id)}
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
                        {#if container.compose_service}
                          <span class="compose-service-badge">{container.compose_service}</span>
                        {/if}
                        <span class="compact-image" title={container.image}>{container.image}</span>
                        <span class="compact-state">{container.status}</span>
                        <div class="compact-actions">
                          {#if container.state === 'running'}
                            <button 
                              class="action-btn-compact stop" 
                              on:click={() => handleAction(container.id, 'stop', container.is_self ?? false)}
                              disabled={container.is_self}
                              title={container.is_self ? '无法停止本应用容器' : '停止'}
                            >
                              ⏸️
                            </button>
                            <button 
                              class="action-btn-compact restart" 
                              on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                              disabled={container.is_self}
                              title={container.is_self ? '无法重启本应用容器' : '重启'}
                            >
                              🔄
                            </button>
                          {:else if ['exited', 'created', 'dead'].includes(container.state)}
                            <button 
                              class="action-btn-compact start" 
                              on:click={() => handleAction(container.id, 'start', container.is_self ?? false)}
                              title="启动"
                            >
                              ▶️
                            </button>
                          {:else}
                            <button 
                              class="action-btn-compact restart" 
                              on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                              disabled={container.is_self}
                              title={container.is_self ? '无法重启本应用容器' : '重启'}
                            >
                              🔄
                            </button>
                          {/if}
                          <a 
                            class="action-btn-compact logs" 
                            href={`#/logs/${container.id}`}
                            title="查看日志"
                          >
                            📋
                          </a>
                          <a 
                            class="action-btn-compact details" 
                            href={`#/container/${container.id}`}
                            title="查看详情"
                          >
                            ℹ️
                          </a>
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
                          {#if container.compose_service}
                            <span class="compose-service-badge">{container.compose_service}</span>
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
                            on:click={() => handleAction(container.id, 'stop', container.is_self ?? false)}
                            disabled={container.is_self}
                            title={container.is_self ? '无法停止本应用容器' : ''}
                          >
                            ⏸️ 停止
                          </button>
                          <button 
                            class="action-btn restart" 
                            on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                            disabled={container.is_self}
                            title={container.is_self ? '无法重启本应用容器' : ''}
                          >
                            🔄 重启
                          </button>
                        {:else if ['exited', 'created', 'dead'].includes(container.state)}
                          <button 
                            class="action-btn start" 
                            on:click={() => handleAction(container.id, 'start', container.is_self ?? false)}
                          >
                            ▶️ 启动
                          </button>
                        {:else}
                          <button 
                            class="action-btn restart" 
                            on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                            disabled={container.is_self}
                            title={container.is_self ? '无法重启本应用容器' : ''}
                          >
                            🔄 重启
                          </button>
                        {/if}
                        <a 
                          class="action-btn logs" 
                          href={`#/logs/${container.id}`}
                        >
                          📋 日志
                        </a>
                        <a 
                          class="action-btn details" 
                          href={`#/container/${container.id}`}
                        >
                          ℹ️ 详情
                        </a>
                      </div>
                    {/if}
                  </div>
                {/each}
              </div>
              {/if}
            </div>
          {/each}
        {/if}
      {:else if groupMode === 'health'}
        <!-- Grouped by health -->
        {@const { grouped, ungrouped } = groupContainersByHealth(containers)}
        
        <!-- Quick navigation sidebar -->
        {#if grouped.size > 0}
          <div class="quick-nav-sidebar">
            <div class="quick-nav-title">快速跳转</div>
            {#each Array.from(grouped.keys()) as healthName}
              <button class="quick-nav-item" on:click={() => scrollToGroup(`health-${healthName}`)}>
                {healthEmojis[healthName] || '⚪'} {healthName}
              </button>
            {/each}
          </div>
        {/if}
        
        {#if grouped.size > 0}
          {#each Array.from(grouped.entries()) as [healthName, healthContainers] (healthName)}
            <div class="compose-group" id="group-health-{healthName}">
              <button 
                class="compose-group-header" 
                class:compact={displayMode === 'compact'}
                on:click={() => toggleGroupCollapse(`health-${healthName}`)}
                aria-expanded={!collapsedGroups.has(`health-${healthName}`)}
                aria-label={`${healthName} health group, ${healthContainers.length} containers`}
              >
                <span class="compose-icon">{healthEmojis[healthName] || '⚪'}</span>
                <h3 class="compose-project-name">{healthName}</h3>
                <span class="compose-count">{healthContainers.length} 个容器</span>
                <span class="collapse-icon" aria-hidden="true">{collapsedGroups.has(`health-${healthName}`) ? '▶' : '▼'}</span>
              </button>
              {#if !collapsedGroups.has(`health-${healthName}`)}
              <div class="container-list" class:compact={displayMode === 'compact'}>
                {#each healthContainers as container (container.id)}
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
                        {#if container.compose_service}
                          <span class="compose-service-badge">{container.compose_service}</span>
                        {/if}
                        <span class="compact-image" title={container.image}>{container.image}</span>
                        <span class="compact-state">{container.status}</span>
                        <div class="compact-actions">
                          {#if container.state === 'running'}
                            <button 
                              class="action-btn-compact stop" 
                              on:click={() => handleAction(container.id, 'stop', container.is_self ?? false)}
                              disabled={container.is_self}
                              title={container.is_self ? '无法停止本应用容器' : '停止'}
                            >
                              ⏸️
                            </button>
                            <button 
                              class="action-btn-compact restart" 
                              on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                              disabled={container.is_self}
                              title={container.is_self ? '无法重启本应用容器' : '重启'}
                            >
                              🔄
                            </button>
                          {:else if ['exited', 'created', 'dead'].includes(container.state)}
                            <button 
                              class="action-btn-compact start" 
                              on:click={() => handleAction(container.id, 'start', container.is_self ?? false)}
                              title="启动"
                            >
                              ▶️
                            </button>
                          {:else}
                            <button 
                              class="action-btn-compact restart" 
                              on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                              disabled={container.is_self}
                              title={container.is_self ? '无法重启本应用容器' : '重启'}
                            >
                              🔄
                            </button>
                          {/if}
                          <a 
                            class="action-btn-compact logs" 
                            href={`#/logs/${container.id}`}
                            title="查看日志"
                          >
                            📋
                          </a>
                          <a 
                            class="action-btn-compact details" 
                            href={`#/container/${container.id}`}
                            title="查看详情"
                          >
                            ℹ️
                          </a>
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
                          {#if container.compose_service}
                            <span class="compose-service-badge">{container.compose_service}</span>
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
                            on:click={() => handleAction(container.id, 'stop', container.is_self ?? false)}
                            disabled={container.is_self}
                            title={container.is_self ? '无法停止本应用容器' : ''}
                          >
                            ⏸️ 停止
                          </button>
                          <button 
                            class="action-btn restart" 
                            on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                            disabled={container.is_self}
                            title={container.is_self ? '无法重启本应用容器' : ''}
                          >
                            🔄 重启
                          </button>
                        {:else if ['exited', 'created', 'dead'].includes(container.state)}
                          <button 
                            class="action-btn start" 
                            on:click={() => handleAction(container.id, 'start', container.is_self ?? false)}
                          >
                            ▶️ 启动
                          </button>
                        {:else}
                          <button 
                            class="action-btn restart" 
                            on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                            disabled={container.is_self}
                            title={container.is_self ? '无法重启本应用容器' : ''}
                          >
                            🔄 重启
                          </button>
                        {/if}
                        <a 
                          class="action-btn logs" 
                          href={`#/logs/${container.id}`}
                        >
                          📋 日志
                        </a>
                        <a 
                          class="action-btn details" 
                          href={`#/container/${container.id}`}
                        >
                          ℹ️ 详情
                        </a>
                      </div>
                    {/if}
                  </div>
                {/each}
              </div>
              {/if}
            </div>
          {/each}
        {/if}
      {:else}
        <!-- Ungrouped list -->
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
                      on:click={() => handleAction(container.id, 'stop', container.is_self ?? false)}
                      disabled={container.is_self}
                      title={container.is_self ? '无法停止本应用容器' : '停止'}
                    >
                      ⏸️
                    </button>
                    <button 
                      class="action-btn-compact restart" 
                      on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                      disabled={container.is_self}
                      title={container.is_self ? '无法重启本应用容器' : '重启'}
                    >
                      🔄
                    </button>
                  {:else if ['exited', 'created', 'dead'].includes(container.state)}
                    <button 
                      class="action-btn-compact start" 
                      on:click={() => handleAction(container.id, 'start', container.is_self ?? false)}
                      title="启动"
                    >
                      ▶️
                    </button>
                  {:else}
                    <button 
                      class="action-btn-compact restart" 
                      on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                      disabled={container.is_self}
                      title={container.is_self ? '无法重启本应用容器' : '重启'}
                    >
                      🔄
                    </button>
                  {/if}
                  <a 
                    class="action-btn-compact logs" 
                    href={`#/logs/${container.id}`}
                    title="查看日志"
                  >
                    📋
                  </a>
                  <a 
                    class="action-btn-compact details" 
                    href={`#/container/${container.id}`}
                    title="查看详情"
                  >
                    ℹ️
                  </a>
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
                    on:click={() => handleAction(container.id, 'stop', container.is_self ?? false)}
                    disabled={container.is_self}
                    title={container.is_self ? '无法停止本应用容器' : ''}
                  >
                    ⏸️ 停止
                  </button>
                  <button 
                    class="action-btn restart" 
                    on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                    disabled={container.is_self}
                    title={container.is_self ? '无法重启本应用容器' : ''}
                  >
                    🔄 重启
                  </button>
                {:else if ['exited', 'created', 'dead'].includes(container.state)}
                  <button 
                    class="action-btn start" 
                    on:click={() => handleAction(container.id, 'start', container.is_self ?? false)}
                  >
                    ▶️ 启动
                  </button>
                {:else}
                  <button 
                    class="action-btn restart" 
                    on:click={() => handleAction(container.id, 'restart', container.is_self ?? false)}
                    disabled={container.is_self}
                    title={container.is_self ? '无法重启本应用容器' : ''}
                  >
                    🔄 重启
                  </button>
                {/if}
                <a 
                  class="action-btn logs" 
                  href={`#/logs/${container.id}`}
                >
                  📋 日志
                </a>
                <a 
                  class="action-btn details" 
                  href={`#/container/${container.id}`}
                >
                  ℹ️ 详情
                </a>
              </div>
            {/if}
          </div>
        {/each}
        </div>
      {/if}
    {/if}
  </main>
</div>

<style>
  .home-container {
    min-height: 100vh;
    background: var(--color-background, #f5f5f4);
  }
  
  /* Floating header that appears when scrolled - positioned inside the main header area */
  .floating-header {
    position: fixed;
    top: 0;
    left: 50%;
    transform: translateX(-50%) translateY(-100%);
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 1.5rem;
    padding: 1rem 2rem;
    background: var(--color-primary, #171717);
    color: var(--color-background, #f5f5f4);
    z-index: 101;
    opacity: 0;
    transition: opacity 0.3s ease-out, transform 0.3s ease-out;
    pointer-events: none;
    border-radius: 0 0 var(--radius, 0.25rem) var(--radius, 0.25rem);
  }
  
  .floating-header.visible {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
    pointer-events: auto;
  }
  
  .floating-header h2 {
    font-size: 1.25rem;
    font-weight: 600;
    margin: 0;
    font-family: var(--font-heading, "Playfair Display", serif);
  }
  
  .floating-header .header-actions {
    display: flex;
    gap: 0.75rem;
  }
  
  .floating-header .mode-toggle,
  .floating-header .refresh-button,
  .floating-header .group-mode-select,
  .floating-header .label-key-select {
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    color: var(--color-background, #f5f5f4);
    padding: 0.4rem 0.75rem;
    font-size: 0.85rem;
  }
  
  .floating-header .mode-toggle:hover,
  .floating-header .refresh-button:hover:not(:disabled),
  .floating-header .group-mode-select:hover,
  .floating-header .label-key-select:hover {
    background: rgba(255, 255, 255, 0.2);
    border-color: rgba(255, 255, 255, 0.3);
  }
  
  .floating-header .group-mode-select option,
  .floating-header .label-key-select option {
    background: var(--color-primary, #171717);
    color: var(--color-background, #f5f5f4);
  }
  
  .main-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
    position: relative;
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
  .refresh-button,
  .group-mode-select,
  .label-key-select {
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
  .refresh-button:hover:not(:disabled),
  .group-mode-select:hover,
  .label-key-select:hover {
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
  
  .action-btn.logs {
    background: var(--color-secondary, #525252);
    color: white;
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }
  
  .action-btn.logs:hover {
    background: #404040;
  }
  
  .action-btn.details {
    background: #0284c7;
    color: white;
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }
  
  .action-btn.details:hover {
    background: #0369a1;
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
  
  .action-btn-compact.logs {
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }
  
  /* Mobile responsive styles */
  @media (max-width: 640px) {
    .main-content {
      padding: 1rem;
    }
    
    .content-header {
      flex-wrap: wrap;
      gap: 1rem;
    }
    
    .content-header h2 {
      font-size: 1.5rem;
    }
    
    /* Standard mode mobile: action buttons wrap to new line */
    .container-list:not(.compact) .container-item {
      flex-wrap: wrap;
    }
    
    .container-list:not(.compact) .container-actions {
      width: 100%;
      margin-top: 0.75rem;
      padding-top: 0.75rem;
      border-top: 1px solid rgba(0, 0, 0, 0.1);
      justify-content: flex-start;
    }
    
    /* Compact mode mobile: actions float on right side above content */
    .container-list.compact .container-item {
      position: relative;
      padding-right: 4.5rem;
    }
    
    .container-compact {
      flex-wrap: nowrap;
      overflow: hidden;
    }
    
    .compact-image {
      display: none;
    }
    
    .compact-state {
      min-width: auto;
      flex: 1;
    }
    
    .compact-actions {
      position: absolute;
      right: 0.5rem;
      top: 50%;
      transform: translateY(-50%);
      background: linear-gradient(to right, transparent, var(--color-surface, #e7e5e4) 20%);
      padding-left: 1rem;
    }

    /* Floating header mobile responsive - smaller and more compact buttons */
    .floating-header {
      padding: 1rem;
      gap: 0.5rem;
    }

    .floating-header h2 {
      font-size: 0.9rem;
      flex-shrink: 0;
    }

    .floating-header .header-actions {
      gap: 0.5rem;
    }

    .floating-header .mode-toggle,
    .floating-header .refresh-button,
    .floating-header .group-mode-select,
    .floating-header .label-key-select {
      padding: 0.35rem 0.5rem;
      font-size: 0.75rem;
      min-width: auto;
    }
  }

  /* Hide title on smaller screens to save space */
  @media (max-width: 480px) {
    .floating-header h2 {
      display: none;
    }
  }

  /* Hide button text on very small screens, show icons only */
  @media (max-width: 400px) {
    .floating-header .mode-toggle .mode-text {
      display: none;
    }

    .floating-header .refresh-button {
      font-size: 0;
      gap: 0;
    }

    .floating-header .refresh-icon {
      font-size: 1rem;
    }
  }
  
  /* Compose group styles */
  .compose-group {
    margin-bottom: 2rem;
  }
  
  .compose-group-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    background: var(--color-primary, #171717);
    color: var(--color-background, #f5f5f4);
    border-radius: var(--radius, 0.25rem) var(--radius, 0.25rem) 0 0;
    margin-bottom: 0.5rem;
    width: 100%;
    border: none;
    cursor: pointer;
    text-align: left;
    font-family: inherit;
    transition: background 0.2s, color 0.2s;
  }
  
  .compose-group-header:hover {
    background: #2a2a2a;
    color: var(--color-text, #0a0a0a);
  }
  
  .compose-group-header.compact {
    padding: 0.5rem 0.75rem;
  }
  
  .compose-icon {
    font-size: 1.25rem;
  }
  
  .compose-project-name {
    font-size: 1.1rem;
    font-weight: 600;
    margin: 0;
    font-family: var(--font-heading, "Playfair Display", serif);
  }
  
  .compose-group-header.compact .compose-project-name {
    font-size: 0.95rem;
  }
  
  .compose-count {
    margin-left: auto;
    font-size: 0.9rem;
    opacity: 0.7;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .compose-group-header.compact .compose-count {
    font-size: 0.8rem;
  }
  
  .collapse-icon {
    font-size: 0.8rem;
    opacity: 0.7;
    margin-left: 0.5rem;
  }
  
  .compose-service-badge {
    background: var(--color-secondary, #525252);
    color: white;
    font-size: 0.7rem;
    padding: 0.15rem 0.5rem;
    border-radius: var(--radius, 0.25rem);
    font-weight: 500;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .compose-group .container-list {
    margin-top: 0;
  }
  
  .compose-group .container-item {
    border-radius: 0;
    margin-bottom: 0;
  }
  
  .compose-group .container-item:first-child {
    border-radius: var(--radius, 0.25rem) var(--radius, 0.25rem) 0 0;
  }
  
  .compose-group .container-item:last-child {
    border-radius: 0 0 var(--radius, 0.25rem) var(--radius, 0.25rem);
  }
  
  .compose-group .container-item:only-child {
    border-radius: var(--radius, 0.25rem);
  }
  
  /* Quick navigation sidebar */
  .quick-nav-sidebar {
    position: sticky;
    top: 5rem;
    float: left;
    margin-left: -260px;
    margin-right: 20px;
    background: var(--color-surface, #e7e5e4);
    border: 1px solid rgba(0, 0, 0, 0.1);
    border-radius: var(--radius, 0.25rem);
    padding: 0.75rem;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
    max-height: calc(100vh - 6rem);
    overflow-y: auto;
    z-index: 50;
  }
  
  /* Desktop: sidebar with fixed width when floating */
  @media (min-width: 1025px) {
    .quick-nav-sidebar {
      max-width: 240px;
    }
  }
  
  .quick-nav-title {
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--color-text, #0a0a0a);
    margin-bottom: 0.5rem;
    padding-bottom: 0.5rem;
    border-bottom: 1px solid rgba(0, 0, 0, 0.1);
    font-family: var(--font-heading, "Playfair Display", serif);
  }
  
  .quick-nav-item {
    display: block;
    width: 100%;
    text-align: left;
    padding: 0.5rem 0.75rem;
    margin-bottom: 0.25rem;
    border: none;
    background: transparent;
    color: var(--color-text, #0a0a0a);
    font-size: 0.85rem;
    font-family: var(--font-body, "Merriweather", serif);
    border-radius: var(--radius, 0.25rem);
    cursor: pointer;
    transition: all 0.2s;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .quick-nav-item:hover {
    background: var(--color-primary, #171717);
    color: var(--color-background, #f5f5f4);
  }
  
  /* Mobile responsive adjustments for quick nav */
  @media (max-width: 1024px) {
    .quick-nav-sidebar {
      position: static;
      float: none;
      margin-left: 0;
      margin-right: 0;
      margin-bottom: 1.5rem;
      max-width: 100%;
      max-height: none;
      display: flex;
      flex-wrap: wrap;
      gap: 0.5rem;
      align-items: center;
      padding: 1rem;
    }
    
    .quick-nav-title {
      width: 100%;
      margin-bottom: 0.5rem;
      padding-bottom: 0.5rem;
    }
    
    .quick-nav-item {
      width: auto;
      flex: 0 0 auto;
      margin-bottom: 0;
      padding: 0.4rem 0.75rem;
      font-size: 0.8rem;
    }
  }
</style>
