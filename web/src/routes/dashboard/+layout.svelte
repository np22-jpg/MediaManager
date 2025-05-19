<script lang="ts">
	import AppSidebar from '$lib/components/app-sidebar.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type {LayoutProps} from './$types';
	import {setContext} from 'svelte';
	import {goto} from '$app/navigation';
	import {base} from "$app/paths";

	let {data, children}: LayoutProps = $props();
	console.log('Received User Data: ', data.user);
	if (!data.user.is_verified) {
		goto(base + '/login/verify')
	}
	setContext('user', () => data.user);
</script>

<Sidebar.Provider>
	<AppSidebar/>
	<Sidebar.Inset>
		{@render children()}
	</Sidebar.Inset>
</Sidebar.Provider>
