<script>
	import { emailAuth } from './helpers';

	let { email = $bindable(''), onSendCode, onClose } = $props();
	let password = $state('');

	let isLoginValid = $derived(email.includes('@') && email.length >= 5 && password.length >= 2);
	let isEmailValid = $derived(email.includes('@') && email.length >= 5);

	async function handleSendCode() {
		const status = await emailAuth(email);
		onSendCode(status);
	}
</script>

<div class="flex flex-col gap-4 mb-6 border-b pb-2">
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
	<button
		onclick={handleSendCode}
		disabled={!isEmailValid}
		class="py-2 px-4 rounded-md text-white
		transition-colors duration-200
		{isEmailValid ? 'bg-blue-600 hover:bg-blue-700 cursor-pointer' : 'bg-gray-400 cursor-not-allowed'}"
	>
		Send Code
	</button>
</div>
<div class="flex flex-col gap-4 mb-6 border-b pb-2">
	<h2 class="text-lg font-bold text-gray-900 text-center">Through Email and Password</h2>
	<p class="text-gray-600 text-center">Create an account or login into an existing account.</p>
	<input
		type="email"
		bind:value={email}
		class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:outline-none"
		placeholder="your@email.com"
	/>
	<input
		type="text"
		bind:value={password}
		class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:outline-none"
		placeholder="your password"
	/>
	<button
		onclick={onClose}
		class="py-2 px-4 rounded-md text-white
        transition-colors duration-200
        {isLoginValid
			? 'bg-blue-600 hover:bg-blue-700 cursor-pointer'
			: 'bg-gray-400 cursor-not-allowed'}"
	>
		Login
	</button>
</div>
<button
	onclick={onClose}
	class="w-full px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-100 cursor-pointer"
>
	Cancel
</button>
