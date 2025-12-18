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
  id?: number;
  username: string;
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
  compose_project?: string; // Docker Compose project name
  compose_service?: string; // Docker Compose service name
  labels?: Record<string, string>; // All container labels
  restart_policy?: RestartPolicy;
  env?: string[];
  networks?: Record<string, NetworkInfo>;
  ports?: PortBinding[];
  mounts?: MountInfo[];
  hostname?: string;
}

export interface RestartPolicy {
  name: string;
  maximum_retry_count?: number;
}

export interface NetworkInfo {
  network_id: string;
  gateway?: string;
  ip_address?: string;
  mac_address?: string;
}

export interface PortBinding {
  container_port: string;
  host_ip?: string;
  host_port?: string;
}

export interface MountInfo {
  type: string; // bind, volume, tmpfs
  source: string;
  destination: string;
  mode?: string;
  rw: boolean;
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

export interface VolumeFileInfo {
  name: string;
  path: string;
  is_directory: boolean;
  size: number;
  mode: string;
  mod_time: string;
}

export interface VolumeFileContent {
  path: string;
  content: string;
  size: number;
}

export interface ContainerAction {
  containerId: string;
  action: 'start' | 'stop' | 'restart';
}

// Config types
export interface SystemConfig {
  docker_socket: string;
  log_level: string;
  volume_explorer_image: string;
  session_max_timeout: number;
  username: string;
}

export interface UpdateConfigRequest {
  docker_socket?: string;
  log_level?: string;
  volume_explorer_image?: string;
  session_max_timeout?: number;
}

// Auth types
export interface AuthState {
  isAuthenticated: boolean;
  user: User | null;
  token: string | null;
}
