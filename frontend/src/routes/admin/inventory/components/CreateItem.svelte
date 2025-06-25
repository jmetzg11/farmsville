<script>
	import { refreshItems } from '$lib/stores/items';
	import { createItem } from '$lib/api_calls/admin_items.js';

	const DEFAULT_ITEM = {
		title: '',
		description: '',
		quantity: 0,
		photo: null
	};

	let newItem = { ...DEFAULT_ITEM };
	let fileInput;

	function triggerFileInput() {
		fileInput.click();
	}

	function handlePhotoChange(event) {
		const file = event.target.files[0];
		if (file) {
			newItem.photo = file;
		}
	}

	async function handleCreate() {
		const formData = new FormData();
		formData.append('title', newItem.title);
		formData.append('description', newItem.description);
		formData.append('quantity', newItem.quantity.toString());

		if (newItem.photo) {
			formData.append('photo', newItem.photo);
		}

		await createItem(formData);
		await refreshItems();
		newItem = { ...DEFAULT_ITEM };
	}

	$: isFormValid =
		newItem.title.trim() !== '' && newItem.description.trim() !== '' && newItem.quantity > 0;
</script>

<div class="section-container">
	<h2 class="text-2xl font-bold text-gray-800 mb-6 border-b pb-2">Create Item</h2>
	<div class="gap-6 p-4">
		<div class="item-card">
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

				<!-- Bottom row for quantity, photo, and create button -->
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
					<div class="flex items-center">
						<input
							type="file"
							accept="image/*"
							bind:this={fileInput}
							onchange={handlePhotoChange}
							class="hidden"
						/>
						<button
							onclick={triggerFileInput}
							class="bg-slate-700 hover:bg-slate-800 text-white font-medium py-1 px-4 rounded text-sm transition-colors cursor-pointer"
							>{newItem.photo ? 'Photo Selected' : 'Upload Photo'}</button
						>
					</div>
					<button
						onclick={handleCreate}
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
