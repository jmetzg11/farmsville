<script>
	import { onMount } from 'svelte';
	import { getBlogs } from '$lib/api_calls/blog';
	import Blog from './components/Blog.svelte';
	let blogs = $state([]);

	onMount(async () => {
		blogs = await getBlogs();
	});
</script>

<div class="blogs-container">
	{#each blogs as blog}
		<article class="blog-display-container">
			<Blog {blog} />
		</article>
	{/each}
</div>
