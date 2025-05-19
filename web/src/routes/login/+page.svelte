<script lang="ts">
    import LoginForm from '$lib/components/login-form.svelte';
    import Logo from '$lib/components/logo-side-by-side.svelte';
</script>

<div class="flex h-screen w-full items-center justify-center px-4 pt-16">
    <div class="flex flex-col items-center">
        <!-- Ensure Logo doesn't inherit conflicting styles -->
        <div class="mb-16 scale-150">
            <Logo/>
        </div>

        <LoginForm/>
    </div>
</div>