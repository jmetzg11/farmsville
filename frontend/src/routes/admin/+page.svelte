<script>
	import { onMount } from 'svelte';
	import { user } from '$lib/stores/auth';
	import AuthModal from '$lib/AuthModal.svelte';
	import Dashboard from '$lib/admin/Dashboard.svelte';
	import { initializeUserStore } from '$lib/stores/auth';
	import { refreshItems } from '$lib/stores/items';
	let showAuthModal = $state(false);

	onMount(async () => {
		await initializeUserStore();
		await refreshItems();
		if (!$user.isAuthenticated) {
			showAuthModal = true;
		}
	});
</script>

{#if $user.admin}<Dashboard />{:else}
	<p>You are not authorized to access this page</p>
{/if}

<AuthModal bind:showAuthModal />
