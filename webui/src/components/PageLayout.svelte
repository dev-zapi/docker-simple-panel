<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    children: Snippet;
  }

  let { children }: Props = $props();
</script>

<div class="page-layout">
  {@render children()}
</div>

<style>
  .page-layout {
    max-width: 80rem;
    margin: 0 auto;
    padding: 2rem;
  }

  @media (max-width: 640px) {
    .page-layout {
      padding: 1rem;
    }
  }
</style>
