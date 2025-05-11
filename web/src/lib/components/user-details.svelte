<script lang="ts">
	import {getContext} from 'svelte';
	import type {User} from '$lib/types';

	const user: () => User = getContext('user');
</script>

<span class="truncate font-semibold">{user().email}</span>
{#if user().is_superuser}
	<span class="truncate text-xs">Administrator</span>
{:else}
	<span class="truncate text-xs">User</span>
{/if}
