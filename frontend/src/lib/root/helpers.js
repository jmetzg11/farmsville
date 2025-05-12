export async function makeClaim(itemId, quantity) {
	try {
		const url = `${import.meta.env.VITE_API_URL}/items/claim`;
		console.log(url);
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				itemId,
				quantity
			}),
			credentials: 'include'
		});
		if (!response.ok) {
			console.error('Error making claim', response.statusText);
			return false;
		}
		return true;
	} catch (error) {
		console.error('Error making claim', error);
		return false;
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
