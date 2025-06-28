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
	<div class="flex-vertical">
		<h2 class="text-title">Send Code to Email</h2>
		<p class="text-info">
			Enter your email address and we'll send you a 6 digit code from a gmail account.
		</p>
		<input type="email" bind:value={email} class="input" placeholder="your@email.com" />
	</div>

	<div class="flex-buttons mt-6">
		<button onclick={handleSendCode} disabled={!isEmailValid} class="btn"> Send Code </button>
		<button onclick={onClose} class="btn-close"> Cancel </button>
	</div>
{:else if status === 'enter-code'}
	<div class="flex-vertical">
		<p class="text-subtitle">You have 15 minutes to enter this code to authenticate.</p>

		<input
			type="text"
			value={code}
			oninput={handleCodeInput}
			onkeypress={(e) => preventNonNumericInput(e)}
			inputmode="numeric"
			pattern="[0-9]*"
			maxlength="6"
			class="input"
			placeholder="Enter 6-digit code"
		/>
	</div>
	<div class="flex justify-end mt-4">
		<button onclick={onClose} class="btn-close"> Cancel </button>
	</div>
{/if}
