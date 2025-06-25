<script>
	let { showPhotoModal = $bindable(false), photoUrl = '' } = $props();

	function handleClose() {
		showPhotoModal = false;
	}
</script>

{#if showPhotoModal}
	<div
		role="button"
		tabindex="0"
		aria-label="Close photo modal background"
		class="modal-container"
		onclick={handleClose}
		onkeydown={(e) => e.key === 'Enter' && handleClose()}
	>
		<div class="modal-content-photo">
			<img
				src={photoUrl}
				alt="Selected item in full size"
				class="max-h-[85vh] max-w-full object-contain"
			/>
			<button aria-label="Close photo modal" class="modal-photo-close-button" onclick={handleClose}>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					class="w-5 h-5"
				>
					<line x1="18" y1="6" x2="6" y2="18"></line>
					<line x1="6" y1="6" x2="18" y2="18"></line>
				</svg>
			</button>
		</div>
	</div>
{/if}
