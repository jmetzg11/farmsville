<script>
	import { logout } from './helpers';
	let { onCancel } = $props();

	async function handleLogout() {
		await logout();
		onCancel();
	}
</script>

<p class="text-gray-600 mb-4">Are you sure you want to log out?</p>
<div class="flex gap-2 justify-between mt-4">
	<button
		onclick={handleLogout}
		class="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-100 cursor-pointer"
		>Logout</button
	>
	<button
		onclick={onCancel}
		class="py-2 px-4 rounded-md text-white
                transition-colors duration-200
                bg-blue-600 hover:bg-blue-700 cursor-pointer"
	>
		cancel
	</button>
</div>
