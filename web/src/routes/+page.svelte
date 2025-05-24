<script lang="ts">
	import {goto} from '$app/navigation';
	import {base} from '$app/paths';
	import {onMount} from 'svelte';

	onMount(() => {
		goto(base + '/dashboard');
	});
</script>

<svelte:head>
	<title>Home</title>
	<meta content="Svelte demo app" name="description"/>
</svelte:head>
<h1>Servas</h1>
