import type { User, LoginCredentials, Container, ContainerAction, RegisterRequest, Volume } from '../types';
import { mockUsers, mockContainers, mockVolumes, mockToken } from './mockData';

// Mock API for development (no real backend required)
// In production, replace with real API calls

let users = [...mockUsers];
let containers = [...mockContainers];
let volumes = [...mockVolumes];

// Helper to simulate network delay
const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

export const mockAuthApi = {
  async login(credentials: LoginCredentials): Promise<{ user: User; token: string }> {
    await delay(500);
    
    const user = users.find(u => u.username === credentials.username && u.password === credentials.password);
    
    if (!user) {
      throw new Error('Invalid credentials');
    }
    
    const { password, ...userWithoutPassword } = user;
    return {
      user: userWithoutPassword,
      token: mockToken
    };
  },
  
  async register(data: RegisterRequest): Promise<User> {
    await delay(500);
    
    const existingUser = users.find(u => u.username === data.username);
    if (existingUser) {
      throw new Error('Username already exists');
    }
    
    const newUser = {
      id: users.length + 1,
      username: data.username,
      nickname: data.nickname,
      password: data.password,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    };
    
    users.push(newUser);
    
    const { password, ...userWithoutPassword } = newUser;
    return userWithoutPassword;
  },
  
  async logout(): Promise<void> {
    await delay(200);
  }
};

export const mockUserApi = {
  async getUsers(): Promise<User[]> {
    await delay(300);
    return users.map(({ password, ...user }) => user);
  },
  
  async createUser(user: Omit<User, 'id'>): Promise<User> {
    await delay(400);
    const newUser = {
      id: users.length + 1,
      username: user.username,
      nickname: user.nickname,
      password: (user as any).password || 'default123',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    };
    users.push(newUser);
    const { password, ...userWithoutPassword } = newUser;
    return userWithoutPassword;
  },
  
  async deleteUser(userId: string | number): Promise<void> {
    await delay(300);
    const id = typeof userId === 'string' ? parseInt(userId) : userId;
    users = users.filter(u => u.id !== id);
  },
  
  async updateUser(userId: string | number, userData: Partial<User>): Promise<User> {
    await delay(400);
    const id = typeof userId === 'string' ? parseInt(userId) : userId;
    const userIndex = users.findIndex(u => u.id === id);
    if (userIndex === -1) {
      throw new Error('User not found');
    }
    
    users[userIndex] = { 
      ...users[userIndex], 
      ...userData,
      updated_at: new Date().toISOString()
    };
    const { password, ...userWithoutPassword } = users[userIndex];
    return userWithoutPassword;
  }
};

export const mockContainerApi = {
  async getContainers(): Promise<Container[]> {
    await delay(400);
    return [...containers] as Container[];
  },
  
  async getContainer(id: string): Promise<Container> {
    await delay(300);
    const container = containers.find(c => c.id === id);
    if (!container) {
      throw new Error('Container not found');
    }
    return container as Container;
  },
  
  async controlContainer(action: ContainerAction): Promise<void> {
    await delay(600);
    const container = containers.find(c => c.id === action.containerId);
    
    if (!container) {
      throw new Error('Container not found');
    }
    
    switch (action.action) {
      case 'start':
        container.state = 'running';
        container.status = 'Up 1 second';
        break;
      case 'stop':
        container.state = 'exited';
        container.status = 'Exited (0) 1 second ago';
        break;
      case 'restart':
        container.state = 'running';
        container.status = 'Up 1 second';
        break;
    }
  }
};

export const mockVolumeApi = {
  async getVolumes(): Promise<Volume[]> {
    await delay(400);
    return [...volumes];
  }
};
