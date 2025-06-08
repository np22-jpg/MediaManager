<script>
    import {getFullyQualifiedShowName} from "$lib/utils.js";
    import logo from "$lib/images/svelte-logo.svg";
    import {env} from "$env/dynamic/public";

    const apiUrl = env.PUBLIC_API_URL;
    let {show} = $props();
</script>

<picture>
    <source
            srcset="{apiUrl}/static/image/{show.id}.avif"
            type="image/avif"
    />
    <source
            srcset="{apiUrl}/static/image/{show.id}.webp"
            type="image/webp"
    />
    <img
            alt="{getFullyQualifiedShowName(show)}'s Poster Image"
            class="aspect-9/16 center h-auto w-full rounded-lg object-cover"
            src="{apiUrl}/static/image/{show.id}.jpeg"
    />
</picture>