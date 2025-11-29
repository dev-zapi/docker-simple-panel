<script lang="ts">
  import { onMount } from 'svelte';
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
    background: #f5f5f5;
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
  }
  
  .content-header h2 {
    font-size: 1.75rem;
    font-weight: 700;
    color: #333;
    margin: 0;
  }
  
  .add-button {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    padding: 0.625rem 1.25rem;
    border-radius: 8px;
    cursor: pointer;
    font-size: 0.95rem;
    font-weight: 600;
    transition: transform 0.2s, box-shadow 0.2s;
  }
  
  .add-button:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
  }
  
  .error-banner {
    background: #fee;
    border: 1px solid #fcc;
    color: #c33;
    padding: 1rem;
    border-radius: 8px;
    margin-bottom: 1.5rem;
  }
  
  .loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 4rem 2rem;
    color: #666;
  }
  
  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid #f3f3f3;
    border-top: 4px solid #667eea;
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
    color: #999;
  }
  
  .empty-icon {
    font-size: 4rem;
    margin-bottom: 1rem;
  }
  
  .users-table {
    background: white;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  }
  
  table {
    width: 100%;
    border-collapse: collapse;
  }
  
  thead {
    background: #f8f9fa;
  }
  
  th {
    padding: 1rem;
    text-align: left;
    font-weight: 600;
    color: #333;
    border-bottom: 2px solid #e0e0e0;
  }
  
  td {
    padding: 1rem;
    border-bottom: 1px solid #f0f0f0;
  }
  
  tr:last-child td {
    border-bottom: none;
  }
  
  tbody tr:hover {
    background: #f8f9fa;
  }
  
  .username {
    font-weight: 500;
    color: #333;
  }
  
  .nickname {
    color: #666;
  }
  
  .actions {
    text-align: right;
  }
  
  .delete-btn {
    background: #e74c3c;
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 6px;
    cursor: pointer;
    font-size: 0.9rem;
    transition: background 0.2s;
  }
  
  .delete-btn:hover {
    background: #c0392b;
  }
  
  /* Modal styles */
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }
  
  .modal {
    background: white;
    border-radius: 12px;
    padding: 2rem;
    width: 90%;
    max-width: 480px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
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
    color: #333;
    margin: 0;
  }
  
  .close-btn {
    background: none;
    border: none;
    font-size: 1.5rem;
    color: #999;
    cursor: pointer;
    padding: 0;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    transition: background 0.2s;
  }
  
  .close-btn:hover {
    background: #f0f0f0;
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
  
  .modal-actions {
    display: flex;
    gap: 0.75rem;
    justify-content: flex-end;
    margin-top: 0.5rem;
  }
  
  .cancel-btn {
    background: white;
    border: 2px solid #e0e0e0;
    padding: 0.625rem 1.25rem;
    border-radius: 8px;
    cursor: pointer;
    font-size: 0.95rem;
    transition: all 0.2s;
  }
  
  .cancel-btn:hover:not(:disabled) {
    border-color: #999;
  }
  
  .cancel-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  .submit-btn {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    padding: 0.625rem 1.25rem;
    border-radius: 8px;
    cursor: pointer;
    font-size: 0.95rem;
    font-weight: 600;
    transition: transform 0.2s, box-shadow 0.2s;
  }
  
  .submit-btn:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
  }
  
  .submit-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
</style>
