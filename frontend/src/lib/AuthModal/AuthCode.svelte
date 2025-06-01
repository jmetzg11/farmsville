<script>
	import { emailAuth, authVerify } from '$lib/api_calls/auth';
	import { preventNonNumericInput } from '$lib/helpers';

	let {
		onClose,
		onSuccess,
		status = $bindable('auth-code'),
		email = $bindable(''),
		message = $bindable(''),
		previousStatus = $bindable('start')
	} = $props();

	let code = $state('');

	let isEmailValid = $derived(email.includes('@') && email.length >= 5);

	async function handleSendCode() {
		const result = await emailAuth(email);
		if (result === 'enter-code') {
			status = 'enter-code';
		} else {
			message = 'Invalid email address. Please try again.';
			previousStatus = 'auth-code';
			status = 'error';
		}
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
			onSuccess(result.user);
		} else {
			status = 'error';
			message = 'Invalid code. Please try again.';
			previousStatus = 'enter-code';
		}
		code = '';
	}
</script>

{#if status === 'auth-code'}
	<div class="flex flex-col gap-4 mb-6 pb-2">
		<h2 class="text-lg font-bold text-gray-900 text-center">Send Code to Email</h2>
		<p class="text-gray-600 text-center">
			Enter your email address and we'll send you a 6 digit code from a gmail account.
		</p>
		<input
			type="email"
			bind:value={email}
			class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:outline-none"
			placeholder="your@email.com"
		/>
	</div>

	<div class="flex justify-between">
		<button
			onclick={handleSendCode}
			disabled={!isEmailValid}
			class="py-2 px-4 rounded-md text-white transition-colors duration-200 {isEmailValid
				? 'bg-blue-600 hover:bg-blue-700 cursor-pointer'
				: 'bg-gray-400 cursor-not-allowed'}"
		>
			Send Code
		</button>
		<button
			onclick={onClose}
			class="py-2 px-4 border border-gray-300 rounded-md hover:bg-gray-100 cursor-pointer"
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
			onclick={onClose}
			class="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-100 cursor-pointer"
		>
			Cancel
		</button>
	</div>
{/if}
