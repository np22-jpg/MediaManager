<script lang="ts">
	import {page} from '$app/state';
	import {Separator} from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import {getContext} from 'svelte';
	import type {PublicSeasonFile, Season, Show} from '$lib/types';
	import CheckmarkX from '$lib/components/checkmark-x.svelte';
	import {getFullyQualifiedShowName, getTorrentQualityString} from '$lib/utils';
	import {env} from "$env/dynamic/public";
	import {browser} from "$app/environment";
	import ShowPicture from "$lib/components/show-picture.svelte";

	const apiUrl = env.PUBLIC_API_URL
	const SeasonNumber = page.params.SeasonNumber;
	let seasonFiles: PublicSeasonFile[] = $state(page.data.files);
	let show: Show = getContext('show');
	let season: Season = $derived(
			show().seasons.find((item) => item.number === parseInt(SeasonNumber))
	);

	console.log('loaded files', seasonFiles);
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
					<Breadcrumb.Link href="/dashboard/tv/{show().id}">
						{show().name}
						{show().year == null ? '' : '(' + show().year + ')'}
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
	{getFullyQualifiedShowName(show())} Season {SeasonNumber}
</h1>
<div class="flex flex-1 flex-col gap-4 p-4">
	<div class="flex items-center gap-2">
		<div class="max-h-50% w-1/3 max-w-sm rounded-xl bg-muted/50">
			<ShowPicture show={show()}/>
		</div>
		<div class="h-full w-1/4 flex-auto rounded-xl bg-muted/50 p-4">
			<p class="leading-7 [&:not(:first-child)]:mt-6">
				{show().overview}
			</p>
		</div>
		<div class="h-full w-1/3 flex-auto rounded-xl bg-muted/50 p-4">
			<Table.Root>
				<Table.Caption>A list of all downloaded/downloading versions of this season.</Table.Caption>
				<Table.Header>
					<Table.Row>
						<Table.Head>Quality</Table.Head>
						<Table.Head>File Path Suffix</Table.Head>
						<Table.Head>Imported</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each seasonFiles as file}
						<Table.Row>
							<Table.Cell class="w-[50px]">
								{getTorrentQualityString(file.quality)}
							</Table.Cell>
							<Table.Cell class="w-[100px]">
								{file.file_path_suffix}
							</Table.Cell>
							<Table.Cell class="w-[10px] font-medium">
								<CheckmarkX state={file.downloaded}/>
							</Table.Cell>
						</Table.Row>
					{:else}
						<span class="font-semibold">You haven't downloaded this season yet.</span>
					{/each}
				</Table.Body>
			</Table.Root>
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
