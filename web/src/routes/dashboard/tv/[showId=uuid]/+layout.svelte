<script lang="ts">
	import { setContext } from 'svelte';
	import type { LayoutProps } from './$types';

	let { data, children }: LayoutProps = $props();

	const showData = $derived(data.showData);
	setContext('show', () => showData);
	const fetchError = $derived((data as { error?: string }).error || null);
</script>

{#if fetchError}
	<p>Error loading show: {fetchError}</p>
{:else if showData}
	{@render children()}
{:else}
	<p>Loading show data...</p>
{/if}
