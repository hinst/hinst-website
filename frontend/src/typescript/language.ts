export enum SupportedLanguage {
	ENGLISH = 'en',
	RUSSIAN = 'ru',
	GERMAN = 'de',
}

export const supportedLanguageNames = {
	[SupportedLanguage.ENGLISH]: 'English',
	[SupportedLanguage.RUSSIAN]: 'Russian',
	[SupportedLanguage.GERMAN]: 'German',
}

export function getSystemLanguage(): SupportedLanguage | undefined {
	for (const lang of navigator.languages) {
		if (lang.startsWith('ru')) return SupportedLanguage.RUSSIAN;
		if (lang.startsWith('en')) return SupportedLanguage.ENGLISH;
		if (lang.startsWith('de')) return SupportedLanguage.GERMAN;
	}
	return undefined;;
}
