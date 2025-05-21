import { refreshUsers } from '$lib/stores/users';

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
		return data;
	} catch (error) {
		console.error('Error fetching items', error);
		return null;
	}
}

export async function deleteUser(user) {
	try {
		const url = `${import.meta.env.VITE_API_URL}/users/remove`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(user.id),
			credentials: 'include'
		});
		await refreshUsers();
		return true;
	} catch (error) {
		console.error('Error deleting user', error);
		return false;
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
		await refreshUsers();
		return true;
	} catch (error) {
		console.error('Error creating user', error);
		return false;
	}
}

export async function updateUser(user) {
	try {
		const url = `${import.meta.env.VITE_API_URL}/users/update`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(user),
			credentials: 'include'
		});
		await refreshUsers();
		return true;
	} catch (error) {
		console.error('Error creating user', error);
		return false;
	}
}
