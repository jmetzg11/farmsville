<script>
	let { showModal = $bindable(false) } = $props();
	let email = $state('');

	let isEmailValid = $derived(email.includes('@') && email.length >= 5);

	function closeModal() {
		showModal = false;
	}

	async function handleSubmit() {
		try {
			const url = `${import.meta.env.VITE_API_URL}/auth`;
			const response = await fetch(url, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ email })
			});
			if (response.ok) {
				console.log('success');
			}
		} catch (error) {
			console.error('Error sending auth email');
		}
	}
</script>

{#if showModal}
	<div class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
		<div class="bg-white p-6 rounded-lg shadow-lg max-w-sm w-full">
			<p class="text-gray-600 mb-4">
				Enter your email address and we'll send you a login link. It'll come from Helen's son's
				account (jmetzg11) and expire in 15 minutes.
			</p>
			<input
				type="email"
				bind:value={email}
				class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:outline-none"
				placeholder="your@email.com"
			/>
			<div class="flex gap-2 justify-between mt-4">
				<button
					onclick={handleSubmit}
					disabled={!isEmailValid}
					class="py-2 px-4 rounded-md text-white
						transition-colors duration-200
						{isEmailValid ? 'bg-blue-600 hover:bg-blue-700 cursor-pointer' : 'bg-gray-400 cursor-not-allowed'}"
				>
					Send Login Link
				</button>
				<button
					onclick={closeModal}
					class="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-100 cursor-pointer"
				>
					Cancel
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- add state to enter code, expires in 15 mins, email from jmetzg11 (Helen's Son) -->
