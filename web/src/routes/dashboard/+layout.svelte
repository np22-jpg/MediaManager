<script lang="ts">
	import AppSidebar from '$lib/components/app-sidebar.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type {LayoutProps} from './$types';
	import {setContext} from 'svelte';

	let {data, children}: LayoutProps = $props();
	console.log('Received User Data: ', data.user);
	setContext('user', () => data.user);
</script>

<Sidebar.Provider>
	<AppSidebar/>
	<Sidebar.Inset>
		{@render children()}
	</Sidebar.Inset>
</Sidebar.Provider>
