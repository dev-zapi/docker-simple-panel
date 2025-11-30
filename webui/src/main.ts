import { mount } from 'svelte'
import './app.css'
import App from './App.svelte'
// Import theme store to initialize theme on app load
import './stores/themeStore'

const app = mount(App, {
  target: document.getElementById('app')!,
})

export default app
