import { writable, get } from 'svelte/store';

export type DisplayMode = 'compact' | 'standard';

export interface CustomAction {
  icon: string;
  label: string;
  onClick: () => void;
}

export interface PageHeaderState {
  title: string;
  showDisplayModeToggle: boolean;
  displayMode: DisplayMode;
  showRefreshButton: boolean;
  refreshing: boolean;
  onToggleDisplayMode: (() => void) | null;
  onRefresh: (() => void) | null;
  // Custom actions (e.g., "Add User" button)
  customActions: CustomAction[];
  // Scroll-based header state
  isScrolled: boolean;
  contentHeaderVisible: boolean;
}

const createDefaultState = (): PageHeaderState => ({
  title: '',
  showDisplayModeToggle: false,
  displayMode: 'standard',
  showRefreshButton: false,
  refreshing: false,
  onToggleDisplayMode: null,
  onRefresh: null,
  customActions: [],
  isScrolled: false,
  contentHeaderVisible: true,
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
    setCustomActions: (actions: CustomAction[]) => {
      update(state => ({ ...state, customActions: actions }));
    },
    setIsScrolled: (isScrolled: boolean) => {
      update(state => ({ ...state, isScrolled }));
    },
    setContentHeaderVisible: (visible: boolean) => {
      update(state => ({ ...state, contentHeaderVisible: visible }));
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
