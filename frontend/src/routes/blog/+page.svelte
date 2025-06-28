<script>
	import { onMount } from 'svelte';
	import Blog from './components/Blog.svelte';
	let blogs = $state([]);

	async function getBlogs() {
		const url = `${import.meta.env.VITE_API_URL}/blogs`;
		const response = await fetch(url);
		const data = await response.json();
		blogs = data.blogs;
	}

	function extractYouTubeId(url) {
		const regex = /(?:youtube\.com\/watch\?v=|youtu\.be\/)([^&\n?#]+)/;
		const match = url.match(regex);
		return match ? match[1] : null;
	}

	onMount(() => {
		getBlogs();
	});
</script>

<div class="max-w-4xl mx-auto p-6">
	{#each blogs as blog}
		<article class="blog-display-container">
			<Blog {blog} />
		</article>
	{/each}
</div>
