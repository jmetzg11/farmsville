<script>
	import { createUser } from '$lib/api_calls/users.js';
	let { showAddModal = $bindable(false) } = $props();
	let user = $state({
		name: '',
		email: '',
		phone: ''
	});

	let isValidUser = $derived(
		(user.email.includes('@') && user.email.length >= 5) || user.name.length > 2
	);

	function closeModal() {
		showAddModal = false;
	}

	async function addUser() {
		showAddModal = false;
		await createUser(user);
	}
</script>

{#if showAddModal}
	<div class="modal-container">
		<div class="modal-content-item">
			<h2 class="text-title mb-4">Add New User</h2>
			<div class="space-y-5">
				<input id="name" class="input" placeholder="Enter full name" bind:value={user.name} />
				<input
					id="email"
					type="email"
					class="input"
					placeholder="Enter email address"
					bind:value={user.email}
				/>
				<input
					id="phone"
					type="tel"
					class="input"
					placeholder="Enter phone number"
					bind:value={user.phone}
				/>
			</div>
			<div class="flex gap-3 justify-between mt-8">
				<button onclick={addUser} disabled={!isValidUser} class="btn"> Add User </button>
				<button onclick={closeModal} class="btn-close"> Cancel </button>
			</div>
		</div>
	</div>
{/if}
