<script lang="ts">
	import { goto } from '$app/navigation';
	import { PUBLIC_BASE_PATH } from '$env/dynamic/public';
	import { browser } from '$app/environment';
	import { redirect } from '@sveltejs/kit';

	if (browser) goto(PUBLIC_BASE_PATH + '/dashboard');
	else throw redirect(307, PUBLIC_BASE_PATH + '/login');
</script>

<svelte:head>
	<title>Home</title>
	<meta content="Svelte demo app" name="description" />
</svelte:head>
<h1>Servas</h1>
