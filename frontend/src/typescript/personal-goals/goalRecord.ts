import { SupportedLanguage } from 'src/typescript/language';

export class GoalHeader {
	constructor(
		public id: string,
		public title: string,
		public titleEnglish: string,
		public titleGerman: string
	) {}

	getTitle(language: SupportedLanguage) {
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
