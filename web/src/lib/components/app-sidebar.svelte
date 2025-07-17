<script lang="ts" module>
	import {
		Bell,
		CircleDollarSign,
		Clapperboard,
		Home,
		Info,
		LifeBuoy,
		Settings,
		TvIcon
	} from 'lucide-svelte';
	import { base } from '$app/paths';

	import { PUBLIC_VERSION } from '$env/static/public';

	const data = {
		navMain: [
			{
				title: 'Dashboard',
				url: base + '/dashboard',
				icon: Home,
				isActive: true
			},
			{
				title: 'TV',
				url: base + '/dashboard/tv',
				icon: TvIcon,
				isActive: true,
				items: [
					{
						title: 'Add a show',
						url: base + '/dashboard/tv/add-show'
					},
					{
						title: 'Torrents',
						url: base + '/dashboard/tv/torrents'
					},
					{
						title: 'Requests',
						url: base + '/dashboard/tv/requests'
					}
				]
			},
			{
				title: 'Movies',
				url: base + '/dashboard/movies',
				icon: Clapperboard,
				isActive: true,
				items: [
					{
						title: 'Add a movie',
						url: base + '/dashboard/movies/add-movie'
					},
					{
						title: 'Torrents',
						url: base + '/dashboard/movies/torrents'
					},
					{
						title: 'Requests',
						url: base + '/dashboard/movies/requests'
					}
				]
			}
		],
		navSecondary: [
			{
				title: 'Notifications',
				url: base + '/dashboard/notifications',
				icon: Bell
			},
			{
				title: 'Settings',
				url: base + '/dashboard/settings',
				icon: Settings
			},
			{
				title: 'Support',
				url: 'https://github.com/maxdorninger/MediaManager/issues',
				icon: LifeBuoy
			},
			{
				title: 'Donate',
				url: 'https://github.com/sponsors/maxdorninger',
				icon: CircleDollarSign
			},
			{
				title: 'About',
				url: base + '/dashboard/about',
				icon: Info
			}
		]
	};
</script>

<script lang="ts">
	import NavMain from '$lib/components/nav-main.svelte';
	import NavSecondary from '$lib/components/nav-secondary.svelte';
	import NavUser from '$lib/components/nav-user.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { ComponentProps } from 'svelte';
	import logo from '$lib/images/logo.svg';

	let { ref = $bindable(null), ...restProps }: ComponentProps<typeof Sidebar.Root> = $props();
</script>

<Sidebar.Root {...restProps} bind:ref variant="inset">
	<Sidebar.Header>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton size="lg">
					{#snippet child({ props })}
						<a href="{base}/dashboard" {...props}>
							<img class="size-12" src={logo} alt="Media Manager Logo" />
							<div class="grid flex-1 text-left text-sm leading-tight">
								<span class="truncate font-semibold">Media Manager</span>
								<span class="truncate text-xs">{PUBLIC_VERSION}</span>
							</div>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Header>
	<Sidebar.Content>
		<NavMain items={data.navMain} />
		<!--  <NavProjects projects={data.projects}/> -->
		<NavSecondary class="mt-auto" items={data.navSecondary} />
	</Sidebar.Content>
	<Sidebar.Footer>
		<NavUser />
	</Sidebar.Footer>
</Sidebar.Root>
