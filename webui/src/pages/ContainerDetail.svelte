<script lang="ts">
  import { onMount } from 'svelte';
  import { push } from 'svelte-spa-router';
  import Header from '../components/Header.svelte';
  import { containerApi } from '../services/api';
  import type { Container } from '../types';
  
  export let params: { id: string };
  
  let container: Container | null = null;
  let loading = true;
  let error = '';
  
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
  
  onMount(async () => {
    await loadContainerDetails();
  });
  
  async function loadContainerDetails() {
    try {
      error = '';
      loading = true;
      container = await containerApi.getContainer(params.id);
    } catch (err) {
      error = '获取容器详情失败';
      console.error('Failed to load container details:', err);
    } finally {
      loading = false;
    }
  }
  
  function goBack() {
    push('/');
  }
  
  function formatRestartPolicy(policy?: { name: string; maximum_retry_count?: number }): string {
    if (!policy || !policy.name) return '无';
    
    switch (policy.name) {
      case 'no':
        return '不重启';
      case 'always':
        return '总是重启';
      case 'unless-stopped':
        return '除非停止';
      case 'on-failure':
        return policy.maximum_retry_count 
          ? `失败时重启 (最多 ${policy.maximum_retry_count} 次)`
          : '失败时重启';
      default:
        return policy.name;
    }
  }
  
  function formatDate(timestamp: number): string {
    const date = new Date(timestamp * 1000);
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
  }
</script>

