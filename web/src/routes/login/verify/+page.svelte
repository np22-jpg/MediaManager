<script lang="ts">
    import {UserCheck} from 'lucide-svelte';
    import {Button} from '$lib/components/ui/button/index.js';
    import Logo from '$lib/components/logo-side-by-side.svelte';

</script>


<div class="flex min-h-screen flex-col items-center justify-center bg-background px-4 py-12 sm:px-6 lg:px-8">
    <div class="absolute top-4 left-4">
        <Logo/>
    </div>
    <div class="mx-auto w-full max-w-md text-center">


    <div class="mb-6">
            <UserCheck class="mx-auto h-16 w-16 text-primary"/>
        </div>
        <h1 class="mt-4 text-3xl font-bold tracking-tight text-foreground sm:text-4xl">
            Account Pending Activation
        </h1>
        <p class="mt-4 text-lg text-muted-foreground">
            Your account has been successfully created, but activation by an
            administrator is required.
        </p>
        <div class="mt-8">
            <Button href="/dashboard">Go to Dashboard</Button>
        </div>
        <p class="mt-10 text-sm text-muted-foreground">
            The above button will only work once your account is verified.
        </p>
        <p class="mt-10 text-sm text-muted-foreground end">
            If you have any questions, please contact an administrator.</p>
    </div>
</div>