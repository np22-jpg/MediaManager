<script lang="ts">
    import {getContext} from "svelte";

    const user: () => App.User = getContext("user")
</script>

<span class="truncate font-semibold">{user().email}</span>
{#if user().is_superuser}
    <span class="truncate text-xs">Administrator</span>
{:else }
    <span class="truncate text-xs">User</span>
{/if}