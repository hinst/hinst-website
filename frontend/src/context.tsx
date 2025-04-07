import { createContext } from 'react';
import { getCurrentLanguage } from 'src/typescript/language';

export const LanguageContext = createContext(getCurrentLanguage());
export const DisplayWidthContext = createContext(window.innerWidth);