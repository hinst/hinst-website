import { createContext } from 'react';
import { getCurrentLanguage } from './language';

export const LanguageContext = createContext(getCurrentLanguage());
export const DisplayWidthContext = createContext(window.innerWidth);