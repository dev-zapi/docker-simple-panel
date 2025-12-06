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
    created: 1705315800, // Unix timestamp for 2024-01-15T10:30:00Z
    compose_project: 'web-stack',
    compose_service: 'nginx',
    restart_policy: {
      name: 'always',
      maximum_retry_count: 0
    },
    env: [
      'PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin',
      'NGINX_VERSION=1.25.3',
      'NJS_VERSION=0.8.2',
      'PKG_RELEASE=1~bookworm',
      'APP_ENV=production',
      'LOG_LEVEL=info'
    ],
    networks: {
      'web-stack_default': {
        network_id: 'abc123def456',
        gateway: '172.20.0.1',
        ip_address: '172.20.0.2',
        mac_address: '02:42:ac:14:00:02'
      }
    },
    ports: [
      { container_port: '80/tcp', host_ip: '0.0.0.0', host_port: '8080' },
      { container_port: '80/tcp', host_ip: '::', host_port: '8080' },
      { container_port: '443/tcp', host_ip: '0.0.0.0', host_port: '8443' }
    ],
    mounts: [
      {
        type: 'volume',
        source: 'nginx-config',
        destination: '/etc/nginx/conf.d',
        mode: 'rw',
        rw: true
      },
      {
        type: 'bind',
        source: '/var/www/html',
        destination: '/usr/share/nginx/html',
        mode: 'ro',
        rw: false
      }
    ],
    hostname: 'nginx-web'
  },
  {
    id: 'container2',
    name: 'mysql-db',
    state: 'running' as const,
    status: 'Up 2 hours',
    health: 'healthy' as const,
    image: 'mysql:8.0',
    created: 1705315500, // Unix timestamp for 2024-01-15T10:25:00Z
    compose_project: 'web-stack',
    compose_service: 'mysql',
    restart_policy: {
      name: 'unless-stopped',
      maximum_retry_count: 0
    },
    env: [
      'PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin',
      'MYSQL_ROOT_PASSWORD=secretpassword',
      'MYSQL_DATABASE=appdb',
      'MYSQL_USER=appuser',
      'MYSQL_PASSWORD=apppass',
      'MYSQL_VERSION=8.0.35'
    ],
    networks: {
      'web-stack_default': {
        network_id: 'abc123def456',
        gateway: '172.20.0.1',
        ip_address: '172.20.0.3',
        mac_address: '02:42:ac:14:00:03'
      }
    },
    ports: [
      { container_port: '3306/tcp', host_ip: '127.0.0.1', host_port: '3306' }
    ],
    mounts: [
      {
        type: 'volume',
        source: 'mysql-data',
        destination: '/var/lib/mysql',
        mode: 'rw',
        rw: true
      }
    ],
    hostname: 'mysql-db'
  },
  {
    id: 'container3',
    name: 'redis-cache',
    state: 'exited' as const,
    status: 'Exited (0) 10 minutes ago',
    health: 'none' as const,
    image: 'redis:7',
    created: 1705315200, // Unix timestamp for 2024-01-15T10:20:00Z
    compose_project: 'web-stack',
    compose_service: 'redis',
    restart_policy: {
      name: 'on-failure',
      maximum_retry_count: 3
    },
    env: [
      'PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin',
      'REDIS_VERSION=7.2.3',
      'REDIS_REPLICATION_MODE=master'
    ],
    networks: {
      'web-stack_default': {
        network_id: 'abc123def456',
        gateway: '172.20.0.1',
        ip_address: '172.20.0.4',
        mac_address: '02:42:ac:14:00:04'
      }
    },
    ports: [
      { container_port: '6379/tcp' }
    ],
    mounts: [
      {
        type: 'volume',
        source: 'redis-data',
        destination: '/data',
        mode: 'rw',
        rw: true
      }
    ],
    hostname: 'redis-cache'
  },
  {
    id: 'container4',
    name: 'postgres-db',
    state: 'running' as const,
    status: 'Up 2 hours (health: starting)',
    health: 'starting' as const,
    image: 'postgres:15',
    created: 1705314900, // Unix timestamp for 2024-01-15T10:15:00Z
    compose_project: 'data-stack',
    compose_service: 'postgres',
    restart_policy: {
      name: 'always',
      maximum_retry_count: 0
    },
    env: [
      'PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin',
      'POSTGRES_VERSION=15.5',
      'POSTGRES_DB=maindb',
      'POSTGRES_USER=dbadmin',
      'POSTGRES_PASSWORD=dbsecret',
      'PGDATA=/var/lib/postgresql/data'
    ],
    networks: {
      'data-stack_default': {
        network_id: 'def456ghi789',
        gateway: '172.21.0.1',
        ip_address: '172.21.0.2',
        mac_address: '02:42:ac:15:00:02'
      }
    },
    ports: [
      { container_port: '5432/tcp', host_ip: '0.0.0.0', host_port: '5432' }
    ],
    mounts: [
      {
        type: 'volume',
        source: 'postgres-data',
        destination: '/var/lib/postgresql/data',
        mode: 'rw',
        rw: true
      },
      {
        type: 'bind',
        source: '/etc/localtime',
        destination: '/etc/localtime',
        mode: 'ro',
        rw: false
      }
    ],
    hostname: 'postgres-db'
  },
  {
    id: 'container5',
    name: 'mongodb-data',
    state: 'running' as const,
    status: 'Up 3 hours',
    health: 'healthy' as const,
    image: 'mongo:6.0',
    created: 1705314600, // Unix timestamp for 2024-01-15T10:10:00Z
    compose_project: 'data-stack',
    compose_service: 'mongodb',
    restart_policy: {
      name: 'unless-stopped',
      maximum_retry_count: 0
    },
    env: [
      'PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin',
      'MONGO_VERSION=6.0.12',
      'MONGO_INITDB_ROOT_USERNAME=root',
      'MONGO_INITDB_ROOT_PASSWORD=mongosecret',
      'MONGO_INITDB_DATABASE=admin'
    ],
    networks: {
      'data-stack_default': {
        network_id: 'def456ghi789',
        gateway: '172.21.0.1',
        ip_address: '172.21.0.3',
        mac_address: '02:42:ac:15:00:03'
      },
      'bridge': {
        network_id: 'bridge000',
        gateway: '172.17.0.1',
        ip_address: '172.17.0.5',
        mac_address: '02:42:ac:11:00:05'
      }
    },
    ports: [
      { container_port: '27017/tcp', host_ip: '0.0.0.0', host_port: '27017' }
    ],
    mounts: [
      {
        type: 'volume',
        source: 'mongodb-data',
        destination: '/data/db',
        mode: 'rw',
        rw: true
      },
      {
        type: 'volume',
        source: 'mongodb-config',
        destination: '/data/configdb',
        mode: 'rw',
        rw: true
      }
    ],
    hostname: 'mongodb-data'
  },
  {
    id: 'container6',
    name: 'elasticsearch-search',
    state: 'running' as const,
    status: 'Up 4 hours',
    health: 'healthy' as const,
    image: 'elasticsearch:8.11',
    created: 1705314300, // Unix timestamp for 2024-01-15T10:05:00Z
    restart_policy: {
      name: 'always',
      maximum_retry_count: 0
    },
    env: [
      'PATH=/usr/share/elasticsearch/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin',
      'ELASTIC_VERSION=8.11.3',
      'discovery.type=single-node',
      'xpack.security.enabled=false',
      'ES_JAVA_OPTS=-Xms512m -Xmx512m'
    ],
    networks: {
      'bridge': {
        network_id: 'bridge000',
        gateway: '172.17.0.1',
        ip_address: '172.17.0.6',
        mac_address: '02:42:ac:11:00:06'
      }
    },
    ports: [
      { container_port: '9200/tcp', host_ip: '0.0.0.0', host_port: '9200' },
      { container_port: '9300/tcp', host_ip: '0.0.0.0', host_port: '9300' }
    ],
    mounts: [
      {
        type: 'volume',
        source: 'elasticsearch-data',
        destination: '/usr/share/elasticsearch/data',
        mode: 'rw',
        rw: true
      }
    ],
    hostname: 'elasticsearch-search'
  },
  {
    id: 'container7',
    name: 'rabbitmq-queue',
    state: 'running' as const,
    status: 'Up 5 hours',
    health: 'healthy' as const,
    image: 'rabbitmq:3-management',
    created: 1705314000, // Unix timestamp for 2024-01-15T10:00:00Z
    compose_project: 'messaging',
    compose_service: 'rabbitmq',
    restart_policy: { name: 'always', maximum_retry_count: 0 },
    env: ['PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin', 'RABBITMQ_VERSION=3.12.10'],
    networks: { 'messaging_default': { network_id: 'msg123', gateway: '172.22.0.1', ip_address: '172.22.0.2', mac_address: '02:42:ac:16:00:02' } },
    ports: [{ container_port: '5672/tcp', host_ip: '0.0.0.0', host_port: '5672' }, { container_port: '15672/tcp', host_ip: '0.0.0.0', host_port: '15672' }],
    mounts: [{ type: 'volume', source: 'rabbitmq-data', destination: '/var/lib/rabbitmq', mode: 'rw', rw: true }],
    hostname: 'rabbitmq-queue'
  },
  {
    id: 'container8',
    name: 'grafana-monitor',
    state: 'running' as const,
    status: 'Up 6 hours',
    health: 'healthy' as const,
    image: 'grafana/grafana:10.0',
    created: 1705313700, // Unix timestamp for 2024-01-15T09:55:00Z
    compose_project: 'monitoring',
    compose_service: 'grafana',
    restart_policy: { name: 'unless-stopped', maximum_retry_count: 0 },
    env: ['PATH=/usr/share/grafana/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin', 'GF_SECURITY_ADMIN_PASSWORD=admin', 'GF_INSTALL_PLUGINS=grafana-clock-panel'],
    networks: { 'monitoring_default': { network_id: 'mon123', gateway: '172.23.0.1', ip_address: '172.23.0.2', mac_address: '02:42:ac:17:00:02' } },
    ports: [{ container_port: '3000/tcp', host_ip: '0.0.0.0', host_port: '3000' }],
    mounts: [{ type: 'volume', source: 'grafana-data', destination: '/var/lib/grafana', mode: 'rw', rw: true }],
    hostname: 'grafana-monitor'
  },
  {
    id: 'container9',
    name: 'prometheus-metrics',
    state: 'running' as const,
    status: 'Up 6 hours',
    health: 'healthy' as const,
    image: 'prom/prometheus:v2.47',
    created: 1705313400, // Unix timestamp for 2024-01-15T09:50:00Z
    compose_project: 'monitoring',
    compose_service: 'prometheus',
    restart_policy: { name: 'always', maximum_retry_count: 0 },
    env: ['PATH=/bin:/usr/local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin'],
    networks: { 'monitoring_default': { network_id: 'mon123', gateway: '172.23.0.1', ip_address: '172.23.0.3', mac_address: '02:42:ac:17:00:03' } },
    ports: [{ container_port: '9090/tcp', host_ip: '0.0.0.0', host_port: '9090' }],
    mounts: [{ type: 'volume', source: 'prometheus-data', destination: '/prometheus', mode: 'rw', rw: true }, { type: 'bind', source: '/etc/prometheus', destination: '/etc/prometheus', mode: 'ro', rw: false }],
    hostname: 'prometheus-metrics'
  },
  {
    id: 'container10',
    name: 'caddy-proxy',
    state: 'running' as const,
    status: 'Up 7 hours',
    health: 'healthy' as const,
    image: 'caddy:2.7',
    created: 1705313100, // Unix timestamp for 2024-01-15T09:45:00Z
    restart_policy: { name: 'always', maximum_retry_count: 0 },
    env: ['PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin', 'CADDY_VERSION=2.7.6'],
    networks: { 'bridge': { network_id: 'bridge000', gateway: '172.17.0.1', ip_address: '172.17.0.10', mac_address: '02:42:ac:11:00:0a' } },
    ports: [{ container_port: '80/tcp', host_ip: '0.0.0.0', host_port: '80' }, { container_port: '443/tcp', host_ip: '0.0.0.0', host_port: '443' }],
    mounts: [{ type: 'bind', source: '/srv/caddy/Caddyfile', destination: '/etc/caddy/Caddyfile', mode: 'ro', rw: false }, { type: 'volume', source: 'caddy-data', destination: '/data', mode: 'rw', rw: true }],
    hostname: 'caddy-proxy'
  },
  {
    id: 'container11',
    name: 'jenkins-ci',
    state: 'exited' as const,
    status: 'Exited (137) 30 minutes ago',
    health: 'none' as const,
    image: 'jenkins/jenkins:lts',
    created: 1705312800, // Unix timestamp for 2024-01-15T09:40:00Z
    restart_policy: { name: 'no', maximum_retry_count: 0 },
    env: ['PATH=/opt/java/openjdk/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin', 'JENKINS_VERSION=2.426.1'],
    networks: { 'bridge': { network_id: 'bridge000', gateway: '172.17.0.1', ip_address: '172.17.0.11', mac_address: '02:42:ac:11:00:0b' } },
    ports: [{ container_port: '8080/tcp', host_ip: '0.0.0.0', host_port: '8081' }],
    mounts: [{ type: 'volume', source: 'jenkins-home', destination: '/var/jenkins_home', mode: 'rw', rw: true }],
    hostname: 'jenkins-ci'
  },
  {
    id: 'container12',
    name: 'sonarqube-analysis',
    state: 'running' as const,
    status: 'Up 8 hours',
    health: 'healthy' as const,
    image: 'sonarqube:10-community',
    created: 1705312500, // Unix timestamp for 2024-01-15T09:35:00Z
    restart_policy: { name: 'unless-stopped', maximum_retry_count: 0 },
    env: ['PATH=/opt/java/openjdk/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin', 'SONAR_VERSION=10.3.0.82913'],
    networks: { 'bridge': { network_id: 'bridge000', gateway: '172.17.0.1', ip_address: '172.17.0.12', mac_address: '02:42:ac:11:00:0c' } },
    ports: [{ container_port: '9000/tcp', host_ip: '0.0.0.0', host_port: '9000' }],
    mounts: [{ type: 'volume', source: 'sonarqube-data', destination: '/opt/sonarqube/data', mode: 'rw', rw: true }],
    hostname: 'sonarqube-analysis'
  },
  {
    id: 'container13',
    name: 'vault-secrets',
    state: 'running' as const,
    status: 'Up 9 hours',
    health: 'healthy' as const,
    image: 'hashicorp/vault:1.15',
    created: 1705312200, // Unix timestamp for 2024-01-15T09:30:00Z
    restart_policy: { name: 'always', maximum_retry_count: 0 },
    env: ['PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin', 'VAULT_VERSION=1.15.4', 'VAULT_ADDR=http://127.0.0.1:8200'],
    networks: { 'bridge': { network_id: 'bridge000', gateway: '172.17.0.1', ip_address: '172.17.0.13', mac_address: '02:42:ac:11:00:0d' } },
    ports: [{ container_port: '8200/tcp', host_ip: '0.0.0.0', host_port: '8200' }],
    mounts: [{ type: 'volume', source: 'vault-file', destination: '/vault/file', mode: 'rw', rw: true }],
    hostname: 'vault-secrets'
  },
  {
    id: 'container14',
    name: 'consul-discovery',
    state: 'running' as const,
    status: 'Up 10 hours',
    health: 'healthy' as const,
    image: 'hashicorp/consul:1.17',
    created: 1705311900, // Unix timestamp for 2024-01-15T09:25:00Z
    restart_policy: { name: 'always', maximum_retry_count: 0 },
    env: ['PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin', 'CONSUL_VERSION=1.17.1'],
    networks: { 'bridge': { network_id: 'bridge000', gateway: '172.17.0.1', ip_address: '172.17.0.14', mac_address: '02:42:ac:11:00:0e' } },
    ports: [{ container_port: '8500/tcp', host_ip: '0.0.0.0', host_port: '8500' }, { container_port: '8600/tcp', host_ip: '0.0.0.0', host_port: '8600' }],
    mounts: [{ type: 'volume', source: 'consul-data', destination: '/consul/data', mode: 'rw', rw: true }],
    hostname: 'consul-discovery'
  },
  {
    id: 'container15',
    name: 'minio-storage',
    state: 'paused' as const,
    status: 'Up 11 hours (Paused)',
    health: 'none' as const,
    image: 'minio/minio:latest',
    created: 1705311600, // Unix timestamp for 2024-01-15T09:20:00Z
    restart_policy: { name: 'unless-stopped', maximum_retry_count: 0 },
    env: ['PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin', 'MINIO_ROOT_USER=minioadmin', 'MINIO_ROOT_PASSWORD=minioadmin'],
    networks: { 'bridge': { network_id: 'bridge000', gateway: '172.17.0.1', ip_address: '172.17.0.15', mac_address: '02:42:ac:11:00:0f' } },
    ports: [{ container_port: '9000/tcp', host_ip: '0.0.0.0', host_port: '9001' }, { container_port: '9001/tcp', host_ip: '0.0.0.0', host_port: '9002' }],
    mounts: [{ type: 'volume', source: 'minio-data', destination: '/data', mode: 'rw', rw: true }],
    hostname: 'minio-storage'
  },
  {
    id: 'container16',
    name: 'keycloak-auth',
    state: 'running' as const,
    status: 'Up 12 hours',
    health: 'healthy' as const,
    image: 'quay.io/keycloak/keycloak:23.0',
    created: 1705311300, // Unix timestamp for 2024-01-15T09:15:00Z
    restart_policy: { name: 'always', maximum_retry_count: 0 },
    env: ['PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin', 'KEYCLOAK_ADMIN=admin', 'KEYCLOAK_ADMIN_PASSWORD=admin', 'KC_DB=postgres'],
    networks: { 'bridge': { network_id: 'bridge000', gateway: '172.17.0.1', ip_address: '172.17.0.16', mac_address: '02:42:ac:11:00:10' } },
    ports: [{ container_port: '8080/tcp', host_ip: '0.0.0.0', host_port: '8082' }],
    mounts: [{ type: 'volume', source: 'keycloak-data', destination: '/opt/keycloak/data', mode: 'rw', rw: true }],
    hostname: 'keycloak-auth'
  },
  {
    id: 'container17',
    name: 'traefik-router',
    state: 'running' as const,
    status: 'Up 1 day',
    health: 'healthy' as const,
    image: 'traefik:v3.0',
    created: 1705311000, // Unix timestamp for 2024-01-15T09:10:00Z
    restart_policy: { name: 'always', maximum_retry_count: 0 },
    env: ['PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin', 'TRAEFIK_VERSION=v3.0.0'],
    networks: { 'bridge': { network_id: 'bridge000', gateway: '172.17.0.1', ip_address: '172.17.0.17', mac_address: '02:42:ac:11:00:11' } },
    ports: [{ container_port: '80/tcp', host_ip: '0.0.0.0', host_port: '8083' }, { container_port: '8080/tcp', host_ip: '0.0.0.0', host_port: '8084' }],
    mounts: [{ type: 'bind', source: '/var/run/docker.sock', destination: '/var/run/docker.sock', mode: 'ro', rw: false }],
    hostname: 'traefik-router'
  },
  {
    id: 'container18',
    name: 'portainer-ui',
    state: 'running' as const,
    status: 'Up 2 days',
    health: 'healthy' as const,
    image: 'portainer/portainer-ce:2.19',
    created: 1705310700, // Unix timestamp for 2024-01-15T09:05:00Z
    restart_policy: { name: 'always', maximum_retry_count: 0 },
    env: ['PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin'],
    networks: { 'bridge': { network_id: 'bridge000', gateway: '172.17.0.1', ip_address: '172.17.0.18', mac_address: '02:42:ac:11:00:12' } },
    ports: [{ container_port: '9000/tcp', host_ip: '0.0.0.0', host_port: '9003' }],
    mounts: [{ type: 'bind', source: '/var/run/docker.sock', destination: '/var/run/docker.sock', mode: 'rw', rw: true }, { type: 'volume', source: 'portainer-data', destination: '/data', mode: 'rw', rw: true }],
    hostname: 'portainer-ui'
  },
  {
    id: 'container19',
    name: 'gitea-code',
    state: 'exited' as const,
    status: 'Exited (0) 1 hour ago',
    health: 'none' as const,
    image: 'gitea/gitea:1.21',
    created: 1705310400, // Unix timestamp for 2024-01-15T09:00:00Z
    restart_policy: { name: 'unless-stopped', maximum_retry_count: 0 },
    env: ['PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin', 'USER_UID=1000', 'USER_GID=1000'],
    networks: { 'bridge': { network_id: 'bridge000', gateway: '172.17.0.1', ip_address: '172.17.0.19', mac_address: '02:42:ac:11:00:13' } },
    ports: [{ container_port: '3000/tcp', host_ip: '0.0.0.0', host_port: '3001' }, { container_port: '22/tcp', host_ip: '0.0.0.0', host_port: '2222' }],
    mounts: [{ type: 'volume', source: 'gitea-data', destination: '/data', mode: 'rw', rw: true }],
    hostname: 'gitea-code'
  },
  {
    id: 'container20',
    name: 'drone-runner',
    state: 'restarting' as const,
    status: 'Restarting (1) 5 seconds ago',
    health: 'unhealthy' as const,
    image: 'drone/drone:2.21',
    created: 1705310100, // Unix timestamp for 2024-01-15T08:55:00Z
    restart_policy: { name: 'on-failure', maximum_retry_count: 5 },
    env: ['PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin', 'DRONE_SERVER_HOST=drone.example.com', 'DRONE_SERVER_PROTO=https'],
    networks: { 'bridge': { network_id: 'bridge000', gateway: '172.17.0.1', ip_address: '172.17.0.20', mac_address: '02:42:ac:11:00:14' } },
    ports: [{ container_port: '80/tcp', host_ip: '0.0.0.0', host_port: '8085' }],
    mounts: [{ type: 'bind', source: '/var/run/docker.sock', destination: '/var/run/docker.sock', mode: 'rw', rw: true }],
    hostname: 'drone-runner'
  }
];

