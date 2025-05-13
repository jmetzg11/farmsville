<script>
	import { onMount } from 'svelte';
	import { refreshUsers, users } from '../users/helpers';
	import { sendTextMessage } from './helpers';

	const validUsers = $derived($users.filter((user) => user.phone.length >= 9));
	let selectedUsers = $state([]);
	let message = $state('');

	$inspect($users);

	onMount(async () => {
		await refreshUsers();
	});

	async function sendMessage() {
		const phoneNumbers = selectedUsers.map((user) => user.phone);
		await sendTextMessage(phoneNumbers, message);
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

<div class="max-w-4xl mx-auto p-4">
	<h2 class="text-2xl font-bold text-slate-700 mb-4">Text Messages</h2>

	<div class="mb-4">
		<textarea
			bind:value={message}
			class="w-full p-3 border rounded-md focus:outline-none focus:ring-2 focus:ring-teal-500 h-32"
			placeholder="Enter your message here..."
		></textarea>
	</div>

	<div class="flex justify-between mb-4">
		<button
			onclick={selectAll}
			class="bg-teal-500 text-white px-6 py-2 rounded-md font-bold hover:bg-teal-600 transition-colors cursor-pointer"
			>Select All</button
		>
		<button
			onclick={clearAll}
			class="bg-teal-500 text-white px-6 py-2 rounded-md font-bold hover:bg-teal-600 transition-colors cursor-pointer"
			>Clear All</button
		>
		<button
			onclick={sendMessage}
			disabled={selectedUsers.length === 0 || !message.trim()}
			class="bg-teal-500 text-white px-6 py-2 rounded-md font-bold hover:bg-teal-600 transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
			>Send</button
		>
	</div>

	<div class="mt-6">
		<div class="text-slate-700 font-medium mb-2">Recipients ({selectedUsers.length} selected)</div>
		<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2">
			{#each validUsers as user}
				<button
					class="p-3 rounded-md border cursor-pointer transition-colors"
					class:bg-teal-100={isSelected(user)}
					class:border-teal-500={isSelected(user)}
					class:border-slate-200={!isSelected(user)}
					onclick={() => toggleUser(user)}
				>
					<div class="font-medium truncate">{user.name ? user.name : user.email}</div>
					<div class="text-sm text-slate-500 truncate">{user.phone}</div>
				</button>
			{/each}
		</div>
	</div>
</div>
