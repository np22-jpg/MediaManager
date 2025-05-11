<script lang="ts">
    import {page} from '$app/state';
    import {Separator} from '$lib/components/ui/separator/index.js';
    import * as Sidebar from '$lib/components/ui/sidebar/index.js';
    import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
    import * as Table from '$lib/components/ui/table/index.js';
    import {getTorrentStatusString} from '$lib/utils'; // Corrected path
    import type {Torrent} from '$lib/types';

    let torrentsPromise: Promise<Torrent[]> = page.data.torrents;
</script>

<header class="flex h-16 shrink-0 items-center gap-2">
    <div class="flex items-center gap-2 px-4">
        <Sidebar.Trigger class="-ml-1"/>
        <Separator orientation="vertical" class="mr-2 h-4"/>
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

<div class="flex w-full flex-1 flex-col gap-4 p-4 pt-0">
    <div
            class="grid w-full auto-rows-min gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5"
    >
        {#await torrentsPromise}
            Loading...
        {:then torrents}
            <Table.Root>
                <Table.Caption>A list of the torrents.</Table.Caption>
                <Table.Header>
                    <Table.Row>
                        <Table.Head class="w-[100px]">Name</Table.Head>
                        <Table.Head>Download Status</Table.Head>
                        <Table.Head>Import Status</Table.Head>
                    </Table.Row>
                </Table.Header>
                <Table.Body>
                    {#each torrents as torrent}
                        <a href={'/dashboard/torrents/' + torrent.id}>
                            <Table.Row>
                                <Table.Cell class="font-medium">{torrent.title}</Table.Cell>
                                <Table.Cell>{getTorrentStatusString(torrent.status)}</Table.Cell>
                                <Table.Cell>{torrent.imported ? 'Yes' : 'No'}</Table.Cell>
                            </Table.Row>
                        </a>
                    {/each}
                </Table.Body>
            </Table.Root>
        {/await}
    </div>
</div>
