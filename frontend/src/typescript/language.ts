export enum SupportedLanguages {
	ENGLISH = 'ENGLISH',
	RUSSIAN = 'RUSSIAN'
}

export function getCurrentLanguage() {
	for (const lang of navigator.languages) {
		if (lang.startsWith('ru')) return SupportedLanguages.RUSSIAN;
		if (lang.startsWith('en')) return SupportedLanguages.ENGLISH;
	}
	return SupportedLanguages.ENGLISH;
}
