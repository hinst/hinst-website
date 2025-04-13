export function compareStrings(a: string, b: string) {
	return a === b ? 0 : a < b ? -1 : 1;
}
