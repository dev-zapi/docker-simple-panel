// Mock data for development
export const mockUsers = [
  { id: 1, username: 'admin', nickname: 'Administrator', password: 'admin123', created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
  { id: 2, username: 'user1', nickname: 'User One', password: 'pass123', created_at: '2024-01-02T00:00:00Z', updated_at: '2024-01-02T00:00:00Z' },
  { id: 3, username: 'user2', nickname: 'User Two', password: 'pass123', created_at: '2024-01-03T00:00:00Z', updated_at: '2024-01-03T00:00:00Z' }
];

export const mockContainers = [
  {
    id: 'container1',
    name: 'nginx-web',
    state: 'running' as const,
    status: 'Up 2 hours',
    health: 'healthy' as const,
    image: 'nginx:latest',
    created: 1705315800 // Unix timestamp for 2024-01-15T10:30:00Z
  },
  {
    id: 'container2',
    name: 'mysql-db',
    state: 'running' as const,
    status: 'Up 2 hours',
    health: 'healthy' as const,
    image: 'mysql:8.0',
    created: 1705315500 // Unix timestamp for 2024-01-15T10:25:00Z
  },
  {
    id: 'container3',
    name: 'redis-cache',
    state: 'exited' as const,
    status: 'Exited (0) 10 minutes ago',
    health: 'none' as const,
    image: 'redis:7',
    created: 1705315200 // Unix timestamp for 2024-01-15T10:20:00Z
  },
  {
    id: 'container4',
    name: 'postgres-db',
    state: 'running' as const,
    status: 'Up 2 hours (health: starting)',
    health: 'starting' as const,
    image: 'postgres:15',
    created: 1705314900 // Unix timestamp for 2024-01-15T10:15:00Z
  }
];

// Mock JWT token
export const mockToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.mock.token';
