<script>
	import { createAccount } from '$lib/api_calls/auth';

	let {
		onSuccess,
		onClose,
		status = $bindable('create-account'),
		message = $bindable(''),
		previousStatus = $bindable('')
	} = $props();

	let email = $state('');
	let password = $state('');
	let name = $state('');
	let phone = $state('');

	let isEmailValid = $derived(email.includes('@') && email.length >= 5);
	let isPasswordValid = $derived(password.length > 0);

	async function handleCreateAccount() {
		email = email.toLowerCase();
		const result = await createAccount(name, phone, email, password);
		if (result.status === 'success') {
			onSuccess(result.user);
		} else {
			message = result.message;
			previousStatus = 'create-account';
			status = 'error';
		}
	}
</script>

<div class="flex flex-col gap-4 mb-6 pb-2">
	<h2 class="text-lg font-bold text-gray-900 text-center">Login with Password</h2>
	<p class="text-gray-600 text-center">* Email and password are required</p>
	<input
		type="text"
		bind:value={name}
		class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:outline-none"
		placeholder="Your Name"
	/>
	<input
		type="text"
		bind:value={phone}
		class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:outline-none"
		placeholder="Phone Number e.g. (123) 456-7890"
	/>
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
		placeholder="password"
	/>
	<button
		onclick={handleCreateAccount}
		disabled={!isEmailValid || !isPasswordValid}
		class="py-2 px-4 rounded-md text-white transition-colors duration-200 {isEmailValid &&
		isPasswordValid
			? 'bg-blue-600 hover:bg-blue-700 cursor-pointer'
			: 'bg-gray-400 cursor-not-allowed'}"
	>
		Create Account
	</button>
	<button
		onclick={onClose}
		class="py-2 px-4 border border-gray-300 rounded-md hover:bg-gray-100 cursor-pointer"
	>
		Cancel
	</button>
</div>
