<script>
	import { setUser } from '$lib/stores/auth';
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

	async function handleSubmit() {
		try {
			const url = `${import.meta.env.VITE_API_URL}/auth`;
			const response = await fetch(url, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ email: email.toLowerCase() })
			});

			if (response.ok) {
				status = 'enter-code';
			} else {
				const errorData = await response.json();
				console.error('Server returned an error:', response.status, errorData);
			}
		} catch (error) {
			console.error('Network error sending auth email:', error);
		}
	}

	function handleCodeInput(e) {
		const digitsOnly = e.target.value.replace(/\D/g, '');
		code = digitsOnly.substring(0, 6);
		if (code.length === 6) {
			handleCodeSubmit();
		}
	}

	function preventNonNumericInput(e) {
		if (!/[0-9]/.test(e.key)) {
			e.preventDefault();
		}
	}

	async function handleCodeSubmit() {
		try {
			const url = `${import.meta.env.VITE_API_URL}/auth/verify`;
			const response = await fetch(url, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				credentials: 'include',
				body: JSON.stringify({ email: email.toLowerCase(), code })
			});

			if (response.ok) {
				const data = await response.json();
				setUser(data.user);
				closeModal();
			} else {
				status = 'error';
			}
		} catch (error) {
			console.error('Network error sending auth code:', error);
		}
	}

	function handleTryAgain() {
		status = 'enter-code';
	}
</script>

{#if showAuthModal}
	<div
		class="fixed inset-0 bg-[rgba(0,0,0,0.5)] z-50 flex items-center justify-center h-screen overflow-hidden p-4"
	>
		<div class="bg-white p-6 rounded-lg shadow-lg max-w-sm w-full">
			{#if status === 'start'}
				<p class="text-gray-600 mb-4">
					Enter your email address and we'll send you a login link. It'll come from Helen's son's
					account (jmetzg11) and expire in 15 minutes.
				</p>
				<input
					type="email"
					bind:value={email}
					class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:outline-none"
					placeholder="your@email.com"
				/>
				<div class="flex gap-2 justify-between mt-4">
					<button
						onclick={handleSubmit}
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
						onkeypress={preventNonNumericInput}
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
			{:else if status === 'error'}
				<p class="text-red-600 mb-4">Invalid code. Please try again.</p>
				<div class="flex gap-2 justify-between mt-4">
					<button
						onclick={handleTryAgain}
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
