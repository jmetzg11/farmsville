<script>
	import PostMessage from './PostMessage.svelte';
	import DeleteMessage from './DeleteMessage.svelte';

	let activeTab = $state('post');
</script>

<div class="section-container">
	<div class="flex-between-center mb-4">
		<button
			class="text-2xl font-bold cursor-pointer {activeTab === 'post'
				? 'text-teal-600 border-b-2 border-teal-600'
				: 'text-slate-700 hover:text-teal-500'}"
			onclick={() => (activeTab = 'post')}
		>
			Post Message
		</button>
		<button
			class="text-2xl font-bold cursor-pointer {activeTab === 'delete'
				? 'text-teal-600 border-b-2 border-teal-600'
				: 'text-slate-700 hover:text-teal-500'}"
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
