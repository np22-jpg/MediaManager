<script lang="ts">
    import Autoplay from 'embla-carousel-autoplay';
    import * as Carousel from '$lib/components/ui/carousel/index.js';
    import type {MetaDataProviderShowSearchResult} from '$lib/types';
    import AddMediaCard from '$lib/components/add-media-card.svelte';
    import {Skeleton} from "$lib/components/ui/skeleton";
    import {Button} from "$lib/components/ui/button";
    import {ChevronDown, ChevronRight} from "lucide-svelte";

    let {media, isShow, isLoading}: {
        media: MetaDataProviderShowSearchResult[],
        isShow: boolean,
        isLoading: boolean
    } = $props();
</script>

<div class="grid w-full gap-4 sm:grid-cols-1
     md:grid-cols-2 lg:grid-cols-3">
    {#if isLoading}
        <Skeleton class="w-full h-[70vh]"/>
        <Skeleton class="w-full h-[70vh]"/>
        <Skeleton class="w-full h-[70vh]"/>

    {:else }
        {#each media.slice(0, 3) as mediaItem}
            <AddMediaCard isShow={isShow} result={mediaItem}/>
        {/each}
    {/if}
    {#if isShow}
        <Button class="md:col-start-2" variant="secondary" href="/dashboard/tv/add-show">
            More recommendations
            <ChevronRight/>
        </Button>
    {:else }
        <Button class="md:col-start-2" variant="secondary" href="/dashboard/movies/add-movie">
            More recommendations
            <ChevronRight/>
        </Button>
    {/if}
</div>

