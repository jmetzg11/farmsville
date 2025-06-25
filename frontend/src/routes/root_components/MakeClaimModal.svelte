<script>
	import { formatDate } from '$lib/helpers.js';
	import { makeClaim } from '$lib/api_calls/customer.js';
	import { refreshItems } from '$lib/stores/items';
	let { showClaimModal = $bindable(false), selectedItem = null } = $props();
	let quantity = $state(1);

	function closeModal() {
		showClaimModal = false;
	}

	async function handleSubmit() {
		const result = await makeClaim(selectedItem.id, quantity);
		if (result) {
			quantity = 1;
			showClaimModal = false;
			refreshItems();
		} else {
			throw new Error('Failed to make claim');
		}
	}
</script>

{#if showClaimModal}
	<div class="modal-container">
		<div class="modal-content">
			<div class="flex justify-between items-center mb-4">
				<h3 class="text-xl font-bold text-gray-800">{selectedItem.name}</h3>
				<span class="text-md text-gray-500">{formatDate(selectedItem.created_at)}</span>
			</div>
			<p class="text-gray-600 mb-3 text-md">{selectedItem.description}</p>
			<div>
				<input
					type="number"
					bind:value={quantity}
					min="1"
					max={selectedItem.remaining_quantity}
					class="w-18 text-lg text-center px-1 py-0 border border-gray-300 rounded-md text-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
				/
				<span class="text-gray-700">{selectedItem.remaining_quantity}</span>
			</div>
			<div class="flex gap-2 justify-between mt-4">
				<button onclick={handleSubmit} disabled={quantity <= 0} class="btn">Claim</button>
				<button onclick={closeModal} class="btn-close">Close</button>
			</div>
		</div>
	</div>
{/if}
