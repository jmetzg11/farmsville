<script>
	import PostMessage from './PostMessage.svelte';
	import DeleteMessage from './DeleteMessage.svelte';

	let activeTab = $state('post');
</script>

<div class="section-container">
	<div class="flex-between-center mb-4">
		<button
			class="tab-button {activeTab === 'post' ? 'tab-button-active' : 'tab-button-inactive'}"
			onclick={() => (activeTab = 'post')}
		>
			Post Message
		</button>
		<button
			class="tab-button {activeTab === 'delete' ? 'tab-button-active' : 'tab-button-inactive'}"
			onclick={() => (activeTab = 'delete')}
		>
			Delete Message
		</button>
	</div>

	{#if activeTab === 'post'}
		<PostMessage />
	{:else}
		<DeleteMessage />
	{/if}
</div>
