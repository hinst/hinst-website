import { HashRouter, Route, Routes } from 'react-router';
import Header from './header';
import { useEffect, useState } from 'react';
import HomePage from './homePage';
import GoalBrowser from './personal-goals/goalBrowser';
import { getCurrentLanguage, SupportedLanguages } from 'src/typescript/language';
import { AppContext } from './context';
import Cookies from 'js-cookie';
import SettingsPage from './settings/settingsPage';
import { APP_TITLE } from 'src/typescript/global';
import { settingsStorage } from 'src/typescript/settings';

export default function App() {
	settingsStorage.initialize();

	const [pageTitle, setPageTitle] = useState(APP_TITLE);
	useEffect(() => {
		document.title = pageTitle;
	}, [pageTitle]);

	const [currentLanguage, setCurrentLanguage] =
		useState<SupportedLanguages>(getCurrentLanguage());
	useEffect(() => {
		const timer = setInterval(() => {
			const newLanguage = getCurrentLanguage();
			if (newLanguage !== currentLanguage) setCurrentLanguage(newLanguage);
		}, 1000);
		return () => clearInterval(timer);
	}, []);

	return (
		<AppContext.Provider
			value={{
				currentLanguage,
				displayWidth: window.innerWidth,
				goalManagerMode: Cookies.get('goalManagerMode') == '1'
			}}
		>
			<div
				style={{
					padding: 10,
					paddingBottom: 0,
					display: 'flex',
					flexDirection: 'column',
					width: '100%',
					maxWidth: '100%',
					maxHeight: '100%'
				}}
			>
				<HashRouter>
					<div style={{ marginBottom: 10 }}>
						<Header title={pageTitle} />
					</div>
					<Routes>
						<Route path='/' element={<HomePage setPageTitle={setPageTitle} />} />
						<Route
							path='/personal-goals/:id'
							element={<GoalBrowser setPageTitle={setPageTitle} />}
						/>
						<Route
							path='/settings'
							element={<SettingsPage setPageTitle={setPageTitle} />}
						/>
					</Routes>
				</HashRouter>
			</div>
		</AppContext.Provider>
	);
}
