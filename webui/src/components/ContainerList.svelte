<script lang="ts">
  import type { Container } from '../types';
  
  interface Props {
    containers: Container[];
    displayMode: 'compact' | 'standard';
    onAction: (containerId: string, action: 'start' | 'stop' | 'restart', isSelf: boolean) => void;
  }
  
  let { containers, displayMode, onAction }: Props = $props();
  
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
</script>

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
        {#if container.compose_service}
          <span class="compose-service-badge">{container.compose_service}</span>
        {/if}
        <span class="compact-image" title={container.image}>{container.image}</span>
        <span class="compact-state">{container.status}</span>
        <div class="compact-actions">
          {#if container.state === 'running'}
            <button 
              class="action-btn-compact stop" 
              onclick={() => onAction(container.id, 'stop', container.is_self ?? false)}
              disabled={container.is_self}
              title={container.is_self ? "无法停止本应用容器" : "停止"}
            >
              ⏸️
            </button>
            <button 
              class="action-btn-compact restart" 
              onclick={() => onAction(container.id, 'restart', container.is_self ?? false)}
              disabled={container.is_self}
              title={container.is_self ? "无法重启本应用容器" : "重启"}
            >
              🔄
            </button>
          {:else if ['exited', 'created', 'dead'].includes(container.state)}
            <button 
              class="action-btn-compact start" 
              onclick={() => onAction(container.id, 'start', container.is_self ?? false)}
              title="启动"
            >
              ▶️
            </button>
          {:else}
            <button 
              class="action-btn-compact restart" 
              onclick={() => onAction(container.id, 'restart', container.is_self ?? false)}
              disabled={container.is_self}
              title={container.is_self ? "无法重启本应用容器" : "重启"}
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
            onclick={() => onAction(container.id, 'stop', container.is_self ?? false)}
            disabled={container.is_self}
            title={container.is_self ? "无法停止本应用容器" : ""}
          >
            ⏸️ 停止
          </button>
          <button 
            class="action-btn restart" 
            onclick={() => onAction(container.id, 'restart', container.is_self ?? false)}
            disabled={container.is_self}
            title={container.is_self ? "无法重启本应用容器" : ""}
          >
            🔄 重启
          </button>
        {:else if ['exited', 'created', 'dead'].includes(container.state)}
          <button 
            class="action-btn start"
            onclick={() => onAction(container.id, 'start', container.is_self ?? false)}
          >
            ▶️ 启动
          </button>
        {:else}
          <button 
            class="action-btn restart" 
            onclick={() => onAction(container.id, 'restart', container.is_self ?? false)}
            disabled={container.is_self}
            title={container.is_self ? "无法重启本应用容器" : ""}
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

<style>
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
  
  .compose-service-badge {
    background: var(--color-secondary, #525252);
    color: white;
    font-size: 0.7rem;
    padding: 0.15rem 0.5rem;
    border-radius: var(--radius, 0.25rem);
    font-weight: 500;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  /* Mobile responsive styles */
  @media (max-width: 640px) {
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
  }
</style>
