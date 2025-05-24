<script lang="ts">
    import {Button} from '$lib/components/ui/button/index.js';
    import * as Card from '$lib/components/ui/card/index.js';
    import {Input} from '$lib/components/ui/input/index.js';
    import {Label} from '$lib/components/ui/label/index.js';
    import {goto} from '$app/navigation';
    import {env} from '$env/dynamic/public';
    import * as Tabs from "$lib/components/ui/tabs/index.js";

    let apiUrl = env.PUBLIC_API_URL;

    let email = '';
    let password = '';
    let errorMessage = '';
    let isLoading = false;
    let tabValue = "login";

    async function handleLogin(event: Event) {
        event.preventDefault();

        isLoading = true;
        errorMessage = '';

        const formData = new URLSearchParams();
        formData.append('username', email);
        formData.append('password', password);
        try {
            const response = await fetch(apiUrl + '/auth/cookie/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                },
                body: formData.toString(),
                credentials: 'include'
            });

            if (response.ok) {
                console.log('Login successful!');
                console.log('Received User Data: ', response);
                goto('/dashboard');
                errorMessage = 'Login successful! Redirecting...';
            } else {
                let errorText = await response.text();
                try {
                    const errorData = JSON.parse(errorText);
                    errorMessage = errorData.message || 'Login failed. Please check your credentials.';
                } catch {
                    errorMessage = errorText || 'Login failed. Please check your credentials.';
                }
                console.error('Login failed:', response.status, errorText);
            }
        } catch (error) {
            console.error('Login request failed:', error);
            errorMessage = 'An error occurred during the login request.';
        } finally {
            isLoading = false;
        }
    }

    async function handleSignup(event: Event) {
        event.preventDefault();

        isLoading = true;
        errorMessage = '';

        try {
            const response = await fetch(apiUrl + '/auth/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    email: email,
                    password: password

                }),
                credentials: 'include'
            });

            if (response.ok) {
                console.log('Registration successful!');
                console.log('Received User Data: ', response);
                tabValue = "login"; // Switch to login tab after successful registration
                errorMessage = 'Registration successful! Redirecting...';
            } else {
                let errorText = await response.text();
                try {
                    const errorData = JSON.parse(errorText);
                    errorMessage = errorData.message || 'Registration failed. Please check your credentials.';
                } catch {
                    errorMessage = errorText || 'Registration failed. Please check your credentials.';
                }
                console.error('Registration failed:', response.status, errorText);
            }
        } catch (error) {
            console.error('Registration request failed:', error);
            errorMessage = 'An error occurred during the Registration request.';
        } finally {
            isLoading = false;
        }
    }
</script>
{#snippet tabSwitcher()}
    <!--    <Tabs.List>-->
    <!--        <Tabs.Trigger value="login">Login</Tabs.Trigger>-->
    <!--        <Tabs.Trigger value="register">Sign up</Tabs.Trigger>-->
    <!--    </Tabs.List>-->
{/snippet}
<Tabs.Root class="w-[400px]" value={tabValue}>
    <Tabs.Content value="login">
        <Card.Root class="mx-auto max-w-sm">
            <Card.Header>
                {@render tabSwitcher()}

                <Card.Title class="text-2xl">Login</Card.Title>
                <Card.Description>Enter your email below to login to your account</Card.Description>
            </Card.Header>
            <Card.Content>
                <form class="grid gap-4" onsubmit={handleLogin}>
                    <div class="grid gap-2">
                        <Label for="email">Email</Label>
                        <Input bind:value={email} id="email" placeholder="m@example.com" required type="email"/>
                    </div>
                    <div class="grid gap-2">
                        <div class="flex items-center">
                            <Label for="password">Password</Label>
                            <!-- TODO: add link to relevant documentation -->
                            <a class="ml-auto inline-block text-sm underline" href="##"> Forgot your password? </a>
                        </div>
                        <Input bind:value={password} id="password" required type="password"/>
                    </div>

                    {#if errorMessage}
                        <p class="text-sm text-red-500">{errorMessage}</p>
                    {/if}

                    <Button class="w-full" disabled={isLoading} type="submit">
                        {#if isLoading}
                            Logging in...
                        {:else}
                            Login
                        {/if}
                    </Button>
                </form>

                <Button class="mt-2 w-full" variant="outline">Login with Google</Button>

                <div class="mt-4 text-center text-sm">
                    Don't have an account?
                    <span class="underline" onclick={tabValue="register"}> Sign up </span>
                </div>
            </Card.Content>
        </Card.Root>
    </Tabs.Content>
    <Tabs.Content value="register">
        <Card.Root class="mx-auto max-w-sm">

            <Card.Header>
                {@render tabSwitcher()}
                <Card.Title class="text-2xl">Sign up</Card.Title>
                <Card.Description>Enter your email and password below to sign up.</Card.Description>
            </Card.Header>
            <Card.Content>

                <form class="grid gap-4" onsubmit={handleSignup}>
                    <div class="grid gap-2">
                        <Label for="email">Email</Label>
                        <Input bind:value={email} id="email" placeholder="m@example.com" required type="email"/>
                    </div>
                    <div class="grid gap-2">
                        <div class="flex items-center">
                            <Label for="password">Password</Label>
                        </div>
                        <Input bind:value={password} id="password" required type="password"/>
                    </div>

                    {#if errorMessage}
                        <p class="text-sm text-red-500">{errorMessage}</p>
                    {/if}

                    <Button class="w-full" disabled={isLoading} type="submit">
                        {#if isLoading}
                            Signing up...
                        {:else}
                            Sign up
                        {/if}
                    </Button>
                </form>

                <Button class="mt-2 w-full" variant="outline">Login with Google</Button>

                <div class="mt-4 text-center text-sm">
                    Already have an account?
                    <span class="underline" onclick={tabValue="login"}>Login </span>
                </div>
            </Card.Content>
        </Card.Root>
    </Tabs.Content>
</Tabs.Root>


