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
			return { status: 'error' };
		}
	} catch (error) {
		console.error('Network error sending auth code:', error);
		return { status: 'error' };
	}
}

export async function login(email, password) {
	try {
		const url = `${import.meta.env.VITE_API_URL}/auth/login`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			credentials: 'include',
			body: JSON.stringify({ email: email.toLowerCase(), password })
		});

		if (response.ok) {
			const data = await response.json();
			return { status: 'success', user: data.user };
		} else {
			return { status: 'error' };
		}
	} catch {
		return { status: 'error' };
	}
}

export async function createAccount(name, phone, email, password) {
	try {
		const url = `${import.meta.env.VITE_API_URL}/auth/create`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			credentials: 'include',
			body: JSON.stringify({ name, phone, email, password })
		});

		return await response.json();
	} catch {
		return { status: 'error', message: 'Network error' };
	}
}

export async function sendCodeToResetPassword(email) {
	try {
		const url = `${import.meta.env.VITE_API_URL}/auth/code-to-reset-password`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ email })
		});
		return await response.json();
	} catch {
		return { status: 'error', message: 'Network error' };
	}
}

export async function resetPassword(email, code, password) {
	console.log('resetPassword', email, typeof email, code, typeof code, password, typeof password);
	try {
		const url = `${import.meta.env.VITE_API_URL}/auth/reset-password`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			credentials: 'include',
			body: JSON.stringify({ email, code, password })
		});
		return await response.json();
	} catch {
		return { status: 'error', message: 'Network error' };
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
