<script>
	import { onMount } from 'svelte';
	import { initializeUserStore } from '$lib/stores/auth';
	import Items from './root_components/Items.svelte';
	import ClaimedItmes from './root_components/ClaimedItmes.svelte';
	import { refreshItems } from '$lib/stores/items';

	onMount(() => {
		refreshItems();
		initializeUserStore();
	});
</script>

<Items />

<ClaimedItmes />
