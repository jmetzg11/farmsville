<script>
	import { onMount } from 'svelte';
	import { getMessages, deleteMessage } from '$lib/api_calls/messages.js';
	import { formatDate } from '$lib/helpers.js';

	let messages = $state([]);
	let showModal = $state(false);
	let messageToDelete = $state(null);

	onMount(async () => {
		const results = await getMessages();
		messages = results.messages;
	});

	function confirmDelete(message) {
		messageToDelete = message;
		showModal = true;
	}

	async function handleDeleteMessage() {
		const result = await deleteMessage(messageToDelete.id);
		if (result.success) {
			messages = messages.filter((message) => message.id !== messageToDelete.id);
		}
		showModal = false;
		messageToDelete = null;
	}

	function cancelDelete() {
		showModal = false;
		messageToDelete = null;
	}
</script>

<div class="space-y-3">
	{#each messages as message}
		<div class="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow">
			<div class="flex justify-between items-center">
				<div class="flex-1">
					<h3 class="text-lg font-semibold text-gray-900 mb-1">{message.title}</h3>
					<p class="text-sm text-gray-600">{formatDate(message.created_at)}</p>
				</div>
				<button
					class="ml-4 px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition-colors font-medium cursor-pointer"
					onclick={() => confirmDelete(message)}
				>
					Delete
				</button>
			</div>
		</div>
	{/each}
</div>

{#if showModal}
	<div class="modal-container">
		<div class="modal-content-item">
			<div class="mb-4">
				<h3 class="text-xl font-bold text-gray-900 mb-2">Confirm Deletion</h3>
				<div class="h-px bg-gray-200"></div>
			</div>

			<div class="mb-6">
				<h4 class="font-semibold text-gray-800 mb-2">{messageToDelete?.title}</h4>
				<div class="bg-gray-50 rounded-md p-4 border-l-4 border-gray-300">
					<p class="text-gray-700 text-sm leading-relaxed">{messageToDelete?.message}</p>
				</div>
			</div>

			<p class="text-red-600 text-sm mb-6 font-medium">⚠️ This action cannot be undone</p>

			<div class="flex justify-end space-x-3">
				<button
					class="px-5 py-2 text-gray-700 bg-gray-200 rounded-md hover:bg-gray-300 transition-colors cursor-pointer font-medium"
					onclick={cancelDelete}
				>
					Cancel
				</button>
				<button
					class="px-5 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition-colors cursor-pointer font-medium shadow-sm"
					onclick={handleDeleteMessage}
				>
					Delete Message
				</button>
			</div>
		</div>
	</div>
{/if}
