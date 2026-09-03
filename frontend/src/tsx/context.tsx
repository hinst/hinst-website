import { createContext } from 'react';
import { settingsStorage } from 'src/typescript/settings';

export const AppContext = createContext({
	currentLanguage: settingsStorage.resolvedLanguage,
	windowWidth: window.innerWidth,
	isAdminMode: false,
	setPageTitle: (title: string) => {}
});
