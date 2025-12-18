<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    title: string;
    children?: Snippet;
  }

  let { title, children }: Props = $props();
</script>

<div class="content-header-wrapper">
  <div class="content-header">
    <h2>{title}</h2>
    {#if children}
      <div class="header-actions">
        {@render children()}
      </div>
    {/if}
  </div>
</div>

<style>
  .content-header-wrapper {
    /* Full screen width - breaking out of container */
    position: sticky;
    top: 0;
    z-index: 50;
    margin-left: calc(-50vw + 50%);
    margin-right: calc(-50vw + 50%);
    margin-bottom: 1.5rem;
  }

  .content-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    background: var(--color-background-blur);
    backdrop-filter: blur(15px);
    -webkit-backdrop-filter: blur(15px);
    padding: 1rem max(1rem, calc((100vw - 80rem) / 2 + 2rem));
    border-radius: 0.5rem;
  }

  .content-header h2 {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--color-text, #0a0a0a);
    margin: 0;
    font-family: var(--font-heading, "Playfair Display", serif);
    width: 12rem;
    flex-shrink: 0;
  }

  .header-actions {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
    align-items: flex-start;
    justify-content: flex-end;
  }
</style>
