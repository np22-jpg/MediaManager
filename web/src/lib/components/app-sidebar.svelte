<script lang="ts" module>
	import LifeBuoy from '@lucide/svelte/icons/life-buoy';
	import Send from '@lucide/svelte/icons/send';
	import TvIcon from '@lucide/svelte/icons/tv';
	import LayoutPanelLeft from '@lucide/svelte/icons/layout-panel-left';
	import DownloadIcon from '@lucide/svelte/icons/download';

	const data = {
		navMain: [
			{
				title: 'TV',
				url: '/dashboard/tv',
				icon: TvIcon,
				isActive: true,
				items: [
					{
						title: 'Add a show',
						url: '/dashboard/tv/add-show'
					},
					{
                        title: 'Torrents',
                        url: '/dashboard/tv/torrents'
                    },
                    {
						title: 'Requests',
						url: '/dashboard/tv/requests'
					}

                ]
			},
		],
		navSecondary: [
			{
				title: 'Support',
				url: '#',
				icon: LifeBuoy
			},
			{
				title: 'Feedback',
				url: '#',
				icon: Send
			}
		],
		projects: [
			{
				name: 'Dashboard',
				url: '/dashboard',
				icon: LayoutPanelLeft
			}
		]
	};
</script>

<script lang="ts">
	import NavMain from '$lib/components/nav-main.svelte';
	import NavProjects from '$lib/components/nav-projects.svelte';
	import NavSecondary from '$lib/components/nav-secondary.svelte';
	import NavUser from '$lib/components/nav-user.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import Command from '@lucide/svelte/icons/command';
	import type {ComponentProps} from 'svelte';
	import logo from '$lib/images/logo.svg';
	import {base} from "$app/paths";

	let {ref = $bindable(null), ...restProps}: ComponentProps<typeof Sidebar.Root> = $props();
</script>

<Sidebar.Root bind:ref variant="inset" {...restProps}>
	<Sidebar.Header>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton size="lg">
					{#snippet child({props})}
						<a href="{base}/dashboard" {...props}>
							<img class="size-12" src={logo} alt="Media Manager Logo"/>
							<div class="grid flex-1 text-left text-sm leading-tight">
								<span class="truncate font-semibold">Media Manager</span>
								<span class="truncate text-xs">version? or smth else?</span>
							</div>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Header>
	<Sidebar.Content>
		<NavMain items={data.navMain}/>
		<NavProjects projects={data.projects}/>
		<NavSecondary items={data.navSecondary} class="mt-auto"/>
	</Sidebar.Content>
	<Sidebar.Footer>
		<NavUser/>
	</Sidebar.Footer>
</Sidebar.Root>
