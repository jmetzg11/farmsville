<script>
	import { onMount } from 'svelte';
	import Items from '$lib/root/Items.svelte';
	import ClaimedItmes from '$lib/root/ClaimedItmes.svelte';
	let items;
	let claimedItems;
	import { user } from '$lib/stores/auth';
	async function getComments() {
		try {
			const url = `${import.meta.env.VITE_API_URL}/items`;
			const response = await fetch(url);

			if (!response.ok) {
				throw new Error('Failed to get items');
			}

			const data = await response.json();
			items = data.items;
			claimedItems = data.claimedItems;
		} catch (error) {
			console.error('Error fetching items', error);
		}
	}
	onMount(() => {
		getComments();
	});
</script>

<h1>Welcome, {$user.name ? $user.name : $user.email}</h1>

<Items {items} />

<ClaimedItmes {claimedItems} />
