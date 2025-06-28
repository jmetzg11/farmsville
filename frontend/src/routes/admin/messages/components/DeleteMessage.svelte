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

<div class="flex-vertical">
	{#each messages as message}
		<div class="item-card p-4">
			<div class="flex-between-center">
				<div class="flex-1">
					<h3 class="card-title">{message.title}</h3>
					<p class="card-text-no-center">{formatDate(message.created_at)}</p>
				</div>
				<button class="btn-danger" onclick={() => confirmDelete(message)}> Delete </button>
			</div>
		</div>
	{/each}
</div>

{#if showModal}
	<div class="modal-container">
		<div class="modal-content-item">
			<h2 class="text-title-small mb-4">Confirm Deletion</h2>

			<div class="mb-6">
				<h4 class="card-title mb-2">{messageToDelete?.title}</h4>
				<div class="bg-gray-50 rounded-md p-4 border-l-4 border-gray-300">
					<p class="card-text">{messageToDelete?.message}</p>
				</div>
			</div>

			<p class="text-error mb-4">⚠️ This action cannot be undone</p>

			<div class="flex-buttons">
				<button class="btn-close" onclick={cancelDelete}> Cancel </button>
				<button class="btn-danger" onclick={handleDeleteMessage}> Delete Message </button>
			</div>
		</div>
	</div>
{/if}
