<script>
	let { onClose, status = $bindable('start') } = $props();

	const authOptions = [
		{
			id: 'auth-code',
			title: 'Authentication Code',
			description: 'email or text',
			status: 'auth-code'
		},
		{
			id: 'login-password',
			title: 'Login with Password',
			description: 'Use your email and password',
			status: 'login-password'
		},
		{
			id: 'create-account',
			title: 'Create Account',
			description: 'First time users',
			status: 'create-account'
		},
		{
			id: 'reset-password',
			title: 'Reset Password',
			description: 'Forgot or never set',
			status: 'reset-password'
		}
	];
</script>

<div class="container mx-auto p-2">
	<h1 class="text-2xl font-bold text-gray-800 mb-6 pb-2 text-center">Choose Login Method</h1>

	<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
		{#each authOptions as option}
			<button
				class="bg-white p-4 rounded-lg shadow-md border border-gray-200 hover:shadow-lg transition-shadow cursor-pointer"
				id={option.id}
				onclick={() => (status = option.status)}
			>
				<div class="text-center">
					<h2 class="text-lg font-bold text-gray-800 mb-2">{option.title}</h2>
					<p class="text-gray-600 text-sm leading-relaxed">
						{option.description}
					</p>
				</div>
			</button>
		{/each}
	</div>

	<div class="flex justify-center mt-8">
		<button
			onclick={onClose}
			class="bg-gray-500 hover:bg-gray-600 text-white font-medium py-2 px-6 rounded transition-colors cursor-pointer w-full"
		>
			Cancel
		</button>
	</div>
</div>
