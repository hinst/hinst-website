import { useEffect, useState } from 'react';
import { HashRouter, Route, Routes } from 'react-router';
import { apiClient } from 'src/typescript/apiClient';
import { APP_TITLE } from 'src/typescript/global';
import type { SupportedLanguage } from 'src/typescript/language';
import { settingsStorage } from 'src/typescript/settings';
import { AppContext } from './appContext';
import Header from './header';
import HomePage from './homePage';
import ManualPingTracker from './manual-ping-tracker/manualPingTracker';
import GoalBrowser from './personal-goals/goalBrowser';
import { PersonalGoalsSearch } from './personal-goals-search/personalGoalsSearch';
import SettingsPage from './settings/settingsPage';

export default function App() {
	settingsStorage.initialize();

	const [currentLanguage, setCurrentLanguage] = useState<SupportedLanguage>(
		settingsStorage.resolvedLanguage
	);
	useEffect(() => {
		const timer = setInterval(() => setCurrentLanguage(settingsStorage.resolvedLanguage), 500);
		return () => clearInterval(timer);
	}, []);

	const [windowWidth, setWindowWidth] = useState(window.innerWidth);
	useEffect(() => {
		const timer = setInterval(() => setWindowWidth(window.innerWidth), 500);
		return () => clearInterval(timer);
	}, []);

	const [pageTitle, setPageTitle] = useState(APP_TITLE);
	useEffect(() => {
		document.title = pageTitle;
	}, [pageTitle]);

	const [isAdminMode, setAdminMode] = useState(false);
	useEffect(() => {
		async function loadAdminMode() {
			setAdminMode(await apiClient.isAdminModeEnabled());
		}
		const _promise = loadAdminMode();
	}, []);

	return (
		<AppContext.Provider
			value={{
				currentLanguage,
				windowWidth: windowWidth,
				isAdminMode: isAdminMode,
				setPageTitle: setPageTitle
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
						<Route path='/' element={<HomePage />} />
						<Route path='/personal-goals/:id' element={<GoalBrowser />} />
						<Route path='/settings' element={<SettingsPage />} />
						<Route path='/manual-ping-tracker' element={<ManualPingTracker />} />
						<Route path='/personal-goals-search' element={<PersonalGoalsSearch />} />
					</Routes>
				</HashRouter>
			</div>
		</AppContext.Provider>
	);
}
