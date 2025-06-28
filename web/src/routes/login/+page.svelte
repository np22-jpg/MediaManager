<script lang="ts">
	import LoginForm from '$lib/components/login-form.svelte';
	import logo from '$lib/images/logo.svg';
	import background from '$lib/images/pawel-czerwinski-NTYYL9Eb9y8-unsplash.jpg?enhanced';
	import { page } from '$app/state';

	let oauthProvider = page.data.oauthProvider;
</script>

<div class="grid min-h-svh lg:grid-cols-2">
	<div class="flex flex-col gap-4 p-6 md:p-10">
		<div class="flex justify-center gap-2 md:justify-start">
			<a class="flex items-center gap-2 font-medium" href="##">
				<div class="flex size-16 items-center justify-center rounded-md text-primary-foreground">
					<img alt="MediaManager Logo" class="size-12" src={logo} />
				</div>
				<h1 class="scale-110">Media Manager</h1>
			</a>
		</div>
		<div class="flex flex-1 items-center justify-center">
			<div class="w-full max-w-[90vw]">
				<LoginForm {oauthProvider} />
			</div>
		</div>
	</div>
	<div class="relative hidden lg:block">
		<enhanced:img
			src={background}
			alt="background"
			class="absolute inset-0 h-full w-full rounded-l-3xl object-cover dark:brightness-[0.8]"
		/>
	</div>
</div>
