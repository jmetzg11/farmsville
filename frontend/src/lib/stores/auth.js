import { writable } from 'svelte/store';

export const user = writable({
	name: null,
	email: null,
	admin: false,
	isAuthenticated: false
});

export function setUser(userData) {
	user.set({
		...userData,
		isAuthenticated: true
	});
}

export async function initializeUserStore() {
	try {
		const url = import.meta.env.VITE_API_URL + '/auth/me';
		const response = await fetch(url, {
			credentials: 'include'
		});
		if (response.ok) {
			const userData = await response.json();
			setUser(userData.user);
			return true;
		}
		return false;
	} catch (error) {
		console.error('Error initializing user store', error);
		return false;
	}
}
