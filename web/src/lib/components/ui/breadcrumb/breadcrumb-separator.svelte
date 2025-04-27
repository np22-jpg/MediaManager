<script lang="ts">
    import type {WithElementRef} from "bits-ui";
    import type {HTMLLiAttributes} from "svelte/elements";
    import {cn} from "$lib/utils.js";

    let {
        ref = $bindable(null),
        class: className,
        children,
        ...restProps
    }: WithElementRef<HTMLLiAttributes> = $props();
</script>

<li
        {...restProps}
        aria-hidden="true"
        bind:this={ref}
        class={cn("[&>svg]:size-3.5", className)}
        role="presentation"
>
    {#if children}
        {@render children?.()}
    {:else}
        <ChevronRight/>
    {/if}
</li>
