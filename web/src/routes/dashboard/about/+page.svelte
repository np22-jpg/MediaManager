<script lang="ts">
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import { base } from '$app/paths';
	import logo from '$lib/images/logo.svg';
	import { PUBLIC_VERSION } from '$env/static/public';
</script>

<header class="flex h-16 shrink-0 items-center gap-2">
	<div class="flex items-center gap-2 px-4">
		<Sidebar.Trigger class="-ml-1" />
		<Separator class="mr-2 h-4" orientation="vertical" />
		<Breadcrumb.Root>
			<Breadcrumb.List>
				<Breadcrumb.Item class="hidden md:block">
					<Breadcrumb.Link href="{base}/dashboard">MediaManager</Breadcrumb.Link>
				</Breadcrumb.Item>
				<Breadcrumb.Separator class="hidden md:block" />
				<Breadcrumb.Item>
					<Breadcrumb.Page>About</Breadcrumb.Page>
				</Breadcrumb.Item>
			</Breadcrumb.List>
		</Breadcrumb.Root>
	</div>
</header>

<div class="flex w-full flex-col items-center justify-center px-4 py-12">
	<img alt="Media Manager Logo" class="mb-4 h-24 w-24" src={logo} />
	<h1 class="mb-2 text-4xl font-bold">About Media Manager</h1>
	<p class="mb-6 max-w-2xl text-center text-lg">
		<strong>Media Manager</strong> is an all-in-one solution for organizing and building your media library.
		Built for simplicity and modernity, it helps you keep track of your favorite shows and movies and
		explore trending content—all in one place.
	</p>
	<p class="mb-2 text-sm text-muted-foreground">
		Version: v{PUBLIC_VERSION}
	</p>
	<div class="my-6 mb-6 flex items-center gap-2 text-sm text-muted-foreground sm:w-1/2 lg:w-1/3">
		<a class="flex items-center gap-2" href="https://www.themoviedb.org/" target="_blank">
			<img
				alt="TMDB Logo"
				class="h-auto w-20"
				src="https://www.themoviedb.org/assets/2/v4/logos/v2/blue_square_2-d537fb228cf3ded904ef09b136fe3fec72548ebc1fea3fbbd1ad9e36364db38b.svg"
			/>
			<span
				>Metadata provided by TMDB. Please consider adding missing information or subscribing.</span
			>
		</a>
	</div>
	<div class="my-6 mb-6 flex items-center gap-2 text-sm text-muted-foreground sm:w-1/2 lg:w-1/3">
		<a class="flex items-center gap-2" href="https://thetvdb.com/subscribe" target="_blank">
			<img
				alt="TheTVDB Logo"
				class="h-auto w-20"
				src="https://www.thetvdb.com/images/attribution/logo2.png"
			/>
			<span
				>Metadata provided by TheTVDB. Please consider adding missing information or subscribing.</span
			>
		</a>
	</div>
</div>
