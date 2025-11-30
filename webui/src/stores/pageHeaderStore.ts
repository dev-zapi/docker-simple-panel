import { writable, get } from 'svelte/store';

export type DisplayMode = 'compact' | 'standard';

export interface PageHeaderState {
  title: string;
  showDisplayModeToggle: boolean;
  displayMode: DisplayMode;
  showRefreshButton: boolean;
  refreshing: boolean;
  onToggleDisplayMode: (() => void) | null;
  onRefresh: (() => void) | null;
}

const createDefaultState = (): PageHeaderState => ({
  title: '',
  showDisplayModeToggle: false,
  displayMode: 'standard',
  showRefreshButton: false,
  refreshing: false,
  onToggleDisplayMode: null,
  onRefresh: null,
});

const createPageHeaderStore = () => {
  const store = writable<PageHeaderState>(createDefaultState());
  const { subscribe, set, update } = store;

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
      update(state => ({ ...state, onToggleDisplayMode: callback }));
    },
    setOnRefresh: (callback: (() => void) | null) => {
      update(state => ({ ...state, onRefresh: callback }));
    },
    triggerToggleDisplayMode: () => {
      const state = get(store);
      if (state.onToggleDisplayMode) {
        state.onToggleDisplayMode();
      }
    },
    triggerRefresh: () => {
      const state = get(store);
      if (state.onRefresh) {
        state.onRefresh();
      }
    },
    reset: () => {
      set(createDefaultState());
    },
  };
};

export const pageHeaderStore = createPageHeaderStore();
