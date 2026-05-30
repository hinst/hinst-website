import { SupportedLanguage as SupportedLanguage } from 'src/typescript/language';

export interface GoalInfo {
	englishTitle: string;
}

export const GOAL_INFOS = new Map<string, GoalInfo>([
	[
		'Кодить каждую неделю 8 часов',
		{
			englishTitle: 'Weekly Coding'
		}
	],
	[
		'Экономные поездки на двух колёсах',
		{
			englishTitle: 'Savvy trips on two wheels'
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
