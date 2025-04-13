import { createContext } from 'react';
import { getCurrentLanguage } from 'src/typescript/language';

export const AppContext = createContext({
	currentLanguage: getCurrentLanguage(),
	displayWidth: window.innerWidth
});
