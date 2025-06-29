<script>
	import { onMount } from 'svelte';
	import { getBlogTitles } from '$lib/api_calls/blog';

	let blogTitles = $state([]);
	$inspect(blogTitles);

	onMount(async () => {
		blogTitles = await getBlogTitles();
	});
</script>

<h1>Edit Blog</h1>
