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
		class="fixed inset-0 bg-black/70 z-50 flex items-center justify-center p-4"
		onclick={handleClose}
		onkeydown={(e) => e.key === 'Enter' && handleClose()}
	>
		<div class="relative max-w-4xl max-h-[90vh] bg-white rounded-lg shadow-xl overflow-hidden">
			<img
				src={photoUrl}
				alt="Selected item in full size"
				class="max-h-[85vh] max-w-full object-contain"
			/>
			<button
				aria-label="Close photo modal"
				class="absolute top-2 right-2 bg-black/50 hover:bg-black/70 text-white rounded-full p-2 transition-colors cursor-pointer"
				onclick={handleClose}
			>
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
