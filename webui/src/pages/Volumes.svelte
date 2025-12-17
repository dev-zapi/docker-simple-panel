<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import Header from '../components/Header.svelte';
  import { volumeApi, containerApi } from '../services/api';
  import type { Volume, Container } from '../types';
  
  let volumes: Volume[] = [];
  let containers: Container[] = [];
  let loading = true;
  let error = '';
  let refreshing = false;
  let deletingVolume: string | null = null;
  let volumeToDelete: string | null = null;
  let deleteTimeoutId: number | null = null;
  

  
  // Load volumes and containers
  onMount(() => {
    loadData();
  });
  
  onDestroy(() => {
  });
  
  async function loadData() {
    try {
      error = '';
      [volumes, containers] = await Promise.all([
        volumeApi.getVolumes(),
        containerApi.getContainers()
      ]);
    } catch (err) {
      error = '获取卷列表失败';
      console.error('Failed to load volumes:', err);
    } finally {
      loading = false;
      refreshing = false;
    }
  }
  
  async function handleRefresh() {
    refreshing = true;
    await loadData();
  }
  
  async function handleDeleteClick(volumeName: string) {
    // First click: set the volume to delete (confirmation state)
    if (volumeToDelete !== volumeName) {
      // Clear any existing timeout
      if (deleteTimeoutId !== null) {
        clearTimeout(deleteTimeoutId);
      }
      
      volumeToDelete = volumeName;
      // Reset confirmation after 3 seconds
      deleteTimeoutId = setTimeout(() => {
        if (volumeToDelete === volumeName) {
          volumeToDelete = null;
          deleteTimeoutId = null;
        }
      }, 3000) as unknown as number;
      return;
    }
    
    // Second click: show confirmation dialog before actually deleting
    // Clear the timeout since we're showing the dialog
    if (deleteTimeoutId !== null) {
      clearTimeout(deleteTimeoutId);
      deleteTimeoutId = null;
    }
    
    // Show native confirmation dialog
    const confirmed = confirm(`确定要删除卷 "${volumeName}" 吗？\n\n此操作无法撤销。`);
    
    if (!confirmed) {
      // User cancelled, reset state
      volumeToDelete = null;
      return;
    }
    
    // User confirmed, proceed with deletion
    try {
      deletingVolume = volumeName;
      error = '';
      await volumeApi.deleteVolume(volumeName);
      volumeToDelete = null;
      // Refresh the volume list
      await loadData();
    } catch (err) {
      error = `删除卷失败: ${err instanceof Error ? err.message : '未知错误'}`;
      console.error('Failed to delete volume:', err);
    } finally {
      deletingVolume = null;
    }
  }
  
  function cancelDelete() {
    // Clear the timeout when cancelling
    if (deleteTimeoutId !== null) {
      clearTimeout(deleteTimeoutId);
      deleteTimeoutId = null;
    }
    volumeToDelete = null;
  }
  
  // Get container names from IDs
  function getContainerNames(containerIds: string[]): string[] {
    return containerIds
      .map(id => {
        const container = containers.find(c => c.id === id);
        return container ? container.name : id;
      })
      .sort();
  }
  
  // Format date
  function formatDate(dateStr: string): string {
    if (!dateStr) return 'N/A';
    try {
      const date = new Date(dateStr);
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      });
    } catch {
      return dateStr;
    }
  }
</script>

