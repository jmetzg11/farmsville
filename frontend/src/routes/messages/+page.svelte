<script>
	import { getMessages } from '$lib/api_calls/messages.js';
	import { onMount } from 'svelte';
	import { formatDate } from '$lib/helpers.js';
	let messages = $state([]);

	onMount(async () => {
		const results = await getMessages();
		messages = results.messages;
	});
</script>

<div class="max-w-4xl mx-auto p-6">
	<h1 class="text-3xl font-bold text-gray-900 mb-8">Messages</h1>

	<div class="space-y-6">
		{#each messages as message}
			<div class="bg-white rounded-lg shadow-md border border-gray-200 p-6">
				<div class="flex justify-between items-start mb-4">
					<h2 class="text-xl font-semibold text-gray-900">{message.title}</h2>
					<span class="text-sm text-gray-500 font-mono">{formatDate(message.created_at)}</span>
				</div>
				<p class="text-gray-700 leading-relaxed">{message.message}</p>
			</div>
		{/each}
	</div>
</div>
