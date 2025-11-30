import { writable } from 'svelte/store';

export type DisplayMode = 'compact' | 'standard';

export interface PageHeaderState {
  title: string;
  showDisplayModeToggle: boolean;
  displayMode: DisplayMode;
  showRefreshButton: boolean;
  refreshing: boolean;
}

const defaultState: PageHeaderState = {
  title: '',
  showDisplayModeToggle: false,
  displayMode: 'standard',
  showRefreshButton: false,
  refreshing: false,
};

const createPageHeaderStore = () => {
  const { subscribe, set, update } = writable<PageHeaderState>(defaultState);

  // Callbacks for actions
  let onToggleDisplayMode: (() => void) | null = null;
  let onRefresh: (() => void) | null = null;

  return {
    subscribe,
    setTitle: (title: string) => {
      update(state => ({ ...state, title }));
    },
    setDisplayMode: (displayMode: DisplayMode) => {
      update(state => ({ ...state, displayMode }));
    },
    setShowDisplayModeToggle: (show: boolean) => {
      update(state => ({ ...state, showDisplayModeToggle: show }));
    },
    setShowRefreshButton: (show: boolean) => {
      update(state => ({ ...state, showRefreshButton: show }));
    },
    setRefreshing: (refreshing: boolean) => {
      update(state => ({ ...state, refreshing }));
    },
    setOnToggleDisplayMode: (callback: (() => void) | null) => {
      onToggleDisplayMode = callback;
    },
    setOnRefresh: (callback: (() => void) | null) => {
      onRefresh = callback;
    },
    triggerToggleDisplayMode: () => {
      if (onToggleDisplayMode) {
        onToggleDisplayMode();
      }
    },
    triggerRefresh: () => {
      if (onRefresh) {
        onRefresh();
      }
    },
    reset: () => {
      set(defaultState);
      onToggleDisplayMode = null;
      onRefresh = null;
    },
  };
};

export const pageHeaderStore = createPageHeaderStore();
