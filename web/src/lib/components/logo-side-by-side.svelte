<script lang="ts">
    import logo from '$lib/images/logo.svg';

    let props = $props();
</script>

<div {...props} class="flex items-center">
    <img alt="Logo" class="h-12 w-12 mr-2" src={logo}/>
    <span class="text-3xl font-bold">Media Manager</span>
</div>