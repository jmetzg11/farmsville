import { writable } from 'svelte/store';

export const defaultUser = {
	name: null,
	email: null,
	admin: false,
	isAuthenticated: false
};

export const user = writable(defaultUser);

export function authenticateUser(userData) {
	user.set(userData);
}

export async function initializeUserStore() {
	try {
		const url = import.meta.env.VITE_API_URL + '/auth/me';
		const response = await fetch(url, {
			credentials: 'include'
		});
		if (response.ok) {
			const userData = await response.json();
			console.log('User data:', userData);
			authenticateUser(userData.user);
			return true;
		}
		return false;
	} catch (error) {
		console.log('Error initializing user store', error);
		console.error('Error initializing user store', error);
		return false;
	}
}
