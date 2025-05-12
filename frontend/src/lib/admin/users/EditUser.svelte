<script>
	import { updateUser } from './helpers';
	let { showEditModal = $bindable(false), formDetails = $bindable({}) } = $props();

	function closeModal() {
		showEditModal = false;
	}

	async function editUser() {
		showEditModal = false;
		await updateUser(formDetails);
	}
</script>

{#if showEditModal}
	<div
		class="fixed inset-0 bg-[rgba(0,0,0,0.5)] z-50 flex items-center justify-center h-screen overflow-hidden p-4"
	>
		<div class="bg-white p-8 rounded-lg shadow-lg max-w-md w-full">
			<h2 class="text-2xl font-semibold text-gray-800 mb-6 text-center">
				Edit {formDetails.name ? formDetails.name : formDetails.email}
			</h2>

			<div class="space-y-5">
				<div class="relative">
					<label for="name" class="block text-sm font-medium text-gray-700 mb-1 ml-1">Name</label>
					<div class="relative">
						<span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-500">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-5 w-5"
								viewBox="0 0 20 20"
								fill="currentColor"
							>
								<path
									fill-rule="evenodd"
									d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
									clip-rule="evenodd"
								/>
							</svg>
						</span>
						<input
							id="name"
							class="w-full pl-10 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none text-gray-700"
							placeholder="Enter full name"
							bind:value={formDetails.name}
						/>
					</div>
				</div>

				<div class="relative">
					<label for="email" class="block text-sm font-medium text-gray-700 mb-1 ml-1"
						>Email Address</label
					>
					<div class="relative">
						<span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-500">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-5 w-5"
								viewBox="0 0 20 20"
								fill="currentColor"
							>
								<path d="M2.003 5.884L10 9.882l7.997-3.998A2 2 0 0016 4H4a2 2 0 00-1.997 1.884z" />
								<path d="M18 8.118l-8 4-8-4V14a2 2 0 002 2h12a2 2 0 002-2V8.118z" />
							</svg>
						</span>
						<input
							id="email"
							type="email"
							class="w-full pl-10 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none text-gray-700"
							placeholder="Enter email address"
							bind:value={formDetails.email}
						/>
					</div>
				</div>

				<div class="relative">
					<label for="phone" class="block text-sm font-medium text-gray-700 mb-1 ml-1"
						>Phone Number</label
					>
					<div class="relative">
						<span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-500">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-5 w-5"
								viewBox="0 0 20 20"
								fill="currentColor"
							>
								<path
									d="M2 3a1 1 0 011-1h2.153a1 1 0 01.986.836l.74 4.435a1 1 0 01-.54 1.06l-1.548.773a11.037 11.037 0 006.105 6.105l.774-1.548a1 1 0 011.059-.54l4.435.74a1 1 0 01.836.986V17a1 1 0 01-1 1h-2C7.82 18 2 12.18 2 5V3z"
								/>
							</svg>
						</span>
						<input
							id="phone"
							type="tel"
							class="w-full pl-10 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none text-gray-700"
							placeholder="Enter phone number"
							bind:value={formDetails.phone}
						/>
					</div>
				</div>

				<div class="flex gap-3 justify-end mt-8">
					<button
						onclick={closeModal}
						class="px-5 py-2.5 border border-gray-300 bg-white text-gray-700 rounded-lg hover:bg-gray-50 font-medium transition-colors duration-200 shadow-sm cursor-pointer"
					>
						Cancel
					</button>
					<button
						onclick={editUser}
						class="px-5 py-2.5 rounded-lg text-white font-medium shadow-sm transition-colors duration-200
						bg-amber-600 hover:bg-amber-700 cursor-pointer"
					>
						Edit User
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}
