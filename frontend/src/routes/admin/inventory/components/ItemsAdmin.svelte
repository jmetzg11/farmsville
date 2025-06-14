<script>
	import { items } from '$lib/stores/items';
	import { formatDate } from '$lib/helpers.js';
	import EditItemModal from './EditItemModal.svelte';
	import ClaimItemModal from './ClaimItemModal.svelte';
	import PhotoModal from '$lib/PhotoModal.svelte';

	let selectedItem = $state(null);
	let showEditModal = $state(false);
	let showClaimModal = $state(false);
	let showPhotoModal = $state(false);
	let photoUrl = $state('');

	function handleEdit(item) {
		selectedItem = item;
		showEditModal = true;
	}

	function handleClaim(item) {
		selectedItem = item;
		showClaimModal = true;
	}

	function openPhotoModal(photoPath) {
		photoUrl = `${import.meta.env.VITE_PHOTO_URL}/${photoPath}`;
		showPhotoModal = true;
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
					{#if item.photo_path}
						<button
							type="button"
							class="h-28 w-full overflow-hidden mb-2 cursor-pointer rounded-lg"
							onclick={() => openPhotoModal(item.photo_path)}
							onkeydown={(e) => e.key === 'Enter' && openPhotoModal(item.photo_path)}
						>
							<img
								src={`${import.meta.env.VITE_PHOTO_URL}/${item.photo_path}`}
								alt={item.name}
								class="w-full h-full object-contain"
							/>
						</button>
					{/if}
					<p class="text-gray-600 mb-3 text-sm {item.photo_path ? 'text-center' : ''}">
						{item.description}
					</p>
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
<PhotoModal bind:showPhotoModal {photoUrl} />
