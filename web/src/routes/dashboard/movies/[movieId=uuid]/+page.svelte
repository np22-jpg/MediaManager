<script lang="ts">
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import { ImageOff } from 'lucide-svelte';
	import { getContext } from 'svelte';
	import type { PublicMovie, RichShowTorrent, User } from '$lib/types.js';
	import { getFullyQualifiedMediaName } from '$lib/utils';
	import { page } from '$app/state';
	import TorrentTable from '$lib/components/torrent-table.svelte';
	import MediaPicture from '$lib/components/media-picture.svelte';
	import DownloadMovieDialog from '$lib/components/download-movie-dialog.svelte';
	import RequestMovieDialog from '$lib/components/request-movie-dialog.svelte';

	let movie: PublicMovie = page.data.movie;
	let user: () => User = getContext('user');
	let torrents: RichShowTorrent = page.data.torrents;
</script>

<svelte:head>
	<title>{getFullyQualifiedMediaName(movie)} - MediaManager</title>
	<meta
		content="View details and manage downloads for {getFullyQualifiedMediaName(
			movie
		)} in MediaManager"
		name="description"
	/>
</svelte:head>

<header class="flex h-16 shrink-0 items-center gap-2">
	<div class="flex items-center gap-2 px-4">
		<Sidebar.Trigger class="-ml-1" />
		<Separator class="mr-2 h-4" orientation="vertical" />
		<Breadcrumb.Root>
			<Breadcrumb.List>
				<Breadcrumb.Item class="hidden md:block">
					<Breadcrumb.Link href="/dashboard">MediaManager</Breadcrumb.Link>
				</Breadcrumb.Item>
				<Breadcrumb.Separator class="hidden md:block" />
				<Breadcrumb.Item>
					<Breadcrumb.Link href="/dashboard">Home</Breadcrumb.Link>
				</Breadcrumb.Item>
				<Breadcrumb.Separator class="hidden md:block" />
				<Breadcrumb.Item>
					<Breadcrumb.Link href="/dashboard/movies">Movies</Breadcrumb.Link>
				</Breadcrumb.Item>
				<Breadcrumb.Separator class="hidden md:block" />
				<Breadcrumb.Item>
					<Breadcrumb.Page>{getFullyQualifiedMediaName(movie)}</Breadcrumb.Page>
				</Breadcrumb.Item>
			</Breadcrumb.List>
		</Breadcrumb.Root>
	</div>
</header>
<h1 class="scroll-m-20 text-center text-4xl font-extrabold tracking-tight lg:text-5xl">
	{getFullyQualifiedMediaName(movie)}
</h1>
<div class="flex w-full flex-1 flex-col gap-4 p-4">
	<div class="flex flex-col gap-4 md:flex-row md:items-stretch">
		<div class="w-full overflow-hidden rounded-xl bg-muted/50 md:w-1/3 md:max-w-sm">
			{#if movie.id}
				<MediaPicture media={movie} />
			{:else}
				<div
					class="aspect-9/16 flex h-auto w-full items-center justify-center rounded-lg bg-gray-200 text-gray-500"
				>
					<ImageOff size={48} />
				</div>
			{/if}
		</div>
		<div class="w-full flex-auto rounded-xl bg-muted/50 p-4 md:w-1/4">
			<p class="leading-7 [&:not(:first-child)]:mt-6">
				{movie.overview}
			</p>
		</div>
		<div class="w-full flex-auto rounded-xl bg-muted/50 p-4 md:w-1/3">
			{#if user().is_superuser}
				<DownloadMovieDialog {movie} />
				<div class="my-4"></div>
			{/if}
			<RequestMovieDialog {movie} />
		</div>
	</div>
	<!-- 	<div class="flex-1 rounded-xl bg-muted/50 p-4">
            <div class="w-full overflow-x-auto">

            </div>
        </div> -->
	<div class="flex-1 rounded-xl bg-muted/50 p-4">
		<div class="w-full overflow-x-auto">
			<TorrentTable isShow={false} torrents={torrents.torrents} />
		</div>
	</div>
</div>
