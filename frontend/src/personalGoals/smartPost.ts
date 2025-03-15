export interface SmartPost {
    /** Can be: 'post' */
    type: string;
    id: string;
    msg: string;
    /** Example: 2023-04-28 09:12:21 */
    date: string;
    comments: Comment[];
    images: {
        url: string;
        dataUrl: string;
    }[];
    count_comments: string;
    username: string;
}

export interface SmartPostExtended extends SmartPost {
    isAutoTranslated: boolean;
    languageName?: string;
}
