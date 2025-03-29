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
		isAuthentiated: true
	});
}
