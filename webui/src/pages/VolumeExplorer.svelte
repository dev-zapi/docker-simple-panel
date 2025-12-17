<script lang="ts">
  import { onMount } from 'svelte';
  import { push } from 'svelte-spa-router';
  import Header from '../components/Header.svelte';
  import { volumeApi } from '../services/api';
  import type { VolumeFileInfo, VolumeFileContent } from '../types';
  
  export let params: { name?: string } = {};
  
  let volumeName = params.name || '';
  let currentPath = '/';
  let files: VolumeFileInfo[] = [];
  let selectedFile: VolumeFileContent | null = null;
  let loading = true;
  let loadingFile = false;
  let error = '';
  let fileError = '';
  let deletingVolume = false;
  let showDeleteConfirm = false;
  let deleteTimeoutId: number | null = null;
  

  
  onMount(() => {
    if (!volumeName) {
      push('/volumes');
      return;
    }
    
    loadFiles();
  });
  
  async function loadFiles() {
    try {
      loading = true;
      error = '';
      files = await volumeApi.exploreVolumeFiles(volumeName, currentPath);
      selectedFile = null;
      fileError = '';
    } catch (err) {
      error = '加载文件列表失败：' + (err instanceof Error ? err.message : '未知错误');
      console.error('Failed to load files:', err);
    } finally {
      loading = false;
    }
  }
  
  async function handleNavigate(file: VolumeFileInfo) {
    if (file.is_directory) {
      currentPath = file.path;
      await loadFiles();
    } else {
      await loadFileContent(file);
    }
  }
  
  async function loadFileContent(file: VolumeFileInfo) {
    try {
      loadingFile = true;
      fileError = '';
      selectedFile = await volumeApi.readVolumeFile(volumeName, file.path);
    } catch (err) {
      fileError = '读取文件失败：' + (err instanceof Error ? err.message : '未知错误');
      console.error('Failed to read file:', err);
    } finally {
      loadingFile = false;
    }
  }
  
  function handleGoUp() {
    if (currentPath === '/') return;
    
    const parts = currentPath.split('/').filter(p => p);
    parts.pop();
    currentPath = parts.length === 0 ? '/' : '/' + parts.join('/');
    loadFiles();
  }
  
  function handleGoToRoot() {
    currentPath = '/';
    loadFiles();
  }
  
  function formatFileSize(bytes: number): string {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
  }
  
  function closeFileViewer() {
    selectedFile = null;
    fileError = '';
  }
  
  async function handleDeleteVolumeClick() {
    // First click: show confirmation
    if (!showDeleteConfirm) {
      showDeleteConfirm = true;
      // Reset confirmation after 3 seconds
      deleteTimeoutId = window.setTimeout(() => {
        showDeleteConfirm = false;
        deleteTimeoutId = null;
      }, 3000);
      return;
    }
    
    // Second click: show native confirmation dialog
    if (deleteTimeoutId !== null) {
      clearTimeout(deleteTimeoutId);
      deleteTimeoutId = null;
    }
    
    const confirmed = confirm(`确定要删除卷 "${volumeName}" 吗？\n\n此操作无法撤销。`);
    
    if (!confirmed) {
      showDeleteConfirm = false;
      return;
    }
    
    // User confirmed, proceed with deletion
    try {
      deletingVolume = true;
      error = '';
      await volumeApi.deleteVolume(volumeName);
      // Navigate back to volumes page
      push('/volumes');
    } catch (err) {
      error = `删除卷失败: ${err instanceof Error ? err.message : '未知错误'}`;
      console.error('Failed to delete volume:', err);
      showDeleteConfirm = false;
    } finally {
      deletingVolume = false;
    }
  }
  
  function cancelDeleteVolume() {
    if (deleteTimeoutId !== null) {
      clearTimeout(deleteTimeoutId);
      deleteTimeoutId = null;
    }
    showDeleteConfirm = false;
  }
</script>

