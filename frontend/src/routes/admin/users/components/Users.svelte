<script>
	import { onMount } from 'svelte';
	import { refreshUsers, users } from '$lib/stores/users';
	import AddUser from './AddUser.svelte';
	import EditUser from './EditUser.svelte';
	import DeleteUser from './DeleteUser.svelte';

	let showAddModal = $state(false);
	let showEditModal = $state(false);
	let showDeleteModal = $state(false);
	let selectedUser = $state({});

	function showAddUser() {
		showAddModal = true;
	}

	function showEditForm(user) {
		showEditModal = true;
		selectedUser = user;
	}

	function removeUser(user) {
		showDeleteModal = true;
		selectedUser = user;
	}

	onMount(async () => {
		await refreshUsers();
	});
</script>

<AddUser bind:showAddModal />
<EditUser bind:showEditModal bind:selectedUser />
<DeleteUser bind:showDeleteModal {selectedUser} />

<button class="btn-secondary mb-4" onclick={showAddUser}> Add User </button>

<table class="table">
	<thead>
		<tr>
			<th class="th">Name</th>
			<th class="th">Email</th>
			<th class="th">Phone</th>
			<th class="th">Role</th>
			<th class="th-last">Actions</th>
		</tr>
	</thead>
	<tbody class="tbody">
		{#each $users as user}
			<tr class="tr-body">
				<td class="td">{user.name}</td>
				<td class="td">{user.email}</td>
				<td class="td">{user.phone}</td>
				<td class="td">
					{#if user.admin}
						<span class="td-span-admin"> Admin </span>
					{:else}
						<span class="td-span-user"> User </span>
					{/if}
				</td>
				<td class="td-last">
					<button class="td-button-edit" onclick={() => showEditForm(user)}> Edit </button>
					<button class="td-button-delete" onclick={() => removeUser(user)}> Delete </button>
				</td>
			</tr>
		{/each}
	</tbody>
</table>
