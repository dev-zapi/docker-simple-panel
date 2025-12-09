import type { 
  User, 
  LoginCredentials, 
  LoginResponse,
  RegisterRequest,
  Container, 
  ContainerAction,
  Volume,
  VolumeFileInfo,
  VolumeFileContent,
  SystemConfig,
  PublicConfig,
  UpdateConfigRequest,
  ApiResponse,
  ApiErrorResponse
} from '../types';
import { mockAuthApi, mockUserApi, mockContainerApi, mockVolumeApi } from './mockApi';

// Use relative path for API calls to support reverse proxy
// In development, Vite will proxy /api to backend server
// In production, nginx will proxy /api to backend server
const API_BASE_URL = import.meta.env.VITE_API_URL || '/api';
const USE_MOCK_API = import.meta.env.VITE_USE_MOCK_API !== 'false'; // Use mock by default

// Helper function to get auth headers
const getAuthHeaders = (): HeadersInit => {
  const token = localStorage.getItem('token');
  return {
    'Content-Type': 'application/json',
    ...(token ? { 'Authorization': `Bearer ${token}` } : {})
  };
};

// Helper to handle API responses
async function handleApiResponse<T>(response: Response): Promise<T> {
  const contentType = response.headers.get('content-type');
  
  if (!response.ok) {
    if (contentType?.includes('application/json')) {
      const errorData = await response.json() as ApiErrorResponse;
      throw new Error(errorData.error || 'Request failed');
    }
    throw new Error(`Request failed with status ${response.status}`);
  }
  
  if (contentType?.includes('application/json')) {
    const apiResponse = await response.json() as ApiResponse<T>;
    if (apiResponse.success && apiResponse.data !== undefined) {
      return apiResponse.data;
    }
    // Some endpoints may return success without data field
    if (apiResponse.success) {
      return apiResponse as unknown as T;
    }
    throw new Error(apiResponse.message || 'Request failed');
  }
  
  throw new Error('Invalid response format');
}

// Auth API
export const authApi = USE_MOCK_API ? mockAuthApi : {
  async login(credentials: LoginCredentials): Promise<{ user: User; token: string }> {
    const response = await fetch(`${API_BASE_URL}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(credentials)
    });
    
    const loginResponse = await handleApiResponse<LoginResponse>(response);
    
    // Convert LoginResponse to expected format with User object
    // Backend doesn't return user ID in login response, using -1 as placeholder
    const user: User = {
      id: -1,
      username: loginResponse.username,
      nickname: loginResponse.nickname
    };
    
    return {
      user,
      token: loginResponse.token
    };
  },
  
  async register(data: RegisterRequest): Promise<User> {
    const response = await fetch(`${API_BASE_URL}/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    });
    
    return handleApiResponse<User>(response);
  },
  
  async logout(): Promise<void> {
    // Backend doesn't have logout endpoint
    // Just clear local storage
  }
};

// User API - Backend has user management endpoints at /api/users
export const userApi = USE_MOCK_API ? mockUserApi : {
  async getUsers(): Promise<User[]> {
    const response = await fetch(`${API_BASE_URL}/users`, {
      headers: getAuthHeaders()
    });
    
    return handleApiResponse<User[]>(response);
  },
  
  async createUser(user: Omit<User, 'id'> & { password: string }): Promise<User> {
    const response = await fetch(`${API_BASE_URL}/users`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify({
        username: user.username,
        password: user.password,
        nickname: user.nickname
      })
    });
    
    return handleApiResponse<User>(response);
  },
  
  async deleteUser(userId: string | number): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/users/${userId}`, {
      method: 'DELETE',
      headers: getAuthHeaders()
    });
    
    await handleApiResponse<void>(response);
  },
  
  async updateUser(userId: string | number, user: Partial<User>): Promise<User> {
    // Backend doesn't support user updates yet
    throw new Error('User update not supported by backend');
  }
};

// Container API
export const containerApi = USE_MOCK_API ? mockContainerApi : {
  async getContainers(): Promise<Container[]> {
    const response = await fetch(`${API_BASE_URL}/containers`, {
      headers: getAuthHeaders()
    });
    
    return handleApiResponse<Container[]>(response);
  },
  
  async getContainer(id: string): Promise<Container> {
    const response = await fetch(`${API_BASE_URL}/containers/${id}`, {
      headers: getAuthHeaders()
    });
    
    return handleApiResponse<Container>(response);
  },
  
  async controlContainer(action: ContainerAction): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/containers/${action.containerId}/${action.action}`, {
      method: 'POST',
      headers: getAuthHeaders()
    });
    
    await handleApiResponse<void>(response);
  }
};

