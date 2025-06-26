<script lang="ts">
    import Autoplay from 'embla-carousel-autoplay';
    import * as Carousel from '$lib/components/ui/carousel/index.js';
    import type {MetaDataProviderShowSearchResult} from '$lib/types';
    import AddMediaCard from '$lib/components/add-media-card.svelte';

    let {media, isShow}: { media: MetaDataProviderShowSearchResult[], isShow: boolean } = $props();
</script>

<div class="grid w-full gap-4 sm:grid-cols-1
     md:grid-cols-2 lg:grid-cols-3">
    {#each media.slice(0, 3) as mediaItem}
        <AddMediaCard isShow={isShow} result={mediaItem}/>
    {/each}
</div>

