<script lang="ts">
    import {useSidebar} from "$lib/components/ui/sidebar/index.js";
    import UserDetails from "$lib/components/user-details.svelte"

    //let { user }: { user: { name: string; email: string; avatar: string } } = $props();
    const sidebar = useSidebar();
</script>

<Sidebar.Menu>
    <Sidebar.MenuItem>
        <DropdownMenu.Root>
            <DropdownMenu.Trigger>
                {#snippet child({props})}
                    <Sidebar.MenuButton
                            size="lg"
                            class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                            {...props}
                    >
                        <div class="grid flex-1 text-left text-sm leading-tight">
                            <UserDetails/>
                        </div>
                        <ChevronsUpDown class="ml-auto size-4"/>
                    </Sidebar.MenuButton>
                {/snippet}
            </DropdownMenu.Trigger>
            <DropdownMenu.Content
                    align="end"
                    class="w-[var(--bits-dropdown-menu-anchor-width)] min-w-56 rounded-lg"
                    side={sidebar.isMobile ? "bottom" : "right"}
                    sideOffset={4}
            >
                <DropdownMenu.Label class="p-0 font-normal">
                    <div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                        <div class="grid flex-1 text-left text-sm leading-tight">
                            <UserDetails/>
                        </div>
                    </div>
                </DropdownMenu.Label>
                <DropdownMenu.Separator/>
                <DropdownMenu.Group>
                    <DropdownMenu.Item>
                        <Sparkles/>
                        Upgrade to Pro
                    </DropdownMenu.Item>
                </DropdownMenu.Group>
                <DropdownMenu.Separator/>
                <DropdownMenu.Group>
                    <DropdownMenu.Item>
                        <BadgeCheck/>
                        Account
                    </DropdownMenu.Item>
                    <DropdownMenu.Item>
                        <CreditCard/>
                        Billing
                    </DropdownMenu.Item>
                    <DropdownMenu.Item>
                        <Bell/>
                        Notifications
                    </DropdownMenu.Item>
                </DropdownMenu.Group>
                <DropdownMenu.Separator/>
                <DropdownMenu.Item>
                    <LogOut/>
                    Log out
                </DropdownMenu.Item>
            </DropdownMenu.Content>
        </DropdownMenu.Root>
    </Sidebar.MenuItem>
</Sidebar.Menu>
