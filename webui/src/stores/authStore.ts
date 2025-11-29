import { writable, derived } from 'svelte/store';
import type { AuthState, User } from '../types';

// Load auth state from localStorage
const loadAuthState = (): AuthState => {
  if (typeof window === 'undefined') {
    return { isAuthenticated: false, user: null, token: null };
  }
  
  const token = localStorage.getItem('token');
  const userStr = localStorage.getItem('user');
  
  if (token && userStr) {
    try {
      const user = JSON.parse(userStr) as User;
      return { isAuthenticated: true, user, token };
    } catch (e) {
      return { isAuthenticated: false, user: null, token: null };
    }
  }
  
  return { isAuthenticated: false, user: null, token: null };
};

// Create the auth store
const createAuthStore = () => {
  const { subscribe, set, update } = writable<AuthState>(loadAuthState());
  
  return {
    subscribe,
    login: (user: User, token: string) => {
      if (typeof window !== 'undefined') {
        localStorage.setItem('token', token);
        localStorage.setItem('user', JSON.stringify(user));
      }
      set({ isAuthenticated: true, user, token });
    },
    logout: () => {
      if (typeof window !== 'undefined') {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
      }
      set({ isAuthenticated: false, user: null, token: null });
    },
    updateUser: (user: User) => {
      update(state => {
        if (typeof window !== 'undefined') {
          localStorage.setItem('user', JSON.stringify(user));
        }
        return { ...state, user };
      });
    }
  };
};

export const authStore = createAuthStore();

// Derived store for authentication status
export const isAuthenticated = derived(
  authStore,
  $authStore => $authStore.isAuthenticated
);