export const mockVolumes = [
  // Database volumes
  {
    name: 'mysql-data',
    driver: 'local',
    mountpoint: '/var/lib/docker/volumes/mysql-data/_data',
    created_at: '2024-01-15T08:00:00Z',
    scope: 'local',
    containers: ['container2']
  },
  {
    name: 'postgres-data',
    driver: 'local',
    mountpoint: '/var/lib/docker/volumes/postgres-data/_data',
    created_at: '2024-01-15T08:15:00Z',
    scope: 'local',
    containers: ['container4']
  },
  {
    name: 'mongodb-data',
    driver: 'local',
    mountpoint: '/var/lib/docker/volumes/mongodb-data/_data',
    created_at: '2024-01-15T08:30:00Z',
    scope: 'local',
    containers: ['container5']
  },
  {
    name: 'elasticsearch-data',
    driver: 'local',
    mountpoint: '/var/lib/docker/volumes/elasticsearch-data/_data',
    created_at: '2024-01-15T08:45:00Z',
    scope: 'local',
    containers: ['container6']
  },
  
  // Message queue and cache volumes
  {
    name: 'rabbitmq-data',
    driver: 'local',
    mountpoint: '/var/lib/docker/volumes/rabbitmq-data/_data',
    created_at: '2024-01-15T09:00:00Z',
    scope: 'local',
    containers: ['container7']
  },
  
  // Monitoring and observability volumes
  {
    name: 'grafana-data',
    driver: 'local',
    mountpoint: '/var/lib/docker/volumes/grafana-data/_data',
    created_at: '2024-01-15T09:15:00Z',
    scope: 'local',
    containers: ['container8']
  },
  {
    name: 'prometheus-data',
    driver: 'local',
    mountpoint: '/var/lib/docker/volumes/prometheus-data/_data',
    created_at: '2024-01-15T09:30:00Z',
    scope: 'local',
    containers: ['container9']
  },
  {
    name: 'shared-logs',
    driver: 'local',
    mountpoint: '/var/lib/docker/volumes/shared-logs/_data',
    created_at: '2024-01-15T07:00:00Z',
    scope: 'local',
    containers: ['container1', 'container8', 'container9'] // Multi-container volume
  },
  
  // CI/CD and development volumes
  {
    name: 'jenkins-home',
    driver: 'local',
    mountpoint: '/var/lib/docker/volumes/jenkins-home/_data',
    created_at: '2024-01-15T09:45:00Z',
    scope: 'local',
    containers: ['container11']
  },
  {
    name: 'sonarqube-data',
    driver: 'local',
    mountpoint: '/var/lib/docker/volumes/sonarqube-data/_data',
    created_at: '2024-01-15T10:00:00Z',
    scope: 'local',
    containers: ['container12']
  },
  {
    name: 'gitea-data',
    driver: 'local',
    mountpoint: '/var/lib/docker/volumes/gitea-data/_data',
    created_at: '2024-01-15T11:15:00Z',
    scope: 'local',
    containers: ['container19']
  },
  
  // Infrastructure volumes
  {
    name: 'vault-file',
    driver: 'local',
    mountpoint: '/var/lib/docker/volumes/vault-file/_data',
    created_at: '2024-01-15T10:15:00Z',
    scope: 'local',
    containers: ['container13']
  },
  {
    name: 'consul-data',
    driver: 'local',
    mountpoint: '/var/lib/docker/volumes/consul-data/_data',
    created_at: '2024-01-15T10:30:00Z',
    scope: 'local',
    containers: ['container14']
  },
  {
    name: 'minio-data',
    driver: 'local',
    mountpoint: '/var/lib/docker/volumes/minio-data/_data',
    created_at: '2024-01-15T10:45:00Z',
    scope: 'local',
    containers: ['container15']
  },
  {
    name: 'keycloak-data',
    driver: 'local',
    mountpoint: '/var/lib/docker/volumes/keycloak-data/_data',
    created_at: '2024-01-15T11:00:00Z',
    scope: 'local',
    containers: ['container16']
  },
  
  // Unused volumes (for testing edge cases)
  {
    name: 'unused-volume',
    driver: 'local',
    mountpoint: '/var/lib/docker/volumes/unused-volume/_data',
    created_at: '2024-01-14T15:00:00Z',
    scope: 'local',
    containers: []
  },
  {
    name: 'backup-volume',
    driver: 'local',
    mountpoint: '/var/lib/docker/volumes/backup-volume/_data',
    created_at: '2024-01-10T08:00:00Z',
    scope: 'local',
    containers: []
  }
];

// Mock JWT token
export const mockToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.mock.token';
