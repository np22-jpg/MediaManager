<script lang="ts">
	import {env} from '$env/dynamic/public';
	import {Separator} from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import {Input} from '$lib/components/ui/input';
	import {Label} from '$lib/components/ui/label';
	import {Button} from '$lib/components/ui/button';
	import {ImageOff} from 'lucide-svelte';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import LoaderCircle from '@lucide/svelte/icons/loader-circle';
	import * as Select from '$lib/components/ui/select/index.js';

	import {buttonVariants} from '$lib/components/ui/button/index.js';

	import {getContext} from 'svelte';
	import {goto} from '$app/navigation';
	import type {PublicIndexerQueryResult, Show} from '$lib/types.js';
	import {getFullyQualifiedShowName} from '$lib/utils';
	import {toOptimizedURL} from 'sveltekit-image-optimize/components';

	import X from '@lucide/svelte/icons/x';
	import Check from '@lucide/svelte/icons/check';


	let show: Show = getContext('show');
	console.log('loaded show:', show);

	let selectedSeasonNumber: number = $state(1);
	let torrents: PublicIndexerQueryResult[] = $state([]);
	let isLoadingTorrents: boolean = $state(false);
	let torrentsError: string | null = $state(null);
	let queryOverride: string = $state('');
	let filePathSuffix: string = $state('');

	async function downloadTorrent(result_id: string) {
		let url = new URL(env.PUBLIC_API_URL + '/tv/torrents');
		url.searchParams.append('public_indexer_result_id', result_id);
		url.searchParams.append('show_id', show.id);
		if (filePathSuffix !== '') {
			url.searchParams.append('file_path_suffix', filePathSuffix);
		}
		try {
			const response = await fetch(url, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				credentials: 'include'
			});

			if (!response.ok) {
				const errorMessage = `Failed to download torrent for show ${show.id} and season ${selectedSeasonNumber}: ${response.statusText}`;
				console.error(errorMessage);
				torrentsError = errorMessage;
				return false;
			}

			const data: PublicIndexerQueryResult[] = await response.json();
			console.log('Downloading torrent:', data);

			return true;
		} catch (err) {
			const errorMessage = `Error downloading torrent: ${err instanceof Error ? err.message : 'An unknown error occurred'}`;
			console.error(errorMessage);
			return false;
		}
	}

	async function getTorrents(
			season_number: number,
			override: boolean = false
	): Promise<PublicIndexerQueryResult[]> {
		isLoadingTorrents = true;
		torrentsError = null;
		torrents = [];

		let url = new URL(env.PUBLIC_API_URL + '/tv/torrents');
		url.searchParams.append('show_id', show.id);
		if (override) {
			url.searchParams.append('search_query_override', queryOverride);
		} else {
			url.searchParams.append('season_number', season_number.toString());
		}

		try {
			const response = await fetch(url, {
				method: 'GET',
				headers: {
					'Content-Type': 'application/json'
				},
				credentials: 'include'
			});

			if (!response.ok) {
				const errorMessage = `Failed to fetch torrents for show ${show.id} and season ${selectedSeasonNumber}: ${response.statusText}`;
				console.error(errorMessage);
				torrentsError = errorMessage;
				return [];
			}

			const data: PublicIndexerQueryResult[] = await response.json();
			console.log('Fetched torrents:', data);

			return data;
		} catch (err) {
			const errorMessage = `Error fetching torrents: ${err instanceof Error ? err.message : 'An unknown error occurred'}`;
			console.error(errorMessage);
			torrentsError = errorMessage;
			return [];
		} finally {
			isLoadingTorrents = false;
		}
	}

	$effect(() => {
		if (show?.id) {
			console.log('selectedSeasonNumber changed:', selectedSeasonNumber);
			getTorrents(selectedSeasonNumber).then((fetchedTorrents) => {
				if (!isLoadingTorrents) {
					torrents = fetchedTorrents;
				} else if (fetchedTorrents.length > 0 || torrentsError) {
					torrents = fetchedTorrents;
				}
			});
		}
	});
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
					<Breadcrumb.Page>{getFullyQualifiedShowName(show)}</Breadcrumb.Page>
				</Breadcrumb.Item>
			</Breadcrumb.List>
		</Breadcrumb.Root>
	</div>
</header>
<h1 class="scroll-m-20 text-center text-4xl font-extrabold tracking-tight lg:text-5xl">
	{getFullyQualifiedShowName(show)}
