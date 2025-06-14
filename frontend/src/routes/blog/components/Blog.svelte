<script>
	import { extractYouTubeId } from '$lib/helpers';
	let { blog } = $props();
</script>

<div class="p-6">
	<h2 class="text-2xl font-semibold mb-6 text-gray-900">{blog.title}</h2>

	<div class="space-y-4">
		{#each blog.blocks.toSorted((a, b) => a.order - b.order) as block}
			{#if block.type === 'text'}
				<div class="prose prose-lg max-w-none">
					<p class="text-gray-700 leading-relaxed whitespace-pre-wrap">{block.media}</p>
				</div>
			{:else if block.type === 'image'}
				<div class="flex justify-center">
					<img
						src={`${import.meta.env.VITE_PHOTO_URL}${block.media}`}
						alt="Blog Content"
						class="max-w-full max-h-64 object-cover h-auto rounded-lg shadow-sm"
					/>
				</div>
			{:else if block.type === 'youtube'}
				<div class="flex justify-center">
					<div class="w-full max-w-2xl aspect-video">
						<iframe
							src="https://www.youtube.com/embed/{extractYouTubeId(block.media)}"
							title="YouTube video"
							class="w-full h-full rounded-lg"
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
