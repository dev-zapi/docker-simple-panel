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
  },
  {
    id: 'container5',
    name: 'mongodb-data',
    state: 'running' as const,
    status: 'Up 3 hours',
    health: 'healthy' as const,
    image: 'mongo:6.0',
    created: 1705314600 // Unix timestamp for 2024-01-15T10:10:00Z
  },
  {
    id: 'container6',
    name: 'elasticsearch-search',
    state: 'running' as const,
    status: 'Up 4 hours',
    health: 'healthy' as const,
    image: 'elasticsearch:8.11',
    created: 1705314300 // Unix timestamp for 2024-01-15T10:05:00Z
  },
  {
    id: 'container7',
    name: 'rabbitmq-queue',
    state: 'running' as const,
    status: 'Up 5 hours',
    health: 'healthy' as const,
    image: 'rabbitmq:3-management',
    created: 1705314000 // Unix timestamp for 2024-01-15T10:00:00Z
  },
  {
    id: 'container8',
    name: 'grafana-monitor',
    state: 'running' as const,
    status: 'Up 6 hours',
    health: 'healthy' as const,
    image: 'grafana/grafana:10.0',
    created: 1705313700 // Unix timestamp for 2024-01-15T09:55:00Z
  },
  {
    id: 'container9',
    name: 'prometheus-metrics',
    state: 'running' as const,
    status: 'Up 6 hours',
    health: 'healthy' as const,
    image: 'prom/prometheus:v2.47',
    created: 1705313400 // Unix timestamp for 2024-01-15T09:50:00Z
  },
  {
    id: 'container10',
    name: 'caddy-proxy',
    state: 'running' as const,
    status: 'Up 7 hours',
    health: 'healthy' as const,
    image: 'caddy:2.7',
    created: 1705313100 // Unix timestamp for 2024-01-15T09:45:00Z
  },
  {
    id: 'container11',
    name: 'jenkins-ci',
    state: 'exited' as const,
    status: 'Exited (137) 30 minutes ago',
    health: 'none' as const,
    image: 'jenkins/jenkins:lts',
    created: 1705312800 // Unix timestamp for 2024-01-15T09:40:00Z
  },
  {
    id: 'container12',
    name: 'sonarqube-analysis',
    state: 'running' as const,
    status: 'Up 8 hours',
    health: 'healthy' as const,
    image: 'sonarqube:10-community',
    created: 1705312500 // Unix timestamp for 2024-01-15T09:35:00Z
  },
  {
    id: 'container13',
    name: 'vault-secrets',
    state: 'running' as const,
    status: 'Up 9 hours',
    health: 'healthy' as const,
    image: 'hashicorp/vault:1.15',
    created: 1705312200 // Unix timestamp for 2024-01-15T09:30:00Z
  },
  {
    id: 'container14',
    name: 'consul-discovery',
    state: 'running' as const,
    status: 'Up 10 hours',
    health: 'healthy' as const,
    image: 'hashicorp/consul:1.17',
    created: 1705311900 // Unix timestamp for 2024-01-15T09:25:00Z
  },
  {
    id: 'container15',
    name: 'minio-storage',
    state: 'paused' as const,
    status: 'Up 11 hours (Paused)',
    health: 'none' as const,
    image: 'minio/minio:latest',
    created: 1705311600 // Unix timestamp for 2024-01-15T09:20:00Z
  },
  {
    id: 'container16',
    name: 'keycloak-auth',
    state: 'running' as const,
    status: 'Up 12 hours',
    health: 'healthy' as const,
    image: 'quay.io/keycloak/keycloak:23.0',
    created: 1705311300 // Unix timestamp for 2024-01-15T09:15:00Z
  },
  {
    id: 'container17',
    name: 'traefik-router',
    state: 'running' as const,
    status: 'Up 1 day',
    health: 'healthy' as const,
    image: 'traefik:v3.0',
    created: 1705311000 // Unix timestamp for 2024-01-15T09:10:00Z
  },
  {
    id: 'container18',
    name: 'portainer-ui',
    state: 'running' as const,
    status: 'Up 2 days',
    health: 'healthy' as const,
    image: 'portainer/portainer-ce:2.19',
    created: 1705310700 // Unix timestamp for 2024-01-15T09:05:00Z
  },
  {
    id: 'container19',
    name: 'gitea-code',
    state: 'exited' as const,
    status: 'Exited (0) 1 hour ago',
    health: 'none' as const,
    image: 'gitea/gitea:1.21',
    created: 1705310400 // Unix timestamp for 2024-01-15T09:00:00Z
  },
  {
    id: 'container20',
    name: 'drone-runner',
    state: 'restarting' as const,
    status: 'Restarting (1) 5 seconds ago',
    health: 'unhealthy' as const,
    image: 'drone/drone:2.21',
    created: 1705310100 // Unix timestamp for 2024-01-15T08:55:00Z
  }
];

// Mock JWT token
export const mockToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.mock.token';
