export async function editItem(item) {
	try {
		const url = `${import.meta.env.VITE_API_URL}/items/update`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(item),
			credentials: 'include'
		});
		if (!response.ok) {
			console.error('Error editing item', response.statusText);
			return false;
		}
		return true;
	} catch (error) {
		console.error('Error editing item', error);
		return false;
	}
}

export async function removeItem(itemID) {
	try {
		const url = `${import.meta.env.VITE_API_URL}/items/remove`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ id: itemID }),
			credentials: 'include'
		});
		if (!response.ok) {
			console.error('Error removing item', response.statusText);
			return false;
		}
		return true;
	} catch (error) {
		console.error('Error removing item', error);
		return false;
	}
}
