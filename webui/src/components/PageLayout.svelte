<script lang="ts">
  import type { Snippet } from 'svelte';
  import ContentHeader from './ContentHeader.svelte';

  interface Props {
    title: string;
    actions?: Snippet;
    children: Snippet;
  }

  let { title, actions, children }: Props = $props();
</script>

<div class="page-layout">
  <ContentHeader {title}>
    {#if actions}
      {@render actions()}
    {/if}
  </ContentHeader>
  
  <div class="page-content">
    {@render children()}
  </div>
</div>

<style>
  .page-layout {
    max-width: 80rem;
    margin: 0 auto;
    padding: 2rem;
  }

  .page-content {
    /* No additional styling needed - inherits width from page-layout */
  }

  @media (max-width: 640px) {
    .page-layout {
      padding: 1rem;
    }
  }
</style>
