<script lang="ts">
	import logo from '$lib/images/logo.svg';

	let props = $props();
</script>

<div {...props} class="flex items-center">
	<img alt="Logo" class="mr-2 h-12 w-12" src={logo} />
	<span class="text-3xl font-bold dark:text-white">Media Manager</span>
</div>
