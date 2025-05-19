<script>
	import { authenticateUser, user } from '$lib/stores/auth';
	import { authVerify, handleTryAgain, preventNonNumericInput, logout, emailAuth } from './helpers';
	let { showAuthModal = $bindable(false) } = $props();
	let email = $state('');
	let code = $state('');
	let status = $state('start');
	let isEmailValid = $derived(email.includes('@') && email.length >= 5);

	function closeModal() {
		showAuthModal = false;
		email = '';
		status = 'start';
		code = '';
	}

	$effect(() => {
		status = $user.isAuthenticated ? 'logout' : 'start';
	});

	async function handleSendCode() {
		status = await emailAuth(email, code);
	}

	function handleCodeInput(e) {
		const digitsOnly = e.target.value.replace(/\D/g, '');
		code = digitsOnly.substring(0, 6);
		if (code.length === 6) {
			handleCodeSubmit();
		}
	}

	async function handleCodeSubmit() {
		const result = await authVerify(email, code);
		if (result.status === 'success') {
			authenticateUser(result.user);
			closeModal();
		} else {
			status = 'error';
		}
	}

	async function handleLogout() {
		await logout();
		closeModal();
	}
</script>

{#if showAuthModal}
	<div
		class="fixed inset-0 bg-[rgba(0,0,0,0.5)] z-50 flex items-center justify-center h-screen overflow-hidden p-4"
	>
		<div class="bg-white p-6 rounded-lg shadow-lg max-w-sm w-full">
			{#if status === 'start'}
				<p class="text-gray-600 mb-4">
					Enter your email address and we'll send you a 6 digit code from a gmail account.
				</p>
				<input
					type="email"
					bind:value={email}
					class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:outline-none"
					placeholder="your@email.com"
				/>
				<div class="flex gap-2 justify-between mt-4">
					<button
						onclick={handleSendCode}
						disabled={!isEmailValid}
						class="py-2 px-4 rounded-md text-white
					transition-colors duration-200
					{isEmailValid ? 'bg-blue-600 hover:bg-blue-700 cursor-pointer' : 'bg-gray-400 cursor-not-allowed'}"
					>
						Send Code
					</button>
					<button
						onclick={closeModal}
						class="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-100 cursor-pointer"
					>
						Cancel
					</button>
				</div>
			{:else if status === 'enter-code'}
				<p class="text-gray-600 mb-4">You have 15 minutes to enter this code to authenticate.</p>
				<div class="flex gap-2 justify-between w-full">
					<input
						type="text"
						value={code}
						oninput={handleCodeInput}
						onkeypress={(e) => preventNonNumericInput(e)}
						inputmode="numeric"
						pattern="[0-9]*"
						maxlength="6"
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:outline-none text-center text-xl tracking-wider"
						placeholder="Enter 6-digit code"
					/>
				</div>
				<div class="flex justify-end mt-4">
					<button
						onclick={closeModal}
						class="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-100 cursor-pointer"
					>
						Cancel
					</button>
				</div>
			{:else if status === 'logout'}
				<p class="text-gray-600 mb-4">Are you sure you want to log out?</p>
				<div class="flex gap-2 justify-between mt-4">
					<button
						onclick={handleLogout}
						class="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-100 cursor-pointer"
						>Logout</button
					>
					<button
						onclick={() => (showAuthModal = false)}
						class="py-2 px-4 rounded-md text-white
					transition-colors duration-200
					bg-blue-600 hover:bg-blue-700 cursor-pointer">cancel</button
					>
				</div>
			{:else if status === 'error'}
				<p class="text-red-600 mb-4">Invalid code. Please try again.</p>
				<div class="flex gap-2 justify-between mt-4">
					<button
						onclick={() => handleTryAgain(status)}
						class="py-2 px-4 rounded-md text-white
					transition-colors duration-200 bg-blue-600 hover:bg-blue-700 cursor-pointer
					"
					>
						Try Again
					</button>
					<button
						onclick={closeModal}
						class="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-100 cursor-pointer"
					>
						Cancel
					</button>
				</div>
			{/if}
		</div>
	</div>
{/if}
