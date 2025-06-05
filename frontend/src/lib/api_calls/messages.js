export async function sendEmail(emails, title, message) {
	try {
		const url = `${import.meta.env.VITE_API_URL}/send-email`;
		console.log(url);
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ emails, title, message }),
			credentials: 'include'
		});
		return true;
	} catch (error) {
		console.error('Error deleting user', error);
		return false;
	}
}
