<script>
	import { formatDate } from '../root/helpers.js';
	import { claimedItems } from '$lib/stores/items';
</script>

<div class="container mx-auto px-4 py-6">
	<h2 class="text-2xl font-bold text-gray-800 mb-6 border-b pb-2">Items Claimed (Admin)</h2>

	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
		{#each $claimedItems as item}
			<div
				class="bg-white rounded-lg shadow-md overflow-hidden border border-gray-200 hover:shadow-lg transition-shadow"
			>
				<div class="p-4">
					<div class="flex justify-between items-center">
						<span class="font-medium text-gray-800">{item.item_name}</span>
						<span class="font-medium bg-blue-100 text-blue-800 px-2 py-1 rounded-full text-sm">
							Qty: {item.quantity}
						</span>
					</div>
					<div class="flex justify-between items-center text-sm text-gray-600 mt-4">
						<span class="font-medium">{item.user_name ? item.user_name : item.user_email}</span>
						<span class="text-gray-500">{formatDate(item.created_at)}</span>
					</div>
				</div>
			</div>
		{/each}
	</div>
</div>
