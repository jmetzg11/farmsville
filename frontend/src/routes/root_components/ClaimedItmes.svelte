<script>
	import { formatDate } from '$lib/helpers.js';
	import { claimedItems } from '$lib/stores/items';
</script>

<div class="section-container">
	<h2 class="text-2xl font-bold text-gray-800 mb-6 border-b pb-2">Items Claimed</h2>

	<div class="grid-display">
		{#each $claimedItems as item}
			<div class="item-card">
				<div class="p-4">
					<div class="item-row">
						<span class="card-title">{item.item_name}</span>
						<span class="card-qty-display">Qty: {item.quantity}</span>
					</div>
					<div class="item-row mt-6">
						<span class="card-text"
							>{item.user_name ? item.user_name : item.user_email.split('@')[0]}</span
						>
						<span class="card-text">{formatDate(item.created_at)}</span>
					</div>
				</div>
			</div>
		{/each}
	</div>
</div>
