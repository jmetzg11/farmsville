<script>
	let { children } = $props();
</script>

<div class="flex flex-col">
	<header class="bg-slate-700 text-white shadow-md w-full">
		<div class="container mx-auto px-4 py-3 flex flex-col sm:flex-row justify-between items-center">
			<nav>
				<ul class="flex space-x-6">
					<li><a href="/admin" class="hover:text-slate-200 transition-colors">Dashboard</a></li>
					<li>
						<a href="/admin/inventory" class="hover:text-slate-200 transition-colors">Inventory</a>
					</li>
					<li><a href="/admin/users" class="hover:text-slate-200 transition-colors">Users</a></li>
					<li>
						<a href="/admin/messages" class="hover:text-slate-200 transition-colors">Messages</a>
					</li>
				</ul>
			</nav>
		</div>
	</header>

	<main class="container mx-auto px-4 py-6">
		{@render children()}
	</main>
</div>
