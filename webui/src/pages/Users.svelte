<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import Header from '../components/Header.svelte';
  import { userApi, authApi } from '../services/api';
  import type { User } from '../types';
  
  let users: User[] = [];
  let loading = true;
  let error = '';
  let showAddModal = false;
  

  
  // Form fields for new user
  let newUsername = '';
  let newPassword = '';
  let newNickname = '';
  let addingUser = false;
  
  async function loadUsers() {
    try {
      error = '';
      loading = true;
      users = await userApi.getUsers();
    } catch (err) {
      error = '获取用户列表失败';
      console.error('Failed to load users:', err);
    } finally {
      loading = false;
    }
  }
  
  function openAddModal() {
    showAddModal = true;
    newUsername = '';
    newPassword = '';
    newNickname = '';
  }
  
  function closeAddModal() {
    showAddModal = false;
  }
  
  async function handleAddUser(e: Event) {
    e.preventDefault();
    
    if (!newUsername || !newPassword) {
      return;
    }
    
    addingUser = true;
    
    try {
      // Note: Backend uses /api/auth/register endpoint for user creation
      await authApi.register({
        username: newUsername,
        password: newPassword,
        nickname: newNickname || newUsername
      });
      closeAddModal();
      await loadUsers();
    } catch (err: any) {
      error = err.message || '添加用户失败';
      console.error('Failed to add user:', err);
    } finally {
      addingUser = false;
    }
  }
  
  async function handleDeleteUser(userId: number, username: string) {
    if (!confirm(`确定要删除用户 "${username}" 吗？`)) {
      return;
    }
    
    try {
      await userApi.deleteUser(userId);
      await loadUsers();
    } catch (err: any) {
      error = err.message || '删除用户失败（后端不支持此功能）';
      console.error('Failed to delete user:', err);
    }
  }
  
  onMount(() => {
    loadUsers();
  });
  
  onDestroy(() => {
  });
</script>

