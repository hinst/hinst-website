import { DateTime } from 'luxon';

export class GoalPostObject {
	constructor(
		public goalId: number = 0,
		/** Unix epoch seconds */
		public dateTime: number = 0,
		/** HTML */
		public text: string = '',
		public isAutoTranslated: boolean = false,
		public isTranslationPending: boolean = false,
		public languageName: string = '',
		public languageTag: string = '',
		public isPublic: boolean = false,
		public searchIndexingEnabled?: boolean,
		public imageCount: number = 0,
	) {}
}

export class GoalPostHeader {
	constructor(
		public goalId: number = 0,
		/** Unix epoch time seconds */
		public dateTime: number = 0,
		public isPublic: boolean = false,
		/** "post" or "comment" */
		public type: string = '',
		public title: string = '',
		public googlePingedAt?: number,
		public publicUrl?: string,
	) {}

	get yearAndMonthText(): string {
		return DateTime.fromMillis(this.dateTime * 1000).toFormat('yyyy-MM');
	}

	get dateText(): string {
		return DateTime.fromMillis(this.dateTime * 1000).toFormat('yyyy-MM-dd');
	}
}
