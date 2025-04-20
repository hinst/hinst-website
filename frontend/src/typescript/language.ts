export enum SupportedLanguage {
	ENGLISH = 'en',
	RUSSIAN = 'ru',
	GERMAN = 'de',
}

export function getSystemLanguage(): SupportedLanguage | undefined {
	for (const lang of navigator.languages) {
		if (lang.startsWith('ru')) return SupportedLanguage.RUSSIAN;
		if (lang.startsWith('en')) return SupportedLanguage.ENGLISH;
		if (lang.startsWith('de')) return SupportedLanguage.GERMAN;
	}
	return undefined;;
}
