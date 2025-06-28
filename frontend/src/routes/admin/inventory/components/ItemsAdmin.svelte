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

<div class="section-container">
	<h2 class="text-title-section">Items Available</h2>
	<div class="grid-display">
		{#each $items as item}
			<div class="item-card">
				<div class="p-4">
					<div class="item-row mb-2">
						<span class="card-title">{item.name}</span>
						<span class="card-text">{formatDate(item.created_at)}</span>
					</div>
					{#if item.photo_path}
						<button
							type="button"
							class="card-image-wrapper"
							onclick={() => openPhotoModal(item.photo_path)}
							onkeydown={(e) => e.key === 'Enter' && openPhotoModal(item.photo_path)}
						>
							<img
								src={`${import.meta.env.VITE_PHOTO_URL}/${item.photo_path}`}
								alt={item.name}
								class="card-image"
							/>
						</button>
					{/if}
					<p class="card-text mb-4">
						{item.description}
					</p>
					<div class="flex-between-center">
						<div>
							<div class="card-quantity-bar-wrapper">
								<div
									class="card-quantity-bar"
									style="width: {(item.remaining_quantity / item.quantity) * 100}%"
								></div>
							</div>
							<span class="card-text">{item.remaining_quantity} / {item.quantity} remaining </span>
						</div>
					</div>
					<div class="flex gap-2 justify-between mt-8">
						<button onclick={() => handleClaim(item)} class="btn"> Claim </button>
						<button onclick={() => handleEdit(item)} class="btn-secondary"> Edit </button>
					</div>
				</div>
			</div>
		{/each}
	</div>
</div>

<EditItemModal bind:showEditModal {selectedItem} />
<ClaimItemModal bind:showClaimModal {selectedItem} />
<PhotoModal bind:showPhotoModal {photoUrl} />
