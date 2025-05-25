<script lang="ts">
    import UserTable from '$lib/components/user-data-table.svelte';
    import {page} from '$app/state';
    import * as Card from "$lib/components/ui/card/index.js";
    import {getContext} from "svelte";
    import UserSettings from '$lib/components/user-settings.svelte';

    let currentUser = getContext("user")
    let users = page.data.users;
</script>

<div class="flex w-full flex-1 flex-col gap-4 p-4 pt-0 max-w-[1000px] mx-auto">
    <h1 class="scroll-m-20 my-6 text-center text-4xl font-extrabold tracking-tight lg:text-5xl">
        Settings
    </h1>
    <Card.Root id="me">
        <Card.Header>
            <Card.Title>You</Card.Title>
            <Card.Description>Change your email or password</Card.Description>
        </Card.Header>
        <Card.Content>
            <UserSettings/>
        </Card.Content>
    </Card.Root>
    {#if currentUser().is_superuser}
        <Card.Root id="users">
            <Card.Header>
                <Card.Title>Users</Card.Title>
                <Card.Description>Edit or delete users</Card.Description>
            </Card.Header>
            <Card.Content>
                <UserTable bind:users={users}/>
            </Card.Content>
        </Card.Root>
    {/if}
</div>

