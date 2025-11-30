import { writable } from 'svelte/store';

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
  const { subscribe, set, update } = writable<PageHeaderState>(createDefaultState());

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
      update(state => {
        if (state.onToggleDisplayMode) {
          state.onToggleDisplayMode();
        }
        return state;
      });
    },
    triggerRefresh: () => {
      update(state => {
        if (state.onRefresh) {
          state.onRefresh();
        }
        return state;
      });
    },
    reset: () => {
      set(createDefaultState());
    },
  };
};

export const pageHeaderStore = createPageHeaderStore();
