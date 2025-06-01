<script>
	import { sendCodeToResetPassword, resetPassword } from '$lib/api_calls/auth';
	import { preventNonNumericInput } from '$lib/helpers';

	let {
		email,
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
			console.log('success', result);
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
		console.log(result);
	}
</script>

<div class="flex flex-col gap-4 mb-6 pb-2">
	{#if status === 'reset-password'}
		<h2 class="text-lg font-bold text-gray-900 text-center">Restset Your Password</h2>
		<p class="text-gray-600 text-center">A authentication code will be sent to your email</p>
		<input
			type="email"
			bind:value={email}
			class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:outline-none"
			placeholder="your@email.com"
		/>
		<button
			onclick={handleSendCode}
			disabled={!isEmailValid}
			class="py-2 px-4 rounded-md text-white transition-colors duration-200 {isEmailValid
				? 'bg-blue-600 hover:bg-blue-700 cursor-pointer'
				: 'bg-gray-400 cursor-not-allowed'}">Send Code</button
		>
		<button
			onclick={onClose}
			class="py-2 px-4 border border-gray-300 rounded-md hover:bg-gray-100 cursor-pointer"
		>
			Cancel
		</button>
	{:else if status === 'enter-code-and-password'}
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
		<input
			type="text"
			bind:value={password}
			class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:outline-none"
			placeholder="password"
		/>
		<input
			type="text"
			bind:value={confirmPassword}
			class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:outline-none"
			placeholder="confirm password"
		/>

		<button
			onclick={handleResetPassword}
			disabled={!codeIsValid || !passwordIsValid}
			class="py-2 px-4 rounded-md text-white transition-colors duration-200 {codeIsValid &&
			passwordIsValid
				? 'bg-blue-600 hover:bg-blue-700 cursor-pointer'
				: 'bg-gray-400 cursor-not-allowed'}"
		>
			Reset Password
		</button>
		<button
			onclick={onClose}
			class="py-2 px-4 border border-gray-300 rounded-md hover:bg-gray-100 cursor-pointer"
		>
			Cancel
		</button>
	{/if}
</div>
