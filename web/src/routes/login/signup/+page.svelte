<script lang="ts">
    import LoginForm from '$lib/components/login-form.svelte';
    import {getContext} from "svelte";

    let oauthProvider: () => {oauth_name: string} = getContext("oauthProvider");
</script>

<svelte:head>
    <title>Login - MediaManager</title>
    <meta
            content="Signup - MediaManager"
            name="description"
    />
</svelte:head>

<LoginForm login={false} oauthProvider={oauthProvider()}/>
