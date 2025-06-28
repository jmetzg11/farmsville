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
	let isPasswordValid = $derived(password.length > 0);
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

<div class="flex-vertical">
	<h2 class="text-title">Login with Password</h2>
	<input type="email" bind:value={email} class="input" placeholder="your@email.com" />
	<input type="password" bind:value={password} class="input" placeholder="password" />
	<button onclick={handleLogin} disabled={!isEmailValid || !isPasswordValid} class="btn">
		Login
	</button>
	<button onclick={onClose} class="btn-close"> Cancel </button>
</div>
