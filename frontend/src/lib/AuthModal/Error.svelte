<script>
	let { status = $bindable('error'), onClose, message, previousStatus } = $props();

	function onTryAgain() {
		status = previousStatus;
	}
</script>

<p class="text-error mb-4">{message}.</p>
<div class="flex-buttons">
	<button onclick={onTryAgain} class="btn"> Try Again </button>
	<button onclick={onClose} class="btn-close"> Cancel </button>
</div>
