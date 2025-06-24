<script lang="ts">
    import Autoplay from 'embla-carousel-autoplay';
    import * as Carousel from '$lib/components/ui/carousel/index.js';
    import type {MetaDataProviderShowSearchResult} from '$lib/types';
    import AddShowCard from '$lib/components/add-show-card.svelte';

    let {shows}: { shows: MetaDataProviderShowSearchResult } = $props();
</script>

<Carousel.Root
        opts={{
		align: 'start',
		loop: true
	}}
        class="max-w-[80vw]"
        plugins={[
		Autoplay({
			delay: 5000,
			stopOnInteraction: false,
			stopOnMouseEnter: true,
			playOnInit: true,
			stopOnFocusIn: true
		})
	]}
>

    <Carousel.Content>
        {#each shows as show}
            <Carousel.Item class="md:basis-1/2 md:max-w-[40vw] xl:max-w-[20vw]">
                <AddShowCard result={show}/>
            </Carousel.Item>
        {/each}
    </Carousel.Content>
    <Carousel.Previous/>
    <Carousel.Next/>
</Carousel.Root>
