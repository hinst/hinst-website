import type { GoalObject } from 'src/typescript/apiTypes';
import { SupportedLanguage } from 'src/typescript/language';

export interface GoalObjectEx extends GoalObject {
	/** Returns the title in the given language. */
	getTitle(language: SupportedLanguage): string;
}

export function goalObjectWithMethods(data: GoalObject): GoalObjectEx {
	return {
		...data,
		getTitle(language: SupportedLanguage) {
			switch (language) {
				case SupportedLanguage.RUSSIAN:
					return data.title;
				case SupportedLanguage.GERMAN:
					return data.titleGerman;
				case SupportedLanguage.ENGLISH:
					return data.titleEnglish;
			}
		}
	};
}
