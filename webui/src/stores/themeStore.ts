import { writable, derived } from 'svelte/store';

export type Theme = 'light' | 'dark' | 'system';
export type ResolvedTheme = 'light' | 'dark';

interface ThemeState {
  theme: Theme;
  resolvedTheme: ResolvedTheme;
}

// Get system preference
const getSystemTheme = (): ResolvedTheme => {
  if (typeof window === 'undefined') return 'light';
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
};

// Load theme from localStorage
const loadTheme = (): Theme => {
  if (typeof window === 'undefined') return 'system';
  const saved = localStorage.getItem('theme');
  if (saved === 'light' || saved === 'dark' || saved === 'system') {
    return saved;
  }
  return 'system';
};

// Resolve theme based on preference
const resolveTheme = (theme: Theme): ResolvedTheme => {
  if (theme === 'system') {
    return getSystemTheme();
  }
  return theme;
};

// Apply theme to document
const applyTheme = (resolvedTheme: ResolvedTheme) => {
  if (typeof document === 'undefined') return;
  document.documentElement.setAttribute('data-theme', resolvedTheme);
};

// Create theme store
const createThemeStore = () => {
  const initialTheme = loadTheme();
  const initialResolved = resolveTheme(initialTheme);
  
  const { subscribe, set, update } = writable<ThemeState>({
    theme: initialTheme,
    resolvedTheme: initialResolved
  });
  
  // Apply initial theme
  applyTheme(initialResolved);
  
  // Listen for system theme changes
  if (typeof window !== 'undefined') {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    mediaQuery.addEventListener('change', () => {
      update(state => {
        if (state.theme === 'system') {
          const newResolved = getSystemTheme();
          applyTheme(newResolved);
          return { ...state, resolvedTheme: newResolved };
        }
        return state;
      });
    });
  }
  
  return {
    subscribe,
    setTheme: (theme: Theme) => {
      if (typeof window !== 'undefined') {
        localStorage.setItem('theme', theme);
      }
      const resolved = resolveTheme(theme);
      applyTheme(resolved);
      set({ theme, resolvedTheme: resolved });
    },
    toggle: () => {
      update(state => {
        // Cycle through: system -> light -> dark -> system
        const nextTheme: Theme = 
          state.theme === 'system' ? 'light' :
          state.theme === 'light' ? 'dark' : 'system';
        
        if (typeof window !== 'undefined') {
          localStorage.setItem('theme', nextTheme);
        }
        const resolved = resolveTheme(nextTheme);
        applyTheme(resolved);
        return { theme: nextTheme, resolvedTheme: resolved };
      });
    }
  };
};

export const themeStore = createThemeStore();

// Derived stores for convenience
export const currentTheme = derived(themeStore, $store => $store.theme);
export const resolvedTheme = derived(themeStore, $store => $store.resolvedTheme);
