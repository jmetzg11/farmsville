<script>
	import { user } from '$lib/stores/auth';
	import AuthModal from '$lib/AuthModal/AuthModal.svelte';
	import MakeClaimModal from './MakeClaimModal.svelte';
	import PhotoModal from '../../lib/PhotoModal.svelte';
	import { formatDate } from '$lib/helpers.js';
	import { items } from '$lib/stores/items';

	let showAuthModal = $state(false);
	let showClaimModal = $state(false);
	let showPhotoModal = $state(false);

	let selectedItem = $state(null);
	let photoUrl = $state('');

	function handleClick(item) {
		selectedItem = item;
		if ($user.isAuthenticated) {
			showClaimModal = true;
		} else {
			showAuthModal = true;
		}
	}

	function openPhotoModal(photoPath) {
		photoUrl = `${import.meta.env.VITE_PHOTO_URL}/${photoPath}`;
		showPhotoModal = true;
	}
</script>

<div class="section-container">
	<h2 class="text-2xl font-bold text-gray-800 mb-6 border-b pb-2">Items Available</h2>
	<div class="grid-display">
		{#each $items as item}
			<div class="item-card">
				<div class="p-4">
					<div class="flex-between-center mb-2">
						<h3 class="card-title">{item.name}</h3>
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
					<p class="card-text-no-center mt-4">{item.description}</p>
					<div class="flex-between-center mt-4">
						<div class="space-y-2">
							<div class="card-quantity-bar-wrapper">
								<div
									class="card-quantity-bar"
									style="width: {(item.remaining_quantity / item.quantity) * 100}%"
								></div>
							</div>
							<span class="card-text">{item.remaining_quantity} / {item.quantity} remaining</span>
						</div>
						<button
							onclick={() => handleClick(item)}
							disabled={item.remaining_quantity <= 0}
							class="btn"
						>
							{item.remaining_quantity <= 0 ? 'Out of Stock' : 'Claim'}
						</button>
					</div>
				</div>
			</div>
		{/each}
	</div>
</div>

<AuthModal bind:showAuthModal />

<MakeClaimModal bind:showClaimModal {selectedItem} />

<PhotoModal bind:showPhotoModal {photoUrl} />
