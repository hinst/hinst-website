import { HashRouter, Route, Routes } from 'react-router';
import Header from './header';
import { useEffect, useState } from 'react';
import HomePage from './homePage';
import GoalBrowser from './personal-goals/goalBrowser';
import { getCurrentLanguage, SupportedLanguages } from 'src/typescript/language';
import { AppContext } from './context';
import Cookies from 'js-cookie';

const PAGE_TITLE = 'Showcase Website';

export default function App() {
	useEffect(() => {
		if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches)
			document.getElementsByTagName('html')[0].setAttribute('data-theme', 'dark');
	}, []);

	const [pageTitle, setPageTitle] = useState(PAGE_TITLE);
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
					</Routes>
				</HashRouter>
			</div>
		</AppContext.Provider>
	);
}
