export interface SmartPostImage {
	url?: string;
	dataUrl: string;
}

export interface GoalPostObject {
	goalId: number;
	/** Unix epoch timestamp seconds */
	dateTime: number;
	/** HTML */
	text: string;
	isAutoTranslated: boolean;
	isTranslationPending: boolean;
	languageName?: string;
	languageTag: string;
	isPublic: boolean;
	searchIndexingEnabled: boolean;
	imageCount: number;
}
