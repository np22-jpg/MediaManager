<script lang="ts" module>
    import BookOpen from "@lucide/svelte/icons/book-open";
    import Bot from "@lucide/svelte/icons/bot";
    import Settings2 from "@lucide/svelte/icons/settings-2";
    import SquareTerminal from "@lucide/svelte/icons/square-terminal";

    const data = {
        navMain: [
            {
                title: "Playground",
                url: "#",
                icon: SquareTerminal,
                isActive: true,
                items: [
                    {
                        title: "History",
                        url: "#",
                    },
                    {
                        title: "Starred",
                        url: "#",
                    },
                    {
                        title: "Settings",
                        url: "#",
                    },
                ],
            },
            {
                title: "Models",
                url: "#",
                icon: Bot,
                items: [
                    {
                        title: "Genesis",
                        url: "#",
                    },
                    {
                        title: "Explorer",
                        url: "#",
                    },
                    {
                        title: "Quantum",
                        url: "#",
                    },
                ],
            },
            {
                title: "Documentation",
                url: "#",
                icon: BookOpen,
                items: [
                    {
                        title: "Introduction",
                        url: "#",
                    },
                    {
                        title: "Get Started",
                        url: "#",
                    },
                    {
                        title: "Tutorials",
                        url: "#",
                    },
                    {
                        title: "Changelog",
                        url: "#",
                    },
                ],
            },
            {
                title: "Settings",
                url: "#",
                icon: Settings2,
                items: [
                    {
                        title: "General",
                        url: "#",
                    },
                    {
                        title: "Team",
                        url: "#",
                    },
                    {
                        title: "Billing",
                        url: "#",
                    },
                    {
                        title: "Limits",
                        url: "#",
                    },
                ],
            },
        ],
    };
</script>

<script lang="ts">
    import NavMain from "$lib/components/nav-main.svelte";
    import NavUser from "$lib/components/nav-user.svelte";
    import * as Sidebar from "$lib/components/ui/sidebar/index.js";
    import type {ComponentProps} from "svelte";

    let {
        ref = $bindable(null),
        collapsible = "icon",
        ...restProps
    }: ComponentProps<typeof Sidebar.Root> = $props();
</script>

<Sidebar.Root {...restProps} bind:ref {collapsible}>
    <Sidebar.Header>
        <h1 class="font-bold">MediaManager</h1>
    </Sidebar.Header>
    <Sidebar.Content>
        <NavMain items={data.navMain}/>
    </Sidebar.Content>
    <Sidebar.Footer>
        <NavUser user={data.user}/>
    </Sidebar.Footer>
    <Sidebar.Rail/>
</Sidebar.Root>
