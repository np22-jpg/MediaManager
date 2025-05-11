<script lang="ts">
    import {cn} from '$lib/utils'; // Assuming you have the cn utility from shadcn-svelte

    let {
        label,
        variant = 'default',
        size = 'default',
        onClose = undefined,
        class: className = ''
    } = $props<{
        label: string;
        variant?: 'default' | 'secondary' | 'outline' | 'destructive';
        size?: 'default' | 'sm' | 'lg';
        onClose?: () => void;
        class?: string;
    }>();

    // Base styles for the chip
    const baseStyles =
        'inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50';

    // Variant styles
    const variantStyles = {
        default:
            'border bg-background text-foreground shadow-sm hover:bg-accent hover:text-accent-foreground',
        secondary: 'bg-secondary text-secondary-foreground shadow-sm hover:bg-secondary/80',
        outline:
            'border border-input bg-background shadow-sm hover:bg-accent hover:text-accent-foreground',
        destructive: 'bg-destructive text-destructive-foreground shadow-sm hover:bg-destructive/90'
    };

    // Size styles
    const sizeStyles = {
        default: 'h-9 px-3 py-0.5', // Adjusted height for New York style
        sm: 'h-7 px-2 py-0.5 text-xs', // Adjusted height for New York style
        lg: 'h-10 px-4 py-0.5' // Adjusted height for New York style
    };

    // Styles for the close button
    const closeButtonStyles =
        'ml-1 inline-flex h-4 w-4 shrink-0 items-center justify-center rounded-full';
</script>

<div class={cn(baseStyles, variantStyles[variant], sizeStyles[size], className)}>
    {label}
    {#if onClose}
        <button
                class={cn(closeButtonStyles, 'hover:bg-accent-foreground/20')}
                onclick={onClose}
                aria-label="remove tag"
        >
            <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    class="h-3 w-3"
            >
                <path d="M18 6L6 18"/>
                <path d="M6 6L18 18"/>
            </svg>
        </button>
    {/if}
</div>
