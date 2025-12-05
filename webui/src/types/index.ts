// API Response wrapper types
export interface ApiResponse<T = any> {
  success: boolean;
  message?: string;
  data?: T;
}

export interface ApiErrorResponse {
  success: false;
  error: string;
}

// User types
export interface User {
  id: number;
  username: string;
  nickname: string;
  created_at?: string;
  updated_at?: string;
}

export interface LoginCredentials {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  username: string;
  nickname: string;
}

export interface RegisterRequest {
  username: string;
  password: string;
  nickname: string;
}

// Docker container types
export interface Container {
  id: string;
  name: string;
  image: string;
  state: 'created' | 'running' | 'paused' | 'restarting' | 'removing' | 'exited' | 'dead';
  status: string; // Human-readable status like "Up 2 hours"
  health?: 'healthy' | 'unhealthy' | 'starting' | 'none';
  created: number; // Unix timestamp
  is_self?: boolean; // Whether this container is running the DSP application
}

// Docker volume types
export interface Volume {
  name: string;
  driver: string;
  mountpoint: string;
  created_at: string;
  scope: string;
  containers: string[]; // List of container IDs using this volume
}

export interface ContainerAction {
  containerId: string;
  action: 'start' | 'stop' | 'restart';
}

// Config types
export interface SystemConfig {
  docker_socket: string;
  disable_registration: boolean;
  log_level: string;
}

export interface PublicConfig {
  disable_registration: boolean;
}

export interface UpdateConfigRequest {
  docker_socket?: string;
  disable_registration?: boolean;
  log_level?: string;
}

// Auth types
export interface AuthState {
  isAuthenticated: boolean;
  user: User | null;
  token: string | null;
}
