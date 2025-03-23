<script>
	async function getComments() {
		try {
			const url = `${import.meta.env.VITE_API_URL}/home`;
			const response = await fetch(`${import.meta.env.VITE_API_URL}/get-unapproved-comments`);

			if (!response.ok) {
				throw new Error('Failed to fetch unapproved comments');
			}

			const data = await response.json();
			comments = data;
		} catch (error) {
			console.error('Error fetching unapproved comments:', error);
		}
	}
</script>
