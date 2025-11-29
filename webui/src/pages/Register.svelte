<script lang="ts">
  import { onMount } from 'svelte';
  import { push } from 'svelte-spa-router';
  import { authApi, configApi } from '../services/api';
  
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
</script>

<div class="register-container">
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
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 1rem;
  }
  
  .register-box {
    background: white;
    border-radius: 16px;
    padding: 2.5rem;
    width: 100%;
    max-width: 420px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  }
  
  .register-header {
    text-align: center;
    margin-bottom: 2rem;
  }
  
  .logo-icon {
    color: #667eea;
    margin-bottom: 1rem;
  }
  
  .register-header h1 {
    font-size: 1.75rem;
    font-weight: 700;
    color: #333;
    margin: 0 0 0.5rem 0;
  }
  
  .register-header p {
    color: #666;
    margin: 0;
    font-size: 0.95rem;
  }
  
  .loading-message {
    text-align: center;
    color: #666;
    padding: 2rem 0;
    font-size: 1rem;
  }
  
  .disabled-message {
    text-align: center;
    padding: 1rem 0;
  }
  
  .disabled-message p {
    color: #c33;
    margin: 0 0 1.5rem 0;
    font-size: 1rem;
  }
  
  .register-form {
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
  
  .success-message {
    background: #efe;
    border: 1px solid #cfc;
    color: #3c3;
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
  
  .form-group label .optional {
    font-weight: 400;
    color: #999;
    font-size: 0.85rem;
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
  
  .register-button {
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
  
  .register-button:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
  }
  
  .register-button:active:not(:disabled) {
    transform: translateY(0);
  }
  
  .register-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  .login-link {
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
