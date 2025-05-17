<script lang="ts">
	import {page} from '$app/state';
	import {env} from '$env/dynamic/public';
	import {Separator} from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import {getContext} from 'svelte';
	import type {Season, Show} from '$lib/types';

	const SeasonNumber = page.params.SeasonNumber;
	let show: Show = getContext('show');
	let season: Season;
	show.seasons.forEach((item) => {
		if (item.number === parseInt(SeasonNumber)) season = item;
	});

	console.log('loaded ', show);
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
					<Breadcrumb.Link href="/dashboard/tv/{show.id}">
						{show.name}
						{show.year == null ? '' : '(' + show.year + ')'}
					</Breadcrumb.Link>
				</Breadcrumb.Item>
				<Breadcrumb.Separator class="hidden md:block"/>
				<Breadcrumb.Item>
					<Breadcrumb.Page>Season {SeasonNumber}</Breadcrumb.Page>
				</Breadcrumb.Item>
			</Breadcrumb.List>
		</Breadcrumb.Root>
	</div>
</header>
<h1 class="scroll-m-20 text-center text-4xl font-extrabold tracking-tight lg:text-5xl">
	{show.name}
	{show.year == null ? '' : '(' + show.year + ')'} Season {SeasonNumber}
</h1>
<div class="flex flex-1 flex-col gap-4 p-4">
	<div class="flex items-center gap-2">
		<div class="max-h-50% w-1/3 max-w-sm rounded-xl bg-muted/50">
			<img
					class="aspect-9/16 h-auto w-full rounded-lg object-cover"
					src="{env.PUBLIC_API_URL}/static/image/{show.id}.jpg"
					alt="{show.name}'s Poster Image"
			/>
		</div>
		<div class="h-full flex-auto rounded-xl bg-muted/50 p-4">
			<p class="leading-7 [&:not(:first-child)]:mt-6">
				{show.overview}
			</p>
		</div>
	</div>
	<div class="min-h-[100vh] flex-1 rounded-xl bg-muted/50 p-4 md:min-h-min">
		<div class="w-full overflow-x-auto">
			<Table.Root>
				<Table.Caption>A list of all episodes.</Table.Caption>
				<Table.Header>
					<Table.Row>
						<Table.Head class="w-[100px]">Number</Table.Head>
						<Table.Head class="min-w-[50px]">Title</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each season.episodes as episode (episode.id)}
						<Table.Row>
							<Table.Cell class="w-[100px] font-medium">{episode.number}</Table.Cell>
							<Table.Cell class="min-w-[50px]">{episode.title}</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</div>
	</div>
</div>
