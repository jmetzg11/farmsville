<script>
	import { formatDate } from '$lib/helpers.js';
	import { claimedItems } from '$lib/stores/items';
	import { removeClaimedItem } from '$lib/api_calls/admin_items.js';
	import { refreshItems } from '$lib/stores/items';

	async function handleRemove(itemId) {
		await removeClaimedItem(itemId);
		await refreshItems();
	}
</script>

<div class="section-container">
	<h2 class="text-title-section">Items Claimed</h2>
	<div class="grid-display">
		{#each $claimedItems as item}
			<div class="item-card">
				<div class="p-4">
					<div class="item-row">
						<span class="card-title">{item.item_name}</span>
						<div class="flex items-center gap-2">
							<span class="card-qty-display">Qty: {item.quantity}</span>
							<button
								on:click={() => handleRemove(item.id)}
								class="card-x-red-circle"
								aria-label="Remove item"
							>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class="h-4 w-4"
									fill="none"
									viewBox="0 0 24 24"
									stroke="currentColor"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M6 18L18 6M6 6l12 12"
									/>
								</svg>
							</button>
						</div>
					</div>
					<div class="item-row mt-4">
						<span class="card-text font-medium"
							>{item.user_name ? item.user_name : item.user_email}</span
						>
						<span class="card-text">{formatDate(item.created_at)}</span>
					</div>
				</div>
			</div>
		{/each}
	</div>
</div>
