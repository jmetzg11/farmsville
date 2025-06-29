<script>
	import PostBlog from './PostBlog.svelte';
	import EditBlogs from './EditBlogs.svelte';

	let activeTab = $state('post');
</script>

<div class="section-container">
	<div class="flex-between-center mb-4">
		<button
			class="tab-button {activeTab === 'post' ? 'tab-button-active' : 'tab-button-inactive'}"
			onclick={() => (activeTab = 'post')}
		>
			Post Blog
		</button>
		<button
			class="tab-button {activeTab === 'edit' ? 'tab-button-active' : 'tab-button-inactive'}"
			onclick={() => (activeTab = 'edit')}
		>
			Edit Blog
		</button>
	</div>

	{#if activeTab === 'post'}
		<PostBlog />
	{:else}
		<EditBlogs />
	{/if}
</div>
