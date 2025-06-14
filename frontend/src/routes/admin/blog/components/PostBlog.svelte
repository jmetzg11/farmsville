<script>
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
	<input
		bind:value={title}
		class="w-full p-3 border rounded-md focus:outline-none focus:ring-2 focus:ring-teal-500"
		placeholder="Enter your title here..."
	/>
</div>

<div class="space-y-4 mb-6">
	{#each content as block, index (block.id)}
		<div class="border rounded-lg p-4 bg-gray-50">
			<div class="flex justify-between items-center mb-2">
				<span class="text-sm font-medium text-gray-600">
					{block.type === 'text' ? 'Text' : block.type === 'image' ? 'Image' : 'YouTube Video'}
				</span>
				<div class="flex gap-4">
					<button
						onclick={() => moveBlock(index, 'up')}
						class="text-gray-500 hover:text-gray-700 disabled:opacity-30 text-lg cursor-pointer"
						disabled={index === 0}>↑</button
					>
					<button
						onclick={() => moveBlock(index, 'down')}
						class="text-gray-500 hover:text-gray-700 disabled:opacity-30 text-lg cursor-pointer"
						disabled={index === content.length - 1}>↓</button
					>
					<button
						onclick={() => removeBlock(block.id)}
						class="text-red-500 hover:text-red-700 text-lg cursor-pointer">×</button
					>
				</div>
			</div>
			{#if block.type === 'text'}
				<textarea
					bind:value={block.media}
					class="w-full p-3 border rounded-md focus:outline-none focus:ring-2 focus:ring-teal-500 h-24 font-mono"
					placeholder="Enter your text here..."
					style="white-space: pre-wrap;"
				></textarea>
			{:else if block.type === 'image'}
				<div class="space-y-3">
					<input
						type="file"
						accept="image/*"
						onchange={(e) => handleFileUpload(block.id, e)}
						class="w-full p-2 border rounded-md"
					/>

					{#if block.preview}
						<div class="mt-2">
							<img
								src={block.preview}
								alt="Preview"
								class="max-w-xs max-h-48 object-contain border rounded"
							/>
						</div>
					{/if}
				</div>
			{:else if block.type === 'youtube'}
				<div class="space-y-3">
					<input
						bind:value={block.media}
						class="w-full p-3 border rounded-md focus:outline-none focus:ring-2 focus:ring-teal-500"
						placeholder="Paste YouTube URL here..."
					/>

					{#if block.media && extractYouTubeId(block.media)}
						<div class="mt-2">
							<iframe
								src="https://www.youtube.com/embed/{extractYouTubeId(block.media)}"
								title="YouTube video player"
								class="w-full h-64 border rounded"
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
		<button
			onclick={() => addBlock('text')}
			class="bg-sky-500 text-white px-4 py-2 rounded-md hover:bg-sky-600 transition-colors cursor-pointer"
			>+ Text</button
		>
		<button
			onclick={() => addBlock('image')}
			class="bg-sky-500 text-white px-4 py-2 rounded-md hover:bg-sky-600 transition-colors cursor-pointer"
			>+ Image</button
		>
		<button
			onclick={() => addBlock('youtube')}
			class="bg-sky-500 text-white px-4 py-2 rounded-md hover:bg-sky-600 transition-colors cursor-pointer"
			>+ YouTube</button
		>
	</div>
	<button
		onclick={handleSubmit}
		class="bg-emerald-500 text-white px-6 py-2 rounded-md hover:bg-emerald-600 transition-colors cursor-pointer"
	>
		Publish
	</button>
</div>
