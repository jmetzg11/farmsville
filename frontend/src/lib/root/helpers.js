export async function getComments() {
	try {
		const url = `${import.meta.env.VITE_API_URL}/items`;
		const response = await fetch(url);

		if (!response.ok) {
			throw new Error('Failed to get items');
		}

		const data = await response.json();
		return {
			items: data.items,
			claimedItems: data.claimedItems
		};
	} catch (error) {
		console.error('Error fetching items', error);
		return null;
	}
}

export function formatDate(dateString) {
	const date = new Date(dateString);
	return date.toLocaleDateString('en-US', {
		year: 'numeric',
		month: 'short',
		day: 'numeric'
	});
}
