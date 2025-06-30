<script>
	import { onMount } from 'svelte';
	import { postBlog } from '$lib/api_calls/blog';
	import { extractYouTubeId } from '$lib/helpers';
	import EditBlog from './EditBlog.svelte';

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

<EditBlog bind:title bind:content submitBlog={handleSubmit} newBlog={true} />
