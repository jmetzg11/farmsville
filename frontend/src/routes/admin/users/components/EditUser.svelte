<script>
	import { updateUser } from '$lib/api_calls/users.js';
	let { showEditModal = $bindable(false), selectedUser = $bindable({}) } = $props();

	function closeModal() {
		showEditModal = false;
	}

	async function editUser() {
		showEditModal = false;
		await updateUser(selectedUser);
	}
</script>

{#if showEditModal}
	<div class="modal-container">
		<div class="modal-content-item">
			<h2 class="text-title mb-4">
				Edit {selectedUser.name ? selectedUser.name : selectedUser.email}
			</h2>
			<div class="space-y-5">
				<input
					id="name"
					class="input"
					placeholder="Enter full name"
					bind:value={selectedUser.name}
				/>

				<input
					id="email"
					type="email"
					class="input"
					placeholder="Enter email address"
					bind:value={selectedUser.email}
				/>

				<input
					id="phone"
					type="tel"
					class="input"
					placeholder="Enter phone number"
					bind:value={selectedUser.phone}
				/>
				<div class="flex-buttons">
					<button onclick={editUser} class="btn-secondary"> Edit User </button>
					<button onclick={closeModal} class="btn-close"> Cancel </button>
				</div>
			</div>
		</div>
	</div>
{/if}
