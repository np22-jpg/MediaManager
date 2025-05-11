<script lang="ts">
    import {setContext} from 'svelte';
    import type {LayoutProps} from './$types';

    let {data}: LayoutProps = $props();

    const showData = $derived(data.showData);
    setContext('show', showData);
    const fetchError = $derived(data.error);
</script>

{#if fetchError}
    <p>Error loading show: {fetchError}</p>
{:else if showData}
    <slot/>
{:else}
    <p>Loading show data...</p>
{/if}
