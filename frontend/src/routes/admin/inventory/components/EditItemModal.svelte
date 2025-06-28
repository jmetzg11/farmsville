<script>
	import { formatDate } from '$lib/helpers.js';
	import { editItem, removeItem } from '$lib/api_calls/admin_items.js';
	import { refreshItems } from '$lib/stores/items';
	let { showEditModal = $bindable(false), selectedItem = null } = $props();

	let fileInput;
	let photoFile = $state(null);

	function triggerFileInput() {
		fileInput.click();
	}

	function handlePhotoChange(event) {
		const file = event.target.files[0];
		if (file) {
			photoFile = file;
		}
	}

	async function handleEdit() {
		await editItem(selectedItem, photoFile);
		await refreshItems();
		showEditModal = false;
	}

	async function handleRemove() {
		await removeItem(selectedItem.id);
		await refreshItems();
		showEditModal = false;
	}

	function handleCancel() {
		showEditModal = false;
	}
</script>

{#if showEditModal}
	<div class="modal-container">
		<div class="modal-content-item">
			<div class="p-4">
				<div class="flex justify-between items-center mb-2">
					<input type="text" bind:value={selectedItem.name} class="input-three-quarter" />
					<span class="card-text">{formatDate(selectedItem.created_at)}</span>
				</div>
				<textarea bind:value={selectedItem.description} class="input" rows="3"></textarea>
				<div class="flex justify-between items-center">
					<div>
						<div class="flex gap-2 items-center">
							<input
								type="number"
								bind:value={selectedItem.remaining_quantity}
								class="input-small"
								min="0"
								max={selectedItem.quantity}
							/>
							/
							<input
								type="number"
								bind:value={selectedItem.quantity}
								class="input-small"
								min={selectedItem.remaining_quantity}
							/>
							<span>remaining</span>
						</div>
					</div>

					<button onclick={triggerFileInput} class="btn-tertiary"
						>{photoFile ? 'Photo Selected' : 'Upload Photo'}</button
					>
					<input
						type="file"
						accept="image/*"
						bind:this={fileInput}
						onchange={handlePhotoChange}
						class="hidden"
					/>
				</div>
				<div class="flex gap-2 justify-between mt-4">
					<button onclick={handleEdit} class="btn"> Save </button>
					<button onclick={handleRemove} class="btn-danger">Remove</button>
					<button onclick={handleCancel} class="btn-close">Cancel</button>
				</div>
			</div>
		</div>
	</div>
{/if}
