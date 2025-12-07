<script lang="ts">
  import { onMount } from 'svelte';
  import Header from '../components/Header.svelte';
  import { configApi } from '../services/api';
  import type { SystemConfig } from '../types';
  
  let config: SystemConfig | null = null;
  let loading = true;
  let saving = false;
  let error = '';
  let successMessage = '';
  
  // Form fields
  let dockerSocket = '';
  let disableRegistration = false;
  let logLevel = 'info';
  let volumeExplorerImage = 'busybox:latest';
  
  const logLevelOptions = [
    { value: 'error', label: '错误 (Error)' },
    { value: 'warn', label: '警告 (Warn)' },
    { value: 'info', label: '信息 (Info)' },
    { value: 'debug', label: '调试 (Debug)' }
  ];
  
  async function loadConfig() {
    try {
      error = '';
      config = await configApi.getConfig();
      dockerSocket = config.docker_socket;
      disableRegistration = config.disable_registration;
      logLevel = config.log_level || 'info';
      volumeExplorerImage = config.volume_explorer_image || 'busybox:latest';
    } catch (err) {
      error = '获取配置失败';
      console.error('Failed to load config:', err);
    } finally {
      loading = false;
    }
  }
  
  async function handleSave() {
    saving = true;
    error = '';
    successMessage = '';
    
    try {
      const updatedConfig = await configApi.updateConfig({
        docker_socket: dockerSocket,
        disable_registration: disableRegistration,
        log_level: logLevel,
        volume_explorer_image: volumeExplorerImage
      });
      
      config = updatedConfig;
      successMessage = '配置保存成功';
      
      // Clear success message after 3 seconds
      setTimeout(() => {
        successMessage = '';
      }, 3000);
    } catch (err) {
      error = '保存配置失败：' + (err instanceof Error ? err.message : '未知错误');
      console.error('Failed to save config:', err);
    } finally {
      saving = false;
    }
  }
  
  function handleReset() {
    if (config) {
      dockerSocket = config.docker_socket;
      disableRegistration = config.disable_registration;
      logLevel = config.log_level || 'info';
      volumeExplorerImage = config.volume_explorer_image || 'busybox:latest';
    }
    error = '';
    successMessage = '';
  }
  
  onMount(() => {
    loadConfig();
  });
</script>

