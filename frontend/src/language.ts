export enum SupportedLanguages {
	ENGLISH = 'ENGLISH',
	RUSSIAN = 'RUSSIAN',
}

export function getCurrentLanguage() {
	if (navigator.languages.some(l => l.startsWith('ru')))
		return SupportedLanguages.RUSSIAN;
	return SupportedLanguages.ENGLISH;
}