<div class="volumes-container">
  <Header />
  
  <main class="main-content">
    <div class="content-header">
      <h2>卷列表</h2>
      <div class="header-actions">
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
    
    {#if loading}
      <div class="loading">
        <div class="spinner"></div>
        <p>加载中...</p>
      </div>
    {:else if volumes.length === 0}
      <div class="empty-state">
        <div class="empty-icon">💾</div>
        <p>暂无卷</p>
      </div>
    {:else}
      <div class="volume-list">
        {#each volumes as volume (volume.name)}
          <div class="volume-item">
            <div class="volume-info">
              <div class="volume-name">
                <span class="name-text">📦 {volume.name}</span>
              </div>
              <div class="volume-meta">
                <div class="meta-item">
                  <span class="meta-label">驱动:</span>
                  <span class="meta-value">{volume.driver}</span>
                </div>
                <div class="meta-item">
                  <span class="meta-label">作用域:</span>
                  <span class="meta-value">{volume.scope}</span>
                </div>
                <div class="meta-item">
                  <span class="meta-label">创建时间:</span>
                  <span class="meta-value">{formatDate(volume.created_at)}</span>
                </div>
              </div>
              <div class="volume-mountpoint">
                <span class="meta-label">挂载点:</span>
                <span class="path">{volume.mountpoint}</span>
              </div>
              {#if volume.containers.length > 0}
                <div class="volume-containers">
                  <span class="meta-label">关联容器 ({volume.containers.length}):</span>
                  <div class="container-tags">
                    {#each getContainerNames(volume.containers) as containerName}
                      <span class="container-tag">{containerName}</span>
                    {/each}
                  </div>
                </div>
              {:else}
                <div class="volume-containers empty">
                  <span class="meta-label">关联容器:</span>
                  <span class="empty-text">未被使用</span>
                </div>
              {/if}
              <div class="volume-actions">
                <button class="explore-button" on:click={() => window.location.hash = `/volumes/${volume.name}/explorer`}>
                  📂 浏览文件
                </button>
              </div>
            </div>
            <div class="volume-actions">
              {#if volumeToDelete === volume.name}
                <button 
                  class="delete-button confirm" 
                  on:click={() => handleDeleteClick(volume.name)}
                  disabled={deletingVolume === volume.name}
                >
                  {deletingVolume === volume.name ? '删除中...' : '确认删除'}
                </button>
                <button 
                  class="cancel-button" 
                  on:click={cancelDelete}
                  disabled={deletingVolume === volume.name}
                >
                  取消
                </button>
              {:else}
                <button 
                  class="delete-button" 
                  on:click={() => handleDeleteClick(volume.name)}
                  disabled={deletingVolume !== null}
                >
                  删除卷
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
  .volumes-container {
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
    position: sticky;
    top: 0;
    background: var(--color-background, #f5f5f4);
    z-index: 50;
    padding: 1rem 0;
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
  
  .refresh-button:hover:not(:disabled) {
    background: var(--color-background, #f5f5f4);
    border-color: var(--color-primary, #171717);
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
  
  .empty-state {
    text-align: center;
    padding: 4rem 2rem;
    color: var(--color-muted, #78716c);
  }
  
  .empty-icon {
    font-size: 4rem;
    margin-bottom: 1rem;
  }
  
  .volume-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .volume-item {
    background: var(--color-surface, #e7e5e4);
    border-radius: var(--radius, 0.25rem);
    padding: 1.25rem;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    transition: box-shadow 0.2s;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .volume-item:hover {
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  }
  
  .volume-info {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    flex: 1;
  }
  
  .volume-actions {
    display: flex;
    gap: 0.75rem;
    justify-content: flex-end;
    padding-top: 0.5rem;
    border-top: 1px solid rgba(0, 0, 0, 0.1);
  }
  
  .delete-button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: var(--color-surface, #e7e5e4);
    border: 1px solid rgba(153, 27, 27, 0.3);
    color: var(--color-error, #991b1b);
    padding: 0.5rem 1rem;
    border-radius: var(--radius, 0.25rem);
    cursor: pointer;
    font-size: 0.9rem;
    transition: all 0.2s;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .delete-button:hover:not(:disabled) {
    background: rgba(153, 27, 27, 0.1);
    border-color: var(--color-error, #991b1b);
  }
  
  .delete-button.confirm {
    background: var(--color-error, #991b1b);
    color: var(--color-background, #f5f5f4);
    border-color: var(--color-error, #991b1b);
  }
  
  .delete-button.confirm:hover:not(:disabled) {
    background: #7f1d1d;
    border-color: #7f1d1d;
  }
  
  .delete-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  .cancel-button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: var(--color-surface, #e7e5e4);
    border: 1px solid rgba(0, 0, 0, 0.1);
    padding: 0.5rem 1rem;
    border-radius: var(--radius, 0.25rem);
    cursor: pointer;
    font-size: 0.9rem;
    transition: all 0.2s;
    color: var(--color-text, #0a0a0a);
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .cancel-button:hover:not(:disabled) {
    background: var(--color-background, #f5f5f4);
    border-color: var(--color-primary, #171717);
  }
  
  .cancel-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  .volume-name {
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--color-text, #0a0a0a);
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-family: var(--font-heading, "Playfair Display", serif);
  }
  
  .volume-meta {
    display: flex;
    flex-wrap: wrap;
    gap: 1.5rem;
    font-size: 0.9rem;
  }
  
  .meta-item {
    display: flex;
    gap: 0.5rem;
  }
  
  .meta-label {
    font-weight: 600;
    color: var(--color-secondary, #525252);
  }
  
  .meta-value {
    color: var(--color-text, #0a0a0a);
  }
  
  .volume-mountpoint {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    font-size: 0.9rem;
  }
  
  .path {
    font-family: monospace;
    color: var(--color-muted, #78716c);
    background: rgba(0, 0, 0, 0.05);
    padding: 0.25rem 0.5rem;
    border-radius: var(--radius, 0.25rem);
    overflow-wrap: break-word;
  }
  
  .volume-containers {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    font-size: 0.9rem;
  }
  
  .volume-containers.empty {
    flex-direction: row;
    gap: 0.5rem;
  }
  
  .empty-text {
    color: var(--color-muted, #78716c);
    font-style: italic;
  }
  
  .container-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }
  
  .container-tag {
    background: var(--color-primary, #171717);
    color: var(--color-background, #f5f5f4);
    padding: 0.25rem 0.75rem;
    border-radius: var(--radius, 0.25rem);
    font-size: 0.85rem;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .volume-actions {
    display: flex;
    gap: 0.75rem;
    margin-top: 0.75rem;
  }
  
  .explore-button {
    background: var(--color-primary, #171717);
    color: var(--color-background, #f5f5f4);
    border: none;
    padding: 0.5rem 1rem;
    border-radius: var(--radius, 0.25rem);
    cursor: pointer;
    font-size: 0.9rem;
    transition: all 0.2s;
    font-family: var(--font-body, "Merriweather", serif);
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  
  .explore-button:hover {
    background: var(--color-secondary, #525252);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
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
    
    .volume-meta {
      flex-direction: column;
      gap: 0.5rem;
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

    .floating-header .refresh-button {
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
    .floating-header .refresh-button {
      font-size: 0;
      gap: 0;
    }

    .floating-header .refresh-icon {
      font-size: 1rem;
    }
  }
</style>
