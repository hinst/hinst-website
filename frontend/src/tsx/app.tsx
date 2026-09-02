import { HashRouter, Route, Routes } from 'react-router';
import Header from './header';
import { useEffect, useState } from 'react';
import HomePage from './homePage';
import GoalBrowser from './personal-goals/goalBrowser';
import { SupportedLanguage } from 'src/typescript/language';
import { AppContext } from './context';
import SettingsPage from './settings/settingsPage';
import { APP_TITLE } from 'src/typescript/global';
import { settingsStorage } from 'src/typescript/settings';
import ManualPingTracker from './manual-ping-tracker/manualPingTracker';
import { PersonalGoalsSearch } from './personal-goals-search/personalGoalsSearch';
import { apiClient } from 'src/typescript/apiClient';

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
		loadAdminMode();
	}, []);

	async function loadAdminMode() {
		setAdminMode(await apiClient.isAdminModeEnabled());
	}

	return (
		<AppContext.Provider
			value={{
				currentLanguage,
				windowWidth: windowWidth,
				isAdminMode: isAdminMode
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
						<Route
							path='/manual-ping-tracker'
							element={<ManualPingTracker setPageTitle={setPageTitle} />}
						/>
						<Route path='/personal-goals-search' element={<PersonalGoalsSearch />} />
					</Routes>
				</HashRouter>
			</div>
		</AppContext.Provider>
	);
}
