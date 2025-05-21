<script>
	import { onMount } from 'svelte';
	import { formatDate } from '$lib/root/helpers.js';
	import { claimItem } from './helpers.js';
	import { refreshItems } from '$lib/stores/items';
	import { users, refreshUsers } from '$lib/stores/users';
	let { showClaimModal = $bindable(false), selectedItem = null } = $props();

	let amount = $state(0);
	let selectedUser = $state(null);

	onMount(async () => {
		await refreshUsers();
	});

	async function handleClaim() {
		await claimItem(selectedUser, selectedItem.id, amount);
		await refreshItems();
		cleanUp();
	}

	function cleanUp() {
		amount = 0;
		selectedUser = null;
		showClaimModal = false;
	}
</script>

{#if showClaimModal}
	<div
		class="fixed inset-0 bg-[rgba(0,0,0,0.5)] z-50 flex items-center justify-center h-screen overflow-hidden p-6"
	>
		<div
			class="bg-white rounded-lg shadow-md overflow-hidden border border-gray-200 hover:shadow-lg transition-shadow max-w-md w-full gap-8"
		>
			<div class="p-4">
				<div class="flex justify-between items-center mb-4">
					<h3 class="text-lg font-bold text-gray-800">{selectedItem.name}</h3>
					<div class="w-1/3 text-sm text-gray-600 mt-1">
						<span class="font-medium">{selectedItem.remaining_quantity}</span> / {selectedItem.quantity}
						remaining
					</div>
					<span class="text-xs text-gray-500">{formatDate(selectedItem.created_at)}</span>
				</div>
				<div class="mb-4">
					<div class="flex justify-between items-center mb-2">
						<div></div>
						<label for="amount-input" class="text-sm font-medium text-gray-700">Amount</label>
					</div>
					<div class="flex justify-between items-center mb-4">
						<div class="w-2/3 pr-2">
							<select bind:value={selectedUser} class="w-full p-2 border-gray-300 rounded-md">
								<option value={null} disabled>Select a User</option>
								{#each $users as user}
									<option value={user.id}>{user.name ? user.name : user.email}</option>
								{/each}
							</select>
						</div>
						<div class="w-1/3 pl-2">
							<input
								type="number"
								bind:value={amount}
								min="1"
								max={selectedItem.remaining_quantity}
								class="w-full p-2 border border-gray-300 rounded-md"
								placeholder="Amount"
							/>
						</div>
					</div>
				</div>
				<div class="flex gap-2 justify-between mt-8">
					<button
						onclick={handleClaim}
						disabled={!selectedUser || !amount}
						class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-1 px-4 rounded text-sm transition-colors cursor-pointer disabled:opacity-50 disabled:bg-gray-400 disabled:cursor-not-allowed disabled:hover:bg-gray-400"
					>
						Claim
					</button>
					<button
						onclick={cleanUp}
						class="bg-gray-300 hover:bg-gray-400 text-gray-800 font-medium py-1 px-4 rounded text-sm transition-colors cursor-pointer"
						>Cancel</button
					>
				</div>
			</div>
		</div>
	</div>
{/if}
