export async function getBlogs() {
	const url = `${import.meta.env.VITE_API_URL}/blogs`;
	const response = await fetch(url);
	const data = await response.json();
	return data.blogs;
}

export async function getBlogTitles() {
	const url = `${import.meta.env.VITE_API_URL}/get-blog-titles`;
	const response = await fetch(url, {
		credentials: 'include'
	});
	const data = await response.json();
	return data.titles;
}

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
