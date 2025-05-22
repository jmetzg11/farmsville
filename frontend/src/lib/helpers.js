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
