<script>
	let { onTryAgain, onClose } = $props();
</script>

<p class="text-red-600 mb-4">Invalid code. Please try again.</p>
<div class="flex gap-2 justify-between mt-4">
	<button
		onclick={onTryAgain}
		class="py-2 px-4 rounded-md text-white
    transition-colors duration-200 bg-blue-600 hover:bg-blue-700 cursor-pointer
    "
	>
		Try Again
	</button>
	<button
		onclick={onClose}
		class="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-100 cursor-pointer"
	>
		Cancel
	</button>
</div>
