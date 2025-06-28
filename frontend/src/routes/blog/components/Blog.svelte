<script>
	import { extractYouTubeId } from '$lib/helpers';
	let { blog } = $props();
</script>

<div class="p-6">
	<h2 class="blog-display-title">{blog.title}</h2>

	<div class="space-y-4">
		{#each blog.blocks.toSorted((a, b) => a.order - b.order) as block}
			{#if block.type === 'text'}
				<div class="prose prose-lg max-w-none">
					<p class="blog-display-text">{block.media}</p>
				</div>
			{:else if block.type === 'image'}
				<div class="flex justify-center">
					<img
						src={`${import.meta.env.VITE_PHOTO_URL}${block.media}`}
						alt="Blog Content"
						class="blog-display-image"
					/>
				</div>
			{:else if block.type === 'youtube'}
				<div class="flex justify-center">
					<div class="blog-display-youtube-container">
						<iframe
							src="https://www.youtube.com/embed/{extractYouTubeId(block.media)}"
							title="YouTube video"
							class="blog-display-youtube-iframe"
							frameborder="0"
							allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
							allowfullscreen
						></iframe>
					</div>
				</div>
			{/if}
		{/each}
	</div>
</div>
