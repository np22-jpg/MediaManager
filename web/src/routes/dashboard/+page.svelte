<script lang="ts">
    import {Separator} from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import RecommendedShowsCarousel from '$lib/components/recommended-shows-carousel.svelte';
	import LoadingBar from '$lib/components/loading-bar.svelte';
    import {base} from '$app/paths';
    import {page} from '$app/state';
    import type {MetaDataProviderShowSearchResult} from '$lib/types';

	let recommendedShows: Promise<MetaDataProviderShowSearchResult[]> = page.data.tvRecommendations;
</script>

<header class="flex h-16 shrink-0 items-center gap-2">
	<div class="flex items-center gap-2 px-4">
        <Sidebar.Trigger class="-ml-1"/>
        <Separator class="mr-2 h-4" orientation="vertical"/>
		<Breadcrumb.Root>
			<Breadcrumb.List>
				<Breadcrumb.Item class="hidden md:block">
					<Breadcrumb.Link href="{base}/dashboard">MediaManager</Breadcrumb.Link>
				</Breadcrumb.Item>
                <Breadcrumb.Separator class="hidden md:block"/>
				<Breadcrumb.Item>
					<Breadcrumb.Page>Home</Breadcrumb.Page>
				</Breadcrumb.Item>
			</Breadcrumb.List>
		</Breadcrumb.Root>
	</div>
</header>
<div class="flex flex-1 flex-col gap-4 p-4 pt-0">
	<h1 class="scroll-m-20 text-center text-4xl font-extrabold tracking-tight lg:text-5xl">
		Dashboard
	</h1>
	<div class="min-h-[100vh] flex-1 items-center justify-center rounded-xl p-4 md:min-h-min">
        <div
                class="mx-auto max-w-[80vw] sm:max-w-[200px] md:max-w-[500px] lg:max-w-[750px] xl:max-w-[1200px]"
        >
			<h3 class="my-4 scroll-m-20 text-center text-2xl font-semibold tracking-tight">
				Trending Shows
			</h3>
			{#await recommendedShows}
                <LoadingBar/>
			{:then recommendations}
                <RecommendedShowsCarousel shows={recommendations}/>
			{/await}
		</div>
	</div>

	<!---
        <div class="grid auto-rows-min gap-4 md:grid-cols-3">
            <div class="aspect-video rounded-xl bg-muted/50"></div>
            <div class="aspect-video rounded-xl bg-muted/50"></div>
            <div class="aspect-video rounded-xl bg-muted/50">
            </div>
        </div>
    -->
</div>
