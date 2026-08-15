import { DateTime } from 'luxon';
import { GoalPostHeader } from 'src/typescript/generated/rest_objects';

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
