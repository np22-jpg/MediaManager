<script lang="ts">
    import type {User} from '$lib/types.js';
    import CheckmarkX from '$lib/components/checkmark-x.svelte';
    import * as Table from '$lib/components/ui/table/index.js';
    import {Button} from "$lib/components/ui/button/index.js";
    import {env} from "$env/dynamic/public";
    import {toast} from "svelte-sonner";
    import * as Dialog from "$lib/components/ui/dialog/index.js";
    import {getTorrentQualityString} from "$lib/utils";
    import {Label} from "$lib/components/ui/label/index.js";
    import * as RadioGroup from "$lib/components/ui/radio-group/index.js";
    import {Input} from "$lib/components/ui/input/index.js";
    import {invalidateAll} from "$app/navigation";
    import {getContext} from 'svelte';

    let newPassword: string = $state('');
    let newEmail: string = $state('');
    let dialogOpen = $state(false);

    async function saveUser() {
        try {
            const response = await fetch(`${env.PUBLIC_API_URL}/users/me`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json'
                },
                credentials: 'include',
                body: JSON.stringify({
                    ...newPassword !== "" && {password: newPassword},
                    ...newEmail !== "" && {password: newEmail}

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
            toast.error('Error updating user: ' + (error instanceof Error ? error.message : String(error)));
        } finally {
            newPassword = '';
            newEmail = '';
        }
    }
</script>

<Dialog.Root bind:open={dialogOpen}>
    <Dialog.Trigger>
        <Button class="w-full" onclick={() => dialogOpen = true} variant="outline">
            Edit my details
        </Button>
    </Dialog.Trigger>
    <Dialog.Content class="max-w-[600px] w-full rounded-lg shadow-lg bg-white p-6">
        <Dialog.Header>
            <Dialog.Title class="text-xl font-semibold mb-1">
                Edit User Details
            </Dialog.Title>
            <Dialog.Description class="text-sm text-gray-500 mb-4">
                Change your email or password. Leave fields empty to not change them.
            </Dialog.Description>
        </Dialog.Header>
        <div class="space-y-6">
            <!-- Email -->
            <div>
                <Label class="block text-sm font-medium mb-1" for="email">Email</Label>
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
                <Label class="block text-sm font-medium mb-1" for="password">Password</Label>
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