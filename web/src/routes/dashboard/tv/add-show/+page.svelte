<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import {env} from '$env/dynamic/public';
	import {Separator} from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import {Input} from '$lib/components/ui/input';
	import {Label} from '$lib/components/ui/label';
	import {Button} from '$lib/components/ui/button';
	import {ChevronDown, ImageOff} from 'lucide-svelte';
	import * as Collapsible from '$lib/components/ui/collapsible/index.js';
	import type {MetaDataProviderShowSearchResult} from '$lib/types.js';
	import * as RadioGroup from '$lib/components/ui/radio-group/index.js';

	let searchTerm: string = $state('');
	let metadataProvider: string = $state('tmdb');
	let results:
			| (MetaDataProviderShowSearchResult & { added?: boolean; downloaded?: boolean })[]
			| null = $state(null);

	async function search() {
		if (searchTerm.length > 0) {
			let url = new URL(env.PUBLIC_API_URL + '/tv/search');
			url.searchParams.append('query', searchTerm);
			url.searchParams.append('metadata_provider', metadataProvider);
			const response = await fetch(url, {
				method: 'GET',
				credentials: 'include'
			});
			results = await response.json();
		} else {
			results = null;
		}
	}

	async function addShow(show: MetaDataProviderShowSearchResult & { added?: boolean }) {
		let url = new URL(env.PUBLIC_API_URL + '/tv/shows');
		url.searchParams.append('show_id', String(show.external_id));
		url.searchParams.append('metadata_provider', show.metadata_provider);
		const response = await fetch(url, {
			method: 'POST',
			credentials: 'include'
		});

		if (response.ok) {
			if (results) {
				const index = results.findIndex(
						(item) =>
								item.external_id === show.external_id &&
								item.metadata_provider === show.metadata_provider
				);
				if (index !== -1) {
					results[index].added = true;
					results = [...results];
				}
			}
		}
		return response;
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

<div class="flex w-full flex-1 flex-col items-center  gap-4 p-4 pt-0">
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
				{#each results as result (result.external_id)}
					<Card.Root class="h-full max-w-sm">
						<Card.Header>
							<Card.Title>
								{result.name}
								{#if result.year != null}
									({result.year})
								{/if}
							</Card.Title>
							<Card.Description class="truncate">{result?.overview}</Card.Description>
						</Card.Header>
						<Card.Content>
							{#if result.poster_path != null}
								<img
										class="h-auto max-w-full rounded-lg object-cover"
										src={result.poster_path}
										alt="{result.name}'s Poster Image"
								/>
							{:else}
								<ImageOff/>
							{/if}
						</Card.Content>
						<Card.Footer>
							<Button onclick={() => addShow(result)} disabled={result.added}>
								{result.added ? 'Show already exists' : 'Add Show'}
							</Button>
						</Card.Footer>
					</Card.Root>
				{/each}
			</div>
		{/if}
	{/if}
</div>
