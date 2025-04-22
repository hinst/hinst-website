import { SupportedLanguage as SupportedLanguage } from 'src/typescript/language';
//@ts-ignore
import codingWeekly from 'images/codingWeekly.jpg';
//@ts-ignore
import bicycle from 'images/bicycle.jpg';

export interface GoalInfo {
	englishTitle: string;
	coverImage: any;
}

export const GOAL_INFOS = new Map<string, GoalInfo>([
	[
		'Кодить каждую неделю 8 часов',
		{
			englishTitle: 'Weekly Coding',
			coverImage: codingWeekly
		}
	],
	[
		'Окупить стоимость велосипеда и самоката',
		{
			englishTitle: 'Bicycle and E-Scooter',
			coverImage: bicycle
		}
	]
]);

export function translateGoalTitle(language: SupportedLanguage, text: string) {
	if ([SupportedLanguage.ENGLISH, SupportedLanguage.GERMAN].includes(language)) {
		const info = GOAL_INFOS.get(text);
		if (info) return info.englishTitle;
	}
	return text;
}
