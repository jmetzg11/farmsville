<script>
	import { onMount } from 'svelte';
	import { initializeUserStore } from '$lib/stores/auth';
	import Items from '$lib/root/Items.svelte';
	import ClaimedItmes from '$lib/root/ClaimedItmes.svelte';
	import { user } from '$lib/stores/auth';
	import { refreshItems } from '$lib/stores/items';

	onMount(() => {
		refreshItems();
		initializeUserStore();
	});
</script>

<h1>Welcome, {$user.name ? $user.name : $user.email ? $user.email : 'Guest'}</h1>

<Items />

<ClaimedItmes />
