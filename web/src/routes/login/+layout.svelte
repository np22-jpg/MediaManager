<script lang="ts">
    import logo from "$lib/images/logo.svg";
    import { Separator } from "$lib/components/ui/separator/index.js";

    import background from '$lib/images/pawel-czerwinski-NTYYL9Eb9y8-unsplash.jpg?enhanced';
    import { PUBLIC_VERSION } from '$env/static/public';
    import { base } from '$app/paths';
    import { page } from '$app/state';
    import {setContext} from "svelte";

    setContext("oauthProvider", () => page.data.oauthProvider)
    let {data, children} = $props();
</script>

<div class="grid min-h-svh lg:grid-cols-2">
    <div class="flex flex-col gap-4 p-6 md:p-10">
        <div class="flex justify-center gap-2 md:justify-start">
            <a class="flex items-center gap-2 " href="{base}/">
                <div class="flex size-16 items-center justify-center rounded-md text-primary-foreground">
                    <img alt="MediaManager Logo" class="size-12" src={logo} />
                </div>
                <div>
                    <h1 class="text-2xl font-bold">Media Manager</h1>
                    <span class="truncate text-xs">{PUBLIC_VERSION}</span>
                </div>
            </a>
        </div>
        <div class="flex flex-1 items-center justify-center">
            <div class="w-full max-w-[90vw]">
                {@render children()}
            </div>
        </div>
        <div class="flex flex-col items-center justify-center gap-3 text-center">
            <a class="underline" href="https://maxdorninger.github.io/MediaManager/troubleshooting.html">
                Trouble logging in?
            </a>
            <div
                class="flex flex-wrap items-center justify-center gap-x-3 gap-y-1 text-sm text-muted-foreground"
            >
                <a class="underline" href="https://github.com/maxdorninger/MediaManager">GitHub</a>
                <Separator class="h-4" orientation="vertical" />
                <a class="underline" href="https://github.com/sponsors/maxdorninger">Donate</a>
                <Separator class="h-4" orientation="vertical" />
                <a
                    class="underline"
                    href="https://unsplash.com/photos/blue-white-and-red-abstract-painting-NTYYL9Eb9y8"
                >
                    Image Credit
                </a>
            </div>
        </div>
    </div>
    <div class="relative hidden lg:block">
        <enhanced:img
                src={background}
                alt="background"
                class="absolute inset-0 h-full w-full rounded-l-3xl object-cover dark:brightness-[0.8]"
        />
    </div>
</div>
