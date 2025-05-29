<script>
	import { login } from '$lib/api_calls/auth';

	let {
		onSuccess,
		status = $bindable('login-password'),
		previousStatus = $bindable('start'),
		message = $bindable(''),
		onClose
	} = $props();

	let email = $state('');
	let password = $state('');

	let isEmailValid = $derived(email.includes('@') && email.length >= 5);
	let isPasswordValid = $derived(password.length >= 0);
	async function handleLogin() {
		const result = await login(email, password);
		if (result.status === 'success') {
			onSuccess(result.user);
		} else {
			status = 'error';
			message = 'Invalid email or password. Please try again.';
			previousStatus = 'login-password';
		}
	}
</script>

<div class="flex flex-col gap-4 mb-6 pb-2">
	<h2 class="text-lg font-bold text-gray-900 text-center">Login with Password</h2>
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
		onclick={handleLogin}
		disabled={!isEmailValid || !isPasswordValid}
		class="py-2 px-4 rounded-md text-white transition-colors duration-200 {isEmailValid &&
		isPasswordValid
			? 'bg-blue-600 hover:bg-blue-700 cursor-pointer'
			: 'bg-gray-400 cursor-not-allowed'}"
	>
		Login
	</button>
	<button
		onclick={onClose}
		class="py-2 px-4 border border-gray-300 rounded-md hover:bg-gray-100 cursor-pointer"
	>
		Cancel
	</button>
</div>
