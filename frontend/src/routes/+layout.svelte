<script>
	import { user } from '$lib/stores/auth';
	import AuthModal from '$lib/AuthModal/AuthModal.svelte';
	import Footer from './Footer.svelte';

	let { children } = $props();
	import '../app.css';

	let showAuthModal = $state(false);
</script>

<div class="min-h-screen flex flex-col">
	<header class="bg-teal-500 text-white shadow-md">
		<div class="container mx-auto px-4 py-3 flex flex-col sm:flex-row justify-between items-center">
			<div class="text-xl font-bold">
				Welcome, {$user.name ? $user.name : $user.email ? $user.email.split('@')[0] : 'Guest'}
			</div>
			<nav>
				<ul class="flex space-x-6">
					<li>
						<a href="/" class="hover:text-teal-200 transition-colors">Inventory</a>
					</li>
					<li>
						<a href="/messages" class="hover:text-teal-200 transition-colors">Messages</a>
					</li>
					<li>
						<a href="/blog" class="hover:text-teal-200 transition-colors">Blog</a>
					</li>
					<li>
						<button
							onclick={() => (showAuthModal = true)}
							class="hover:text-teal-200 transition-colors cursor-pointer"
						>
							{$user.isAuthenticated ? 'Logout' : 'Login'}
						</button>
					</li>
					{#if $user.admin}
						<li><a href="/admin" class="hover:text-teal-200 transition-colors">Admin</a></li>
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
