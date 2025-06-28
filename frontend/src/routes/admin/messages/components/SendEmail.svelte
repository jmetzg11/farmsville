<script>
	import { onMount } from 'svelte';
	import { refreshUsers, users } from '$lib/stores/users.js';
	import { sendEmail } from '$lib/api_calls/messages.js';

	const validUsers = $derived($users.filter((user) => user.email !== null));
	let selectedUsers = $state([]);
	let message = $state('');
	let title = $state('');

	onMount(async () => {
		await refreshUsers();
	});

	async function handleSendEmail() {
		const emails = selectedUsers.map((user) => user.email);
		await sendEmail(emails, title, message);
		title = '';
		message = '';
		selectedUsers = [];
	}

	function selectAll() {
		selectedUsers = [...validUsers];
	}

	function clearAll() {
		selectedUsers = [];
	}

	function toggleUser(user) {
		const index = selectedUsers.findIndex((u) => u.id === user.id);
		if (index === -1) {
			selectedUsers = [...selectedUsers, user];
		} else {
			selectedUsers = selectedUsers.filter((u) => u.id !== user.id);
		}
	}

	function isSelected(user) {
		return selectedUsers.some((u) => u.id === user.id);
	}
</script>

<div class="section-container">
	<h2 class="text-title mb-4">Send Email Message</h2>
	<div class="mb-4">
		<input bind:value={title} class="input" placeholder="Enter your title here..." />
	</div>

	<div class="mb-4">
		<textarea bind:value={message} class="input h-32" placeholder="Enter your message here..."
		></textarea>
	</div>

	<div class="flex justify-between mb-4">
		<button onclick={selectAll} class="btn">Select All</button>
		<button onclick={clearAll} class="btn-secondary">Clear All</button>
		<button
			onclick={handleSendEmail}
			disabled={selectedUsers.length === 0 || !message.trim()}
			class="btn-tertiary">Send</button
		>
	</div>

	<div class="mt-6">
		<div class="text-subtitle mb-4">Recipients ({selectedUsers.length} selected)</div>
		<div class="gird-tight-display">
			{#each validUsers as user}
				<button
					class="p-3 rounded-md border cursor-pointer transition-colors"
					class:bg-teal-100={isSelected(user)}
					class:border-teal-500={isSelected(user)}
					class:border-slate-200={!isSelected(user)}
					onclick={() => toggleUser(user)}
				>
					<div class="font-medium truncate">{user.name ? user.name : user.email}</div>
				</button>
			{/each}
		</div>
	</div>
</div>
