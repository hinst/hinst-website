import { HashRouter, Route, Routes } from 'react-router';
import Header from './header';
import { useEffect, useState } from 'react';
import HomePage from './homePage';
import GoalBrowser from './personal-goals/goalBrowser';
import { SupportedLanguage } from 'src/typescript/language';
import { AppContext } from './context';
import Cookies from 'js-cookie';
import SettingsPage from './settings/settingsPage';
import { APP_TITLE } from 'src/typescript/global';
import { settingsStorage } from 'src/typescript/settings';
import { debounce } from 'lodash';

export default function App() {
	settingsStorage.initialize();

	const [currentLanguage, setCurrentLanguage] = useState<SupportedLanguage>(
		settingsStorage.resolvedLanguage
	);
	useEffect(() => {
		const timer = setInterval(() => {
			const newLanguage = settingsStorage.resolvedLanguage;
			if (newLanguage !== currentLanguage) setCurrentLanguage(newLanguage);
		}, 500);
		return () => clearInterval(timer);
	}, []);

	const [windowWidth, setWindowWidth] = useState(window.innerWidth);
	useEffect(() => {
		const timer = setInterval(() => {
			const newWindowWidth = window.innerWidth;
			if (windowWidth != newWindowWidth) setWindowWidth(newWindowWidth);
		}, 500);
		return () => clearInterval(timer);
	}, []);

	const [pageTitle, setPageTitle] = useState(APP_TITLE);
	useEffect(() => {
		document.title = pageTitle;
	}, [pageTitle]);

	return (
		<AppContext.Provider
			value={{
				currentLanguage,
				windowWidth: windowWidth,
				goalManagerMode: Cookies.get('goalManagerMode') === '1'
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
