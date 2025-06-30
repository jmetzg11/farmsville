<script>
	import { user } from '$lib/stores/auth';
	import AuthModal from '$lib/AuthModal/AuthModal.svelte';
	import Footer from './Footer.svelte';

	let { children } = $props();
	import '../app.css';

	let showAuthModal = $state(false);
</script>

<div class="header-container">
	<header class="header">
		<div class="header-wrapper">
			<div class="text-xl font-bold">
				Welcome, {$user.name ? $user.name : $user.email ? $user.email.split('@')[0] : 'Guest'}
			</div>
			<nav>
				<ul class="header-ul">
					<li>
						<a href="/" class="header-li">Inventory</a>
					</li>
					<li>
						<a href="/messages" class="header-li">Messages</a>
					</li>
					<li>
						<a href="/blog" class="header-li">Blog</a>
					</li>
					<li>
						<button onclick={() => (showAuthModal = true)} class="header-li cursor-pointer">
							{$user.isAuthenticated ? 'Logout' : 'Login'}
						</button>
					</li>
					{#if $user.admin}
						<li><a href="/admin" class="header-li">Admin</a></li>
					{/if}
				</ul>
			</nav>
		</div>
	</header>

	<AuthModal bind:showAuthModal />

	<main class="flex-grow">
		{@render children()}
	</main>
	<Footer />
</div>
