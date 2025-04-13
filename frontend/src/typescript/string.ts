export function compareStrings(a: string, b: string) {
	return a === b ? 0 : a < b ? -1 : 1;
}

export function getHashFromString(s: string) {
	let h = 0, l = s.length, i = 0;
	if (l > 0)
		while (i < l)
			h = (h << 5) - h + s.charCodeAt(i++) | 0;
	return h;
}