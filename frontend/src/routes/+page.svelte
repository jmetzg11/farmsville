<script>
	import { onMount } from 'svelte';
	import { initializeUserStore } from '$lib/stores/auth';
	import { getComments } from '$lib/root/helpers';
	import Items from '$lib/root/Items.svelte';
	import ClaimedItmes from '$lib/root/ClaimedItmes.svelte';
	let items = [];
	let claimedItems = [];
	import { user } from '$lib/stores/auth';

	onMount(async () => {
		const results = await getComments();
		if (results) {
			items = results.items;
			claimedItems = results.claimedItems;
		}
		initializeUserStore();
	});
</script>

<h1>Welcome, {$user.name ? $user.name : $user.email}</h1>

<Items {items} />

<ClaimedItmes {claimedItems} />
