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
    <svg class="logo-icon" width="32" height="32" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect x="4" y="8" width="24" height="16" rx="2" stroke="currentColor" stroke-width="2"/>
      <rect x="8" y="12" width="4" height="8" fill="currentColor"/>
      <rect x="14" y="12" width="4" height="8" fill="currentColor"/>
      <rect x="20" y="12" width="4" height="8" fill="currentColor"/>
      <circle cx="16" cy="4" r="2" fill="currentColor"/>
    </svg>
    <h1 class="title">DSP</h1>
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
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
  
  .header-left {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }
  
  .logo-icon {
    color: white;
  }
  
  .title {
    font-size: 1.5rem;
    font-weight: 700;
    margin: 0;
    letter-spacing: 0.05em;
  }
  
  .header-right {
    display: flex;
    align-items: center;
  }
  
  .user-menu {
    position: relative;
  }
  
  .user-button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: rgba(255, 255, 255, 0.2);
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 8px;
    color: white;
    cursor: pointer;
    font-size: 0.95rem;
    transition: background 0.2s;
  }
  
  .user-button:hover {
    background: rgba(255, 255, 255, 0.3);
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
    background: white;
    border-radius: 8px;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
    min-width: 200px;
    overflow: hidden;
    z-index: 1000;
  }
  
  .menu-item {
    display: block;
    width: 100%;
    padding: 0.75rem 1rem;
    border: none;
    background: white;
    color: #333;
    text-align: left;
    cursor: pointer;
    font-size: 0.95rem;
    transition: background 0.2s;
  }
  
  .menu-item:hover {
    background: #f5f5f5;
  }
  
  .menu-item.logout {
    border-top: 1px solid #eee;
    color: #e74c3c;
  }
  
  .menu-item.logout:hover {
    background: #fee;
  }
</style>
