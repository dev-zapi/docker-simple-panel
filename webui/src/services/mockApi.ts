import type { User, LoginCredentials, Container, ContainerAction, RegisterRequest, Volume, VolumeFileInfo, VolumeFileContent } from '../types';
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
  },
  
  async exploreVolumeFiles(volumeName: string, path: string = '/'): Promise<VolumeFileInfo[]> {
    await delay(500);
    
    // Mock file structure for different volumes and paths
    const mockFileStructure: Record<string, Record<string, VolumeFileInfo[]>> = {
      'mysql-data': {
        '/': [
          { name: 'ibdata1', path: '/ibdata1', is_directory: false, size: 12582912, mode: '-rw-r-----', mod_time: '2024-01-15 10:30:00 +0000' },
          { name: 'mysql', path: '/mysql', is_directory: true, size: 0, mode: 'drwxr-x---', mod_time: '2024-01-15 10:30:00 +0000' },
          { name: 'performance_schema', path: '/performance_schema', is_directory: true, size: 0, mode: 'drwxr-x---', mod_time: '2024-01-15 10:30:00 +0000' },
          { name: 'appdb', path: '/appdb', is_directory: true, size: 0, mode: 'drwxr-x---', mod_time: '2024-01-15 10:30:00 +0000' }
        ],
        '/mysql': [
          { name: 'user.frm', path: '/mysql/user.frm', is_directory: false, size: 10816, mode: '-rw-r-----', mod_time: '2024-01-15 10:30:00 +0000' },
          { name: 'db.frm', path: '/mysql/db.frm', is_directory: false, size: 9582, mode: '-rw-r-----', mod_time: '2024-01-15 10:30:00 +0000' }
        ],
        '/appdb': [
          { name: 'users.ibd', path: '/appdb/users.ibd', is_directory: false, size: 98304, mode: '-rw-r-----', mod_time: '2024-01-15 10:30:00 +0000' },
          { name: 'posts.ibd', path: '/appdb/posts.ibd', is_directory: false, size: 131072, mode: '-rw-r-----', mod_time: '2024-01-15 10:30:00 +0000' }
        ]
      },
      'nginx-config': {
        '/': [
          { name: 'default.conf', path: '/default.conf', is_directory: false, size: 1024, mode: '-rw-r--r--', mod_time: '2024-01-15 09:00:00 +0000' },
          { name: 'ssl', path: '/ssl', is_directory: true, size: 0, mode: 'drwxr-xr-x', mod_time: '2024-01-15 09:00:00 +0000' },
          { name: 'sites-enabled', path: '/sites-enabled', is_directory: true, size: 0, mode: 'drwxr-xr-x', mod_time: '2024-01-15 09:00:00 +0000' }
        ],
        '/ssl': [
          { name: 'server.crt', path: '/ssl/server.crt', is_directory: false, size: 2048, mode: '-rw-r--r--', mod_time: '2024-01-15 09:00:00 +0000' },
          { name: 'server.key', path: '/ssl/server.key', is_directory: false, size: 1704, mode: '-rw-------', mod_time: '2024-01-15 09:00:00 +0000' }
        ]
      },
      'unused-volume': {
        '/': [
          { name: 'README.txt', path: '/README.txt', is_directory: false, size: 256, mode: '-rw-r--r--', mod_time: '2024-01-14 15:00:00 +0000' },
          { name: 'data', path: '/data', is_directory: true, size: 0, mode: 'drwxr-xr-x', mod_time: '2024-01-14 15:00:00 +0000' },
          { name: 'config.json', path: '/config.json', is_directory: false, size: 512, mode: '-rw-r--r--', mod_time: '2024-01-14 15:00:00 +0000' }
        ],
        '/data': [
          { name: 'sample.txt', path: '/data/sample.txt', is_directory: false, size: 128, mode: '-rw-r--r--', mod_time: '2024-01-14 15:00:00 +0000' },
          { name: 'backup.tar.gz', path: '/data/backup.tar.gz', is_directory: false, size: 1048576, mode: '-rw-r--r--', mod_time: '2024-01-14 15:00:00 +0000' }
        ]
      }
    };
    
    // Return default structure if volume not in mock data
    const volumeFiles = mockFileStructure[volumeName] || {
      '/': [
        { name: 'data', path: '/data', is_directory: true, size: 0, mode: 'drwxr-xr-x', mod_time: '2024-01-15 12:00:00 +0000' },
        { name: 'config', path: '/config', is_directory: true, size: 0, mode: 'drwxr-xr-x', mod_time: '2024-01-15 12:00:00 +0000' },
        { name: 'README.md', path: '/README.md', is_directory: false, size: 1024, mode: '-rw-r--r--', mod_time: '2024-01-15 12:00:00 +0000' }
      ],
      '/data': [
        { name: 'file1.txt', path: '/data/file1.txt', is_directory: false, size: 256, mode: '-rw-r--r--', mod_time: '2024-01-15 12:00:00 +0000' }
      ],
      '/config': [
        { name: 'settings.json', path: '/config/settings.json', is_directory: false, size: 512, mode: '-rw-r--r--', mod_time: '2024-01-15 12:00:00 +0000' }
      ]
    };
    
    const files = volumeFiles[path] || [];
    
    if (files.length === 0) {
      throw new Error('Directory not found or empty');
    }
    
    return files;
  },
  
  async readVolumeFile(volumeName: string, filePath: string): Promise<VolumeFileContent> {
    await delay(400);
    
    // Mock file contents
    const mockFileContents: Record<string, Record<string, string>> = {
      'mysql-data': {
        '/mysql/user.frm': '-- MySQL User table schema\n-- Binary data (not readable as text)',
        '/appdb/users.ibd': '-- MySQL InnoDB data file\n-- Binary data (not readable as text)'
      },
      'nginx-config': {
        '/default.conf': `server {
    listen 80;
    server_name localhost;
    
    location / {
        root /usr/share/nginx/html;
        index index.html index.htm;
    }
    
    error_page 500 502 503 504 /50x.html;
    location = /50x.html {
        root /usr/share/nginx/html;
    }
}`,
        '/ssl/server.crt': `-----BEGIN CERTIFICATE-----
MIIDXTCCAkWgAwIBAgIJAKL0UG+mRzKhMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMjQwMTE1MDAwMDAwWhcNMjUwMTE1MDAwMDAwWjBF
...
-----END CERTIFICATE-----`
      },
      'unused-volume': {
        '/README.txt': `This is an unused Docker volume for testing purposes.

You can use this volume to test the volume explorer feature.
It contains sample files and directories.`,
        '/config.json': `{
  "version": "1.0.0",
  "name": "test-app",
  "database": {
    "host": "localhost",
    "port": 3306,
    "name": "testdb"
  },
  "features": {
    "logging": true,
    "monitoring": true,
    "caching": false
  }
}`,
        '/data/sample.txt': `Sample text file content.
This is line 2.
This is line 3.

End of file.`,
        '/data/backup.tar.gz': '-- Binary file (gzip compressed tar archive)\n-- Not readable as text'
      }
    };
    
    // Default content for unknown files
    const volumeContents = mockFileContents[volumeName] || {};
    const content = volumeContents[filePath] || `This is a sample file from ${volumeName}.
Path: ${filePath}

This is mock content for testing the volume file explorer.
You can see the file structure and read file contents.`;
    
    return {
      path: filePath,
      content: content,
      size: content.length
    };
  }
};
