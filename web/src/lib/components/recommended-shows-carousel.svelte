<script lang="ts">
    import * as Card from "$lib/components/ui/card/index.js";
    import Autoplay from "embla-carousel-autoplay";
    import * as Carousel from "$lib/components/ui/carousel/index.js";
    import type {MetaDataProviderShowSearchResult} from "$lib/types";
    import AddShowCard from "$lib/components/add-show-card.svelte";

    let {shows}: { shows: MetaDataProviderShowSearchResult } = $props();
</script>
<Carousel.Root
        opts={{
  align: "start",
  loop: true,
 }}
        plugins={[
    Autoplay({
      delay: 2000,
      stopOnInteraction: false,
      stopOnMouseEnter: true,
      playOnInit: true,

    }),
  ]}
>
    <Carousel.Content class="-ml-1">
        {#each shows as show}
            <Carousel.Item class="pl-1 md:basis-1/2 lg:basis-1/3">
                <div class="p-1">
                    <AddShowCard result={show}/>
                </div>
            </Carousel.Item>
        {/each}
    </Carousel.Content>
    <Carousel.Previous/>
    <Carousel.Next/>
</Carousel.Root>
