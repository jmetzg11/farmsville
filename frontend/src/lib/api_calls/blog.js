export async function postBlog(formData) {
	try {
		const url = `${import.meta.env.VITE_API_URL}/post-blog`;
		const response = await fetch(url, {
			method: 'POST',
			body: formData,
			credentials: 'include'
		});
		return response.ok;
	} catch (error) {
		console.error('Error posting blog', error);
		return false;
	}
}
