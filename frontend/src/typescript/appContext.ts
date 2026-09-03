import { createContext } from 'react';
import { PageTitle } from 'src/typescript/pageTitle';
import { settingsStorage } from 'src/typescript/settings';

export const AppContext = createContext({
	currentLanguage: settingsStorage.resolvedLanguage,
	windowWidth: window.innerWidth,
	isAdminMode: false,
	setPageTitle: (_title: PageTitle) => {},
});
