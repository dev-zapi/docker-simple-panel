<script lang="ts">
  import Router from 'svelte-spa-router';
  import { get } from 'svelte/store';
  import { isAuthenticated } from './stores/authStore';
  
  import Login from './pages/Login.svelte';
  import Register from './pages/Register.svelte';
  import Home from './pages/Home.svelte';
  import Volumes from './pages/Volumes.svelte';
  import Users from './pages/Users.svelte';
  import Profile from './pages/Profile.svelte';
  import Settings from './pages/Settings.svelte';
  
  // Define routes - using direct component assignment for compatibility
  const routes = {
    '/': Home,
    '/login': Login,
    '/register': Register,
    '/volumes': Volumes,
    '/users': Users,
    '/profile': Profile,
    '/settings': Settings
  };
  
  // Public routes that don't require authentication
  const publicRoutes = ['#/login', '#/register'];
  
  // Route guard - redirect to login if not authenticated
  function conditionsFailed() {
    window.location.hash = '/login';
  }
  
  // Check authentication on route change
  $: if (!publicRoutes.includes(window.location.hash) && !get(isAuthenticated)) {
    window.location.hash = '/login';
  }
</script>

<Router {routes} />

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
      'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
      sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }
  
  :global(*) {
    box-sizing: border-box;
  }
</style>
