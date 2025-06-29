<script>
	import { onMount } from 'svelte';
	import { postBlog } from '$lib/api_calls/blog';
	import { extractYouTubeId } from '$lib/helpers';
	let title = $state('');
	let content = $state([{ type: 'text', media: '', id: crypto.randomUUID() }]);

	async function handleSubmit() {
		const formData = new FormData();
		formData.append('title', title);

		content.forEach((block, index) => {
			formData.append(`content[${index}][type]`, block.type);
			if (block.type === 'image' && block.media instanceof File) {
				formData.append(`content[${index}][file]`, block.media);
			} else {
				formData.append(`content[${index}][media]`, block.media || '');
			}
		});

		const success = await postBlog(formData);
		if (success) {
			title = '';
			content = [{ type: 'text', media: '', id: crypto.randomUUID() }];
		} else {
			throw new Error('Failed to post blog');
		}
	}

	function addBlock(type) {
		content.push({
			type: type,
			media: type === 'text' ? '' : type === 'image' ? null : '',
			id: crypto.randomUUID()
		});
	}

	function moveBlock(index, direction) {
		if (direction === 'up' && index > 0) {
			[content[index], content[index - 1]] = [content[index - 1], content[index]];
		} else if (direction === 'down' && index < content.length - 1) {
			[content[index], content[index + 1]] = [content[index + 1], content[index]];
		}
	}

	function removeBlock(id) {
		content = content.filter((block) => block.id !== id);
		if (content.length === 0) {
			content.push({ type: 'text', media: '', id: crypto.randomUUID() });
		}
	}

	function handleFileUpload(blockId, event) {
		const file = event.target.files[0];
		if (!file) return;

		const block = content.find((b) => b.id === blockId);
		if (block) {
			block.media = file;
			block.preview = URL.createObjectURL(file);
		}
	}
</script>

<div class="mb-4">
	<input bind:value={title} class="blog-title-input" placeholder="Enter your title here..." />
</div>

<div class="space-y-4 mb-6">
	{#each content as block, index (block.id)}
		<div class="blog-create-section-container">
			<div class="blog-create-section-wrapper">
				<span class="blog-create-section-title">
					{block.type === 'text' ? 'Text' : block.type === 'image' ? 'Image' : 'YouTube Video'}
				</span>
				<div class="flex gap-4">
					<button
						onclick={() => moveBlock(index, 'up')}
						class="blog-reorder-button"
						disabled={index === 0}>↑</button
					>
					<button
						onclick={() => moveBlock(index, 'down')}
						class="blog-reorder-button"
						disabled={index === content.length - 1}>↓</button
					>
					<button onclick={() => removeBlock(block.id)} class="blog-remove-section-button">×</button
					>
				</div>
			</div>
			{#if block.type === 'text'}
				<textarea
					bind:value={block.media}
					class="blog-create-textarea"
					placeholder="Enter your text here..."
					style="white-space: pre-wrap;"
				></textarea>
			{:else if block.type === 'image'}
				<div class="space-y-3">
					<input
						type="file"
						accept="image/*"
						onchange={(e) => handleFileUpload(block.id, e)}
						class="blog-create-impage-input"
					/>

					{#if block.preview}
						<div class="mt-2">
							<img src={block.preview} alt="Preview" class="blog-create-image-preview" />
						</div>
					{/if}
				</div>
			{:else if block.type === 'youtube'}
				<div class="space-y-3">
					<input
						bind:value={block.media}
						class="blog-create-youtube-input"
						placeholder="Paste YouTube URL here..."
					/>

					{#if block.media && extractYouTubeId(block.media)}
						<div class="mt-2">
							<iframe
								src="https://www.youtube.com/embed/{extractYouTubeId(block.media)}"
								title="YouTube video player"
								class="blog-create-youtube-preview"
								frameborder="0"
								allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
								allowfullscreen
							></iframe>
						</div>
					{/if}
				</div>
			{/if}
		</div>
	{/each}
</div>

<div class="flex justify-between gap-2">
	<div class="flex gap-2">
		<button onclick={() => addBlock('text')} class="blog-btn-add-section">+ Text</button>
		<button onclick={() => addBlock('image')} class="blog-btn-add-section">+ Image</button>
		<button onclick={() => addBlock('youtube')} class="blog-btn-add-section">+ YouTube</button>
	</div>
	<button onclick={handleSubmit} class="blog-btn-publish"> Publish </button>
</div>
