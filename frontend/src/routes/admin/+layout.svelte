<script>
	let { children } = $props();
	import { user } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	import { initializeUserStore } from '$lib/stores/auth';
	import { refreshItems } from '$lib/stores/items';
	import AuthModal from '$lib/AuthModal/AuthModal.svelte';
	let showAuthModal = $state(false);

	onMount(async () => {
		await initializeUserStore();
		await refreshItems();
		if (!$user.isAuthenticated) {
			showAuthModal = true;
		}
	});
</script>

{#if $user.admin}
	<div class="flex flex-col">
		<header class="admin-header">
			<div class="header-wrapper">
				<nav>
					<ul class="header-ul">
						<li><a href="/admin" class="admin-header-li">Dashboard</a></li>
						<li>
							<a href="/admin/inventory" class="admin-header-li">Inventory</a>
						</li>
						<li><a href="/admin/users" class="admin-header-li">Users</a></li>
						<li>
							<a href="/admin/messages" class="admin-header-li">Messages</a>
						</li>
						<li>
							<a href="/admin/blog" class="admin-header-li">Blog</a>
						</li>
					</ul>
				</nav>
			</div>
		</header>

		<main class="section-container">
			{@render children()}
		</main>
	</div>
{:else}
	<p>You are not authorized to access this page</p>
{/if}

<AuthModal bind:showAuthModal />
