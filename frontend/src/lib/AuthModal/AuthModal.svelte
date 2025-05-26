<script>
	import { authenticateUser, user } from '$lib/stores/auth';
	import Start from './Start.svelte';
	import AuthenticationCode from './AuthenticationCode.svelte';
	import EnterCode from './EnterCode.svelte';
	import Logout from './Logout.svelte';
	import Error from './Error.svelte';
	let { showAuthModal = $bindable(false) } = $props();
	let email = $state('');
	let status = $state('start');

	function closeModal() {
		showAuthModal = false;
		email = '';
		status = 'start';
	}

	$effect(() => {
		status = $user.isAuthenticated ? 'logout' : 'start';
	});
</script>

{#if showAuthModal}
	<div
		class="fixed inset-0 bg-[rgba(0,0,0,0.5)] z-50 flex items-center justify-center h-screen overflow-hidden p-4"
	>
		<div class="bg-white p-6 rounded-lg shadow-lg max-w-xl w-full">
			{#if status === 'start'}
				<Start onClose={closeModal} bind:status />
			{:else if status === 'auth-code'}
				<AuthenticationCode onClose={closeModal} bind:status />
			{:else if status === 'enter-code'}
				<EnterCode
					{email}
					onSuccess={(result) => {
						authenticateUser(result.user);
						closeModal();
					}}
					onError={() => (status = 'error')}
					onClose={closeModal}
				/>
			{:else if status === 'logout'}
				<Logout onCancel={() => (showAuthModal = false)} />
			{:else if status === 'error'}
				<Error onTryAgain={() => (status = 'start')} onClose={closeModal} />
			{/if}
		</div>
	</div>
{/if}
