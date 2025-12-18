<script lang="ts">
  import { onMount } from 'svelte';
  import { push } from 'svelte-spa-router';
  import { authStore } from '../stores/authStore';
  import { themeStore, getThemeIcon } from '../stores/themeStore';
  import { authApi } from '../services/api';
  
  let username = '';
  let password = '';
  let error = '';
  let sessionExpiredMessage = '';
  let loading = false;
  
  onMount(async () => {
    // Check if session expired
    if (localStorage.getItem('sessionExpired') === 'true') {
      sessionExpiredMessage = '会话已过期，请重新登录';
      localStorage.removeItem('sessionExpired');
    }
  });
  
  async function handleSubmit(e: Event) {
    e.preventDefault();
    error = '';
    sessionExpiredMessage = '';
    loading = true;
    
    try {
      const response = await authApi.login({ username, password });
      authStore.login(response.user, response.token);
      push('/');
    } catch (err) {
      error = '登录失败：用户名或密码错误';
      console.error('Login error:', err);
    } finally {
      loading = false;
    }
  }
  
  function toggleTheme() {
    themeStore.toggle();
  }
</script>

<div class="login-container">
  <button 
    class="theme-toggle" 
    on:click={toggleTheme}
    title="切换主题"
    aria-label="切换主题"
  >
    <span class="theme-icon">{getThemeIcon($themeStore.theme)}</span>
  </button>
  
  <div class="login-box">
    <div class="login-header">
      <svg class="logo-icon" width="48" height="48" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
        <rect x="4" y="8" width="24" height="16" rx="2" stroke="currentColor" stroke-width="2"/>
        <rect x="8" y="12" width="4" height="8" fill="currentColor"/>
        <rect x="14" y="12" width="4" height="8" fill="currentColor"/>
        <rect x="20" y="12" width="4" height="8" fill="currentColor"/>
        <circle cx="16" cy="4" r="2" fill="currentColor"/>
      </svg>
      <h1>Docker Simple Panel</h1>
      <p>登录以管理您的容器</p>
    </div>
    
    <form on:submit={handleSubmit} class="login-form">
      {#if sessionExpiredMessage}
        <div class="info-message">
          {sessionExpiredMessage}
        </div>
      {/if}
      
      {#if error}
        <div class="error-message">
          {error}
        </div>
      {/if}
      
      <div class="form-group">
        <label for="username">用户名</label>
        <input
          type="text"
          id="username"
          bind:value={username}
          required
          disabled={loading}
          placeholder="请输入用户名"
        />
      </div>
      
      <div class="form-group">
        <label for="password">密码</label>
        <input
          type="password"
          id="password"
          bind:value={password}
          required
          disabled={loading}
          placeholder="请输入密码"
        />
      </div>
      
      <button type="submit" class="login-button" disabled={loading}>
        {loading ? '登录中...' : '登录'}
      </button>
    </form>
  </div>
</div>

<style>
  .login-container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--color-background, #f5f5f4);
    position: relative;
    padding: 1rem;
  }
  
  .login-box {
    background: var(--color-surface, #e7e5e4);
    border-radius: var(--radius, 0.25rem);
    padding: 2.5rem;
    width: 100%;
    max-width: 420px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4);
  }
  
  .login-header {
    text-align: center;
    margin-bottom: 2rem;
  }
  
  .logo-icon {
    color: var(--color-primary, #171717);
    margin-bottom: 1rem;
  }
  
  .login-header h1 {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--color-text, #0a0a0a);
    margin: 0 0 0.5rem 0;
    font-family: var(--font-heading, "Playfair Display", serif);
  }
  
  .login-header p {
    color: var(--color-muted, #78716c);
    margin: 0;
    font-size: 0.95rem;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .login-form {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }
  
  .error-message {
    background: rgba(153, 27, 27, 0.1);
    border: 1px solid var(--color-error, #991b1b);
    color: var(--color-error, #991b1b);
    padding: 0.75rem;
    border-radius: var(--radius, 0.25rem);
    font-size: 0.9rem;
  }
  
  .info-message {
    background: rgba(245, 158, 11, 0.1);
    border: 1px solid #f59e0b;
    color: #d97706;
    padding: 0.75rem;
    border-radius: var(--radius, 0.25rem);
    font-size: 0.9rem;
  }
  
  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .form-group label {
    font-weight: 600;
    color: var(--color-text, #0a0a0a);
    font-size: 0.95rem;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .form-group input {
    padding: 0.75rem;
    border: 1px solid rgba(0, 0, 0, 0.2);
    border-radius: var(--radius, 0.25rem);
    font-size: 1rem;
    transition: border-color 0.2s, box-shadow 0.2s;
    background: var(--color-background, #f5f5f4);
    color: var(--color-text, #0a0a0a);
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .form-group input:focus {
    outline: none;
    border-color: var(--color-primary, #171717);
    box-shadow: 0 0 0 2px rgba(23, 23, 23, 0.2);
  }
  
  .form-group input:disabled {
    background: rgba(0, 0, 0, 0.05);
    cursor: not-allowed;
  }
  
  .login-button {
    background: var(--color-primary, #171717);
    color: var(--color-background, #f5f5f4);
    border: none;
    padding: 0.875rem;
    border-radius: var(--radius, 0.25rem);
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.2s, transform 0.2s;
    margin-top: 0.5rem;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .login-button:hover:not(:disabled) {
    background: var(--color-secondary, #525252);
  }
  
  .login-button:active:not(:disabled) {
    transform: translateY(1px);
  }
  
  .login-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  .register-link {
    text-align: center;
    color: var(--color-muted, #78716c);
    font-size: 0.9rem;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .link-button {
    background: none;
    border: none;
    color: var(--color-accent, #991b1b);
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    padding: 0;
    text-decoration: underline;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .link-button:hover {
    color: var(--color-primary, #171717);
  }
  
  .theme-toggle {
    position: absolute;
    top: 1rem;
    right: 1rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--color-surface, #e7e5e4);
    border: var(--border-width, 0px) solid var(--color-muted, #78716c);
    padding: 0.5rem;
    border-radius: var(--radius, 0.25rem);
    cursor: pointer;
    transition: background 0.2s;
    width: 40px;
    height: 40px;
  }
  
  .theme-toggle:hover {
    background: var(--color-muted, #78716c);
  }
  
  .theme-icon {
    font-size: 1.25rem;
    line-height: 1;
  }
</style>
