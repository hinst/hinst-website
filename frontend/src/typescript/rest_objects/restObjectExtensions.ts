import { DateTime } from 'luxon';
import { GoalObject, GoalPostHeader } from 'src/typescript/generated/rest_objects';
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

export interface GoalPostHeaderEx extends GoalPostHeader {
	/** "yyyy-MM" */
	yearAndMonthText: string;
	/** "yyyy-MM-dd" */
	dateText: string;
}

export function goalPostHeaderWithMethods(data: GoalPostHeader): GoalPostHeaderEx {
	return {
		...data,
		yearAndMonthText: DateTime.fromMillis(data.dateTime * 1000).toFormat('yyyy-MM'),
		dateText: DateTime.fromMillis(data.dateTime * 1000).toFormat('yyyy-MM-dd')
	};
}
