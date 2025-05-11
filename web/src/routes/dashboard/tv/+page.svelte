<script lang="ts">
    import {page} from '$app/state';
    import * as Card from '$lib/components/ui/card/index.js';
    import {env} from '$env/dynamic/public';
    import {Separator} from '$lib/components/ui/separator/index.js';
    import * as Sidebar from '$lib/components/ui/sidebar/index.js';
    import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
    import {toOptimizedURL} from 'sveltekit-image-optimize/components';

    let tvShowsPromise = page.data.tvShows;
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
                    <Breadcrumb.Page>Shows</Breadcrumb.Page>
                </Breadcrumb.Item>
            </Breadcrumb.List>
        </Breadcrumb.Root>
    </div>
</header>

<div class="flex w-full flex-1 flex-col gap-4 p-4 pt-0">
    <div
            class="grid w-full auto-rows-min gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5"
    >
        {#await tvShowsPromise}
            Loading...
        {:then tvShowsJson}
            {#await tvShowsJson.json()}
                Loading...
            {:then tvShows}
                {#each tvShows as show}
                    <a href={'/dashboard/tv/' + show.id}>
                        <Card.Root>
                            <Card.Header>
                                <Card.Title>{show.name}</Card.Title>
                                <Card.Description class="truncate">{show.overview}</Card.Description>
                            </Card.Header>
                            <Card.Content>
                                <img
                                        class="aspect-9/16 h-auto max-w-full rounded-lg object-cover"
                                        src={toOptimizedURL(`${env.PUBLIC_API_URL}/static/image/${show.id}.jpg`)}
                                        alt="{show.name}'s Poster Image"
                                />
                            </Card.Content>
                            <Card.Footer>
                                <p>Card Footer</p>
                            </Card.Footer>
                        </Card.Root>
                    </a>
                {/each}
            {/await}
        {/await}
    </div>
</div>
