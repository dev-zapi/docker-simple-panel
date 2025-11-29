<script lang="ts">
  import { onMount } from 'svelte';
  import { push } from 'svelte-spa-router';
  import { authStore } from '../stores/authStore';
  import { authApi, configApi } from '../services/api';
  
  let username = '';
  let password = '';
  let error = '';
  let loading = false;
  let registrationEnabled = false;
  let configLoading = true;
  
  onMount(async () => {
    try {
      const config = await configApi.getPublicConfig();
      registrationEnabled = !config.disable_registration;
    } catch (err) {
      // Default to showing registration link if config check fails
      registrationEnabled = true;
    } finally {
      configLoading = false;
    }
  });
  
  async function handleSubmit(e: Event) {
    e.preventDefault();
    error = '';
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
  
  function goToRegister() {
    push('/register');
  }
</script>

<div class="login-container">
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
      
      {#if !configLoading && registrationEnabled}
        <div class="register-link">
          没有账户？<button type="button" on:click={goToRegister} class="link-button">立即注册</button>
        </div>
      {/if}
    </form>
  </div>
</div>

<style>
  .login-container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 1rem;
  }
  
  .login-box {
    background: white;
    border-radius: 16px;
    padding: 2.5rem;
    width: 100%;
    max-width: 420px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  }
  
  .login-header {
    text-align: center;
    margin-bottom: 2rem;
  }
  
  .logo-icon {
    color: #667eea;
    margin-bottom: 1rem;
  }
  
  .login-header h1 {
    font-size: 1.75rem;
    font-weight: 700;
    color: #333;
    margin: 0 0 0.5rem 0;
  }
  
  .login-header p {
    color: #666;
    margin: 0;
    font-size: 0.95rem;
  }
  
  .login-form {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }
  
  .error-message {
    background: #fee;
    border: 1px solid #fcc;
    color: #c33;
    padding: 0.75rem;
    border-radius: 8px;
    font-size: 0.9rem;
  }
  
  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .form-group label {
    font-weight: 600;
    color: #333;
    font-size: 0.95rem;
  }
  
  .form-group input {
    padding: 0.75rem;
    border: 2px solid #e0e0e0;
    border-radius: 8px;
    font-size: 1rem;
    transition: border-color 0.2s;
  }
  
  .form-group input:focus {
    outline: none;
    border-color: #667eea;
  }
  
  .form-group input:disabled {
    background: #f5f5f5;
    cursor: not-allowed;
  }
  
  .login-button {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    padding: 0.875rem;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: transform 0.2s, box-shadow 0.2s;
    margin-top: 0.5rem;
  }
  
  .login-button:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
  }
  
  .login-button:active:not(:disabled) {
    transform: translateY(0);
  }
  
  .login-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  .register-link {
    text-align: center;
    color: #666;
    font-size: 0.9rem;
  }
  
  .link-button {
    background: none;
    border: none;
    color: #667eea;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    padding: 0;
    text-decoration: underline;
  }
  
  .link-button:hover {
    color: #764ba2;
  }
</style>
