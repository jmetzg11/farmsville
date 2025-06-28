<script>
	import { deleteUser } from '$lib/api_calls/users.js';
	let { showDeleteModal = $bindable(false), selectedUser } = $props();

	function closeModal() {
		showDeleteModal = false;
	}

	async function handleDeleteUser() {
		await deleteUser(selectedUser);
		showDeleteModal = false;
	}
</script>

{#if showDeleteModal}
	<div class="modal-container">
		<div class="modal-content">
			<h2 class="text-title">
				Are you sure you want to delete
				{selectedUser.name ? selectedUser.name : selectedUser.email}?
			</h2>
			<div class="flex-buttons mt-8">
				<button onclick={handleDeleteUser} class="btn-danger"> Delete User </button>
				<button onclick={closeModal} class="btn-close"> Cancel </button>
			</div>
		</div>
	</div>
{/if}
