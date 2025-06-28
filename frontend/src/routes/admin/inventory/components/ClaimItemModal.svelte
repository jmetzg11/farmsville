<script>
	import { onMount } from 'svelte';
	import { formatDate } from '$lib/helpers.js';
	import { claimItem } from '$lib/api_calls/admin_items.js';
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
	<div class="modal-container">
		<div class="modal-content">
			<div class="p-4">
				<div class="flex-between-center mb-4">
					<h3 class="text-lg font-bold text-gray-800">{selectedItem.name}</h3>
					<div class="w-1/3 text-sm text-gray-600 mt-1">
						<span class="font-medium">{selectedItem.remaining_quantity}</span> / {selectedItem.quantity}
						remaining
					</div>
				</div>
				<div class="mb-4">
					<div class="flex-between-center mb-2">
						<div></div>
						<label for="amount-input" class="text-sm font-medium text-gray-700">Amount</label>
					</div>
					<div class="flex-between-center mb-4">
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
								class="input border-gray-300"
								placeholder="Amount"
							/>
						</div>
					</div>
				</div>
				<div class="flex gap-2 justify-between mt-8">
					<button onclick={handleClaim} disabled={!selectedUser || !amount} class="btn">
						Claim
					</button>
					<button onclick={cleanUp} class="btn-close">Cancel</button>
				</div>
			</div>
		</div>
	</div>
{/if}
