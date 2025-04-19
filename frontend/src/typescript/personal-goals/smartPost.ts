export interface SmartPostImage {
	url?: string;
	dataUrl: string;
}

export interface SmartPost {
	/** Can be: 'post' */
	type: string;
	id: string;
	/** Goal id */
	obj_id: string;
	msg: string;
	/** Example: 2023-04-28 09:12:21 */
	date: string;
	comments: Comment[];
	images: SmartPostImage[];
	count_comments: string;
	username: string;
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
}

export interface GoalPostObjectExtended extends GoalPostObject {
	images: string[];
}
