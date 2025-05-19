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
				url: '#',
				icon: TvIcon,
				isActive: true,
				items: [
					{
						title: 'Shows',
						url: '/dashboard/tv'
					},
					{
						title: 'Add a show',
						url: '/dashboard/tv/add-show'
					},
					{
                        title: 'Torrents',
                        url: '/dashboard/tv/torrents'
                    },
                    {
						title: 'Settings',
						url: '#'
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

	let {ref = $bindable(null), ...restProps}: ComponentProps<typeof Sidebar.Root> = $props();
</script>

<Sidebar.Root bind:ref variant="inset" {...restProps}>
	<Sidebar.Header>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton size="lg">
					{#snippet child({props})}
						<a href="##" {...props}>
							<div
									class="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground"
							>
								<Command class="size-4"/>
							</div>
							<div class="grid flex-1 text-left text-sm leading-tight">
								<span class="truncate font-semibold">Acme Inc</span>
								<span class="truncate text-xs">Enterprise</span>
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
