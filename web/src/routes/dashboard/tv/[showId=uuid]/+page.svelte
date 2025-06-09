<script lang="ts">
    import {env} from '$env/dynamic/public';
    import {Separator} from '$lib/components/ui/separator/index.js';
    import * as Sidebar from '$lib/components/ui/sidebar/index.js';
    import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
    import {goto} from '$app/navigation';
    import {ImageOff} from 'lucide-svelte';
    import * as Table from '$lib/components/ui/table/index.js';
    import {getContext} from 'svelte';
    import type {RichShowTorrent, Show, User} from '$lib/types.js';
    import {getFullyQualifiedShowName} from '$lib/utils';
    import DownloadSeasonDialog from '$lib/components/download-season-dialog.svelte';
    import CheckmarkX from '$lib/components/checkmark-x.svelte';
    import {page} from '$app/state';
    import TorrentTable from '$lib/components/torrent-table.svelte';
    import RequestSeasonDialog from '$lib/components/request-season-dialog.svelte';
    import {browser} from "$app/environment";
    import ShowPicture from "$lib/components/show-picture.svelte";

    const apiUrl = env.PUBLIC_API_URL
    let show: Show = getContext('show');
    let user: User = getContext('user');
    let torrents: RichShowTorrent = page.data.torrentsData;
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
                    <Breadcrumb.Page>{getFullyQualifiedShowName(show())}</Breadcrumb.Page>
                </Breadcrumb.Item>
            </Breadcrumb.List>
        </Breadcrumb.Root>
    </div>
</header>
<h1 class="scroll-m-20 text-center text-4xl font-extrabold tracking-tight lg:text-5xl">
    {getFullyQualifiedShowName(show())}
</h1>
<div class="flex flex-1 w-full flex-col gap-4 p-4">
    <div class="flex flex-col md:flex-row md:items-stretch gap-4">
        <div class="w-full md:w-1/3 md:max-w-sm rounded-xl bg-muted/50 overflow-hidden">
            {#if show().id}
                <ShowPicture show={show()}/>
            {:else}
                <div
                        class="aspect-9/16 flex h-auto w-full items-center justify-center rounded-lg bg-gray-200 text-gray-500"
                >
                    <ImageOff size={48}/>
                </div>
            {/if}
        </div>
        <div class="w-full md:w-1/4 flex-auto rounded-xl bg-muted/50 p-4">
            <p class="leading-7 [&:not(:first-child)]:mt-6">
                {show().overview}
            </p>
        </div>
        <div
                class="w-full md:w-1/3 flex-auto rounded-xl bg-muted/50 p-4"
        >
            {#if user().is_superuser}
                <DownloadSeasonDialog show={show()}/>
            {/if}
            <RequestSeasonDialog show={show()}/>
        </div>
    </div>
    <div class="flex-1 rounded-xl bg-muted/50 p-4">
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
                    {#if show().seasons.length > 0}
                        {#each show().seasons as season (season.id)}
                            <Table.Row
                                    link={true}
                                    onclick={() => goto('/dashboard/tv/' + show().id + '/' + season.id)}
                            >
                                <Table.Cell class="min-w-[10px] font-medium">{season.number}</Table.Cell>
                                <Table.Cell class="min-w-[10px] font-medium">
                                    <CheckmarkX state={season.downloaded}/>
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
    <div class="flex-1 rounded-xl bg-muted/50 p-4">
        <div class="w-full overflow-x-auto">
            <TorrentTable torrents={torrents.torrents}/>
        </div>
    </div>
</div>
