<script lang="ts">
    import LoginForm from '$lib/components/login-form.svelte';
    import logo from '$lib/images/logo.svg';
    import background from '$lib/images/pawel-czerwinski-NTYYL9Eb9y8-unsplash.jpg';
    import {toOptimizedURL} from "sveltekit-image-optimize/components";
</script>

<div class="grid min-h-svh lg:grid-cols-2">
    <div class="flex flex-col gap-4 p-6 md:p-10">
        <div class="flex justify-center gap-2 md:justify-start">
            <a class="flex items-center gap-2 font-medium" href="##">
                <div
                        class="text-primary-foreground flex size-16 items-center justify-center rounded-md"
                >
                    <img alt="MediaManager Logo" class="size-12" src={logo}/>
                </div>
                <h1 class="scale-110">Media Manager</h1>
            </a>
        </div>
        <div class="flex flex-1 items-center justify-center">
            <div class="w-full max-w-xs">
                <LoginForm/>
            </div>
        </div>
    </div>
    <div class="relative hidden lg:block">
        <img
                alt="background"
                class="absolute inset-0 h-full w-full object-cover dark:brightness-[0.8]  rounded-l-3xl "
                src="{toOptimizedURL(background)}"
        />
    </div>
</div>