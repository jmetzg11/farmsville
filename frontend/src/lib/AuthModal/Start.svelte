<script>
	let { onClose, status = $bindable('start') } = $props();

	const authOptions = [
		{
			id: 'auth-code',
			title: 'Authentication Code',
			description: 'email or text'
		},
		{
			id: 'login-password',
			title: 'Login with Password',
			description: 'Use your email and password'
		},
		{
			id: 'create-account',
			title: 'Create Account',
			description: 'First time users'
		},
		{
			id: 'reset-password',
			title: 'Reset Password',
			description: 'Forgot or never set'
		}
	];
</script>

<div class="container mx-auto p-2">
	<h1 class="text-title mb-4">Choose Login Method</h1>

	<div class="two-col mb-6">
		{#each authOptions as option}
			<button class="clickable-card" id={option.id} onclick={() => (status = option.id)}>
				<div class="text-center">
					<h2 class="text-title-small mb-2">{option.title}</h2>
					<p class="text-gray-600 text-sm leading-relaxed">
						{option.description}
					</p>
				</div>
			</button>
		{/each}
	</div>

	<div class="flex justify-end">
		<button onclick={onClose} class="btn-close"> Cancel </button>
	</div>
</div>
