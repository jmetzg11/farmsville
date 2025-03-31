<script>
	import { onMount } from 'svelte';
	import { initializeUserStore } from '$lib/stores/auth';
	import Items from '$lib/root/Items.svelte';
	import ClaimedItmes from '$lib/root/ClaimedItmes.svelte';
	import { refreshItems } from '$lib/stores/items';

	onMount(() => {
		refreshItems();
		initializeUserStore();
	});
</script>

<Items />

<ClaimedItmes />
