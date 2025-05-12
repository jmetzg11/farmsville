import { writable } from 'svelte/store';

export const users = writable([]);

export async function refreshUsers() {
	const freshUsers = await getUsers();
	console.log(freshUsers);
	users.set(freshUsers);
}

export async function getUsers() {
	try {
		const url = `${import.meta.env.VITE_API_URL}/users`;
		const response = await fetch(url, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			},
			credentials: 'include'
		});

		if (!response.ok) {
			throw new Error('Failed to get items');
		}

		const data = await response.json();
		console.log(data);
		return data;
	} catch (error) {
		console.error('Error fetching items', error);
		return null;
	}
}

export async function createUser(user) {
	try {
		const url = `${import.meta.env.VITE_API_URL}/users/create`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(user),
			credentials: 'include'
		});
		refreshUsers();
	} catch (error) {
		console.error('Error creating user', error);
		return false;
	}
}
