<script>
	import { formatDate } from '../root/helpers.js';
	import { editItem, removeItem } from './helpers.js';
	import { refreshItems } from '$lib/stores/items';
	let { showModal = $bindable(false), selectedItem = null } = $props();

	async function handleEdit() {
		await editItem(selectedItem);
		await refreshItems();
		showModal = false;
	}

	async function handleRemove() {
		await removeItem(selectedItem.id);
		await refreshItems();
		showModal = false;
	}

	function handleCancel() {
		showModal = false;
	}
</script>

{#if showModal}
	<div
		class="fixed inset-0 bg-[rgba(0,0,0,0.5)] z-50 flex items-center justify-center h-screen overflow-hidden p-4"
	>
		<div
			class="bg-white rounded-lg shadow-md overflow-hidden border border-gray-200 hover:shadow-lg transition-shadow"
		>
			<div class="p-4">
				<div class="flex justify-between items-center mb-2">
					<input
						type="text"
						bind:value={selectedItem.name}
						class="text-lg font-bold text-gray-800 border border-gray-300 rounded px-2 py-1 w-3/4"
					/>
					<span class="text-xs text-gray-500">{formatDate(selectedItem.created_at)}</span>
				</div>
				<textarea
					bind:value={selectedItem.description}
					class="text-gray-600 mb-3 text-sm border border-gray-300 rounded px-2 py-1 w-full"
					rows="3"
				></textarea>
				<div class="flex justify-between items-center">
					<div>
						<div class="flex gap-2 items-center">
							<input
								type="number"
								bind:value={selectedItem.remaining_quantity}
								class="font-medium w-16 border border-gray-300 rounded px-2 py-1"
								min="0"
								max={selectedItem.quantity}
							/>
							/
							<input
								type="number"
								bind:value={selectedItem.quantity}
								class="w-16 border border-gray-300 rounded px-2 py-1"
								min={selectedItem.remaining_quantity}
							/>
							<span>remaining</span>
						</div>
					</div>
				</div>
				<div class="flex gap-2 justify-between mt-4">
					<button
						onclick={handleEdit}
						class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-1 px-4 rounded text-sm transition-colors cursor-pointer"
					>
						Edit
					</button>
					<button
						onclick={handleRemove}
						class="bg-red-600 hover:bg-red-700 text-white font-medium py-1 px-4 rounded text-sm transition-colors cursor-pointer"
						>Remove</button
					>
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
