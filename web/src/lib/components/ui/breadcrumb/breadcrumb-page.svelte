<script lang="ts">
	import type {WithElementRef} from "bits-ui";
	import type {HTMLAttributes} from "svelte/elements";
	import {cn} from "$lib/utils.js";

	let {
        ref = $bindable(null),
        class: className,
        children,
        ...restProps
    }: WithElementRef<HTMLAttributes<HTMLSpanElement>> = $props();
</script>

<span
        {...restProps}
        aria-current="page"
        aria-disabled="true"
        bind:this={ref}
        class={cn("text-foreground font-normal", className)}
        role="link"
>
	{@render children?.()}
</span>
