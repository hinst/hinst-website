import { DateTime } from 'luxon';
import { GoalObject, GoalPostHeader } from 'src/typescript/generated/rest_objects';
import { SupportedLanguage } from 'src/typescript/language';

export class GoalObjectWithMethods implements GoalObject {
	id!: number;
	title!: string;
	titleEnglish!: string;
	titleGerman!: string;

	constructor(data: GoalObject) {
		Object.assign(this, data);
	}

	/** Returns the title in the given language. */
	getTitle(language: SupportedLanguage): string {
		switch (language) {
			case SupportedLanguage.RUSSIAN:
				return this.title;
			case SupportedLanguage.GERMAN:
				return this.titleGerman;
			case SupportedLanguage.ENGLISH:
				return this.titleEnglish;
		}
	}
}

export class GoalPostHeaderWithMethods implements GoalPostHeader {
	goalId!: number;
	dateTime!: number;
	isPublic!: boolean;
	type!: string;
	title!: string;
	/** "yyyy-MM" */
	yearAndMonthText: string;
	/** "yyyy-MM-dd" */
	dateText: string;

	constructor(data: GoalPostHeader) {
		Object.assign(this, data);
		this.yearAndMonthText = DateTime.fromMillis(data.dateTime * 1000).toFormat('yyyy-MM');
		this.dateText = DateTime.fromMillis(data.dateTime * 1000).toFormat('yyyy-MM-dd');
	}
}
