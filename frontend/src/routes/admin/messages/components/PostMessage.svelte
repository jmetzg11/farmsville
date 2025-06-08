<script>
	import { postMessage } from '$lib/api_calls/messages.js';

	let title = $state('');
	let message = $state('');

	function sendMessage() {
		postMessage(title, message);
		title = '';
		message = '';
	}
</script>

<div class="mb-4">
	<input
		bind:value={title}
		class="w-full p-3 border rounded-md focus:outline-none focus:ring-2 focus:ring-teal-500"
		placeholder="Enter your title here..."
	/>
</div>

<div class="mb-4">
	<textarea
		bind:value={message}
		class="w-full p-3 border rounded-md focus:outline-none focus:ring-2 focus:ring-teal-500 h-32"
		placeholder="Enter your message here..."
	></textarea>
</div>

<div class="flex justify-between mb-4">
	<button
		onclick={sendMessage}
		disabled={message.length === 0 || title.length === 0}
		class="bg-teal-500 text-white px-6 py-2 rounded-md font-bold hover:bg-teal-600 transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
		>Send</button
	>
</div>
