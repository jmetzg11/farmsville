<script>
	import { onMount } from 'svelte';
	import { refreshUsers, users } from './helpers';
	import AddUser from './AddUser.svelte';
	import EditUser from './EditUser.svelte';
	import DeleteUser from './DeleteUser.svelte';

	let showAddModal = $state(false);
	let showEditModal = $state(false);
	let showDeleteModal = $state(false);
	let formDetails = $state({});

	function showAddUser() {
		console.log('I was called ');
		showAddModal = true;
	}

	function showEditForm(user) {
		showEditModal = true;
		formDetails = user;
	}

	function removeUser(user) {
		showDeleteModal = true;
		formDetails = user;
	}

	onMount(async () => {
		await refreshUsers();
	});

	$effect(() => {
		$inspect(users);
	});
</script>

<AddUser bind:showAddModal />
<EditUser bind:showEditModal />
<DeleteUser bind:showDeleteModal />
<button
	class="px-4 py-2 border bg-green-400 rounded-md hover:bg-green-500 cursor-pointer"
	onclick={showAddUser}>Add User</button
>

<table class="min-w-full divide-y divide-gray-200">
	<thead>
		<tr>
			<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
			<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Email</th>
			<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Phone</th>
			<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Role</th>
			<th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Actions</th>
		</tr>
	</thead>
	<tbody class="bg-white divide-y divide-gray-200">
		{#each $users as user}
			<tr class="hover:bg-gray-50">
				<td class="px-6 py-4 whitespace-nowrap">{user.name}</td>
				<td class="px-6 py-4 whitespace-nowrap">{user.email}</td>
				<td class="px-6 py-4 whitespace-nowrap">{user.phone}</td>
				<td class="px-6 py-4 whitespace-nowrap">
					{#if user.admin}
						<span
							class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800"
						>
							Admin
						</span>
					{:else}
						<span
							class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-gray-100 text-gray-800"
						>
							User
						</span>
					{/if}
				</td>
				<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
					<button
						class="text-blue-500 hover:text-blue-700 mr-3 cursor-pointer"
						onclick={() => showEditForm(user)}
					>
						Edit
					</button>
					<button
						class="text-red-500 hover:text-red-700 cursor-pointer"
						onclick={() => removeUser(user)}
					>
						Delete
					</button>
				</td>
			</tr>
		{/each}
	</tbody>
</table>
