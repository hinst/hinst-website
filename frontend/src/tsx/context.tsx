import { createContext } from 'react';
import Cookies from 'js-cookie';
import { settingsStorage } from 'src/typescript/settings';

export const AppContext = createContext({
	currentLanguage: settingsStorage.resolvedLanguage,
	windowWidth: window.innerWidth,
	goalManagerMode: Cookies.get('goalManagerMode') === '1'
});
