<script lang="ts">
  import { authStore } from '../stores/authStore';
  import { push } from 'svelte-spa-router';
  
  let showMenu = false;
  
  function toggleMenu() {
    showMenu = !showMenu;
  }
  
  function handleLogout() {
    authStore.logout();
    push('/login');
  }
  
  function goToUserManagement() {
    showMenu = false;
    push('/users');
  }
  
  function goToProfile() {
    showMenu = false;
    push('/profile');
  }
  
  function goToSettings() {
    showMenu = false;
    push('/settings');
  }
  
  function goToHome() {
    push('/');
  }
  
  // Close menu when clicking outside
  function handleClickOutside(event: MouseEvent) {
    if (showMenu && !(event.target as HTMLElement).closest('.user-menu')) {
      showMenu = false;
    }
  }
</script>

<svelte:window on:click={handleClickOutside} />

<header class="header">
  <div class="header-left">
    <button class="logo-button" on:click={goToHome}>
      <svg class="logo-icon" width="32" height="32" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
        <rect x="4" y="8" width="24" height="16" rx="2" stroke="currentColor" stroke-width="2"/>
        <rect x="8" y="12" width="4" height="8" fill="currentColor"/>
        <rect x="14" y="12" width="4" height="8" fill="currentColor"/>
        <rect x="20" y="12" width="4" height="8" fill="currentColor"/>
        <circle cx="16" cy="4" r="2" fill="currentColor"/>
      </svg>
      <h1 class="title">DSP</h1>
    </button>
  </div>
  
  {#if $authStore.isAuthenticated && $authStore.user}
    <div class="header-right">
      <div class="user-menu">
        <button class="user-button" on:click={toggleMenu}>
          <span class="user-name">{$authStore.user.nickname || $authStore.user.username}</span>
          <span class="menu-icon">▼</span>
        </button>
        
        {#if showMenu}
          <div class="dropdown-menu">
            <button class="menu-item" on:click={goToProfile}>
              👤 编辑个人信息
            </button>
            <button class="menu-item" on:click={goToUserManagement}>
              👥 用户管理
            </button>
            <button class="menu-item" on:click={goToSettings}>
              ⚙️ 系统设置
            </button>
            <button class="menu-item logout" on:click={handleLogout}>
              🚪 登出
            </button>
          </div>
        {/if}
      </div>
    </div>
  {/if}
</header>

<style>
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 2rem;
    background: var(--color-primary, #171717);
    color: var(--color-background, #f5f5f4);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
    position: sticky;
    top: 0;
    z-index: 100;
  }
  
  .header-left {
    display: flex;
    align-items: center;
    flex-shrink: 0;
  }
  
  .logo-button {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    background: none;
    border: none;
    cursor: pointer;
    color: inherit;
    padding: 0;
  }
  
  .logo-button:hover {
    opacity: 0.9;
  }
  
  .logo-icon {
    color: var(--color-background, #f5f5f4);
  }
  
  .title {
    font-size: 1.5rem;
    font-weight: 700;
    margin: 0;
    letter-spacing: 0.1em;
    font-family: var(--font-heading, "Playfair Display", serif);
  }
  
  .header-right {
    display: flex;
    align-items: center;
    flex-shrink: 0;
  }
  
  .user-menu {
    position: relative;
  }
  
  .user-button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    padding: 0.5rem 1rem;
    border-radius: var(--radius, 0.25rem);
    color: var(--color-background, #f5f5f4);
    cursor: pointer;
    font-size: 0.95rem;
    transition: background 0.2s;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .user-button:hover {
    background: rgba(255, 255, 255, 0.2);
  }
  
  .user-name {
    font-weight: 500;
  }
  
  .menu-icon {
    font-size: 0.7rem;
  }
  
  .dropdown-menu {
    position: absolute;
    top: 100%;
    right: 0;
    margin-top: 0.5rem;
    background: var(--color-surface, #e7e5e4);
    border-radius: var(--radius, 0.25rem);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
    min-width: 200px;
    overflow: hidden;
    z-index: 1000;
  }
  
  .menu-item {
    display: block;
    width: 100%;
    padding: 0.75rem 1rem;
    border: none;
    background: var(--color-surface, #e7e5e4);
    color: var(--color-text, #0a0a0a);
    text-align: left;
    cursor: pointer;
    font-size: 0.95rem;
    transition: background 0.2s;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .menu-item:hover {
    background: var(--color-background, #f5f5f4);
  }
  
  .menu-item.logout {
    border-top: 1px solid rgba(0, 0, 0, 0.1);
    color: var(--color-error, #991b1b);
  }
  
  .menu-item.logout:hover {
    background: rgba(153, 27, 27, 0.1);
  }
</style>
