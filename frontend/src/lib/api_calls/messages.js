export async function sendTextMessage(numbers, message) {
	try {
		const url = `${import.meta.env.VITE_API_URL}/messages`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ numbers, message }),
			credentials: 'include'
		});
		return true;
	} catch (error) {
		console.error('Error deleting user', error);
		return false;
	}
}
