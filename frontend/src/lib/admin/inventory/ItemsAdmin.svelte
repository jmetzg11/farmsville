<script>
	import { items } from '$lib/stores/items';
	import { formatDate } from '$lib/root/helpers';
	import EditItemModal from './EditItemModal.svelte';
	import ClaimItemModal from './ClaimItemModal.svelte';
	let selectedItem = $state(null);
	let showEditModal = $state(false);
	let showClaimModal = $state(false);

	function handleEdit(item) {
		selectedItem = item;
		showEditModal = true;
	}

	function handleClaim(item) {
		selectedItem = item;
		showClaimModal = true;
	}
</script>

<div class="container mx-auto px-4 py-6">
	<h2 class="text-2xl font-bold text-gray-800 mb-6 border-b pb-2">Items Available</h2>
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 p-4">
		{#each $items as item}
			<div
				class="bg-white rounded-lg shadow-md overflow-hidden border border-gray-200 hover:shadow-lg transition-shadow"
			>
				<div class="p-4">
					<div class="flex justify-between items-center mb-2">
						<h3 class="text-lg font-bold text-gray-800">{item.name}</h3>
						<span class="text-xs text-gray-500">{formatDate(item.created_at)}</span>
					</div>
					<p class="text-gray-600 mb-3 text-sm">{item.description}</p>
					<div class="flex justify-between items-center mt-4">
						<div>
							<div class="w-full bg-gray-200 rounded-full h-2.5">
								<div
									class="bg-green-600 h-2.5 rounded-full"
									style="width: {(item.remaining_quantity / item.quantity) * 100}%"
								></div>
							</div>
							<div class="text-sm text-gray-600 mt-1">
								<span class="font-medium">{item.remaining_quantity}</span> / {item.quantity} remaining
							</div>
						</div>
					</div>
					<div class="flex gap-2 justify-between mt-8">
						<button
							onclick={() => handleClaim(item)}
							class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-1 px-4 rounded text-sm transition-colors cursor-pointer"
						>
							Claim
						</button>
						<button
							onclick={() => handleEdit(item)}
							class="bg-orange-500 hover:bg-orange-600 text-white font-medium py-1 px-4 rounded text-sm transition-colors cursor-pointer"
						>
							Edit
						</button>
					</div>
				</div>
			</div>
		{/each}
	</div>
</div>

<EditItemModal bind:showEditModal {selectedItem} />
<ClaimItemModal bind:showClaimModal {selectedItem} />
