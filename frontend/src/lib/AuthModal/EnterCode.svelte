<script>
	import { authVerify } from '$lib/api_calls/auth';
	import { preventNonNumericInput } from '$lib/helpers';

	let { email, onSuccess, onError, onClose } = $props();

	let code = $state('');

	function handleCodeInput(e) {
		const digitsOnly = e.target.value.replace(/\D/g, '');
		code = digitsOnly.substring(0, 6);
		if (code.length === 6) {
			handleCodeSubmit();
		}
	}

	async function handleCodeSubmit() {
		const result = await authVerify(email, code);
		if (result.status === 'success') {
			onSuccess(result);
		} else {
			onError();
		}
		code = '';
	}
</script>

<p class="text-gray-600 mb-4">You have 15 minutes to enter this code to authenticate.</p>
<div class="flex gap-2 justify-between w-full">
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
</div>
<div class="flex justify-end mt-4">
	<button
		onclick={onClose}
		class="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-100 cursor-pointer"
	>
		Cancel
	</button>
</div>
