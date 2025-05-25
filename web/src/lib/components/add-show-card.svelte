<script lang="ts">
    import {Button} from '$lib/components/ui/button/index.js';
    import {env} from '$env/dynamic/public';
    import * as Card from '$lib/components/ui/card/index.js';
    import {ImageOff} from 'lucide-svelte';
    import {goto} from '$app/navigation';
    import {base} from '$app/paths';
    import type {MetaDataProviderShowSearchResult} from '$lib/types.js';
    import {toOptimizedURL} from "sveltekit-image-optimize/components";

    let loading = $state(false);
    let errorMessage = $state(null);
    let {result}: { result: MetaDataProviderShowSearchResult } = $props();
    console.log('Add Show Card Result: ', result);

    async function addShow() {
        loading = true;
        let url = new URL(env.PUBLIC_API_URL + '/tv/shows');
        url.searchParams.append('show_id', String(result.external_id));
        url.searchParams.append('metadata_provider', result.metadata_provider);
        const response = await fetch(url, {
            method: 'POST',
            credentials: 'include'
        });
        let responseData = await response.json();
        console.log('Added Show: Response Data: ', responseData);
        if (response.ok) {
            await goto(base + '/dashboard/tv/' + responseData.id);
        } else {
            errorMessage = 'Error occurred: ' + responseData;
        }
        loading = false;
    }
</script>

<Card.Root class="h-full max-w-sm">
    <Card.Header>
        <Card.Title class="h-12 overflow-hidden leading-tight flex items-center">
            {result.name}
            {#if result.year != null}
                ({result.year})
            {/if}
        </Card.Title>
        <Card.Description
                class="truncate">{result.overview !== "" ? result.overview : "No overview available"}</Card.Description>
    </Card.Header>
    <Card.Content class="w-full h-96 flex items-center justify-center">
        {#if result.poster_path != null}
            <img
                    class="max-h-full max-w-full object-contain  rounded-lg"
                    src={toOptimizedURL(result.poster_path)}
                    alt="{result.name}'s Poster Image"
            />
        {:else}
            <div class="w-full h-full flex items-center justify-center">
                <ImageOff class="w-12 h-12 text-gray-400"/>
            </div>
        {/if}
    </Card.Content>
    <Card.Footer class="flex flex-col gap-2 items-start p-4 bg-card rounded-b-lg border-t">
        <Button
                class="w-full font-semibold"
                disabled={result.added || loading}
                onclick={() => addShow(result)}
        >
            {#if loading}
                <span class="animate-pulse">Loading...</span>
            {:else}
                {result.added ? 'Show already exists' : 'Add Show'}
            {/if}
        </Button>
        <div class="flex items-center gap-2 w-full">
            {#if result.vote_average != null}
            <span class="text-sm text-yellow-600 font-medium flex items-center">
                <svg class="w-4 h-4 mr-1 text-yellow-400" fill="currentColor" viewBox="0 0 20 20"><path
                        d="M10 15l-5.878 3.09 1.122-6.545L.488 6.91l6.561-.955L10 0l2.951 5.955 6.561.955-4.756 4.635 1.122 6.545z"/></svg>
                Rating: {result.vote_average}/10
            </span>
            {/if}
        </div>
        {#if errorMessage}
            <p class="text-xs text-red-500 bg-red-50 rounded px-2 py-1 w-full">{errorMessage}</p>
        {/if}
    </Card.Footer>
</Card.Root>
