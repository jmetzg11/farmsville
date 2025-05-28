<script>
	import { emailAuth } from '$lib/api_calls/auth';

	let { status = $bindable('auth-code'), email = $bindable(''), onClose } = $props();

	let isEmailValid = $derived(email.includes('@') && email.length >= 5);

	async function handleSendCode() {
		status = await emailAuth(email);
	}
</script>

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
