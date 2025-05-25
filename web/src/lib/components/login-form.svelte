<script lang="ts">
	import {Button} from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import {Input} from '$lib/components/ui/input/index.js';
	import {Label} from '$lib/components/ui/label/index.js';
	import {goto} from '$app/navigation';
	import {env} from '$env/dynamic/public';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import {toast} from 'svelte-sonner';
	import LoadingBar from '$lib/components/loading-bar.svelte';

	let apiUrl = env.PUBLIC_API_URL;

	let {oauthProvider} = $props();
	let oauthProviderName = $derived(oauthProvider.oauth_name);

	let email = $state('');
	let password = $state('');
	let errorMessage = $state('');
	let isLoading = $state(false);
	let tabValue = $state('login');

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
				errorMessage = 'Login successful! Redirecting...';
				toast.success(errorMessage);
				goto('/dashboard');
			} else {
				let errorText = await response.text();
				try {
					const errorData = JSON.parse(errorText);
					errorMessage = errorData.message || 'Login failed. Please check your credentials.';
				} catch {
					errorMessage = errorText || 'Login failed. Please check your credentials.';
				}
				toast.error(errorMessage);
				console.error('Login failed:', response.status, errorText);
			}
		} catch (error) {
			console.error('Login request failed:', error);
			errorMessage = 'An error occurred during the login request.';
			toast.error(errorMessage);
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
				tabValue = 'login'; // Switch to login tab after successful registration
				errorMessage = 'Registration successful! Please login.';
				toast.success(errorMessage);
			} else {
				let errorText = await response.text();
				try {
					const errorData = JSON.parse(errorText);
					errorMessage = errorData.message || 'Registration failed. Please check your credentials.';
				} catch {
					errorMessage = errorText || 'Registration failed. Please check your credentials.';
				}
				toast.error(errorMessage);
				console.error('Registration failed:', response.status, errorText);
			}
		} catch (error) {
			console.error('Registration request failed:', error);
			errorMessage = 'An error occurred during the Registration request.';
			toast.error(errorMessage);
		} finally {
			isLoading = false;
		}
	}

	async function handleOauth() {
		try {
			const response = await fetch(
					apiUrl + '/auth/cookie/' + oauthProviderName + '/authorize?scopes=email',
					{
						method: 'GET',
						headers: {
							'Content-Type': 'application/json'
						}
					}
			);
			if (response.ok) {
				let result = await response.json();
				console.log(
						'Redirecting to OAuth provider:',
						oauthProviderName,
						'url: ',
						result.authorization_url
				);
				toast.success('Redirecting to ' + oauthProviderName + ' for authentication...');
				window.location = result.authorization_url;
			} else {
				let errorText = await response.text();
				toast.error(errorMessage);
				console.error('Login failed:', response.status, errorText);
			}
		} catch (error) {
			console.error('Login request failed:', error);
			errorMessage = 'An error occurred during the login request.';
			toast.error(errorMessage);
		}
	}
</script>

{#snippet oauthLogin()}
	{#await oauthProvider}
		<LoadingBar/>
	{:then result}
		{#if result.oauth_name != null}
			<Button class="mt-2 w-full" onclick={() => handleOauth()} variant="outline"
			>Login with {result.oauth_name}</Button
			>
		{/if}
	{/await}
{/snippet}
<Tabs.Root class="w-[400px]" value={tabValue}>
	<Tabs.Content value="login">
		<Card.Root class="mx-auto max-w-sm">
			<Card.Header>
				<Card.Title class="text-2xl">Login</Card.Title>
				<Card.Description>Enter your email below to login to your account</Card.Description>
			</Card.Header>
			<Card.Content>
				<form class="grid gap-4" onsubmit={handleLogin}>
					<div class="grid gap-2">
						<Label for="email">Email</Label>
						<Input
								bind:value={email}
								id="email"
								placeholder="m@example.com"
								required
								type="email"
						/>
					</div>
					<div class="grid gap-2">
						<div class="flex items-center">
							<Label for="password">Password</Label>
							<!-- TODO: add link to relevant documentation -->
							<a class="ml-auto inline-block text-sm underline" href="##">
								Forgot your password?
							</a>
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

				{@render oauthLogin()}

				<div class="mt-4 text-center text-sm">
					<Button onclick={() => (tabValue = 'register')} variant="link">
						Don't have an account? Sign up
					</Button>
				</div>
			</Card.Content>
		</Card.Root>
	</Tabs.Content>
	<Tabs.Content value="register">
		<Card.Root class="mx-auto max-w-sm">
			<Card.Header>
				<Card.Title class="text-2xl">Sign up</Card.Title>
				<Card.Description>Enter your email and password below to sign up.</Card.Description>
			</Card.Header>
			<Card.Content>
				<form class="grid gap-4" onsubmit={handleSignup}>
					<div class="grid gap-2">
						<Label for="email2">Email</Label>
						<Input
								bind:value={email}
								id="email2"
								placeholder="m@example.com"
								required
								type="email"
						/>
					</div>
					<div class="grid gap-2">
						<div class="flex items-center">
							<Label for="password2">Password</Label>
						</div>
						<Input bind:value={password} id="password2" required type="password"/>
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
				<!-- TODO: dynamically display oauth providers based on config -->
				{@render oauthLogin()}

				<div class="mt-4 text-center text-sm">
					<Button onclick={() => (tabValue = 'login')} variant="link"
					>Already have an account? Login
					</Button>
				</div>
			</Card.Content>
		</Card.Root>
	</Tabs.Content>
</Tabs.Root>
