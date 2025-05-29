<script>
	import { authenticateUser, user } from '$lib/stores/auth';
	import Start from './Start.svelte';
	import AuthenticationCode from './AuthenticationCode.svelte';
	import LoginWithPassword from './LoginWithPassword.svelte';
	import EnterCode from './EnterCode.svelte';
	import CreateAccount from './CreateAccount.svelte';
	import Logout from './Logout.svelte';
	import Error from './Error.svelte';
	let { showAuthModal = $bindable(false) } = $props();
	let email = $state('');
	let message = $state('');
	let status = $state('start');
	let previousStatus = $state('start');

	function closeModal() {
		showAuthModal = false;
		email = '';
		message = '';
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
				<AuthenticationCode
					onClose={() => (status = 'start')}
					bind:status
					bind:email
					bind:message
					bind:previousStatus
				/>
			{:else if status === 'login-password'}
				<LoginWithPassword
					onSuccess={(user) => {
						authenticateUser(user);
						closeModal();
					}}
					bind:status
					bind:previousStatus
					bind:message
					onClose={() => (status = 'start')}
				/>
			{:else if status === 'enter-code'}
				<EnterCode
					{email}
					onSuccess={(user) => {
						authenticateUser(user);
						closeModal();
					}}
					onClose={() => (status = 'start')}
					bind:status
					bind:message
					bind:previousStatus
				/>
			{:else if status == 'create-account'}
				<CreateAccount onClose={() => (status = 'start')} bind:status />
			{:else if status === 'logout'}
				<Logout onCancel={() => (showAuthModal = false)} />
			{:else if status === 'error'}
				<Error bind:status onClose={closeModal} {message} {previousStatus} />
			{/if}
		</div>
	</div>
{/if}
