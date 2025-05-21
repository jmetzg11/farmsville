<script>
	let { showDeleteModal = $bindable(false), selectedUser } = $props();
	import { deleteUser } from './helpers';

	function closeModal() {
		showDeleteModal = false;
	}

	async function handleDeleteUser() {
		await deleteUser(selectedUser);
		showDeleteModal = false;
	}
</script>

{#if showDeleteModal}
	<div
		class="fixed inset-0 bg-[rgba(0,0,0,0.5)] z-50 flex items-center justify-center h-screen overflow-hidden p-4"
	>
		<div class="bg-white p-8 rounded-lg shadow-lg max-w-md w-full">
			<h2 class="text-2xl font-semibold text-gray-800 mb-6 text-center">
				Are you sure you want to delete
				{selectedUser.name ? selectedUser.name : selectedUser.email}?
			</h2>
			<div class="flex gap-3 justify-between mt-8">
				<button
					onclick={handleDeleteUser}
					class="px-5 py-2.5 rounded-lg text-white font-medium shadow-sm transition-colors duration-200 bg-red-400 hover:bg-red-500"
				>
					Delete User
				</button>
				<button
					onclick={closeModal}
					class="px-5 py-2.5 border border-gray-300 bg-white text-gray-700 rounded-lg hover:bg-gray-50 font-medium transition-colors duration-200 shadow-sm cursor-pointer"
				>
					Cancel
				</button>
			</div>
		</div>
	</div>
{/if}
