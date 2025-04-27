<script lang="ts">
	import type {WithElementRef, WithoutChildren} from "bits-ui";
	import type {HTMLAttributes} from "svelte/elements";
	import {cn} from "$lib/utils.js";

	let {
        ref = $bindable(null),
        class: className,
        ...restProps
    }: WithoutChildren<WithElementRef<HTMLAttributes<HTMLSpanElement>>> = $props();
</script>

<span
        {...restProps}
        aria-hidden="true"
        bind:this={ref}
        class={cn("flex size-9 items-center justify-center", className)}
        role="presentation"
>
	<Ellipsis class="size-4"/>
	<span class="sr-only">More</span>
</span>
