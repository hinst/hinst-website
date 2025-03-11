import { getCurrentLanguage, SupportedLanguages } from '../language';
//@ts-ignore
import codingWeekly from '../../images/codingWeekly.jpg';
//@ts-ignore
import bicycle from '../../images/bicycle.jpg';

export interface GoalInfo {
	englishTitle: string;
	coverImage: any;
}

export const GOAL_INFOS = new Map<string, GoalInfo>([
	['Кодить каждую неделю 8 часов', {
		englishTitle: 'Weekly Coding',
		coverImage: codingWeekly,
	}],
	['Окупить стоимость велосипеда и самоката', {
		englishTitle: 'My Bicycle and E-Scooter',
		coverImage: bicycle,
	}],
]);

export function translateGoalTitle(text: string) {
	if (getCurrentLanguage() === SupportedLanguages.ENGLISH) {
		const info = GOAL_INFOS.get(text);
		if (info)
			return info.englishTitle;
	}
	return text;
}
