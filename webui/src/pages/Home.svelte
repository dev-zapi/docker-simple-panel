<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { slide } from 'svelte/transition';
  import Header from '../components/Header.svelte';
  import PageLayout from '../components/PageLayout.svelte';
  import ContainerList from '../components/ContainerList.svelte';
  import { containerApi } from '../services/api';
  import type { Container } from '../types';
  
  let containers: Container[] = [];
  let loading = true;
  let error = '';
  let refreshing = false;
  let displayMode: 'compact' | 'standard' = 'standard';
  let groupMode: 'none' | 'compose' | 'label' | 'status-health' = 'none';
  let sortMode: 'none' | 'name' | 'created' | 'state-health' | 'compose' = 'none';
  let actionError = '';
  let collapsedGroups: Set<string> = new Set();
  let selectedLabelKey: string = '';
  let availableLabelKeys: string[] = [];
  let filterText: string = '';
  
  // Track loading state for individual container actions
  let loadingActions = new Map<string, 'start' | 'stop' | 'restart'>();
  
  // Track if component has mounted to prevent reload during initial setup
  let isMounted = false;
  let filterDebounceTimer: number | undefined = undefined;
  const FILTER_DEBOUNCE_DELAY = 500; // milliseconds to wait before reloading after filter text changes
  const REFRESH_DELAY = 200; // milliseconds to wait after action before refreshing
  

  
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
  
  // Load display mode, group mode, and sort mode from localStorage
  onMount(() => {
    const savedMode = localStorage.getItem('displayMode');
    if (savedMode === 'compact' || savedMode === 'standard') {
      displayMode = savedMode;
    }
    const savedGroupMode = localStorage.getItem('groupMode');
    if (savedGroupMode === 'none' || savedGroupMode === 'compose' || savedGroupMode === 'label' || savedGroupMode === 'status-health') {
      groupMode = savedGroupMode;
    } else if (savedGroupMode === 'status' || savedGroupMode === 'health') {
      // Migrate old modes to new combined mode
      groupMode = 'status-health';
      localStorage.setItem('groupMode', 'status-health');
    }
    const savedSortMode = localStorage.getItem('sortMode');
    if (savedSortMode === 'none' || savedSortMode === 'name' || savedSortMode === 'created' || savedSortMode === 'state-health' || savedSortMode === 'compose') {
      sortMode = savedSortMode;
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
    const savedFilterText = localStorage.getItem('containerFilterText');
    if (savedFilterText) {
      filterText = savedFilterText;
    }
    loadContainers();
    
    // Mark as mounted after initial load completes
    isMounted = true;
  });
  
  onDestroy(() => {
    // Clear debounce timer if it exists
    clearTimeout(filterDebounceTimer);
  });
  
  function toggleDisplayMode() {
    displayMode = displayMode === 'compact' ? 'standard' : 'compact';
    localStorage.setItem('displayMode', displayMode);
  }
  
  function handleGroupModeChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    groupMode = target.value as 'none' | 'compose' | 'label' | 'status-health';
    localStorage.setItem('groupMode', groupMode);
  }
  
  function handleSortModeChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    sortMode = target.value as 'none' | 'name' | 'created' | 'state-health' | 'compose';
    localStorage.setItem('sortMode', sortMode);
  }
  
  function handleFilterTextChange(event: Event) {
    const target = event.target as HTMLInputElement;
    filterText = target.value;
    localStorage.setItem('containerFilterText', filterText);
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
  
  // Sort containers based on selected sort mode
  function sortContainers(containers: Container[]): Container[] {
    if (sortMode === 'none') {
      return containers;
    }
    
    const sorted = [...containers];
    
    switch (sortMode) {
      case 'name':
        sorted.sort((a, b) => a.name.localeCompare(b.name));
        break;
      case 'created':
        // Sort by created timestamp, newest first
        sorted.sort((a, b) => b.created - a.created);
        break;
      case 'state-health':
        // Sort by state first, then by health
        sorted.sort((a, b) => {
          const stateCompare = a.state.localeCompare(b.state);
          if (stateCompare !== 0) return stateCompare;
          const healthA = a.health || 'none';
          const healthB = b.health || 'none';
          return healthA.localeCompare(healthB);
        });
        break;
      case 'compose':
        // Sort by compose project name, then by compose service, then by name
        sorted.sort((a, b) => {
          const composeA = a.compose_project || '';
          const composeB = b.compose_project || '';
          const composeCompare = composeA.localeCompare(composeB);
          if (composeCompare !== 0) return composeCompare;
          
          const serviceA = a.compose_service || '';
          const serviceB = b.compose_service || '';
          const serviceCompare = serviceA.localeCompare(serviceB);
          if (serviceCompare !== 0) return serviceCompare;
          
          return a.name.localeCompare(b.name);
        });
        break;
    }
    
    return sorted;
  }
  
  // Filter containers by name
  function filterContainersByName(containers: Container[]): Container[] {
    if (!filterText.trim()) {
      return containers;
    }
    const lowerFilter = filterText.toLowerCase();
    return containers.filter(container => 
      container.name.toLowerCase().includes(lowerFilter)
    );
  }
  
  // Group containers by compose project
  function groupContainersByCompose(containers: Container[]) {
    const sorted = sortContainers(containers);
    const grouped = new Map<string, Container[]>();
    const ungrouped: Container[] = [];
    
    for (const container of sorted) {
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
    const sorted = sortContainers(containers);
    const grouped = new Map<string, Container[]>();
    const ungrouped: Container[] = [];
    
    for (const container of sorted) {
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
  
  // Group containers by combined status and health
  function groupContainersByStatusHealth(containers: Container[]) {
    const sorted = sortContainers(containers);
    const grouped = new Map<string, Container[]>();
    
    for (const container of sorted) {
      const state = container.state;
      const health = container.health || 'none';
      // Create a combined key like "running-healthy" or "exited-none"
      const groupKey = `${state}-${health}`;
      const existing = grouped.get(groupKey) || [];
      existing.push(container);
      grouped.set(groupKey, existing);
    }
    
    return { grouped, ungrouped: [] };
  }
  
  // Helper to get display info for combined status-health group
  function getStatusHealthDisplay(groupKey: string) {
    const [state, health] = groupKey.split('-');
    const stateEmoji = stateEmojis[state] || '⚪';
    const healthEmoji = health !== 'none' ? (healthEmojis[health] || '') : '';
    const displayName = health !== 'none' ? `${state} (${health})` : state;
    return { stateEmoji, healthEmoji, displayName, state, health };
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
  
  // Derived filtered containers - updates when containers or filterText changes
  $: filteredContainers = filterContainersByName(containers);
  
  // Derived sorted containers - updates when filteredContainers or sortMode changes
  $: sortedContainers = sortContainers(filteredContainers);
  
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
  
  // Auto-reload containers when sort mode changes (after initial mount)
  $: if (isMounted && sortMode) {
    loadContainers();
  }
  
  // Auto-reload containers when group mode changes (after initial mount)
  $: if (isMounted && groupMode) {
    loadContainers();
  }
  
  // Auto-reload containers when filter text changes with debouncing (after initial mount)
  $: if (isMounted && filterText !== undefined) {
    // Clear existing timer
    clearTimeout(filterDebounceTimer);
    // Set new timer to reload after user stops typing
    filterDebounceTimer = setTimeout(() => {
      loadContainers();
    }, FILTER_DEBOUNCE_DELAY);
  }
  
  function handleLabelKeyChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    selectedLabelKey = target.value;
    localStorage.setItem('selectedLabelKey', selectedLabelKey);
    // Reload containers when user manually changes label key
    if (isMounted) {
      loadContainers();
    }
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
    
    // Show confirmation dialog
    const actionText = action === 'start' ? '启动' : action === 'stop' ? '停止' : '重启';
    const container = containers.find(c => c.id === containerId);
    const containerName = container?.name || containerId;
    
    if (!confirm(`确定要${actionText}容器 "${containerName}" 吗？`)) {
      return;
    }
    
    try {
      actionError = '';
      // Set loading state for this specific action
      loadingActions = new Map(loadingActions).set(containerId, action);
      
      await containerApi.controlContainer({ containerId, action });
      
      // Wait before refreshing to allow container state to stabilize
      await new Promise(resolve => setTimeout(resolve, REFRESH_DELAY));
      
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
    } finally {
      // Clear loading state for this container
      const newLoadingActions = new Map(loadingActions);
      newLoadingActions.delete(containerId);
      loadingActions = newLoadingActions;
    }
  }
  
  async function handleRefresh() {
    refreshing = true;
    await loadContainers();
  }
</script>

<div class="home-container">
  <Header />
  
  <PageLayout title="容器列表">
    {#snippet actions()}
      <input
        type="text"
        class="filter-input"
        placeholder="按名称筛选..."
        value={filterText}
        oninput={handleFilterTextChange}
        aria-label="按容器名称筛选"
      />
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
      <select 
        class="group-mode-select" 
        value={groupMode} 
        onchange={handleGroupModeChange}
        aria-label="选择分组方式"
      >
        <option value="none">不分组</option>
        <option value="compose">按 Compose 分组</option>
        <option value="label">按标签分组</option>
        <option value="status-health">按状态和健康分组</option>
      </select>
      {#if groupMode === 'label' && availableLabelKeys.length > 0}
        <select 
          class="label-key-select" 
          value={selectedLabelKey} 
          onchange={handleLabelKeyChange}
          aria-label="选择标签"
        >
          {#each availableLabelKeys as labelKey}
            <option value={labelKey}>{labelKey}</option>
          {/each}
        </select>
      {/if}
      <select 
        class="sort-mode-select" 
        value={sortMode} 
        onchange={handleSortModeChange}
        aria-label="选择排序方式"
      >
        <option value="none">不排序</option>
        <option value="name">按名称</option>
        <option value="created">按创建时间</option>
        <option value="state-health">按状态和健康</option>
        <option value="compose">按 Compose 名称</option>
      </select>
      <button class="refresh-button" onclick={handleRefresh} disabled={refreshing}>
        <span class="refresh-icon" class:spinning={refreshing}>🔄</span>
        刷新
      </button>
    {/snippet}
    
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
        {@const { grouped, ungrouped } = groupContainersByCompose(sortedContainers)}
        
        <!-- Quick navigation sidebar -->
        {#if grouped.size > 0 || ungrouped.length > 0}
          <div class="quick-nav-sidebar">
            <div class="quick-nav-title">快速跳转</div>
            {#each Array.from(grouped.keys()) as projectName}
              <button class="quick-nav-item" onclick={() => scrollToGroup(projectName)}>
                {projectName}
              </button>
            {/each}
            {#if ungrouped.length > 0}
              <button class="quick-nav-item" onclick={() => scrollToGroup('_ungrouped_')}>
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
                onclick={() => toggleGroupCollapse(projectName)}
                aria-expanded={!collapsedGroups.has(projectName)}
                aria-label={`${projectName} compose group, ${projectContainers.length} containers`}
              >
                <span class="compose-icon">📚</span>
                <h3 class="compose-project-name">{projectName}</h3>
                <span class="compose-count">{projectContainers.length} 个容器</span>
                <span class="collapse-icon" aria-hidden="true">{collapsedGroups.has(projectName) ? '▶' : '▼'}</span>
              </button>
              {#if !collapsedGroups.has(projectName)}
              <div class="group-content" transition:slide={{ duration: 300 }}>
                <ContainerList
                  containers={projectContainers}
                  displayMode={displayMode}
                  onAction={handleAction}
                  loadingActions={loadingActions}
                />
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
              onclick={() => toggleGroupCollapse('_ungrouped_')}
              aria-expanded={!collapsedGroups.has('_ungrouped_')}
              aria-label={`独立容器 group, ${ungrouped.length} containers`}
            >
              <span class="compose-icon">📦</span>
              <h3 class="compose-project-name">独立容器</h3>
              <span class="compose-count">{ungrouped.length} 个容器</span>
              <span class="collapse-icon" aria-hidden="true">{collapsedGroups.has('_ungrouped_') ? '▶' : '▼'}</span>
            </button>
            {#if !collapsedGroups.has('_ungrouped_')}
            <div class="group-content" transition:slide={{ duration: 300 }}>
              <ContainerList
                containers={ungrouped}
                displayMode={displayMode}
                onAction={handleAction}
                loadingActions={loadingActions}
              />
            </div>
            {/if}
          </div>
        {/if}
      {:else if groupMode === 'label' && selectedLabelKey}
        <!-- Grouped by selected label -->
        {@const { grouped, ungrouped } = groupContainersByLabel(sortedContainers, selectedLabelKey)}
        
        <!-- Quick navigation sidebar -->
        {#if grouped.size > 0 || ungrouped.length > 0}
          <div class="quick-nav-sidebar">
            <div class="quick-nav-title">快速跳转</div>
            {#each Array.from(grouped.keys()) as groupName}
              <button class="quick-nav-item" onclick={() => scrollToGroup(`label-${groupName}`)}>
                {groupName}
              </button>
            {/each}
            {#if ungrouped.length > 0}
              <button class="quick-nav-item" onclick={() => scrollToGroup('_ungrouped_label_')}>
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
                onclick={() => toggleGroupCollapse(`label-${labelValue}`)}
                aria-expanded={!collapsedGroups.has(`label-${labelValue}`)}
                aria-label={`${labelValue} label group, ${labelContainers.length} containers`}
              >
                <span class="compose-icon">🏷️</span>
                <h3 class="compose-project-name">{labelValue}</h3>
                <span class="compose-count">{labelContainers.length} 个容器</span>
                <span class="collapse-icon" aria-hidden="true">{collapsedGroups.has(`label-${labelValue}`) ? '▶' : '▼'}</span>
              </button>
              {#if !collapsedGroups.has(`label-${labelValue}`)}
              <div class="group-content" transition:slide={{ duration: 300 }}>
                <ContainerList
                  containers={labelContainers}
                  displayMode={displayMode}
                  onAction={handleAction}
                  loadingActions={loadingActions}
                />
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
              onclick={() => toggleGroupCollapse('_ungrouped_label_')}
              aria-expanded={!collapsedGroups.has('_ungrouped_label_')}
              aria-label={`未分组容器, ${ungrouped.length} containers`}
            >
              <span class="compose-icon">📦</span>
              <h3 class="compose-project-name">未分组容器</h3>
              <span class="compose-count">{ungrouped.length} 个容器</span>
              <span class="collapse-icon" aria-hidden="true">{collapsedGroups.has('_ungrouped_label_') ? '▶' : '▼'}</span>
            </button>
            {#if !collapsedGroups.has('_ungrouped_label_')}
            <div class="group-content" transition:slide={{ duration: 300 }}>
              <ContainerList
                containers={ungrouped}
                displayMode={displayMode}
                onAction={handleAction}
                loadingActions={loadingActions}
              />
            </div>
            {/if}
          </div>
        {/if}
      {:else if groupMode === 'status-health'}
        <!-- Grouped by combined status and health -->
        {@const { grouped, ungrouped } = groupContainersByStatusHealth(sortedContainers)}
        
        <!-- Quick navigation sidebar -->
        {#if grouped.size > 0}
          <div class="quick-nav-sidebar">
            <div class="quick-nav-title">快速跳转</div>
            {#each Array.from(grouped.keys()) as groupKey}
              {@const { stateEmoji, healthEmoji, displayName } = getStatusHealthDisplay(groupKey)}
              <button class="quick-nav-item" onclick={() => scrollToGroup(`status-health-${groupKey}`)}>
                {stateEmoji}{healthEmoji} {displayName}
              </button>
            {/each}
          </div>
        {/if}
        
        {#if grouped.size > 0}
          {#each Array.from(grouped.entries()) as [groupKey, groupContainers] (groupKey)}
            {@const { stateEmoji, healthEmoji, displayName } = getStatusHealthDisplay(groupKey)}
            <div class="compose-group" id="group-status-health-{groupKey}">
              <button 
                class="compose-group-header" 
                class:compact={displayMode === 'compact'}
                onclick={() => toggleGroupCollapse(`status-health-${groupKey}`)}
                aria-expanded={!collapsedGroups.has(`status-health-${groupKey}`)}
                aria-label={`${displayName} group, ${groupContainers.length} containers`}
              >
                <span class="compose-icon">{stateEmoji}{healthEmoji}</span>
                <h3 class="compose-project-name">{displayName}</h3>
                <span class="compose-count">{groupContainers.length} 个容器</span>
                <span class="collapse-icon" aria-hidden="true">{collapsedGroups.has(`status-health-${groupKey}`) ? '▶' : '▼'}</span>
              </button>
              {#if !collapsedGroups.has(`status-health-${groupKey}`)}
              <div class="group-content" transition:slide={{ duration: 300 }}>
                <ContainerList
                  containers={groupContainers}
                  displayMode={displayMode}
                  onAction={handleAction}
                  loadingActions={loadingActions}
                />
              </div>
              {/if}
            </div>
          {/each}
        {/if}
      {:else}
        <!-- Ungrouped list -->
        <ContainerList
          containers={sortedContainers}
          displayMode={displayMode}
          onAction={handleAction}
          loadingActions={loadingActions}
        />
      {/if}
    {/if}
  </PageLayout>
</div>

<style>
  .home-container {
    min-height: 100vh;
    background: var(--color-background, #f5f5f4);
  }
  
  .mode-toggle,
  .refresh-button,
  .group-mode-select,
  .label-key-select,
  .sort-mode-select,
  .filter-input {
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
  
  .filter-input {
    min-width: 200px;
    cursor: text;
  }
  
  .filter-input::placeholder {
    color: var(--color-muted, #78716c);
    opacity: 0.7;
  }
  
  .mode-toggle:hover,
  .refresh-button:hover:not(:disabled),
  .group-mode-select:hover,
  .label-key-select:hover,
  .sort-mode-select:hover,
  .filter-input:hover {
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
    margin-bottom: 0;
    width: 100%;
    border: none;
    cursor: pointer;
    text-align: left;
    font-family: inherit;
    transition: background 0.2s, color 0.2s, border-radius 0.2s ease 0s, margin-bottom 0.2s ease 0s;
  }
  
  .compose-group-header[aria-expanded="false"] {
    border-radius: var(--radius, 0.25rem);
    margin-bottom: 0;
    transition: background 0.2s, color 0.2s, border-radius 0.2s ease 0.15s, margin-bottom 0.2s ease 0.15s;
  }

  .group-content {
    overflow: hidden;
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
  
  /* Quick navigation sidebar */
  .quick-nav-sidebar {
    /* Responsive width: clamp between 140px and 200px based on available space */
    --sidebar-width: clamp(140px, 15vw, 200px);
    --sidebar-gap: 2rem;
    
    position: sticky;
    top: 5rem;
    float: left;
    width: var(--sidebar-width);
    /* Position close to main content with minimal gap */
    margin-left: calc(-1 * var(--sidebar-width) - var(--sidebar-gap));
    margin-right: var(--sidebar-gap);
    background: var(--color-surface, #e7e5e4);
    border: 1px solid rgba(0, 0, 0, 0.1);
    border-radius: var(--radius, 0.25rem);
    padding: 0.5rem;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
    max-height: calc(100vh - 6rem);
    overflow-y: auto;
    z-index: 50;
  }
  
  .quick-nav-title {
    font-size: 0.8rem;
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
    padding: 0.4rem 0.5rem;
    margin-bottom: 0.15rem;
    border: none;
    background: transparent;
    color: var(--color-text, #0a0a0a);
    font-size: 0.75rem;
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
  
  /* Hide sidebar on medium screens (no space for floating sidebar) */
  @media (max-width: 1280px) {
    .quick-nav-sidebar {
      display: none;
    }
  }
  
  /* Extra large screens: slightly larger sidebar */
  @media (min-width: 1600px) {
    .quick-nav-sidebar {
      --sidebar-width: clamp(180px, 12vw, 220px);
      padding: 0.75rem;
    }
    
    .quick-nav-title {
      font-size: 0.85rem;
    }
    
    .quick-nav-item {
      font-size: 0.8rem;
      padding: 0.5rem 0.6rem;
    }
  }
</style>