<div class="detail-container">
  <Header />
  
  <main class="main-content">
    <div class="content-header">
      <button class="back-button" on:click={goBack}>
        ← 返回列表
      </button>
      <h2>容器详情</h2>
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
    {:else if container}
      <div class="detail-content">
        <!-- Basic Information -->
        <section class="detail-section">
          <h3 class="section-title">📦 基本信息</h3>
          <div class="info-grid">
            <div class="info-item">
              <span class="info-label">容器名称</span>
              <span class="info-value">{container.name}</span>
            </div>
            <div class="info-item">
              <span class="info-label">容器ID</span>
              <span class="info-value mono">{container.id}</span>
            </div>
            <div class="info-item">
              <span class="info-label">状态</span>
              <span class="info-value">
                <span class="status-emoji">{stateEmojis[container.state] || '⚪'}</span>
                {container.status}
              </span>
            </div>
            {#if container.health && container.health !== 'none'}
              <div class="info-item">
                <span class="info-label">健康状态</span>
                <span class="info-value">
                  <span class="health-emoji">{healthEmojis[container.health]}</span>
                  {container.health}
                </span>
              </div>
            {/if}
            <div class="info-item">
              <span class="info-label">镜像地址</span>
              <span class="info-value mono">{container.image}</span>
            </div>
            {#if container.hostname}
              <div class="info-item">
                <span class="info-label">主机名</span>
                <span class="info-value mono">{container.hostname}</span>
              </div>
            {/if}
            <div class="info-item">
              <span class="info-label">创建时间</span>
              <span class="info-value">{formatDate(container.created)}</span>
            </div>
            {#if container.compose_project}
              <div class="info-item">
                <span class="info-label">所属 Compose</span>
                <span class="info-value">{container.compose_project}</span>
              </div>
            {/if}
            {#if container.compose_service}
              <div class="info-item">
                <span class="info-label">Compose 服务</span>
                <span class="info-value">{container.compose_service}</span>
              </div>
            {/if}
            <div class="info-item">
              <span class="info-label">重启策略</span>
              <span class="info-value">{formatRestartPolicy(container.restart_policy)}</span>
            </div>
          </div>
        </section>
        
        <!-- Network Information -->
        {#if container.networks && Object.keys(container.networks).length > 0}
          <section class="detail-section">
            <h3 class="section-title">🌐 网络信息</h3>
            <div class="networks-list">
              {#each Object.entries(container.networks) as [networkName, networkInfo]}
                <div class="network-item">
                  <div class="network-header">
                    <span class="network-name">{networkName}</span>
                  </div>
                  <div class="network-details">
                    {#if networkInfo.ip_address}
                      <div class="network-detail">
                        <span class="detail-label">IP 地址:</span>
                        <span class="detail-value mono">{networkInfo.ip_address}</span>
                      </div>
                    {/if}
                    {#if networkInfo.gateway}
                      <div class="network-detail">
                        <span class="detail-label">网关:</span>
                        <span class="detail-value mono">{networkInfo.gateway}</span>
                      </div>
                    {/if}
                    {#if networkInfo.mac_address}
                      <div class="network-detail">
                        <span class="detail-label">MAC 地址:</span>
                        <span class="detail-value mono">{networkInfo.mac_address}</span>
                      </div>
                    {/if}
                    <div class="network-detail">
                      <span class="detail-label">网络 ID:</span>
                      <span class="detail-value mono small">{networkInfo.network_id}</span>
                    </div>
                  </div>
                </div>
              {/each}
            </div>
          </section>
        {/if}
        
        <!-- Port Mappings -->
        {#if container.ports && container.ports.length > 0}
          <section class="detail-section">
            <h3 class="section-title">🔌 端口映射</h3>
            <div class="ports-list">
              {#each container.ports as port}
                <div class="port-item">
                  {#if port.host_port}
                    <span class="port-mapping">
                      <span class="host-port">{port.host_ip || '0.0.0.0'}:{port.host_port}</span>
                      <span class="port-arrow">→</span>
                      <span class="container-port">{port.container_port}</span>
                    </span>
                  {:else}
                    <span class="port-mapping">
                      <span class="container-port">{port.container_port}</span>
                      <span class="port-note">(已暴露)</span>
                    </span>
                  {/if}
                </div>
              {/each}
            </div>
          </section>
        {/if}
        
        <!-- Mount Information -->
        {#if container.mounts && container.mounts.length > 0}
          <section class="detail-section">
            <h3 class="section-title">💾 存储映射</h3>
            <div class="mounts-list">
              {#each container.mounts as mount}
                <div class="mount-item">
                  <div class="mount-type-badge" class:bind-type={mount.type === 'bind'} class:volume-type={mount.type === 'volume'}>
                    {mount.type}
                  </div>
                  <div class="mount-details">
                    <div class="mount-path">
                      <span class="path-label">源:</span>
                      <span class="path-value mono">{mount.source}</span>
                    </div>
                    <div class="mount-path">
                      <span class="path-label">目标:</span>
                      <span class="path-value mono">{mount.destination}</span>
                    </div>
                    <div class="mount-meta">
                      <span class="mount-rw" class:readonly={!mount.rw}>
                        {mount.rw ? '读写' : '只读'}
                      </span>
                      {#if mount.mode}
                        <span class="mount-mode mono">{mount.mode}</span>
                      {/if}
                    </div>
                  </div>
                </div>
              {/each}
            </div>
          </section>
        {/if}
        
        <!-- Environment Variables -->
        {#if container.env && container.env.length > 0}
          <section class="detail-section">
            <h3 class="section-title">🔧 环境变量</h3>
            <div class="env-list">
              {#each container.env as envVar}
                {@const [key, ...valueParts] = envVar.split('=')}
                {@const value = valueParts.join('=')}
                <div class="env-item">
                  <span class="env-key mono">{key}</span>
                  <span class="env-equals">=</span>
                  <span class="env-value mono">{value || ''}</span>
                </div>
              {/each}
            </div>
          </section>
        {/if}
        
        <!-- Labels -->
        {#if container.labels && Object.keys(container.labels).length > 0}
          <section class="detail-section">
            <h3 class="section-title">🏷️ 标签</h3>
            <div class="labels-list">
              {#each Object.entries(container.labels) as [key, value]}
                <div class="label-item">
                  <span class="label-key mono">{key}</span>
                  <span class="label-equals">=</span>
                  <span class="label-value mono">{value}</span>
                </div>
              {/each}
            </div>
          </section>
        {/if}
      </div>
    {:else}
      <div class="empty-state">
        <div class="empty-icon">📦</div>
        <p>容器不存在</p>
      </div>
    {/if}
  </main>
</div>

<style>
  .detail-container {
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
    align-items: center;
    gap: 1rem;
    margin-bottom: 1.5rem;
    position: sticky;
    top: 0;
    background: var(--color-background-blur, rgba(245, 245, 244, 0.8));
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
    z-index: 50;
    padding: 1rem 0;
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
  
  .empty-state {
    text-align: center;
    padding: 4rem 2rem;
    color: var(--color-muted, #78716c);
  }
  
  .empty-icon {
    font-size: 4rem;
    margin-bottom: 1rem;
  }
  
  .detail-content {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
  
  .detail-section {
    background: var(--color-surface, #e7e5e4);
    border-radius: var(--radius, 0.25rem);
    padding: 1.5rem;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  }
  
  .section-title {
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--color-text, #0a0a0a);
    margin: 0 0 1rem 0;
    font-family: var(--font-heading, "Playfair Display", serif);
  }
  
  .info-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1rem;
  }
  
  .info-item {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }
  
  .info-label {
    font-size: 0.85rem;
    font-weight: 500;
    color: var(--color-muted, #78716c);
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }
  
  .info-value {
    font-size: 1rem;
    color: var(--color-text, #0a0a0a);
    word-break: break-word;
  }
  
  .mono {
    font-family: monospace;
    font-size: 0.95em;
  }
  
  .small {
    font-size: 0.85em;
  }
  
  .status-emoji, .health-emoji {
    margin-right: 0.25rem;
  }
  
  .networks-list, .mounts-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .network-item {
    background: var(--color-background, #f5f5f4);
    border-radius: var(--radius, 0.25rem);
    padding: 1rem;
    border: 1px solid rgba(0, 0, 0, 0.1);
  }
  
  .network-header {
    margin-bottom: 0.75rem;
  }
  
  .network-name {
    font-weight: 600;
    font-size: 1.05rem;
    color: var(--color-text, #0a0a0a);
    font-family: var(--font-heading, "Playfair Display", serif);
  }
  
  .network-details {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .network-detail {
    display: flex;
    gap: 0.5rem;
    align-items: baseline;
  }
  
  .detail-label {
    font-size: 0.9rem;
    color: var(--color-muted, #78716c);
    min-width: 80px;
  }
  
  .detail-value {
    font-size: 0.9rem;
    color: var(--color-text, #0a0a0a);
    word-break: break-all;
  }
  
  .ports-list {
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
  }
  
  .port-item {
    background: var(--color-background, #f5f5f4);
    border-radius: var(--radius, 0.25rem);
    padding: 0.5rem 0.75rem;
    border: 1px solid rgba(0, 0, 0, 0.1);
  }
  
  .port-mapping {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-family: monospace;
    font-size: 0.9rem;
  }
  
  .host-port {
    color: var(--color-secondary, #525252);
    font-weight: 600;
  }
  
  .port-arrow {
    color: var(--color-muted, #78716c);
  }
  
  .container-port {
    color: var(--color-text, #0a0a0a);
    font-weight: 600;
  }
  
  .port-note {
    color: var(--color-muted, #78716c);
    font-size: 0.85em;
  }
  
  .mount-item {
    background: var(--color-background, #f5f5f4);
    border-radius: var(--radius, 0.25rem);
    padding: 1rem;
    border: 1px solid rgba(0, 0, 0, 0.1);
    display: flex;
    gap: 1rem;
  }
  
  .mount-type-badge {
    padding: 0.25rem 0.5rem;
    border-radius: var(--radius, 0.25rem);
    font-size: 0.75rem;
    font-weight: 600;
    text-transform: uppercase;
    align-self: flex-start;
    background: var(--color-secondary, #525252);
    color: white;
  }
  
  .mount-type-badge.volume-type {
    background: var(--color-secondary, #525252);
  }
  
  .mount-type-badge.bind-type {
    background: var(--color-warning, #b45309);
  }
  
  .mount-details {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .mount-path {
    display: flex;
    gap: 0.5rem;
  }
  
  .path-label {
    font-size: 0.9rem;
    color: var(--color-muted, #78716c);
    min-width: 40px;
  }
  
  .path-value {
    font-size: 0.9rem;
    color: var(--color-text, #0a0a0a);
    word-break: break-all;
  }
  
  .mount-meta {
    display: flex;
    gap: 0.75rem;
    margin-top: 0.25rem;
  }
  
  .mount-rw {
    font-size: 0.85rem;
    padding: 0.15rem 0.4rem;
    border-radius: var(--radius, 0.25rem);
    background: var(--color-success, #15803d);
    color: white;
    font-weight: 500;
  }
  
  .mount-rw.readonly {
    background: var(--color-muted, #78716c);
  }
  
  .mount-mode {
    font-size: 0.85rem;
    color: var(--color-muted, #78716c);
  }
  
  .env-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    max-height: 400px;
    overflow-y: auto;
  }
  
  .env-item {
    background: var(--color-background, #f5f5f4);
    border-radius: var(--radius, 0.25rem);
    padding: 0.5rem 0.75rem;
    border: 1px solid rgba(0, 0, 0, 0.1);
    display: flex;
    gap: 0.5rem;
    align-items: baseline;
    font-size: 0.9rem;
  }
  
  .env-key {
    color: var(--color-secondary, #525252);
    font-weight: 600;
  }
  
  .env-equals {
    color: var(--color-muted, #78716c);
  }
  
  .env-value {
    color: var(--color-text, #0a0a0a);
    word-break: break-all;
  }
  
  .labels-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    max-height: 400px;
    overflow-y: auto;
  }
  
  .label-item {
    background: var(--color-background, #f5f5f4);
    border-radius: var(--radius, 0.25rem);
    padding: 0.5rem 0.75rem;
    border: 1px solid rgba(0, 0, 0, 0.1);
    display: flex;
    gap: 0.5rem;
    align-items: baseline;
    font-size: 0.9rem;
  }
  
  .label-key {
    color: var(--color-secondary, #525252);
    font-weight: 600;
  }
  
  .label-equals {
    color: var(--color-muted, #78716c);
  }
  
  .label-value {
    color: var(--color-text, #0a0a0a);
    word-break: break-all;
  }
  
  @media (max-width: 640px) {
    .main-content {
      padding: 1rem;
    }
    
    .info-grid {
      grid-template-columns: 1fr;
    }
    
    .mount-item {
      flex-direction: column;
      gap: 0.75rem;
    }
  }
</style>
