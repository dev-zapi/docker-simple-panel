<script lang="ts">
  import { push } from 'svelte-spa-router';
  import Header from '../components/Header.svelte';
  import { authStore } from '../stores/authStore';
  import { userApi } from '../services/api';
  
  let nickname = $authStore.user?.nickname || '';
  let currentPassword = '';
  let newPassword = '';
  let confirmPassword = '';
  let saving = false;
  let error = '';
  let success = '';
  
  async function handleUpdateProfile(e: Event) {
    e.preventDefault();
    error = '';
    success = '';
    
    if (!$authStore.user) return;
    
    if (newPassword && newPassword !== confirmPassword) {
      error = '两次输入的密码不一致';
      return;
    }
    
    saving = true;
    
    try {
      const updateData: any = { nickname };
      
      if (newPassword) {
        updateData.password = newPassword;
      }
      
      const updatedUser = await userApi.updateUser($authStore.user.id, updateData);
      authStore.updateUser(updatedUser);
      success = '个人信息更新成功';
      
      // Clear password fields
      currentPassword = '';
      newPassword = '';
      confirmPassword = '';
    } catch (err: any) {
      error = err.message || '更新失败（后端暂不支持用户信息修改功能）';
      console.error('Failed to update profile:', err);
    } finally {
      saving = false;
    }
  }
  
  function handleCancel() {
    push('/');
  }
</script>

<div class="profile-container">
  <Header />
  
  <main class="main-content">
    <div class="profile-box">
      <h2>编辑个人信息</h2>
      
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
      
      <form on:submit={handleUpdateProfile} class="profile-form">
        <div class="form-group">
          <label for="username">用户名</label>
          <input
            type="text"
            id="username"
            value={$authStore.user?.username}
            disabled
            class="disabled-input"
          />
          <span class="help-text">用户名不可修改</span>
        </div>
        
        <div class="form-group">
          <label for="nickname">昵称</label>
          <input
            type="text"
            id="nickname"
            bind:value={nickname}
            required
            disabled={saving}
            placeholder="请输入昵称"
          />
        </div>
        
        <div class="divider"></div>
        
        <div class="section-title">修改密码（可选）</div>
        
        <div class="form-group">
          <label for="new-password">新密码</label>
          <input
            type="password"
            id="new-password"
            bind:value={newPassword}
            disabled={saving}
            placeholder="留空表示不修改密码"
          />
        </div>
        
        <div class="form-group">
          <label for="confirm-password">确认新密码</label>
          <input
            type="password"
            id="confirm-password"
            bind:value={confirmPassword}
            disabled={saving}
            placeholder="请再次输入新密码"
          />
        </div>
        
        <div class="form-actions">
          <button type="button" class="cancel-btn" on:click={handleCancel} disabled={saving}>
            取消
          </button>
          <button type="submit" class="save-btn" disabled={saving}>
            {saving ? '保存中...' : '保存'}
          </button>
        </div>
      </form>
    </div>
  </main>
</div>

<style>
  .profile-container {
    min-height: 100vh;
    background: var(--color-background, #f5f5f4);
  }
  
  .main-content {
    max-width: 600px;
    margin: 0 auto;
    padding: 2rem;
  }
  
  .profile-box {
    background: var(--color-surface, #e7e5e4);
    border-radius: var(--radius, 0.25rem);
    padding: 2rem;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  }
  
  .profile-box h2 {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--color-text, #0a0a0a);
    margin: 0 0 1.5rem 0;
    font-family: var(--font-heading, "Playfair Display", serif);
  }
  
  .error-message {
    background: rgba(153, 27, 27, 0.1);
    border: 1px solid var(--color-error, #991b1b);
    color: var(--color-error, #991b1b);
    padding: 0.75rem;
    border-radius: var(--radius, 0.25rem);
    margin-bottom: 1rem;
    font-size: 0.9rem;
  }
  
  .success-message {
    background: rgba(21, 128, 61, 0.1);
    border: 1px solid var(--color-success, #15803d);
    color: var(--color-success, #15803d);
    padding: 0.75rem;
    border-radius: var(--radius, 0.25rem);
    margin-bottom: 1rem;
    font-size: 0.9rem;
  }
  
  .profile-form {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
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
  
  .form-group input:focus:not(:disabled) {
    outline: none;
    border-color: var(--color-primary, #171717);
    box-shadow: 0 0 0 2px rgba(23, 23, 23, 0.2);
  }
  
  .form-group input:disabled,
  .disabled-input {
    background: rgba(0, 0, 0, 0.05);
    cursor: not-allowed;
    color: var(--color-muted, #78716c);
  }
  
  .help-text {
    font-size: 0.85rem;
    color: var(--color-muted, #78716c);
  }
  
  .divider {
    height: 1px;
    background: rgba(0, 0, 0, 0.1);
    margin: 0.5rem 0;
  }
  
  .section-title {
    font-weight: 600;
    color: var(--color-secondary, #525252);
    font-size: 0.95rem;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .form-actions {
    display: flex;
    gap: 0.75rem;
    justify-content: flex-end;
    margin-top: 0.5rem;
  }
  
  .cancel-btn {
    background: transparent;
    border: 1px solid rgba(0, 0, 0, 0.2);
    padding: 0.625rem 1.5rem;
    border-radius: var(--radius, 0.25rem);
    cursor: pointer;
    font-size: 0.95rem;
    transition: all 0.2s;
    color: var(--color-text, #0a0a0a);
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .cancel-btn:hover:not(:disabled) {
    background: rgba(0, 0, 0, 0.05);
  }
  
  .cancel-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  .save-btn {
    background: var(--color-primary, #171717);
    color: var(--color-background, #f5f5f4);
    border: none;
    padding: 0.625rem 1.5rem;
    border-radius: var(--radius, 0.25rem);
    cursor: pointer;
    font-size: 0.95rem;
    font-weight: 600;
    transition: background 0.2s, transform 0.2s;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .save-btn:hover:not(:disabled) {
    background: var(--color-secondary, #525252);
  }
  
  .save-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
</style>
