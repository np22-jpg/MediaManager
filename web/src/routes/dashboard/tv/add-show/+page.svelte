<script lang="ts">
	import {env} from '$env/dynamic/public';
	import {Separator} from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import {Input} from '$lib/components/ui/input';
	import {Label} from '$lib/components/ui/label';
	import {Button} from '$lib/components/ui/button';
	import {ChevronDown} from 'lucide-svelte';
	import * as Collapsible from '$lib/components/ui/collapsible/index.js';
	import type {MetaDataProviderShowSearchResult} from '$lib/types.js';
	import * as RadioGroup from '$lib/components/ui/radio-group/index.js';
	import AddShowCard from '$lib/components/add-show-card.svelte';
	import {toast} from 'svelte-sonner';

	let searchTerm: string = $state('');
	let metadataProvider: string = $state('tmdb');
	let results: MetaDataProviderShowSearchResult[] | null = $state(null);

	async function search() {
		if (searchTerm.length > 0) {
			let url = new URL(env.PUBLIC_API_URL + '/tv/search');
			url.searchParams.append('query', searchTerm);
			url.searchParams.append('metadata_provider', metadataProvider);
			toast.info(`Searching for "${searchTerm}" using ${metadataProvider.toUpperCase()}...`);
			try {
				const response = await fetch(url, {
					method: 'GET',
					credentials: 'include'
				});
				if (!response.ok) {
					const errorText = await response.text();
					throw new Error(`Search failed: ${response.status} ${errorText || response.statusText}`);
				}
				results = await response.json();
				if (results && results.length > 0) {
					toast.success(`Found ${results.length} result(s) for "${searchTerm}".`);
				} else {
					toast.info(`No results found for "${searchTerm}".`);
				}
			} catch (error) {
				const errorMessage =
						error instanceof Error ? error.message : 'An unknown error occurred during search.';
				console.error('Search error:', error);
				toast.error(errorMessage);
				results = null; // Clear previous results on error
			}
		} else {
			toast.warning('Please enter a search term.');
			results = null;
		}
	}
</script>

<header class="flex h-16 shrink-0 items-center gap-2">
	<div class="flex items-center gap-2 px-4">
		<Sidebar.Trigger class="-ml-1"/>
		<Separator class="mr-2 h-4" orientation="vertical"/>
		<Breadcrumb.Root>
			<Breadcrumb.List>
				<Breadcrumb.Item class="hidden md:block">
					<Breadcrumb.Link href="/dashboard">MediaManager</Breadcrumb.Link>
				</Breadcrumb.Item>
				<Breadcrumb.Separator class="hidden md:block"/>
				<Breadcrumb.Item>
					<Breadcrumb.Link href="/dashboard">Home</Breadcrumb.Link>
				</Breadcrumb.Item>
				<Breadcrumb.Separator class="hidden md:block"/>
				<Breadcrumb.Item>
					<Breadcrumb.Link href="/dashboard/tv">Shows</Breadcrumb.Link>
				</Breadcrumb.Item>
				<Breadcrumb.Separator class="hidden md:block"/>
				<Breadcrumb.Item>
					<Breadcrumb.Page>Add a Show</Breadcrumb.Page>
				</Breadcrumb.Item>
			</Breadcrumb.List>
		</Breadcrumb.Root>
	</div>
</header>

<div class="flex w-full flex-1 flex-col items-center gap-4 p-4 pt-0">
	<div class="grid w-full max-w-sm items-center gap-12">
		<h1 class="scroll-m-20 text-center text-4xl font-extrabold tracking-tight lg:text-5xl">
			Add a show
		</h1>
		<section>
			<Label for="search-box">Show Name</Label>
			<Input bind:value={searchTerm} id="search-box" placeholder="Show Name" type="text"/>
			<p class="text-sm text-muted-foreground">Search for a Show to add.</p>
		</section>
		<section>
			<Collapsible.Root class="w-[350px] space-y-2">
				<Collapsible.Trigger>
					<div class="flex items-center justify-between space-x-4 px-4">
						<h4 class="text-sm font-semibold">Advanced Settings</h4>
						<Button class="w-9 p-0" size="sm" variant="ghost">
							<ChevronDown/>
							<span class="sr-only">Toggle</span>
						</Button>
					</div>
				</Collapsible.Trigger>
				<Collapsible.Content class="space-y-2">
					<Label for="metadata-provider-selector">Choose which Metadata Provider to query.</Label>
					<RadioGroup.Root bind:value={metadataProvider} id="metadata-provider-selector">
						<div class="flex items-center space-x-2">
							<RadioGroup.Item id="option-one" value="tmdb"/>
							<Label for="option-one">TMDB (Recommended)</Label>
						</div>
						<div class="flex items-center space-x-2">
							<RadioGroup.Item id="option-two" value="tvdb"/>
							<Label for="option-two">TVDB</Label>
						</div>
					</RadioGroup.Root>
				</Collapsible.Content>
			</Collapsible.Root>
		</section>
		<section>
			<Button onclick={search} type="submit">Search</Button>
		</section>
	</div>

	<Separator class="my-8"/>

	{#if results != null}
		{#if results.length === 0}
			<h3 class="mx-auto">No Shows found.</h3>
		{:else}
			<div
					class="grid w-full max-w-full auto-rows-min gap-4 sm:grid-cols-1
             md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5"
			>
				{#each results as result}
					<AddShowCard {result}/>
				{/each}
			</div>
		{/if}
	{/if}
</div>
