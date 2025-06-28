<script lang="ts">
	import {Button} from '$lib/components/ui/button/index.js';
	import {env} from '$env/dynamic/public';
	import {toast} from 'svelte-sonner';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import {Label} from '$lib/components/ui/label/index.js';
	import {Input} from '$lib/components/ui/input/index.js';

	const apiUrl = env.PUBLIC_API_URL;
	let newPassword: string = $state('');
	let newEmail: string = $state('');
	let dialogOpen = $state(false);

	async function saveUser() {
		try {
			const response = await fetch(`${apiUrl}/users/me`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json'
				},
				credentials: 'include',
				body: JSON.stringify({
					...(newPassword !== '' && {password: newPassword}),
					...(newEmail !== '' && {password: newEmail})
				})
			});

			if (response.ok) {
				toast.success(`Updated details successfully.`);
				dialogOpen = false;
			} else {
				const errorText = await response.text();
				console.error(`Failed to update user: ${response.statusText}`, errorText);
				toast.error(`Failed to update user: ${response.statusText}`);
			}
		} catch (error) {
			console.error('Error updating user:', error);
			toast.error(
					'Error updating user: ' + (error instanceof Error ? error.message : String(error))
			);
		} finally {
			newPassword = '';
			newEmail = '';
		}
	}
</script>

<Dialog.Root bind:open={dialogOpen}>
	<Dialog.Trigger>
		<Button class="w-full" onclick={() => (dialogOpen = true)} variant="outline">
			Edit my details
		</Button>
	</Dialog.Trigger>
	<Dialog.Content class="w-full max-w-[600px] rounded-lg bg-white p-6 shadow-lg">
		<Dialog.Header>
			<Dialog.Title class="mb-1 text-xl font-semibold">Edit User Details</Dialog.Title>
			<Dialog.Description class="mb-4 text-sm text-gray-500">
				Change your email or password. Leave fields empty to not change them.
			</Dialog.Description>
		</Dialog.Header>
		<div class="space-y-6">
			<!-- Email -->
			<div>
				<Label class="mb-1 block text-sm font-medium" for="email">Email</Label>
				<Input
						bind:value={newEmail}
						class="w-full"
						id="email"
						placeholder="Keep empty to not change the email"
						type="email"
				/>
			</div>
			<!-- Password -->
			<div>
				<Label class="mb-1 block text-sm font-medium" for="password">Password</Label>
				<Input
						bind:value={newPassword}
						class="w-full"
						id="password"
						placeholder="Keep empty to not change the password"
						type="password"
				/>
			</div>
		</div>
		<div class="mt-8 flex justify-end gap-2">
			<Button onclick={() => saveUser()} variant="destructive">Save</Button>
		</div>
	</Dialog.Content>
</Dialog.Root>
