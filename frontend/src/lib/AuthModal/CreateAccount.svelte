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
		if (result.success) {
			onSuccess(result.user);
		} else {
			message = result.message;
			previousStatus = 'create-account';
			status = 'error';
		}
	}
</script>

<div class="flex-vertical">
	<h2 class="text-title">Login with Password</h2>
	<p class="text-subtitle">* Email and password are required</p>
	<input type="text" bind:value={name} class="input" placeholder="Your Name" />
	<input
		type="text"
		bind:value={phone}
		class="input"
		placeholder="Phone Number e.g. (123) 456-7890"
	/>
	<input type="email" bind:value={email} class="input" placeholder="your@email.com" />
	<input type="text" bind:value={password} class="input" placeholder="password" />
	<button onclick={handleCreateAccount} disabled={!isEmailValid || !isPasswordValid} class="btn">
		Create Account
	</button>
	<button onclick={onClose} class="btn-close"> Cancel </button>
</div>
