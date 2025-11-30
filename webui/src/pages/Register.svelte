<script lang="ts">
  import { onMount } from 'svelte';
  import { push } from 'svelte-spa-router';
  import { authApi, configApi } from '../services/api';
  import { themeStore, getThemeIcon } from '../stores/themeStore';
  
  let username = '';
  let password = '';
  let confirmPassword = '';
  let nickname = '';
  let error = '';
  let success = '';
  let loading = false;
  let registrationEnabled = true;
  let configLoading = true;
  
  onMount(async () => {
    try {
      const config = await configApi.getPublicConfig();
      registrationEnabled = !config.disable_registration;
    } catch (err) {
      // Default to allowing registration if config check fails
      registrationEnabled = true;
    } finally {
      configLoading = false;
    }
  });
  
  async function handleSubmit(e: Event) {
    e.preventDefault();
    error = '';
    success = '';
    
    // Validate passwords match
    if (password !== confirmPassword) {
      error = '两次输入的密码不一致';
      return;
    }
    
    // Validate password length
    if (password.length < 6) {
      error = '密码长度至少为6位';
      return;
    }
    
    loading = true;
    
    try {
      await authApi.register({
        username,
        password,
        nickname: nickname || username
      });
      success = '注册成功！正在跳转到登录页面...';
      setTimeout(() => {
        push('/login');
      }, 1500);
    } catch (err: any) {
      error = err.message || '注册失败，请重试';
      console.error('Register error:', err);
    } finally {
      loading = false;
    }
  }
  
  function goToLogin() {
    push('/login');
  }
  
  function toggleTheme() {
    themeStore.toggle();
  }
</script>

<div class="register-container">
  <button 
    class="theme-toggle" 
    on:click={toggleTheme}
    title="切换主题"
    aria-label="切换主题"
  >
    <span class="theme-icon">{getThemeIcon($themeStore.theme)}</span>
  </button>
  
  <div class="register-box">
    <div class="register-header">
      <svg class="logo-icon" width="48" height="48" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
        <rect x="4" y="8" width="24" height="16" rx="2" stroke="currentColor" stroke-width="2"/>
        <rect x="8" y="12" width="4" height="8" fill="currentColor"/>
        <rect x="14" y="12" width="4" height="8" fill="currentColor"/>
        <rect x="20" y="12" width="4" height="8" fill="currentColor"/>
        <circle cx="16" cy="4" r="2" fill="currentColor"/>
      </svg>
      <h1>Docker Simple Panel</h1>
      <p>创建新账户</p>
    </div>
    
    {#if configLoading}
      <div class="loading-message">
        加载中...
      </div>
    {:else if !registrationEnabled}
      <div class="disabled-message">
        <p>注册功能已被管理员禁用</p>
        <div class="login-link">
          <button type="button" on:click={goToLogin} class="link-button">返回登录</button>
        </div>
      </div>
    {:else}
      <form on:submit={handleSubmit} class="register-form">
        {#if error}
          <div class="error-message">
            {error}
          </div>
        {/if}
        
        {#if success}
          <div class="success-message">
            {success}
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
          <label for="nickname">昵称 <span class="optional">(可选)</span></label>
          <input
            type="text"
            id="nickname"
            bind:value={nickname}
            disabled={loading}
            placeholder="请输入昵称"
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
            placeholder="请输入密码（至少6位）"
          />
        </div>
        
        <div class="form-group">
          <label for="confirmPassword">确认密码</label>
          <input
            type="password"
            id="confirmPassword"
            bind:value={confirmPassword}
            required
            disabled={loading}
            placeholder="请再次输入密码"
          />
        </div>
        
        <button type="submit" class="register-button" disabled={loading}>
          {loading ? '注册中...' : '注册'}
        </button>
        
        <div class="login-link">
          已有账户？<button type="button" on:click={goToLogin} class="link-button">立即登录</button>
        </div>
      </form>
    {/if}
  </div>
</div>

<style>
  .register-container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--color-primary, #171717);
    padding: 1rem;
    position: relative;
  }
  
  .register-box {
    background: var(--color-surface, #e7e5e4);
    border-radius: var(--radius, 0.25rem);
    padding: 2.5rem;
    width: 100%;
    max-width: 420px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4);
  }
  
  .register-header {
    text-align: center;
    margin-bottom: 2rem;
  }
  
  .logo-icon {
    color: var(--color-primary, #171717);
    margin-bottom: 1rem;
  }
  
  .register-header h1 {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--color-text, #0a0a0a);
    margin: 0 0 0.5rem 0;
    font-family: var(--font-heading, "Playfair Display", serif);
  }
  
  .register-header p {
    color: var(--color-muted, #78716c);
    margin: 0;
    font-size: 0.95rem;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .loading-message {
    text-align: center;
    color: var(--color-muted, #78716c);
    padding: 2rem 0;
    font-size: 1rem;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .disabled-message {
    text-align: center;
    padding: 1rem 0;
  }
  
  .disabled-message p {
    color: var(--color-error, #991b1b);
    margin: 0 0 1.5rem 0;
    font-size: 1rem;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .register-form {
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
  
  .success-message {
    background: rgba(21, 128, 61, 0.1);
    border: 1px solid var(--color-success, #15803d);
    color: var(--color-success, #15803d);
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
  
  .form-group label .optional {
    font-weight: 400;
    color: var(--color-muted, #78716c);
    font-size: 0.85rem;
  }
  
  .form-group input {
    padding: 0.75rem;
    border: 1px solid rgba(0, 0, 0, 0.2);
    border-radius: var(--radius, 0.25rem);
    font-size: 1rem;
    transition: border-color 0.2s, box-shadow 0.2s;
    background: var(--color-background, #f5f5f4);
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
  
  .register-button {
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
  
  .register-button:hover:not(:disabled) {
    background: var(--color-secondary, #525252);
  }
  
  .register-button:active:not(:disabled) {
    transform: translateY(1px);
  }
  
  .register-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  .login-link {
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
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    padding: 0.5rem;
    border-radius: var(--radius, 0.25rem);
    cursor: pointer;
    transition: background 0.2s;
    width: 40px;
    height: 40px;
  }
  
  .theme-toggle:hover {
    background: rgba(255, 255, 255, 0.2);
  }
  
  .theme-icon {
    font-size: 1.25rem;
    line-height: 1;
  }
</style>