<div class="users-container">
  <Header />
  
  <main class="main-content">
    <div class="content-header">
      <h2>用户管理</h2>
      <button class="add-button" on:click={openAddModal}>
        ➕ 添加用户
      </button>
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
    {:else if users.length === 0}
      <div class="empty-state">
        <div class="empty-icon">👥</div>
        <p>暂无用户</p>
      </div>
    {:else}
      <div class="users-table">
        <table>
          <thead>
            <tr>
              <th>用户名</th>
              <th>昵称</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {#each users as user (user.id)}
              <tr>
                <td class="username">{user.username}</td>
                <td class="nickname">{user.nickname}</td>
                <td class="actions">
                  <button class="delete-btn" on:click={() => handleDeleteUser(user.id, user.username)}>
                    🗑️ 删除
                  </button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </main>
</div>

{#if showAddModal}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div class="modal-overlay" on:click={closeAddModal}>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div class="modal" on:click|stopPropagation>
      <div class="modal-header">
        <h3>添加新用户</h3>
        <button class="close-btn" on:click={closeAddModal}>✕</button>
      </div>
      
      <form on:submit={handleAddUser} class="modal-form">
        <div class="form-group">
          <label for="username">用户名 *</label>
          <input
            type="text"
            id="username"
            bind:value={newUsername}
            required
            disabled={addingUser}
            placeholder="请输入用户名"
          />
        </div>
        
        <div class="form-group">
          <label for="password">密码 *</label>
          <input
            type="password"
            id="password"
            bind:value={newPassword}
            required
            disabled={addingUser}
            placeholder="请输入密码"
          />
        </div>
        
        <div class="form-group">
          <label for="nickname">昵称</label>
          <input
            type="text"
            id="nickname"
            bind:value={newNickname}
            disabled={addingUser}
            placeholder="请输入昵称（可选）"
          />
        </div>
        
        <div class="modal-actions">
          <button type="button" class="cancel-btn" on:click={closeAddModal} disabled={addingUser}>
            取消
          </button>
          <button type="submit" class="submit-btn" disabled={addingUser}>
            {addingUser ? '添加中...' : '添加'}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

<style>
  .users-container {
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
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    position: sticky;
    top: 0;
    background: rgba(245, 245, 244, 0.8);
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
    z-index: 50;
    padding: 1rem 0;
  }
  
  .content-header h2 {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--color-text, #0a0a0a);
    margin: 0;
    font-family: var(--font-heading, "Playfair Display", serif);
  }
  
  .add-button {
    background: var(--color-primary, #171717);
    color: var(--color-background, #f5f5f4);
    border: none;
    padding: 0.625rem 1.25rem;
    border-radius: var(--radius, 0.25rem);
    cursor: pointer;
    font-size: 0.95rem;
    font-weight: 600;
    transition: background 0.2s;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .add-button:hover {
    background: var(--color-secondary, #525252);
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
  
  .users-table {
    background: var(--color-surface, #e7e5e4);
    border-radius: var(--radius, 0.25rem);
    overflow: hidden;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  }
  
  table {
    width: 100%;
    border-collapse: collapse;
  }
  
  thead {
    background: rgba(0, 0, 0, 0.05);
  }
  
  th {
    padding: 1rem;
    text-align: left;
    font-weight: 600;
    color: var(--color-text, #0a0a0a);
    border-bottom: 2px solid rgba(0, 0, 0, 0.1);
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  td {
    padding: 1rem;
    border-bottom: 1px solid rgba(0, 0, 0, 0.05);
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  tr:last-child td {
    border-bottom: none;
  }
  
  tbody tr:hover {
    background: rgba(0, 0, 0, 0.03);
  }
  
  .username {
    font-weight: 500;
    color: var(--color-text, #0a0a0a);
  }
  
  .nickname {
    color: var(--color-muted, #78716c);
  }
  
  .actions {
    text-align: right;
  }
  
  .delete-btn {
    background: var(--color-error, #991b1b);
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: var(--radius, 0.25rem);
    cursor: pointer;
    font-size: 0.9rem;
    transition: background 0.2s;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .delete-btn:hover {
    background: #7f1d1d;
  }
  
  /* Modal styles */
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }
  
  .modal {
    background: var(--color-surface, #e7e5e4);
    border-radius: var(--radius, 0.25rem);
    padding: 2rem;
    width: 90%;
    max-width: 480px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4);
  }
  
  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
  }
  
  .modal-header h3 {
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--color-text, #0a0a0a);
    margin: 0;
    font-family: var(--font-heading, "Playfair Display", serif);
  }
  
  .close-btn {
    background: none;
    border: none;
    font-size: 1.5rem;
    color: var(--color-muted, #78716c);
    cursor: pointer;
    padding: 0;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius, 0.25rem);
    transition: background 0.2s;
  }
  
  .close-btn:hover {
    background: rgba(0, 0, 0, 0.1);
  }
  
  .modal-form {
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
  
  .form-group input:focus {
    outline: none;
    border-color: var(--color-primary, #171717);
    box-shadow: 0 0 0 2px rgba(23, 23, 23, 0.2);
  }
  
  .form-group input:disabled {
    background: rgba(0, 0, 0, 0.05);
    cursor: not-allowed;
  }
  
  .modal-actions {
    display: flex;
    gap: 0.75rem;
    justify-content: flex-end;
    margin-top: 0.5rem;
  }
  
  .cancel-btn {
    background: transparent;
    border: 1px solid rgba(0, 0, 0, 0.2);
    padding: 0.625rem 1.25rem;
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
  
  .submit-btn {
    background: var(--color-primary, #171717);
    color: var(--color-background, #f5f5f4);
    border: none;
    padding: 0.625rem 1.25rem;
    border-radius: var(--radius, 0.25rem);
    cursor: pointer;
    font-size: 0.95rem;
    font-weight: 600;
    transition: background 0.2s;
    font-family: var(--font-body, "Merriweather", serif);
  }
  
  .submit-btn:hover:not(:disabled) {
    background: var(--color-secondary, #525252);
  }
  
  .submit-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
</style>
