<script>
	import { logout } from '$lib/api_calls/auth';
	let { onCancel } = $props();

	async function handleLogout() {
		await logout();
		onCancel();
	}
</script>

<p class="text-title mb-8">Are you sure you want to log out?</p>
<div class="flex-buttons">
	<button onclick={handleLogout} class="btn">Logout</button>
	<button onclick={onCancel} class="btn-close"> Cancel </button>
</div>
