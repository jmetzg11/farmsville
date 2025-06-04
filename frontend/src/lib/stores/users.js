import { writable } from 'svelte/store';
import { getUsers } from '$lib/api_calls/users.js';

export const users = writable([]);

export async function refreshUsers() {
	const freshUsers = await getUsers();
	users.set(freshUsers);
}
