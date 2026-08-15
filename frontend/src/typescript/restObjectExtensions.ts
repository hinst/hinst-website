import { DateTime } from 'luxon';
import { GoalObject, GoalPostHeader } from 'src/typescript/generated/rest_objects';
import { SupportedLanguage } from 'src/typescript/language';

/**
 * Attaches frontend-only methods to the generated REST objects.
 *
 * The generated interfaces are extended via TypeScript module augmentation
 * (interface merging), so the generated file itself stays untouched.
 * The methods are plain interface members without runtime code, so instances
 * deserialized from JSON must be hydrated with the `...WithMethods` helpers
 * below (done in `apiClient` and in components that fetch directly).
 */
declare module 'src/typescript/generated/rest_objects' {
	interface GoalObject {
		/** Returns the title in the given language. */
		getTitle(language: SupportedLanguage): string;
	}
	interface GoalPostHeader {
		/** "yyyy-MM" */
		yearAndMonthText: string;
		/** "yyyy-MM-dd" */
		dateText: string;
	}
}

export function goalObjectWithMethods(object: GoalObject): GoalObject {
	object.getTitle = (language: SupportedLanguage) => {
		switch (language) {
			case SupportedLanguage.RUSSIAN:
				return object.title;
			case SupportedLanguage.GERMAN:
				return object.titleGerman;
			case SupportedLanguage.ENGLISH:
				return object.titleEnglish;
		}
	};
	return object;
}

export function goalPostHeaderWithMethods<T extends GoalPostHeader>(object: T): T {
	object.yearAndMonthText = DateTime.fromMillis(object.dateTime * 1000).toFormat('yyyy-MM');
	object.dateText = DateTime.fromMillis(object.dateTime * 1000).toFormat('yyyy-MM-dd');
	return object;
}