<div class="settings-container">
  <Header />
  
  <main class="main-content">
    <div class="content-header">
      <h2>系统设置</h2>
    </div>
    
    {#if error}
      <div class="error-banner">
        {error}
      </div>
    {/if}
    
    {#if successMessage}
      <div class="success-banner">
        {successMessage}
      </div>
    {/if}
    
    {#if loading}
      <div class="loading">
        <div class="spinner"></div>
        <p>加载中...</p>
      </div>
    {:else}
      <div class="settings-form">
        <div class="form-section">
          <h3>Docker 配置</h3>
          
          <div class="form-group">
            <label for="dockerSocket">Docker Socket 路径</label>
            <input
              type="text"
              id="dockerSocket"
              bind:value={dockerSocket}
              placeholder="/var/run/docker.sock"
              disabled={saving}
            />
            <p class="form-help">Docker 守护进程的 Unix socket 路径</p>
          </div>
        </div>
        
        <div class="form-section">
          <h3>用户注册</h3>
          
          <div class="form-group checkbox-group">
            <label class="checkbox-label">
              <input
                type="checkbox"
                bind:checked={disableRegistration}
                disabled={saving}
              />
              <span class="checkbox-text">禁用用户注册</span>
            </label>
            <p class="form-help">启用后，新用户将无法通过注册页面创建账户</p>
          </div>
        </div>
        
        <div class="form-section">
          <h3>日志配置</h3>
          
          <div class="form-group">
            <label for="logLevel">日志级别</label>
            <select id="logLevel" bind:value={logLevel} disabled={saving}>
              {#each logLevelOptions as option}
                <option value={option.value}>{option.label}</option>
              {/each}
            </select>
            <p class="form-help">设置系统日志的详细程度</p>
          </div>
        </div>
        
        <div class="form-section">
          <h3>卷管理配置</h3>
          
          <div class="form-group">
            <label for="volumeExplorerImage">卷浏览器镜像</label>
            <input
              type="text"
              id="volumeExplorerImage"
              bind:value={volumeExplorerImage}
              placeholder="busybox:latest"
              disabled={saving}
            />
            <p class="form-help">用于浏览卷文件的临时容器镜像（推荐：busybox:latest 或 alpine:latest）</p>
          </div>
        </div>
        
        <div class="form-actions">
          <button class="btn-secondary" on:click={handleReset} disabled={saving}>
            重置
          </button>
          <button class="btn-primary" on:click={handleSave} disabled={saving}>
            {saving ? '保存中...' : '保存配置'}
          </button>
        </div>
      </div>
    {/if}
  </main>
</div>

<style>
  .settings-container {
    min-height: 100vh;
    background: var(--color-background, #f5f5f4);
  }
  
  .main-content {
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem;
  }
  
  .content-header {
    margin-bottom: 1.5rem;
  }
  
  .content-header h2 {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--color-text, #0a0a0a);
    margin: 0;
    font-family: var(--font-heading, "Playfair Display", serif);
  }
  
  .error-banner {
    background: rgba(153, 27, 27, 0.1);
    border: 1px solid var(--color-error, #991b1b);
    color: var(--color-error, #991b1b);
    padding: 1rem;
    border-radius: var(--radius, 0.25rem);
    margin-bottom: 1.5rem;
  }
  
  .success-banner {
    background: rgba(21, 128, 61, 0.1);
    border: 1px solid var(--color-success, #15803d);
    color: var(--color-success, #15803d);
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
  
  .settings-form {
    background: var(--color-surface, #e7e5e4);
    border-radius: var(--radius, 0.25rem);
    padding: 1.5rem;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  }
  
  .form-section {
    margin-bottom: 2rem;
    padding-bottom: 1.5rem;
    border-bottom: 1px solid rgba(0, 0, 0, 0.1);
  }
  
  .form-section:last-of-type {
    border-bottom: none;
    margin-bottom: 0;
    padding-bottom: 0;
  }
  
  .form-section h3 {
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--color-text, #0a0a0a);
    margin: 0 0 1rem 0;
    font-family: var(--font-heading, "Playfair Display", serif);
  }
  
  .form-group {
    margin-bottom: 1rem;
  }
  
  .form-group label {
    display: block;
    font-weight: 500;
    color: var(--color-text, #0a0a0a);
    margin-bottom: 0.5rem;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .form-group input[type="text"],
  .form-group select {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid rgba(0, 0, 0, 0.2);
    border-radius: var(--radius, 0.25rem);
    font-size: 1rem;
    background: var(--color-background, #f5f5f4);
    color: var(--color-text, #0a0a0a);
    font-family: var(--font-body, "Merriweather", serif);
    transition: border-color 0.2s, box-shadow 0.2s;
  }
  
  .form-group input[type="text"]:focus,
  .form-group select:focus {
    outline: none;
    border-color: var(--color-primary, #171717);
    box-shadow: 0 0 0 2px rgba(23, 23, 23, 0.2);
  }
  
  .form-group input:disabled,
  .form-group select:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  .form-help {
    font-size: 0.875rem;
    color: var(--color-muted, #78716c);
    margin: 0.5rem 0 0 0;
  }
  
  .checkbox-group {
    display: flex;
    flex-direction: column;
  }
  
  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    cursor: pointer;
  }
  
  .checkbox-label input[type="checkbox"] {
    width: 1.25rem;
    height: 1.25rem;
    accent-color: var(--color-primary, #171717);
    cursor: pointer;
  }
  
  .checkbox-text {
    font-weight: 500;
    color: var(--color-text, #0a0a0a);
  }
  
  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 2rem;
    padding-top: 1.5rem;
    border-top: 1px solid rgba(0, 0, 0, 0.1);
  }
  
  .btn-primary,
  .btn-secondary {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: var(--radius, 0.25rem);
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .btn-primary {
    background: var(--color-primary, #171717);
    color: white;
  }
  
  .btn-primary:hover:not(:disabled) {
    background: var(--color-secondary, #525252);
  }
  
  .btn-primary:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  .btn-secondary {
    background: transparent;
    color: var(--color-text, #0a0a0a);
    border: 1px solid rgba(0, 0, 0, 0.2);
  }
  
  .btn-secondary:hover:not(:disabled) {
    background: rgba(0, 0, 0, 0.05);
  }
  
  .btn-secondary:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
</style>
