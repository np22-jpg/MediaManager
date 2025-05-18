<script lang="ts">
	import {page} from '$app/state';
	import {Separator} from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import {getTorrentQualityString, getTorrentStatusString} from '$lib/utils';
	import type {RichShowTorrent} from '$lib/types';
	import {getFullyQualifiedShowName} from '$lib/utils';
	import * as Accordion from '$lib/components/ui/accordion/index.js';
	import CheckmarkX from '$lib/components/checkmark-x.svelte';

	let showsPromise: Promise<RichShowTorrent[]> = $state(page.data.shows);
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
					<Breadcrumb.Page>Torrents</Breadcrumb.Page>
				</Breadcrumb.Item>
			</Breadcrumb.List>
		</Breadcrumb.Root>
	</div>
</header>

<div class="flex w-full flex-1 flex-col items-center gap-4 p-4 pt-0">
	{#await showsPromise}
		Loading...
	{:then shows}
		<Accordion.Root type="single" class="w-full lg:max-w-[70%]">
			{#each shows as show}
				<div class="w-full rounded-xl bg-muted/50 p-6">
					<Accordion.Item>
						<Accordion.Trigger>
							<h3 class="scroll-m-20 text-2xl font-semibold tracking-tight">
								{getFullyQualifiedShowName(show)}
							</h3>
						</Accordion.Trigger>
						<Accordion.Content>
							<Table.Root>
								<Table.Header>
									<Table.Row>
										<Table.Head class="w-[500px]">Name</Table.Head>
										<Table.Head>Seasons</Table.Head>
										<Table.Head>Download Status</Table.Head>
										<Table.Head>Quality</Table.Head>
										<Table.Head>File Path Suffix</Table.Head>
										<Table.Head>Imported</Table.Head>
									</Table.Row>
								</Table.Header>
								<Table.Body>
									{#each show.torrents as torrent}
										<Table.Row>
											<Table.Cell class="font-medium">
												<a href={'/dashboard/torrents/' + torrent.torrent_id}>
													{torrent.torrent_title}
												</a>
											</Table.Cell>
											<Table.Cell>
												<a href={'/dashboard/torrents/' + torrent.torrent_id}>
													{torrent.seasons}
												</a>
											</Table.Cell>
											<Table.Cell>
												<a href={'/dashboard/torrents/' + torrent.torrent_id}>
													{getTorrentStatusString(torrent.status)}
												</a>
											</Table.Cell>
											<Table.Cell class="font-medium">
												<a href={'/dashboard/torrents/' + torrent.torrent_id}>
													{getTorrentQualityString(torrent.quality)}
												</a>
											</Table.Cell>
											<Table.Cell class="font-medium">
												<a href={'/dashboard/torrents/' + torrent.torrent_id}>
													{torrent.file_path_suffix}
												</a>
											</Table.Cell>
											<Table.Cell>
												<a href={'/dashboard/torrents/' + torrent.torrent_id}>
													<CheckmarkX state={torrent.imported}/>
												</a>
											</Table.Cell>
										</Table.Row>
									{/each}
								</Table.Body>
							</Table.Root>
						</Accordion.Content>
					</Accordion.Item>
				</div>
			{:else}
				<h3 class="scroll-m-20 text-2xl font-semibold tracking-tight">
					You've not added any torrents yet.
				</h3>
			{/each}
		</Accordion.Root>
	{/await}
</div>
