<script>
	import { onMount } from 'svelte';
	import { formatDate } from '$lib/root/helpers.js';
	import { claimItem } from './helpers.js';
	import { refreshItems } from '$lib/stores/items';
	import { users, refreshUsers } from '../users/helpers.js';
	let { showClaimModal = $bindable(false), selectedItem = null } = $props();

	let amount = $state(0);

	onMount(async () => {
		await refreshUsers();
	});

	let allUsers = $state([]);

	async function handleClaim() {
		await claimItem(1, 1, 1);
		showClaimModal = false;
	}

	function handleCancel() {
		showClaimModal = false;
	}
</script>

{#if showClaimModal}
	<div
		class="fixed inset-0 bg-[rgba(0,0,0,0.5)] z-50 flex items-center justify-center h-screen overflow-hidden p-6"
	>
		<div
			class="bg-white rounded-lg shadow-md overflow-hidden border border-gray-200 hover:shadow-lg transition-shadow max-w-md w-full"
		>
			<div class="p-4">
				<h1 class="bg-amber-300 p-4">Under Construction</h1>
				<div class="flex justify-between items-center mb-4">
					<h3 class="text-lg font-bold text-gray-800">{selectedItem.name}</h3>
					<span class="text-xs text-gray-500">{formatDate(selectedItem.created_at)}</span>
				</div>
				<div class="flex justify-between items-center mb-4">
					<div>Users</div>
					<div>Amount</div>
				</div>

				<div class="flex justify-between items-center">
					<div>
						<div class="w-full bg-gray-200 rounded-full h-2.5">
							<div
								class="bg-green-600 h-2.5 rounded-full"
								style="width: {(selectedItem.remaining_quantity / selectedItem.quantity) * 100}%"
							></div>
						</div>
						<div class="text-sm text-gray-600 mt-1">
							<span class="font-medium">{selectedItem.remaining_quantity}</span> / {selectedItem.quantity}
							remaining
						</div>
					</div>
				</div>
				<div class="flex gap-2 justify-between mt-4">
					<button
						onclick={handleClaim}
						class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-1 px-4 rounded text-sm transition-colors cursor-pointer"
					>
						Claim
					</button>
					<button
						onclick={handleCancel}
						class="bg-gray-300 hover:bg-gray-400 text-gray-800 font-medium py-1 px-4 rounded text-sm transition-colors cursor-pointer"
						>Cancel</button
					>
				</div>
			</div>
		</div>
	</div>
{/if}
