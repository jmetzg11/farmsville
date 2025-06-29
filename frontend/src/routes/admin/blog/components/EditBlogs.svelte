<script>
	import { onMount } from 'svelte';
	import { getBlogTitles, getBlogById, editBlog } from '$lib/api_calls/blog';
	import { formatDate } from '$lib/helpers';
	import EditBlog from './EditBlog.svelte';

	let blogTitles = $state([]);
	let focusedBlog = $state(null);

	onMount(async () => {
		blogTitles = await getBlogTitles();
	});

	async function handleBlogClick(id) {
		focusedBlog = await getBlogById(id);
	}

	async function submitBlog(formData) {
		formData.append('id', focusedBlog.id);
		const success = await editBlog(formData);
		if (success) {
			focusedBlog = null;
			blogTitles = await getBlogTitles();
		}
	}
</script>

{#if focusedBlog}
	<EditBlog
		bind:title={focusedBlog.title}
		bind:content={focusedBlog.blocks}
		{submitBlog}
		newBlog={false}
	/>
{:else}
	<div class="section-container">
		<div class="space-y-6">
			{#each blogTitles as blogTitle}
				<button class="clickable-card w-full" onclick={() => handleBlogClick(blogTitle.id)}>
					<div class="flex justify-between items-start mb-4">
						<h2 class="card-title-larger">{blogTitle.title}</h2>
						<span class="text-sm text-gray-500 font-mono">{formatDate(blogTitle.created_at)}</span>
					</div>
				</button>
			{/each}
		</div>
	</div>
{/if}
