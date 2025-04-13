import { createContext } from 'react';
import { getCurrentLanguage } from 'src/typescript/language';
import Cookies from 'js-cookie';

export const AppContext = createContext({
	currentLanguage: getCurrentLanguage(),
	displayWidth: window.innerWidth,
	goalManagerMode: Cookies.get('goalManagerMode') == '1'
});
