import { writable } from 'svelte/store';

export const items = writable([]);
export const claimedItems = writable([]);

export async function refreshItems() {
	try {
		const url = `${import.meta.env.VITE_API_URL}/items`;
		const response = await fetch(url);

		if (!response.ok) {
			throw new Error('Failed to get items');
		}

		const data = await response.json();
		items.set(data.items);
		claimedItems.set(data.claimedItems);
	} catch (error) {
		console.error('Error fetching items', error);
		return null;
	}
}
