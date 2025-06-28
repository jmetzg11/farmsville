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
	<input bind:value={title} class="input" placeholder="Enter your title here..." />
</div>

<div class="mb-4">
	<textarea bind:value={message} class="input h-32" placeholder="Enter your message here..."
	></textarea>
</div>

<div class="flex justify-between mb-4">
	<button
		onclick={sendMessage}
		disabled={message.length === 0 || title.length === 0}
		class="btn-secondary">Send</button
	>
</div>
