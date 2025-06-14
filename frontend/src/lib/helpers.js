export function formatDate(dateString) {
	const date = new Date(dateString);
	return date.toLocaleDateString('en-US', {
		year: 'numeric',
		month: 'short',
		day: 'numeric'
	});
}

export function preventNonNumericInput(e) {
	if (!/[0-9]/.test(e.key)) {
		e.preventDefault();
	}
}

export function extractYouTubeId(url) {
	const regex = /(?:youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/embed\/)([^&\n?#]+)/;
	const match = url.match(regex);
	return match ? match[1] : null;
}
