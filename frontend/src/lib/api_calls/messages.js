export async function postMessage(title, message) {
	try {
		const url = `${import.meta.env.VITE_API_URL}/post-message`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ title, message }),
			credentials: 'include'
		});
		return true;
	} catch (error) {
		console.error('Error posting message', error);
		return false;
	}
}

export async function getMessages() {
	try {
		const url = `${import.meta.env.VITE_API_URL}/messages`;
		const response = await fetch(url, {
			method: 'GET'
		});

		return await response.json();
	} catch (error) {
		console.error('Error getting messages', error);
		return [];
	}
}

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
