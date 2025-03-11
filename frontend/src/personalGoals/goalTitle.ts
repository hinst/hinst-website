import { getCurrentLanguage, SupportedLanguages } from '../language';

export function translateGoalTitle(text: string) {
	if (getCurrentLanguage() === SupportedLanguages.ENGLISH) {
		switch (text) {
			case 'Кодить каждую неделю 8 часов':
				return 'Weekly Coding';
			case 'Окупить стоимость велосипеда и самоката':
				return 'My Bicycle and E-Scooter';
		}
	}
	return text;
}
