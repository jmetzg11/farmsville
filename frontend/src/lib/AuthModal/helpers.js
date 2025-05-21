import { authenticateUser } from '$lib/stores/auth';
import { defaultUser } from '$lib/stores/auth';

export async function emailAuth(email) {
	try {
		const url = `${import.meta.env.VITE_API_URL}/auth`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ email: email.toLowerCase() })
		});

		if (response.ok) {
			return 'enter-code';
		} else {
			const errorData = await response.json();
			console.error('Server returned an error:', response.status, errorData);
			return 'error';
		}
	} catch (error) {
		console.error('Network error sending auth email:', error);
		return 'error';
	}
}

export async function authVerify(email, code) {
	try {
		const url = `${import.meta.env.VITE_API_URL}/auth/verify`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			credentials: 'include',
			body: JSON.stringify({ email: email.toLowerCase(), code })
		});

		if (response.ok) {
			const data = await response.json();
			return { status: 'success', user: data.user };
		} else {
			return { status: 'error', message: 'Invalid code' };
		}
	} catch (error) {
		console.error('Network error sending auth code:', error);
		return { status: 'error', message: 'Network error sending auth code' };
	}
}

export async function logout() {
	try {
		const url = `${import.meta.env.VITE_API_URL}/auth/logout`;
		await fetch(url, {
			credentials: 'include'
		});
		authenticateUser(defaultUser);
	} catch (error) {
		console.error('Network error logging out:', error);
	}
}

export function preventNonNumericInput(e) {
	if (!/[0-9]/.test(e.key)) {
		e.preventDefault();
	}
}
