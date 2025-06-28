<script lang="ts">
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';
	import { browser } from '$app/environment';
	import { redirect } from '@sveltejs/kit';

	if (browser) goto(base + '/dashboard');
	else throw redirect(307, base + '/login');
</script>

<svelte:head>
	<title>Home</title>
	<meta content="Svelte demo app" name="description" />
</svelte:head>
<h1>Servas</h1>
