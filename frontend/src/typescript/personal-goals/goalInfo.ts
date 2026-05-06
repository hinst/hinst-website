import { SupportedLanguage as SupportedLanguage } from 'src/typescript/language';
//@ts-ignore
import codingWeekly from 'url:images/codingWeekly.jpg';
//@ts-ignore
import eScooter from 'url:images/Navee.jpg';

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
		'Экономные поездки на двух колёсах',
		{
			englishTitle: 'Savvy trips on two wheels',
			coverImage: eScooter
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
