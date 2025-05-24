<script>
    import {Button} from '$lib/components/ui/button/index.js';
    import {env} from '$env/dynamic/public';
    import * as Card from '$lib/components/ui/card/index.js';
    import {ImageOff} from 'lucide-svelte';
    import {goto} from '$app/navigation';
    import {base} from '$app/paths';

    let loading = $state(false);
    let errorMessage = $state(null);
    let {result} = $props();

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
        <Card.Title>
            {result.name}
            {#if result.year != null}
                ({result.year})
            {/if}
        </Card.Title>
        <Card.Description class="truncate">{result?.overview}</Card.Description>
    </Card.Header>
    <Card.Content>
        {#if result.poster_path != null}
            <img
                    class="h-auto max-w-full rounded-lg object-cover"
                    src={result.poster_path}
                    alt="{result.name}'s Poster Image"
            />
        {:else}
            <ImageOff/>
        {/if}
    </Card.Content>
    <Card.Footer>
        <Button disabled={result.added || loading} onclick={() => addShow(result)}>
            {#if loading}
                Loading...
            {:else}
                {result.added ? 'Show already exists' : 'Add Show'}
            {/if}
        </Button>
        {#if errorMessage}
            <p class="text-sm text-red-500">{errorMessage}</p>
        {/if}
    </Card.Footer>
</Card.Root>
