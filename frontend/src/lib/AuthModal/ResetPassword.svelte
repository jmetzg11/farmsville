<script>
	import { sendCodeToResetPassword, resetPassword } from '$lib/api_calls/auth';
	import { preventNonNumericInput } from '$lib/helpers';

	let {
		email,
		onSuccess,
		onClose,
		status = $bindable('reset-password'),
		message = $bindable(''),
		previousStatus = $bindable('')
	} = $props();

	let code = $state('');
	let password = $state('');
	let confirmPassword = $state('');

	let codeIsValid = $state(false);
	let passwordIsValid = $derived(password.length > 0 && password === confirmPassword);
	let isEmailValid = $derived(email.includes('@') && email.length >= 5);

	async function handleSendCode() {
		email = email.toLowerCase();
		const result = await sendCodeToResetPassword(email);
		if (result.success) {
			status = 'enter-code-and-password';
		} else {
			message = result.message;
			previousStatus = 'reset-password';
			status = 'error';
		}
	}

	function handleCodeInput(e) {
		const digitsOnly = e.target.value.replace(/\D/g, '');
		code = digitsOnly.substring(0, 6);
		codeIsValid = code.length === 6;
	}

	async function handleResetPassword() {
		const result = await resetPassword(email, code, password);
		if (result.success) {
			onSuccess(result.user);
		} else {
			message = result.message;
			previousStatus = 'reset-password';
			status = 'error';
		}
	}
</script>

<div class="flex-vertical">
	{#if status === 'reset-password'}
		<h2 class="text-title">Restset Your Password</h2>
		<p class="text-subtitle">A authentication code will be sent to your email</p>
		<input type="email" bind:value={email} class="input" placeholder="your@email.com" />
		<button onclick={handleSendCode} disabled={!isEmailValid} class="btn">Send Code</button>
		<button onclick={onClose} class="btn-close"> Cancel </button>
	{:else if status === 'enter-code-and-password'}
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
		<input type="text" bind:value={password} class="input" placeholder="password" />
		<input type="text" bind:value={confirmPassword} class="input" placeholder="confirm password" />

		<button onclick={handleResetPassword} disabled={!codeIsValid || !passwordIsValid} class="btn">
			Reset Password
		</button>
		<button onclick={onClose} class="btn-close"> Cancel </button>
	{/if}
</div>
