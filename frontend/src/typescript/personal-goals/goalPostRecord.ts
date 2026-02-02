import { DateTime } from 'luxon';

export class GoalPostRecord {
	constructor(
		public goalId: number = 0,
		/** Unix timestamp seconds */
		public dateTime: number = 0,
		public isPublic: boolean = false,
		public type: string = '',
		public title: string = '',
	) {}

	get yearAndMonthText(): string {
		return DateTime.fromMillis(this.dateTime * 1000).toFormat('yyyy-MM');
	}
}
