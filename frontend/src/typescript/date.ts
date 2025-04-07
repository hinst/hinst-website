export function parseMonthlyDate(text: string) {
	const parts = text.split('-');
	return new Date(
		parseInt(parts[0]),
		parseInt(parts[1]) - 1,
		1
	);
}

export function getMonthName(date: Date) {
	return date.toLocaleString('en', { month: 'long' });
}