</h1>
<div class="flex flex-1 flex-col gap-4 p-4">
	<div class="flex items-center gap-2">
		<div class="max-h-50% w-1/3 max-w-sm rounded-xl bg-muted/50">
			{#if show?.id}
				<img
						class="aspect-9/16 h-auto w-full rounded-lg object-cover"
						src={toOptimizedURL(`${env.PUBLIC_API_URL}/static/image/${show.id}.jpg`)}
						alt="{show.name}'s Poster Image"
				/>
			{:else}
				<div
						class="aspect-9/16 flex h-auto w-full items-center justify-center rounded-lg bg-gray-200 text-gray-500"
				>
					<ImageOff size={48}/>
				</div>
			{/if}
		</div>
		<div class="h-full flex-auto rounded-xl bg-muted/50 p-4">
			<p class="leading-7 [&:not(:first-child)]:mt-6">
				{show.overview}
			</p>
		</div>
		<div class="h-full flex-auto rounded-xl bg-muted/50 p-4">
			<Dialog.Root>
				<Dialog.Trigger class={buttonVariants({ variant: 'default' })}
				>Download Seasons
				</Dialog.Trigger>
				<Dialog.Content class="max-h-[90vh] w-fit min-w-[80vw] overflow-y-auto">
					<Dialog.Header>
						<Dialog.Title>Download a Season</Dialog.Title>
						<Dialog.Description>
							Search and download torrents for a specific season or season packs.
						</Dialog.Description>
					</Dialog.Header>
					<Tabs.Root class="w-full" value="basic">
						<Tabs.List>
							<Tabs.Trigger value="basic">Standard Mode</Tabs.Trigger>
							<Tabs.Trigger value="advanced">Advanced Mode</Tabs.Trigger>
						</Tabs.List>
						<Tabs.Content value="basic">
							<div class="grid w-full items-center gap-1.5">
								{#if show?.seasons?.length > 0}
									<Label for="season-number"
									>Enter a season number from 1 to {show.seasons.at(-1).number}</Label
									>
									<Input
											type="number"
											class="max-w-sm"
											id="season-number"
											bind:value={selectedSeasonNumber}
											max={show.seasons.at(-1).number}
									/>
									<p class="text-sm text-muted-foreground">
										Enter the season's number you want to search for. The first, usually 1, or the
										last season number usually yield the most season packs. Note that only Seasons
										which are listed in the "Seasons" cell will be imported!
									</p>
									<Label for="file-suffix">Filepath suffix</Label>
									<Select.Root type="single" bind:value={filePathSuffix} id="file-suffix">
										<Select.Trigger class="w-[180px]">{filePathSuffix}</Select.Trigger>
										<Select.Content>
											<Select.Item value="">None</Select.Item>
											<Select.Item value="2160P">2160p</Select.Item>
											<Select.Item value="1080P">1080p</Select.Item>
											<Select.Item value="720P">720p</Select.Item>
											<Select.Item value="480P">480p</Select.Item>
											<Select.Item value="360P">360p</Select.Item>
										</Select.Content>
									</Select.Root>
									<p class="text-sm text-muted-foreground">
										This is necessary to differentiate between versions of the same season/show, for
										example a 1080p and a 4K version of a season.
									</p>
									<Label for="file-suffix-display"
									>The files will be saved in the following directory:</Label
									>
									<p class="text-sm text-muted-foreground" id="file-suffix-display">
										/{getFullyQualifiedShowName(show)} [{show.metadata_provider}
										id-{show.external_id}
										]/Season XX/{show.name} SXXEXX {filePathSuffix === ''
											? ''
											: ' - ' + filePathSuffix}.mkv
									</p>
								{:else}
									<p class="text-sm text-muted-foreground">
										No season information available for this show.
									</p>
								{/if}
							</div>
						</Tabs.Content>
						<Tabs.Content value="advanced">
							<div class="grid w-full items-center gap-1.5">
								{#if show?.seasons?.length > 0}
									<Label for="query-override">Enter a custom query</Label>
									<div class="flex w-full max-w-sm items-center space-x-2">
										<Input type="text" id="query-override" bind:value={queryOverride}/>
										<Button
												variant="secondary"
												onclick={async () => {
												isLoadingTorrents = true;
												torrentsError = null;
												torrents = [];
												try {
													const fetchedTorrents = await getTorrents(selectedSeasonNumber, true);
													torrents = fetchedTorrents;
												} catch (error) {
													console.log(error);
												} finally {
													isLoadingTorrents = false;
												}
											}}
										>
											Search
										</Button>
									</div>
									<p class="text-sm text-muted-foreground">
										The custom query will override the default search string like "The Simpsons
										Season 3". Note that only Seasons which are listed in the "Seasons" cell will be
										imported!
									</p>
									<Label for="file-suffix">Filepath suffix</Label>
									<Input
											type="text"
											class="max-w-sm"
											id="file-suffix"
											bind:value={filePathSuffix}
											placeholder="1080P"
									/>
									<p class="text-sm text-muted-foreground">
										This is necessary to differentiate between versions of the same season/show, for
										example a 1080p and a 4K version of a season.
									</p>

									<Label for="file-suffix-display"
									>The files will be saved in the following directory:</Label
									>
									<p class="text-sm text-muted-foreground" id="file-suffix-display">
										/{getFullyQualifiedShowName(show)} [{show.metadata_provider}
										id-{show.external_id}
										]/Season XX/{show.name} SXXEXX {filePathSuffix === ''
											? ''
											: ' - ' + filePathSuffix}.mkv
									</p>
								{:else}
									<p class="text-sm text-muted-foreground">
										No season information available for this show.
									</p>
								{/if}
							</div>
						</Tabs.Content>
					</Tabs.Root>
					<div class="mt-4 items-center">
						{#if isLoadingTorrents}
							<div class="flex w-full max-w-sm items-center space-x-2">
								<LoaderCircle class="animate-spin"/>
								<p>Loading torrents...</p>
							</div>
						{:else if torrentsError}
							<p class="text-red-500">Error: {torrentsError}</p>
						{:else if torrents.length > 0}
							<h3 class="mb-2 text-lg font-semibold">Found Torrents:</h3>
							<div class="max-h-[200px] overflow-y-auto rounded-md border p-2">
								<Table.Root>
									<Table.Header>
										<Table.Row>
											<Table.Head>Title</Table.Head>
											<Table.Head>Size</Table.Head>
											<Table.Head>Seeders</Table.Head>
											<Table.Head>Indexer Flags</Table.Head>
											<Table.Head>Seasons</Table.Head>
											<Table.Head class="text-right">Actions</Table.Head>
										</Table.Row>
									</Table.Header>
									<Table.Body>
										{#each torrents as torrent (torrent.id)}
											<Table.Row>
												<Table.Cell
														class="max-w-[300px] font-medium">{torrent.title}</Table.Cell>
												<Table.Cell>{(torrent.size / 1024 / 1024 / 1024).toFixed(2)}GB
												</Table.Cell>
												<Table.Cell>{torrent.seeders}</Table.Cell>
												<Table.Cell>
													{#each torrent.flags as flag}
														{flag},&nbsp;
													{/each}
												</Table.Cell>
												<Table.Cell>
													{#if torrent.season.length === 1}
														{torrent.season[0]}
													{:else}
														{torrent.season.at(0)}−{torrent.season.at(-1)}
													{/if}
												</Table.Cell>
												<Table.Cell class="text-right">
													<Button
															size="sm"
															variant="outline"
															onclick={() => {
															downloadTorrent(torrent.id);
														}}
													>
														Download
													</Button>
												</Table.Cell>
											</Table.Row>
										{/each}
									</Table.Body>
								</Table.Root>
							</div>
						{:else if show?.seasons?.length > 0}
							<p>No torrents found for season {selectedSeasonNumber}. Try a different season.</p>
						{/if}
					</div>
				</Dialog.Content>
			</Dialog.Root>
		</div>
	</div>
	<div class="min-h-[100vh] flex-1 rounded-xl bg-muted/50 p-4 md:min-h-min">
		<div class="w-full overflow-x-auto">
			<Table.Root>
				<Table.Caption>A list of all seasons.</Table.Caption>
				<Table.Header>
					<Table.Row>
						<Table.Head>Number</Table.Head>
						<Table.Head>Exists on file</Table.Head>
						<Table.Head>Title</Table.Head>
						<Table.Head>Overview</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#if show?.seasons?.length > 0}
						{#each show.seasons as season (season.id)}
							<Table.Row
									link={true}
									onclick={() => goto('/dashboard/tv/' + show.id + '/' + season.number)}
							>
								<Table.Cell class="min-w-[10px] font-medium">{season.number}</Table.Cell>
								<Table.Cell class="min-w-[10px] font-medium">
									{#if season.downloaded}
										<Check class="stroke-green-500"/>
									{:else}
										<X class="stroke-rose-600"/>
									{/if}

								</Table.Cell>

								<Table.Cell class="min-w-[50px]">{season.name}</Table.Cell>
								<Table.Cell class="max-w-[300px] truncate">{season.overview}</Table.Cell>
							</Table.Row>
						{/each}
					{:else}
						<Table.Row>
							<Table.Cell colspan="3" class="text-center">No season data available.</Table.Cell>
						</Table.Row>
					{/if}
				</Table.Body>
			</Table.Root>
		</div>
	</div>
</div>
