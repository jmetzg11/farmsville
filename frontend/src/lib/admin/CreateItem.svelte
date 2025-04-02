<script>
	import { refreshItems } from '$lib/stores/items';
	import { createItem } from './helpers.js';

	const DEFAULT_ITEM = {
		title: '',
		description: '',
		quantity: 0
	};

	let newItem = { ...DEFAULT_ITEM };

	async function handleCreate() {
		await createItem(newItem);
		await refreshItems();
		newItem = { ...DEFAULT_ITEM };
	}

	$: isFormValid =
		newItem.title.trim() !== '' && newItem.description.trim() !== '' && newItem.quantity > 0;
</script>

<div class="container mx-auto px-4 py-6">
	<h2 class="text-2xl font-bold text-gray-800 mb-6 border-b pb-2">Create Item</h2>
	<div class="gap-6 p-4">
		<div
			class="bg-white rounded-lg shadow-md overflow-hidden border border-gray-200 hover:shadow-lg transition-shadow"
		>
			<div class="p-4">
				<!-- Top row for title -->
				<div class="mb-4">
					<input
						type="text"
						placeholder="Title"
						bind:value={newItem.title}
						class="text-lg font-bold text-gray-800 border border-gray-300 rounded px-2 py-1 w-full"
					/>
				</div>

				<!-- Middle section for description -->
				<div class="mb-4">
					<textarea
						placeholder="Description"
						bind:value={newItem.description}
						class="text-gray-600 text-sm border border-gray-300 rounded px-2 py-1 w-full"
						rows="3"
					></textarea>
				</div>

				<!-- Bottom row for quantity and create button -->
				<div class="flex items-center justify-between">
					<div class="flex items-center">
						<label for="quantity" class="text-sm text-gray-700 mr-2">Quantity:</label>
						<input
							id="quantity"
							type="number"
							bind:value={newItem.quantity}
							class="font-medium w-16 border border-gray-300 rounded px-2 py-1"
							min="0"
						/>
					</div>
					<button
						on:click={handleCreate}
						disabled={!isFormValid}
						class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-1 px-4 rounded text-sm transition-colors cursor-pointer
						disabled:bg-blue-300 disabled:cursor-not-allowed"
					>
						Create
					</button>
				</div>
			</div>
		</div>
	</div>
</div>