// Volume API
export const volumeApi = USE_MOCK_API ? mockVolumeApi : {
  async getVolumes(): Promise<Volume[]> {
    const response = await fetch(`${API_BASE_URL}/volumes`, {
      headers: getAuthHeaders()
    });
    
    return handleApiResponse<Volume[]>(response);
  },
  
  async exploreVolumeFiles(volumeName: string, path: string = '/'): Promise<VolumeFileInfo[]> {
    const encodedPath = encodeURIComponent(path);
    const response = await fetch(`${API_BASE_URL}/volumes/${volumeName}/files?path=${encodedPath}`, {
      headers: getAuthHeaders()
    });
    
    return handleApiResponse<VolumeFileInfo[]>(response);
  },
  
  async readVolumeFile(volumeName: string, filePath: string): Promise<VolumeFileContent> {
    const encodedPath = encodeURIComponent(filePath);
    const response = await fetch(`${API_BASE_URL}/volumes/${volumeName}/file?path=${encodedPath}`, {
      headers: getAuthHeaders()
    });
    
    return handleApiResponse<VolumeFileContent>(response);
  },
  
  async deleteVolumeFile(volumeName: string, filePath: string): Promise<void> {
    const encodedPath = encodeURIComponent(filePath);
    const response = await fetch(`${API_BASE_URL}/volumes/${volumeName}/file?path=${encodedPath}`, {
      method: 'DELETE',
      headers: getAuthHeaders()
    });
    
    await handleApiResponse<void>(response);
  },
  
  async deleteVolume(volumeName: string): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/volumes/${encodeURIComponent(volumeName)}`, {
      method: 'DELETE',
      headers: getAuthHeaders()
    });
    
    await handleApiResponse<void>(response);
  }
};

// Config API
export const configApi = {
  async getConfig(): Promise<SystemConfig> {
    const response = await fetch(`${API_BASE_URL}/config`, {
      headers: getAuthHeaders()
    });
    
    return handleApiResponse<SystemConfig>(response);
  },
  
  // Public endpoint to check if registration is enabled (no auth required)
  // Note: Callers should handle errors - consider defaulting to allow registration on failure
  async getPublicConfig(): Promise<PublicConfig> {
    const response = await fetch(`${API_BASE_URL}/config/public`, {
      headers: { 'Content-Type': 'application/json' }
    });
    
    return handleApiResponse<PublicConfig>(response);
  },
  
  async updateConfig(config: UpdateConfigRequest): Promise<SystemConfig> {
    const response = await fetch(`${API_BASE_URL}/config`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(config)
    });
    
    return handleApiResponse<SystemConfig>(response);
  },
  
  async patchConfig(config: UpdateConfigRequest): Promise<SystemConfig> {
    const response = await fetch(`${API_BASE_URL}/config`, {
      method: 'PATCH',
      headers: getAuthHeaders(),
      body: JSON.stringify(config)
    });
    
    return handleApiResponse<SystemConfig>(response);
  }
};

// Health API
export const healthApi = {
  async checkHealth(): Promise<{ status: string }> {
    const response = await fetch(`${API_BASE_URL}/health`);
    return handleApiResponse(response);
  },
  
  async checkDockerHealth(): Promise<{ status: string }> {
    const response = await fetch(`${API_BASE_URL}/docker/health`, {
      headers: getAuthHeaders()
    });
    return handleApiResponse(response);
  }
};