<div class="explorer-container">
  <Header />
  
  <main class="main-content">
    <div class="content-header">
      <div class="header-top">
        <h2>📦 {volumeName}</h2>
        <div class="header-actions">
          {#if showDeleteConfirm}
            <button 
              class="delete-volume-button confirm" 
              on:click={handleDeleteVolumeClick}
              disabled={deletingVolume}
            >
              {deletingVolume ? '删除中...' : '确认删除卷'}
            </button>
            <button 
              class="cancel-button" 
              on:click={cancelDeleteVolume}
              disabled={deletingVolume}
            >
              取消
            </button>
          {:else}
            <button 
              class="delete-volume-button" 
              on:click={handleDeleteVolumeClick}
              disabled={deletingVolume}
            >
              🗑️ 删除卷
            </button>
          {/if}
          <button class="back-button" on:click={() => push('/volumes')}>
            ← 返回卷列表
          </button>
        </div>
      </div>
      <div class="breadcrumb">
        <button class="breadcrumb-btn" on:click={handleGoToRoot} title="根目录">🏠</button>
        {#if currentPath !== '/'}
          <span class="separator">/</span>
          <span class="path-text">{currentPath}</span>
        {/if}
      </div>
    </div>
    
    {#if error}
      <div class="error-banner">
        {error}
      </div>
    {/if}
    
    <div class="explorer-layout">
      <!-- File list panel -->
      <div class="file-list-panel">
        {#if loading}
          <div class="loading">
            <div class="spinner"></div>
            <p>加载中...</p>
          </div>
        {:else if files.length === 0}
          <div class="empty-state">
            <div class="empty-icon">📁</div>
            <p>此目录为空</p>
          </div>
        {:else}
          <div class="file-list">
            {#if currentPath !== '/'}
              <button class="file-item directory" on:click={handleGoUp}>
                <span class="file-icon">📂</span>
                <span class="file-name">..</span>
                <span class="file-meta">返回上级</span>
              </button>
            {/if}
            
            {#each files as file (file.path)}
              <button class="file-item" class:directory={file.is_directory} on:click={() => handleNavigate(file)}>
                <span class="file-icon">{file.is_directory ? '📁' : '📄'}</span>
                <div class="file-info">
                  <span class="file-name">{file.name}</span>
                  <div class="file-meta">
                    <span class="file-mode">{file.mode}</span>
                    {#if !file.is_directory}
                      <span class="file-size">{formatFileSize(file.size)}</span>
                    {/if}
                    <span class="file-time">{file.mod_time}</span>
                  </div>
                </div>
              </button>
            {/each}
          </div>
        {/if}
      </div>
      
      <!-- File viewer panel -->
      <div class="file-viewer-panel" class:visible={selectedFile !== null}>
        {#if selectedFile}
          <div class="viewer-header">
            <div class="viewer-title">
              <span class="viewer-icon">📄</span>
              <span class="viewer-path">{selectedFile.path}</span>
              <span class="viewer-size">({formatFileSize(selectedFile.size)})</span>
            </div>
            <button class="close-btn" on:click={closeFileViewer}>✕</button>
          </div>
          <div class="viewer-content">
            {#if loadingFile}
              <div class="loading">
                <div class="spinner"></div>
                <p>读取文件中...</p>
              </div>
            {:else if fileError}
              <div class="file-error">
                {fileError}
              </div>
            {:else}
              <pre class="code-content">{selectedFile.content}</pre>
            {/if}
          </div>
        {:else}
          <div class="viewer-placeholder">
            <div class="placeholder-icon">👈</div>
            <p>选择一个文件查看内容</p>
          </div>
        {/if}
      </div>
    </div>
  </main>
</div>

<style>
  .explorer-container {
    min-height: 100vh;
    background: var(--color-background, #f5f5f4);
  }
  
  .main-content {
    max-width: 1400px;
    margin: 0 auto;
    padding: 2rem;
  }
  
  .content-header {
    margin-bottom: 1.5rem;
    position: sticky;
    top: 0;
    background: var(--color-background-blur, rgba(245, 245, 244, 0.8));
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
    z-index: 50;
    padding: 1rem 0;
  }
  
  .header-top {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }
  
  .header-actions {
    display: flex;
    gap: 0.75rem;
    align-items: center;
    flex-wrap: wrap;
  }
  
  .content-header h2 {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--color-text, #0a0a0a);
    margin: 0;
    font-family: var(--font-heading, "Playfair Display", serif);
  }
  
  .delete-volume-button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: var(--color-surface, #e7e5e4);
    border: 1px solid rgba(153, 27, 27, 0.3);
    color: var(--color-error, #991b1b);
    padding: 0.5rem 1rem;
    border-radius: var(--radius, 0.25rem);
    cursor: pointer;
    font-size: 0.95rem;
    transition: all 0.2s;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .delete-volume-button:hover:not(:disabled) {
    background: rgba(153, 27, 27, 0.1);
    border-color: var(--color-error, #991b1b);
  }
  
  .delete-volume-button.confirm {
    background: var(--color-error, #991b1b);
    color: var(--color-background, #f5f5f4);
    border-color: var(--color-error, #991b1b);
  }
  
  .delete-volume-button.confirm:hover:not(:disabled) {
    background: #7f1d1d;
    border-color: #7f1d1d;
  }
  
  .delete-volume-button:disabled {
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
    font-size: 0.95rem;
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
  
  .back-button {
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
  
  .back-button:hover {
    background: var(--color-background, #f5f5f4);
    border-color: var(--color-primary, #171717);
  }
  
  .breadcrumb {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.95rem;
    color: var(--color-muted, #78716c);
    font-family: monospace;
  }
  
  .breadcrumb-btn {
    background: transparent;
    border: none;
    padding: 0.25rem;
    cursor: pointer;
    font-size: 1.2rem;
    transition: transform 0.2s;
  }
  
  .breadcrumb-btn:hover {
    transform: scale(1.2);
  }
  
  .separator {
    color: var(--color-muted, #78716c);
  }
  
  .path-text {
    font-weight: 500;
  }
  
  .error-banner {
    background: rgba(153, 27, 27, 0.1);
    border: 1px solid var(--color-error, #991b1b);
    color: var(--color-error, #991b1b);
    padding: 1rem;
    border-radius: var(--radius, 0.25rem);
    margin-bottom: 1.5rem;
  }
  
  .explorer-layout {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
    height: calc(100vh - 250px);
    min-height: 500px;
  }
  
  .file-list-panel {
    background: var(--color-surface, #e7e5e4);
    border-radius: var(--radius, 0.25rem);
    padding: 1rem;
    overflow-y: auto;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  }
  
  .file-viewer-panel {
    background: var(--color-surface, #e7e5e4);
    border-radius: var(--radius, 0.25rem);
    overflow: hidden;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    display: flex;
    flex-direction: column;
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
    border: 4px solid var(--color-background, #f5f5f4);
    border-top: 4px solid var(--color-primary, #171717);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 1rem;
  }
  
  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
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
  
  .file-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .file-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem;
    background: var(--color-background, #f5f5f4);
    border: 1px solid rgba(0, 0, 0, 0.1);
    border-radius: var(--radius, 0.25rem);
    cursor: pointer;
    transition: all 0.2s;
    text-align: left;
    width: 100%;
  }
  
  .file-item:hover {
    background: var(--color-surface, #e7e5e4);
    border-color: var(--color-primary, #171717);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
  
  .file-item.directory {
    font-weight: 500;
  }
  
  .file-icon {
    font-size: 1.5rem;
    flex-shrink: 0;
  }
  
  .file-info {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }
  
  .file-name {
    font-weight: 500;
    color: var(--color-text, #0a0a0a);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  
  .file-meta {
    display: flex;
    gap: 1rem;
    font-size: 0.85rem;
    color: var(--color-muted, #78716c);
    font-family: monospace;
  }
  
  .viewer-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    border-bottom: 1px solid rgba(0, 0, 0, 0.1);
    background: var(--color-background, #f5f5f4);
  }
  
  .viewer-title {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.95rem;
    color: var(--color-text, #0a0a0a);
    font-family: monospace;
  }
  
  .viewer-icon {
    font-size: 1.2rem;
  }
  
  .viewer-path {
    font-weight: 500;
  }
  
  .viewer-size {
    color: var(--color-muted, #78716c);
  }
  
  .close-btn {
    background: transparent;
    border: none;
    padding: 0.25rem 0.5rem;
    cursor: pointer;
    font-size: 1.5rem;
    color: var(--color-muted, #78716c);
    transition: color 0.2s;
  }
  
  .close-btn:hover {
    color: var(--color-error, #991b1b);
  }
  
  .viewer-content {
    flex: 1;
    overflow: auto;
    padding: 1rem;
    background: var(--color-background, #f5f5f4);
  }
  
  .file-error {
    color: var(--color-error, #991b1b);
    padding: 1rem;
    background: rgba(153, 27, 27, 0.1);
    border: 1px solid var(--color-error, #991b1b);
    border-radius: var(--radius, 0.25rem);
  }
  
  .code-content {
    margin: 0;
    padding: 1rem;
    background: #1e1e1e;
    color: #d4d4d4;
    border-radius: var(--radius, 0.25rem);
    font-family: 'Courier New', monospace;
    font-size: 0.9rem;
    line-height: 1.5;
    overflow-x: auto;
    white-space: pre;
  }
  
  .viewer-placeholder {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: var(--color-muted, #78716c);
  }
  
  .placeholder-icon {
    font-size: 4rem;
    margin-bottom: 1rem;
  }
  
  @media (max-width: 1024px) {
    .explorer-layout {
      grid-template-columns: 1fr;
      height: auto;
    }
    
    .file-viewer-panel {
      min-height: 400px;
    }
  }
  
  @media (max-width: 640px) {
    .main-content {
      padding: 1rem;
    }
    
    .header-top {
      flex-direction: column;
      gap: 1rem;
      align-items: flex-start;
    }
  }
</style>
