<script>
	import { onMount } from 'svelte';
	import { initializeUserStore } from '$lib/stores/auth';
	import { getItems } from '$lib/root/helpers';
	import Items from '$lib/root/Items.svelte';
	import ClaimedItmes from '$lib/root/ClaimedItmes.svelte';
	let items = $state([]);
	let claimedItems = $state([]);
	import { user } from '$lib/stores/auth';

	onMount(async () => {
		const results = await getItems();
		if (results) {
			items = results.items;
			claimedItems = results.claimedItems;
		}
		initializeUserStore();
	});
</script>

<h1>Welcome, {$user.name ? $user.name : $user.email ? $user.email : 'Guest'}</h1>

<Items {items} />

<ClaimedItmes {claimedItems} />
